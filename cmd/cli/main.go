package main

import (
	"fmt"
	"log"

	"github.com/vninomtz/pkms/internal"
)

func main() {
	Bookmarks()
}

func Notes() {
	r := internal.NewNoteRepository("")
	notes, err := r.All()
	if err != nil {
		log.Fatalf("AllNotes: %v", err)
	}

	for i, n := range notes {
		fmt.Printf("%d: %s - %s\n", i, n.Type, n.Title)

	}
}
func Bookmarks() {
	r := internal.NewNoteRepository("")
	links, err := r.AllBookmarks()
	if err != nil {
		log.Fatalf("Bookmarks: %v", err)
	}

	for i, v := range links {
		fmt.Printf("%d: %s \n", i, v)

	}
}
