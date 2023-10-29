package main

// ListNode Definition for singly-linked list.
type ListNode struct {
	Val  int
	Next *ListNode
}

func detectCycle(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return nil
	}

	// Step 1: Check if cycle exists and find meeting point
	slow := head
	fast := head
	meetingPoint := (*ListNode)(nil)

	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next

		if slow == fast {
			meetingPoint = slow
			break
		}
	}

	// Step 2: If no cycle exists, return null
	if meetingPoint == nil {
		return nil
	}

	// Step 3: Find the start of cycle
	slow = head
	for slow != meetingPoint {
		slow = slow.Next
		meetingPoint = meetingPoint.Next
	}
	return slow
}
