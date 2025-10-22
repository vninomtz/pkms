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
	fs.Parse(args)

	FindByFilename(cfg.NotesDir, *filename)
}

func FindByFilename(dir, filename string) {
	srv := notes.New(dir)

	n, err := srv.GetFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(n)

}
