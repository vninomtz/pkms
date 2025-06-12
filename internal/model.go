package internal

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/adrg/frontmatter"
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

type INode interface {
	Name() string
}

type Node struct {
	Id       string   `json:"id"`
	Filename string   `json:"-"`
	Path     string   `json:"-"`
	Title    string   `json:"title"`
	Parent   string   `json:"-"`
	Bytes    []byte   `json:"-"`
	Content  string   `json:"content"`
	Size     int64    `json:"-"`
	Type     NodeType `json:"-"`
	Meta     Metadata `json:"metadata"`
	Html     string   `json:"html"`
}
type FileNode struct {
	Filename string
	Path     string
	Content  []byte
	Parent   string
	Size     int64
	Meta     Meta
}

func (n *FileNode) Name() string {
	return strings.TrimSuffix(n.Filename, filepath.Ext(n.Filename))
}
func (n *FileNode) ToMap() map[string]string {
	m := make(map[string]string)

	m["Name"] = n.Name()

	return m
}
func (n *FileNode) Links() ([]string, error) {
	links := []string{}
	scanner := bufio.NewScanner(bytes.NewBuffer(n.Content))
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		word := scanner.Text()
		u, err := url.Parse(word)
		if err == nil && u.Scheme != "" && u.Host != "" {
			links = append(links, word)
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		return links, err
	}

	return links, nil
}

type Filter struct {
	Field string
	Value string
}

type Meta struct {
	Title string
	Tags  any
}
type Metadata struct {
	Tags []string `json:"tags"`
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

func ExtractMetadata(raw []byte) ([]byte, Metadata, error) {
	var metadata Metadata
	var data map[string]interface{}
	content, err := frontmatter.Parse(bytes.NewReader(raw), &data)
	tags, ok := data["tags"]
	if err != nil {
		return content, metadata, err
	}
	if !ok {
		return content, metadata, nil
	}
	switch v := tags.(type) {
	case string:
		metadata.Tags = append(metadata.Tags, v)
	case []interface{}:
		var stringTags []string
		for _, item := range v {
			str, ok := item.(string)
			if ok {
				stringTags = append(stringTags, str)
			}
		}
		metadata.Tags = stringTags
	}

	return content, metadata, nil
}

func GetYaml(str string) string {
	res := ""
	queue := []string{}
	lines := strings.Split(str, "\n")
	open := false
	end := false
	for _, line := range lines {
		tmp := strings.TrimSpace(line)
		if tmp == "---" {
			queue = append(queue, tmp)
			if !open {
				open = true
			} else {
				end = true
				break
			}
		} else if open {
			queue = append(queue, line)
		}
	}

	if end {
		res = strings.Join(queue, "\n")
	}

	return res
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

func tagsToSet(tags []string) map[string]bool {
	set := map[string]bool{}

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
	ListAllTags() (map[string]int, error)
	GetByTitle(title string) (Node, error)
	Find([]Filter) ([]Node, error)
}
