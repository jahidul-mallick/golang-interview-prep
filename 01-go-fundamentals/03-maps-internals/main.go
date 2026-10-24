// Shows nil maps, comma-ok, random iteration order, sets, and safe concurrent use.
// Run: go run ./01-go-fundamentals/03-maps-internals
package main

import (
	"fmt"
	"maps"
	"slices"
	"sync"
)

type SafeCounter struct {
	mu sync.Mutex
	m  map[string]int
}

func (c *SafeCounter) Inc(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[key]++
}

func main() {
	fmt.Println("=== nil map: read ok, write panics ===")
	var nilMap map[string]int
	fmt.Println("read:", nilMap["a"], "len:", len(nilMap))
	func() {
		defer func() { fmt.Println("recovered:", recover()) }()
		nilMap["a"] = 1
	}()

	fmt.Println("=== comma-ok: missing key vs zero value ===")
	ages := map[string]int{"amit": 0, "riya": 30}
	v, ok := ages["amit"]
	fmt.Println("amit:", v, ok)
	v, ok = ages["zed"]
	fmt.Println("zed:", v, ok)

	fmt.Println("=== iteration order is random: sort the keys ===")
	for _, k := range slices.Sorted(maps.Keys(ages)) {
		fmt.Println(k, ages[k])
	}

	fmt.Println("=== map values are not addressable ===")
	type user struct{ Score int }
	users := map[string]user{"a": {1}}
	u := users["a"] // users["a"].Score++ does not compile
	u.Score++
	users["a"] = u
	fmt.Println("users:", users)

	fmt.Println("=== set with map[T]struct{} ===")
	set := map[string]struct{}{}
	for _, w := range []string{"go", "rust", "go"} {
		set[w] = struct{}{}
	}
	fmt.Println("unique words:", len(set))

	fmt.Println("=== delete during range is allowed ===")
	nums := map[int]bool{1: true, 2: true, 3: true, 4: true}
	for k := range nums {
		if k%2 == 0 {
			delete(nums, k)
		}
	}
	fmt.Println("left:", slices.Sorted(maps.Keys(nums)))

	fmt.Println("=== concurrent writes need a mutex ===")
	c := SafeCounter{m: map[string]int{}}
	var wg sync.WaitGroup
	for range 50 {
		wg.Go(func() { c.Inc("hits") })
	}
	wg.Wait()
	fmt.Println("hits:", c.m["hits"])
}
