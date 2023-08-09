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
	//DB_PATH := os.Getenv(DB_NOTES_PATH)
	DIR_PATH := os.Getenv(DIR_NOTES_PATH)
	logger := log.New(os.Stdout, "INFO: ", log.Ltime)

	//repo := internal.NewSqliteNodeRepo(DB_PATH)
	repo := internal.NewFsRepo(logger, DIR_PATH)

	srv := internal.NewNoteService(logger, repo)

	// CLI commander

	// Add command
	cmdAdd := flag.NewFlagSet("add", flag.ExitOnError)
	addContent := cmdAdd.String("c", "", "Inline note content")
	addTitle := cmdAdd.String("title", "", "Note title")

	// Ls command
	cmdLs := flag.NewFlagSet("ls", flag.ExitOnError)

	// Get Command
	cmdGet := flag.NewFlagSet("get", flag.ExitOnError)

	// Find Command
	cmdFind := flag.NewFlagSet("find", flag.ExitOnError)
	findTitle := cmdFind.String("n", "", "Note title")
	findTags := cmdFind.String("t", "", "Note tags")

	if len(os.Args) < 2 {
		logger.Fatalln("Expected one subcommand")
	}

	switch os.Args[1] {
	case "add":
		cmdAdd.Parse(os.Args[2:])
		var str string
		if *addContent == "" {
			input, err := ReadInputFromEditor()
			if err != nil {
				logger.Fatalln(err)
			}
			str = string(input)
		} else {
			str = *addContent
		}

		note, err := srv.New(*addTitle, str)
		if err != nil {
			logger.Fatalln(err)
		}
		fmt.Printf("%s", note.Title)
	case "get":
		cmdGet.Parse(os.Args[2:])
		var title string
		if len(cmdGet.Args()) > 0 {
			title = cmdGet.Arg(0)
		}
		note, err := srv.GetByTitle(title)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(note.Content)
	case "find":
		cmdFind.Parse(os.Args[2:])
		filters := []internal.Filter{}
		filters = append(filters, internal.Filter{Field: "title", Value: *findTitle})
		filters = append(filters, internal.Filter{Field: "tags", Value: *findTags})

		notes, err := srv.Find(filters)
		if err != nil {
			log.Fatal(err)
		}
		for i, n := range notes {
			fmt.Printf("%d. %s\n", i+1, n.Title)
		}
	case "ls":
		cmdLs.Parse(os.Args[2:])
		notes, err := srv.ListAll()
		if err != nil {
			logger.Fatalln(err)
		}
		for i, n := range notes {
			fmt.Printf("%d     %s\n", i+1, n.Title)
		}
	default:
		logger.Fatalln("Unknow subcommand")
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
