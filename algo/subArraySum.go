package main

func subarraySum(nums []int, k int) int {
	var sum, count int
	var sumMap = make(map[int]int)
	sumMap[0] = 1
	for i := 0; i < len(nums); i++ {
		sum += nums[i]
		if j, ok := sumMap[sum-k]; ok {
			count += j
		}
		sumMap[sum] += 1
	}
	return count
}
