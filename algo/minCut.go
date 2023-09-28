package main

import (
	"fmt"
	"math"
)

func main() {
	fmt.Println(minCut("abacba"))
}
func minCut(s string) int {
	l := len(s)
	isPalindrome := make([][]bool, l)
	for i := range isPalindrome {
		isPalindrome[i] = make([]bool, l)
		isPalindrome[i][i] = true
	}
	for length := 2; length <= l; length++ {
		for i := 0; i <= l-length; i++ {
			j := i + length - 1
			if s[i] == s[j] {
				if length == 2 {
					isPalindrome[i][j] = true
				} else {
					isPalindrome[i][j] = isPalindrome[i+1][j-1]
				}
			} else {
				isPalindrome[i][j] = false
			}
		}
	}
	dp := make([]int, l)
	for i := range dp {
		dp[i] = math.MaxInt32
	}
	for i := 0; i < l; i++ {
		if isPalindrome[0][i] {
			dp[i] = 0
		} else {
			for j := 0; j < i; j++ {
				if isPalindrome[j+1][i] {
					dp[i] = min(dp[i], dp[j]+1)
				}
			}
		}
	}
	return dp[l-1]
}
