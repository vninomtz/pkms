package cmd

import (
	"flag"
	"log"

	"github.com/vninomtz/pkms/internal/loader"
)

func DocsCommand(args []string) {
	fs := flag.NewFlagSet("test", flag.ExitOnError)
	load := fs.String("path", "", "Load all documents from a path")

	fs.Parse(args)

	loader := loader.New(*load)
	err := loader.Load()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%d documents loaded\n", len(loader.Documents))
}
