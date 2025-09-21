package crawler

import "time"

type Page struct {
	URL        string            `json:"url"`
	StatusCode int               `json:"status_code"`
	Headers    map[string]string `json:"headers"`
	FetchedAt  time.Time         `json:"fetched_at"`
	HTML       []byte            `json:"-"`
}
