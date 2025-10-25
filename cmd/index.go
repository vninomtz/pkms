package cmd

import (
	"flag"
	"fmt"
	"log"

	"github.com/vninomtz/pkms/internal/config"
	"github.com/vninomtz/pkms/internal/notes"
	"github.com/vninomtz/pkms/internal/store"
)

func IndexCommand(args []string) {
	cfg := config.New()
	cfg.Load()
	fs := flag.NewFlagSet("index", flag.ExitOnError)
	fs.Parse(args)

	srv := notes.New(cfg.NotesDir)

	notes, err := srv.GetAll()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Indexing %d notes\n", len(notes))
	st, err := store.New(cfg.SQLiteFile)
	if err != nil {
		log.Fatal(err)
	}
	success := 0
	for _, n := range notes {
		docId, err := st.SaveDocument(n.Entry)
		if err != nil {
			log.Printf("Error to index %s: %w\n", n.Entry.Filename, err)
			continue
		}
		_, err = st.SaveNote(n, docId)
		success++
	}
	fmt.Printf("%d docs indexed\n", success)
}
