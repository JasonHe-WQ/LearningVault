package main

func longestPalindromeSubseq(s string) (string, int) {
	n := len(s)
	dp := make([][]int, n)
	for i := range dp {
		dp[i] = make([]int, n)
		dp[i][i] = 1
	}

	for l := 2; l <= n; l++ {
		for i := 0; i <= n-l; i++ {
			j := i + l - 1
			if s[i] == s[j] {
				if l == 2 {
					dp[i][j] = 2
				} else {
					dp[i][j] = dp[i+1][j-1] + 2
				}
			} else {
				dp[i][j] = max(dp[i+1][j], dp[i][j-1])
			}
		}
	}

	// Reconstruct the longest palindromic subsequence
	i, j := 0, n-1
	var res string
	for i <= j {
		if s[i] == s[j] {
			if i == j {
				res += string(s[i])
			} else {
				res = string(s[i]) + res + string(s[j])
			}
			i++
			j--
		} else {
			if dp[i+1][j] > dp[i][j-1] {
				i++
			} else {
				j--
			}
		}
	}

	return res, dp[0][n-1]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
