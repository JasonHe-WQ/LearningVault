package main

import (
	"sort"
)

func deleteCharacter(s string) int {
	freqMp := make(map[rune]int)
	for _, v := range s {
		freqMp[v]++
	}

	freqArr := make([]int, 0, len(freqMp))
	for _, v := range freqMp {
		freqArr = append(freqArr, v)
	}

	sort.Ints(freqArr)

	res := 0
	used := make(map[int]bool)

	for _, freq := range freqArr {
		if !used[freq] {
			used[freq] = true
			continue
		}

		newFreq := freq
		for used[newFreq] {
			newFreq--
			res++
		}

		if newFreq > 0 {
			used[newFreq] = true
		}
	}

	return res
}
