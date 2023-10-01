package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	firstLine := scanner.Text()
	sepFirstLine := strings.Split(firstLine, " ")
	n, _ := strconv.Atoi(sepFirstLine[0])
	m, _ := strconv.Atoi(sepFirstLine[1])
	matrix := make([][]int, n)
	for i := 0; i < n; i++ {
		scanner.Scan()
		line := scanner.Text()
		sepLine := strings.Split(line, " ")
		matrix[i] = make([]int, m)
		for j := 0; j < m; j++ {
			matrix[i][j], _ = strconv.Atoi(sepLine[j])
		}
	}
	sum := maxSquareMatrix(matrix)
	fmt.Print(sum)
}

func maxSquareMatrix(matrix [][]int) int {
	m, n := len(matrix), len(matrix[0])
	// 初始化前缀和矩阵
	prefixSum := make([][]int, m+1)
	for i := range prefixSum {
		prefixSum[i] = make([]int, n+1)
	}

	// 填充前缀和矩阵
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			prefixSum[i][j] = matrix[i-1][j-1] + prefixSum[i-1][j] + prefixSum[i][j-1] - prefixSum[i-1][j-1]
		}
	}

	maxSum := matrix[0][0]
	for l := 1; l <= min(m, n); l++ { // l为方阵的大小
		for i := 0; i+l-1 < m; i++ {
			for j := 0; j+l-1 < n; j++ {
				sum := prefixSum[i+l][j+l] - prefixSum[i][j+l] - prefixSum[i+l][j] + prefixSum[i][j]
				if sum > maxSum {
					maxSum = sum
				}
			}
		}
	}
	return maxSum
}
