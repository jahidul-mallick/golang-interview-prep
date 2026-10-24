// Timer, Ticker, AfterFunc, and select with a timeout (Go 1.23+ semantics).
// Run: go run ./02-concurrency/15-timers-tickers-time
package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("=== Timer (fires once) ===")
	t := time.NewTimer(50 * time.Millisecond)
	<-t.C
	fmt.Println("timer fired")
	t.Reset(20 * time.Millisecond) // Go 1.23+: safe, no manual drain needed
	<-t.C
	fmt.Println("timer fired again after Reset")

	fmt.Println("=== Ticker (repeats) ===")
	tick := time.NewTicker(30 * time.Millisecond)
	defer tick.Stop() // still good practice
	for i := 1; i <= 3; i++ {
		<-tick.C
		fmt.Println("tick", i)
	}

	fmt.Println("=== AfterFunc (runs f in its own goroutine) ===")
	done := make(chan struct{})
	time.AfterFunc(40*time.Millisecond, func() {
		fmt.Println("AfterFunc fired")
		close(done)
	})
	<-done

	fmt.Println("=== select with timeout ===")
	ch := make(chan int)
	select {
	case v := <-ch:
		fmt.Println("got", v)
	case <-time.After(20 * time.Millisecond):
		fmt.Println("timed out (expected)")
	}

	fmt.Println("=== time.Since ===")
	start := time.Now()
	time.Sleep(10 * time.Millisecond)
	fmt.Println("elapsed >= 10ms:", time.Since(start) >= 10*time.Millisecond)
}
