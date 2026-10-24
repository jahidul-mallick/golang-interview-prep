# 13 · Testing, Benchmarks & Fuzzing

[← Section index](../README.md) · Prev: [12 · Packages, Modules & Initialization Order](../12-packages-modules-init-order/README.md) · Next: [Iterators & range-over-func (Go 1.23+)](../14-iterators-range-over-func/README.md)

**Interview weight: 1.0%**  ·  Revise in: 10 min  ·  Run: `go run ./01-go-fundamentals/13-testing-benchmarks-fuzzing`

## In one line
Go tests are plain functions in `_test.go` files, usually table-driven with subtests, run by `go test`, with benchmarks and fuzzing built in.

## Must know
- Test file ends in `_test.go`; functions are `TestXxx(t *testing.T)`, `BenchmarkXxx(b *testing.B)`, `FuzzXxx(f *testing.F)`.
- Table-driven tests with `t.Run(name, ...)` subtests are the standard style.
- `t.Errorf` continues the test, `t.Fatalf` stops it. `t.Helper()` fixes the reported line.
- Mock by depending on small interfaces and passing in a fake; no mocking framework is required.
- `net/http/httptest` gives `NewRecorder` for handlers and `NewServer` for clients.
- Benchmarks: `go test -bench . -benchmem`; since Go 1.24 prefer `for b.Loop() { ... }`.
- Run `go test -race ./...` in CI; use `-cover` for coverage and `-count=1` to skip the cache.
- `t.Cleanup`, `t.TempDir()` and `t.Setenv` tidy up; `testing.Short()` plus `t.Skip` skips slow tests.

## Code to remember
```go
func TestReverse(t *testing.T) {
	tests := []struct{ name, in, want string }{
		{"empty", "", ""},
		{"ascii", "abc", "cba"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Reverse(tt.in); got != tt.want {
				t.Errorf("Reverse(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func BenchmarkReverse(b *testing.B) {
	for b.Loop() {
		Reverse("hello, world")
	}
}
```

## Common traps
- Comparing error strings instead of `errors.Is` / `errors.As`.
- Calling `t.Fatal` from a goroutine other than the test goroutine.
- Shared global state between tests, so they pass alone but fail together or with `-race`.
- Benchmarking code the compiler can optimize away, or timing setup work.

## Interview Q&A

**Q1. What is a table-driven test?**

**A:** A slice of cases (name, input, expected) looped with `t.Run`. Adding a case is one line, and each case reports separately.

**Q2. How do you test code that calls a database or an API?**

**A:** Depend on an interface and pass a fake in the test. For HTTP clients use `httptest.NewServer`.

**Q3. How do you run a benchmark and read it?**

**A:** `go test -bench . -benchmem`. It prints ns/op, B/op and allocs/op. Compare runs with `benchstat`.

**Q4. What is fuzzing?**

**A:** `go test -fuzz FuzzXxx` generates random inputs to find panics and broken properties, such as `Reverse(Reverse(s)) == s`. Failing inputs are saved in `testdata/`.

**Q5. How do you run tests in parallel and catch races?**

**A:** `t.Parallel()` lets a test run alongside other parallel tests. `go test -race` reports data races that actually happen during the run.

**Q6. `t.Error` vs `t.Fatal`?**

**A:** `t.Error` marks the test failed and carries on. `t.Fatal` marks it failed and stops that test at once, so use it when later steps need the earlier result.
