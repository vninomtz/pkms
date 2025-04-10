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

type INode interface {
	Name() string
}

type Node struct {
	Id      int32
	Title   string
	Content string
	Type    NodeType
	checked bool
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
	err := yaml.Unmarshal([]byte(GetYaml(content)), &md)
	if err != nil {
		return Metadata{}, nil
	}
	return md, nil
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

func (m *Metadata) GetTags() []string {
	_tags := []string{}
	if m.Tags == "" {
		return _tags
	}
	for _, v := range strings.Split(m.Tags, ",") {
		_tags = append(_tags, strings.TrimSpace(v))
	}
	return _tags
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
	ListAllTags() (map[string]int, error)
	GetByTitle(title string) (Node, error)
	Find([]Filter) ([]Node, error)
}
