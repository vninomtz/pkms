package internal

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type NodeType string

const (
	TypeNote NodeType = "note"
)

const (
	// YYYY-MM-DD: 2022-03-23
	YYYYMMDD = "2006-01-02"
	// 24h hh:mm:ss: 14:23:20
	HHMMSS24h = "15:04:05"
)

type Node struct {
	Id      int32
	Title   string
	Content string
	Type    NodeType
	checked bool
}

type Filter struct {
	Field string
	Value string
}

type Metadata struct {
	Title string
	Tags  string
}

func NewNote(title, content string) (Node, error) {
	if content == "" {
		return Node{}, errors.New("Empty content")
	}
	if title == "" {
		title = NewTimeId()
	}
	return Node{
		Title:   title,
		Content: content,
		Type:    TypeNote,
	}, nil
}

func NewTimeId() string {
	t := time.Now()
	date := strings.Join(strings.Split(t.Format(YYYYMMDD), "-"), "")
	timeF := strings.Join(strings.Split(t.Format(HHMMSS24h), ":"), "")
	return fmt.Sprintf("%s%s", date, timeF)
}

func ExtractMetadata(content string) (Metadata, error) {
	md := Metadata{}
	err := yaml.Unmarshal([]byte(content), &md)
	if err != nil {
		return Metadata{}, nil
	}
	return md, nil
}

// Include tags check that all tags are included in the metadata
func (m *Metadata) IncludeTags(str string) bool {
	set := tagsToSet(m.Tags)
	for _, v := range strings.Split(str, ",") {
		if !set[strings.TrimSpace(v)] {
			return false
		}
	}
	return true
}

func tagsToSet(str string) map[string]bool {
	set := map[string]bool{}

	tags := strings.Split(str, ",")
	for _, v := range tags {
		set[strings.TrimSpace(v)] = true
	}
	return set
}

type NodeRepository interface {
	Save(Node) error
	GetNodes() ([]Node, error)
}

type NoteService interface {
	New(title, content string) (Node, error)
	ListAll() ([]Node, error)
	GetByTitle(title string) (Node, error)
	Find([]Filter) ([]Node, error)
}
