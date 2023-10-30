package main

func swapPairs(head *ListNode) *ListNode {
	// 创建哑节点，并初始化 prev 指针
	dummy := &ListNode{Next: head}
	prev := dummy

	for head != nil && head.Next != nil {
		// 初始化 first 和 second 指针
		first := head
		second := head.Next

		// 进行交换
		prev.Next = second
		first.Next = second.Next
		second.Next = first

		// 移动 prev、first 和 second 指针
		prev = first
		head = first.Next
	}

	return dummy.Next
}
