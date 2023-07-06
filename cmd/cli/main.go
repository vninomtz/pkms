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
	os.Setenv(DB_NOTES_PATH, "/Users/vnino/github.com/vninomtz/swe-notes/database/nodes.db")
	os.Setenv(DIR_NOTES_PATH, "/Users/vnino/github.com/vninomtz/vnotes/docs")
	DB_PATH := os.Getenv(DB_NOTES_PATH)
	DIR_PATH := os.Getenv(DIR_NOTES_PATH)

	create := flag.Bool("new", false, "New note")
	fs := flag.Bool("fs", false, "Use file system repo")
	list := flag.Bool("ls", false, "List notes")
	title := flag.String("name", "", "Note title")
	content := flag.String("c", "", "Note inline content")

	flag.Parse()

	repo := internal.NewSqliteNodeRepo(DB_PATH)

	if *fs {
		repo = internal.NewFsRepo(DIR_PATH)
	}

	srv := internal.NewNoteService(repo)

	if *create {
		if *content != "" {
			srv.New(*title, *content)
			return
		}
		input, err := ReadInputFromEditor()
		if err != nil {
			fmt.Println(err)
			panic(1)
		}
		srv.New(*title, string(input))
	}

	if *list {
		notes, err := srv.ListAll()
		if err != nil {
			log.Fatal(err)
		}
		for _, n := range notes {
			fmt.Printf("> %s\n", n.Title)
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
