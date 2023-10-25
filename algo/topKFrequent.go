package main

import (
	"container/heap"
)

type ElementFrequency struct {
	element int
	freq    int
}
type MinHeap []ElementFrequency

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i].freq < h[j].freq }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MinHeap) Push(x interface{}) {
	*h = append(*h, x.(ElementFrequency))
}

func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func topKFrequent(nums []int, k int) []int {
	frequencyMap := make(map[int]int)
	for _, num := range nums {
		frequencyMap[num]++
	}

	h := &MinHeap{}
	heap.Init(h)

	for element, freq := range frequencyMap {
		heap.Push(h, ElementFrequency{element, freq})
		if h.Len() > k {
			heap.Pop(h)
		}
	}

	result := make([]int, k)
	for i := 0; i < k; i++ {
		result[k-i-1] = heap.Pop(h).(ElementFrequency).element
	}

	return result
}
