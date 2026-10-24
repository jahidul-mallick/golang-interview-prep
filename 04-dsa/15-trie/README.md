# Trie (Prefix Tree)

[← Section index](../README.md) · Prev: [14 · Dynamic Programming](../14-dynamic-programming/README.md) · Next: [LRU & LFU Cache Design](../16-lru-lfu-cache-design/README.md)

**Interview weight: 0.3%**  ·  Revise in: 3 min  ·  Run: `go run ./04-dsa/15-trie`
Optional / low priority

## In one line
Prefix questions (autocomplete, "starts with", word dictionary): store words as a tree of characters.

## Must know
- Each node has `children` (array of 26 or a map) and an `end` flag.
- Insert, search, and prefix check are all O(L), where L is the word length.
- Use `[26]*node` for lowercase letters only; use `map[rune]*node` for general text.
- Search needs `end == true`. A prefix check only needs the path to exist.

## Code to remember
```go
n := &t.root
for _, c := range word {
	i := c - 'a'
	if n.next[i] == nil {
		n.next[i] = &node{}
	}
	n = n.next[i]
}
n.end = true
```

## Common traps
- Treating a prefix as a full word (forgetting the `end` flag).
- Indexing `c - 'a'` with uppercase or non-letter input.

## Interview Q&A

**Q1. Trie vs hash map of words?**

**A:** A map finds exact words in O(L). A trie also answers prefix queries and lists words by prefix.

**Q2. Time and space complexity of a trie?**

**A:** Insert, search and prefix are O(L). Space is O(total characters), times 26 per node when you use an array of children.
