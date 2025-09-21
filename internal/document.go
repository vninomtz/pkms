package internal

import "time"

type Document struct {
	Name      string
	Content   []byte
	Size      int64
	UpdatedAt time.Time
	Path      string
	Ext       string
}

type DocumentRepository interface {
	Save(Document) error
	All() ([]Document, error)
}
