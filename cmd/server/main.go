package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/vninomtz/pkms/internal"
	"github.com/vninomtz/pkms/internal/search"
)

const (
	DB_NOTES_PATH  = "DB_NOTES_PATH"
	DIR_NOTES_PATH = "DIR_NOTES_PATH"
)

type templateHandler struct {
	tmpl *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Host": r.Host,
	}
	t.tmpl.ExecuteTemplate(w, "layout", data)
}

func main() {
	if err := FileServerRun(); err != nil {
		log.Printf("Error in server %v\n", err)
	}
}

func FileServerRun() error {
	lp := filepath.Join("./templates/", "*.html")
	tmpl := template.Must(template.New("pkms").ParseGlob(lp))
	port := flag.String("p", "8100", "port to serve on")
	directory := flag.String("dir", "", "Directory to serve content")
	flag.Parse()

	collector := internal.NewCollector(*directory, "")
	nodes, err := collector.Collect()
	if err != nil {
		log.Println(err)
		return err
	}
	searcher := internal.NewSearcher(nodes)

	//http.Handle("/", http.FileServer(http.Dir(*directory)))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		content := struct {
			Title string
		}{
			Title: "vic.aware",
		}
		if err := tmpl.ExecuteTemplate(w, "home", content); err != nil {
			log.Println(err)
			w.Write([]byte("Unexpected error"))
		}
	})
	http.HandleFunc("/writings", func(w http.ResponseWriter, r *http.Request) {
		content := struct {
			Title string
			Items []map[string]string
		}{
			Title: "Writings | vic.aware",
			Items: collector.ToMaps(),
		}
		if err := tmpl.ExecuteTemplate(w, "writings", content); err != nil {
			log.Println(err)
			w.Write([]byte("Unexpected error"))
		}
	})
	http.HandleFunc("/writings/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		n, err := searcher.File(slug)
		if err != nil {
			log.Printf("Error %v throw by File %s\n", err, slug)
			w.Write([]byte("Not found"))
			return
		}
		content, err := internal.MDToHTML(n.Content)
		if err != nil {
			log.Printf("Error %v to parse %s\n", err, slug)
			w.Write([]byte("Unexpected error"))
			return
		}
		pContent := struct {
			Title string
			Body  template.HTML
		}{
			Title: n.Name(),
			Body:  template.HTML(content),
		}

		if err := tmpl.ExecuteTemplate(w, "layout", pContent); err != nil {
			w.Write([]byte("Unexpected error"))
		}
	})

	log.Printf("Serving %s on HTTP port: %s\n", *directory, *port)
	return http.ListenAndServe(":"+*port, nil)
}

func Server() {
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
