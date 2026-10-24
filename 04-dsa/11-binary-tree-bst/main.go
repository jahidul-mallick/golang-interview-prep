// Binary tree classics: max depth (DFS), inorder, level order (BFS), validate BST.
// Run: go run ./04-dsa/11-binary-tree-bst
package main

import (
	"fmt"
	"math"
)

type Node struct {
	Val         int
	Left, Right *Node
}

func maxDepth(n *Node) int {
	if n == nil {
		return 0
	}
	return 1 + max(maxDepth(n.Left), maxDepth(n.Right))
}

// inorder visits left, node, right: sorted output for a BST.
func inorder(n *Node, out *[]int) {
	if n == nil {
		return
	}
	inorder(n.Left, out)
	*out = append(*out, n.Val)
	inorder(n.Right, out)
}

// levelOrder is BFS: process one queue snapshot (one level) per loop.
func levelOrder(root *Node) [][]int {
	var res [][]int
	if root == nil {
		return res
	}
	queue := []*Node{root}
	for len(queue) > 0 {
		var level []int
		for size := len(queue); size > 0; size-- {
			n := queue[0]
			queue = queue[1:]
			level = append(level, n.Val)
			if n.Left != nil {
				queue = append(queue, n.Left)
			}
			if n.Right != nil {
				queue = append(queue, n.Right)
			}
		}
		res = append(res, level)
	}
	return res
}

// isValidBST passes down the allowed (lo, hi) range, not just the parent.
func isValidBST(n *Node, lo, hi int) bool {
	if n == nil {
		return true
	}
	if n.Val <= lo || n.Val >= hi {
		return false
	}
	return isValidBST(n.Left, lo, n.Val) && isValidBST(n.Right, n.Val, hi)
}

func main() {
	//       4
	//      / \
	//     2   6
	//    / \
	//   1   3
	root := &Node{4, &Node{2, &Node{Val: 1}, &Node{Val: 3}}, &Node{Val: 6}}

	fmt.Println("=== max depth ===")
	fmt.Println(maxDepth(root))
	fmt.Println("=== level order ===")
	fmt.Println(levelOrder(root))
	fmt.Println("=== inorder (sorted for a BST) ===")
	var in []int
	inorder(root, &in)
	fmt.Println(in)
	fmt.Println("=== valid BST ===")
	fmt.Println(isValidBST(root, math.MinInt, math.MaxInt))
	root.Left.Right.Val = 5 // 5 is in the left subtree of 4 but greater than 4
	fmt.Println(isValidBST(root, math.MinInt, math.MaxInt))
}
