// GC basics: GOGC changes how often the GC runs, runtime.GC frees garbage, MemStats shows it.
// Run: go run ./03-runtime-internals/05-garbage-collector
package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

var (
	live [][]byte // data that stays reachable
	sink []byte   // forces each allocation onto the heap
)

func heapMB() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.HeapAlloc >> 20
}

// gcCycles allocates 200 MiB of garbage with the given GOGC and counts GC cycles.
func gcCycles(gogc int) uint32 {
	old := debug.SetGCPercent(gogc)
	defer debug.SetGCPercent(old)
	runtime.GC() // start from a clean heap
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	for range 200 {
		sink = make([]byte, 1<<20)
	}
	runtime.ReadMemStats(&after)
	return after.NumGC - before.NumGC
}

func main() {
	for range 20 {
		live = append(live, make([]byte, 1<<20)) // 20 MiB of live data
	}

	fmt.Println("=== GOGC: higher means fewer GC cycles (more memory) ===")
	fmt.Println("GOGC=100:", gcCycles(100), "cycles")
	fmt.Println("GOGC=800:", gcCycles(800), "cycles")

	// GOMEMLIMIT is a SOFT limit: the GC works harder near it, it is not a hard cap.
	debug.SetMemoryLimit(512 << 20)
	fmt.Println("soft memory limit set to 512 MiB")

	fmt.Println("=== drop live data, then runtime.GC() ===")
	fmt.Println("heap with data alive (MB):", heapMB())
	live, sink = nil, nil
	runtime.GC()
	fmt.Println("heap after runtime.GC() (MB):", heapMB())
	fmt.Println("tip: GODEBUG=gctrace=1 prints one line per GC cycle")
}
