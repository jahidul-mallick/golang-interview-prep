# 12 · Heap / Priority Queue / Top-K

[← Section index](../README.md) · Prev: [11 · Binary Tree & BST](../11-binary-tree-bst/README.md) · Next: [13 · Graphs: BFS, DFS, Topological Sort, Union-Find, Dijkstra](../13-graphs-bfs-dfs-topological-sort-union-find/README.md)

**Interview weight: 0.9%**  ·  Revise in: 8 min  ·  Run: `go run ./04-dsa/12-heap-priority-queue-topk`

## In one line
"Kth largest", "top K", "merge K sorted", "always need the smallest/largest next": use a heap of size K.

## Must know
- Go has no ready heap type. Implement `Len, Less, Swap, Push, Pop` on a slice and use `container/heap`.
- Call `heap.Push(h, x)` and `heap.Pop(h)`, never `h.Push` / `h.Pop` directly.
- Top K largest: keep a MIN-heap of size K. The root is the Kth largest.
- Top K smallest: keep a MAX-heap of size K (flip `Less`).
- Complexity: push/pop O(log k), whole pass O(n log k), space O(k).
- `heap.Init` builds a heap from a slice in O(n).
- `heap.Fix(h, i)` re-sorts one element after its priority changes (used in Dijkstra-style code).
- Sorting everything is O(n log n); the heap wins when k is much smaller than n.

## Code to remember
```go
type IntHeap []int // min-heap
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
- Using `h.Pop()` instead of `heap.Pop(h)` (skips the re-heapify).
- Wrong heap direction: min-heap for top K largest, not max-heap.
- `Push`/`Pop` need pointer receivers; `Len/Less/Swap` can be value receivers.
- Ranging over a map gives random order, so sort the result if the output must be stable.

## Interview Q&A

**Q1. Kth largest element, best approach?**

**A:** Min-heap of size k, O(n log k). Quickselect is O(n) average but O(n^2) worst case.

**Q2. Why a min-heap for the largest K?**

**A:** The root is the smallest of the K kept, so it is the one to evict when a bigger value arrives.

**Q3. Time to build a heap from n items?**

**A:** O(n) with `heap.Init`, versus O(n log n) for n pushes.

**Q4. Merge K sorted lists?**

**A:** Push each list head into a min-heap, pop the smallest, push its next node. O(N log k).

**Q5. Top K frequent elements?**

**A:** Count with a map, then keep a min-heap of size k ordered by count. O(n log k). Bucket sort by count gives O(n).
