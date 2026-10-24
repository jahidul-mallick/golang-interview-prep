// sync.Cond, sync.Pool and sync.Map.
// Run: go run ./02-concurrency/06-sync-cond-pool-map
package main

import (
	"bytes"
	"fmt"
	"sync"
)

func main() {
	fmt.Println("=== sync.Cond: wait for a condition, wake all with Broadcast ===")
	var mu sync.Mutex
	cond := sync.NewCond(&mu)
	ready := false
	var wg sync.WaitGroup
	for i := range 3 {
		wg.Go(func() {
			mu.Lock()
			for !ready { // always a loop: re-check after every wake-up
				cond.Wait() // unlocks mu while waiting, re-locks before returning
			}
			mu.Unlock()
			fmt.Println("waiter", i, "released")
		})
	}
	mu.Lock()
	ready = true
	mu.Unlock()
	cond.Broadcast()
	wg.Wait()

	fmt.Println("\n=== sync.Pool: reuse temporary buffers ===")
	pool := sync.Pool{New: func() any { return new(bytes.Buffer) }}
	for i := range 3 {
		buf := pool.Get().(*bytes.Buffer)
		buf.Reset() // a pooled object may hold old data
		fmt.Fprintf(buf, "request-%d", i)
		fmt.Println(buf.String())
		pool.Put(buf) // the pool may drop it at any GC: never rely on it
	}

	fmt.Println("\n=== sync.Map: no locking code, untyped ===")
	var m sync.Map
	m.Store("a", 1)
	actual, loaded := m.LoadOrStore("a", 99) // key exists: returns the old value
	fmt.Println("LoadOrStore:", actual, loaded)
	if v, ok := m.Load("a"); ok {
		fmt.Println("Load:", v.(int)+1) // values are any: type assertion needed
	}
	m.Range(func(k, v any) bool {
		fmt.Println("Range:", k, v)
		return true
	})
	m.Delete("a")
}
