package main

import (
	"fmt"
	"math/rand"
	"time"
)

func quickSelect(nums []int, left, right, k int) int {
	pivotIndex := randInt(left, right)
	pivotIndex = partition(nums, left, right, pivotIndex)
	if pivotIndex == k {
		return nums[pivotIndex]
	} else if k < pivotIndex {
		return quickSelect(nums, left, pivotIndex-1, k)
	} else {
		return quickSelect(nums, pivotIndex+1, right, k)
	}
}

func partition(nums []int, left, right, pivotIndex int) int {
	pivot := nums[pivotIndex]
	// Move pivot to end
	nums[pivotIndex], nums[right] = nums[right], nums[pivotIndex]
	storeIndex := left
	for i := left; i < right; i++ {
		if nums[i] > pivot {
			nums[storeIndex], nums[i] = nums[i], nums[storeIndex]
			storeIndex++
		}
	}
	// Move pivot to its final place
	nums[right], nums[storeIndex] = nums[storeIndex], nums[right]
	return storeIndex
}

func randInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min+1)
}

func findKthLargest(nums []int, k int) int {
	return quickSelect(nums, 0, len(nums)-1, k-1)
}

func main() {
	nums := []int{3, 2, 1, 5, 6, 4}
	k := 2
	result := findKthLargest(nums, k)
	fmt.Println("The", k, "th largest element is:", result)
}
