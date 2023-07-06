package internal

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type noteService struct {
	repo NodeRepository
}

func NewNoteService(repo NodeRepository) NoteService {
	return &noteService{
		repo: repo,
	}
}

func (s *noteService) New(title, content string) error {
	if content == "" {
		return errors.New("Empty content")
	}
	if title == "" {
		title = NewTimeId()
	}
	note := Node{
		Title:   title,
		Content: content,
		Type:    "Note",
	}
	return s.repo.Save(note)
}

func (s *noteService) ListAll() ([]Node, error) {
	notes, err := s.repo.GetNodes()
	if err != nil {
		return nil, errors.New("Error to consult notes")
	}
	return notes, nil
}

const (
	// YYYY-MM-DD: 2022-03-23
	YYYYMMDD = "2006-01-02"
	// 24h hh:mm:ss: 14:23:20
	HHMMSS24h = "15:04:05"
)

func NewTimeId() string {
	t := time.Now()

	date := strings.Join(strings.Split(t.Format(YYYYMMDD), "-"), "")
	timeF := strings.Join(strings.Split(t.Format(HHMMSS24h), ":"), "")

	return fmt.Sprintf("%s%s", date, timeF)
}
