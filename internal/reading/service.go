package reading

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

type Service struct {
	filePath string
}

func New(filePath string) *Service {
	return &Service{filePath: filePath}
}

func (s *Service) LoadResources() (*ReadingResources, error) {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	var resources ReadingResources
	if err := yaml.Unmarshal(data, &resources); err != nil {
		return nil, fmt.Errorf("error parsing YAML: %w", err)
	}

	return &resources, nil
}

func (s *Service) SaveResources(resources *ReadingResources) error {
	data, err := yaml.Marshal(resources)
	if err != nil {
		return fmt.Errorf("error marshaling YAML: %w", err)
	}

	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	return nil
}

func (s *Service) FindBook(title, author string) *Book {
	resources, err := s.LoadResources()
	if err != nil {
		return nil
	}

	for _, book := range resources.Resources {
		if book.Type == "book" &&
			strings.EqualFold(book.Title, title) &&
			strings.EqualFold(book.Author, author) {
			return &book
		}
	}

	return nil
}

func (s *Service) AddBook(book *Book) error {
	resources, err := s.LoadResources()
	if err != nil {
		return err
	}

	book.ID = generateUUID()
	book.Type = "book"

	resources.Resources = append(resources.Resources, *book)

	if book.Status == "completed" {
		resources.Metadata.Breakdown.Books.Completed++
	} else {
		resources.Metadata.Breakdown.Books.Pending++
	}

	resources.Metadata.TotalResources++
	resources.Metadata.Updated = time.Now().Format("2006-01-02")

	return s.SaveResources(resources)
}

func generateUUID() string {
	return uuid.New().String()[:8]
}

func PromptString(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func PromptStringWithDefault(prompt, defaultValue string) string {
	fullPrompt := fmt.Sprintf("%s [%s]: ", prompt, defaultValue)
	input := PromptString(fullPrompt)
	if input == "" {
		return defaultValue
	}
	return input
}

func PromptBool(prompt string) bool {
	for {
		response := PromptString(prompt + " (s/n): ")
		if strings.ToLower(response) == "s" {
			return true
		} else if strings.ToLower(response) == "n" {
			return false
		}
		fmt.Println("Por favor, responde 's' o 'n'")
	}
}

func PromptCategories(prompt string) []string {
	input := PromptString(prompt + " (separadas por comas, opcional): ")
	if input == "" {
		return []string{}
	}
	parts := strings.Split(input, ",")
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}
	return parts
}
