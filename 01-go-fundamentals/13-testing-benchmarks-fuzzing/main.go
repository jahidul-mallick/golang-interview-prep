// Code under test for topic 13. The real tests and benchmark are snippets in
// the README; here we run a table check and measure allocations by hand.
// Run: go run ./01-go-fundamentals/13-testing-benchmarks-fuzzing
package main

import (
	"fmt"
	"testing"
)

var sink string // keeps the compiler from optimizing the call away

// Reverse reverses s rune by rune.
func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func main() {
	fmt.Println("=== table-driven check (same shape as a Test function) ===")
	tests := []struct {
		name, in, want string
	}{
		{"empty", "", ""},
		{"ascii", "abc", "cba"},
		{"unicode", "héllo", "olléh"},
	}
	for _, tt := range tests {
		got := Reverse(tt.in)
		status := "PASS"
		if got != tt.want {
			status = "FAIL"
		}
		fmt.Printf("%-8s Reverse(%q) = %q  %s\n", tt.name, tt.in, got, status)
	}

	fmt.Println("=== allocations per call (what -benchmem reports as allocs/op) ===")
	allocs := testing.AllocsPerRun(100, func() { sink = Reverse("hello, world") })
	fmt.Println("allocs per Reverse call:", allocs)
	fmt.Println("full benchmark: go test -bench . -benchmem")
}
