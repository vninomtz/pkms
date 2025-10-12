package store

import (
	"strings"

	"github.com/vninomtz/pkms/internal/notes"
)

func (s *Store) SaveNote(note notes.Note, docId int64) (int64, error) {
	q := `INSERT INTO parse_documents(document_id,title,content,is_public,tags,type) VALUES (?,?,?,?,?,?)`

	isPublic := 0
	if note.Public {
		isPublic = 1
	}
	tags := ""
	if len(note.Tags) > 0 {
		tags = strings.Join(note.Tags, ",")
	}
	res, err := s.DB().Exec(q,
		docId,
		note.Title,
		note.Content,
		isPublic,
		tags,
		note.Type,
	)
	if err != nil {
		return 0, err
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	if len(note.Links) > 0 {
		for _, l := range note.Links {
			q = "INSERT INTO bookmarks(url, note_title) VALUES(?,?)"
			_, err = s.DB().Exec(q, l, note.Title)
			if err != nil {
				return lastID, err
			}
		}
	}
	return lastID, nil
}

func (s *Store) All() ([]notes.Note, error) {
	var res []notes.Note
	rows, err := s.DB().Query("SELECT title,content,is_public,tags,type FROM notes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var n notes.Note
		var tags string
		if err := rows.Scan(&n.Title, &n.Content, &n.Public, &tags, &n.Type); err != nil {
			return nil, err
		}
		n.Tags = strings.Split(tags, ",")
		res = append(res, n)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Store) AllBookmarks() ([]string, error) {
	var links []string
	rows, err := s.db.Query("SELECT DISTINCT url FROM bookmarks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var l string
		if err := rows.Scan(&l); err != nil {
			return nil, err
		}
		links = append(links, l)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return links, nil
}
