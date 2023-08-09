package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/vninomtz/swe-notes/internal"
)

const (
	DB_NOTES_PATH  = "DB_NOTES_PATH"
	DIR_NOTES_PATH = "DIR_NOTES_PATH"
)

func main() {
	DB_PATH := os.Getenv(DB_NOTES_PATH)
	DIR_PATH := os.Getenv(DIR_NOTES_PATH)

	// commands
	create := flag.Bool("new", false, "New note")
	list := flag.Bool("ls", false, "List notes")
	get := flag.Bool("get", false, "Get note")
	find := flag.Bool("find", false, "Find notes")
	// conf
	fs := flag.Bool("fs", false, "Use file system repo")
	// values
	title := flag.String("name", "", "Note title")
	tags := flag.String("tags", "", "Note tags")
	content := flag.String("c", "", "Note inline content")

	flag.Parse()

	repo := internal.NewSqliteNodeRepo(DB_PATH)

	if *fs {
		repo = internal.NewFsRepo(DIR_PATH)
	}
	logger := log.New(os.Stdout, "INFO: ", log.Ltime)
	srv := internal.NewNoteService(logger, repo)

	if *create {
		var str string
		if *content == "" {
			input, err := ReadInputFromEditor()
			if err != nil {
				fmt.Println(err)
				panic(1)
			}
			str = string(input)
		} else {
			str = *content
		}

		note, err := srv.New(*title, str)
		if err != nil {
			fmt.Println(err)
			panic(1)
		}
		fmt.Printf("%s", note.Title)
	}

	if *list {
		notes, err := srv.ListAll()
		if err != nil {
			log.Fatal(err)
		}
		for i, n := range notes {
			fmt.Printf("%d     %s\n", i+1, n.Title)
		}
	}
	if *get {
		note, err := srv.GetByTitle(*title)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(note.Content)
	}
	if *find {
		filters := []internal.Filter{}
		filters = append(filters, internal.Filter{Field: "title", Value: *title})
		filters = append(filters, internal.Filter{Field: "tags", Value: *tags})

		notes, err := srv.Find(filters)
		if err != nil {
			log.Fatal(err)
		}
		for i, n := range notes {
			fmt.Printf("%d. %s\n", i+1, n.Title)
		}
	}
}

func ReadInputFromEditor() ([]byte, error) {
	file, err := ioutil.TempFile(os.TempDir(), "swenotes")
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
