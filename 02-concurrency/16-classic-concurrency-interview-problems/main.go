// The classic live-coding problems: alternate print, producer-consumer,
// bounded parallelism with results, merge channels, timeout a call.
// Run: go run ./02-concurrency/16-classic-concurrency-interview-problems
package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Two goroutines pass a token back and forth; the buffer of 1 absorbs the last handover.
func alternate(n int) {
	odd, even := make(chan struct{}, 1), make(chan struct{}, 1)
	odd <- struct{}{} // odd goes first
	var wg sync.WaitGroup
	wg.Go(func() {
		for i := 1; i <= n; i += 2 {
			<-odd
			fmt.Print(i, " ")
			even <- struct{}{}
		}
	})
	wg.Go(func() {
		for i := 2; i <= n; i += 2 {
			<-even
			fmt.Print(i, " ")
			odd <- struct{}{}
		}
	})
	wg.Wait()
	fmt.Println()
}

// runAll runs tasks with at most limit at once and returns results in input order.
func runAll(limit int, tasks []func() int) []int {
	out := make([]int, len(tasks))
	sem := make(chan struct{}, limit)
	var wg sync.WaitGroup
	for i, t := range tasks {
		sem <- struct{}{} // acquire before go: at most limit goroutines exist
		wg.Go(func() {
			defer func() { <-sem }()
			out[i] = t() // each goroutine owns one index, so no lock
		})
	}
	wg.Wait()
	return out
}

// merge is fan-in: one forwarder per input, close out after all are done.
func merge(ins ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	for _, in := range ins {
		wg.Go(func() {
			for v := range in {
				out <- v
			}
		})
	}
	go func() { wg.Wait(); close(out) }()
	return out
}

func emit(vals ...int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for _, v := range vals {
			ch <- v
		}
	}()
	return ch
}

// callWithTimeout gives up after d. The result channel has room for 1, so the
// slow goroutine can still send and exit instead of blocking forever.
func callWithTimeout(d time.Duration, f func() string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()
	res := make(chan string, 1)
	go func() { res <- f() }()
	select {
	case s := <-res:
		return s, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func main() {
	fmt.Println("=== alternate print 1..6 ===")
	alternate(6)

	fmt.Println("=== producer-consumer, bounded buffer (cap 2) ===")
	buf := make(chan int, 2)
	go func() { // the producer closes the channel
		defer close(buf)
		for i := 1; i <= 5; i++ {
			buf <- i // blocks while the buffer is full
		}
	}()
	sum := 0
	for v := range buf {
		sum += v
	}
	fmt.Println("consumed sum:", sum, "(want 15)")

	fmt.Println("=== run 6 tasks, 2 at a time, ordered results ===")
	var tasks []func() int
	for i := range 6 {
		tasks = append(tasks, func() int { time.Sleep(10 * time.Millisecond); return i * i })
	}
	fmt.Println(runAll(2, tasks))

	fmt.Println("=== merge two channels (order not guaranteed, sum is) ===")
	total := 0
	for v := range merge(emit(1, 2, 3), emit(10, 20)) {
		total += v
	}
	fmt.Println("merged sum:", total, "(want 36)")

	fmt.Println("=== timeout a slow call ===")
	slow := func() string { time.Sleep(200 * time.Millisecond); return "done" }
	fmt.Println(callWithTimeout(20*time.Millisecond, slow))
	fmt.Println(callWithTimeout(500*time.Millisecond, slow))
}
