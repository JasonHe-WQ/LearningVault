package main

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func maxProfit(prices []int) int {
	n := len(prices)
	if n < 2 {
		return 0
	}

	buy := make([]int, n)
	sell := make([]int, n)
	frozen := make([]int, n)

	buy[0] = -prices[0]
	sell[0] = 0
	frozen[0] = 0

	for i := 1; i < n; i++ {
		buy[i] = max(buy[i-1], frozen[i-1]-prices[i])
		sell[i] = max(sell[i-1], buy[i-1]+prices[i])
		frozen[i] = max(frozen[i-1], sell[i-1])
	}

	return max(sell[n-1], frozen[n-1])
}
