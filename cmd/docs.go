package cmd

import (
	"flag"
	"fmt"
	"log"

	"github.com/vninomtz/pkms/internal"
	"github.com/vninomtz/pkms/internal/loader"
)

func DocsCommand(args []string) {
	fs := flag.NewFlagSet("docs", flag.ExitOnError)
	load := fs.String("path", "", "Documents from a path")
	filename := fs.String("name", "", "Name of the document")
	export := fs.Bool("export", false, "Export document")
	fs.Parse(args)

	loader := loader.New(*load)
	err := loader.Load()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d documents loaded\n", len(loader.Documents))

	if *filename != "" {
		document := loader.FindByName(*filename)
		if document == nil {
			fmt.Printf("Document with name %s not found\n", *filename)
		} else {
			document.Print()
			if *export {
				parser := internal.NewTemplateParser("templates", "index.html")
				html, err := parser.Parse(*document)
				if err != nil {
					log.Fatal(err)
				}
				internal.WriteHtml(document.Name(), html)

			}
		}

	}

}
