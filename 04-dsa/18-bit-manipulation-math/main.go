// Bit tricks and math: single number, count set bits, power of two, gcd.
// Run: go run ./04-dsa/18-bit-manipulation-math
package main

import "fmt"

// singleNumber: x ^ x = 0 and x ^ 0 = x, so pairs cancel out.
func singleNumber(nums []int) int {
	res := 0
	for _, n := range nums {
		res ^= n
	}
	return res
}

// countBits: n & (n-1) clears the lowest set bit.
func countBits(n uint) int {
	count := 0
	for ; n != 0; n &= n - 1 {
		count++
	}
	return count
}

func isPowerOfTwo(n int) bool { return n > 0 && n&(n-1) == 0 }

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func main() {
	fmt.Println("=== single number [4 1 2 1 2] ===")
	fmt.Println(singleNumber([]int{4, 1, 2, 1, 2}))
	fmt.Println("=== count set bits of 13 (1101) ===")
	fmt.Println(countBits(13))
	fmt.Println("=== power of two: 16, 18 ===")
	fmt.Println(isPowerOfTwo(16), isPowerOfTwo(18))
	fmt.Println("=== gcd(48, 18) ===")
	fmt.Println(gcd(48, 18))
}
