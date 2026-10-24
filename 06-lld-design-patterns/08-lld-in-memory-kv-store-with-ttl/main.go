// In-memory KV store with TTL: lazy expiry on Get plus a background sweeper.
// Run: go run ./06-lld-design-patterns/08-lld-in-memory-kv-store-with-ttl
package main

import (
	"fmt"
	"sync"
	"time"
)

type entry struct {
	val     any
	expires time.Time // zero value means no expiry
}

type Store struct {
	mu   sync.RWMutex
	data map[string]entry
	stop chan struct{}
}

func NewStore(sweepEvery time.Duration) *Store {
	s := &Store{data: map[string]entry{}, stop: make(chan struct{})}
	go s.sweep(sweepEvery)
	return s
}

func (s *Store) Set(key string, val any, ttl time.Duration) {
	e := entry{val: val}
	if ttl > 0 {
		e.expires = time.Now().Add(ttl)
	}
	s.mu.Lock()
	s.data[key] = e
	s.mu.Unlock()
}

// Get treats an expired entry as missing (lazy expiry); the sweeper frees it later.
func (s *Store) Get(key string) (any, bool) {
	s.mu.RLock()
	e, ok := s.data[key]
	s.mu.RUnlock()
	if !ok || (!e.expires.IsZero() && time.Now().After(e.expires)) {
		return nil, false
	}
	return e.val, true
}

func (s *Store) Delete(key string) {
	s.mu.Lock()
	delete(s.data, key)
	s.mu.Unlock()
}

func (s *Store) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data)
}

func (s *Store) sweep(every time.Duration) {
	t := time.NewTicker(every)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			now := time.Now()
			s.mu.Lock()
			for k, e := range s.data {
				if !e.expires.IsZero() && now.After(e.expires) {
					delete(s.data, k)
				}
			}
			s.mu.Unlock()
		case <-s.stop:
			return
		}
	}
}

func (s *Store) Close() { close(s.stop) }

func main() {
	fmt.Println("=== KV store with TTL ===")
	s := NewStore(200 * time.Millisecond)
	defer s.Close()

	s.Set("name", "asha", 0)
	s.Set("session", "tok123", 50*time.Millisecond)
	v, ok := s.Get("session")
	fmt.Println("  session:", v, ok, "| entries:", s.Len())

	time.Sleep(80 * time.Millisecond)
	v, ok = s.Get("session")
	fmt.Println("  after TTL (lazy):", v, ok, "| entries still stored:", s.Len())

	time.Sleep(250 * time.Millisecond)
	fmt.Println("  after sweep, entries:", s.Len())

	s.Delete("name")
	_, ok = s.Get("name")
	fmt.Println("  name after delete:", ok)
}
