// LRU cache: hash map + doubly linked list, O(1) Get and Put.
// Run: go run ./04-dsa/16-lru-lfu-cache-design
package main

import (
	"container/list"
	"fmt"
)

type entry struct{ key, val int }

// LRU is not safe for concurrent use; wrap calls in a sync.Mutex if needed.
type LRU struct {
	cap   int
	items map[int]*list.Element // key -> node in order
	order *list.List            // front = most recent, back = least recent
}

func NewLRU(cap int) *LRU {
	return &LRU{cap: cap, items: map[int]*list.Element{}, order: list.New()}
}

func (c *LRU) Get(key int) int {
	el, ok := c.items[key]
	if !ok {
		return -1
	}
	c.order.MoveToFront(el)
	return el.Value.(*entry).val
}

func (c *LRU) Put(key, val int) {
	if el, ok := c.items[key]; ok {
		el.Value.(*entry).val = val
		c.order.MoveToFront(el)
		return
	}
	if c.order.Len() == c.cap {
		oldest := c.order.Back()
		c.order.Remove(oldest)
		delete(c.items, oldest.Value.(*entry).key) // node stores the key for this
	}
	c.items[key] = c.order.PushFront(&entry{key, val})
}

func main() {
	fmt.Println("=== LRU cache, capacity 2 ===")
	c := NewLRU(2)
	c.Put(1, 10)
	c.Put(2, 20)
	fmt.Println("Get(1) =", c.Get(1)) // 10, key 1 is now most recent
	c.Put(3, 30)                      // evicts key 2
	fmt.Println("Get(2) =", c.Get(2)) // -1
	c.Put(4, 40)                      // evicts key 1
	fmt.Println("Get(1) =", c.Get(1), "Get(3) =", c.Get(3), "Get(4) =", c.Get(4))
}
