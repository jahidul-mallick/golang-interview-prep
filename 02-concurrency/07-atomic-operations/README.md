# sync/atomic

[← Section index](../README.md) · Prev: [sync.Cond, sync.Pool, sync.Map](../06-sync-cond-pool-map/README.md) · Next: [Context: Cancellation, Timeouts, Values](../08-context-cancellation-timeouts-values/README.md)

**Interview weight: 1.0%**  ·  Revise in: 5 min  ·  Run: `go run ./02-concurrency/07-atomic-operations`

## In one line
Atomics make a single read, write or read-modify-write on one variable indivisible without a lock, which is ideal for counters, flags and swapping an immutable pointer.

## Must know
- `x++` is load, add, store, so it is a data race across goroutines. `atomic.Int64.Add(1)` is safe.
- Prefer typed atomics (Go 1.19+): `atomic.Int32`, `Int64`, `Uint64`, `Bool`, `Pointer[T]`. `atomic.Value` is the older untyped one.
- Operations: `Load`, `Store`, `Add`, `Swap`, `CompareAndSwap` (CAS).
- CAS loop: load, compute, `CompareAndSwap(old, new)`, retry if it failed.
- `atomic.Pointer[T]` plus copy-on-write (never edit, always replace): build a new config, store the pointer, readers `Load` a consistent snapshot.
- Atomics protect one variable only. Two related fields need a mutex (or one pointer to an immutable struct).
- Atomics also order memory: what a goroutine wrote before an atomic `Store` is visible to one that sees that value via `Load`.
- Typed atomics must not be copied after first use.

## Code to remember
```go
var hits atomic.Int64
hits.Add(1)
n := hits.Load()

var cfg atomic.Pointer[Config]
cfg.Store(&Config{Version: 1})
c := cfg.Load() // never modify *c, publish a new one instead

var closed atomic.Bool
if closed.CompareAndSwap(false, true) {
    // exactly one caller runs this
}
```

## Common traps
- Mixing atomic and plain access to the same variable (still a race).
- Using atomics for multi-field invariants, such as "balance and count must change together".
- Writing a spin lock from CAS. It wastes CPU. Use `sync.Mutex`.
- Mutating the struct behind an `atomic.Pointer` instead of replacing it.

## Interview Q&A

**Q1. Atomic vs Mutex?**

**A:** Atomic for one simple variable (counter, flag, pointer swap), faster and no blocking. Mutex when several fields or a multi-step operation must be consistent.

**Q2. What is compare-and-swap?**

**A:** Atomically "if the value is still `old`, set it to `new` and return true". It is the building block for lock-free updates and one-time flags.

**Q3. How do you hot-reload config without locking readers?**

**A:** Keep it in an `atomic.Pointer[Config]`. The reloader builds a fresh Config and calls `Store`. Readers call `Load` once and use that snapshot.

**Q4. Why can `counter++` be wrong with goroutines?**

**A:** It is three steps, load, add, store, and two goroutines can interleave and lose an update. Use a mutex or `atomic.Add`.
