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
	bagCap, _ := strconv.Atoi(scanner.Text())
	scanner.Scan()
	textline := scanner.Text()
	weightStr := strings.Split(textline, " ")
	weight := make([]int, 0)
	for _, v := range weightStr {
		temp, _ := strconv.Atoi(v)
		weight = append(weight, temp)
	}
	scanner.Scan()
	textline = scanner.Text()
	valueStr := strings.Split(textline, " ")
	value := make([]int, 0)
	for _, v := range valueStr {
		temp, _ := strconv.Atoi(v)
		value = append(value, temp)
	}
	dp := make([][]int, len(value)+1)
	dp[0] = make([]int, bagCap+1) // 代表从前i个物品中选，背包容量为j时，最大的价值
	for i := 1; i <= len(value); i++ {
		// 一共有len(value)个物品，对应第i-1个物品
		dp[i] = make([]int, bagCap+1)
		for j := 1; j <= bagCap; j++ {
			dp[i][j] = dp[i-1][j]
			if j >= weight[i-1] {
				dp[i][j] = max(dp[i][j], dp[i-1][j-weight[i-1]]+value[i-1])
			}
		}
	}
	fmt.Println(dp[len(value)][bagCap])
}
