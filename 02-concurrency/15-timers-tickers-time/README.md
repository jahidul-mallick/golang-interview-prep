# Timers, Tickers & time Package

[← Section index](../README.md) · Prev: [Race Conditions, Deadlocks & Detection](../14-race-conditions-deadlocks-detection/README.md) · Next: [Classic Concurrency Interview Problems](../16-classic-concurrency-interview-problems/README.md)

**Interview weight: 0.7%**  ·  Revise in: 4 min  ·  Run: `go run ./02-concurrency/15-timers-tickers-time`

## In one line
A Timer fires once, a Ticker fires repeatedly, and since Go 1.23 neither leaks when dropped and `Reset` needs no manual draining.

## Must know
- `time.NewTimer(d)` fires once on `t.C`; `time.NewTicker(d)` fires every `d`
- `time.AfterFunc(d, f)` runs `f` in its own goroutine, no channel
- `time.After(d)` is `NewTimer(d).C`. Fine in a `select`; in a hot loop prefer one reused Timer
- Go 1.23+: unreferenced Timers and Tickers are garbage collected even if not stopped, so the old "time.After leaks until it fires" advice is pre-1.23
- Go 1.23+: `Stop` and `Reset` no longer leave a stale value in the channel, so drain-before-Reset is not needed
- Still call `ticker.Stop()` as good practice
- For cancellable waits, select on `ctx.Done()` as well as the timer
- `time.Now()` carries a monotonic reading, so `time.Since(start)` ignores wall clock jumps

## Code to remember
```go
ticker := time.NewTicker(time.Second)
defer ticker.Stop()
for {
	select {
	case <-ticker.C:
		doWork()
	case <-ctx.Done():
		return
	}
}
```

## Common traps
- Using `time.Sleep` in a goroutine that cannot be cancelled; select on a timer and `ctx.Done()` instead
- Ticker with slow work: ticks are dropped, not queued, so work can run late
- `AfterFunc` callback running concurrently with your own code, so shared state needs a lock
- Comparing wall-clock times across machines or after clock changes for durations

## Interview Q&A

**Q1. Timer versus Ticker?**

**A:** A Timer fires once and can be `Reset`. A Ticker fires every interval until stopped.

**Q2. Does `time.After` in a loop leak?**

**A:** Before Go 1.23 each unfired timer stayed alive until it expired. Since Go 1.23 unreferenced timers are collected, so it is no longer a leak, though reusing one Timer is still cheaper in hot loops.

**Q3. How do you implement a timeout for a channel receive?**

**A:** `select` on the channel and `time.After(d)`, or better, `context.WithTimeout` and `ctx.Done()`.

**Q4. How do you run a job every 5 seconds until shutdown?**

**A:** A `Ticker` in a `for`/`select` with `ctx.Done()`, and `defer ticker.Stop()`. If the job is slower than the interval, ticks are dropped, not queued.
