package main

func dailyTemperatures(temperatures []int) []int {
	n := len(temperatures)
	answer := make([]int, n)
	stack := make([]int, 0)

	for i := 0; i < n; i++ {
		// 当栈不为空且当前温度大于栈顶温度
		for len(stack) > 0 && temperatures[i] > temperatures[stack[len(stack)-1]] {
			// 弹出栈顶元素
			lastIndex := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			// 计算天数
			answer[lastIndex] = i - lastIndex
		}
		// 将当前天的索引压入栈
		stack = append(stack, i)
	}
	return answer
}
