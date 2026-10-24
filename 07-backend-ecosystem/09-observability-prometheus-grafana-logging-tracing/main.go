// Shows the three pillars in stdlib form: structured logs (slog), a metric in
// Prometheus text format, and a trace/request id carried through context.
// Run: go run ./07-backend-ecosystem/09-observability-prometheus-grafana-logging-tracing
package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
)

type ctxKey struct{}

// counter is a tiny Prometheus Counter with one label (status code).
type counter struct {
	mu sync.Mutex
	n  map[string]int
}

func (c *counter) inc(code string) {
	c.mu.Lock()
	c.n[code]++
	c.mu.Unlock()
}

func (c *counter) write(w http.ResponseWriter) {
	c.mu.Lock()
	defer c.mu.Unlock()
	fmt.Fprintln(w, "# TYPE http_requests_total counter")
	for _, code := range []string{"200", "500"} {
		fmt.Fprintf(w, "http_requests_total{code=%q} %d\n", code, c.n[code])
	}
}

func main() {
	// JSON logs are easy for Loki/ELK to parse. Drop the time so output is stable.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	}))
	requests := &counter{n: map[string]int{}}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /work", func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID") // same id in logs, downstream calls and traces
		ctx := context.WithValue(r.Context(), ctxKey{}, id)
		code := "200"
		if r.URL.Query().Get("fail") != "" {
			code = "500"
			w.WriteHeader(http.StatusInternalServerError)
		}
		requests.inc(code)
		logger.InfoContext(ctx, "handled", "request_id", ctx.Value(ctxKey{}), "code", code)
	})
	mux.HandleFunc("GET /metrics", func(w http.ResponseWriter, r *http.Request) { requests.write(w) })

	srv := httptest.NewServer(mux)
	defer srv.Close()

	fmt.Println("=== logs: one JSON line per request ===")
	for _, path := range []string{"/work", "/work", "/work?fail=1"} {
		req, _ := http.NewRequest("GET", srv.URL+path, nil)
		req.Header.Set("X-Request-ID", "req-123")
		if resp, err := http.DefaultClient.Do(req); err == nil {
			resp.Body.Close()
		}
	}

	fmt.Println("=== metrics: what Prometheus scrapes from /metrics ===")
	resp, err := http.Get(srv.URL + "/metrics")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Print(string(body))
}
