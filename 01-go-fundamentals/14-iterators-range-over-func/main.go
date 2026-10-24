// Shows range over int, a custom iter.Seq, and slices/maps iterator helpers.
// Run: go run ./01-go-fundamentals/14-iterators-range-over-func
package main

import (
	"fmt"
	"iter"
	"maps"
	"slices"
)

func fibonacci(n int) iter.Seq[int] {
	return func(yield func(int) bool) {
		a, b := 0, 1
		for range n {
			if !yield(a) {
				return // the loop body did break
			}
			a, b = b, a+b
		}
	}
}

func main() {
	fmt.Println("=== range over int ===")
	for i := range 3 {
		fmt.Print(i, " ")
	}
	fmt.Println()

	fmt.Println("=== custom iterator with early break ===")
	for v := range fibonacci(100) {
		if v > 20 {
			break
		}
		fmt.Print(v, " ")
	}
	fmt.Println()

	fmt.Println("=== Seq2 and collecting back to a slice ===")
	for i, w := range slices.All([]string{"go", "is", "fast"}) {
		fmt.Printf("[%d]=%s ", i, w)
	}
	fmt.Println()
	fmt.Println(slices.Collect(fibonacci(6)), slices.Sorted(maps.Keys(map[string]int{"b": 1, "a": 2})))
}
