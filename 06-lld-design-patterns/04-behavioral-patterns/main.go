// Behavioral patterns: strategy, observer, state.
// Run: go run ./06-lld-design-patterns/04-behavioral-patterns
package main

import (
	"fmt"
	"sync"
	"time"
)

// Strategy: the algorithm is a func value chosen by the caller.
type Backoff func(attempt int) time.Duration

func Constant(d time.Duration) Backoff { return func(int) time.Duration { return d } }
func Exponential(base time.Duration) Backoff {
	return func(n int) time.Duration { return base << n }
}

// Retry prints the wait instead of sleeping so the demo is instant.
func Retry(op func() error, tries int, b Backoff) error {
	var err error
	for i := 0; i < tries; i++ {
		if err = op(); err == nil {
			return nil
		}
		fmt.Printf("  attempt %d failed, would wait %v\n", i+1, b(i))
	}
	return fmt.Errorf("gave up after %d tries: %w", tries, err)
}

// Observer: publishers know only a topic, not who listens.
type Bus struct {
	mu   sync.RWMutex
	subs map[string][]func(string)
}

func NewBus() *Bus { return &Bus{subs: map[string][]func(string){}} }

func (b *Bus) Subscribe(topic string, h func(string)) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subs[topic] = append(b.subs[topic], h)
}

// Publish copies the handler list under the lock and calls handlers after
// unlocking, so a slow or re-entrant handler cannot deadlock the bus.
func (b *Bus) Publish(topic, msg string) {
	b.mu.RLock()
	hs := b.subs[topic]
	b.mu.RUnlock()
	for _, h := range hs {
		h(msg)
	}
}

// State: a transition table; a missing key is an illegal transition.
type State string
type Event string

var next = map[State]map[Event]State{
	"created": {"pay": "paid"},
	"paid":    {"ship": "shipped", "cancel": "cancelled"},
	"shipped": {"deliver": "delivered"},
}

func Apply(s State, e Event) (State, error) {
	if to, ok := next[s][e]; ok {
		return to, nil
	}
	return s, fmt.Errorf("cannot %s from %s", e, s)
}

func main() {
	fmt.Println("=== strategy: retry with backoff ===")
	calls := 0
	flaky := func() error {
		if calls++; calls < 3 {
			return fmt.Errorf("boom")
		}
		return nil
	}
	fmt.Println("  result:", Retry(flaky, 5, Exponential(100*time.Millisecond)))
	fmt.Println("  result:", Retry(func() error { return fmt.Errorf("down") }, 2, Constant(time.Second)))

	fmt.Println("=== observer ===")
	bus := NewBus()
	bus.Subscribe("order.paid", func(m string) { fmt.Println("  email:", m) })
	bus.Subscribe("order.paid", func(m string) { fmt.Println("  ledger:", m) })
	bus.Publish("order.paid", "order 42")
	bus.Publish("order.unknown", "nobody listens")

	fmt.Println("=== state ===")
	s := State("created")
	for _, e := range []Event{"pay", "cancel", "ship"} {
		var err error
		if s, err = Apply(s, e); err != nil {
			fmt.Println(" ", err)
			continue
		}
		fmt.Printf("  %s -> %s\n", e, s)
	}
}
