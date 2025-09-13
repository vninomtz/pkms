package internal

import (
	"database/sql"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type repository struct {
	path string
	db   *sql.DB
}

func NewRepository(path string) DocumentRepository {
	return &repository{
		path: path,
	}
}

func (r *repository) Init() error {
	bytes, err := os.ReadFile(filepath.Join(r.path, "schema.sql"))
	if err != nil {
		return err
	}
	err = r.Open()
	if err != nil {
		return err
	}

	err = r.Exec(string(bytes))
	if err != nil {
		return err
	}

	return r.Close()
}

func (r *repository) Open() error {
	db, err := sql.Open("sqlite3", "file:tmp.db")
	if err != nil {
		return err
	}
	r.db = db
	return nil
}
func (r *repository) Close() error {
	return r.db.Close()
}
func (repo *repository) Exec(query string) error {
	_, err := repo.db.Exec(query)
	return err
}

func (r *repository) Save(doc Document) error {
	err := r.Open()
	if err != nil {
		return err
	}

	q := `INSERT INTO documents(name, bytes, size, path, ext, updated_at) VALUES (?,?,?,?,?,?)`

	_, err = r.db.Exec(q,
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
	err = r.Close()
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) All() ([]Document, error) {
	err := r.Open()
	if err != nil {
		return nil, err
	}

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
	err = r.Close()
	if err != nil {
		return nil, err
	}
	return nodes, nil
}
