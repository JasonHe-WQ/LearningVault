package main

import (
	"fmt"
	"sort"
)

func leastInterval(tasks []byte, n int) int {
	count := make([]int, 26) // 假设任务是从A-Z
	for _, t := range tasks {
		count[t-'A']++
	}
	// 按任务数量降序排序
	sort.Slice(count, func(i, j int) bool {
		return count[i] > count[j]
	})
	// 最多的任务数量
	maxVal := count[0] - 1
	idleSlots := maxVal * n
	// 减去其他任务填补的空闲槽
	for i := 1; i < 26; i++ {
		idleSlots -= min(count[i], maxVal)
	}
	// 如果idleSlots小于0，说明没有空闲槽
	if idleSlots < 0 {
		return len(tasks)
	}
	// 计算与maxVal相等的任务数量
	equalMaxCount := 0
	for _, c := range count {
		if c == count[0] {
			equalMaxCount++
		}
	}
	return len(tasks) + idleSlots + equalMaxCount - 1
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	tasks := []byte{'A', 'A', 'A', 'B', 'B', 'B'}
	n := 2
	fmt.Println(leastInterval(tasks, n))
}
