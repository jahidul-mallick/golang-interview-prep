// Bounded concurrency: channel semaphore, WaitGroup.Go, and a tiny errgroup.
// golang.org/x/sync/errgroup is not stdlib, so miniGroup mimics what it does.
// Run: go run ./02-concurrency/11-errgroup-semaphore-bounded-concurrency
package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// miniGroup: first error cancels the shared ctx; Wait returns that first error.
type miniGroup struct {
	wg     sync.WaitGroup
	once   sync.Once
	err    error
	cancel context.CancelFunc
}

func withContext(ctx context.Context) (*miniGroup, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &miniGroup{cancel: cancel}, ctx
}

func (g *miniGroup) Go(f func() error) {
	g.wg.Go(func() {
		if err := f(); err != nil {
			g.once.Do(func() { g.err = err; g.cancel() })
		}
	})
}

func (g *miniGroup) Wait() error {
	g.wg.Wait()
	g.cancel()
	return g.err
}

func main() {
	fmt.Println("=== semaphore: buffered channel limits to 3 at once ===")
	sem := make(chan struct{}, 3)
	var running, peak atomic.Int32
	var wg sync.WaitGroup
	for range 9 {
		sem <- struct{}{} // acquire before starting the goroutine
		wg.Go(func() {
			defer func() { <-sem }() // release
			n := running.Add(1)
			for p := peak.Load(); n > p && !peak.CompareAndSwap(p, n); p = peak.Load() {
			}
			time.Sleep(20 * time.Millisecond)
			running.Add(-1)
		})
	}
	wg.Wait()
	fmt.Println("peak concurrent:", peak.Load())

	fmt.Println("=== errgroup pattern: first error cancels siblings ===")
	g, ctx := withContext(context.Background())
	for i := range 4 {
		g.Go(func() error {
			if i == 2 {
				return errors.New("task 2 failed")
			}
			select {
			case <-time.After(500 * time.Millisecond):
				return nil
			case <-ctx.Done():
				fmt.Println("task", i, "cancelled")
				return ctx.Err()
			}
		})
	}
	fmt.Println("Wait returned:", g.Wait())
}
