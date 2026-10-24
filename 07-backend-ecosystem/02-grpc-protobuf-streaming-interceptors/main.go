// Shows gRPC ideas with the stdlib only: unary interceptor chain, deadline,
// server streaming (a channel here) and protobuf varint encoding.
// Run: go run ./07-backend-ecosystem/02-grpc-protobuf-streaming-interceptors
package main

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"time"
)

type Handler func(ctx context.Context, req string) (string, error)

// Interceptor has the same shape as grpc.UnaryServerInterceptor (simplified).
type Interceptor func(ctx context.Context, req string, next Handler) (string, error)

func chain(h Handler, ics ...Interceptor) Handler {
	for i := len(ics) - 1; i >= 0; i-- {
		ic, next := ics[i], h
		h = func(ctx context.Context, req string) (string, error) { return ic(ctx, req, next) }
	}
	return h
}

func logging(ctx context.Context, req string, next Handler) (string, error) {
	resp, err := next(ctx, req)
	fmt.Printf("  interceptor log: req=%q err=%v\n", req, err)
	return resp, err
}

func slowHandler(ctx context.Context, req string) (string, error) {
	select {
	case <-time.After(200 * time.Millisecond):
		return "hello " + req, nil
	case <-ctx.Done(): // honour the client deadline
		return "", ctx.Err()
	}
}

// listNumbers is a server-streaming RPC: many responses for one request.
func listNumbers(ctx context.Context, n int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out) // closing = io.EOF on the client
		for i := 1; i <= n; i++ {
			select {
			case out <- i:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

func main() {
	h := chain(slowHandler, logging)

	fmt.Println("=== unary call with enough deadline ===")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := h(ctx, "bob")
	cancel()
	fmt.Println("  resp:", resp, err)

	fmt.Println("=== deadline exceeded (gRPC code DeadlineExceeded) ===")
	ctx, cancel = context.WithTimeout(context.Background(), 50*time.Millisecond)
	_, err = h(ctx, "slow")
	cancel()
	fmt.Println("  is deadline error:", errors.Is(err, context.DeadlineExceeded))

	fmt.Println("=== server streaming ===")
	for v := range listNumbers(context.Background(), 3) {
		fmt.Println("  received", v)
	}

	fmt.Println("=== protobuf varint: field 1 = 150 ===")
	// tag = (field number << 3) | wire type 0 (varint) = 0x08
	buf := []byte{0x08}
	buf = binary.AppendUvarint(buf, 150)
	fmt.Printf("  bytes: % x (3 bytes; JSON {\"id\":150} is 10)\n", buf)
}
