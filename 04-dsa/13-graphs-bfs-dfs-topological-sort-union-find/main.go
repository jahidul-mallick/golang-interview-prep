// Graph classics: BFS shortest path, DFS islands, topological sort, union-find.
// Run: go run ./04-dsa/13-graphs-bfs-dfs-topological-sort-union-find
package main

import "fmt"

// bfsDist: unweighted shortest distance from src (-1 = unreachable).
// dist doubles as the visited set.
func bfsDist(adj [][]int, src int) []int {
	dist := make([]int, len(adj))
	for i := range dist {
		dist[i] = -1
	}
	dist[src] = 0
	queue := []int{src}
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		for _, v := range adj[u] {
			if dist[v] == -1 {
				dist[v] = dist[u] + 1
				queue = append(queue, v)
			}
		}
	}
	return dist
}

// numIslands: DFS that sinks each island it visits (the grid is the visited set).
func numIslands(grid [][]byte) int {
	var sink func(r, c int)
	sink = func(r, c int) {
		if r < 0 || r >= len(grid) || c < 0 || c >= len(grid[0]) || grid[r][c] != '1' {
			return
		}
		grid[r][c] = '0'
		sink(r+1, c)
		sink(r-1, c)
		sink(r, c+1)
		sink(r, c-1)
	}
	count := 0
	for r := range grid {
		for c := range grid[r] {
			if grid[r][c] == '1' {
				count++
				sink(r, c)
			}
		}
	}
	return count
}

// topoOrder is Kahn's algorithm. Edge {a, b} means a must come before b.
// If fewer than n nodes come out, there is a cycle (returns nil).
func topoOrder(n int, edges [][2]int) []int {
	adj := make([][]int, n)
	indeg := make([]int, n)
	for _, e := range edges {
		adj[e[0]] = append(adj[e[0]], e[1])
		indeg[e[1]]++
	}
	var queue, order []int
	for i, d := range indeg {
		if d == 0 {
			queue = append(queue, i)
		}
	}
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		order = append(order, u)
		for _, v := range adj[u] {
			if indeg[v]--; indeg[v] == 0 {
				queue = append(queue, v)
			}
		}
	}
	if len(order) < n {
		return nil
	}
	return order
}

// DSU is union-find with path compression: near O(1) per operation.
type DSU struct{ parent []int }

func NewDSU(n int) *DSU {
	p := make([]int, n)
	for i := range p {
		p[i] = i
	}
	return &DSU{p}
}

func (d *DSU) Find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.Find(d.parent[x])
	}
	return d.parent[x]
}

// Union returns false if a and b were already connected.
func (d *DSU) Union(a, b int) bool {
	ra, rb := d.Find(a), d.Find(b)
	if ra == rb {
		return false
	}
	d.parent[ra] = rb
	return true
}

func main() {
	fmt.Println("=== BFS shortest distance from 0 ===")
	adj := [][]int{{1, 2}, {0, 3}, {0, 3}, {1, 2, 4}, {3}, {}}
	fmt.Println(bfsDist(adj, 0)) // node 5 is unreachable

	fmt.Println("=== number of islands ===")
	grid := [][]byte{
		[]byte("11000"),
		[]byte("11000"),
		[]byte("00100"),
		[]byte("00011"),
	}
	fmt.Println(numIslands(grid))

	fmt.Println("=== topological sort (course order) ===")
	fmt.Println(topoOrder(4, [][2]int{{0, 1}, {0, 2}, {1, 3}, {2, 3}}))
	fmt.Println("cycle detected:", topoOrder(2, [][2]int{{0, 1}, {1, 0}}) == nil)

	fmt.Println("=== union-find: connected components ===")
	d := NewDSU(5)
	comps := 5
	for _, e := range [][2]int{{0, 1}, {1, 2}, {3, 4}, {0, 2}} {
		if d.Union(e[0], e[1]) {
			comps--
		}
	}
	fmt.Println(comps)
}
