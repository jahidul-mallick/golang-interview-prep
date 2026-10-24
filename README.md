# Go Backend Interview Prep (Mid-Level, 4-5 years)

**Author:** Jahidul Mallick

The minimum you need to crack mid-level Go backend interviews. Every topic is one short README plus one small runnable `main.go`. Nothing here is advanced on purpose.

**Weight** = the share of your total preparation effort a topic deserves. All 100 weighted topics add up to **100%**, and so do the 8 sections below. A bigger weight means it is asked more often in mid-level Go backend interviews, so study from the top weight down. The numbers are my judgement, not measured data.

## How to use this repo

1. Start with the highest weights first (the "Must study" list below).
2. For each topic: read the README, run `go run ./<section>/<topic>`, then answer the Q&A out loud.
3. Last two days: re-read only the **Must know** blocks and the rapid-fire bank in section 08.
4. Requirements: Go 1.25+, standard library only, no setup.

## Sections at a glance

| # | Section | Topics | Share of prep | Priority |
|---|---|---|---|---|
| 02 | [Concurrency](02-concurrency/README.md) | 16 | 18.1% | Highest |
| 01 | [Go fundamentals](01-go-fundamentals/README.md) | 16 | 16.0% | Highest |
| 08 | [Resume and behavioral](08-resume-and-behavioral/README.md) | 5 | 4.2% | Highest |
| 05 | [System design (HLD)](05-system-design-hld/README.md) | 17 | 17.2% | High |
| 07 | [Backend ecosystem](07-backend-ecosystem/README.md) | 11 | 11.5% | High |
| 03 | [Runtime internals](03-runtime-internals/README.md) | 8 | 7.2% | High |
| 04 | [DSA in Go](04-dsa/README.md) | 18 | 16.3% | Medium (easy and medium problems only) |
| 06 | [LLD and design patterns](06-lld-design-patterns/README.md) | 10 | 9.5% | Medium |

## Must study first (the highest-weight topics)

Concurrency: goroutines, channels, select, mutex, WaitGroup/Once, context, worker pool, races and deadlocks, classic problems.
Fundamentals: slices, interfaces, closures and defer, error handling, gotchas.
Runtime: GMP scheduler. Ecosystem: REST with net/http, PostgreSQL, Redis, Kafka/NATS. System design: databases, caching, queues, resilience, URL shortener, interview framework, your multi-cloud project. DSA: hashing, two pointers, sliding window. Behavioral: resume walkthrough, rapid-fire bank.

## Weight table

### 02 Concurrency
| Topic | Weight |
|---|---|
| [Goroutines basics](02-concurrency/01-goroutines-basics/README.md) | 1.4% |
| [Channels](02-concurrency/02-channels-unbuffered-buffered-closing/README.md) | 1.4% |
| [Select, timeouts, default](02-concurrency/03-select-timeouts-default/README.md) | 1.2% |
| [Mutex and RWMutex](02-concurrency/04-sync-mutex-rwmutex/README.md) | 1.3% |
| [WaitGroup and Once](02-concurrency/05-sync-waitgroup-once/README.md) | 1.2% |
| [Cond, Pool, Map](02-concurrency/06-sync-cond-pool-map/README.md) | 0.8% |
| [Atomic operations](02-concurrency/07-atomic-operations/README.md) | 1.0% |
| [Context](02-concurrency/08-context-cancellation-timeouts-values/README.md) | 1.3% |
| [Worker pool](02-concurrency/09-worker-pool/README.md) | 1.3% |
| [Pipeline, fan-out, fan-in](02-concurrency/10-pipeline-fan-out-fan-in/README.md) | 1.1% |
| [errgroup and semaphore](02-concurrency/11-errgroup-semaphore-bounded-concurrency/README.md) | 1.0% |
| [Graceful shutdown](02-concurrency/12-graceful-shutdown-http-server/README.md) | 1.0% |
| [Rate limiter (token bucket)](02-concurrency/13-rate-limiter-token-bucket/README.md) | 0.9% |
| [Races and deadlocks](02-concurrency/14-race-conditions-deadlocks-detection/README.md) | 1.3% |
| [Timers and tickers](02-concurrency/15-timers-tickers-time/README.md) | 0.7% |
| [Classic concurrency problems](02-concurrency/16-classic-concurrency-interview-problems/README.md) | 1.2% |

### 01 Go fundamentals
| Topic | Weight |
|---|---|
| [Variables, types, constants, iota](01-go-fundamentals/01-variables-types-constants-iota/README.md) | 0.8% |
| [Arrays and slices](01-go-fundamentals/02-arrays-slices-internals/README.md) | 1.3% |
| [Maps](01-go-fundamentals/03-maps-internals/README.md) | 1.2% |
| [Strings, runes, bytes](01-go-fundamentals/04-strings-runes-bytes/README.md) | 0.9% |
| [Structs, methods, embedding](01-go-fundamentals/05-structs-methods-embedding/README.md) | 1.1% |
| [Interfaces and type switch](01-go-fundamentals/06-interfaces-type-assertion-type-switch/README.md) | 1.3% |
| [Pointers and receivers](01-go-fundamentals/07-pointers-receivers-addressability/README.md) | 1.2% |
| [Functions, closures, defer](01-go-fundamentals/08-functions-closures-defer/README.md) | 1.3% |
| [Error handling](01-go-fundamentals/09-error-handling-wrapping/README.md) | 1.3% |
| [Panic and recover](01-go-fundamentals/10-panic-recover/README.md) | 1.0% |
| [Generics](01-go-fundamentals/11-generics/README.md) | 0.8% |
| [Packages, modules, init order](01-go-fundamentals/12-packages-modules-init-order/README.md) | 0.7% |
| [Testing and benchmarks](01-go-fundamentals/13-testing-benchmarks-fuzzing/README.md) | 1.0% |
| [Iterators (range over func)](01-go-fundamentals/14-iterators-range-over-func/README.md) | 0.3% |
| [JSON and reflection](01-go-fundamentals/15-json-reflection-struct-tags/README.md) | 0.6% |
| [Gotchas and tricky questions](01-go-fundamentals/16-gotchas-and-tricky-questions/README.md) | 1.2% |

### 03 Runtime internals
| Topic | Weight |
|---|---|
| [GMP scheduler](03-runtime-internals/01-gmp-scheduler/README.md) | 1.2% |
| [Memory model](03-runtime-internals/02-memory-model-happens-before/README.md) | 0.8% |
| [Stack, heap, escape analysis](03-runtime-internals/03-stack-heap-escape-analysis/README.md) | 1.0% |
| [Memory allocator](03-runtime-internals/04-memory-allocator-mcache-mcentral-mheap/README.md) | 0.6% |
| [Garbage collector](03-runtime-internals/05-garbage-collector/README.md) | 1.0% |
| [pprof and tracing](03-runtime-internals/06-pprof-profiling-tracing/README.md) | 1.0% |
| [Data structure memory layout](03-runtime-internals/07-data-structure-internals-memory-layout/README.md) | 0.5% |
| [Runtime Q&A](03-runtime-internals/08-runtime-interview-questions/README.md) | 1.1% |

### 04 DSA
| Topic | Weight |
|---|---|
| [Big-O and Go toolkit](04-dsa/01-big-o-and-go-toolkit-for-dsa/README.md) | 0.8% |
| [Two pointers, prefix sum](04-dsa/02-arrays-two-pointers-prefix-sum/README.md) | 1.2% |
| [Sliding window](04-dsa/03-sliding-window/README.md) | 1.2% |
| [Hashing](04-dsa/04-hashing-hashmap-problems/README.md) | 1.2% |
| [Strings](04-dsa/05-strings-problems/README.md) | 1.1% |
| [Linked list](04-dsa/06-linked-list/README.md) | 1.0% |
| [Stack and queue](04-dsa/07-stack-queue-monotonic/README.md) | 1.0% |
| [Binary search](04-dsa/08-binary-search/README.md) | 1.0% |
| [Sorting](04-dsa/09-sorting-algorithms/README.md) | 0.8% |
| [Recursion and backtracking](04-dsa/10-recursion-backtracking/README.md) | 0.8% |
| [Trees and BST](04-dsa/11-binary-tree-bst/README.md) | 1.1% |
| [Heap and top-K](04-dsa/12-heap-priority-queue-topk/README.md) | 0.9% |
| [Graphs](04-dsa/13-graphs-bfs-dfs-topological-sort-union-find/README.md) | 1.0% |
| [Dynamic programming](04-dsa/14-dynamic-programming/README.md) | 0.8% |
| [Trie](04-dsa/15-trie/README.md) | 0.3% |
| [LRU cache](04-dsa/16-lru-lfu-cache-design/README.md) | 1.1% |
| [Intervals and greedy](04-dsa/17-intervals-greedy/README.md) | 0.7% |
| [Bit manipulation](04-dsa/18-bit-manipulation-math/README.md) | 0.3% |

### 05 System design (HLD)
| Topic | Weight |
|---|---|
| [Fundamentals](05-system-design-hld/01-fundamentals-scalability-availability-latency/README.md) | 1.1% |
| [Load balancing and caching](05-system-design-hld/02-load-balancing-and-caching-strategies/README.md) | 1.2% |
| [Databases, replication, sharding](05-system-design-hld/03-databases-sql-nosql-indexing-replication-sharding/README.md) | 1.2% |
| [CAP and consistency](05-system-design-hld/04-cap-consistency-models-consensus/README.md) | 1.0% |
| [Message queues](05-system-design-hld/05-message-queues-kafka-nats-pubsub/README.md) | 1.2% |
| [Microservices and resilience](05-system-design-hld/06-microservices-communication-and-resilience/README.md) | 1.2% |
| [API design](05-system-design-hld/07-api-design-rest-grpc-graphql-versioning/README.md) | 1.1% |
| [Consistent hashing](05-system-design-hld/08-consistent-hashing-and-distributed-cache/README.md) | 0.7% |
| [Design a rate limiter](05-system-design-hld/09-design-rate-limiter/README.md) | 1.0% |
| [Design a URL shortener](05-system-design-hld/10-design-url-shortener/README.md) | 1.2% |
| [Design a notification system](05-system-design-hld/11-design-notification-system/README.md) | 0.8% |
| [Design a job scheduler](05-system-design-hld/12-design-distributed-job-scheduler/README.md) | 0.7% |
| [Design: multi-cloud discovery (your project)](05-system-design-hld/13-design-multi-cloud-resource-discovery-platform/README.md) | 1.3% |
| [Design a chat system](05-system-design-hld/14-design-chat-system-websockets/README.md) | 0.6% |
| [Design metrics and alerting](05-system-design-hld/15-design-metrics-monitoring-alerting-system/README.md) | 0.5% |
| [Interview framework](05-system-design-hld/16-system-design-interview-framework-and-checklist/README.md) | 1.3% |
| [Distributed design patterns](05-system-design-hld/17-distributed-design-patterns/README.md) | 1.1% |

### 06 LLD and design patterns
| Topic | Weight |
|---|---|
| [SOLID in Go](06-lld-design-patterns/01-solid-principles-in-go/README.md) | 1.2% |
| [Creational patterns](06-lld-design-patterns/02-creational-patterns/README.md) | 1.0% |
| [Structural patterns](06-lld-design-patterns/03-structural-patterns/README.md) | 0.8% |
| [Behavioral patterns](06-lld-design-patterns/04-behavioral-patterns/README.md) | 0.9% |
| [Idiomatic Go architecture](06-lld-design-patterns/05-idiomatic-go-architecture-patterns/README.md) | 1.0% |
| [LLD: parking lot](06-lld-design-patterns/06-lld-parking-lot/README.md) | 1.1% |
| [LLD: rule engine](06-lld-design-patterns/07-lld-rule-engine/README.md) | 0.8% |
| [LLD: KV store with TTL](06-lld-design-patterns/08-lld-in-memory-kv-store-with-ttl/README.md) | 0.9% |
| [LLD: state machines](06-lld-design-patterns/09-lld-state-machines-vending-machine-elevator/README.md) | 0.7% |
| [LLD approach and checklist](06-lld-design-patterns/10-lld-interview-approach-and-checklist/README.md) | 1.1% |

### 07 Backend ecosystem
| Topic | Weight |
|---|---|
| [REST with net/http](07-backend-ecosystem/01-rest-api-in-go-net-http-middleware/README.md) | 1.2% |
| [gRPC and protobuf](07-backend-ecosystem/02-grpc-protobuf-streaming-interceptors/README.md) | 1.1% |
| [PostgreSQL](07-backend-ecosystem/03-postgresql-for-go-developers/README.md) | 1.2% |
| [Redis](07-backend-ecosystem/04-redis-patterns/README.md) | 1.2% |
| [Kafka, NATS, Pub/Sub](07-backend-ecosystem/05-kafka-nats-pubsub-deep-dive/README.md) | 1.2% |
| [Neo4j and Cypher](07-backend-ecosystem/06-neo4j-graph-modeling-cypher/README.md) | 1.0% |
| [MongoDB](07-backend-ecosystem/07-mongodb-basics-and-schema-design/README.md) | 0.8% |
| [Kubernetes basics](07-backend-ecosystem/08-kubernetes-for-go-developers/README.md) | 1.0% |
| [Observability](07-backend-ecosystem/09-observability-prometheus-grafana-logging-tracing/README.md) | 1.0% |
| [Security: mTLS, JWT, Istio](07-backend-ecosystem/10-security-mtls-jwt-oauth-istio/README.md) | 1.0% |
| [Docker and CI/CD](07-backend-ecosystem/11-docker-ci-cd-for-go/README.md) | 0.8% |

### 08 Resume and behavioral
| Topic | Weight |
|---|---|
| [Resume walkthrough and project deep dive](08-resume-and-behavioral/01-resume-walkthrough-and-project-deep-dive/README.md) | 1.4% |
| [STAR stories](08-resume-and-behavioral/02-star-stories/README.md) | 1.0% |
| [Interview rounds and questions to ask](08-resume-and-behavioral/03-interview-rounds-and-questions-to-ask/README.md) | 0.6% |
| [Rapid-fire question bank](08-resume-and-behavioral/04-rapid-fire-question-bank/README.md) | 1.2% |
| [Study plan](08-resume-and-behavioral/05-30-day-study-plan/README.md) | plan |

## Commands

```powershell
go run ./02-concurrency/09-worker-pool   # run one topic
go vet ./...                              # sanity check everything
```

Diagrams are in [assets/](assets/): original SVG files, plus two PNG overviews for the Go runtime. Anything that needs a real client library (pgx, go-redis, grpc, neo4j, mongo) is shown only as a README snippet marked with the module it needs.
