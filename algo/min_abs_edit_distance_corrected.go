package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	textA := strings.Split(scanner.Text(), " ")
	listNumA := make([]int, len(textA))
	for index := range textA {
		listNumA[index], _ = strconv.Atoi(textA[index])
	}
	scanner.Scan()
	textB := strings.Split(scanner.Text(), " ")
	listNumB := make([]int, len(textB))
	for index := range textA {
		listNumB[index], _ = strconv.Atoi(textB[index])
	}
	m, n := len(listNumA), len(listNumB)
	dp := make([][]int, m+1)
	dp[0] = make([]int, n+1)
	for index := 1; index < m+1; index++ {
		dp[index] = make([]int, n+1)
		dp[index][0] = dp[index-1][0] + int(math.Abs(float64(listNumA[index-1])))
	}
	for index := 1; index < n+1; index++ {
		dp[0][index] = dp[0][index-1] + int(math.Abs(float64(listNumB[index-1])))
	}
	for i := 1; i < m+1; i++ {
		for j := 1; j < n+1; j++ {
			if listNumA[i-1] == listNumB[j-1] {
				dp[i][j] = dp[i-1][j-1]
				continue
			} else {
				dp[i][j] = min(dp[i-1][j]+int(math.Abs(float64(listNumA[i-1]))), dp[i][j-1]+int(math.Abs(float64(listNumB[j-1]))))
				dp[i][j] = min(dp[i][j], dp[i-1][j-1]+int(math.Abs(float64(listNumA[i-1]-listNumB[j-1]))))
			}
		}
	}
	fmt.Println(dp[m][n])
}
