// Binary search: classic, boundary, rotated array, search on the answer.
// Run: go run ./04-dsa/08-binary-search
package main

import "fmt"

// search is the classic version. mid = lo + (hi-lo)/2 avoids int overflow.
func search(nums []int, target int) int {
	lo, hi := 0, len(nums)-1
	for lo <= hi {
		mid := lo + (hi-lo)/2
		switch {
		case nums[mid] == target:
			return mid
		case nums[mid] < target:
			lo = mid + 1
		default:
			hi = mid - 1
		}
	}
	return -1
}

// lowerBound returns the first index with nums[i] >= target (len if none).
// Half-open [lo, hi): this is the template that handles most variants.
func lowerBound(nums []int, target int) int {
	lo, hi := 0, len(nums)
	for lo < hi {
		mid := lo + (hi-lo)/2
		if nums[mid] < target {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	return lo
}

// searchRotated: in a rotated sorted array one half is always sorted.
// Check whether the target lies in the sorted half, else go the other way.
func searchRotated(nums []int, target int) int {
	lo, hi := 0, len(nums)-1
	for lo <= hi {
		mid := lo + (hi-lo)/2
		if nums[mid] == target {
			return mid
		}
		if nums[lo] <= nums[mid] { // left half sorted
			if nums[lo] <= target && target < nums[mid] {
				hi = mid - 1
			} else {
				lo = mid + 1
			}
		} else { // right half sorted
			if nums[mid] < target && target <= nums[hi] {
				lo = mid + 1
			} else {
				hi = mid - 1
			}
		}
	}
	return -1
}

// minEatingSpeed (Koko): binary search on the answer. "Can she finish at
// speed k within h hours" is monotone in k, so find the smallest true k.
func minEatingSpeed(piles []int, h int) int {
	lo, hi := 1, 0
	for _, p := range piles {
		hi = max(hi, p)
	}
	for lo < hi {
		mid := lo + (hi-lo)/2
		hours := 0
		for _, p := range piles {
			hours += (p + mid - 1) / mid // ceil(p/mid)
		}
		if hours <= h {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return lo
}

func main() {
	nums := []int{1, 3, 5, 7, 9, 11}
	fmt.Println("=== classic search ===")
	fmt.Println(search(nums, 7), search(nums, 4))

	fmt.Println("=== lower bound (first >= target) ===")
	fmt.Println(lowerBound(nums, 6), lowerBound(nums, 100)) // 3 6

	fmt.Println("=== first and last position (lowerBound twice) ===")
	dup := []int{5, 7, 7, 8, 8, 10}
	fmt.Println(lowerBound(dup, 8), lowerBound(dup, 9)-1) // 3 4

	fmt.Println("=== search in rotated sorted array ===")
	fmt.Println(searchRotated([]int{4, 5, 6, 7, 0, 1, 2}, 0)) // 4

	fmt.Println("=== binary search on the answer (Koko eating bananas) ===")
	fmt.Println(minEatingSpeed([]int{3, 6, 7, 11}, 8)) // 4
}
