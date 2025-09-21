package internal

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type repository struct {
	path string
	db   *sql.DB
}

func NewRepository(path string, db *sql.DB) DocumentRepository {
	return &repository{
		path: path,
		db:   db,
	}
}

func (r *repository) Save(doc Document) error {
	q := `INSERT INTO documents(name, bytes, size, path, ext, updated_at) VALUES (?,?,?,?,?,?)`

	_, err := r.db.Exec(q,
		doc.Name,
		doc.Content,
		doc.Size,
		doc.Path,
		doc.Ext,
		doc.UpdatedAt.Format(time.RFC3339),
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) All() ([]Document, error) {
	var nodes []Document
	rows, err := r.db.Query("SELECT name, bytes, size, path, ext, updated_at FROM documents")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var n Document
		var updated string
		if err := rows.Scan(
			&n.Name,
			&n.Content,
			&n.Size,
			&n.Path,
			&n.Ext,
			&updated,
		); err != nil {
			return nil, err
		}
		n.UpdatedAt, err = time.Parse(time.RFC3339, updated)
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
