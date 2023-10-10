package main

import (
	"fmt"
)

func maxSubArray(nums []int) int {
	currentSum := nums[0]
	maxSum := nums[0]

	for i := 1; i < len(nums); i++ {
		currentSum = max(nums[i], currentSum+nums[i])
		maxSum = max(maxSum, currentSum)
	}

	return maxSum
}

func main() {
	// 示例 1
	nums1 := []int{-2, 1, -3, 4, -1, 2, 1, -5, 4}
	fmt.Println("示例 1 输出:", maxSubArray(nums1)) // 输出应为 6

	// 示例 2
	nums2 := []int{1}
	fmt.Println("示例 2 输出:", maxSubArray(nums2)) // 输出应为 1

	// 示例 3
	nums3 := []int{5, 4, -1, 7, 8}
	fmt.Println("示例 3 输出:", maxSubArray(nums3)) // 输出应为 23
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
