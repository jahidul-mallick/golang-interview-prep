// Shows zero values, := shadowing, untyped constants and iota.
// Run: go run ./01-go-fundamentals/01-variables-types-constants-iota
package main

import (
	"fmt"
	"strconv"
)

type Weekday int

const (
	Sunday Weekday = iota
	Monday
	Tuesday
)

func (d Weekday) String() string { return [...]string{"Sunday", "Monday", "Tuesday"}[d] }

const (
	_  = iota // skip 0
	KB = 1 << (10 * iota)
	MB
	GB
)

type Perm uint8

const (
	Read Perm = 1 << iota
	Write
	Exec
)

func main() {
	fmt.Println("=== zero values ===")
	var (
		n int
		s string
		p *int
		m map[string]int
		l []int
	)
	fmt.Printf("int=%d string=%q ptr nil=%v map nil=%v slice nil=%v\n", n, s, p == nil, m == nil, l == nil)
	fmt.Println("read from nil map:", m["x"], "| append to nil slice:", append(l, 1))

	fmt.Println("=== := shadowing ===")
	var err error
	if true {
		_, err := strconv.Atoi("abc") // declares a NEW err in the if block
		_ = err
	}
	fmt.Println("outer err is still:", err)

	fmt.Println("=== untyped constants ===")
	const big = 1 << 100 // fine: untyped constants have arbitrary precision
	fmt.Println("big >> 98 =", big>>98)
	fmt.Println("7/2 =", 7/2, "| 7.0/2 =", 7.0/2)
	var u uint8 = 255
	u++ // typed integers wrap around at runtime
	fmt.Println("uint8 255+1 =", u, "| int size in bits:", strconv.IntSize)

	fmt.Println("=== iota ===")
	fmt.Println(Sunday, Monday, Tuesday)
	fmt.Println("KB, MB, GB =", KB, MB, GB)
	perm := Read | Exec
	fmt.Printf("perm=%03b can write=%v can exec=%v\n", perm, perm&Write != 0, perm&Exec != 0)
}
