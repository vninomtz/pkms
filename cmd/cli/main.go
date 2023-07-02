package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/vninomtz/swe-notes/internal"
)

func main() {
	DB_PATH := "./database/nodes.db"

	repo := internal.NewRepository(DB_PATH)

	create := flag.Bool("new", false, "New note")
	list := flag.Bool("ls", false, "List notes")
	purge := flag.Bool("purge", false, "Purge notes")
	title := flag.String("t", "", "Note tile")
	content := flag.String("c", "", "Note content")

	flag.Parse()

	if *create {
		NewNote(repo, *title, *content)
	}

	if *list {
		ListNotes(repo)
	}
	if *purge {
		CheckError(repo.Clean())
	}
}

func NewNote(repo *internal.SQLiteRepository, title, content string) {
	if content == "" {
		log.Fatal("Empty content")
	}
	if title == "" {
		title = NewTimeId()
	}
	note := internal.Node{
		Title:       title,
		Description: content,
		Type:        "Note",
	}
	err := repo.Save(note)
	CheckError(err)
}

func ListNotes(repo *internal.SQLiteRepository) {
	notes, err := repo.GetNodes()
	CheckError(err)

	for _, n := range notes {
		fmt.Printf("%d %s: %s\n", n.Id, n.Title, n.Description)
	}
}

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
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
