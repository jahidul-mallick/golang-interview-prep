// Linked list patterns: reverse, fast/slow pointers, dummy head.
// Run: go run ./04-dsa/06-linked-list
package main

import "fmt"

type ListNode struct {
	Val  int
	Next *ListNode
}

func fromSlice(vals ...int) *ListNode {
	dummy := &ListNode{}
	tail := dummy
	for _, v := range vals {
		tail.Next = &ListNode{Val: v}
		tail = tail.Next
	}
	return dummy.Next
}

func toSlice(head *ListNode) []int {
	var out []int
	for ; head != nil; head = head.Next {
		out = append(out, head.Val)
	}
	return out
}

// reverse flips the links in place with prev/cur/next. O(n) time, O(1) space.
func reverse(head *ListNode) *ListNode {
	var prev *ListNode
	for cur := head; cur != nil; {
		next := cur.Next
		cur.Next = prev
		prev, cur = cur, next
	}
	return prev
}

// hasCycle is Floyd's tortoise and hare: fast moves 2, slow moves 1.
// If they ever meet there is a cycle. O(n) time, O(1) space.
func hasCycle(head *ListNode) bool {
	slow, fast := head, head
	for fast != nil && fast.Next != nil {
		slow, fast = slow.Next, fast.Next.Next
		if slow == fast {
			return true
		}
	}
	return false
}

// mergeSorted merges two sorted lists. The dummy head removes the
// "is the result empty yet" special case.
func mergeSorted(a, b *ListNode) *ListNode {
	dummy := &ListNode{}
	tail := dummy
	for a != nil && b != nil {
		if a.Val <= b.Val {
			tail.Next, a = a, a.Next
		} else {
			tail.Next, b = b, b.Next
		}
		tail = tail.Next
	}
	if a != nil {
		tail.Next = a
	} else {
		tail.Next = b
	}
	return dummy.Next
}

// removeNthFromEnd: move fast n steps ahead, then walk both. When fast hits
// the end, slow is just before the node to delete. One pass.
func removeNthFromEnd(head *ListNode, n int) *ListNode {
	dummy := &ListNode{Next: head}
	slow, fast := dummy, dummy
	for range n {
		fast = fast.Next
	}
	for fast.Next != nil {
		slow, fast = slow.Next, fast.Next
	}
	slow.Next = slow.Next.Next
	return dummy.Next
}

func main() {
	fmt.Println("=== reverse ===")
	fmt.Println(toSlice(reverse(fromSlice(1, 2, 3, 4, 5))))

	fmt.Println("=== cycle detection ===")
	loop := fromSlice(1, 2, 3)
	loop.Next.Next.Next = loop // 3 points back to 1
	fmt.Println("with cycle:", hasCycle(loop), "without:", hasCycle(fromSlice(1, 2, 3)))

	fmt.Println("=== merge two sorted lists ===")
	fmt.Println(toSlice(mergeSorted(fromSlice(1, 3, 5), fromSlice(2, 4, 6))))

	fmt.Println("=== remove 2nd node from the end ===")
	fmt.Println(toSlice(removeNthFromEnd(fromSlice(1, 2, 3, 4, 5), 2)))
}
