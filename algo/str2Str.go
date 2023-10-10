package main

import "fmt"

//小红拿到了一个字符串s。她可以进行任意次以下操作
//
//选择字符串中的一个字母ch1和任意一个字母ch2 (ch2可以不在字符串中出现)，将字符串S中的所有ch1变成ch2。
//
//小红想知道，自己能否通过一些操作将字符串s变成t?

func solve() {
	var s, t string
	fmt.Scan(&s, &t)

	mp := make(map[rune]map[rune]bool)
	ok := true

	for i, ch1 := range s {
		ch2 := rune(t[i])

		if _, exists := mp[ch1]; !exists {
			mp[ch1] = make(map[rune]bool)
		}

		mp[ch1][ch2] = true

		if len(mp[ch1]) > 1 {
			ok = false
			break
		}
	}

	if ok {
		fmt.Println("Yes")
	} else {
		fmt.Println("No")
	}
}

func main() {
	var q int
	fmt.Scan(&q)

	for i := 0; i < q; i++ {
		solve()
	}
}
