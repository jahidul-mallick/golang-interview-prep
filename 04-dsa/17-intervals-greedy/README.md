# Intervals & Greedy

[← Section index](../README.md) · Prev: [LRU & LFU Cache Design](../16-lru-lfu-cache-design/README.md) · Next: [Bit Manipulation & Math](../18-bit-manipulation-math/README.md)

**Interview weight: 0.7%**  ·  Revise in: 6 min  ·  Run: `go run ./04-dsa/17-intervals-greedy`

## In one line
Ranges, meetings, or "can I reach the end": sort first, then sweep once and make the locally best choice.

## Must know
- Merge intervals: sort by start; if `cur.start <= last.end`, extend `last.end`, else append.
- Max non-overlapping (activity selection): sort by END, take an interval if it starts at or after the last end.
- Meeting rooms (can one person attend all): sort by start, check `next.start < prev.end`.
- Rooms needed at once: sort starts and ends separately, or use a min-heap of end times.
- Greedy = pick the best local choice and never revisit. It needs a reason why it is safe.
- Complexity: sorting dominates, O(n log n). The sweep is O(n).
- Jump game: track the farthest reachable index in one pass, O(n).

## Code to remember
```go
sort.Slice(iv, func(i, j int) bool { return iv[i][0] < iv[j][0] })
out := [][2]int{iv[0]}
for _, cur := range iv[1:] {
	last := &out[len(out)-1]
	if cur[0] <= last[1] {
		last[1] = max(last[1], cur[1])
	} else {
		out = append(out, cur)
	}
}
```

## Common traps
- Sorting by start when the problem needs sort by end (activity selection).
- Writing `last.end = cur.end` instead of `max(...)` (a nested interval shrinks the result).
- Touching intervals (`[1,2]`, `[2,3]`): check whether the problem counts them as overlapping.
- Empty input: `iv[0]` panics, so return early when `len(iv) == 0`.

## Interview Q&A

**Q1. How do you recognise an interval problem?**

**A:** Inputs are `[start, end]` pairs and the question is about overlap, merging, or scheduling. Sort first.

**Q2. Why sort by end for activity selection?**

**A:** The earliest-finishing interval leaves the most room for the rest.

**Q3. Greedy vs DP?**

**A:** If one local choice is provably safe, greedy is simpler and faster. If not, use DP.
