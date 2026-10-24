// Shows recover in a deferred func, runtime panics turned into errors,
// per-goroutine recovery and re-panic.
// Run: go run ./01-go-fundamentals/10-panic-recover
package main

import (
	"fmt"
	"sync"
)

// safeCall converts a panic into an error using a named result.
func safeCall(f func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered: %v", r)
		}
	}()
	f()
	return nil
}

func safeGo(wg *sync.WaitGroup, id int, f func()) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("worker %d recovered: %v\n", id, r)
			}
		}()
		f()
	}()
}

func main() {
	fmt.Println("=== defers run while unwinding (LIFO) ===")
	_ = safeCall(func() {
		defer fmt.Println("deferred 1")
		defer fmt.Println("deferred 2")
		panic("boom")
	})

	fmt.Println("=== common runtime panics ===")
	var m map[string]int
	var p *struct{ X int }
	s := []int{1, 2, 3}
	i, zero := 10, 0
	cases := []struct {
		name string
		f    func()
	}{
		{"nil map write", func() { m["a"] = 1 }},
		{"nil pointer", func() { _ = p.X }},
		{"index range", func() { _ = s[i] }},
		{"divide by 0", func() { _ = 1 / zero }},
	}
	for _, c := range cases {
		fmt.Printf("%-14s -> %v\n", c.name, safeCall(c.f))
	}

	fmt.Println("=== each goroutine needs its own recover ===")
	var wg sync.WaitGroup
	for id := 1; id <= 3; id++ {
		safeGo(&wg, id, func() {
			if id == 2 {
				panic("bad job")
			}
		})
	}
	wg.Wait()

	fmt.Println("=== re-panic what you cannot handle ===")
	err := safeCall(func() {
		defer func() {
			r := recover()
			panic(fmt.Sprint("re-panic of ", r))
		}()
		panic("first")
	})
	fmt.Println(err)
}
