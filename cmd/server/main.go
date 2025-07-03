package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
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
type SharedNote struct {
	Content string
}

var (
	shared = make(map[string]SharedNote)
	mu     sync.Mutex
)

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

	indexSearch := search.NewSercher(*directory)
	err = indexSearch.Index()
	if err != nil {
		log.Println(err)
		return err
	}

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
	http.HandleFunc("/api/search", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		text := query.Get("q")
		log.Printf("Search by: %s\n", text)
		notes := indexSearch.Search(text)

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
	http.HandleFunc("/api/bookmarks", func(w http.ResponseWriter, r *http.Request) {
		result, err := searcher.GetBookmarks()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		res := &struct {
			Body    interface{}
			Records int
		}{
			Body:    result,
			Records: len(result),
		}
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}

	})
	http.HandleFunc("/api/share", HandleShare)
	http.HandleFunc("/api/share/{uuid}", HandleShareDetail)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	log.Printf("Serving %s on HTTP port: %s\n", *directory, *port)
	return http.ListenAndServe(":"+*port, nil)
}

func HandleShare(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Content string `json:"content"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Content == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	id := uuid.New()
	mu.Lock()
	shared[id.String()] = SharedNote{Content: req.Content}
	mu.Unlock()

	json.NewEncoder(w).Encode(map[string]string{
		"result": id.String(),
	})
}

func HandleShareDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	id := r.PathValue("uuid")
	mu.Lock()
	note, ok := shared[id]
	mu.Unlock()

	if !ok {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"result": note.Content,
	})
}
