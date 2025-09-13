package internal

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	Extensions []string
	Dir        string
}

func (c Config) Include(ext string) bool {
	for _, e := range c.Extensions {
		if e == ext {
			return true
		}
	}
	return false
}

func Load(cfg Config) ([]Document, error) {
	docs := []Document{}
	err := filepath.Walk(cfg.Dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error to access the path %s: %v\n", path, err)
			return err
		}

		if info.IsDir() {
			return nil
		}
		if !cfg.Include(filepath.Ext(path)) {
			return nil
		}

		b, err := os.ReadFile(path)

		d := Document{
			Name:      info.Name(),
			Path:      path,
			Content:   b,
			Size:      info.Size(),
			UpdatedAt: info.ModTime(),
			Ext:       filepath.Ext(path),
		}
		docs = append(docs, d)
		return nil
	})

	return docs, err
}
