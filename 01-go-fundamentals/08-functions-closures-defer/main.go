// Shows multiple returns, closures, defer order, defer argument timing and named returns.
// Run: go run ./01-go-fundamentals/08-functions-closures-defer
package main

import (
	"errors"
	"fmt"
)

func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("divide by zero")
	}
	return a / b, nil
}

func counter() func() int {
	n := 0
	return func() int { n++; return n } // n lives on after counter returns
}

func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// A deferred closure can change a named result.
func safeDivide(a, b int) (result int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered: %v", r)
		}
	}()
	return a / b, nil
}

func deferOrder() {
	for i := 1; i <= 3; i++ {
		defer fmt.Println("deferred", i) // runs last-in, first-out
	}
}

func deferArgs() {
	x := 1
	defer fmt.Println("defer saw x =", x) // arguments are evaluated NOW
	x = 2
	fmt.Println("x is", x)
}

func main() {
	fmt.Println("=== multiple returns ===")
	fmt.Println(divide(7, 2))
	fmt.Println(divide(1, 0))

	fmt.Println("=== closures keep state ===")
	next := counter()
	fmt.Println(next(), next(), next())
	other := counter() // separate state
	fmt.Println(other())

	fmt.Println("=== loop variable (Go 1.22+: one per iteration) ===")
	var fns []func() int
	for i := 0; i < 3; i++ {
		fns = append(fns, func() int { return i })
	}
	for _, f := range fns {
		fmt.Print(f(), " ")
	}
	fmt.Println()

	fmt.Println("=== variadic ===")
	fmt.Println(sum(), sum(1, 2, 3), sum([]int{4, 5}...))

	fmt.Println("=== defer order and argument timing ===")
	deferOrder()
	deferArgs()

	fmt.Println("=== named result changed by defer ===")
	fmt.Println(safeDivide(6, 3))
	fmt.Println(safeDivide(1, 0))
}
