package bookmarks

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type Bookmark struct {
	Id          int    `json:"id"`
	Url         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

func GetBookmarkFromUrl(url string) (Bookmark, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Bookmark{}, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return Bookmark{}, err
	}
	log.Printf("StatusCode: %s\n", resp.Status)
	log.Printf("Header Content-Type: %s\n", resp.Header.Get("Content-Type"))
	log.Printf("Host: %s", resp.Request.URL.Host)

	if resp.StatusCode > 300 {
		return Bookmark{}, errors.New(fmt.Sprintf("Error to fetch url: %s", resp.Status))
	}

	defer resp.Body.Close()
	bk := ExtractMetadata(resp.Body)
	bk.Url = url

	return bk, nil
}

func ExtractMetadata(r io.Reader) Bookmark {
	bk := Bookmark{}
	z := html.NewTokenizer(r)

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return bk
		case html.StartTagToken, html.SelfClosingTagToken:
			t := z.Token()
			if t.Data == "body" {
				return bk
			}
			if t.Data == "title" {
				if z.Next() == html.TextToken {
					bk.Title = strings.TrimSpace(z.Token().Data)
				}
			}
			if t.Data == "meta" {
				var name, property, content string
				for _, attr := range t.Attr {
					switch strings.ToLower(attr.Key) {
					case "name":
						name = strings.ToLower(attr.Val)
					case "property":
						property = strings.ToLower(attr.Val)
					case "content":
						content = attr.Val
					}
				}
				if name == "description" {
					bk.Description = content
				}
				if property == "og:title" && bk.Title == "" {
					bk.Title = content
				}
				if property == "og:description" && bk.Description == "" {
					bk.Description = content
				}
				if property == "og:url" {
					bk.Url = content
				}
				if property == "og:image" {
					u, err := url.Parse(content)
					if err == nil && u.Scheme != "" && u.Host != "" {
						bk.Image = content
					}
				}

			}
		}
	}
}
