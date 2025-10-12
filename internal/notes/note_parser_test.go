package notes

import (
	"slices"
	"testing"
	"time"
)

func TestExtractMetadata(t *testing.T) {
	content := `---
title: Un arma contra los días solitarios
created: 2025-06-29
updated: 2025-08-10
public: true
type: writing
tags:
  - software-engineering
---
Extra content`
	exp_cre, _ := time.Parse("2006-01-02", "2025-06-29")
	exp_upd, _ := time.Parse("2006-01-02", "2025-08-10")
	exp := Note{
		Title:   "Un arma contra los días solitarios",
		Public:  true,
		Type:    "writing",
		Tags:    []string{"software-engineering"},
		Content: "Extra content",
		Created: exp_cre,
		Updated: exp_upd,
	}

	note, err := Parse([]byte(content))
	if err != nil {
		t.Errorf("No error expected: %s\n", err)
	}
	if note.Title != exp.Title {
		t.Errorf("Expected %s, got %s instead", exp.Title, note.Title)
	}
	if note.Public != exp.Public {
		t.Errorf("Expected %v, got %v instead", exp.Public, note.Public)
	}
	if note.Type != exp.Type {
		t.Errorf("Expected %s, got %s instead", exp.Type, note.Type)
	}
	if !slices.Equal(note.Tags, exp.Tags) {
		t.Errorf("Expected %s, got %s instead", exp.Tags, note.Tags)
	}
	if !exp.Created.Equal(note.Created) {
		t.Errorf("Expected %s, got %s instead", exp.Created, note.Created)
	}
	if !exp.Updated.Equal(note.Updated) {
		t.Errorf("Expected %s, got %s instead", exp.Created, note.Created)
	}
	if exp.Content != note.Content {
		t.Errorf("Expected '%s', got '%s' instead", exp.Content, note.Content)
	}
}
