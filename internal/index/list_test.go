package index

import (
	"testing"
)

func TestLinkedList(t *testing.T) {
	list := NewLinkList()

	list.Insert(1)
	list.Insert(12)
	list.Insert(21)
	list.Insert(4)
	list.Insert(3)
	list.Insert(22)
	list.Insert(23)
	list.Insert(0)
	list.Print()

	n := list.Search(21)
	if n == nil {
		t.Errorf("Expected found element %d\n", 21)
		return
	}
	n = list.Search(12)
	if n == nil {
		t.Errorf("Expected found element %d\n", 12)
		return
	}
	n = list.Search(99)
	if n != nil {
		t.Errorf("Expected not found element %d\n", 99)
		return
	}
}
