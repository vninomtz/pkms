package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/vninomtz/pkms/internal"
	"github.com/vninomtz/pkms/internal/loader"
	"github.com/vninomtz/pkms/internal/notes"
	"github.com/vninomtz/pkms/internal/store"
)

func DocsCommand(args []string) {
	fs := flag.NewFlagSet("docs", flag.ExitOnError)
	add := fs.Bool("add", false, "Add new document (note)")
	path := fs.String("path", "", "Documents from a path")
	filename := fs.String("get", "", "Name of the document")
	export := fs.Bool("export", false, "Export document")
	index := fs.Bool("index", false, "Index documents loaded into a sqlite db")
	fs.Parse(args)

	if *path == "" {
		*path = internal.NotesPath()
	}
	if *path == "" {
		log.Fatalf("Notes directory no provided, set the %s env variable\n", internal.PKMS_NOTES_DIR)
	}

	if *add && *path != "" {
		AddDocument(*path)
		return
	}

	if *index && *path != "" {
		IndexDocuments(*path)
		return
	}
	if *filename != "" {
		GetDocument(*filename, *export)
		return
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
	note, err := notes.ParseMarkdown(doc.Content)
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

func AddDocument(dir string) {
	input, err := ReadInputFromEditor()
	if err != nil {
		log.Fatal(err)
	}
	filename := loader.NewTimeId()
	path, err := internal.WriteNote(filename, input, dir)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Document created %s\n", path)

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
		note, err := notes.ParseMarkdown(doc.Content)
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
func ReadInputFromEditor() ([]byte, error) {
	file, err := os.CreateTemp(os.TempDir(), "pkms")
	if err != nil {
		return []byte{}, err
	}
	filename := file.Name()
	defer os.Remove(filename)

	if err = file.Close(); err != nil {
		return []byte{}, err
	}

	if err = OpenFileInEditor(filename); err != nil {
		return []byte{}, err
	}

	bytes, err := os.ReadFile(filename)
	if err != nil {
		return []byte{}, err
	}

	return bytes, nil
}

func OpenFileInEditor(filename string) error {
	editor := "vim"
	path, err := exec.LookPath(editor)
	if err != nil {
		return err
	}
	cmd := exec.Command(path, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
