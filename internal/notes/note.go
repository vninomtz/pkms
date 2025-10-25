package notes

import (
	"fmt"
	"strings"
	"time"
)

const (
	// YYYY-MM-DD: 2022-03-23
	YYYYMMDD = "2006-01-02"
	// 24h hh:mm:ss: 14:23:20
	HHMMSS24h = "15:04:05"
)

type Note struct {
	Title   string
	Content string
	Type    string
	Created time.Time
	Updated time.Time
	Public  bool
	Tags    []string
	Links   []string
	Notes   []string
	Entry   Entry
}

type Entry struct {
	Filename  string
	Path      string
	Content   []byte
	Size      int64
	UpdatedAt time.Time
	Ext       string
}

func (n Note) Print() {
	fmt.Printf("Title: %s\n", n.Title)
	fmt.Printf("IsPublic: %v\n", n.Public)
	fmt.Printf("Type: %s\n", n.Type)
	fmt.Printf("Tags: %s\n", strings.Join(n.Tags, ","))
	fmt.Printf("Links: %d\n", len(n.Links))
	fmt.Println()
}
func newTimeId() string {
	t := time.Now()
	date := strings.Join(strings.Split(t.Format(YYYYMMDD), "-"), "")
	timeF := strings.Join(strings.Split(t.Format(HHMMSS24h), ":"), "")
	return fmt.Sprintf("%s%s", date, timeF)
}
