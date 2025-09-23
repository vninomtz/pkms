package cmd

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/vninomtz/pkms/internal"
	"github.com/vninomtz/pkms/internal/loader"
)

const PKMS_SHARE_URL = "PKMS_SHARE_URL"

func ShareCommand(args []string) {
	share_url := os.Getenv(PKMS_SHARE_URL)
	fs := flag.NewFlagSet("share", flag.ExitOnError)
	path := fs.String("path", "", "Documents from a path")
	filename := fs.String("name", "", "Name of the note")

	fs.Parse(os.Args[2:])

	if *filename == "" {
		log.Fatal("Missing filename")
	}

	loader := loader.New(*path)
	err := loader.Load()
	if err != nil {
		log.Fatal(err)
	}
	doc := loader.FindByName(*filename)
	if doc == nil {
		fmt.Printf("Document with name %s not found\n", *filename)
		return
	}

	key := internal.RandomKey(16)
	encrypted, err := internal.Encrypt(doc.Content, []byte(key))
	if err != nil {
		log.Fatal(err)
	}
	res, err := ShareNote(share_url, encrypted, 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s/%s#%s", share_url, res, key)
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
