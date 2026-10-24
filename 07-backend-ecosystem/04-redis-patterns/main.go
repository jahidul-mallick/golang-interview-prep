// Shows cache-aside with a TTL cache and stampede protection (singleflight idea).
// Run: go run ./07-backend-ecosystem/04-redis-patterns
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type entry struct {
	val string
	exp time.Time
}

// Cache stands in for Redis GET / SET key val EX ttl.
type Cache struct {
	mu   sync.Mutex
	data map[string]entry
}

func (c *Cache) Get(k string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.data[k]
	if !ok || time.Now().After(e.exp) {
		delete(c.data, k)
		return "", false
	}
	return e.val, true
}

func (c *Cache) Set(k, v string, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[k] = entry{v, time.Now().Add(ttl)}
}

func (c *Cache) Del(k string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, k)
}

var dbCalls atomic.Int32

func queryDB(id string) string {
	dbCalls.Add(1)
	time.Sleep(50 * time.Millisecond) // slow query
	return "user-" + id
}

// call is one in-flight load that concurrent callers share.
type call struct {
	wg  sync.WaitGroup
	val string
}

var (
	flightMu sync.Mutex
	flight   = map[string]*call{}
)

// load runs queryDB once per key even if many goroutines miss at the same time.
func load(cache *Cache, key string) string {
	flightMu.Lock()
	if f, ok := flight[key]; ok {
		flightMu.Unlock()
		f.wg.Wait()
		return f.val
	}
	f := &call{}
	f.wg.Add(1)
	flight[key] = f
	flightMu.Unlock()

	f.val = queryDB(key)
	cache.Set(key, f.val, time.Minute) // fill the cache before leaving the flight map
	flightMu.Lock()
	delete(flight, key)
	flightMu.Unlock()
	f.wg.Done()
	return f.val
}

func getUser(c *Cache, id string) string {
	if v, ok := c.Get(id); ok {
		return v // hit
	}
	return load(c, id) // miss: read the source of truth and fill the cache
}

func main() {
	c := &Cache{data: map[string]entry{}}

	fmt.Println("=== cache-aside: miss then hit ===")
	getUser(c, "1")
	getUser(c, "1")
	fmt.Println("  db calls after 2 reads:", dbCalls.Load())

	fmt.Println("=== write: update DB, then delete the cache key ===")
	c.Del("1")
	getUser(c, "1")
	fmt.Println("  db calls:", dbCalls.Load())

	fmt.Println("=== stampede: 20 goroutines miss at once ===")
	dbCalls.Store(0)
	var wg sync.WaitGroup
	for range 20 {
		wg.Add(1)
		go func() { defer wg.Done(); getUser(c, "2") }()
	}
	wg.Wait()
	fmt.Println("  db calls:", dbCalls.Load(), "(not 20)")
}
