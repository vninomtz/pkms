package notes

import (
	"fmt"
	"strings"
	"time"
)

type Note struct {
	Title     string
	Content   string
	Type      string
	CreatedAt time.Time
	IsPublic  bool
	Tags      []string
	Links     []string
}

func (n Note) Print() {
	fmt.Printf("Title: %s\n", n.Title)
	fmt.Printf("IsPublic: %v\n", n.IsPublic)
	fmt.Printf("Type: %s\n", n.Type)
	fmt.Printf("Tags: %s\n", strings.Join(n.Tags, ","))
	fmt.Printf("Links: %d\n", len(n.Links))
	fmt.Println()
}
