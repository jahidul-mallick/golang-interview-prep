# pprof Profiling & Tracing

[← Section index](../README.md) · Prev: [Garbage Collector](../05-garbage-collector/README.md) · Next: [Data Structure Internals & Memory Layout](../07-data-structure-internals-memory-layout/README.md)

**Interview weight: 1.0%**  ·  Revise in: 6 min  ·  Run: `go run ./03-runtime-internals/06-pprof-profiling-tracing`

## In one line
Use pprof to find where CPU time and memory go, and the execution tracer to see how goroutines are scheduled over time.

## Must know
- Symptom to profile: high CPU = `profile`; memory growing = `heap`; goroutines growing = `goroutine`; slow but CPU idle = `block`, `mutex` or the tracer.
- Block and mutex profiles are off by default (`runtime.SetBlockProfileRate`, `runtime.SetMutexProfileFraction`).
- Add `import _ "net/http/pprof"` and serve on an internal-only port. Endpoints live under `/debug/pprof/`.
- Analyze with `go tool pprof <file or url>`: `top`, `list <func>`, `web`. Flame graph: `go tool pprof -http=:8080 <file>`.
- Heap profile: `inuse_space` = live memory now (leaks), `alloc_space` = total allocated since start (GC pressure).
- Goroutine leak: take `/debug/pprof/goroutine?debug=1` twice and compare counts and stacks.
- Benchmarks can write profiles: `go test -bench . -cpuprofile cpu.out -memprofile mem.out`.
- Execution tracer: `go test -trace trace.out` or `runtime/trace`, view with `go tool trace trace.out`.
- Process: measure, find the hot spot, change one thing, measure again.

## Code to remember
```go
import _ "net/http/pprof" // registers handlers on http.DefaultServeMux

go http.ListenAndServe("127.0.0.1:6060", nil) // keep it off the public port
// go tool pprof http://127.0.0.1:6060/debug/pprof/profile?seconds=30
```

## Common traps
- Exposing `/debug/pprof` publicly.
- Profiling an idle service, or a run too short to collect samples.
- Optimizing without a profile. The hot spot is rarely where you guess.

## Interview Q&A

**Q1. How do you debug high CPU in production?**

**A:** Pull a 30 second CPU profile from `/debug/pprof/profile`, open it with `go tool pprof`, look at `top` and the flame graph, fix the hot function, then compare before and after.

**Q2. How do you find a goroutine leak?**

**A:** Watch `runtime.NumGoroutine()` or fetch `/debug/pprof/goroutine?debug=1` over time. Growing counts with the same stack mean goroutines stuck on a channel or lock.

**Q3. Memory keeps growing. How do you find the leak?**

**A:** Take two heap profiles some time apart and compare: `go tool pprof -diff_base old.pb.gz new.pb.gz`, with `inuse_space`. Also check the goroutine count.

**Q4. pprof vs trace?**

**A:** pprof samples where time or memory goes. The tracer records events (scheduling, GC, syscalls) on a timeline, useful for latency and scheduling problems.

**Q5. Block profile vs mutex profile?**

**A:** Block shows where goroutines wait (channels, locks, `select`). Mutex shows contention on locks, attributed to the code holding them. Both are off by default, so enable them first.
