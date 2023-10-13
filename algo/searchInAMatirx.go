package main

import "fmt"

func binarySearchRow(matrix [][]int, target int) int {
	low, high := 0, len(matrix)-1
	ans := 0
	for low <= high {
		mid := low + (high-low)/2
		if matrix[mid][0] == target {
			return mid
		} else if matrix[mid][0] < target {
			ans = mid
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return ans
}

func binarySearchInRow(row []int, target int) bool {
	low, high := 0, len(row)-1
	for low <= high {
		mid := low + (high-low)/2
		if row[mid] == target {
			return true
		} else if row[mid] < target {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return false
}

func searchMatrix(matrix [][]int, target int) bool {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return false
	}

	row := binarySearchRow(matrix, target)
	if row == -1 {
		return false
	}

	return binarySearchInRow(matrix[row], target)
}

func main() {
	matrix := [][]int{
		{1},
		{3},
	}
	target := 3
	result := searchMatrix(matrix, target)
	fmt.Println("Result:", result) // Output should be true
}
