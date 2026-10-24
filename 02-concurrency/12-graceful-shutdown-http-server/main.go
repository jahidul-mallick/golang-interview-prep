// Graceful shutdown: signal -> Shutdown(ctx) drains in-flight requests.
// A timer stands in for SIGTERM so the demo exits by itself.
// Run: go run ./02-concurrency/12-graceful-shutdown-http-server
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	// Real services stop on SIGINT/SIGTERM; here we also cancel after a delay.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	mux := http.NewServeMux()
	mux.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(300 * time.Millisecond) // simulated work
		fmt.Fprint(w, "slow request finished")
	})
	srv := &http.Server{Handler: mux}

	ln, err := net.Listen("tcp", "127.0.0.1:0") // random free port
	if err != nil {
		fmt.Println("listen:", err)
		return
	}
	serveErr := make(chan error, 1)
	go func() { serveErr <- srv.Serve(ln) }()

	fmt.Println("=== background worker stops on the same ctx ===")
	var wg sync.WaitGroup
	wg.Go(func() {
		<-ctx.Done()
		fmt.Println("worker: ctx cancelled, exiting")
	})

	fmt.Println("=== in-flight request when shutdown starts ===")
	reply := make(chan string, 1)
	go func() {
		resp, err := http.Get("http://" + ln.Addr().String() + "/slow")
		if err != nil {
			reply <- "client error: " + err.Error()
			return
		}
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		reply <- string(b)
	}()

	<-ctx.Done() // 1. signal received
	fmt.Println("=== shutdown: stop accepting, drain in-flight ===")
	shutCtx, done := context.WithTimeout(context.Background(), 2*time.Second)
	defer done()
	if err := srv.Shutdown(shutCtx); err != nil { // 2. drain with a deadline
		fmt.Println("shutdown error:", err)
		srv.Close() // last resort: drop connections
	}
	wg.Wait() // 3. wait for background goroutines
	fmt.Println("client got:", <-reply)
	fmt.Println("Serve returned ErrServerClosed:", errors.Is(<-serveErr, http.ErrServerClosed))
}
