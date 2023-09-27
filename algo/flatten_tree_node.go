/*
*
* Definition for a binary tree node.

	type TreeNode struct {
	    Val int
	    Left *TreeNode
	    Right *TreeNode
	}
*/
package main

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func flatten(root *TreeNode) {
	if root == nil {
		return
	}

	// 先递归地展开左右子树
	flatten(root.Left)
	flatten(root.Right)

	// 保存右子树
	tmp := root.Right

	// 左子树移到右子树
	root.Right = root.Left
	root.Left = nil

	// 找到现在右子树的最后一个节点
	tail := root
	for tail.Right != nil {
		tail = tail.Right
	}

	// 将原来的右子树接到当前右子树的尾部
	tail.Right = tmp
}
