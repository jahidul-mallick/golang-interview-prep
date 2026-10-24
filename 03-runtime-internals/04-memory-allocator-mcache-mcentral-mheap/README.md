# 04 · The Memory Allocator: mcache → mcentral → mheap

[← Section index](../README.md) · Prev: [03 · Stack vs Heap & Escape Analysis](../03-stack-heap-escape-analysis/README.md) · Next: [Garbage Collector](../05-garbage-collector/README.md)

**Interview weight: 0.6%**  ·  Revise in: 4 min  ·  Run: `go run ./03-runtime-internals/04-memory-allocator-mcache-mcentral-mheap`

## In one line
Go allocates from size-class buckets: a per-P cache first (no lock), then a shared per-class list, then the page heap.

## Must know
- Small objects (up to 32 KiB) are rounded up to a fixed size class. A 25-byte request uses a 32-byte slot.
- Lookup order: `mcache` (one per P, no lock), then `mcentral` (shared, one per size class), then `mheap` (pages from the OS).
- Memory is grouped in spans: runs of pages holding objects of one size class.
- Large objects (over 32 KiB) go straight to `mheap`. Tiny pointer-free objects (under 16 bytes) are packed together.
- There is no `free()`: the GC sweeps unused objects, and idle pages return to the OS lazily.

## Common traps
- Many tiny allocations cost CPU and GC work. Pre-size slices and reuse buffers.
- A big `[]*T` is expensive for the GC to scan, a `[]byte` of the same size is not.

## Interview Q&A

**Q1. Why is Go allocation fast?**

**A:** Each P has its own cache of spans per size class, so most allocations need no lock.

**Q2. What are `mcache`, `mcentral` and `mheap`?**

**A:** `mcache` is the per-P local cache, `mcentral` is the shared pool of spans for one size class, `mheap` manages pages obtained from the OS.

**Q3. Why can memory (RSS) stay high after the GC runs?**

**A:** The GC frees objects inside spans, but pages go back to the OS lazily, and the heap may grow before the next cycle.
