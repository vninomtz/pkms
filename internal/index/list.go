package index

import "fmt"

type Node struct {
	key  int
	next *Node
}

type LinkedList struct {
	head *Node
}

func NewLinkList() *LinkedList {
	return &LinkedList{
		head: nil,
	}
}

func (l *LinkedList) Insert(value int) {
	n := &Node{key: value}
	if l.head == nil {
		l.head = n
		return
	}
	if n.key < l.head.key {
		n.next = l.head
		l.head = n
		return
	}

	x := l.head

	for x != nil && x.next != nil && n.key > x.next.key {
		x = x.next
	}
	n.next = x.next
	x.next = n
}
func (l *LinkedList) Search(value int) *Node {
	n := l.head
	for n != nil && n.key != value {
		n = n.next
	}
	return n
}

func (l *LinkedList) Print() {
	n := l.head
	for n != nil {
		fmt.Printf("%d -> ", n.key)
		n = n.next
	}
	fmt.Println()
}
