package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/vninomtz/pkms/internal/crawler"
)

func InspectCommand(args []string) {
	fs := flag.NewFlagSet("inspect", flag.ExitOnError)
	url := fs.String("url", "", "Link to get bookmark")
	fs.Parse(args)

	urls := strings.Split(*url, ",")

	if len(urls) == 0 {
		log.Fatal("Missing urls")
	}

	pages, err := crawler.FetchMultiple(urls)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Urls inspected: ", len(pages))

	for _, page := range pages {
		fmt.Println("URL:", page.URL)
		fmt.Println("Status:", page.StatusCode)
		fmt.Println("Length HTML:", len(page.HTML))
		fmt.Println()
	}
	ExportPages(pages)

}
func ExportPages(items []*crawler.Page) {
	res, err := json.Marshal(items)
	if err != nil {
		log.Fatal("Unexpected error parsing json: ", err)
	}
	err = os.WriteFile("pages.json", res, 0644)
	if err != nil {
		log.Fatal("Unexpected error writing file: ", err)
	}
}
