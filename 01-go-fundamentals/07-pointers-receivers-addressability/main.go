// Shows value vs pointer receivers, method sets, addressability and loop-variable copies.
// Run: go run ./01-go-fundamentals/07-pointers-receivers-addressability
package main

import "fmt"

type Counter struct{ N int }

func (c Counter) IncValue()    { c.N++ } // works on a copy
func (c *Counter) IncPointer() { c.N++ } // works on the original

type Incrementer interface{ IncPointer() }

func resetCopy(c Counter) { c.N = 0 }
func resetPtr(c *Counter) { c.N = 0 }

func main() {
	fmt.Println("=== value vs pointer receiver ===")
	c := Counter{}
	c.IncValue()
	fmt.Println("after IncValue:", c.N)
	c.IncPointer() // Go rewrites this to (&c).IncPointer() because c is addressable
	fmt.Println("after IncPointer:", c.N)

	fmt.Println("=== function arguments are always copies ===")
	resetCopy(c)
	fmt.Println("after resetCopy:", c.N)
	resetPtr(&c)
	fmt.Println("after resetPtr:", c.N)

	fmt.Println("=== method sets and interfaces ===")
	var inc Incrementer = &c // Counter (a value) would not compile: IncPointer is not in its method set
	inc.IncPointer()
	fmt.Println("through interface:", c.N)

	fmt.Println("=== not addressable: map elements ===")
	m := map[string]Counter{"a": {}}
	// m["a"].IncPointer() does not compile: cannot take the address of a map element
	tmp := m["a"]
	tmp.IncPointer()
	m["a"] = tmp
	fmt.Println("map value after copy-modify-store:", m["a"].N)
	pm := map[string]*Counter{"a": {}}
	pm["a"].IncPointer() // fine: the pointer is copied, the Counter is shared
	fmt.Println("map of pointers:", pm["a"].N)

	fmt.Println("=== range copies the element ===")
	list := []Counter{{1}, {2}}
	for _, item := range list {
		item.N *= 10 // changes the copy
	}
	fmt.Println("after range copy:", list)
	for i := range list {
		list[i].N *= 10
	}
	fmt.Println("after index write:", list)

	fmt.Println("=== new and returning a pointer to a local ===")
	p := new(int)
	*p = 7
	fmt.Println(*p, "| from newCounter:", newCounter().N)
}

func newCounter() *Counter {
	c := Counter{N: 5}
	return &c // safe: escape analysis moves c to the heap
}
