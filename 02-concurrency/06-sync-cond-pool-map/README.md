# sync.Cond, sync.Pool, sync.Map

[← Section index](../README.md) · Prev: [sync.WaitGroup and sync.Once](../05-sync-waitgroup-once/README.md) · Next: [sync/atomic](../07-atomic-operations/README.md)

**Interview weight: 0.8%**  ·  Revise in: 6 min  ·  Run: `go run ./02-concurrency/06-sync-cond-pool-map`

## In one line
Cond waits for a condition under a lock, Pool recycles temporary objects to cut allocations, and sync.Map is a concurrent map for a few specific access patterns.

## Must know
- `sync.Cond`: `c.L.Lock(); for !cond { c.Wait() }; ...; c.L.Unlock()`. Wait unlocks, sleeps, re-locks. Always use a `for` loop, never `if`.
- `Signal` wakes one waiter, `Broadcast` wakes all. For a one-shot signal a closed channel is simpler.
- `sync.Pool`: `Get`/`Put` of temporary objects (for example `bytes.Buffer`). It is a cache, not storage.
- The GC can empty a Pool at any time, so never store things that must persist or must be closed.
- Reset an object before reuse. Pool needs a `New` func or `Get` returns nil.
- `sync.Map`: methods `Load`, `Store`, `LoadOrStore`, `Delete`, `Range`. Keys and values are `any`.
- Use `sync.Map` for write-once-read-many keys (caches that only grow) or goroutines touching disjoint keys. Otherwise a `map` with `RWMutex` is typed and usually clearer.
- Never copy a `Cond`, `Pool` or `Map` after first use.

## Code to remember
```go
var bufPool = sync.Pool{New: func() any { return new(bytes.Buffer) }}

func render(s string) string {
    b := bufPool.Get().(*bytes.Buffer)
    defer bufPool.Put(b)
    b.Reset()
    b.WriteString(s)
    return b.String()
}
```

## Common traps
- `if !cond { c.Wait() }` instead of `for`: spurious or stale wake-ups break correctness.
- Putting an object in a Pool and still using it afterwards (it may be handed to someone else).
- Pooling objects that grow huge: one big buffer stays in the pool and wastes memory.
- Using `sync.Map` as the default map and losing type safety for no gain.

## Interview Q&A

**Q1. What is sync.Pool for?**

**A:** Reducing allocation and GC pressure for short-lived, same-shaped objects, for example buffers in a hot request path. `fmt` and `encoding/json` use this idea internally.

**Q2. Why a `for` loop around `cond.Wait()`?**

**A:** Another goroutine may have changed the state between the wake-up and you re-acquiring the lock, so the condition must be re-checked.

**Q3. sync.Map vs map plus mutex?**

**A:** `sync.Map` wins when keys are written once and read many times, or goroutines use disjoint keys. For general read/write mixes, map plus RWMutex is type-safe and clearer.

**Q4. Cond vs channels?**

**A:** Channels cover most cases. Cond fits when many goroutines wait on shared state that is already protected by a mutex, and you need Broadcast on a changing condition.
