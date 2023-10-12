package main

import "sort"

func merge(intervals [][]int) [][]int {
	if len(intervals) == 0 {
		return [][]int{}
	}

	// 按照每个区间的起始点排序
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	var result [][]int
	// 把第一个区间添加到结果集
	result = append(result, intervals[0])

	for i := 1; i < len(intervals); i++ {
		currStart, currEnd := intervals[i][0], intervals[i][1]
		prevEnd := result[len(result)-1][1]

		if currStart <= prevEnd {
			// 区间重叠，合并
			if currEnd > prevEnd {
				result[len(result)-1][1] = currEnd
			}
		} else {
			// 区间不重叠，直接添加到结果集
			result = append(result, intervals[i])
		}
	}

	return result
}
