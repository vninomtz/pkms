package reading

type Book struct {
	ID            string   `yaml:"id"`
	Type          string   `yaml:"type"`
	Status        string   `yaml:"status"`
	Title         string   `yaml:"title"`
	Author        string   `yaml:"author"`
	Language      string   `yaml:"language"`
	Categories    []string `yaml:"categories"`
	CompletedDate string   `yaml:"completed_date"`
	URL           string   `yaml:"url,omitempty"`
	Notes         string   `yaml:"notes"`
	Priority      string   `yaml:"priority,omitempty"`
}

type Metadata struct {
	Title          string `yaml:"title"`
	Description    string `yaml:"description"`
	Created        string `yaml:"created"`
	Updated        string `yaml:"updated"`
	TotalResources int    `yaml:"total_resources"`
	Breakdown      struct {
		Books struct {
			Completed int `yaml:"completed"`
			Pending   int `yaml:"pending"`
		} `yaml:"books"`
		Blogs struct {
			Completed int `yaml:"completed"`
			Pending   int `yaml:"pending"`
		} `yaml:"blogs"`
	} `yaml:"breakdown"`
}

type ReadingResources struct {
	Metadata  Metadata `yaml:"metadata"`
	Resources []Book   `yaml:"resources"`
}
