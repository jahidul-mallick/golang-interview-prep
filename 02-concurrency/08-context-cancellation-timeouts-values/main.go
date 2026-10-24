// context: cancel tree, timeout, values, cooperative worker.
// Run: go run ./02-concurrency/08-context-cancellation-timeouts-values
package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type ctxKey struct{} // unexported key type avoids collisions between packages

func main() {
	fmt.Println("=== cancel flows down the tree, never up ===")
	root, cancelRoot := context.WithCancel(context.Background())
	defer cancelRoot()
	child, cancelChild := context.WithCancel(root)
	defer cancelChild()
	sibling, cancelSibling := context.WithTimeout(root, time.Hour)
	defer cancelSibling()

	cancelChild()
	fmt.Println("after cancelling child:")
	fmt.Println("  root    :", root.Err())
	fmt.Println("  child   :", child.Err())
	fmt.Println("  sibling :", sibling.Err())
	cancelRoot()
	<-sibling.Done()
	fmt.Println("after cancelling root, sibling:", sibling.Err())

	fmt.Println("\n=== timeout: ctx.Err() is DeadlineExceeded ===")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel() // always call cancel to release the timer
	err := slowOp(ctx, time.Second)
	fmt.Println("err:", err, "| is deadline:", errors.Is(err, context.DeadlineExceeded))

	fmt.Println("\n=== values: request-scoped data only ===")
	vctx := context.WithValue(context.Background(), ctxKey{}, "req-7f3a")
	if id, ok := vctx.Value(ctxKey{}).(string); ok {
		fmt.Println("request id:", id)
	}

	fmt.Println("\n=== cancellation is cooperative ===")
	wctx, stop := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	items := 0
	wg.Go(func() {
		for {
			select {
			case <-wctx.Done():
				return // the worker must check ctx itself
			default:
				time.Sleep(time.Millisecond)
				items++
			}
		}
	})
	time.Sleep(10 * time.Millisecond)
	stop()
	wg.Wait()
	fmt.Println("worker stopped, reason:", wctx.Err(), "| did some work:", items > 0)
}

func slowOp(ctx context.Context, d time.Duration) error {
	select {
	case <-time.After(d):
		return nil
	case <-ctx.Done():
		return fmt.Errorf("slowOp: %w", ctx.Err())
	}
}
