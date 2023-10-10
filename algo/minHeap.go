package main

import (
	"container/heap"
	"fmt"
)

// An IntHeap is a min-heap of ints.
type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

// Push and Pop use pointer receivers because they modify the slice's length,
// not just its contents.
func (h *IntHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func findKthNumber(k int) int {
	h := &IntHeap{1}
	heap.Init(h)
	uniqueElements := make(map[int]bool)
	uniqueElements[1] = true

	val := 0
	count := 0

	for count < k {
		//先取出最小的元素
		val = heap.Pop(h).(int)
		if val != 1 && val != 2 && val != 3 && val != 5 && val != 7 {
			count++
		}

		for _, factor := range []int{2, 3, 5, 7} {
			n := val * factor
			if !uniqueElements[n] {
				heap.Push(h, n)
				uniqueElements[n] = true
			}
		}
	}

	return val
}

func main() {
	testValues := []int{1, 2, 3, 4, 5, 6, 7}
	for _, k := range testValues {
		result := findKthNumber(k)
		fmt.Printf("The %dth number is %d \n", k, result)
	}
}
