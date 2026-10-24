// Token bucket rate limiter (lazy refill, no goroutine) and a per-key limiter.
// golang.org/x/time/rate is the production choice; this shows how it works.
// Run: go run ./02-concurrency/13-rate-limiter-token-bucket
package main

import (
	"fmt"
	"sync"
	"time"
)

type Bucket struct {
	mu       sync.Mutex
	tokens   float64
	capacity float64
	rate     float64 // tokens per second
	last     time.Time
	now      func() time.Time // injectable so the demo is deterministic
}

func NewBucket(rate float64, capacity int, now func() time.Time) *Bucket {
	return &Bucket{tokens: float64(capacity), capacity: float64(capacity), rate: rate, last: now(), now: now}
}

func (b *Bucket) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	t := b.now()
	b.tokens = min(b.capacity, b.tokens+t.Sub(b.last).Seconds()*b.rate) // refill lazily
	b.last = t
	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

// PerKey keeps one bucket per client (IP, API key). Real code must evict idle keys.
type PerKey struct {
	mu      sync.Mutex
	buckets map[string]*Bucket
	rate    float64
	burst   int
	now     func() time.Time
}

func (p *PerKey) Allow(key string) bool {
	p.mu.Lock()
	b, ok := p.buckets[key]
	if !ok {
		b = NewBucket(p.rate, p.burst, p.now)
		p.buckets[key] = b
	}
	p.mu.Unlock()
	return b.Allow()
}

func main() {
	clock := time.Now()
	now := func() time.Time { return clock }

	fmt.Println("=== burst of 3, then 2 tokens/sec ===")
	b := NewBucket(2, 3, now)
	for i := 1; i <= 5; i++ {
		fmt.Printf("request %d allowed: %v\n", i, b.Allow())
	}
	clock = clock.Add(500 * time.Millisecond) // fake time: 1 token refilled
	fmt.Println("after 500ms allowed:", b.Allow())
	fmt.Println("immediately again  :", b.Allow())

	fmt.Println("=== per-key: one noisy client does not block another ===")
	p := &PerKey{buckets: map[string]*Bucket{}, rate: 1, burst: 1, now: now}
	fmt.Println("alice:", p.Allow("alice"), p.Allow("alice"))
	fmt.Println("bob  :", p.Allow("bob"))
}
