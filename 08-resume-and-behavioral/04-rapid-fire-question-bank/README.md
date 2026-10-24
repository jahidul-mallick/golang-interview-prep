# Rapid-Fire Question Bank

[← Section index](../README.md) · Prev: [Interview Rounds and Questions to Ask](../03-interview-rounds-and-questions-to-ask/README.md) · Next: [30-Day Study Plan](../05-30-day-study-plan/README.md)

**Interview weight: 1.2%**  ·  Revise in: 20 min

## In one line
Answer each question out loud in 30 seconds or less; if you stall, open the topic folder in brackets.

## Go fundamentals (01-go-fundamentals)
- What are the zero values of map, slice, channel, pointer, interface? [01/01]
- What is in a slice header, and what happens on append past capacity? [01/02]
- Can you read from and write to a nil map? [01/03]
- Why can iterating a string give surprising results? (bytes vs runes) [01/04]
- Value receiver vs pointer receiver: when to use which? [01/07]
- When is an interface value nil but not equal to nil? [01/06]
- How do defer, panic and recover work together? [01/08, 01/10]
- How do you wrap and check errors? (`%w`, `errors.Is`, `errors.As`) [01/09]
- What problem do generics solve, and what is a type constraint? [01/11]

## Concurrency (02-concurrency)
- When does a goroutine leak? Name three causes. [02/01]
- What happens when you close a closed channel, or send on a closed one? [02/02]
- Buffered vs unbuffered channel: when to use each? [02/02]
- What is `select` and what is `default` for? [02/03]
- Mutex vs RWMutex vs atomic: how to choose? [02/04, 02/07]
- How do you cancel a tree of goroutines? [02/08]
- Describe a worker pool. [02/09]
- What do `sync.WaitGroup` and `sync.Once` guarantee? [02/05]
- How do you do graceful HTTP shutdown? [02/12]
- What is a data race and how do you detect it? [02/14]

## Runtime (03-runtime-internals)
- What are G, M and P? [03/01]
- When does a value escape to the heap? [03/03]
- How does Go's garbage collector work at a high level? [03/05]
- What do GOGC and GOMEMLIMIT do? [03/05]
- How do you find a CPU or memory problem with pprof? [03/06]

## DSA (04-dsa)
- When do you use two pointers or a sliding window? [04/02, 04/03]
- Which data structures make an LRU cache O(1)? [04/16]
- BFS vs DFS: when to use each? [04/13]

## System design and LLD (05, 06)
- Cache-aside vs write-through? [05/02]
- Replication vs sharding: what problem does each solve? [05/03]
- What does the CAP theorem say? [05/04]
- How do you make a consumer idempotent? [05/05]
- Circuit breaker: what are its states? [05/06]
- Token bucket: how does it work? [05/09]
- Saga vs two-phase commit? Orchestration vs choreography? [05/17]
- What is the outbox pattern and what problem does it solve? [05/17]
- Is a Redis lock safe? What is a fencing token? [05/17]
- When would you use CQRS? [05/17]
- How would you design a multi-cloud resource discovery platform? [05/13]
- Strategy, factory, functional options: where do you use them in Go? [06/02, 06/04]

## Backend ecosystem (07-backend-ecosystem)
- How do you write HTTP middleware in Go? [07/01]
- What is an index in PostgreSQL and when does it not help? [07/03]
- What is a Kubernetes liveness vs readiness probe? [07/08]
- How do mTLS and JWT differ? [07/10]

## Common traps
- Memorising answers word for word. Understand the idea so you can handle follow-ups.
- Rambling. Give a one-line answer first, then add detail if asked.

## Interview Q&A

**Q1. How should I use this list?**

**A:** Cover the answers, say each one out loud, and mark the ones you missed. Revise those folders first.
