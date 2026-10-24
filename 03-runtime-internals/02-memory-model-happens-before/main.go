// Correct ways to publish data between goroutines (each creates a happens-before edge).
// Run: go run ./03-runtime-internals/02-memory-model-happens-before
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Config struct {
	Name  string
	Limit int
}

func main() {
	fmt.Println("=== channel: send happens-before the receive completes ===")
	ch := make(chan *Config)
	go func() {
		ch <- &Config{Name: "prod", Limit: 10} // build fully, then send
	}()
	c := <-ch
	fmt.Println("receiver sees all fields:", c.Name, c.Limit)

	fmt.Println("=== atomic.Pointer: Store happens-before a Load that sees it ===")
	var cur atomic.Pointer[Config]
	cur.Store(&Config{Name: "v1", Limit: 1})
	done := make(chan struct{})
	go func() {
		defer close(done)
		for cur.Load().Name != "v2" { // plain "for !ready {}" on a bool would be a race
			time.Sleep(time.Millisecond)
		}
		fmt.Println("reader saw the new config:", *cur.Load())
	}()
	cur.Store(&Config{Name: "v2", Limit: 2})
	<-done

	fmt.Println("=== WaitGroup + Mutex: Done happens-before Wait returns ===")
	var (
		wg    sync.WaitGroup
		mu    sync.Mutex
		total int
	)
	for i := 1; i <= 8; i++ {
		wg.Go(func() {
			mu.Lock()
			total += i
			mu.Unlock()
		})
	}
	wg.Wait()
	fmt.Println("total (safe to read after Wait):", total)

	fmt.Println("=== buffered channel as a semaphore (cap 3) ===")
	sem := make(chan struct{}, 3)
	var running, maxSeen atomic.Int32
	for range 12 {
		wg.Go(func() {
			sem <- struct{}{}
			n := running.Add(1)
			for {
				m := maxSeen.Load()
				if n <= m || maxSeen.CompareAndSwap(m, n) {
					break
				}
			}
			time.Sleep(2 * time.Millisecond)
			running.Add(-1)
			<-sem
		})
	}
	wg.Wait()
	fmt.Println("max concurrent:", maxSeen.Load(), "(never above 3)")
}
