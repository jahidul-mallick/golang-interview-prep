# 03 · Maps Internals

[← Section index](../README.md) · Prev: [02 · Arrays & Slices Internals](../02-arrays-slices-internals/README.md) · Next: [04 · Strings, Runes & Bytes](../04-strings-runes-bytes/README.md)

**Interview weight: 1.2%**  ·  Revise in: 6 min  ·  Run: `go run ./01-go-fundamentals/03-maps-internals`

## In one line
A map is a hash table with random iteration order that is not safe for concurrent writes, and writing to a nil map panics.

## Must know
- Keys must be comparable (`==` works): strings, numbers, bools, arrays, structs of those. Slices, maps and funcs cannot be keys.
- A map variable refers to the table. Assigning or passing it shares the same data.
- Reading a missing key returns the zero value. Use `v, ok := m[k]` to tell "missing" from "zero".
- A nil map can be read and ranged, but writing panics. Create it with `make` or a literal.
- Iteration order is randomised on purpose. Sort the keys if you need stable output.
- Deleting during `range` is safe. Inserting during `range` may or may not show up.
- Map elements are not addressable (no `&m[k]`): `m[k].field = v` does not compile. Copy, modify, store back, or store pointers.
- Concurrent write (or read plus write) is a fatal error, not a recoverable panic. Use `sync.Mutex`, `sync.RWMutex`, or `sync.Map` for special cases.

## Code to remember
```go
type SafeCounter struct {
	mu sync.Mutex
	m  map[string]int
}

func (c *SafeCounter) Inc(k string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[k]++
}
```

## Common traps
- Writing to a `var m map[string]int` that was never made.
- Relying on iteration order, or assuming a map is safe "because only one goroutine writes".
- Using a map of structs and trying to update a field in place.
- Forgetting that `make(map[K]V, n)` only pre-sizes, it does not set `len`.

## Interview Q&A

**Q1. How do you know if a key exists?**

**A:** `v, ok := m[k]`. `ok` is false when the key is missing, even if the zero value would look valid.

**Q2. Is a map safe for concurrent use?**

**A:** No. Concurrent writes crash the program with a fatal error. Guard with a mutex, or use `sync.Map` for read-heavy caches.

**Q3. How does a map work internally, and why is iteration order random?**

**A:** It is a hash table: the key's hash picks a slot group, and the table grows and moves entries when it gets full (Swiss tables since Go 1.24). The order is randomised on purpose so nobody relies on it.

**Q4. How do you implement a set?**

**A:** `map[T]struct{}`. The empty struct takes no space and you check membership with `_, ok := set[k]`.

**Q5. Why can I not update `m["a"].Score++`?**

**A:** Map entries can move when the table grows, so Go does not allow taking their address. Use `map[string]*User` or copy and reassign.

**Q6. How do you iterate a map in sorted key order?**

**A:** Collect and sort the keys, then look each up: `keys := slices.Sorted(maps.Keys(m))` (Go 1.23+), or `append` the keys to a slice and call `sort.Strings`.
