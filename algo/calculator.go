package main

import (
	"fmt"
)

func calculate(s string) int {
	var nums []int
	for i := 0; i < len(s); {
		if s[i] == ' ' {
			i++
			continue
		}
		if isDigit(s[i]) {
			j := i
			for j < len(s) && isDigit(s[j]) {
				j++
			}
			num := int(s[i]) - '0'
			nums = append(nums, num)
			i = j
		} else {
			switch s[i] {
			case '(':
				nums = append(nums, -1)
			case ')':
				sum := 0
				for nums[len(nums)-1] != -1 {
					sum += nums[len(nums)-1]
					nums = nums[:len(nums)-1]
				}
				nums = nums[:len(nums)-1]
				nums = append(nums, sum)
			case '+':
			}
			i++
		}
	}

	sum := 0
	for _, num := range nums {
		sum += num
	}
	return sum
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func main() {
	fmt.Println(calculate("1 + 1"))
	fmt.Println(calculate("(1+(1+1+1)+1)+1"))
}
