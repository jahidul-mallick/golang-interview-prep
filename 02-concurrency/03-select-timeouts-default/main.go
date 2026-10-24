// select: timeout, non-blocking default, for-select loop with labeled break.
// Run: go run ./02-concurrency/03-select-timeouts-default
package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("=== timeout with time.After ===")
	slow := make(chan string, 1) // buffer 1: the late sender never blocks
	go func() {
		time.Sleep(50 * time.Millisecond)
		slow <- "late result"
	}()
	select {
	case r := <-slow:
		fmt.Println("got", r)
	case <-time.After(10 * time.Millisecond):
		fmt.Println("timed out after 10ms")
	}

	fmt.Println("\n=== default: non-blocking receive and send ===")
	ch := make(chan int, 1)
	select {
	case v := <-ch:
		fmt.Println("received", v)
	default:
		fmt.Println("nothing to receive, moved on")
	}
	for i := 1; i <= 2; i++ {
		select {
		case ch <- i:
			fmt.Println("sent", i)
		default:
			fmt.Println("buffer full, dropped", i)
		}
	}

	fmt.Println("\n=== for-select loop with a label ===")
	events := make(chan int)
	quit := make(chan struct{})
	go func() {
		for i := 1; i <= 3; i++ {
			events <- i
		}
		close(quit)
	}()
	tick := time.NewTicker(time.Hour) // stand-in for a periodic job
	defer tick.Stop()
loop:
	for {
		select {
		case e := <-events:
			fmt.Println("event", e)
		case <-tick.C:
			fmt.Println("tick")
		case <-quit:
			break loop // a bare break would only leave the select
		}
	}
	fmt.Println("loop exited")

	fmt.Println("\n=== several ready cases: one is chosen at random ===")
	a, b := make(chan int, 1), make(chan int, 1)
	var na, nb int
	for range 1000 {
		a <- 1
		b <- 1
		select {
		case <-a:
			na++
			<-b
		case <-b:
			nb++
			<-a
		}
	}
	fmt.Println("a and b both chosen:", na > 0 && nb > 0)
}
