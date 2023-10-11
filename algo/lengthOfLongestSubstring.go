package main

import "fmt"

func lengthOfLongestSubstring(strIn string) string {
	tmpMap := make(map[rune]int)
	l := len(strIn)
	lenMAX := 0
	lenMAX1 := 0
	lenMAX2 := 0
	slow := 0
	for i := 0; i < l; i++ {
		valueI := rune(strIn[i])
		_, ok := tmpMap[valueI]
		if ok {
			tmpMap[valueI] += 1
			for len(tmpMap) <= i-slow {
				valueSlow := rune(strIn[slow])
				slow += 1
				slowV, _ := tmpMap[valueSlow]
				if slowV == 1 {
					delete(tmpMap, valueSlow)
				} else if slowV > 1 {
					tmpMap[valueSlow] -= 1
				}
			}
		} else {
			tmpMap[valueI] = 1
		}
		if i-slow+1 > lenMAX {
			lenMAX2 = i
			lenMAX1 = slow
			lenMAX = i - slow + 1
		}

	}
	return strIn[lenMAX1 : lenMAX2+1]
}

func main() {
	fmt.Println(lengthOfLongestSubstring("abbbcccdea"))
}

// abbbcccdea
