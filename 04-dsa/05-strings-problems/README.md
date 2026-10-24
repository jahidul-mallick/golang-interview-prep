# String Problems

[← Section index](../README.md) · Prev: [04 · Hashing & HashMap Problems](../04-hashing-hashmap-problems/README.md) · Next: [06 · Linked List](../06-linked-list/README.md)

**Interview weight: 1.1%**  ·  Revise in: 8 min  ·  Run: `go run ./04-dsa/05-strings-problems`

## In one line
Palindrome, anagram or substring questions mean counting letters, two pointers, or expanding around a center, always mindful that Go strings are immutable bytes.

## Must know
- A Go string is an immutable byte sequence. `len(s)` is bytes; `for i, r := range s` yields runes.
- Convert to `[]rune` (or `[]byte` for ASCII-only input) to index or swap characters.
- Anagram: count letters with `[26]int` for lowercase ASCII. O(n) time, O(1) space.
- Palindrome check: two pointers from both ends, O(n). Filter and lowercase first if the problem asks.
- Longest palindromic substring: expand around each center (odd and even). O(n^2) time, O(1) space.
- Build strings with `strings.Builder`, split with `strings.Fields` or `strings.Split`, join with `strings.Join`.
- Longest substring problems usually use a sliding window (folder 03).
- Always ask: ASCII only, or Unicode? That decides bytes versus runes.

## Code to remember
```go
func expand(s string, l, r int) int { // palindrome length around a center
	for l >= 0 && r < len(s) && s[l] == s[r] {
		l--
		r++
	}
	return r - l - 1
}
// for each i: max(expand(s, i, i), expand(s, i, i+1))
```

## Common traps
- Indexing a string with non-ASCII text and cutting a character in half.
- `s += x` in a loop (O(n^2)): use `strings.Builder`.
- Forgetting the even-length center in palindrome expansion.
- Trying to assign `s[i] = 'x'`: strings are immutable, convert to a slice first.

## Interview Q&A

**Q1. Check if two strings are anagrams?**

**A:** Same length, then count letters up for one string and down for the other. All counts must end at zero. O(n).

**Q2. Why not just use `len(s)` for the character count?**

**A:** It returns bytes. For Unicode use `utf8.RuneCountInString(s)` or `len([]rune(s))`.

**Q3. Longest palindromic substring approach?**

**A:** Expand around each of the 2n-1 centers and keep the longest. O(n^2) time, O(1) space. Manacher gets O(n) but is rarely expected.

**Q4. How do you reverse the words in a sentence?**

**A:** `strings.Fields` to split on any whitespace, swap from both ends, `strings.Join` with a space.
