package main

type TrieNode struct {
	children map[rune]*TrieNode
	isEnd    bool
}

type Trie struct {
	root *TrieNode
}

// Constructor  Initialize your data structure here.
func Constructor() Trie {
	return Trie{root: &TrieNode{children: make(map[rune]*TrieNode)}}
}

// Insert  Inserts a word into the trie.
func (tmpTrie *Trie) Insert(word string) {
	node := tmpTrie.root
	for _, ch := range word {
		if _, ok := node.children[ch]; !ok {
			node.children[ch] = &TrieNode{children: make(map[rune]*TrieNode)}
		}
		node = node.children[ch]
	}
	node.isEnd = true
}

// Search  Returns if the word is in the trie.
func (tmpTrie *Trie) Search(word string) bool {
	node := tmpTrie.root
	for _, ch := range word {
		if _, ok := node.children[ch]; !ok {
			return false
		}
		node = node.children[ch]
	}
	return node.isEnd
}

// StartsWith Returns if there is any word in the trie that starts with the given prefix.
func (tmpTrie *Trie) StartsWith(prefix string) bool {
	node := tmpTrie.root
	for _, ch := range prefix {
		if _, ok := node.children[ch]; !ok {
			return false
		}
		node = node.children[ch]
	}
	return true
}
