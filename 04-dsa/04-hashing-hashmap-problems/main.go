// Classic hash map problems.
// Run: go run ./04-dsa/04-hashing-hashmap-problems
package main

import (
	"fmt"
	"slices"
)

// twoSum: store value -> index, look up the complement. O(n) time and space.
func twoSum(nums []int, target int) []int {
	idx := map[int]int{}
	for i, v := range nums {
		if j, ok := idx[target-v]; ok {
			return []int{j, i}
		}
		idx[v] = i
	}
	return nil
}

// groupAnagrams: a [26]int letter count is a comparable array, so it works
// directly as a map key (slices cannot be map keys).
func groupAnagrams(words []string) [][]string {
	groups := map[[26]int][]string{}
	var order [][26]int // keeps the output order stable (map order is random)
	for _, w := range words {
		var key [26]int
		for _, c := range w {
			key[c-'a']++
		}
		if _, ok := groups[key]; !ok {
			order = append(order, key)
		}
		groups[key] = append(groups[key], w)
	}
	var out [][]string
	for _, k := range order {
		out = append(out, groups[k])
	}
	return out
}

// longestConsecutive: only start counting at the beginning of a run
// (v-1 absent), so every number is visited at most twice. O(n).
func longestConsecutive(nums []int) int {
	set := map[int]struct{}{}
	for _, v := range nums {
		set[v] = struct{}{}
	}
	best := 0
	for v := range set {
		if _, ok := set[v-1]; ok {
			continue
		}
		n := 1
		for {
			if _, ok := set[v+n]; !ok {
				break
			}
			n++
		}
		best = max(best, n)
	}
	return best
}

// topKFrequent: count with a map, then sort keys by count. O(n log n).
// A heap or bucket sort gets it lower (see 12-heap-priority-queue-topk).
func topKFrequent(nums []int, k int) []int {
	count := map[int]int{}
	for _, v := range nums {
		count[v]++
	}
	keys := make([]int, 0, len(count))
	for v := range count {
		keys = append(keys, v)
	}
	slices.SortFunc(keys, func(a, b int) int { return count[b] - count[a] })
	return keys[:k]
}

func main() {
	fmt.Println("=== two sum ===")
	fmt.Println(twoSum([]int{2, 7, 11, 15}, 9))

	fmt.Println("=== group anagrams ===")
	fmt.Println(groupAnagrams([]string{"eat", "tea", "tan", "ate", "nat", "bat"}))

	fmt.Println("=== longest consecutive sequence ===")
	fmt.Println(longestConsecutive([]int{100, 4, 200, 1, 3, 2})) // 4

	fmt.Println("=== top k frequent ===")
	fmt.Println(topKFrequent([]int{1, 1, 1, 2, 2, 3}, 2)) // [1 2]
}
