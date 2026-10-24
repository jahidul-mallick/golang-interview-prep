// Channels: unbuffered, buffered, close, range, nil channel.
// Run: go run ./02-concurrency/02-channels-unbuffered-buffered-closing
package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	fmt.Println("=== unbuffered: send waits for a receiver ===")
	ch := make(chan string)
	done := make(chan struct{})
	go func() {
		defer close(done)
		fmt.Println("receiver got:", <-ch)
	}()
	ch <- "hello" // blocks until the receiver takes it
	<-done

	fmt.Println("\n=== buffered: blocks only when full or empty ===")
	buf := make(chan int, 2)
	buf <- 1
	buf <- 2
	fmt.Println("len =", len(buf), "cap =", cap(buf)) // a 3rd send would block

	fmt.Println("\n=== close: drain buffer first, then zero value and ok=false ===")
	close(buf)
	for range 3 {
		v, ok := <-buf
		fmt.Println("v =", v, "ok =", ok)
	}

	fmt.Println("\n=== producer closes, consumer ranges ===")
	fmt.Println("sum =", sum(produce(5)))

	fmt.Println("\n=== nil channel disables a select case (merge) ===")
	for v := range merge(produce(2), produce(3)) {
		fmt.Print(v, " ")
	}
	fmt.Println()

	fmt.Println("\n=== leak: sender blocked after receiver gave up ===")
	base := runtime.NumGoroutine()
	slow := func(out chan int) { time.Sleep(20 * time.Millisecond); out <- 42 }
	unbuf, oneSlot := make(chan int), make(chan int, 1)
	go slow(unbuf)
	go slow(oneSlot)
	time.Sleep(50 * time.Millisecond) // the caller timed out and never received
	fmt.Printf("goroutines %d -> %d (unbuffered sender stuck, buffered one finished)\n", base, runtime.NumGoroutine())
	<-unbuf // release the stuck sender
}

func produce(n int) <-chan int { // the sender owns the channel and closes it
	out := make(chan int)
	go func() {
		defer close(out)
		for i := range n {
			out <- i
		}
	}()
	return out
}

func sum(in <-chan int) int {
	total := 0
	for v := range in { // ends when closed and drained
		total += v
	}
	return total
}

func merge(a, b <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for a != nil || b != nil {
			select {
			case v, ok := <-a:
				if !ok {
					a = nil // a nil channel is never ready
					continue
				}
				out <- v
			case v, ok := <-b:
				if !ok {
					b = nil
					continue
				}
				out <- v
			}
		}
	}()
	return out
}
