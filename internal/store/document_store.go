package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/vninomtz/pkms/internal/loader"
	"github.com/vninomtz/pkms/internal/notes"
)

func (s *Store) SaveDocument(doc notes.Entry) (int64, error) {
	q := `INSERT INTO documents(name, bytes, size, path, ext, updated_at) VALUES (?,?,?,?,?,?)`

	res, err := s.DB().Exec(q,
		doc.Filename,
		doc.Content,
		doc.Size,
		doc.Path,
		doc.Ext,
		doc.UpdatedAt.Format(time.RFC3339),
	)
	if err != nil {
		return 0, err
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

func (s *Store) Documents() ([]loader.Document, error) {
	var nodes []loader.Document
	rows, err := s.DB().Query("SELECT name, bytes, size, path, ext, updated_at FROM documents")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		n, err := s.parseDocument(rows)
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, n)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return nodes, nil
}
func (s *Store) FindDocumetByName(name string) (loader.Document, error) {
	var docs []loader.Document
	rows, err := s.DB().Query("SELECT name, bytes, size, path, ext, updated_at FROM documents where name = ?", name)
	if err != nil {
		return loader.Document{}, err
	}
	defer rows.Close()

	for rows.Next() {
		doc, err := s.parseDocument(rows)
		if err != nil {
			return loader.Document{}, err
		}

		docs = append(docs, doc)
	}
	if err := rows.Err(); err != nil {
		return loader.Document{}, err
	}
	if len(docs) == 0 {
		return loader.Document{}, fmt.Errorf("Not found")
	}
	return docs[0], nil
}

func (s *Store) parseDocument(rows *sql.Rows) (loader.Document, error) {
	var n loader.Document
	var updated string
	if err := rows.Scan(
		&n.Filename,
		&n.Content,
		&n.Size,
		&n.Path,
		&n.Ext,
		&updated,
	); err != nil {
		return loader.Document{}, err
	}
	d, err := time.Parse(time.RFC3339, updated)
	if err != nil {
		return loader.Document{}, err
	}
	n.UpdatedAt = d
	return n, err
}
