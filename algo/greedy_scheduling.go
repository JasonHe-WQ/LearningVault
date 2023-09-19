package main

import (
	"fmt"
	"sort"
)

func leastInterval(tasks []byte, n int) int {
	count := make([]int, 26)
	for _, t := range tasks {
		count[t-'A']++
	}
	sort.Slice(count, func(i int, j int) bool {
		return count[i] > count[j]
	})
	maxVal := count[0] - 1
	idleSlots := maxVal * n
	for i := 1; i < 26; i++ {
		idleSlots -= min(count[i], maxVal)
	}
	if idleSlots < 0 {
		idleSlots = 0
	}

	return len(tasks) + idleSlots
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
