# 10 · Recursion & Backtracking

[← Section index](../README.md) · Prev: [Sorting Algorithms](../09-sorting-algorithms/README.md) · Next: [11 · Binary Tree & BST](../11-binary-tree-bst/README.md)

**Interview weight: 0.8%**  ·  Revise in: 8 min  ·  Run: `go run ./04-dsa/10-recursion-backtracking`

## In one line
If the question says "all subsets / permutations / combinations / valid arrangements", try a choice, recurse, undo the choice.

## Must know
- Template: choose, explore (recurse), un-choose.
- Base case first. Every recursive call must move toward it.
- Copy the path before saving it: `append([]int{}, path...)`. Otherwise later edits corrupt the result.
- Subsets: loop from `start`, recurse with `i+1` (tree below). Permutations: use a `used` array.
- Reuse allowed (combination sum): recurse with `i`, not `i+1`.
- Prune early (sort, then stop when `cands[i] > remain`).
- Complexity: subsets O(n * 2^n), permutations O(n * n!). Space O(n) recursion depth plus output.
- Duplicates in input: sort, then skip `nums[i] == nums[i-1]` at the same depth.

## Diagram
![Recursion tree of subsets of [1,2,3]: each node is one saved subset](../../assets/04-subsets-backtracking.svg)

## Code to remember
```go
var dfs func(start int)
dfs = func(start int) {
	res = append(res, append([]int{}, path...)) // save a copy
	for i := start; i < len(nums); i++ {
		path = append(path, nums[i]) // choose
		dfs(i + 1)                   // explore
		path = path[:len(path)-1]    // un-choose
	}
}
```

## Common traps
- Saving `path` without copying (all results end up identical).
- Forgetting to un-choose (or un-mark `used[i]`).
- Declaring a recursive closure with `:=`. Declare `var dfs func(...)` first.

## Interview Q&A

**Q1. How do you recognise backtracking?**

**A:** The answer is a list of all valid configurations, and you build each one step by step.

**Q2. Why copy the path when saving?**

**A:** A slice shares its backing array, so later append and truncate calls would overwrite the saved result.

**Q3. How do you avoid duplicate results?**

**A:** Sort the input and skip a value equal to the previous one at the same recursion level.

**Q4. How is it different from DP?**

**A:** Backtracking enumerates every answer. DP counts or optimises when subproblems repeat.

**Q5. Time and space complexity of subsets and permutations?**

**A:** Subsets O(n * 2^n), permutations O(n * n!), because each result is copied. Extra space is O(n) for the recursion depth.

**Q6. How do you speed up a backtracking solution?**

**A:** Prune: stop a branch as soon as it cannot lead to a valid answer. Sorting first (combination sum) makes the cut-off easy.
