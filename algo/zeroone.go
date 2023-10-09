package main

type Point struct {
	X int // 重量
	Y int // 价值
}

func maxIncomeProducts(products []*Point, months int) []*Point {
	l := len(products)
	dp := make([][]int, l+1)
	for i := 0; i <= l; i++ {
		dp[i] = make([]int, months+1)
	}
	dp2 := make([][]int, l+1)
	for i := 0; i <= l; i++ {
		dp2[i] = make([]int, months+1)
	}
	for i := 1; i <= l; i++ {
		for j := 1; j <= months; j++ {
			noTake := dp[i-1][j]
			take := 0
			if j >= products[i-1].X {
				take = dp[i-1][j-products[i-1].X] + products[i-1].Y
			}
			if take > noTake {
				dp[i][j] = take
				dp2[i][j] = 1
			} else {
				dp[i][j] = noTake
				dp2[i][j] = 0
			}
		}
	}
	var res []*Point
	m := months
	for i := l; i > 0; i-- {
		if dp2[i][m] == 1 {
			res = append(res, products[i-1])
			m -= products[i-1].X
		}
	}
	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}
	return res
}
