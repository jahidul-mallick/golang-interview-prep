// Size classes: the allocator rounds each request up to a fixed size class.
// Run: go run ./03-runtime-internals/04-memory-allocator-mcache-mcentral-mheap
package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

var sink [][]byte

// perObject allocates n objects of the given size and returns heap bytes used per object.
func perObject(size, n int) float64 {
	sink = make([][]byte, 0, n)
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	for range n {
		sink = append(sink, make([]byte, size))
	}
	runtime.ReadMemStats(&after)
	return float64(after.HeapAlloc-before.HeapAlloc) / float64(n)
}

func main() {
	debug.SetGCPercent(-1) // keep the GC out of the measurement

	fmt.Println("=== requested size -> bytes really used per object ===")
	for _, size := range []int{24, 25, 100, 1000} {
		fmt.Printf("request %4d B  ->  about %.0f B\n", size, perObject(size, 10_000))
	}

	fmt.Println("=== large objects (> 32 KiB) skip the size classes ===")
	fmt.Printf("request 40000 B  ->  about %.0f B (whole 8 KiB pages)\n", perObject(40_000, 100))
}
