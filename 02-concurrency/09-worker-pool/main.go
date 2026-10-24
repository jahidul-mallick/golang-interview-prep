// Worker pool: the canonical pool and a cancellable variant.
// Run: go run ./02-concurrency/09-worker-pool
package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type result struct{ id, out int }

// pool runs n workers and returns squares in input order. Each job is its own
// index here; real code would send {index, payload}.
func pool(n int, jobs []int) []int {
	jobCh := make(chan int)
	resCh := make(chan result)

	var wg sync.WaitGroup
	for range n {
		wg.Add(1) // Add before go, in the spawning goroutine
		go func() {
			defer wg.Done()
			for j := range jobCh {
				resCh <- result{j, j * j}
			}
		}()
	}

	go func() { // producer owns and closes jobCh
		defer close(jobCh)
		for _, j := range jobs {
			jobCh <- j
		}
	}()
	go func() { wg.Wait(); close(resCh) }() // close results only after all workers exit

	out := make([]int, len(jobs))
	for r := range resCh {
		out[r.id] = r.out // job value doubles as its index here
	}
	return out
}

// poolCtx stops early when ctx is cancelled; no goroutine is left blocked.
func poolCtx(ctx context.Context, n, total int) int {
	jobCh := make(chan int)
	var wg sync.WaitGroup
	var mu sync.Mutex
	done := 0

	for range n {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range jobCh {
				select {
				case <-time.After(30 * time.Millisecond): // simulated work
					mu.Lock()
					done++
					mu.Unlock()
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	go func() {
		defer close(jobCh)
		for j := range total {
			select {
			case jobCh <- j:
			case <-ctx.Done():
				return
			}
		}
	}()

	wg.Wait()
	return done
}

func main() {
	fmt.Println("=== canonical pool (3 workers, ordered output) ===")
	jobs := []int{0, 1, 2, 3, 4, 5, 6, 7}
	fmt.Println("squares:", pool(3, jobs))

	fmt.Println("=== cancellable pool (cancel after 100ms) ===")
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	done := poolCtx(ctx, 2, 100)
	fmt.Println("jobs done before cancel:", done, "of 100, stopped early:", done < 100)
}
