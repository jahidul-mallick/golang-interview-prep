// Sliding window: fixed size and variable size.
// Run: go run ./04-dsa/03-sliding-window
package main

import "fmt"

// maxSumK: fixed window of size k. Add the new element, drop the old one. O(n).
func maxSumK(nums []int, k int) int {
	sum := 0
	for i := 0; i < k; i++ {
		sum += nums[i]
	}
	best := sum
	for i := k; i < len(nums); i++ {
		sum += nums[i] - nums[i-k]
		best = max(best, sum)
	}
	return best
}

// longestUnique: longest substring without repeats. The window grows on the
// right and its left edge jumps past the previous copy of a repeated byte.
func longestUnique(s string) int {
	last := map[byte]int{} // byte -> last index seen
	best, left := 0, 0
	for right := 0; right < len(s); right++ {
		if i, ok := last[s[right]]; ok && i >= left {
			left = i + 1
		}
		last[s[right]] = right
		best = max(best, right-left+1)
	}
	return best
}

// minSubArrayLen: shortest subarray with sum >= target (positive numbers).
// Grow right until valid, then shrink left while still valid. O(n).
func minSubArrayLen(target int, nums []int) int {
	best, sum, left := 0, 0, 0
	for right, v := range nums {
		sum += v
		for sum >= target {
			if best == 0 || right-left+1 < best {
				best = right - left + 1
			}
			sum -= nums[left]
			left++
		}
	}
	return best
}

func main() {
	fmt.Println("=== max sum of window size k (fixed) ===")
	fmt.Println(maxSumK([]int{2, 1, 5, 1, 3, 2}, 3)) // 9

	fmt.Println("=== longest substring without repeating chars ===")
	fmt.Println(longestUnique("abcabcbb"), longestUnique("pwwkew")) // 3 3

	fmt.Println("=== minimum size subarray with sum >= 7 ===")
	fmt.Println(minSubArrayLen(7, []int{2, 3, 1, 2, 4, 3})) // 2 ([4 3])
}
