package bookmarks

type Bookmark struct {
	ID          string   `yaml:"id"`
	URL         string   `yaml:"url"`
	Title       string   `yaml:"title"`
	Description string   `yaml:"description"`
	Image       string   `yaml:"image,omitempty"`
	Tags        []string `yaml:"tags,omitempty"`
	SavedDate   string   `yaml:"saved_date"`
	LastScraped string   `yaml:"last_scraped"`
	Status      string   `yaml:"status"`        // "ok" | "failed"
	FetchError  string   `yaml:"fetch_error,omitempty"`
}

type BookmarkMetadata struct {
	Title         string `yaml:"title"`
	Updated       string `yaml:"updated"`
	Total         int    `yaml:"total"`
	CacheTTLHours int    `yaml:"cache_ttl_hours"`
}

type BookmarkList struct {
	Metadata  BookmarkMetadata `yaml:"metadata"`
	Bookmarks []Bookmark       `yaml:"bookmarks"`
}

type BookmarkStats struct {
	Discovered int
	Fetched    int
	Cached     int
	Updated    int
	Failed     int
}
