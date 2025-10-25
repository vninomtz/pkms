package notes

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type NoteService interface {
	New(content []byte) (string, error)
	GetFile(filename string) (Note, error)
	GetPublic() ([]Note, error)
	GetAll() ([]Note, error)
}

type noteService struct {
	BaseDir string
}

func New(dir string) NoteService {
	return &noteService{
		BaseDir: dir,
	}
}

func (n *noteService) New(content []byte) (string, error) {
	filename := fmt.Sprintf("%s.md", newTimeId())
	err := os.MkdirAll(n.BaseDir, 0755)
	if err != nil {
		return "", err
	}
	out := filepath.Join(n.BaseDir, filename)
	err = os.WriteFile(out, content, 0644)
	return out, err
}

func (n *noteService) GetFile(filename string) (Note, error) {
	entries, err := n.load()
	if err != nil {
		return Note{}, err
	}
	for _, doc := range entries {
		if strings.TrimSuffix(doc.Filename, ".md") == filename {
			n, err := Parse(doc.Content)
			if err != nil {
				return Note{}, err
			}
			return n, nil
		}
	}

	return Note{}, errors.New("File not found")
}
func (n *noteService) GetPublic() ([]Note, error) {
	var notes []Note

	entries, err := n.load()
	if err != nil {
		return notes, err
	}
	for _, d := range entries {
		note, err := Parse(d.Content)
		if err != nil {
			continue
		}
		if note.Public {
			note.Entry = d
			notes = append(notes, note)
		}
	}
	return notes, nil
}

func (n *noteService) GetAll() ([]Note, error) {
	var notes []Note

	entries, err := n.load()
	if err != nil {
		return notes, err
	}
	for _, d := range entries {
		note, err := Parse(d.Content)
		if err != nil {
			continue
		}
		note.Entry = d
		notes = append(notes, note)
	}
	return notes, nil
}

func (n *noteService) load() ([]Entry, error) {
	docs := []Entry{}
	err := filepath.Walk(n.BaseDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error to access path %s: %w", err)
		}

		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".md" {
			return nil
		}

		b, err := os.ReadFile(path)

		d := Entry{
			Filename:  info.Name(),
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
