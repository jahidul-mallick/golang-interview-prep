# Goroutines Basics

[← Section index](../README.md) · Next: [Channels: Unbuffered, Buffered, Closing](../02-channels-unbuffered-buffered-closing/README.md)

**Interview weight: 1.4%**  ·  Revise in: 6 min  ·  Run: `go run ./02-concurrency/01-goroutines-basics`

## In one line
A goroutine is a function running on a small growable stack, scheduled by the Go runtime onto a few OS threads, so you can have thousands of them cheaply.

## Must know
- Start one with `go f(x)`: `f` and `x` are evaluated immediately, the call runs later.
- Stack starts at about 2 KB and grows as needed. OS thread stacks are fixed and much larger (MBs).
- The runtime runs many goroutines (G) on few OS threads (M) via scheduler slots (P). Switching needs no kernel call, so it is cheap.
- When `main` returns, the program exits. Other goroutines are killed, no defers run.
- A panic in any goroutine crashes the whole process unless that goroutine recovers it itself.
- Go 1.22+: each loop iteration has its own loop variable, so capturing `i` in a goroutine is safe.
- `wg.Go(f)` (Go 1.25) does `Add(1)`, starts the goroutine and calls `Done` for you.
- A goroutine leak is a goroutine blocked forever (channel, lock, no ctx). Give every goroutine an exit path.

## Code to remember
```go
var wg sync.WaitGroup
for i := range 3 {
    wg.Add(1)               // Add BEFORE the go statement
    go func() {
        defer wg.Done()     // Done via defer, runs even on panic
        fmt.Println("worker", i)
    }()
}
wg.Wait()
```

## Common traps
- `wg.Add(1)` inside the goroutine: `Wait` may run before it. Always Add in the parent.
- Forgetting that `main` exiting ends everything. Wait with a WaitGroup, not `time.Sleep`.
- `recover()` in the parent does not catch a panic in the child goroutine.
- Spawning a goroutine per item with no bound and no cancellation. This leaks and exhausts memory.

## Interview Q&A

**Q1. Goroutine vs OS thread?**

**A:** Goroutines are managed by the Go runtime, start at about 2 KB of stack and switch in user space. Threads are managed by the kernel, have MB stacks and costly context switches. Many goroutines share few threads.

**Q2. Can you kill a goroutine from outside? What happens when `main` exits?**

**A:** You cannot kill it. Signal it with a `context` or a done channel and have it return. If `main` exits first, all goroutines are dropped, so wait with a WaitGroup.

**Q3. How do you find a goroutine leak?**

**A:** Watch `runtime.NumGoroutine()` grow, and dump stacks with pprof `goroutine` profile (or SIGQUIT). Look for many goroutines parked on the same line.

**Q4. Is the classic loop-variable closure bug still a thing?**

**A:** Not since Go 1.22 for `for` loops declared in go.mod 1.22+. In older code, shadow with `i := i` or pass `i` as an argument.

**Q5. What does `GOMAXPROCS` control?**

**A:** The number of Ps, meaning how many goroutines can run Go code at the same instant. It defaults to the number of CPUs.

**Q6. What are the common causes of goroutine leaks?**

**A:** Send or receive on a channel nobody serves, a `range` over a channel never closed, a missing `ctx.Done()` case, and a lock never released. Fix: give each goroutine an exit path, and buffer a result channel of size 1 when the caller may time out.
