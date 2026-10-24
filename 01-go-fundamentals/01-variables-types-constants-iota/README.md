# 01 · Variables, Types, Constants & iota

[← Section index](../README.md) · Next: [02 · Arrays & Slices Internals](../02-arrays-slices-internals/README.md)

**Interview weight: 0.8%**  ·  Revise in: 5 min  ·  Run: `go run ./01-go-fundamentals/01-variables-types-constants-iota`

## In one line
Every Go variable starts at a useful zero value, `:=` declares (and can shadow), and constants are untyped until used, with `iota` generating enums.

## Must know
- Zero values: `0`, `""`, `false`, `nil` (pointer, slice, map, chan, func, interface). No uninitialised memory.
- `var x T` works anywhere; `x := v` only inside functions, and needs at least one new name on the left.
- `:=` in an inner scope creates a NEW variable (shadowing). `if _, err := f(); ...` does not touch an outer `err`.
- `int` and `uint` are 32 or 64 bits depending on the platform (64 on amd64). Use sized types (`int32`, `uint64`) for file and network formats.
- Integers wrap on overflow at runtime. Integer division truncates (`7/2 == 3`).
- Untyped constants have arbitrary precision (`1 << 100` compiles) and become typed only when assigned or used.
- `iota` starts at 0 in each `const` block and increments per line; use `1 << iota` for bit flags.
- A defined type (`type Weekday int`) is a different type from `int`; add a `String()` method for readable output.

## Code to remember
```go
type Weekday int

const (
	Sunday Weekday = iota // 0
	Monday                // 1
	Tuesday               // 2
)

func (d Weekday) String() string { return [...]string{"Sun", "Mon", "Tue"}[d] }
```

## Common traps
- Shadowing `err` with `:=` inside an `if` or `for`, so the outer `err` stays nil and the error is lost.
- Assuming `int` is 64 bits everywhere.
- `iota` enums without a reserved zero value, so an unset field looks like a valid first value.
- Comparing floats with `==`; money should use integer minor units (paise/cents), not `float64`.

## Interview Q&A

**Q1. What is a zero value and why does it matter?**

**A:** Every type has a default value, so variables are always initialised. Good APIs make the zero value usable, like `sync.Mutex` or `bytes.Buffer`.

**Q2. What is the difference between `var x = 5` and `x := 5`?**

**A:** Same result, but `:=` is function-scope only and needs a new variable on the left. At package level you must use `var`.

**Q3. Why does `const big = 1 << 100` compile but `var x int = 1 << 100` not?**

**A:** Constants are untyped and exact. The assignment fails because the value does not fit in `int`.

**Q4. How does `iota` work?**

**A:** It is the line index inside a `const` block, starting at 0. Lines with no expression repeat the previous one, so `KB = 1 << (10 * iota)` then `MB`, `GB` gives powers of 1024.

**Q5. What does shadowing look like in a real bug?**

**A:** `data, err := load()` inside an `if` block when `err` was declared outside. The inner `err` is set, the outer one never changes, and the function returns nil.

**Q6. Can you add an `int` and an `int64` without a conversion?**

**A:** No. Go has no implicit numeric conversion, so write `int64(x) + y`. Untyped constants like `5` adapt to either type.
