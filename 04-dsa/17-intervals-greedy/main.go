// Intervals and greedy: merge intervals, min removals for no overlap, jump game.
// Run: go run ./04-dsa/17-intervals-greedy
package main

import (
	"fmt"
	"sort"
)

// merge: sort by start, then extend the last merged interval or start a new one.
func merge(iv [][2]int) [][2]int {
	if len(iv) == 0 {
		return nil
	}
	sort.Slice(iv, func(i, j int) bool { return iv[i][0] < iv[j][0] })
	out := [][2]int{iv[0]}
	for _, cur := range iv[1:] {
		last := &out[len(out)-1]
		if cur[0] <= last[1] {
			last[1] = max(last[1], cur[1])
		} else {
			out = append(out, cur)
		}
	}
	return out
}

// minRemovals: sort by END, keep every interval that starts after the last kept one ends.
func minRemovals(iv [][2]int) int {
	sort.Slice(iv, func(i, j int) bool { return iv[i][1] < iv[j][1] })
	kept, end := 0, -1<<31
	for _, cur := range iv {
		if cur[0] >= end {
			kept++
			end = cur[1]
		}
	}
	return len(iv) - kept
}

// canJump: track the farthest index reachable so far.
func canJump(nums []int) bool {
	far := 0
	for i, n := range nums {
		if i > far {
			return false
		}
		far = max(far, i+n)
	}
	return true
}

func main() {
	fmt.Println("=== merge intervals ===")
	fmt.Println(merge([][2]int{{1, 3}, {8, 10}, {2, 6}, {15, 18}, {17, 20}}))
	fmt.Println("=== min removals for no overlap ===")
	fmt.Println(minRemovals([][2]int{{1, 2}, {2, 3}, {3, 4}, {1, 3}}))
	fmt.Println("=== jump game ===")
	fmt.Println(canJump([]int{2, 3, 1, 1, 4}), canJump([]int{3, 2, 1, 0, 4}))
}
