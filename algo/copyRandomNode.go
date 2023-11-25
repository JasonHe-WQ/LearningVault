package main

func copyRandomList(head *Node) *Node {
	if head == nil {
		return nil
	}

	// 第一步：复制每个节点，并将副本插入到原节点和其下一个节点之间
	for curr := head; curr != nil; curr = curr.Next.Next {
		copy := &Node{Val: curr.Val, Next: curr.Next}
		curr.Next = copy
	}

	// 第二步：更新复制节点的随机指针
	for curr := head; curr != nil; curr = curr.Next.Next {
		if curr.Random != nil {
			curr.Next.Random = curr.Random.Next
		}
	}

	// 第三步：拆分两个链表
	pseudoHead := &Node{}
	copyCurr := pseudoHead
	for curr := head; curr != nil; curr = curr.Next {
		next := curr.Next.Next
		copy := curr.Next

		copyCurr.Next = copy
		copyCurr = copy

		curr.Next = next
	}

	return pseudoHead.Next
}
