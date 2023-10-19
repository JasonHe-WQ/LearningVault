package main

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func postorderTraversal(root *TreeNode) []int {
	if root == nil {
		return nil
	}
	res := make([]int, 0)
	stack := make([]*TreeNode, 0)
	var lastVisited *TreeNode // 记录上一次访问的节点
	for len(stack) > 0 || root != nil {
		if root != nil {
			stack = append(stack, root)
			root = root.Left
		} else {
			node := stack[len(stack)-1] // 查看栈顶元素但不出栈
			if node.Right != nil && lastVisited != node.Right {
				// 如果右子节点存在且未被访问过
				root = node.Right
			} else {
				res = append(res, node.Val)
				lastVisited = node
				stack = stack[:len(stack)-1] // 出栈
			}
		}
	}
	return res
}

// 对于多叉树，思路类似

type Node struct {
	Val      int
	Children []*Node
}

func postorder(root *Node) []int {
	if root == nil {
		return nil
	}

	var res []int
	stack := []*Node{root}
	var lastVisited *Node

	for len(stack) > 0 {
		// 查看栈顶元素但不出栈
		node := stack[len(stack)-1]

		// 如果该节点没有子节点，或者其子节点都被访问过
		if len(node.Children) == 0 || (lastVisited != nil && inChildren(lastVisited, node.Children)) {
			res = append(res, node.Val)
			lastVisited = node
			stack = stack[:len(stack)-1] // 出栈
		} else {
			// 将子节点逆序压入栈，这样栈顶总是下一个应该访问的节点
			for i := len(node.Children) - 1; i >= 0; i-- {
				stack = append(stack, node.Children[i])
			}
		}
	}

	return res
}

func inChildren(node *Node, children []*Node) bool {
	for _, child := range children {
		if node == child {
			return true
		}
	}
	return false
}
