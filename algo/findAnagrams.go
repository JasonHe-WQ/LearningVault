package main

func findAnagrams(s string, p string) []int {
	// 存储结果的数组
	var result []int
	// 存储目标字符串p中字符出现次数的哈希表
	pCount := make(map[byte]int)
	// 存储当前窗口中字符出现次数的哈希表
	windowCount := make(map[byte]int)

	// 初始化目标字符串的哈希表
	for i := 0; i < len(p); i++ {
		pCount[p[i]]++
	}

	// 定义两个指针，用于表示滑动窗口的边界
	left, right := 0, 0

	// 移动右指针
	for right < len(s) {
		// 更新窗口哈希表
		windowCount[s[right]]++
		// 当窗口大小等于目标字符串长度时
		if right-left+1 == len(p) {
			// 比较两个哈希表
			if isMatch(pCount, windowCount) {
				result = append(result, left)
			}
			// 移动左指针，并更新窗口哈希表
			if windowCount[s[left]] == 1 {
				delete(windowCount, s[left])
			} else {
				windowCount[s[left]]--
			}
			left++
		}
		right++
	}
	return result
}

// 比较两个哈希表是否相等
func isMatch(a, b map[byte]int) bool {
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	for k, v := range b {
		if a[k] != v {
			return false
		}
	}
	return true
}
