// sync.WaitGroup (collecting results) and sync.Once / OnceValue.
// Run: go run ./02-concurrency/05-sync-waitgroup-once
package main

import (
	"fmt"
	"sync"
)

type Config struct{ Region string }

func main() {
	fmt.Println("=== WaitGroup: results by index (no lock needed) ===")
	var wg sync.WaitGroup
	res := make([]int, 5)
	for i := range res {
		wg.Go(func() { res[i] = i * 10 }) // each goroutine writes its own slot
	}
	wg.Wait()
	fmt.Println(res)

	fmt.Println("\n=== WaitGroup: results channel closed by a separate closer ===")
	results := make(chan int)
	for i := range 5 {
		wg.Go(func() { results <- i * i })
	}
	go func() {
		wg.Wait() // never close from a worker: another may still send
		close(results)
	}()
	sum := 0
	for v := range results {
		sum += v
	}
	fmt.Println("sum of squares =", sum)

	fmt.Println("\n=== sync.Once: 10 goroutines, one init ===")
	var once sync.Once
	inits := 0
	for range 10 {
		wg.Go(func() { once.Do(func() { inits++ }) })
	}
	wg.Wait()
	fmt.Println("init ran", inits, "time")

	fmt.Println("\n=== sync.OnceValue: lazy singleton (Go 1.21) ===")
	loadConfig := sync.OnceValue(func() Config {
		fmt.Println("loading config once...")
		return Config{Region: "eu-west-1"}
	})
	fmt.Println(loadConfig())
	fmt.Println(loadConfig()) // cached, no second load

	fmt.Println("\n=== Once cannot retry: a failed init is final ===")
	var o sync.Once
	attempts := 0
	for range 3 {
		o.Do(func() { attempts++ }) // even if this failed, Do never runs it again
	}
	fmt.Println("attempts:", attempts, "(use a mutex + nil check if init must be retried)")
}
