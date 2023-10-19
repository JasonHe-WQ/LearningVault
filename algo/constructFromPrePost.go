package main

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func constructFromPrePost(preorder []int, postorder []int) *TreeNode {
	if len(preorder) == 0 || len(postorder) == 0 {
		return nil
	}

	root := &TreeNode{Val: preorder[0]}

	if len(preorder) == 1 {
		return root
	}

	leftRootVal := preorder[1]
	index := 0

	// 找到左子树根节点在postorder中的位置
	for i, v := range postorder {
		if v == leftRootVal {
			index = i
			break
		}
	}

	// 分割preorder和postorder数组，递归构建左右子树
	root.Left = constructFromPrePost(preorder[1:index+2], postorder[:index+1])
	root.Right = constructFromPrePost(preorder[index+2:], postorder[index+1:len(postorder)-1])

	return root
}
