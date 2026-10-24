// Shows what Kubernetes expects from a Go service: /healthz, /readyz,
// config from env vars, and draining on shutdown (what SIGTERM triggers).
// Run: go run ./07-backend-ecosystem/08-kubernetes-for-go-developers
package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func status(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		return "connection refused"
	}
	resp.Body.Close()
	return resp.Status
}

func main() {
	fmt.Println("=== config from env (ConfigMap/Secret become env vars) ===")
	fmt.Println("  LOG_LEVEL =", getenv("LOG_LEVEL", "info"))

	var ready atomic.Bool
	mux := http.NewServeMux()
	// liveness: "the process is alive". Keep it cheap; do not check dependencies.
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })
	// readiness: "send me traffic". 503 removes the pod from the Service.
	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) {
		if !ready.Load() {
			http.Error(w, "not ready", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(200)
	})

	ln, _ := net.Listen("tcp", "127.0.0.1:0")
	base := "http://" + ln.Addr().String()
	srv := &http.Server{Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	go srv.Serve(ln)

	fmt.Println("=== startup: alive but not ready yet ===")
	fmt.Println("  /healthz:", status(base+"/healthz"))
	fmt.Println("  /readyz :", status(base+"/readyz"))

	ready.Store(true) // caches warmed, DB connected
	fmt.Println("=== ready ===")
	fmt.Println("  /readyz :", status(base+"/readyz"))

	// On SIGTERM: fail readiness first, wait a moment, then Shutdown.
	// Real code: signal.NotifyContext(ctx, syscall.SIGTERM, os.Interrupt)
	fmt.Println("=== SIGTERM: fail readiness, drain, shutdown ===")
	ready.Store(false)
	fmt.Println("  /readyz :", status(base+"/readyz"))
	time.Sleep(100 * time.Millisecond) // real value: a few seconds for endpoints to update
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	fmt.Println("  shutdown err:", srv.Shutdown(ctx))
	fmt.Println("  after shutdown:", status(base+"/healthz"))
}
