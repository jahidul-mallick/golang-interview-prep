# sync.WaitGroup and sync.Once

[← Section index](../README.md) · Prev: [sync.Mutex and sync.RWMutex](../04-sync-mutex-rwmutex/README.md) · Next: [sync.Cond, sync.Pool, sync.Map](../06-sync-cond-pool-map/README.md)

**Interview weight: 1.2%**  ·  Revise in: 6 min  ·  Run: `go run ./02-concurrency/05-sync-waitgroup-once`

## In one line
WaitGroup waits for a group of goroutines to finish, and Once runs an initialization exactly one time even when many goroutines race to trigger it.

## Must know
- `Add(n)` before starting the goroutine, `Done()` (use `defer`) inside it, `Wait()` in the parent.
- Go 1.25: `wg.Go(func(){...})` does Add, go and deferred Done in one call.
- Never copy a WaitGroup: pass a pointer or capture it in a closure.
- Counter below zero panics: "negative WaitGroup counter".
- Collect results safely: write to your own slice index, guard an append with a mutex, or send on a channel that a separate goroutine closes after `wg.Wait()`.
- `sync.Once.Do(f)`: f runs once, and every other caller blocks until f returns. Calling `Do` inside `f` deadlocks.
- If f panics, Once still counts it as done. A Once cannot be reset or retried.
- `sync.OnceFunc`, `OnceValue`, `OnceValues` (Go 1.21) give a lazy function (runs on first call, result cached). `OnceValues` caches the error too.
- WaitGroup has no error handling and no cancellation. For that use `errgroup` (golang.org/x/sync, not stdlib).

## Code to remember
```go
var wg sync.WaitGroup
results := make(chan int)
for i := range 5 {
    wg.Go(func() { results <- i * i })
}
go func() { wg.Wait(); close(results) }() // separate closer
for v := range results { fmt.Println(v) }

var once sync.Once
once.Do(initConfig) // safe from many goroutines
```

## Common traps
- `wg.Add(1)` inside the goroutine: `Wait` can return before it runs.
- Closing the results channel from a worker while others still send (panic).
- Using Once for init that can fail and must be retried: the failure is cached forever.
- Calling `Wait` before all `Add` calls have happened.

## Interview Q&A

**Q1. Why must Add be called before `go`?**

**A:** Otherwise `Wait` might see counter zero and return before the goroutine registers itself.

**Q2. How do you collect results from many goroutines?**

**A:** Pre-sized slice with one index per goroutine (no lock), or a results channel closed by one goroutine after `wg.Wait()`. A mutex around `append` also works.

**Q3. How does sync.Once work and what is it used for?**

**A:** It uses an atomic flag plus a mutex so the function runs once and all callers see its effects. Used for lazy singletons, one-time setup, closing resources once.

**Q4. What if the function passed to Once panics or fails?**

**A:** It is still marked done and will not run again. For retry, use a mutex with a nil check.

**Q5. What happens if you pass a `WaitGroup` by value to a goroutine?**

**A:** The goroutine calls `Done` on a copy, so the original counter never reaches zero and `Wait` blocks forever. Pass `*sync.WaitGroup` or use a closure. `go vet` flags it.
