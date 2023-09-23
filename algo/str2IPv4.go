package main

import (
	"fmt"
	"strconv"
)

func restoreIpAddresses(s string) []string {
	if len(s) < 4 || len(s) > 12 {
		return []string{}
	}
	var result []string
	track := make([]string, 0)
	backtrack(s, 0, track, &result)
	return result
}

func backtrack(s string, start int, track []string, result *[]string) {
	if len(track) == 4 && start == len(s) {
		*result = append(*result, track[0]+"."+track[1]+"."+track[2]+"."+track[3])
		return
	}
	if len(track) == 4 && start < len(s) {
		return
	}
	for i := 1; i <= 3 && start+i <= len(s); i++ {
		sub := s[start : start+i]
		num, _ := strconv.Atoi(sub)
		if (len(sub) > 1 && sub[0] == '0') || num > 255 {
			continue
		}
		track = append(track, sub)
		backtrack(s, start+i, track, result)
		track = track[:len(track)-1]
	}
}

func main() {
	s := "25525511135"
	fmt.Println(restoreIpAddresses(s))
}
