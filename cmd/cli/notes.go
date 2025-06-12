package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/vninomtz/pkms/internal"
)

type StoreType int

const (
	StoreFS = iota
	StoreSQLite
)

const (
	PKMS_STORE_PATH = "PKMS_STORE_PATH"
	STORE_TYPE      = "PKMS_STORE_TYPE"
)

func main() {
	logger := log.New(os.Stdout, "INFO: ", log.Ltime)

	source_path := os.Getenv(PKMS_STORE_PATH)
	store_type, err := strconv.Atoi(os.Getenv(STORE_TYPE))
	if err != nil {
		logger.Fatal("Error parsing STORE_TYPE env")
	}

	var repo internal.NodeRepository

	if store_type == StoreFS {
		repo = internal.NewFsRepo(logger, source_path)
	}
	if store_type == StoreSQLite {
		repo = internal.NewSqliteNodeRepo(source_path)
	}
	if repo == nil {
		logger.Fatal("Error: Unknown Store Type")
	}

	srv := internal.NewNoteService(logger, repo)

	// CLI commander

	// Add command
	cmdAdd := flag.NewFlagSet("add", flag.ExitOnError)
	addContent := cmdAdd.String("c", "", "Inline note content")
	addTitle := cmdAdd.String("title", "", "Note title")

	// Ls command
	cmdLs := flag.NewFlagSet("ls", flag.ExitOnError)
	lsTags := cmdLs.Bool("t", false, "List all tags")

	// Get Command
	cmdGet := flag.NewFlagSet("get", flag.ExitOnError)

	// Find Command
	cmdFind := flag.NewFlagSet("find", flag.ExitOnError)
	findTitle := cmdFind.String("n", "", "Note title")
	findTags := cmdFind.String("t", "", "Note tags")
	hasToExport := cmdFind.Bool("exp", false, "Export result")

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
		if *hasToExport {
			ExportNotes(notes)

		}
	case "ls":
		cmdLs.Parse(os.Args[2:])
		if *lsTags {
			tags, err := srv.ListAllTags()
			if err != nil {
				logger.Fatalln(err)
			}
			for k, v := range tags {
				fmt.Printf("- %s:%d\n", k, v)
			}
		} else {
			notes, err := srv.ListAll()
			if err != nil {
				logger.Fatalln(err)
			}
			for i, n := range notes {
				fmt.Printf("%d     %s\n", i+1, n.Title)
			}
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

func ExportNotes(notes []internal.Node) {
	res, err := json.Marshal(notes)
	if err != nil {
		log.Fatal("Unexpected error parsing json: ", err)
	}
	err = os.WriteFile("output.json", res, 0644)
	if err != nil {
		log.Fatal("Unexpected error writing file: ", err)
	}
}
