package internal

import "testing"

func TestExtractMetadata(t *testing.T) {
	content := `---
title: test
tags: tag1, tag2
---

  Content
  `
	expected := "tag1, tag2"
	meta, err := ExtractMetadata(content)
	if err != nil {
		t.Errorf("No error expected get %s instead", err)
	}
	if meta.Tags != expected {
		t.Errorf("Get %s expected %s instead", meta.Tags, expected)
	}
}

func TestIncludeTags(t *testing.T) {
	meta := Metadata{Title: "", Tags: "tag1, tag2"}

	if !meta.IncludeTags("tag1") {
		t.Errorf("Expected %v, get %v instead", true, meta.IncludeTags("tag1"))
	}
	if !meta.IncludeTags("tag1, tag2") {
		t.Errorf("Expected %v, get %v instead", true, meta.IncludeTags("tag1,tag2"))
	}
	if meta.IncludeTags("tag3") {
		t.Errorf("Expected %v, get %v instead", false, meta.IncludeTags("tag3"))
	}
	if meta.IncludeTags("tag1, tag3") {
		t.Errorf("Expected %v, get %v instead", false, meta.IncludeTags("tag1,tag3"))
	}
}
