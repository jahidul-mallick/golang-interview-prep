# 08 · Binary Search

[← Section index](../README.md) · Prev: [07 · Stack, Queue & Monotonic Stack/Deque](../07-stack-queue-monotonic/README.md) · Next: [Sorting Algorithms](../09-sorting-algorithms/README.md)

**Interview weight: 1.0%**  ·  Revise in: 8 min  ·  Run: `go run ./04-dsa/08-binary-search`

## In one line
Sorted input, a rotated sorted array, or "find the minimum value that works" with a monotone yes/no check means binary search, O(log n).

## Must know
- Precondition: the data, or the yes/no answer, is monotone (all false then all true).
- Classic: `lo <= hi`, `mid = lo + (hi-lo)/2`, return `mid` when found, else move `lo` or `hi` past `mid`.
- Boundary version (first index with `nums[i] >= target`): `lo < hi`, `hi = mid` or `lo = mid + 1`. It handles first/last occurrence and insert position.
- Rotated sorted array: one half is always sorted. Check if the target is inside it, else go to the other half.
- Binary search on the answer: define `canDo(x)`, search the smallest or largest x where it holds (Koko bananas, ship packages).
- Stdlib: `sort.Search(n, func(i int) bool {...})` and `slices.BinarySearch(s, target)` returning `(index, found)`.
- Complexity: O(log n) time, O(1) space, plus the cost of `canDo` for answer-search problems.

## Code to remember
```go
lo, hi := 0, len(nums) // first index with nums[i] >= target
for lo < hi {
	mid := lo + (hi-lo)/2
	if nums[mid] < target {
		lo = mid + 1
	} else {
		hi = mid
	}
}
return lo
```

## Common traps
- Infinite loop from `lo = mid` with `mid` rounding down; pair `lo = mid + 1` with `hi = mid`.
- Mixing the closed `[lo, hi]` and half-open `[lo, hi)` styles in one function.
- Using `(lo+hi)/2`, which can overflow on huge ranges.
- Applying it to data that is not monotone.

## Interview Q&A

**Q1. Search in a rotated sorted array?**

**A:** Compare `nums[lo]` with `nums[mid]` to find the sorted half, check if the target falls inside it, and narrow to the half that can hold it. O(log n).

**Q2. What is "binary search on the answer"?**

**A:** Search the range of possible answers using a check that is false below some threshold and true above it, for example the minimum eating speed.

**Q3. How do you find the first or last occurrence?**

**A:** Use the lower-bound template. First occurrence is `lowerBound(target)`; last is `lowerBound(target+1) - 1` for integers.

**Q4. Standard library helpers?**

**A:** `sort.Search` for custom predicates, `slices.BinarySearch` for sorted slices.

**Q5. Can you binary search an unsorted array?**

**A:** No, it needs sorted (or monotone) data. For one query just scan in O(n); for many queries sort once in O(n log n) and then search in O(log n).
