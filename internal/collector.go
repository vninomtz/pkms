package internal

import (
	"bytes"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/adrg/frontmatter"
)

type collector struct {
	Root  string
	Ext   string
	nodes []FileNode
}

func NewCollector(root, ext string) *collector {
	if ext == "" {
		ext = ".md"
	}
	return &collector{
		Root: root,
		Ext:  ext,
	}
}

func (cfg *collector) ToMaps() []map[string]string {
	list := []map[string]string{}

	for _, n := range cfg.nodes {
		list = append(list, n.ToMap())
	}
	return list
}

func (cfg *collector) Collect() ([]FileNode, error) {
	nodes := []FileNode{}

	err := filepath.Walk(cfg.Root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error to access the path %s: %v\n", path, err)
			return err
		}

		if cfg.Ext != "" && filepath.Ext(path) == cfg.Ext {
			dir := filepath.Dir(path)
			parent := filepath.Base(dir)
			if parent == "." {
				parent = ""
			}
			b, err := os.ReadFile(path)

			var meta Meta
			content, err := frontmatter.Parse(bytes.NewReader(b), &meta)
			if err != nil {
				log.Printf("Error to parse the path %s: %v\n", path, err)
				return err
			}
			n := FileNode{
				Filename: info.Name(),
				Path:     path,
				Content:  content,
				Parent:   parent,
				Size:     info.Size(),
				Meta:     meta,
			}
			nodes = append(nodes, n)
		}
		return nil
	})

	cfg.nodes = nodes
	return nodes, err
}
