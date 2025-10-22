package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/vninomtz/pkms/internal"
	"github.com/vninomtz/pkms/internal/loader"
	"github.com/vninomtz/pkms/internal/notes"
	"github.com/vninomtz/pkms/internal/store"
)

func FindCommand(args []string) {
	fs := flag.NewFlagSet("docs", flag.ExitOnError)
	path := fs.String("path", "", "Documents from a path")
	filename := fs.String("get", "", "Name of the document")
	export := fs.Bool("export", false, "Export document")
	index := fs.Bool("index", false, "Index documents loaded into a sqlite db")
	public := fs.Bool("public", false, "Get public docs")
	copyOutput := fs.String("copy", "", "Copy files readed from stdin to the output passed")
	fs.Parse(args)

	if *path == "" {
		*path = internal.NotesPath()
	}
	if *path == "" {
		log.Fatalf("Notes directory no provided, set the %s env variable\n", internal.PKMS_NOTES_DIR)
	}

	if *index && *path != "" {
		IndexDocuments(*path)
		return
	}
	if *filename != "" {
		GetDocument(*filename, *export)
		return
	}
	if *public {
		GetPublicDocs(*path)
		return
	}
	if *copyOutput != "" {
		DocsCopyTo(*copyOutput)
	}

}

func GetDocument(filename string, export bool) {
	st, err := store.New(internal.DatabasePath())
	if err != nil {
		log.Fatal(err)
	}
	doc, err := st.FindDocumetByName(filename)
	if err != nil {
		log.Fatal(err)
	}
	note, err := notes.Parse(doc.Content)
	if note.Title == "" {
		note.Title = doc.Name()
	}
	if err != nil {
		log.Fatal(err)
	}
	note.Print()
	if export {
		parser := internal.NewTemplateParser("templates", "index.html")
		html, err := parser.Parse(doc.Name(), doc.Content)
		if err != nil {
			log.Fatal(err)
		}
		internal.WriteHtml(doc.Name(), html)
	}
}
func GetPublicDocs(dir string) {
	loader := loader.New(dir)
	err := loader.Load()
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	for _, d := range loader.Documents {
		note, err := notes.Parse(d.Content)
		if err != nil {
			continue
		}
		if note.Public {
			fmt.Println(d.Path)
		}
	}
}
func DocsCopyTo(output string) {
	scanner := bufio.NewScanner(os.Stdin)

	read := 0
	copied := 0

	for scanner.Scan() {
		path := strings.TrimSpace(scanner.Text())
		if path == "" {
			continue
		}
		read++
		cmd := exec.Command("cp", path, output)
		_, err := cmd.Output()
		if err != nil {
			fmt.Println(err)
		} else {
			copied++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d copied of %d readed\n", copied, read)
}

func IndexDocuments(dir string) {
	log.Printf("Indexing %s\n", dir)
	loader := loader.New(dir)
	err := loader.Load()
	if err != nil {
		log.Fatal(err)
	}
	st, err := store.New(internal.DatabasePath())
	if err != nil {
		log.Fatal(err)
	}
	success := 0
	for _, doc := range loader.Documents {
		note, err := notes.Parse(doc.Content)
		if note.Title == "" {
			note.Title = doc.Name()
		}
		if err != nil {
			log.Printf("Error to parse %s: %w\n", doc.Filename, err)
			continue
		}
		docId, err := st.SaveDocument(doc)
		if err != nil {
			log.Printf("Error to index %s: %w\n", doc.Filename, err)
			continue
		}
		_, err = st.SaveNote(note, docId)
		success++
	}
	fmt.Printf("%d docs indexed of %d\n", success, len(loader.Documents))
}
