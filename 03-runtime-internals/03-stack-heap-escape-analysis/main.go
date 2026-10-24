// Which patterns allocate on the heap? Measured with testing.AllocsPerRun.
// See the compiler's reasoning with: go run -gcflags='-m' ./03-runtime-internals/03-stack-heap-escape-analysis
// Run: go run ./03-runtime-internals/03-stack-heap-escape-analysis
package main

import (
	"fmt"
	"testing"
)

type point struct{ X, Y int }

var (
	sinkInt int
	sinkPtr *point
	sinkAny any
)

//go:noinline
func byValue() point { return point{1, 2} } // lives in the caller's frame

//go:noinline
func byPointer() *point { return &point{1, 2} } // &point escapes to heap

//go:noinline
func boxed(i int) any { return i } // converting to interface may allocate

//go:noinline
func fixedSlice() int {
	s := make([]int, 0, 100) // size known at compile time: stack
	s = append(s, 1)
	return len(s)
}

//go:noinline
func dynamicSlice(n int) int {
	s := make([]int, 0, n) // non-constant size: heap
	s = append(s, 1)
	return len(s)
}

func report(name string, f func()) {
	fmt.Printf("%-22s allocs/op = %.0f\n", name, testing.AllocsPerRun(100, f))
}

func main() {
	fmt.Println("=== allocations per call ===")
	report("return struct value", func() { sinkInt = byValue().X })
	report("return &struct", func() { sinkPtr = byPointer() })
	report("any(7) (small int)", func() { sinkAny = boxed(7) })
	report("any(1<<20)", func() { sinkAny = boxed(1 << 20) })
	report("make cap 100 (const)", func() { sinkInt = fixedSlice() })
	report("make cap n (runtime)", func() { sinkInt = dynamicSlice(100) })
	fmt.Println("=== takeaway ===")
	fmt.Println("the compiler decides stack vs heap by whether a value outlives its function")
}
