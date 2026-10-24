// Shows generic functions, constraints, a generic type and stdlib helpers.
// Run: go run ./01-go-fundamentals/11-generics
package main

import (
	"cmp"
	"fmt"
	"maps"
	"slices"
	"strconv"
)

type Number interface{ ~int | ~int64 | ~float64 }

type Celsius float64

func Sum[T Number](xs []T) T {
	var total T // zero value of T
	for _, x := range xs {
		total += x
	}
	return total
}

func Map[T, U any](in []T, f func(T) U) []U {
	out := make([]U, 0, len(in))
	for _, v := range in {
		out = append(out, f(v))
	}
	return out
}

func Filter[T any](in []T, keep func(T) bool) []T {
	var out []T
	for _, v := range in {
		if keep(v) {
			out = append(out, v)
		}
	}
	return out
}

func Max[T cmp.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func Keys[K comparable, V any](m map[K]V) []K {
	out := make([]K, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

func Zero[T any]() T {
	var z T
	return z
}

type Stack[T any] struct{ items []T }

func (s *Stack[T]) Push(v T) { s.items = append(s.items, v) }

func (s *Stack[T]) Pop() (T, bool) {
	var zero T
	if len(s.items) == 0 {
		return zero, false
	}
	v := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return v, true
}

func main() {
	fmt.Println("=== constraints: Sum works for named types because of ~ ===")
	fmt.Println(Sum([]int{1, 2, 3}), Sum([]float64{1.5, 2.5}), Sum([]Celsius{20, 1.5}))

	fmt.Println("=== Map and Filter, types inferred ===")
	nums := []int{1, 2, 3, 4}
	fmt.Println(Map(nums, strconv.Itoa), Filter(nums, func(n int) bool { return n%2 == 0 }))

	fmt.Println("=== cmp.Ordered and comparable ===")
	fmt.Println(Max(3, 7), Max("go", "rust"))
	keys := Keys(map[string]int{"b": 2, "a": 1})
	slices.Sort(keys)
	fmt.Println(keys, slices.Sorted(maps.Keys(map[int]bool{3: true, 1: true})))

	fmt.Println("=== explicit type argument when T is only in the result ===")
	fmt.Printf("%q %v\n", Zero[string](), Zero[*int]())

	fmt.Println("=== generic type ===")
	var s Stack[string]
	s.Push("a")
	s.Push("b")
	v, ok := s.Pop()
	fmt.Println(v, ok)
	var empty Stack[int]
	v2, ok := empty.Pop()
	fmt.Println(v2, ok)
}
