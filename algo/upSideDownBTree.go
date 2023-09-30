/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */
package main

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func upsideDownBinaryTree(root *TreeNode) *TreeNode {
	if root == nil {
		return nil
	}
	newRoot := upsideDownBinaryTree(root.Left)
	if newRoot != nil {
		root.Left.Left = root.Right
		root.Left.Right = root
		root.Left = nil
		root.Right = nil
	}
	if newRoot != nil {
		return newRoot
	}
	return root
}
