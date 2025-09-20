package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	// "strconv"

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
	PKMS_SHARE_URL  = "PKMS_SHARE_URL"
)

func GetNotes() {
	r := internal.NewNoteRepository("")
	notes, err := r.All()
	if err != nil {
		log.Fatalf("AllNotes: %v", err)
	}

	for i, n := range notes {
		fmt.Printf("%d: %s - %s\n", i, n.Type, n.Title)

	}
}
func GetBookmarks() {
	r := internal.NewNoteRepository("")
	links, err := r.AllBookmarks()
	if err != nil {
		log.Fatalf("Bookmarks: %v", err)
	}

	for i, v := range links {
		fmt.Printf("%d: %s \n", i, v)

	}
}

func notes() {
	logger := log.New(os.Stdout, "INFO: ", log.Ltime)

	// source_path := os.Getenv(PKMS_STORE_PATH)
	// store_type, err := strconv.Atoi(os.Getenv(STORE_TYPE))
	// if err != nil {
	// 	logger.Fatal("Error parsing STORE_TYPE env")
	// }

	// var repo internal.NodeRepository
	//
	// if store_type == StoreFS {
	// 	repo = internal.NewFsRepo(logger, source_path)
	// }
	// if store_type == StoreSQLite {
	// 	repo = internal.NewSqliteNodeRepo(source_path)
	// }
	// if repo == nil {
	// 	logger.Fatal("Error: Unknown Store Type")
	// }
	//
	// srv := internal.NewNoteService(logger, repo)
	//
	// // CLI commander
	//
	// // Add command
	// cmdAdd := flag.NewFlagSet("add", flag.ExitOnError)
	// addContent := cmdAdd.String("c", "", "Inline note content")
	// addTitle := cmdAdd.String("title", "", "Note title")
	//
	// // Ls command
	// cmdLs := flag.NewFlagSet("ls", flag.ExitOnError)
	// lsTags := cmdLs.Bool("t", false, "List all tags")
	// lsPublic := cmdLs.Bool("public", false, "List public notes")
	// lsMoveTo := cmdLs.String("cp", "", "Copy notes to")
	//
	// // Get Command
	// cmdGet := flag.NewFlagSet("get", flag.ExitOnError)
	//
	// // Find Command
	// cmdFind := flag.NewFlagSet("find", flag.ExitOnError)
	// findTitle := cmdFind.String("n", "", "Note title")
	// findTags := cmdFind.String("t", "", "Note tags")
	// hasToExport := cmdFind.Bool("exp", false, "Export result")
	// moveTo := cmdFind.String("cp", "", "Copy notes to")
	//
	// // Share Command
	// cmdShare := flag.NewFlagSet("share", flag.ExitOnError)
	// findTitle = cmdShare.String("n", "", "Note title")
	// shareLimit := cmdShare.Int("limit", 1, "Limit of views allowed")
	// shareLink := cmdShare.String("link", "", "Link of note")
	//
	// cmdBM := flag.NewFlagSet("bookmark", flag.ExitOnError)
	cmdDocs := flag.NewFlagSet("docs", flag.ExitOnError)

	if len(os.Args) < 2 {
		logger.Fatalln("Expected one subcommand")
	}

	switch os.Args[1] {
	// case "add":
	// 	cmdAdd.Parse(os.Args[2:])
	// 	var str string
	// 	if *addContent == "" {
	// 		input, err := ReadInputFromEditor()
	// 		if err != nil {
	// 			logger.Fatalln(err)
	// 		}
	// 		str = string(input)
	// 	} else {
	// 		str = *addContent
	// 	}
	//
	// 	note, err := srv.New(*addTitle, str)
	// 	if err != nil {
	// 		logger.Fatalln(err)
	// 	}
	// 	fmt.Printf("%s", note.Title)
	// case "get":
	// 	cmdGet.Parse(os.Args[2:])
	// 	var title string
	// 	if len(cmdGet.Args()) > 0 {
	// 		title = cmdGet.Arg(0)
	// 	}
	// 	note, err := srv.GetByTitle(title)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Println(note.Content)
	// case "find":
	// 	cmdFind.Parse(os.Args[2:])
	// 	filters := []internal.Filter{}
	// 	filters = append(filters, internal.Filter{Field: "title", Value: *findTitle})
	// 	filters = append(filters, internal.Filter{Field: "tags", Value: *findTags})
	//
	// 	notes, err := srv.Find(filters)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	if *hasToExport {
	// 		//ExportNotes(notes)
	// 		ExportToDB(notes)
	// 	}
	// 	if *moveTo != "" {
	// 		MoveNotes(notes, *moveTo)
	// 	}
	case "ls":
		// cmdLs.Parse(os.Args[2:])
		// repo := internal.NewNoteRepository("")
		// notes, err := repo.All()
		// if err != nil {
		// 	log.Fatalf("Erro reading docs: %v", err)
		// }
		// for i, d := range notes {
		// 	log.Printf("%d: %s %s %s", i, d.Title, d.IsPublic, d.Type)
		// }
		// return
		// if *lsTags {
		// 	tags, err := srv.ListAllTags()
		// 	if err != nil {
		// 		logger.Fatalln(err)
		// 	}
		// 	for k, v := range tags {
		// 		fmt.Printf("- %s:%d\n", k, v)
		// 	}
		// } else if *lsPublic {
		// 	notes, err := srv.GetPublicNotes()
		// 	if err != nil {
		// 		logger.Fatalln(err)
		// 	}
		// 	if *lsMoveTo != "" {
		// 		MoveNotes(notes, *lsMoveTo)
		// 		return
		// 	}
		// 	for i, n := range notes {
		// 		fmt.Printf("%d %v     %s\n", i+1, n.Meta.IsPublic, n.Title)
		// 	}
		//
		// } else {
		// 	notes, err := srv.ListAll()
		// 	if err != nil {
		// 		logger.Fatalln(err)
		// 	}
		// 	for i, n := range notes {
		// 		fmt.Printf("%d %v     %s\n", i+1, n.Meta.IsPublic, n.Title)
		// 	}
		// }
	// case "share":
	// 	share_url := os.Getenv(PKMS_SHARE_URL)
	// 	cmdShare.Parse(os.Args[2:])
	// 	if *findTitle != "" {
	// 		n, err := srv.GetByTitle(*findTitle)
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}
	// 		key := internal.RandomKey(16)
	// 		encrypted, err := internal.Encrypt([]byte(n.Html), []byte(key))
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}
	// 		res, err := ShareNote(share_url, encrypted, *shareLimit)
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}
	// 		fmt.Printf("%s/%s#%s", share_url, res, key)
	// 		return
	// 	}
	// 	if *shareLink != "" {
	// 		GetSharedNote(share_url, *shareLink)
	// 		return
	// 	}
	// case "bookmark":
	// 	BookmarkCommand(cmdBM, os.Args, srv)
	// 	return
	case "docs":
		DocsCommand(cmdDocs, os.Args)
		return

	default:
		logger.Fatalln("Unknow subcommand")
	}
}

func ReadInputFromEditor() ([]byte, error) {
	file, err := os.CreateTemp(os.TempDir(), "swenotes")
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

func ShareNote(base_url, content string, limit int) (string, error) {
	url := fmt.Sprint(base_url, "/share")
	var request struct {
		Content string `json:"Content"`
		Limit   int    `json:"Limit"`
	}
	request.Content = content
	request.Limit = limit

	body, err := json.Marshal(&request)
	if err != nil {
		return "", err
	}

	reader := bytes.NewReader(body)

	res, err := http.Post(url, "application/json", reader)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", nil
	}

	var response struct {
		Result string `json:"result"`
	}

	err = json.Unmarshal(resBody, &response)
	if err != nil {
		return "", nil
	}

	return response.Result, nil
}

func GetSharedNote(base_url, id string) {
	url := fmt.Sprintf("%s/notes/%s", base_url, id)
	fmt.Println(url)

	var response struct {
		Result string `json:"Result"`
	}
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode == http.StatusNotFound {
		log.Printf("Note with id %s not found", id)
		return
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response.Result)
}
