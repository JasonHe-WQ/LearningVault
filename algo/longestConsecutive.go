package main

func longestConsecutive(nums []int) int {
	numSet := make(map[int]bool)
	longestStreak := 0

	// 将数组元素存入集合
	for _, num := range nums {
		numSet[num] = true
	}

	// 遍历集合
	for num := range numSet {
		if !numSet[num-1] {
			currentNum := num
			currentStreak := 1

			// 查找连续的数字
			for numSet[currentNum+1] {
				currentNum++
				currentStreak++
			}

			// 更新最长的连续序列长度
			if currentStreak > longestStreak {
				longestStreak = currentStreak
			}
		}
	}

	return longestStreak
}
