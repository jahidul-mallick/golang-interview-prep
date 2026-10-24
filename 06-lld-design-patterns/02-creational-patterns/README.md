# Creational Patterns in Go

[← Section index](../README.md) · Prev: [SOLID Principles in Go](../01-solid-principles-in-go/README.md) · Next: [03 · Structural Patterns in Go](../03-structural-patterns/README.md)

**Interview weight: 1.0%**  ·  Revise in: 8 min  ·  Run: `go run ./06-lld-design-patterns/02-creational-patterns`

## In one line
Go creates objects with plain constructors; the patterns you must write are singleton with `sync.Once`, a factory function, functional options, and a builder.

## Must know
- **Singleton**: package-level value plus `sync.Once`. Thread-safe lazy init, runs exactly once.
- **Factory**: `NewX(kind) (Interface, error)`. Callers depend on the interface, not the concrete type.
- **Functional options**: `type Option func(*T)` and `New(required, opts ...Option)`. The idiomatic Go answer for many optional settings.
- **Builder**: chained methods and a final `Build()` that validates. Use it when construction has steps or must be validated as a whole.
- Constructors are just functions named `NewX`. A usable zero value is even better.
- Prefer options over a builder in Go unless you need step-by-step assembly (for example SQL queries).
- Global singletons hurt testing. Prefer passing the dependency in.

## Code to remember
```go
type Option func(*Server)

func WithPort(p int) Option { return func(s *Server) { s.port = p } }

func NewServer(host string, opts ...Option) *Server {
	s := &Server{host: host, port: 80} // defaults first
	for _, o := range opts {
		o(s)
	}
	return s
}
```

## Common traps
- Double-checked locking by hand. Use `sync.Once`, it already handles memory ordering.
- Forgetting that `once.Do` runs the function once even if it panics or fails. Store the error if init can fail.
- Using `init()` or globals for everything. It hides dependencies and makes tests order-dependent.
- Options that mutate shared state. Each call to `New` must build its own value.

## Interview Q&A

**Q1. How do you write a thread-safe singleton in Go?**

**A:** A package variable plus `sync.Once`: `once.Do(func(){ inst = &T{} })`, and return `inst`.

**Q2. Why functional options instead of a config struct?**

**A:** Defaults stay in one place, callers pass only what they change, and new options do not break existing callers.

**Q3. Factory vs builder?**

**A:** A factory picks which type to create. A builder assembles one complex object step by step.

**Q4. What does the factory return in Go?**

**A:** An interface (plus an error for unknown kinds). That is the case where returning an interface is justified.

**Q5. When is a singleton a bad idea?**

**A:** When it carries mutable state or is hard to fake in tests. Inject it through a constructor instead.

**Q6. What if the singleton's init can fail?**

**A:** `once.Do` never retries, so store the error next to the instance. Or use `sync.OnceValues` (Go 1.21+), which caches the value and the error.
