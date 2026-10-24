// LLD rule engine: rules loaded from JSON, a Condition tree, priorities,
// validate at load time, atomic hot swap. Run: go run ./06-lld-design-patterns/07-lld-rule-engine
package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync/atomic"
)

type Facts map[string]any

// Condition is the composite: leaves compare, And/Or/Not combine.
type Condition interface{ Eval(Facts) bool }

type Compare struct {
	Field string
	Op    string // eq, gt, in
	Value any
}

func (c Compare) Eval(f Facts) bool {
	got, ok := f[c.Field]
	if !ok {
		return false // missing fact means no match, never a panic
	}
	switch c.Op {
	case "eq":
		return fmt.Sprint(got) == fmt.Sprint(c.Value)
	case "gt":
		a, aok := got.(float64)
		b, bok := c.Value.(float64)
		return aok && bok && a > b
	case "in":
		list, _ := c.Value.([]any)
		for _, v := range list {
			if fmt.Sprint(v) == fmt.Sprint(got) {
				return true
			}
		}
	}
	return false
}

type And []Condition
type Or []Condition
type Not struct{ C Condition }

func (a And) Eval(f Facts) bool {
	for _, c := range a {
		if !c.Eval(f) {
			return false
		}
	}
	return true
}

func (o Or) Eval(f Facts) bool {
	for _, c := range o {
		if c.Eval(f) {
			return true
		}
	}
	return false
}

func (n Not) Eval(f Facts) bool { return !n.C.Eval(f) }

type Rule struct {
	ID       string
	Priority int
	When     Condition
	Action   string
}

// JSON shape. Conditions are polymorphic, so decode to a raw node first.
type rawCond struct {
	All   []rawCond `json:"all"`
	Any   []rawCond `json:"any"`
	Not   *rawCond  `json:"not"`
	Field string    `json:"field"`
	Op    string    `json:"op"`
	Value any       `json:"value"`
}

func (r rawCond) build() (Condition, error) {
	kids := func(rs []rawCond) ([]Condition, error) {
		out := make([]Condition, 0, len(rs))
		for _, k := range rs {
			c, err := k.build()
			if err != nil {
				return nil, err
			}
			out = append(out, c)
		}
		return out, nil
	}
	switch {
	case len(r.All) > 0:
		k, err := kids(r.All)
		return And(k), err
	case len(r.Any) > 0:
		k, err := kids(r.Any)
		return Or(k), err
	case r.Not != nil:
		k, err := r.Not.build()
		return Not{k}, err
	case r.Op == "eq" || r.Op == "gt" || r.Op == "in":
		return Compare{r.Field, r.Op, r.Value}, nil
	}
	return nil, fmt.Errorf("bad condition (op %q)", r.Op)
}

// LoadRules validates everything up front so Evaluate can trust the rules.
func LoadRules(src string) ([]Rule, error) {
	var raws []struct {
		ID       string  `json:"id"`
		Priority int     `json:"priority"`
		When     rawCond `json:"when"`
		Action   string  `json:"action"`
	}
	if err := json.Unmarshal([]byte(src), &raws); err != nil {
		return nil, err
	}
	var rules []Rule
	for _, r := range raws {
		c, err := r.When.build()
		if err != nil {
			return nil, fmt.Errorf("rule %s: %w", r.ID, err)
		}
		rules = append(rules, Rule{r.ID, r.Priority, c, r.Action})
	}
	sort.SliceStable(rules, func(i, j int) bool { return rules[i].Priority > rules[j].Priority })
	return rules, nil
}

// Engine keeps an immutable rule slice behind an atomic pointer: readers
// take no lock, and a reload is one pointer store.
type Engine struct{ rules atomic.Pointer[[]Rule] }

func (e *Engine) Swap(r []Rule) { e.rules.Store(&r) }

func (e *Engine) Evaluate(f Facts, firstOnly bool) []string {
	var out []string
	for _, r := range *e.rules.Load() {
		if r.When.Eval(f) {
			out = append(out, r.ID+"->"+r.Action)
			if firstOnly {
				break
			}
		}
	}
	return out
}

const policy = `[
 {"id":"open-db","priority":100,"action":"alert-critical","when":{"all":[
   {"field":"open_to_world","op":"eq","value":true},
   {"field":"type","op":"in","value":["rds","aurora"]}]}},
 {"id":"open-any","priority":50,"action":"alert-warning","when":{"all":[
   {"field":"open_to_world","op":"eq","value":true},
   {"not":{"field":"role","op":"eq","value":"public-web"}}]}},
 {"id":"big","priority":10,"action":"tag-review","when":{"field":"vcpus","op":"gt","value":32}}
]`

func main() {
	rules, err := LoadRules(policy)
	if err != nil {
		panic(err)
	}
	var eng Engine
	eng.Swap(rules)
	db := Facts{"type": "rds", "open_to_world": true, "role": "db", "vcpus": 64.0}
	web := Facts{"type": "ec2", "open_to_world": true, "role": "public-web", "vcpus": 4.0}

	fmt.Println("=== all-match (audit) ===")
	fmt.Println("  db :", eng.Evaluate(db, false))
	fmt.Println("  web:", eng.Evaluate(web, false))

	fmt.Println("=== first-match (highest priority wins) ===")
	fmt.Println("  db :", eng.Evaluate(db, true))

	fmt.Println("=== hot swap ===")
	r2, _ := LoadRules(`[{"id":"quarantine","priority":1,"action":"isolate",
	  "when":{"field":"open_to_world","op":"eq","value":true}}]`)
	eng.Swap(r2)
	fmt.Println("  web:", eng.Evaluate(web, false))

	fmt.Println("=== invalid rule rejected at load ===")
	_, err = LoadRules(`[{"id":"bad","when":{"field":"x","op":"regex","value":"a"}}]`)
	fmt.Println("  error:", err)
}
