package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

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

	port := flag.String("port", "8000", "Port for http server")
	host := flag.String("host", "", "Server host")
	fs := flag.Bool("fs", false, "Use file system repo")

	flag.Parse()

	repo := internal.NewSqliteNodeRepo(DB_PATH)

	if *fs {
		repo = internal.NewFsRepo(DIR_PATH)
	}

	http.HandleFunc("/api/notes", func(w http.ResponseWriter, r *http.Request) {
		srv := internal.NewNoteService(repo)

		notes, err := srv.ListAll()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

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
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
