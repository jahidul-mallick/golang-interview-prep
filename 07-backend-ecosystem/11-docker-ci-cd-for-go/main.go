// Shows build metadata you inject at build time in CI/Docker.
// Run: go run ./07-backend-ecosystem/11-docker-ci-cd-for-go
// Inject a version: go run -ldflags "-X main.version=1.2.3" ./07-backend-ecosystem/11-docker-ci-cd-for-go
package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
)

// version is overwritten at link time with -ldflags "-X main.version=...".
var version = "dev"

func main() {
	fmt.Println("=== build info ===")
	fmt.Println("  version:", version)
	fmt.Println("  Go     :", runtime.Version(), runtime.GOOS+"/"+runtime.GOARCH)
	if bi, ok := debug.ReadBuildInfo(); ok {
		fmt.Println("  module :", bi.Main.Path)
	}

	fmt.Println("=== config comes from env vars (12-factor) ===")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("  PORT =", port)
}
