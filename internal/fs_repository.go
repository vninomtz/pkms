package internal

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type fsRepo struct {
	path string
	ext  string
	size int64
}

func NewFsRepo(dirRoot string) NodeRepository {
	return &fsRepo{
		path: dirRoot,
		ext:  ".md",
	}
}

func (r *fsRepo) buildName(filename string) string {
	name := filename + r.ext
	return filepath.Join(r.path, name)
}

func (r *fsRepo) Save(node Node) error {
	log.Println(r.buildName(node.Title))
	err := os.WriteFile(r.buildName(node.Title), []byte(node.Content), 0644)
	if err != nil {
		return err
	}
	return nil
}

func (r *fsRepo) GetNodes() ([]Node, error) {
	var nodes []Node
	err := filepath.Walk(r.path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error walking: %s\n", err)
			return err
		}
		if filepath.Ext(path) == r.ext {
			data, err := os.ReadFile(path)
			if err != nil {
				log.Printf("Error reading file: %s", err)
				return err
			}
			n := Node{
				Title:   strings.TrimSuffix(info.Name(), r.ext),
				Content: string(data),
			}
			nodes = append(nodes, n)
		}
		return nil
	})
	if err != nil {
		log.Println("Error reading nodes")
		return nil, err
	}
	return nodes, nil
}
