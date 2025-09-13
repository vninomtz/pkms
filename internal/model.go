package internal

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/adrg/frontmatter"
)

type NodeType string

const (
	TypeNote    NodeType = "note"
	TypeQuote   NodeType = "quote"
	TypeWriting NodeType = "writing"
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

func (n *Node) Links() ([]string, error) {
	links := []string{}
	scanner := bufio.NewScanner(bytes.NewBuffer(n.Bytes))
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		word := scanner.Text()
		if IsUrl(word) {
			links = append(links, word)
		} else {
			re := regexp.MustCompile(`\((.*?)\)`)
			matches := re.FindAllStringSubmatch(word, -1)
			if len(matches) > 0 {
				for _, m := range matches {
					is := IsUrl(m[1])
					if is {
						links = append(links, m[1])
					}
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return links, err
	}
	return links, nil
}

func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
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
	Tags     []string `json:"tags"`
	Title    string
	IsPublic bool
	Type     string
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
	if err != nil {
		return content, metadata, err
	}
	tags, ok := data["tags"]
	if !ok {
		return content, metadata, nil
	}
	isPublic, ok := data["isPublic"]
	if ok {
		metadata.IsPublic = ToBool(isPublic)
	}
	noteType, ok := data["type"]
	if ok {
		metadata.Type = noteType.(string)
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
func ToBool(v interface{}) bool {
	if b, ok := v.(bool); ok {
		return b
	}
	return false
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

type NoteService interface {
	New(title, content string) (Node, error)
	ListAll() ([]Node, error)
	GetPublicNotes() ([]Node, error)
	ListAllTags() (map[string]int, error)
	GetByTitle(title string) (Node, error)
	Find([]Filter) ([]Node, error)
	GetBookmarks() ([]string, error)
}
