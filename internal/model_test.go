package internal

import (
	"testing"
)

const (
	inputFile       = "../testdata/markdownFullSyntax.md"
	metadataSection = `---
title: "Hello Markdown"
author: "Awesome Me"
date: "2018-02-14"
output: html_document
tags: tag1, tag2, tag3
---`

	LinksExamples = `
- [[pkms]]
- [[project-bycuriosity]]
- [Example test](https://example.com/test/test-long/)
- [Example test](teste no link)

https://example2.com/test2/test2-long/
	`
)

func TestGetLinks(t *testing.T) {
	expLinks := 2
	node := Node{
		Bytes: []byte(LinksExamples),
	}

	links, err := node.Links()
	if err != nil {
		t.Errorf("Unexpected error %v\n", err)
	}

	if expLinks != len(links) {
		t.Errorf("Expected %d, get %d instead\n", expLinks, len(links))
	}

	for i, l := range links {
		t.Logf("%d: %s\n", i, l)
	}

}

func TestIncludeTags(t *testing.T) {
	meta := Metadata{
		Title: "",
		Tags:  []string{"tag1", "tag2"},
	}

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
