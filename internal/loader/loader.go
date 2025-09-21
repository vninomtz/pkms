package loader

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/vninomtz/pkms/internal"
)

type Loader struct {
	ext       []string
	path      string
	Documents []internal.Document
}

func New(path string) *Loader {
	return &Loader{
		ext:  []string{".md"},
		path: path,
	}
}
func (l *Loader) Load() error {
	docs := []internal.Document{}
	err := filepath.Walk(l.path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error to access path %s: %w", err)
		}

		if info.IsDir() {
			return nil
		}
		if !l.AllowedExt(filepath.Ext(path)) {
			return nil
		}

		b, err := os.ReadFile(path)

		d := internal.Document{
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

	l.Documents = docs
	return err
}
func (l *Loader) AllowedExt(ext string) bool {
	for _, e := range l.ext {
		if e == ext {
			return true
		}
	}
	return false
}
