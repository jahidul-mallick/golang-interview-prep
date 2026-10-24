// Shows Kafka ideas in memory: key -> partition, per-partition order,
// consumer offsets, and at-least-once delivery with an idempotent handler.
// Run: go run ./07-backend-ecosystem/05-kafka-nats-pubsub-deep-dive
package main

import (
	"fmt"
	"hash/fnv"
)

type Msg struct{ Key, Val string }

// Topic is a set of append-only partitions.
type Topic struct{ parts [][]Msg }

func (t *Topic) Produce(m Msg) int {
	h := fnv.New32a()
	h.Write([]byte(m.Key))
	p := int(h.Sum32() % uint32(len(t.parts))) // same key -> same partition -> ordered
	t.parts[p] = append(t.parts[p], m)
	return p
}

// Consumer owns one partition and tracks the committed offset (next to read).
type Consumer struct {
	topic     *Topic
	part      int
	committed int
	seen      map[string]bool // idempotency: processed message ids
}

// Poll processes up to n messages and commits only if commit is true.
// Not committing simulates a crash before the commit.
func (c *Consumer) Poll(n int, commit bool) {
	msgs := c.topic.parts[c.part]
	for off := c.committed; off < len(msgs) && off < c.committed+n; off++ {
		id := fmt.Sprintf("p%d-o%d", c.part, off)
		if c.seen[id] {
			fmt.Printf("  %s duplicate, skipped\n", id)
			continue
		}
		c.seen[id] = true
		fmt.Printf("  %s key=%s val=%s processed\n", id, msgs[off].Key, msgs[off].Val)
	}
	if commit {
		c.committed = min(c.committed+n, len(msgs))
	}
}

func main() {
	t := &Topic{parts: make([][]Msg, 3)}

	fmt.Println("=== same key always lands in the same partition ===")
	for _, m := range []Msg{{"u1", "created"}, {"u2", "created"}, {"u1", "paid"}, {"u3", "created"}, {"u1", "shipped"}} {
		fmt.Printf("  key=%-3s val=%-8s -> partition %d\n", m.Key, m.Val, t.Produce(m))
	}
	p := t.Produce(Msg{"u1", "closed"})

	fmt.Println("=== at-least-once: crash before commit, then redelivery ===")
	c := &Consumer{topic: t, part: p, seen: map[string]bool{}}
	c.Poll(2, false) // processed, but the commit never happened
	fmt.Println("  consumer restarts, committed offset =", c.committed)
	c2 := &Consumer{topic: t, part: p, seen: c.seen} // seen stands in for a DB unique key
	c2.Poll(2, true)
	fmt.Println("  committed offset after success =", c2.committed)
}
