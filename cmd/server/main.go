package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/vninomtz/pkms/internal/search"
)

const (
	DB_NOTES_PATH  = "DB_NOTES_PATH"
	DIR_NOTES_PATH = "DIR_NOTES_PATH"
)

func main() {
	os.Setenv(DB_NOTES_PATH, "/Users/vnino/github.com/vninomtz/pkms/database/nodes.db")
	os.Setenv(DIR_NOTES_PATH, "/Users/vnino/Library/Mobile Documents/iCloud~md~obsidian/Documents/vnotes/docs")
	//DB_PATH := os.Getenv(DB_NOTES_PATH)
	DIR_PATH := os.Getenv(DIR_NOTES_PATH)

	port := flag.String("port", "8000", "Port for http server")
	host := flag.String("host", "", "Server host")
	//fs := flag.Bool("fs", false, "Use file system repo")
	//logger := log.New(os.Stdout, "INFO: ", log.Ltime)

	flag.Parse()

	//repo := internal.NewSqliteNodeRepo(DB_PATH)
	//repo := internal.NewFsRepo(logger, DIR_PATH)
	searcher := search.NewSercher(DIR_PATH)

	err := searcher.Index()
	if err != nil {
		log.Fatal(err)
	}
	//srv := internal.NewNoteService(logger, repo)

	http.HandleFunc("/api/notes", func(w http.ResponseWriter, r *http.Request) {

		query := r.URL.Query()

		text := query.Get("q")

		notes := searcher.Search(text)
		//notes, err := srv.ListAll()
		/*if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}*/

		res := &struct {
			Body    interface{}
			Records int
		}{
			Body:    notes,
			Records: len(notes),
		}
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}

	})

	addr := fmt.Sprintf("%s:%s", *host, *port)

	fmt.Printf("Server runing at: http://localhost:%s", *port)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
