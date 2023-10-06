func maxProfit(prices []int, fee int) int {
	n := len(prices)
	if n < 2 {
		return 0
	}

	// 定义两个数组来存储DP值。
	// cash[i] 是在第i天或更早结束一次销售的任何交易序列的最大利润。
	// hold[i] 是在第i天或更早结束一次购买的任何交易序列的最大利润。
	cash := make([]int, n)
	hold := make([]int, n)

	// 初始化
	cash[0] = 0
	hold[0] = -prices[0]

	// 动态规划找到最大利润
	for i := 1; i < n; i++ {
		cash[i] = max(cash[i-1], hold[i-1]+prices[i]-fee)
		hold[i] = max(hold[i-1], cash[i-1]-prices[i])
	}

	return cash[n-1]
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}