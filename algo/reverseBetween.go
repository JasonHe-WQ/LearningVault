package main

func reverseBetween(head *ListNode, left int, right int) *ListNode {
	if head == nil {
		return nil
	}

	// 创建哑节点并初始化 prev 和 curr 指针
	dummy := &ListNode{Next: head}
	prev := dummy
	curr := head

	// 移动 prev 和 curr 到适当的起始位置
	for i := 0; i < left-1; i++ {
		prev = curr
		curr = curr.Next
	}

	// 保存 prev 和 curr 的位置
	con := prev
	tail := curr

	// 反转链表段
	for i := 0; i < right-left+1; i++ {
		tmp := curr.Next
		curr.Next = prev
		prev = curr
		curr = tmp // 这里指向了下一个节点
	}

	// 重新连接链表
	con.Next = prev
	tail.Next = curr

	return dummy.Next
}
