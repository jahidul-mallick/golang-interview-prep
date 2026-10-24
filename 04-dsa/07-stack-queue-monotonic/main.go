// Stack, queue and the monotonic stack / deque.
// Run: go run ./04-dsa/07-stack-queue-monotonic
package main

import "fmt"

// validParentheses: push openers, each closer must match the top. O(n).
func validParentheses(s string) bool {
	pairs := map[rune]rune{')': '(', ']': '[', '}': '{'}
	var stack []rune
	for _, c := range s {
		if open, isCloser := pairs[c]; isCloser {
			if len(stack) == 0 || stack[len(stack)-1] != open {
				return false
			}
			stack = stack[:len(stack)-1]
		} else {
			stack = append(stack, c)
		}
	}
	return len(stack) == 0
}

// dailyTemperatures: monotonic decreasing stack of indices. A warmer day
// resolves every colder day waiting on the stack. Each index is pushed and
// popped once, so O(n).
func dailyTemperatures(temps []int) []int {
	res := make([]int, len(temps))
	var stack []int // indices with no warmer day found yet
	for i, t := range temps {
		for len(stack) > 0 && temps[stack[len(stack)-1]] < t {
			j := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			res[j] = i - j
		}
		stack = append(stack, i)
	}
	return res
}

// maxSlidingWindow: a deque of indices whose values decrease, so the front
// is always the window maximum. O(n). Popping the front with dq = dq[1:] is
// fine for interviews.
func maxSlidingWindow(nums []int, k int) []int {
	var res, dq []int
	for i, v := range nums {
		for len(dq) > 0 && nums[dq[len(dq)-1]] <= v {
			dq = dq[:len(dq)-1] // smaller values can never be the max again
		}
		dq = append(dq, i)
		if dq[0] <= i-k { // front left the window
			dq = dq[1:]
		}
		if i >= k-1 {
			res = append(res, nums[dq[0]])
		}
	}
	return res
}

// MyQueue is a FIFO built from two stacks. Amortized O(1) per operation.
type MyQueue struct{ in, out []int }

func (q *MyQueue) Push(x int) { q.in = append(q.in, x) }

func (q *MyQueue) Pop() int {
	if len(q.out) == 0 { // refill only when empty, to keep the order right
		for len(q.in) > 0 {
			q.out = append(q.out, q.in[len(q.in)-1])
			q.in = q.in[:len(q.in)-1]
		}
	}
	x := q.out[len(q.out)-1]
	q.out = q.out[:len(q.out)-1]
	return x
}

func main() {
	fmt.Println("=== valid parentheses ===")
	fmt.Println(validParentheses("{[()]}"), validParentheses("(]"))

	fmt.Println("=== daily temperatures (monotonic stack) ===")
	fmt.Println(dailyTemperatures([]int{73, 74, 75, 71, 69, 72, 76, 73}))

	fmt.Println("=== sliding window maximum (monotonic deque) ===")
	fmt.Println(maxSlidingWindow([]int{1, 3, -1, -3, 5, 3, 6, 7}, 3))

	fmt.Println("=== queue from two stacks ===")
	q := &MyQueue{}
	q.Push(1)
	q.Push(2)
	fmt.Print(q.Pop(), " ")
	q.Push(3)
	fmt.Println(q.Pop(), q.Pop())
}
