package cmd

import (
	"flag"
	"fmt"
	"log"

	"github.com/vninomtz/pkms/internal"
	"github.com/vninomtz/pkms/internal/loader"
)

func SearchCommand(args []string) {
	fs := flag.NewFlagSet("search", flag.ExitOnError)
	filename := fs.String("filename", "", "Search note by filename")

	fs.Parse(args)

	notes_dir := internal.NotesPath()
	if notes_dir == "" {
		log.Fatalf("Notes directory no provided, set the %s env variable\n", internal.PKMS_NOTES_DIR)
	}
	FindByFilename(notes_dir, *filename)
}

func FindByFilename(dir, filename string) {
	load := loader.New(dir)
	err := load.Load()
	if err != nil {
		log.Fatal(err)
	}

	for _, doc := range load.Documents {
		if doc.Name() == filename {
			doc.Print()
			return
		}
	}

	fmt.Printf("Note with name %s not found \n", filename)
}
