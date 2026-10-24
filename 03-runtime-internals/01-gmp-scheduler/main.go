// Shows G/M/P facts you can observe: GOMAXPROCS, preemption, cheap goroutines.
// Run: go run ./03-runtime-internals/01-gmp-scheduler
package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	fmt.Println("=== runtime facts ===")
	fmt.Println("NumCPU     =", runtime.NumCPU())
	fmt.Println("GOMAXPROCS =", runtime.GOMAXPROCS(0), "(number of Ps)")

	fmt.Println("=== async preemption (1 P, tight loop) ===")
	runtime.GOMAXPROCS(1)
	go func() {
		for { // never yields, never calls a function
		}
	}()
	start := time.Now()
	time.Sleep(20 * time.Millisecond)
	// Before Go 1.14 this could hang forever on 1 P. Now the loop is preempted.
	fmt.Println("main woke up on the only P while another goroutine spins:", time.Since(start) >= 20*time.Millisecond)
	runtime.GOMAXPROCS(runtime.NumCPU())

	fmt.Println("=== 10,000 parked goroutines ===")
	block := make(chan struct{})
	var wg sync.WaitGroup
	for range 10_000 {
		wg.Go(func() { <-block }) // parked on a channel: no OS thread needed
	}
	time.Sleep(50 * time.Millisecond)
	fmt.Println("goroutines alive:", runtime.NumGoroutine() > 10_000)
	close(block)
	wg.Wait()
	fmt.Println("all released")
}
