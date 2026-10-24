// Tricky "what does this print?" snippets. Each one prints expected vs actual.
// Run: go run ./01-go-fundamentals/16-gotchas-and-tricky-questions
package main

import (
	"fmt"
	"unicode/utf8"
)

type T struct{}

func check(name, expected, actual string) {
	status := "OK"
	if expected != actual {
		status = "DIFF"
	}
	fmt.Printf("=== %s ===\n  expected: %s | actual: %s [%s]\n", name, expected, actual, status)
}

func modify(s []int) {
	s[0] = 9
	s = append(s, 4) // caller's slice header is unchanged
}

func deferArg() (out string) {
	x := 1
	defer func(v int) { out = fmt.Sprint(v) }(x) // argument evaluated now
	x = 2
	return out
}

func deferOrder() (out string) {
	for i := range 3 {
		defer func() { out += fmt.Sprint(i) }()
	}
	return out
}

func panics(f func()) (msg string) {
	defer func() { msg = fmt.Sprint(recover()) }()
	f()
	return ""
}

func main() {
	var p *T
	var i any = p
	check("typed nil inside an interface == nil", "false", fmt.Sprint(i == nil))

	a := []int{1, 2, 3}
	b := a[:2]
	b = append(b, 99) // spare capacity: overwrites a[2]
	check("append shares the backing array", "[1 2 99]", fmt.Sprint(a))

	a2 := []int{1, 2, 3}
	b2 := a2[:2:2] // cap 2 forces append to copy
	b2 = append(b2, 99)
	check("full slice expression s[:2:2]", "[1 2 3] [1 2 99]", fmt.Sprint(a2, b2))

	a3 := []int{1, 2, 3}
	modify(a3)
	check("slice passed to a func: writes shared, append not", "[9 2 3]", fmt.Sprint(a3))

	check("defer evaluates arguments at defer time", "1", deferArg())
	check("defers run last in, first out (Go 1.22 loop var)", "210", deferOrder())

	type item struct{ N int }
	items := []item{{1}, {2}}
	for _, v := range items {
		v.N++ // v is a copy
	}
	check("range value is a copy", "1", fmt.Sprint(items[0].N))

	var funcs []func() int
	for i := range 3 {
		funcs = append(funcs, func() int { return i })
	}
	check("closure captures per-iteration i (Go 1.22+)", "0 1 2", fmt.Sprint(funcs[0](), funcs[1](), funcs[2]()))

	s := "héllo"
	check("len counts bytes, not runes", "6 5", fmt.Sprint(len(s), utf8.RuneCountInString(s)))

	x := 1
	if x := 2; x > 1 {
		_ = x // inner x shadows the outer one
	}
	check("shadowing in if scope", "1", fmt.Sprint(x))

	var nilSlice []int
	check("nil slice vs empty slice", "true false", fmt.Sprint(nilSlice == nil, []int{} == nil))

	var m map[string]int
	check("read from nil map is fine", "0", fmt.Sprint(m["a"]))
	check("write to nil map panics", "assignment to entry in nil map", panics(func() { m["a"] = 1 }))

	ch := make(chan int, 1)
	close(ch)
	v, ok := <-ch
	check("receive from closed channel", "0 false", fmt.Sprint(v, ok))
	check("close of closed channel panics", "close of closed channel", panics(func() { close(ch) }))
}
