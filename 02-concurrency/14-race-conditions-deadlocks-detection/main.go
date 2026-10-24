// Fixing races, avoiding lock-order deadlocks, spotting goroutine leaks.
// The racy version is only shown in a comment: run with -race to see reports.
// Run: go run ./02-concurrency/14-race-conditions-deadlocks-detection
package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type account struct {
	mu      sync.Mutex
	id      int
	balance int
}

// transfer locks the lower id first. Without a fixed order, transfer(a,b)
// and transfer(b,a) running together can each hold one lock and deadlock.
func transfer(from, to *account, amount int) {
	first, second := from, to
	if second.id < first.id {
		first, second = second, first
	}
	first.mu.Lock()
	defer first.mu.Unlock()
	second.mu.Lock()
	defer second.mu.Unlock()
	from.balance -= amount
	to.balance += amount
}

func main() {
	fmt.Println("=== data race fixed: mutex and atomic ===")
	// Racy: `n++` from many goroutines is read-modify-write, not atomic.
	var mu sync.Mutex
	var wg sync.WaitGroup
	n := 0
	var atom atomic.Int64
	for range 100 {
		wg.Go(func() {
			for range 100 {
				mu.Lock()
				n++
				mu.Unlock()
				atom.Add(1)
			}
		})
	}
	wg.Wait()
	fmt.Println("mutex:", n, "atomic:", atom.Load(), "(want 10000)")

	fmt.Println("=== lock ordering: opposite transfers do not deadlock ===")
	a, b := &account{id: 1, balance: 100}, &account{id: 2, balance: 100}
	for range 50 {
		wg.Go(func() { transfer(a, b, 1) })
		wg.Go(func() { transfer(b, a, 1) })
	}
	wg.Wait()
	fmt.Println("total balance:", a.balance+b.balance, "(want 200)")

	fmt.Println("=== goroutine leak: count keeps growing ===")
	before := runtime.NumGoroutine()
	for range 5 {
		ch := make(chan int) // nobody ever sends: the goroutine blocks forever
		go func() { <-ch }()
	}
	time.Sleep(10 * time.Millisecond)
	fmt.Println("goroutines before:", before, "after:", runtime.NumGoroutine())
	fmt.Println("fix: give the goroutine a ctx.Done() or close the channel")
}
