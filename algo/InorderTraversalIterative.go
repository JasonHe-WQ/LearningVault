package main

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// InorderTraversalIterative 使用迭代和一个栈进行中序遍历
func inorderTraversal(root *TreeNode) []int {
	stack := []*TreeNode{}
	current := root
	var res []int
	for current != nil || len(stack) > 0 {
		// 先将所有左子节点压入栈
		for current != nil {
			stack = append(stack, current)
			current = current.Left
		}

		// 弹出栈顶元素并访问
		n := len(stack) - 1
		node := stack[n]
		stack = stack[:n]

		res = append(res, node.Val)

		// 移动到右子节点
		current = node.Right
	}
	return res
}
