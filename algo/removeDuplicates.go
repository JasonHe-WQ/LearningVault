package main

func removeDuplicates(nums []int) int {
	if len(nums) < 3 {
		return len(nums)
	}

	i := 1
	for j := 2; j < len(nums); j++ {
		if nums[j] != nums[i-1] {
			i++
			nums[i] = nums[j]
		}
	}
	return i + 1
}
