package internal

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type NoteRepository interface {
	Save(Note) error
	All() ([]Note, error)
}

type notesRepo struct {
	path string
	db   *sql.DB
}

func NewNoteRepository(path string) NoteRepository {
	return &notesRepo{
		path: path,
	}
}
func (r *notesRepo) Init() error {
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

func (r *notesRepo) Open() error {
	db, err := sql.Open("sqlite3", "file:tmp.db")
	if err != nil {
		return err
	}
	r.db = db
	return nil
}

func (r *notesRepo) Close() error {
	return r.db.Close()
}

func (repo *notesRepo) Exec(query string) error {
	_, err := repo.db.Exec(query)
	return err
}

func (r *notesRepo) Save(nt Note) error {
	err := r.Open()
	if err != nil {
		return err
	}
	q := `
	INSERT INTO notes (title,content,is_public,tags,type)
	VALUES (?,?,?,?,?)`

	isPublic := 0
	if nt.IsPublic {
		isPublic = 1
	}
	_, err = r.db.Exec(q,
		nt.Title,
		nt.Content,
		isPublic,
		strings.Join(nt.Tags, ","),
		nt.Type,
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

func (r *notesRepo) All() ([]Note, error) {
	err := r.Open()
	if err != nil {
		return nil, err
	}
	var notes []Note
	rows, err := r.db.Query("SELECT title,content,is_public,tags,type FROM notes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var n Note
		var tags string
		if err := rows.Scan(&n.Title, &n.Content, &n.IsPublic, &tags, &n.Type); err != nil {
			return nil, err
		}
		n.Tags = strings.Split(tags, ",")
		notes = append(notes, n)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	err = r.Close()
	if err != nil {
		return nil, err
	}
	return notes, nil
}
