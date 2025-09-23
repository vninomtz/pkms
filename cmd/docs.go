package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/vninomtz/pkms/internal"
	"github.com/vninomtz/pkms/internal/loader"
)

func DocsCommand(args []string) {
	fs := flag.NewFlagSet("docs", flag.ExitOnError)
	add := fs.Bool("add", false, "Add new document (note)")
	load := fs.String("path", "", "Documents from a path")
	filename := fs.String("name", "", "Name of the document")
	export := fs.Bool("export", false, "Export document")
	fs.Parse(args)

	if *add && *load != "" {
		AddDocument(*load)
		return
	}

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

func AddDocument(dir string) {
	input, err := ReadInputFromEditor()
	if err != nil {
		log.Fatal(err)
	}
	filename := internal.NewTimeId()
	path, err := internal.WriteNote(filename, input, dir)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Document created %s\n", path)

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
