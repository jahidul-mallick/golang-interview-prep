// Sorting: quick sort, merge sort, and the standard library.
// Run: go run ./04-dsa/09-sorting-algorithms
package main

import (
	"fmt"
	"slices"
	"sort"
)

// quickSort sorts in place. Lomuto partition with the last element as pivot.
// Average O(n log n), worst O(n^2) (already sorted input with this pivot).
func quickSort(a []int) {
	if len(a) < 2 {
		return
	}
	pivot := a[len(a)-1]
	i := 0
	for j := 0; j < len(a)-1; j++ {
		if a[j] < pivot {
			a[i], a[j] = a[j], a[i]
			i++
		}
	}
	a[i], a[len(a)-1] = a[len(a)-1], a[i]
	quickSort(a[:i])
	quickSort(a[i+1:])
}

// mergeSort returns a new sorted slice. Stable, always O(n log n), O(n) space.
func mergeSort(a []int) []int {
	if len(a) < 2 {
		return a
	}
	left, right := mergeSort(a[:len(a)/2]), mergeSort(a[len(a)/2:])
	out := make([]int, 0, len(a))
	i, j := 0, 0
	for i < len(left) && j < len(right) {
		if left[i] <= right[j] { // <= keeps equal elements in order (stable)
			out = append(out, left[i])
			i++
		} else {
			out = append(out, right[j])
			j++
		}
	}
	out = append(out, left[i:]...)
	return append(out, right[j:]...)
}

type person struct {
	name string
	age  int
}

func main() {
	fmt.Println("=== quick sort ===")
	q := []int{5, 2, 9, 1, 5, 6}
	quickSort(q)
	fmt.Println(q)

	fmt.Println("=== merge sort ===")
	fmt.Println(mergeSort([]int{5, 2, 9, 1, 5, 6}))

	fmt.Println("=== standard library ===")
	s := []int{5, 2, 9, 1}
	slices.Sort(s)
	fmt.Println("slices.Sort:", s)

	people := []person{{"A", 30}, {"B", 25}, {"C", 30}, {"D", 25}}
	// SliceStable keeps equal ages in original order; plain sort.Slice does not promise that.
	sort.SliceStable(people, func(i, j int) bool { return people[i].age < people[j].age })
	fmt.Println("stable by age:", people)
}
