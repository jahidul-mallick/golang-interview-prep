// Shows package initialization order and build info.
// Run: go run ./01-go-fundamentals/12-packages-modules-init-order
// Try: go run -ldflags "-X main.version=1.2.3" ./01-go-fundamentals/12-packages-modules-init-order
package main

import (
	"fmt"
	"runtime/debug"
)

// Set at link time with -ldflags "-X main.version=..." (strings only).
var version = "dev"

var steps []string

func record(msg string) int {
	steps = append(steps, msg)
	return len(steps)
}

// a is declared first but depends on b, so b is initialized first.
var a = record("var a") + b
var b = record("var b")

func init() { record("init #1") }
func init() { record("init #2") } // many init functions are allowed

func main() {
	fmt.Println("=== initialization order (before main runs) ===")
	for i, s := range steps {
		fmt.Printf("%d. %s\n", i+1, s)
	}
	fmt.Println("a =", a, "b =", b)

	fmt.Println("=== link-time version ===")
	fmt.Println("version:", version)

	fmt.Println("=== build info ===")
	if bi, ok := debug.ReadBuildInfo(); ok {
		fmt.Println("module:", bi.Main.Path, "| go:", bi.GoVersion)
	}
}
