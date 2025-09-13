package internal

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type fsRepo struct {
	path   string
	ext    string
	size   int64
	logger *log.Logger
}

func NewFsRepo(logger *log.Logger, dirRoot string) *fsRepo {
	return &fsRepo{
		path:   dirRoot,
		ext:    ".md",
		logger: logger,
	}
}

func (r *fsRepo) buildName(filename string) string {
	name := filename + r.ext
	return filepath.Join(r.path, name)
}
func (r *fsRepo) Restore() error {
	return nil
}

func (r *fsRepo) Save(node Node) error {
	err := os.WriteFile(r.buildName(node.Title), []byte(node.Content), 0644)
	if err != nil {
		return err
	}
	r.logger.Println(r.buildName(node.Title))
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
			dir := filepath.Dir(path)
			parent := filepath.Base(dir)
			raw, err := os.ReadFile(path)

			var meta Metadata
			content, meta, err := ExtractMetadata(raw)

			if err != nil {
				log.Printf("Error unmarshaling metadata of file %s: Error %s", path, err)
				content = raw
			}
			noteName := strings.TrimSuffix(info.Name(), r.ext)
			n := Node{
				Id:       noteName,
				Filename: info.Name(),
				Title:    noteName,
				Bytes:    raw,
				Path:     path,
				Parent:   parent,
				Content:  string(content),
				Meta:     meta,
				Size:     info.Size(),
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
