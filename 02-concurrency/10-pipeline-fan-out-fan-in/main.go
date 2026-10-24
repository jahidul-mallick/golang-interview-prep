// Pipeline, fan-out and fan-in with cancellation.
// Run: go run ./02-concurrency/10-pipeline-fan-out-fan-in
package main

import (
	"context"
	"fmt"
	"sort"
	"sync"
)

// Each stage owns and closes its output channel, and stops on ctx.Done().
func gen(ctx context.Context, nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			select {
			case out <- n:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

func square(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			select {
			case out <- n * n:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

// merge is fan-in: one forwarder per input, close out after all finish.
func merge(ctx context.Context, ins ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	for _, in := range ins {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for v := range in {
				select {
				case out <- v:
				case <-ctx.Done():
					return
				}
			}
		}()
	}
	go func() { wg.Wait(); close(out) }()
	return out
}

func main() {
	fmt.Println("=== pipeline: gen -> square ===")
	ctx, cancel := context.WithCancel(context.Background())
	for v := range square(ctx, gen(ctx, 1, 2, 3, 4)) {
		fmt.Println(v)
	}
	cancel()

	fmt.Println("=== fan-out (3 workers on one channel) + fan-in (merge) ===")
	ctx, cancel = context.WithCancel(context.Background())
	src := gen(ctx, 1, 2, 3, 4, 5, 6)
	workers := []<-chan int{square(ctx, src), square(ctx, src), square(ctx, src)}
	var got []int
	for v := range merge(ctx, workers...) {
		got = append(got, v)
	}
	sort.Ints(got) // fan-in order is not deterministic
	fmt.Println("sorted results:", got)
	cancel()

	fmt.Println("=== early stop: take first 2, cancel the rest ===")
	ctx, cancel = context.WithCancel(context.Background())
	out := square(ctx, gen(ctx, 1, 2, 3, 4, 5, 6))
	fmt.Println(<-out, <-out)
	cancel() // upstream goroutines unblock via ctx.Done() instead of leaking
	fmt.Println("cancelled; no goroutine left blocked on a send")
}
