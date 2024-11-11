package internal

import "errors"

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
