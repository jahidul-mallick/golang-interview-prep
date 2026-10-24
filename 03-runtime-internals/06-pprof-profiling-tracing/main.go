// Two ways to get profiles: runtime/pprof in code, net/http/pprof over HTTP.
// Run: go run ./03-runtime-internals/06-pprof-profiling-tracing
package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	_ "net/http/pprof" // registers /debug/pprof/* on http.DefaultServeMux
	"runtime/pprof"
	"strings"
	"time"
)

func busy() int {
	n := 0
	for i := range 50_000_000 {
		n += i % 7
	}
	return n
}

func main() {
	fmt.Println("=== runtime/pprof: CPU profile into a buffer ===")
	var buf bytes.Buffer
	if err := pprof.StartCPUProfile(&buf); err != nil {
		fmt.Println("start failed:", err)
		return
	}
	busy()
	pprof.StopCPUProfile()
	fmt.Println("cpu profile written:", buf.Len() > 0, "(analyze: go tool pprof cpu.prof)")

	fmt.Println("=== net/http/pprof: goroutine profile over HTTP ===")
	ln, err := net.Listen("tcp", "127.0.0.1:0") // random free port
	if err != nil {
		fmt.Println(err)
		return
	}
	go http.Serve(ln, nil)
	time.Sleep(20 * time.Millisecond)
	resp, err := http.Get("http://" + ln.Addr().String() + "/debug/pprof/goroutine?debug=1")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(strings.SplitN(string(body), "\n", 2)[0]) // "goroutine profile: total N"

	fmt.Println("=== profiles available ===")
	for _, p := range pprof.Profiles() {
		fmt.Println(" ", p.Name())
	}
}
