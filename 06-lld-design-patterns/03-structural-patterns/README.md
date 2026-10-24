# 03 · Structural Patterns in Go

[← Section index](../README.md) · Prev: [Creational Patterns in Go](../02-creational-patterns/README.md) · Next: [04 · Behavioral Patterns in Go](../04-behavioral-patterns/README.md)

**Interview weight: 0.8%**  ·  Revise in: 6 min  ·  Run: `go run ./06-lld-design-patterns/03-structural-patterns`

## In one line
Adapter makes a foreign type fit your interface; decorator (middleware) wraps a handler with extra behaviour behind the same interface.

## Must know
- **Adapter**: a small struct that holds the third-party type and exposes the interface your code expects.
- **Decorator**: a function or struct that takes an interface value and returns the same interface, adding behaviour before or after.
- In `net/http`, middleware is `func(http.Handler) http.Handler`. This is the decorator pattern.
- Order matters. In `Chain(h, Logging, Auth)` logging runs first, then auth, then the handler.
- Middleware can short-circuit by not calling `next` (auth failure returns 401).
- Embedding a struct is the cheap decorator: override one method, the rest are promoted.
- Test handlers with `httptest.NewRecorder` and `httptest.NewRequest`, no network needed.

## Code to remember
```go
type Middleware func(http.Handler) http.Handler

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
```

## Common traps
- Forgetting to call `next.ServeHTTP`, so the request silently stops.
- Applying middleware in the wrong order (auth after the handler, logging that misses rejected requests).
- Writing the response header or body in a middleware and then again in the handler.
- Over-adapting: if you own the type, change its method names instead of writing an adapter.

## Interview Q&A

**Q1. What is the decorator pattern in Go?**

**A:** A wrapper with the same interface as the thing it wraps. HTTP middleware is the standard example.

**Q2. Adapter vs decorator?**

**A:** Adapter changes the interface to fit a caller. Decorator keeps the interface and adds behaviour.

**Q3. How do you chain middleware?**

**A:** Loop over the list in reverse and wrap: `h = ms[i](h)`, so the first middleware ends up outermost.

**Q4. Show a decorator that is not HTTP.**

**A:** `type loggingStore struct{ Store }` embeds the inner `Store`, overrides `Get` to log and then calls `s.Store.Get`. Everything else is promoted.

**Q5. When would you write an adapter?**

**A:** When a library type (a legacy logger, an SDK client) does not match the interface your service depends on.

**Q6. Where do recovery and logging middleware go in the chain?**

**A:** Outermost. Recovery catches panics from everything inside it, and logging then sees requests that auth rejected.
