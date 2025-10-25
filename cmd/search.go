package cmd

import (
	"flag"
	"fmt"
	"log"

	"github.com/vninomtz/pkms/internal/config"
	"github.com/vninomtz/pkms/internal/notes"
)

func SearchCommand(args []string) {
	cfg := config.New()
	cfg.Load()

	fs := flag.NewFlagSet("search", flag.ExitOnError)
	filename := fs.String("filename", "", "Search note by filename")
	public := fs.Bool("public", false, "Search public notes")
	fs.Parse(args)

	if *filename != "" {
		FindByFilename(cfg.NotesDir, *filename)
	}
	if *public {
		SearchPublic(cfg.NotesDir)
	}
}

func FindByFilename(dir, filename string) {
	srv := notes.New(dir)

	n, err := srv.GetFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(n)
}

func SearchPublic(dir string) {
	srv := notes.New(dir)

	res, err := srv.GetPublic()
	if err != nil {
		log.Fatal(err)
	}

	for _, n := range res {
		fmt.Println(n.Entry.Path)
	}
}
