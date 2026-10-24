// Dynamic programming: climbing stairs, house robber, coin change, LCS.
// Run: go run ./04-dsa/14-dynamic-programming
package main

import "fmt"

// climbStairs: dp[i] = dp[i-1] + dp[i-2]; only the last two values are kept.
func climbStairs(n int) int {
	a, b := 1, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

// rob: best loot without taking two adjacent houses.
func rob(nums []int) int {
	prev, cur := 0, 0
	for _, n := range nums {
		prev, cur = cur, max(cur, prev+n)
	}
	return cur
}

// coinChange: fewest coins for amount, -1 if impossible. dp[a] = 1 + min(dp[a-coin]).
func coinChange(coins []int, amount int) int {
	dp := make([]int, amount+1)
	for a := 1; a <= amount; a++ {
		dp[a] = amount + 1 // "infinity"
		for _, c := range coins {
			if c <= a {
				dp[a] = min(dp[a], dp[a-c]+1)
			}
		}
	}
	if dp[amount] > amount {
		return -1
	}
	return dp[amount]
}

// lcs: longest common subsequence length with a 2D table.
func lcs(a, b string) int {
	dp := make([][]int, len(a)+1)
	for i := range dp {
		dp[i] = make([]int, len(b)+1)
	}
	for i := 1; i <= len(a); i++ {
		for j := 1; j <= len(b); j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
			}
		}
	}
	return dp[len(a)][len(b)]
}

func main() {
	fmt.Println("=== climb stairs (n=5) ===")
	fmt.Println(climbStairs(5))
	fmt.Println("=== house robber [2 7 9 3 1] ===")
	fmt.Println(rob([]int{2, 7, 9, 3, 1}))
	fmt.Println("=== coin change [1 2 5], amount 11 ===")
	fmt.Println(coinChange([]int{1, 2, 5}, 11), coinChange([]int{2}, 3))
	fmt.Println("=== LCS abcde / ace ===")
	fmt.Println(lcs("abcde", "ace"))
}
