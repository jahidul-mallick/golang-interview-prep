# Bit Manipulation & Math

[← Section index](../README.md) · Prev: [Intervals & Greedy](../17-intervals-greedy/README.md)

**Interview weight: 0.3%**  ·  Revise in: 3 min  ·  Run: `go run ./04-dsa/18-bit-manipulation-math`
Optional / low priority

## In one line
"Find the one that appears once", "power of two", "count bits": XOR cancels pairs, and `n & (n-1)` clears the lowest set bit.

## Must know
- XOR: `x ^ x = 0`, `x ^ 0 = x`. XOR of all items leaves the unpaired one.
- `n & (n-1)` clears the lowest set bit. Loop until zero to count bits. Zero check: `n > 0 && n&(n-1) == 0` means power of two.
- `n & 1` is the last bit. `n >> 1` halves. `1 << k` is 2^k.
- Go has `&^` (AND NOT) to clear bits, and `math/bits.OnesCount` for popcount.
- GCD with Euclid: `a, b = b, a%b`, O(log min(a,b)).

## Code to remember
```go
for ; n != 0; n &= n - 1 {
	count++
}
```

## Common traps
- In Go `&` binds tighter than `==`, so `n&(n-1) == 0` is fine; in C-like languages it is not.
- Right-shifting negative ints keeps the sign bit. Use unsigned types for bit counting.

## Interview Q&A

**Q1. Find the number that appears once while the others appear twice?**

**A:** XOR everything. O(n) time, O(1) space.

**Q2. Check if n is a power of two?**

**A:** `n > 0 && n&(n-1) == 0`.
