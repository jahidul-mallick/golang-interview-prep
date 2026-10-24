# 01 · Big-O and the Go Toolkit for DSA

[← Section index](../README.md) · Next: [02 · Arrays, Two Pointers, Prefix Sum](../02-arrays-two-pointers-prefix-sum/README.md)

**Interview weight: 0.8%**  ·  Revise in: 8 min  ·  Run: `go run ./04-dsa/01-big-o-and-go-toolkit-for-dsa`

## In one line
Read the constraints first to pick the complexity class, then code with the standard Go tools (slices, maps, container/heap) without hesitation.

## Must know
- Common classes, best to worst: O(1), O(log n), O(n), O(n log n), O(n^2), O(2^n), O(n!).
- Rule of thumb (about 10^8 simple ops per second): n up to 10^5 needs O(n log n) or better; n up to 10^3 allows O(n^2); n up to 20 allows O(2^n).
- Say complexity for both time and space. Recursion uses stack space proportional to its depth.
- `append` is amortized O(1): capacity grows geometrically, so n appends cost O(n) overall.
- Trade memory for time: a hash map turns an O(n^2) pair search into O(n).
- Sorting with `slices.Sort` / `slices.SortFunc` is O(n log n). `sort.Slice` is not stable, `sort.SliceStable` is.
- Set is `map[T]struct{}`. Stack is a slice with `append` and `s[:len(s)-1]`. Queue is a slice plus `q = q[1:]`.
- `container/heap` needs 5 methods (`Len`, `Less`, `Swap`, `Push`, `Pop`); `Push` and `Pop` use a pointer receiver.

## Code to remember
```go
type IntHeap []int // min-heap; flip Less to > for a max-heap

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *IntHeap) Pop() any {
	old := *h
	x := old[len(old)-1]
	*h = old[:len(old)-1]
	return x
}
```

## Common traps
- `s += t` in a loop is O(n^2); use `strings.Builder`.
- `len(s)` on a string counts bytes, not characters; use `[]rune(s)` when you need characters.
- `lo + hi` can overflow in other languages; write `lo + (hi-lo)/2` out of habit.
- Sub-slices share the backing array: use `slices.Clone` when you need an independent copy.

## Interview Q&A

**Q1. What is the complexity of appending to a slice?**

**A:** Amortized O(1). Occasionally it reallocates and copies, but growth is geometric, so the average per append stays constant.

**Q2. How do you decide the approach from the constraints?**

**A:** n around 10^5 means O(n log n) or O(n); n around 20 hints at bitmask or backtracking; n around 10^3 allows O(n^2).

**Q3. Time vs space trade-off example?**

**A:** Two sum: brute force is O(n^2) time and O(1) space; a hash map gives O(n) time and O(n) space.

**Q4. Why is `map[T]struct{}` used for sets?**

**A:** `struct{}` takes zero bytes, so you store only the keys.
