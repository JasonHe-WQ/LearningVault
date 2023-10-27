package main

import "fmt"

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

/**
 * Note: 类名、方法名、参数名已经指定，请勿修改
 *
 *
 *
 * @param nodes TreeNode类
 * @return int整型一维数组
 */
func getMostGold(nodes *TreeNode) []int {
	// write code here
	if nodes == nil {
		return []int{}
	}
	root := nodes
	path := make([]int, 0)
	_, path = dfs(root, path)
	l := len(path)
	for i := 0; i < l/2; i++ {
		path[i], path[l-i-1] = path[l-i-1], path[i]
	}
	return path
}
func dfs(node *TreeNode, inner_path []int) (int, []int) {
	if node == nil {
		return 0, make([]int, 0)
	}
	left, leftPath := dfs(node.Left, inner_path)
	fmt.Println(left, leftPath)
	right, rightPath := dfs(node.Right, inner_path)
	fmt.Println(right, rightPath)
	if left > right {
		inner_path = append(inner_path, leftPath...)
	} else {
		inner_path = append(inner_path, rightPath...)
	}
	inner_path = append(inner_path, node.Val)
	return max(left, right) + node.Val, inner_path
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
