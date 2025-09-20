package internal

import "time"

type Note struct {
	Title     string
	Content   string
	Type      string
	CreatedAt time.Time
	IsPublic  bool
	Tags      []string
	Links     []string
}
