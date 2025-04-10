package index

import (
	"fmt"
	"math/rand"
	"time"
)

type SkipListNode struct {
	key  int
	next []*SkipListNode
}

type SkipList struct {
	head     *SkipListNode
	maxLevel int
	p        float64
}

func (sl *SkipList) randomLevel() int {
	level := 1
	for rand.Float64() < sl.p && level < sl.maxLevel {
		level++
	}
	return level
}

func (sl *SkipList) Print() {
	fmt.Printf("Head len %d\n", len(sl.head.next))
	for i := 0; i < sl.maxLevel; i++ {
		if sl.head.next[i] != nil {
			n := sl.head.next[i]
			fmt.Printf("Node: %d, Children: %d\n", n.key, len(n.next))
		} else {
			fmt.Printf("No Node at %d\n", i)
		}
	}
}

func (sl *SkipList) Search(key int) bool {
	current := sl.head
	fmt.Printf("Start searching with head: %d, max line %d\n\n", current.key, len(current.next))
	for level := sl.maxLevel - 1; level >= 0; level-- {
		for current.next[level] != nil && current.next[level].key < key {
			current = current.next[level]
			fmt.Printf("Node: %d, Level: %d\n", current.key, level)
		}
	}
	current = current.next[0]

	return current != nil && current.key == key
}

func (sl *SkipList) Insert(key int) {
	update := make([]*SkipListNode, sl.maxLevel)
	current := sl.head

	fmt.Printf("- Insert key %d\n", key)

	for level := sl.maxLevel - 1; level >= 0; level-- {
		fmt.Printf("  - Start level %d\n", level)
		for current.next[level] != nil && current.next[level].key < key {
			current = current.next[level]
			fmt.Printf("    - Node with Key %d\n", current.key)
		}

		fmt.Printf("  - Update level %d with Node key %d\n\n", level, current.key)
		update[level] = current
	}

	level := sl.randomLevel()
	node := &SkipListNode{key: key, next: make([]*SkipListNode, level)}

	fmt.Printf("  - Insert Node %d at level %d\n", node.key, level)
	for i := 0; i < level; i++ {
		node.next[i] = update[i].next[i]
		update[i].next[i] = node

	}

}

func NewSkipList(maxLevel int, p float64) *SkipList {
	rand.Seed(time.Now().UnixNano())
	head := &SkipListNode{
		key:  -1,
		next: make([]*SkipListNode, maxLevel),
	}
	return &SkipList{head: head, maxLevel: maxLevel, p: p}
}
