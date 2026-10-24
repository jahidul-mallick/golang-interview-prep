// Shows slice header sharing, append growth, the aliasing bug, nil vs empty, and memory retention.
// Run: go run ./01-go-fundamentals/02-arrays-slices-internals
package main

import (
	"encoding/json"
	"fmt"
	"slices"
)

func appendInside(s []int) { s = append(s, 99) } // caller never sees this
func setFirst(s []int)     { s[0] = 99 }         // caller sees this

func main() {
	fmt.Println("=== arrays are values, slices share ===")
	a := [3]int{1, 2, 3}
	b := a // full copy
	b[0] = 100
	fmt.Println("array a:", a, "b:", b)
	s1 := []int{1, 2, 3}
	s2 := s1 // copies only the header
	s2[0] = 100
	fmt.Println("slice s1:", s1, "s2:", s2)

	fmt.Println("=== len and cap ===")
	arr := [8]int{0, 1, 2, 3, 4, 5, 6, 7}
	s := arr[2:5]
	fmt.Println("s =", s, "len =", len(s), "cap =", cap(s))
	dst := make([]int, 2)
	fmt.Println("copy copied", copy(dst, s), "elements:", dst) // limited by len(dst)

	fmt.Println("=== append growth (exact numbers are an implementation detail) ===")
	var g []int
	last := -1
	for i := 0; i < 2000; i++ {
		g = append(g, i)
		if cap(g) != last {
			last = cap(g)
			fmt.Print(last, " ")
		}
	}
	fmt.Println()

	fmt.Println("=== append aliasing bug ===")
	base := make([]int, 3, 10)
	x := append(base, 1)
	y := append(base, 2) // reuses the same spare slot
	fmt.Println("x[3] =", x[3], "y[3] =", y[3], "(x was overwritten)")
	x = append(base[:3:3], 1) // full slice expression forces a new array
	y = append(base[:3:3], 2)
	fmt.Println("fixed: x[3] =", x[3], "y[3] =", y[3])

	fmt.Println("=== slices as parameters ===")
	p := []int{1, 2, 3}
	appendInside(p)
	fmt.Println("after append in callee:", p)
	setFirst(p)
	fmt.Println("after element write in callee:", p)

	fmt.Println("=== nil vs empty ===")
	var n []int
	e := []int{}
	nj, _ := json.Marshal(n)
	ej, _ := json.Marshal(e)
	fmt.Println("nil slice == nil:", n == nil, "| empty slice == nil:", e == nil, "| json:", string(nj), string(ej))

	fmt.Println("=== memory retention ===")
	big := make([]byte, 1<<20)
	small := big[:4]              // keeps the whole 1 MiB array alive
	safe := slices.Clone(big[:4]) // copies 4 bytes, big can be freed
	fmt.Println("small cap:", cap(small), "| safe cap:", cap(safe))
}
