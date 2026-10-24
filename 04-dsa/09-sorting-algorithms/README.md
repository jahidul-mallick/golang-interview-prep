# Sorting Algorithms

[← Section index](../README.md) · Prev: [08 · Binary Search](../08-binary-search/README.md) · Next: [10 · Recursion & Backtracking](../10-recursion-backtracking/README.md)

**Interview weight: 0.8%**  ·  Revise in: 6 min  ·  Run: `go run ./04-dsa/09-sorting-algorithms`

## In one line
Know quick sort and merge sort well enough to write and compare them, and use the standard library in real code.

## Must know
- Quick sort: partition around a pivot, recurse. Average O(n log n), worst O(n^2), O(log n) extra space, not stable.
- Merge sort: split, sort halves, merge. Always O(n log n), O(n) extra space, stable.
- Heap sort: O(n log n), O(1) extra space, not stable.
- Counting sort: O(n + k) for integers in a small range k.
- Stable means equal elements keep their original order.
- Go: `slices.Sort` for ordered types, `slices.SortFunc` with a comparison returning negative, zero or positive, `sort.SliceStable` for stable sorting.
- Comparison sorts cannot beat O(n log n) in the general case.
- Merge sort suits linked lists (no random access needed).

## Code to remember
```go
func merge(l, r []int) []int { // the core of merge sort
	out := make([]int, 0, len(l)+len(r))
	for len(l) > 0 && len(r) > 0 {
		if l[0] <= r[0] { // <= keeps it stable
			out, l = append(out, l[0]), l[1:]
		} else {
			out, r = append(out, r[0]), r[1:]
		}
	}
	return append(append(out, l...), r...)
}
```

## Common traps
- Quick sort with a first or last pivot degrades to O(n^2) on sorted data (and, with Lomuto, on many duplicates); pick a random or median pivot.
- Assuming `sort.Slice` is stable: it is not.
- A comparison function that is not consistent (for example `<=` in `Less`) can misorder results.
- Writing a custom sort when the interviewer just wants you to call the library.

## Interview Q&A

**Q1. Quick sort versus merge sort?**

**A:** Quick sort is in place and usually faster in practice but has an O(n^2) worst case. Merge sort is stable with guaranteed O(n log n) but needs O(n) extra space.

**Q2. What is a stable sort and when does it matter?**

**A:** Equal keys keep their input order. It matters when sorting by several keys in passes, like sorting people by age after sorting by name.

**Q3. When would you use counting sort?**

**A:** When values are integers in a small range k much smaller than n. O(n + k) time with O(k) extra space.

**Q4. How do you sort a slice of structs in Go?**

**A:** `slices.SortFunc(s, func(a, b T) int { return cmp.Compare(a.Age, b.Age) })`, or `sort.SliceStable` when order of ties matters.
