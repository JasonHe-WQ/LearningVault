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
	target, _ := strconv.Atoi(scanner.Text())
	scanner.Scan()
	firstLine := strings.Split(scanner.Text(), " ")
	list := make([]int, 0)
	for i := 0; i < len(firstLine); i++ {
		temp, _ := strconv.Atoi(firstLine[i])
		list = append(list, temp)
	}
	dp := make([]int, target+1)
	for index := range dp {
		dp[index] = math.MaxInt32
	}
	dp[0] = 0
	for _, value := range list {
		for i := value; i <= target; i++ {
			dp[i] = min(dp[i], dp[i-value]+1)
		}
	}
	fmt.Println(dp[target])
}
