package internal

import (
	"fmt"
	"os"
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
)

func TestExtractMetadata(t *testing.T) {
	expected := "tag1, tag2, tag3"
	meta, err := ExtractMetadata(metadataSection)
	if err != nil {
		t.Errorf("No error expected get %s instead", err)
	}
	if meta.Tags != expected {
		t.Errorf("Expected %s, get instead %s", meta.Tags, expected)
	}
}

func TestGetYaml(t *testing.T) {
	input, err := os.ReadFile(inputFile)
	if err != nil {
		t.Fatal(err)
	}

	yaml := GetYaml(string(input))
	if yaml != metadataSection {
		fmt.Println(yaml)
		t.Errorf("Expected %s, get instead %s", metadataSection, yaml)
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
