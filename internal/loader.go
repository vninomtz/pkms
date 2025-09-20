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

func LoadDocuments(dir, db string) {
	repo := NewRepository(db)
	noteRepo := NewNoteRepository(db)

	err := repo.Init()
	if err != nil {
		log.Fatalf("Error to Init database: %v\n", err)
		return
	}

	docs, err := Load(Config{Extensions: []string{".md", ".png"}, Dir: dir})
	if err != nil {
		log.Fatalf("Error to load documents of dir: %v\n", err)
		return
	}
	size := len(docs)

	log.Printf("Documents loaded %d\n", size)

	for _, d := range docs {
		err := repo.Save(d)
		if err != nil {
			log.Fatalf("Error to save %s file: %v", d.Name, err)
			return
		}
		note, err := ParseDocument(d)
		if err != nil {
			log.Fatalf("Error to parse document %s: %v", d.Name, err)
			return
		}
		err = noteRepo.Save(note)
		if err != nil {
			log.Fatalf("Error to save note %s: %v", d.Name, err)
			return
		}
	}
}
