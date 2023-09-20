package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	text_A := strings.Split(scanner.Text(), " ")
	list_num_A := make([]int, len(text_A))
	for index := range text_A {
		list_num_A[index], _ = strconv.Atoi(text_A[index])
	}
	scanner.Scan()
	text_B := strings.Split(scanner.Text(), " ")
	list_num_B := make([]int, len(text_B))
	for index := range text_A {
		list_num_B[index], _ = strconv.Atoi(text_B[index])
	}
}
