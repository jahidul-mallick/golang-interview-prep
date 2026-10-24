# errgroup, Semaphore & Bounded Concurrency

[← Section index](../README.md) · Prev: [Pipeline, Fan-out & Fan-in](../10-pipeline-fan-out-fan-in/README.md) · Next: [Graceful Shutdown (HTTP Server)](../12-graceful-shutdown-http-server/README.md)

**Interview weight: 1.0%**  ·  Revise in: 6 min  ·  Run: `go run ./02-concurrency/11-errgroup-semaphore-bounded-concurrency`

## In one line
errgroup runs goroutines that can fail, cancels the rest on the first error, and a buffered channel (a semaphore: a counter of free slots) is the stdlib way to cap how many run at once.

## Must know
- `golang.org/x/sync/errgroup` is not stdlib. `g.Go(func() error)` starts a task, `g.Wait()` returns the first non-nil error
- `errgroup.WithContext(ctx)` gives a ctx that is cancelled on the first error or when `Wait` returns
- `g.SetLimit(n)` caps live goroutines; `g.Go` then blocks until a slot is free
- Tasks must watch `ctx.Done()`, otherwise cancellation does nothing
- Semaphore: `make(chan struct{}, n)`; send to acquire, receive to release
- `golang.org/x/sync/semaphore` adds weighted, context-aware `Acquire`
- Go 1.25 `sync.WaitGroup.Go` removes the `Add`/`Done` boilerplate but has no error or cancel handling
- `main.go` here has a small stdlib errgroup clone so you can see how it works

## Code to remember
```go
sem := make(chan struct{}, 3) // at most 3 at once
var wg sync.WaitGroup
for _, u := range urls {
	sem <- struct{}{} // acquire BEFORE go: the loop blocks, no goroutine pile-up
	wg.Go(func() {
		defer func() { <-sem }() // release
		fetch(u)
	})
}
wg.Wait()
```
With errgroup (needs the x/sync module): `g, ctx := errgroup.WithContext(ctx); g.SetLimit(3); g.Go(...); err := g.Wait()`.

## Common traps
- Tasks that ignore the group ctx, so siblings keep running after a failure
- Forgetting that `Wait` returns only the first error; the others are dropped
- Acquiring inside the goroutine: one goroutine per task still exists, only their work is limited
- Releasing a semaphore slot you did not acquire, or leaking one on an early return (use `defer`)

## Interview Q&A

**Q1. errgroup or WaitGroup?**

**A:** errgroup when tasks can fail and you want the first error plus cancellation of the rest. WaitGroup when you only need to wait, or you handle errors yourself.

**Q2. How does a channel semaphore limit concurrency?**

**A:** Capacity is the limit. A send blocks when the buffer is full, so the next goroutine waits until another receives and frees a slot.

**Q3. How do you process 10,000 items with at most 20 in flight?**

**A:** `errgroup` with `SetLimit(20)`, or a pool of 20 workers reading from a channel. Do not start 10,000 goroutines and hope.

**Q4. What happens if a task panics?**

**A:** An unrecovered panic in any goroutine crashes the whole program. Recover inside the task if you want to turn it into an error.
