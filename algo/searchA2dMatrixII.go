package main

//matrix = [[1,4,7,11,15],
//			[2,5,8,12,19],
//			[3,6,9,16,22],
//			[10,13,14,17,24],
//			[18,21,23,26,30]],
//			target = 5

// 从右上角开始搜索
// 如果当前元素等于目标值，则返回 true。
// 如果当前元素大于目标值，则移到左边一列。
// 如果当前元素小于目标值，则移到下边一行。
func searchMatrixII(matrix [][]int, target int) bool {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return false
	}

	row, col := 0, len(matrix[0])-1

	for row < len(matrix) && col >= 0 {
		if matrix[row][col] == target {
			return true
		} else if matrix[row][col] < target {
			row++
		} else {
			col--
		}
	}

	return false
}
