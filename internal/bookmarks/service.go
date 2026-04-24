package bookmarks

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

type Service struct {
	filePath  string
	cacheTTL  time.Duration
}

func New(filePath string) *Service {
	return &Service{
		filePath: filePath,
		cacheTTL: 24 * time.Hour,
	}
}

func (s *Service) LoadBookmarks() (*BookmarkList, error) {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return &BookmarkList{
				Metadata: BookmarkMetadata{
					Title:         "Bookmarks",
					Updated:       time.Now().Format("2006-01-02"),
					Total:         0,
					CacheTTLHours: 24,
				},
				Bookmarks: []Bookmark{},
			}, nil
		}
		return nil, fmt.Errorf("error reading bookmarks: %w", err)
	}

	var bookmarks BookmarkList
	if err := yaml.Unmarshal(data, &bookmarks); err != nil {
		return nil, fmt.Errorf("error parsing YAML: %w", err)
	}

	// Migrate legacy bookmarks that predate the status field
	for i := range bookmarks.Bookmarks {
		if bookmarks.Bookmarks[i].Status == "" {
			if bookmarks.Bookmarks[i].LastScraped != "" {
				bookmarks.Bookmarks[i].Status = "ok"
			}
		}
	}

	return &bookmarks, nil
}

func (s *Service) SaveBookmarks(bookmarks *BookmarkList) error {
	bookmarks.Metadata.Updated = time.Now().Format("2006-01-02")
	bookmarks.Metadata.Total = len(bookmarks.Bookmarks)

	data, err := yaml.Marshal(bookmarks)
	if err != nil {
		return fmt.Errorf("error marshaling YAML: %w", err)
	}

	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		return fmt.Errorf("error writing bookmarks: %w", err)
	}

	return nil
}

func (s *Service) GetByURL(url string) *Bookmark {
	bookmarks, err := s.LoadBookmarks()
	if err != nil {
		return nil
	}

	for _, b := range bookmarks.Bookmarks {
		if strings.EqualFold(b.URL, url) {
			return &b
		}
	}

	return nil
}

func (s *Service) IsCached(bookmark *Bookmark) bool {
	if bookmark.LastScraped == "" {
		return false
	}

	t, err := time.Parse(time.RFC3339, bookmark.LastScraped)
	if err != nil {
		return false
	}

	return time.Since(t) < s.cacheTTL
}

func (s *Service) UpdateMetadata(bookmark *Bookmark, metadata map[string]string) {
	if title, ok := metadata["title"]; ok && title != "" {
		bookmark.Title = sanitizeString(title)
	}

	if desc, ok := metadata["description"]; ok && desc != "" {
		bookmark.Description = sanitizeString(desc)
	}

	if ogDesc, ok := metadata["og:description"]; ok && ogDesc != "" && bookmark.Description == "" {
		bookmark.Description = sanitizeString(ogDesc)
	}

	if image, ok := metadata["og:image"]; ok && image != "" {
		bookmark.Image = sanitizeString(image)
	}

	bookmark.LastScraped = time.Now().Format(time.RFC3339)
	bookmark.Status = "ok"
	bookmark.FetchError = ""
}

func (s *Service) MarkFailed(bookmark *Bookmark, errMsg string) {
	bookmark.LastScraped = time.Now().Format(time.RFC3339)
	bookmark.Status = "failed"
	bookmark.FetchError = sanitizeString(errMsg)
}

func (s *Service) AddFailedBookmark(rawURL string, errMsg string) *Bookmark {
	return &Bookmark{
		ID:          generateID(),
		URL:         rawURL,
		SavedDate:   time.Now().Format("2006-01-02"),
		LastScraped: time.Now().Format(time.RFC3339),
		Status:      "failed",
		FetchError:  sanitizeString(errMsg),
	}
}

// sanitizeString removes problematic characters that break YAML parsing
func sanitizeString(s string) string {
	// Remove newlines and extra whitespace
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")
	s = strings.ReplaceAll(s, "\t", " ")
	// Collapse multiple spaces
	s = strings.Join(strings.Fields(s), " ")
	// Limit length to avoid huge YAML
	if len(s) > 500 {
		s = s[:500]
	}
	return s
}

func (s *Service) AddBookmark(rawURL string, metadata map[string]string) *Bookmark {
	bookmark := &Bookmark{
		ID:        generateID(),
		URL:       rawURL,
		SavedDate: time.Now().Format("2006-01-02"),
	}

	s.UpdateMetadata(bookmark, metadata)

	return bookmark
}

func generateID() string {
	return uuid.New().String()[:8]
}
