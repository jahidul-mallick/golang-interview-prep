# LLD Interview Approach & Checklist

[← Section index](../README.md) · Prev: [LLD: State Machines — Vending Machine](../09-lld-state-machines-vending-machine-elevator/README.md)

**Interview weight: 1.1%**  ·  Revise in: 8 min  ·  Run: `go run ./06-lld-design-patterns/10-lld-interview-approach-and-checklist`

## In one line
Clarify, list entities and operations, define interfaces, code the core path, then talk about concurrency and extension.

## Must know
- **0-5 min, clarify**: operations needed, scale hints, edge cases, what is out of scope.
- **5-10 min, model**: nouns become structs, verbs become methods.
- **10-15 min, interfaces**: small interfaces at the points that will vary (pricing, storage, notifier).
- **15-30 min, code**: the 2 or 3 core operations with error handling. Say what you skip.
- **Last minutes, discuss**: thread safety, complexity, how to extend, how to test.
- Checklist: SOLID, small interfaces, errors for not found / invalid / full, locks where state is shared.
- Think out loud and write compiling Go. Simple names, no clever tricks.
- Typical data structures: map plus slice or heap, `container/list` for LRU, a transition table for state machines.

## Code to remember
```go
// Splitwise: net balance per person, then match biggest debtor with biggest creditor.
amt := min(cred[0].amt, debt[0].amt)
out = append(out, Transfer{debt[0].name, cred[0].name, amt})
cred[0].amt -= amt
debt[0].amt -= amt
```

## Common traps
- Starting to code before clarifying requirements.
- Over-engineering with patterns nobody asked for. Add an interface only where behaviour varies.
- Ignoring concurrency when the problem has many users. At least mention a mutex.
- Map iteration order is random. Sort keys when you need deterministic output.

## Interview Q&A

**Q1. Splitwise: how do you reduce the number of transactions?**

**A:** Compute each person's net balance, then repeatedly settle the largest debtor against the largest creditor. This gives a small number of transfers; the true minimum is a harder problem.

**Q2. How do you make a struct thread-safe?**

**A:** Guard shared fields with a `sync.Mutex` or `RWMutex`, or give one goroutine ownership and talk to it over channels.

**Q3. Which pattern fits which requirement?**

**A:** Swappable rule (pricing, retry): strategy. Notify many listeners: observer. Lifecycle (order, vending): state. Create by kind: factory. Wrap behaviour: decorator or middleware. Many optional settings: functional options.

**Q4. Which LLD problems come up most?**

**A:** Parking lot, LRU cache, KV store with TTL, rate limiter, vending machine, rule engine, Splitwise. Practise the first four.
