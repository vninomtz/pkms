package internal

import (
	"errors"
	"fmt"
)

type Searcher interface {
	File(filename string) (FileNode, error)
}

type searcher struct {
	nodes []FileNode
}

func NewSearcher(nodes []FileNode) *searcher {
	return &searcher{
		nodes: nodes,
	}
}

func (s *searcher) File(filename string) (FileNode, error) {
	for _, n := range s.nodes {
		if n.Name() == filename {
			return n, nil
		}
	}
	return FileNode{}, errors.New("NotFound")
}
func (s *searcher) GetBookmarks() ([]string, error) {
	var bookmarks []string

	for _, n := range s.nodes {
		links, err := n.Links()
		if err != nil {
			fmt.Println(err)
		} else {
			bookmarks = append(bookmarks, links...)

		}

	}
	return bookmarks, nil
}
