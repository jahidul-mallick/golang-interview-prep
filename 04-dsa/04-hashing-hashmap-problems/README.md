# 04 · Hashing & HashMap Problems

[← Section index](../README.md) · Prev: [03 · Sliding Window](../03-sliding-window/README.md) · Next: [String Problems](../05-strings-problems/README.md)

**Interview weight: 1.2%**  ·  Revise in: 10 min  ·  Run: `go run ./04-dsa/04-hashing-hashmap-problems`

## In one line
"Have I seen this before?", "count occurrences" or "find the complement" means a map or set that turns an O(n^2) search into O(n).

## Must know
- Lookup, insert and delete in a Go map are O(1) on average. A map is not safe for concurrent writes.
- Store what you need to look up: value to index (two sum), value to count (frequency), key to group (anagrams).
- Map keys must be comparable: strings, ints, structs and arrays like `[26]int` work; slices do not.
- Missing key returns the zero value, so `count[x]++` works without an initial check. Use `v, ok := m[k]` when absent and zero differ.
- Set: `map[T]struct{}`.
- Map iteration order is random. Sort the keys or keep a separate order slice if output order matters.
- Typical cost is O(n) time and O(n) space.
- Longest consecutive sequence: put everything in a set, start counting only where `v-1` is absent. O(n).

## Code to remember
```go
func twoSum(nums []int, target int) []int {
	idx := map[int]int{} // value -> index
	for i, v := range nums {
		if j, ok := idx[target-v]; ok {
			return []int{j, i}
		}
		idx[v] = i // insert after the lookup so an element cannot pair with itself
	}
	return nil
}
```

## Common traps
- Inserting before checking the complement, so one element is used twice.
- Using a slice as a map key (compile error). Use an array or a string key.
- Relying on map iteration order in output.
- Reading from a nil map is fine, but writing to it panics: initialise with `make` or a literal.

## Interview Q&A

**Q1. Two sum in O(n)?**

**A:** One pass with a map from value to index. For each number check whether `target - v` was already seen.

**Q2. How do you group anagrams?**

**A:** Use a canonical key: the sorted string, or a `[26]int` letter count. Words with the same key go in the same group.

**Q3. Why is a hash map lookup O(1) on average, not always?**

**A:** Collisions. Hashing spreads keys over buckets; the worst case degrades, but with a good hash it is constant on average.

**Q4. Top k frequent elements?**

**A:** Count with a map, then pick k by sorting (O(n log n)), a min-heap of size k (O(n log k)), or bucket sort (O(n)).

**Q5. Longest consecutive sequence in O(n) without sorting?**

**A:** Put all numbers in a set. Only start counting at a number whose `v-1` is absent, then walk up while `v+1` exists. Each number is visited once, so O(n).
