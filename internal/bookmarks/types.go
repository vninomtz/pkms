package bookmarks

type Bookmark struct {
	ID          string   `yaml:"id"`
	URL         string   `yaml:"url"`
	Title       string   `yaml:"title"`
	Description string   `yaml:"description"`
	Image       string   `yaml:"image,omitempty"`
	Tags        []string `yaml:"tags,omitempty"`
	SavedDate   string   `yaml:"saved_date"`    // When discovered in notes
	LastScraped string   `yaml:"last_scraped"`  // When metadata was fetched (RFC3339)
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
