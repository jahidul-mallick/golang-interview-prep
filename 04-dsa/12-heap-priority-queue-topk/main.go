// Heap problems with container/heap: kth largest, top-K frequent.
// Run: go run ./04-dsa/12-heap-priority-queue-topk
package main

import (
	"container/heap"
	"fmt"
	"sort"
)

// IntHeap is a min-heap. Push/Pop are for the heap package; use heap.Push/heap.Pop.
type IntHeap []int

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

// kthLargest keeps a min-heap of the k largest values; the root is the answer.
func kthLargest(nums []int, k int) int {
	h := &IntHeap{}
	for _, n := range nums {
		heap.Push(h, n)
		if h.Len() > k {
			heap.Pop(h) // drop the smallest
		}
	}
	return (*h)[0]
}

type pair struct{ val, count int }
type pairHeap []pair

func (h pairHeap) Len() int           { return len(h) }
func (h pairHeap) Less(i, j int) bool { return h[i].count < h[j].count }
func (h pairHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *pairHeap) Push(x any)        { *h = append(*h, x.(pair)) }
func (h *pairHeap) Pop() any {
	old := *h
	x := old[len(old)-1]
	*h = old[:len(old)-1]
	return x
}

func topKFrequent(nums []int, k int) []int {
	freq := map[int]int{}
	for _, n := range nums {
		freq[n]++
	}
	h := &pairHeap{}
	for v, c := range freq {
		heap.Push(h, pair{v, c})
		if h.Len() > k {
			heap.Pop(h) // drop the least frequent
		}
	}
	res := make([]int, 0, k)
	for h.Len() > 0 {
		res = append(res, heap.Pop(h).(pair).val)
	}
	sort.Ints(res) // fixed order so the output is stable
	return res
}

func main() {
	fmt.Println("=== kth largest (k=2) in [3 2 1 5 6 4] ===")
	fmt.Println(kthLargest([]int{3, 2, 1, 5, 6, 4}, 2))
	fmt.Println("=== top 2 frequent in [1 1 1 2 2 3] ===")
	fmt.Println(topKFrequent([]int{1, 1, 1, 2, 2, 3}, 2))
}
