// Big-O in practice plus the Go toolkit every DSA answer needs.
// Run: go run ./04-dsa/01-big-o-and-go-toolkit-for-dsa
package main

import (
	"cmp"
	"container/heap"
	"fmt"
	"slices"
	"strings"
)

// countPairsBrute is O(n^2) time, O(1) space.
func countPairsBrute(nums []int, target int) int {
	count := 0
	for i := range nums {
		for j := i + 1; j < len(nums); j++ {
			if nums[i]+nums[j] == target {
				count++
			}
		}
	}
	return count
}

// countPairsHash is O(n) time, O(n) space: trade memory for time.
func countPairsHash(nums []int, target int) int {
	seen := map[int]int{}
	count := 0
	for _, x := range nums {
		count += seen[target-x] // missing key gives 0, no ok-check needed
		seen[x]++
	}
	return count
}

// IntHeap is a min-heap. Push and Pop need a pointer receiver.
type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] } // use > for a max-heap
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *IntHeap) Pop() any {
	old := *h
	x := old[len(old)-1]
	*h = old[:len(old)-1]
	return x
}

// Set is the idiomatic Go set: struct{} values take no space.
type Set[T comparable] map[T]struct{}

func (s Set[T]) Add(v T)      { s[v] = struct{}{} }
func (s Set[T]) Has(v T) bool { _, ok := s[v]; return ok }

func main() {
	fmt.Println("=== O(n^2) vs O(n): count pairs summing to 6 ===")
	nums := []int{1, 5, 7, -1, 5}
	fmt.Println("brute:", countPairsBrute(nums, 6), "hash:", countPairsHash(nums, 6))

	fmt.Println("=== append is amortized O(1): capacity grows geometrically ===")
	var s []int
	var caps []int
	for i := 0; i < 20; i++ {
		s = append(s, i)
		if len(caps) == 0 || cap(s) != caps[len(caps)-1] {
			caps = append(caps, cap(s))
		}
	}
	fmt.Println("capacities seen:", caps)

	fmt.Println("=== sort with multiple keys (cmp.Or) ===")
	words := []string{"go", "rust", "c", "zig", "ada"}
	slices.SortFunc(words, func(a, b string) int {
		return cmp.Or(cmp.Compare(len(a), len(b)), cmp.Compare(a, b))
	})
	fmt.Println(words)

	fmt.Println("=== container/heap: 3 smallest ===")
	h := &IntHeap{7, 2, 9, 4, 1}
	heap.Init(h)
	for range 3 {
		fmt.Print(heap.Pop(h), " ")
	}
	fmt.Println()

	fmt.Println("=== stack, set, strings.Builder ===")
	stack := []int{1, 2, 3}
	top := stack[len(stack)-1]
	stack = stack[:len(stack)-1]
	set := Set[string]{}
	set.Add("a")
	var sb strings.Builder // avoids O(n^2) from s += t in a loop
	for range 3 {
		sb.WriteString("x")
	}
	fmt.Println("popped:", top, "left:", stack, "has a:", set.Has("a"), "has b:", set.Has("b"), "built:", sb.String())
}
