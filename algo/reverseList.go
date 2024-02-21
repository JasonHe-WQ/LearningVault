package main

// ListNode /**
type ListNode struct {
	Val  int
	Next *ListNode
}

func reverseList(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}

	start := &ListNode{Next: head, Val: 0}
	curr := head.Next
	head.Next = nil // 防止形成循环

	for curr != nil {
		next := curr.Next // 保存当前节点的下一个节点

		// 将当前节点插入到开始节点之后
		curr.Next = start.Next
		start.Next = curr

		// 移动curr指针
		curr = next
	}

	return start.Next
}
