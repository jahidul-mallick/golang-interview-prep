// Quick runtime facts and a sync.Pool example used in runtime interview answers.
// Run: go run ./03-runtime-internals/08-runtime-interview-questions
package main

import (
	"bytes"
	"fmt"
	"runtime"
	"sync"
)

var bufPool = sync.Pool{New: func() any { return new(bytes.Buffer) }}

func render(name string) string {
	buf := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(buf)
	buf.Reset() // always reset before reuse
	buf.WriteString("hello, " + name)
	return buf.String()
}

func depth(n int) int { // deep recursion: the goroutine stack grows automatically
	if n == 0 {
		return 0
	}
	return 1 + depth(n-1)
}

func main() {
	fmt.Println("=== scheduler facts ===")
	fmt.Println("GOMAXPROCS:", runtime.GOMAXPROCS(0), " NumCPU:", runtime.NumCPU())
	fmt.Println("goroutines now:", runtime.NumGoroutine())

	fmt.Println("=== sync.Pool: reuse buffers ===")
	fmt.Println(render("pool"))

	fmt.Println("=== stack grows on demand ===")
	fmt.Println("depth(1_000_000) =", depth(1_000_000))

	fmt.Println("=== GC counters ===")
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Println("GC cycles so far:", m.NumGC)
}
