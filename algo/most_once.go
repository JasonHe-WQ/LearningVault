package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func isPerfectSquare(num int) bool {
	sqrtNum := int(math.Sqrt(float64(num)))
	return sqrtNum*sqrtNum == num
}

func findNonSquare(a []int, l int, r int, x int) int {
	for i := l; i <= r; i++ {
		product := a[i] * x
		if !isPerfectSquare(product) {
			return a[i]
		}
	}
	return -1
}

func main() {
	scn := bufio.NewScanner(os.Stdin)
	scn.Scan()
	ntStr := scn.Text()
	ntstr := strings.Split(ntStr, " ")
	n, _ := strconv.Atoi(ntstr[0])
	t, _ := strconv.Atoi(ntstr[1])
	numList := make([]int, n)
	scn.Scan()
	numstr := scn.Text()
	numTokens := strings.Split(numstr, " ")
	for i, token := range numTokens {
		numList[i], _ = strconv.Atoi(token)
	}
	for i := 0; i < t; i++ {
		scn.Scan()
		lrx := scn.Text()
		lrxS := strings.Split(lrx, " ")
		l, _ := strconv.Atoi(lrxS[0])
		r, _ := strconv.Atoi(lrxS[1])
		x, _ := strconv.Atoi(lrxS[2])
		fmt.Println(findNonSquare(numList, l-1, r-1, x))
	}
}
