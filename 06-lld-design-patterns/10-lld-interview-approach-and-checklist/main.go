// Worked LLD example: Splitwise-style settlement (net balances, greedy match).
// Run: go run ./06-lld-design-patterns/10-lld-interview-approach-and-checklist
package main

import (
	"fmt"
	"sort"
)

type Transfer struct {
	From, To string
	Amount   int
}

// Settle takes net balances (positive = is owed, negative = owes) and returns
// few transfers by always matching the biggest debtor with the biggest creditor.
func Settle(balances map[string]int) []Transfer {
	type p struct {
		name string
		amt  int
	}
	var cred, debt []p
	for n, b := range balances {
		if b > 0 {
			cred = append(cred, p{n, b})
		} else if b < 0 {
			debt = append(debt, p{n, -b})
		}
	}
	byAmt := func(s []p) {
		sort.Slice(s, func(i, j int) bool {
			if s[i].amt != s[j].amt {
				return s[i].amt > s[j].amt
			}
			return s[i].name < s[j].name // stable output despite map order
		})
	}
	var out []Transfer
	for len(cred) > 0 && len(debt) > 0 {
		byAmt(cred)
		byAmt(debt)
		amt := min(cred[0].amt, debt[0].amt)
		out = append(out, Transfer{debt[0].name, cred[0].name, amt})
		cred[0].amt -= amt
		debt[0].amt -= amt
		if cred[0].amt == 0 {
			cred = cred[1:]
		}
		if debt[0].amt == 0 {
			debt = debt[1:]
		}
	}
	return out
}

func main() {
	fmt.Println("=== Splitwise settlement ===")
	// Alice paid 90 for three people; Bob paid 30.
	balances := map[string]int{"alice": 50, "bob": -10, "carol": -40}
	for _, t := range Settle(balances) {
		fmt.Printf("  %s pays %s %d\n", t.From, t.To, t.Amount)
	}
}
