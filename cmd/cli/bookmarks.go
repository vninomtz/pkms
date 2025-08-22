package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/vninomtz/pkms/internal"
	"github.com/vninomtz/pkms/internal/bookmarks"
)

func BookmarkCommand(cmd *flag.FlagSet, args []string, srv internal.NoteService) {
	url := cmd.String("url", "", "Link to get bookmark")
	ls := cmd.Bool("ls", false, "Only display values")
	cmd.Parse(args[2:])

	if *url != "" {
		meta, err := bookmarks.GetBookmarkFromUrl(*url)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Título:", meta.Title)
		fmt.Println("Descripción:", meta.Description)
		fmt.Println("Image:", meta.Image)

		return
	}

	links, err := srv.GetBookmarks()
	if *ls {
		log.Printf("Bookmarks: %d Total\n", len(links))
	} else {
		if err != nil {
			log.Fatal(err)
		}
		var wg sync.WaitGroup
		var mu sync.Mutex
		items := []bookmarks.Bookmark{}

		for i, l := range links {
			wg.Add(1)
			go FetchBookmark(l, i, &wg, &items, &mu)
		}

		wg.Wait()

		log.Printf("Bookmarks: %d Total, %d success\n", len(links), len(items))

		ExportBookmarks(items)

	}
}

func FetchBookmark(url string, id int, wg *sync.WaitGroup, results *[]bookmarks.Bookmark, mu *sync.Mutex) {
	defer wg.Done()

	bk, err := bookmarks.GetBookmarkFromUrl(url)
	if err != nil {
		log.Printf("Error to get bookmark from %s: %v", url, err)
		return
	}
	bk.Id = id

	mu.Lock()
	*results = append(*results, bk)
	mu.Unlock()

}

func ExportBookmarks(items []bookmarks.Bookmark) {
	res, err := json.Marshal(items)
	if err != nil {
		log.Fatal("Unexpected error parsing json: ", err)
	}
	err = os.WriteFile("bookmarks.json", res, 0644)
	if err != nil {
		log.Fatal("Unexpected error writing file: ", err)
	}
}
