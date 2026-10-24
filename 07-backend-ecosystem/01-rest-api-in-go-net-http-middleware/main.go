// Shows a net/http middleware chain with Go 1.22+ method routing.
// Run: go run ./07-backend-ecosystem/01-rest-api-in-go-net-http-middleware
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

type Middleware func(http.Handler) http.Handler

// Chain wraps h so the first middleware listed runs first (outermost).
func Chain(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

// statusRecorder remembers the status code so logging can print it.
type statusRecorder struct {
	http.ResponseWriter
	code int
}

func (s *statusRecorder) WriteHeader(c int) {
	s.code = c
	s.ResponseWriter.WriteHeader(c)
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &statusRecorder{ResponseWriter: w, code: http.StatusOK}
		next.ServeHTTP(rec, r)
		fmt.Printf("  server log: %s %s -> %d\n", r.Method, r.URL.Path, rec.code)
	})
}

func recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if v := recover(); v != nil {
				http.Error(w, "internal error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer demo" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return // short-circuit: next is never called
		}
		next.ServeHTTP(w, r)
	})
}

type Item struct {
	Name string `json:"name"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /items/{id}", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"id": r.PathValue("id")})
	})
	mux.HandleFunc("POST /items", func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 1<<10) // cap body size
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		var it Item
		if err := dec.Decode(&it); err != nil {
			http.Error(w, "bad json: "+err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(it)
	})
	mux.HandleFunc("GET /boom", func(http.ResponseWriter, *http.Request) { panic("bug") })

	srv := httptest.NewServer(Chain(mux, logging, recoverer, auth))
	defer srv.Close()

	call := func(method, path, body, token string) {
		req, _ := http.NewRequest(method, srv.URL+path, strings.NewReader(body))
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("  error:", err)
			return
		}
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		fmt.Printf("  client got: %d %s\n", resp.StatusCode, strings.TrimSpace(string(b)))
	}

	fmt.Println("=== no token: auth short-circuits ===")
	call("GET", "/items/7", "", "")
	fmt.Println("=== path wildcard ===")
	call("GET", "/items/7", "", "demo")
	fmt.Println("=== strict JSON decoding ===")
	call("POST", "/items", `{"name":"pen"}`, "demo")
	call("POST", "/items", `{"nme":"typo"}`, "demo")
	fmt.Println("=== wrong method (mux gives 405) ===")
	call("DELETE", "/items/7", "", "demo")
	fmt.Println("=== panic is recovered ===")
	call("GET", "/boom", "", "demo")
}
