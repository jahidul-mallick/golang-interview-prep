# Iterators & range-over-func (Go 1.23+)

[← Section index](../README.md) · Prev: [13 · Testing, Benchmarks & Fuzzing](../13-testing-benchmarks-fuzzing/README.md) · Next: [JSON, Reflection & Struct Tags](../15-json-reflection-struct-tags/README.md)

**Interview weight: 0.3%**  ·  Revise in: 3 min  ·  Run: `go run ./01-go-fundamentals/14-iterators-range-over-func`
Optional / low priority

## In one line
Since Go 1.23 you can `range` over a function of type `iter.Seq`, which pushes values to the loop body through a `yield` callback.

## Must know
- `for i := range 5` ranges over an integer (Go 1.22).
- `iter.Seq[V]` is `func(yield func(V) bool)`; `iter.Seq2[K, V]` yields pairs.
- The iterator must stop and return as soon as `yield` returns false (the loop did `break`).
- `slices.All`, `slices.Values`, `maps.Keys` return iterators; `slices.Collect` and `slices.Sorted` consume them.

## Code to remember
```go
func Count(n int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := range n {
			if !yield(i) {
				return
			}
		}
	}
}
```

## Interview Q&A

**Q1. Why use an iterator instead of returning a slice?**

**A:** It is lazy (values are produced one at a time, on demand): no full slice in memory, and it can be infinite or stream data.

**Q2. What happens on `break` inside the loop?**

**A:** `yield` returns false and the iterator must return.

**Q3. Iterator or channel?**

**A:** An iterator runs in the caller's goroutine: no extra goroutine, no leak if the loop stops early. Use a channel only for real concurrency.
