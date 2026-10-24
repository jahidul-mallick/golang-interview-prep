// Structural patterns: adapter and decorator (middleware).
// Run: go run ./06-lld-design-patterns/03-structural-patterns
package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

// Adapter: make a third-party type satisfy the interface your code wants.
type Logger interface{ Log(msg string) }

type legacyLogger struct{} // pretend this is a library we cannot change

func (legacyLogger) Write(level, text string) { fmt.Printf("  [%s] %s\n", level, text) }

type legacyAdapter struct{ l legacyLogger }

func (a legacyAdapter) Log(msg string) { a.l.Write("INFO", msg) }

// Decorator: wrap a handler with extra behaviour, same interface in and out.
type Middleware func(http.Handler) http.Handler

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("  log:", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Chain applies middlewares so the first one listed runs first.
func Chain(h http.Handler, ms ...Middleware) http.Handler {
	for i := len(ms) - 1; i >= 0; i-- {
		h = ms[i](h)
	}
	return h
}

func main() {
	fmt.Println("=== adapter ===")
	var log Logger = legacyAdapter{}
	log.Log("service started")

	fmt.Println("=== decorator / middleware ===")
	hello := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello")
	})
	h := Chain(hello, Logging, Auth)
	for _, token := range []string{"", "Bearer abc"} {
		req := httptest.NewRequest("GET", "/ping", nil)
		if token != "" {
			req.Header.Set("Authorization", token)
		}
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		fmt.Printf("  token=%q -> %d %q\n", token, rec.Code, rec.Body.String())
	}
}
