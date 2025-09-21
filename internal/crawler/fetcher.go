package crawler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

func Fetch(url string) (*Page, error) {
	client := &http.Client{
		Timeout: 60 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body: %w", err)
	}

	headers := make(map[string]string)
	for k, v := range resp.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}
	page := &Page{
		URL:        resp.Request.URL.String(),
		StatusCode: resp.StatusCode,
		Headers:    headers,
		FetchedAt:  time.Now(),
		HTML:       body,
	}

	return page, nil
}
func FetchMultiple(urls []string) ([]*Page, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	results := []*Page{}

	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			page, err := Fetch(u)
			if err != nil {
				log.Printf("Error to inspect %s: %w\n", url, err)
				return
			}

			mu.Lock()
			results = append(results, page)
			mu.Unlock()
		}(url)
	}

	wg.Wait()
	return results, nil
}
