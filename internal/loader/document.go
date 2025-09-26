package loader

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

type Document struct {
	Filename  string
	Content   []byte
	Size      int64
	UpdatedAt time.Time
	Path      string
	Ext       string
}

func (d Document) Name() string {
	return strings.TrimSuffix(d.Filename, d.Ext)
}
func (d Document) Print() {
	fmt.Printf("Document: %s\n", d.Filename)
	fmt.Printf("Size: %d bytes\n", d.Size)
	fmt.Println("Content:")
	fmt.Println(string(d.Content))
}

func NewTimeId() string {
	t := time.Now()
	date := strings.Join(strings.Split(t.Format(YYYYMMDD), "-"), "")
	timeF := strings.Join(strings.Split(t.Format(HHMMSS24h), ":"), "")
	return fmt.Sprintf("%s%s", date, timeF)
}
