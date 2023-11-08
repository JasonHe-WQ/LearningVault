package main

import "unicode"

func decodeString(s string) string {
	var numStack []int    // 存储数字的栈
	var strStack []string // 存储字符串的栈
	var currentNum int    // 当前的数字
	var currentStr string // 当前的字符串

	for _, char := range s {
		if unicode.IsDigit(char) {
			// 构建数字
			currentNum = currentNum*10 + int(char-'0')
		} else if char == '[' {
			// 遇到 '[', 把当前的数字和字符串压入栈
			numStack = append(numStack, currentNum)
			strStack = append(strStack, currentStr)
			// 重置当前的数字和字符串
			currentNum = 0
			currentStr = ""
		} else if char == ']' {
			// 遇到 ']', 把栈顶的数字和字符串弹出
			repeatCount := numStack[len(numStack)-1]
			numStack = numStack[:len(numStack)-1]
			// 重复栈顶字符串
			tempStr := strStack[len(strStack)-1]
			strStack = strStack[:len(strStack)-1]
			for i := 0; i < repeatCount; i++ {
				tempStr += currentStr
			}
			currentStr = tempStr
		} else {
			// 构建当前的字符串
			currentStr += string(char)
		}
	}

	return currentStr
}
