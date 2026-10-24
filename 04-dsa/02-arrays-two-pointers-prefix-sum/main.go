// Two pointers and prefix sum on arrays.
// Run: go run ./04-dsa/02-arrays-two-pointers-prefix-sum
package main

import "fmt"

// twoSumSorted: sorted input, one pointer at each end. O(n) time, O(1) space.
func twoSumSorted(nums []int, target int) (int, int) {
	l, r := 0, len(nums)-1
	for l < r {
		switch sum := nums[l] + nums[r]; {
		case sum == target:
			return l, r
		case sum < target:
			l++
		default:
			r--
		}
	}
	return -1, -1
}

// maxArea (container with most water): move the shorter wall, because moving
// the taller one can never increase the area. O(n).
func maxArea(h []int) int {
	best, l, r := 0, 0, len(h)-1
	for l < r {
		best = max(best, min(h[l], h[r])*(r-l))
		if h[l] < h[r] {
			l++
		} else {
			r--
		}
	}
	return best
}

// moveZeroes: the write pointer marks the next slot for a non-zero. In place, O(n).
func moveZeroes(nums []int) {
	w := 0
	for _, v := range nums {
		if v != 0 {
			nums[w] = v
			w++
		}
	}
	for ; w < len(nums); w++ {
		nums[w] = 0
	}
}

// subarraySumK counts subarrays summing to k with a prefix-sum map.
// sum(i..j) = prefix[j] - prefix[i-1], so look for (prefix - k) seen earlier. O(n).
func subarraySumK(nums []int, k int) int {
	seen := map[int]int{0: 1} // the empty prefix
	prefix, count := 0, 0
	for _, v := range nums {
		prefix += v
		count += seen[prefix-k]
		seen[prefix]++
	}
	return count
}

func main() {
	fmt.Println("=== two sum, sorted array ===")
	i, j := twoSumSorted([]int{1, 2, 4, 7, 11, 15}, 15)
	fmt.Println("indices:", i, j) // 2 and 4: 4+11

	fmt.Println("=== container with most water ===")
	fmt.Println(maxArea([]int{1, 8, 6, 2, 5, 4, 8, 3, 7})) // 49

	fmt.Println("=== move zeroes in place ===")
	z := []int{0, 1, 0, 3, 12}
	moveZeroes(z)
	fmt.Println(z)

	fmt.Println("=== subarray sum equals k (prefix sum + map) ===")
	fmt.Println(subarraySumK([]int{1, 2, 3, -2, 2}, 3)) // [1 2] [3] [2 3 -2] [3 -2 2]
}
