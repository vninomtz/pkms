package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/vninomtz/pkms/internal"
)

func main() {
	//cmd := flag.String("cmd", "", "Command to execute")
	path := flag.String("p", "", "Path name")
	tFile := flag.String("t", "", "Template file to parse")
	preview := flag.String("pr", "", "Filename to preview if exists")
	out := flag.String("o", "test.html", "File output")
	build := flag.String("build", "", "Dir to build")
	flag.Parse()

	if err := Run(*path, *tFile, *out, *preview, *build); err != nil {
		log.Println(err)
	}
}

func Run(path, tFile, out, filename, build string) error {
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

		parser := internal.NewTemplateParser(tFile, "index.html")
		html, err := parser.Parse(n) //internal.ParseNodeToHTML(n, tFile)
		err = os.WriteFile(out, html, 0644)
	}
	if build != "" && tFile != "" {
		return Save(nodes, tFile, build)
	}

	return nil
}

func Save(nodes []internal.FileNode, tFile, outputDir string) error {
	os.RemoveAll(outputDir)
	err := os.Mkdir(outputDir, 0755)
	if err != nil {
		return nil
	}

	for _, n := range nodes {
		html, err := internal.ParseNodeToHTML(n, tFile)
		if err != nil {
			return nil
		}
		p := filepath.Join(outputDir, n.Parent+"_"+n.Name()+".html")
		log.Println(p)
		err = os.WriteFile(p, html, 0644)
		if err != nil {
			return nil
		}

	}
	return nil
}
