package main

import "fmt"

func rotate(matrix [][]int) {
	n := len(matrix)

	// Step 1: Transpose matrix
	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			matrix[i][j], matrix[j][i] = matrix[j][i], matrix[i][j]
		}
	}

	// Step 2: Reverse each row
	for i := 0; i < n; i++ {
		for j := 0; j < n/2; j++ {
			matrix[i][j], matrix[i][n-j-1] = matrix[i][n-j-1], matrix[i][j]
		}
	}
}

func main() {
	matrix := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}

	fmt.Println("Before rotation:")
	for _, row := range matrix {
		fmt.Println(row)
	}

	rotate(matrix)

	fmt.Println("After rotation:")
	for _, row := range matrix {
		fmt.Println(row)
	}
}
