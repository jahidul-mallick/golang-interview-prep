// sync/atomic: counter, CAS loop, atomic.Pointer config swap.
// Run: go run ./02-concurrency/07-atomic-operations
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Config struct {
	Version int
	Name    string
}

func main() {
	fmt.Println("=== atomic counter (x++ would be a data race) ===")
	var hits atomic.Int64
	var wg sync.WaitGroup
	for range 100 {
		wg.Go(func() {
			for range 1000 {
				hits.Add(1)
			}
		})
	}
	wg.Wait()
	fmt.Println("hits =", hits.Load())

	fmt.Println("\n=== CompareAndSwap loop: lock-free running max ===")
	var maxSeen atomic.Int64
	for g := range 8 {
		wg.Go(func() {
			for i := range 1000 {
				storeMax(&maxSeen, int64(g*1000+i))
			}
		})
	}
	wg.Wait()
	fmt.Println("max =", maxSeen.Load())

	fmt.Println("\n=== CAS on a flag: do something exactly once ===")
	var closed atomic.Bool
	var winners atomic.Int32
	for range 5 {
		wg.Go(func() {
			if closed.CompareAndSwap(false, true) {
				winners.Add(1) // only one goroutine gets here
			}
		})
	}
	wg.Wait()
	fmt.Println("goroutines that won the CAS:", winners.Load())

	fmt.Println("\n=== atomic.Pointer: swap a whole config, readers need no lock ===")
	var current atomic.Pointer[Config]
	current.Store(&Config{Version: 1, Name: "v1"})
	cfg := current.Load()                          // readers take a snapshot and never modify it
	current.Store(&Config{Version: 2, Name: "v2"}) // writers publish a NEW object
	fmt.Println("old snapshot:", *cfg, " current:", *current.Load())
}

func storeMax(m *atomic.Int64, v int64) {
	for {
		cur := m.Load()
		if v <= cur || m.CompareAndSwap(cur, v) {
			return
		}
		// CAS failed: another goroutine changed it, retry with the fresh value
	}
}
