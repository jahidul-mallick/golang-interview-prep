// Memory layout of common Go types and why field order matters.
// Run: go run ./03-runtime-internals/07-data-structure-internals-memory-layout
package main

import (
	"fmt"
	"unsafe"
)

type Padded struct { // padding is inserted to align each field
	a bool
	b int64
	c bool
}

type Packed struct { // same fields, largest first
	b int64
	a bool
	c bool
}

func main() {
	fmt.Println("=== struct padding ===")
	fmt.Println("Padded:", unsafe.Sizeof(Padded{}), "bytes")
	fmt.Println("Packed:", unsafe.Sizeof(Packed{}), "bytes")

	fmt.Println("=== headers (64-bit) ===")
	fmt.Println("slice:", unsafe.Sizeof([]int{}), "bytes (pointer, len, cap)")
	fmt.Println("string:", unsafe.Sizeof(""), "bytes (pointer, len)")
	fmt.Println("interface:", unsafe.Sizeof(any(0)), "bytes (type, data)")
	fmt.Println("map:", unsafe.Sizeof(map[int]int{}), "bytes (one pointer)")
}
