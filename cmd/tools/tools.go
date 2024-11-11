package main

import (
	"errors"
	"flag"
	"log"
	"os"

	"github.com/vninomtz/pkms/internal"
)

func main() {
	path := flag.String("p", "", "Path name")
	tFile := flag.String("t", "", "Template file to parse")
	preview := flag.String("pr", "", "Filename to preview if exists")
	out := flag.String("o", "test.html", "File output")
	flag.Parse()

	if err := Run(*path, *tFile, *out, *preview); err != nil {
		log.Println(err)
	}
}

func Run(path, tFile, out, filename string) error {
	if path == "" {
		return errors.New("Path is required")
	}
	collector := internal.NewCollector(path, "")
	nodes, err := collector.Collect()
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Path %s includes %d nodes\n", path, len(nodes))

	if filename != "" && out != "" && tFile != "" {
		searcher := internal.NewSearcher(nodes)
		n, err := searcher.File(filename)
		if err != nil {
			log.Printf("Error %v throw by File %s\n", err, filename)
			return err
		}

		html, err := internal.ParseNodeToHTML(n.Content, tFile)
		err = os.WriteFile(out, html, 0644)
	}

	return nil
}
