package internal

import (
	"errors"
	"log"
	"strings"
)

type noteService struct {
	repo   NodeRepository
	logger *log.Logger
}

func NewNoteService(logger *log.Logger, repo NodeRepository) *noteService {
	return &noteService{
		repo:   repo,
		logger: logger,
	}
}

func (s *noteService) New(title, content string) (Node, error) {
	note, err := NewNote(title, content)
	if err != nil {
		return Node{}, err
	}
	return note, s.repo.Save(note)
}

func (s *noteService) ListAll() ([]Node, error) {
	notes, err := s.repo.GetNodes()
	if err != nil {
		return nil, errors.New("Error to consult notes")
	}
	return notes, nil
}

func (s *noteService) ListAllTags() (map[string]int, error) {
	tags := make(map[string]int)
	notes, err := s.ListAll()
	if err != nil {
		return nil, err
	}
	for _, note := range notes {
		meta := note.Meta
		for _, t := range meta.Tags {
			_, ok := tags[t]
			if ok {
				tags[t] = tags[t] + 1
			} else {
				tags[t] = 1
			}
		}
	}
	return tags, nil
}

func (s *noteService) GetBookmarks() ([]string, error) {
	dic := make(map[string]int)
	result := []string{}
	notes, err := s.ListAll()
	if err != nil {
		return nil, err
	}
	for _, note := range notes {
		links, err := note.Links()
		if err != nil {
			log.Printf("Error to extract links of Note %s\n", note.Title)
			continue
		}
		for _, link := range links {
			_, ok := dic[link]
			if !ok {
				dic[link] = 0
				result = append(result, link)
			}
			dic[link]++

		}
	}
	return result, nil
}

func (s *noteService) GetByTitle(title string) (Node, error) {
	notes, err := s.ListAll()
	if err != nil {
		return Node{}, err
	}
	found := -1
	for i := 0; i < len(notes); i++ {
		if notes[i].Title == title {
			found = i
			break
		}
	}
	if found == -1 {
		return Node{}, errors.New("Note note found")
	}

	n := notes[found]
	html, err := MDToHTML([]byte(n.Content))
	if err != nil {
		log.Printf("Error parsing MD to Html for note %s", n.Title)
	} else {
		n.Html = html
	}

	return n, nil
}

func filtersToMap(filters []Filter) map[string]string {
	fields := map[string]bool{"title": true, "tags": true}
	mapFilters := map[string]string{}
	for _, v := range filters {
		if v.Field != "" && fields[strings.ToLower(v.Field)] {
			if v.Value != "" {
				mapFilters[strings.ToLower(v.Field)] = v.Value
			}
		}
	}
	return mapFilters
}

func (s *noteService) Find(_filters []Filter) ([]Node, error) {
	notes, err := s.ListAll()
	if err != nil {
		return nil, err
	}
	filters := filtersToMap(_filters)

	var founds []Node
	for _, note := range notes {
		if IncludeNote(filters, note, s.logger) {
			html, err := MDToHTML([]byte(note.Content))
			if err != nil {
				log.Printf("Error parsing MD to Html for note %s", note.Title)
			} else {
				note.Html = html
			}
			founds = append(founds, note)
		}
	}
	return founds, nil
}

func (s *noteService) GetPublicNotes() ([]Node, error) {
	notes, err := s.ListAll()
	if err != nil {
		return nil, err
	}
	var founds []Node
	for _, note := range notes {
		if note.Meta.IsPublic {
			founds = append(founds, note)
		}
	}
	return founds, nil
}

func IncludeNote(filters map[string]string, note Node, logger *log.Logger) bool {
	val, ok := filters["tags"]
	if ok {
		meta := note.Meta
		if !meta.IncludeTags(val) {
			return false
		}
	}
	val, ok = filters["title"]
	if ok {
		source := strings.ToLower(note.Title)
		target := strings.ToLower(val)
		if !strings.Contains(source, target) {
			return false
		}
	}
	return true
}
