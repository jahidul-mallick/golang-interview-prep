// Goroutine basics: go statement, WaitGroup, loop variable, leak, panic.
// Run: go run ./02-concurrency/01-goroutines-basics
package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	fmt.Println("=== go statement: arguments are evaluated now ===")
	x := 1
	var wg sync.WaitGroup
	wg.Add(1)
	go func(v int) { // v is a copy taken at the go statement
		defer wg.Done()
		fmt.Println("goroutine sees v =", v)
	}(x)
	x = 2
	wg.Wait()
	fmt.Println("main sees x =", x)

	fmt.Println("\n=== loop variable is per-iteration since Go 1.22 ===")
	squares := make([]int, 5)
	var wg2 sync.WaitGroup
	for i := range 5 {
		wg2.Go(func() { squares[i] = i * i }) // wg.Go is Go 1.25: Add(1) + go + defer Done()
	}
	wg2.Wait()
	fmt.Println("squares:", squares)

	fmt.Println("\n=== goroutine leak: parked forever on a channel ===")
	base := runtime.NumGoroutine()
	never := make(chan int)
	go func() { <-never }() // nobody ever sends
	time.Sleep(5 * time.Millisecond)
	fmt.Printf("goroutines: %d -> %d (leaked 1)\n", base, runtime.NumGoroutine())
	close(never) // fix: the owner gives every goroutine an exit path

	fmt.Println("\n=== panic must be recovered inside the same goroutine ===")
	var ran atomic.Bool
	var wg3 sync.WaitGroup
	wg3.Go(func() {
		defer func() { fmt.Println("recovered:", recover()) }()
		panic("boom")
	})
	wg3.Go(func() { ran.Store(true) })
	wg3.Wait()
	fmt.Println("process survived, other goroutine ran:", ran.Load())
	// Without the recover, a panic in any goroutine crashes the whole process.
}
