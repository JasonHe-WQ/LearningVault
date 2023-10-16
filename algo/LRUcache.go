package main

import (
	"container/list"
)

type LRUCache struct {
	capacity int
	cache    map[int]*list.Element
	lruList  *list.List
}

type Pair struct {
	key   int
	value int
}

func Constructor(capacity int) LRUCache {
	return LRUCache{
		capacity: capacity,
		cache:    make(map[int]*list.Element),
		lruList:  list.New(),
	}
}

func (this *LRUCache) Get(key int) int {
	if elem, found := this.cache[key]; found {
		this.lruList.MoveToFront(elem)
		return elem.Value.(Pair).value
	}
	return -1
}

func (this *LRUCache) Put(key int, value int) {
	if elem, found := this.cache[key]; found {
		this.lruList.MoveToFront(elem)
		elem.Value = Pair{key, value}
	} else {
		if this.lruList.Len() == this.capacity {
			delete(this.cache, this.lruList.Back().Value.(Pair).key)
			this.lruList.Remove(this.lruList.Back())
		}
		this.cache[key] = this.lruList.PushFront(Pair{key, value})
	}
}

/**
 * Your LRUCache object will be instantiated and called as such:
 * obj := Constructor(capacity);
 * param_1 := obj.Get(key);
 * obj.Put(key,value);
 */
