// Trie: insert, search, prefix check.
// Run: go run ./04-dsa/15-trie
package main

import "fmt"

type node struct {
	next [26]*node
	end  bool
}

type Trie struct{ root node }

func (t *Trie) Insert(word string) {
	n := &t.root
	for _, c := range word {
		i := c - 'a'
		if n.next[i] == nil {
			n.next[i] = &node{}
		}
		n = n.next[i]
	}
	n.end = true
}

// walk returns the node reached by s, or nil.
func (t *Trie) walk(s string) *node {
	n := &t.root
	for _, c := range s {
		if n = n.next[c-'a']; n == nil {
			return nil
		}
	}
	return n
}

func (t *Trie) Search(w string) bool     { n := t.walk(w); return n != nil && n.end }
func (t *Trie) StartsWith(p string) bool { return t.walk(p) != nil }

func main() {
	fmt.Println("=== trie ===")
	t := &Trie{}
	t.Insert("apple")
	fmt.Println("search apple, app; startsWith app:", t.Search("apple"), t.Search("app"), t.StartsWith("app"))
	t.Insert("app")
	fmt.Println("search app after insert:", t.Search("app"))
}
