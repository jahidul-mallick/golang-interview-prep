// Backtracking: subsets, permutations, combination sum.
// Run: go run ./04-dsa/10-recursion-backtracking
package main

import "fmt"

func subsets(nums []int) [][]int {
	var res [][]int
	var path []int
	var dfs func(start int)
	dfs = func(start int) {
		res = append(res, append([]int{}, path...)) // copy: path is reused
		for i := start; i < len(nums); i++ {
			path = append(path, nums[i]) // choose
			dfs(i + 1)                   // explore
			path = path[:len(path)-1]    // un-choose
		}
	}
	dfs(0)
	return res
}

func permutations(nums []int) [][]int {
	var res [][]int
	var path []int
	used := make([]bool, len(nums))
	var dfs func()
	dfs = func() {
		if len(path) == len(nums) {
			res = append(res, append([]int{}, path...))
			return
		}
		for i, n := range nums {
			if used[i] {
				continue
			}
			used[i] = true
			path = append(path, n)
			dfs()
			path = path[:len(path)-1]
			used[i] = false
		}
	}
	dfs()
	return res
}

// combinationSum: candidates may be reused; passing i (not i+1) allows reuse.
func combinationSum(cands []int, target int) [][]int {
	var res [][]int
	var path []int
	var dfs func(start, remain int)
	dfs = func(start, remain int) {
		if remain == 0 {
			res = append(res, append([]int{}, path...))
			return
		}
		for i := start; i < len(cands) && cands[i] <= remain; i++ { // cands sorted
			path = append(path, cands[i])
			dfs(i, remain-cands[i])
			path = path[:len(path)-1]
		}
	}
	dfs(0, target)
	return res
}

func main() {
	fmt.Println("=== subsets [1 2 3] ===")
	fmt.Println(subsets([]int{1, 2, 3}))
	fmt.Println("=== permutations [1 2 3] ===")
	fmt.Println(permutations([]int{1, 2, 3}))
	fmt.Println("=== combination sum [2 3 6 7], target 7 ===")
	fmt.Println(combinationSum([]int{2, 3, 6, 7}, 7))
}
