// Shows what a graph query does: friends-of-friends as a 2-hop traversal.
// Cypher equivalent: MATCH (me {name:'ann'})-[:FOLLOWS*2]->(x) RETURN x
// Run: go run ./07-backend-ecosystem/06-neo4j-graph-modeling-cypher
package main

import (
	"fmt"
	"sort"
)

// follows is an adjacency list: node -> nodes it points to.
var follows = map[string][]string{
	"ann": {"bob", "cat"},
	"bob": {"dan"},
	"cat": {"dan", "eve"},
	"dan": {"ann"},
}

// atDepth returns nodes first reached at exactly depth hops (BFS, no revisits).
func atDepth(start string, depth int) []string {
	seen := map[string]bool{start: true}
	level := []string{start}
	for d := 0; d < depth; d++ {
		var next []string
		for _, n := range level {
			for _, m := range follows[n] {
				if !seen[m] {
					seen[m] = true
					next = append(next, m)
				}
			}
		}
		level = next
	}
	sort.Strings(level)
	return level
}

func main() {
	fmt.Println("=== direct follows (1 hop) ===")
	fmt.Println(" ", atDepth("ann", 1))
	fmt.Println("=== friends of friends (2 hops) ===")
	fmt.Println(" ", atDepth("ann", 2))
	fmt.Println("=== 3 hops: nothing new, the edge back to ann is ignored ===")
	fmt.Println(" ", atDepth("ann", 3))
}
