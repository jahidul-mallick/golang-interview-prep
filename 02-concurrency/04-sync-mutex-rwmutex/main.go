// sync.Mutex, sync.RWMutex, and lock ordering to avoid deadlock.
// Run: go run ./02-concurrency/04-sync-mutex-rwmutex
package main

import (
	"fmt"
	"sync"
)

type Counter struct {
	mu sync.Mutex // guards n
	n  int
}

func (c *Counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.n++
}

func (c *Counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.n
}

type Cache struct {
	mu sync.RWMutex // guards m
	m  map[string]string
}

func (c *Cache) Get(k string) (string, bool) {
	c.mu.RLock() // many readers at once
	defer c.mu.RUnlock()
	v, ok := c.m[k]
	return v, ok
}

func (c *Cache) Set(k, v string) {
	c.mu.Lock() // one writer, no readers
	defer c.mu.Unlock()
	c.m[k] = v
}

type Account struct {
	mu      sync.Mutex
	id      int
	balance int
}

// Always lock the lower id first so two opposite transfers cannot deadlock.
func transfer(from, to *Account, amount int) {
	first, second := from, to
	if second.id < first.id {
		first, second = second, first
	}
	first.mu.Lock()
	defer first.mu.Unlock()
	second.mu.Lock()
	defer second.mu.Unlock()
	from.balance -= amount
	to.balance += amount
}

func main() {
	fmt.Println("=== Mutex: 100 goroutines x 1000 increments ===")
	var c Counter
	var wg sync.WaitGroup
	for range 100 {
		wg.Go(func() {
			for range 1000 {
				c.Inc()
			}
		})
	}
	wg.Wait()
	fmt.Println("counter =", c.Value())

	fmt.Println("\n=== RWMutex: parallel readers, exclusive writer ===")
	cache := &Cache{m: map[string]string{"region": "us-east-1"}}
	for i := range 4 {
		wg.Go(func() { cache.Get("region") })
		wg.Go(func() { cache.Set("region", fmt.Sprint("eu-west-", i)) })
	}
	wg.Wait()
	_, ok := cache.Get("region")
	fmt.Println("no race, key present:", ok)

	fmt.Println("\n=== lock ordering: opposite transfers do not deadlock ===")
	a, b := &Account{id: 1, balance: 100}, &Account{id: 2, balance: 100}
	for range 1000 {
		wg.Go(func() { transfer(a, b, 1) })
		wg.Go(func() { transfer(b, a, 1) })
	}
	wg.Wait()
	fmt.Println("total balance:", a.balance+b.balance, "(still 200)")
}
