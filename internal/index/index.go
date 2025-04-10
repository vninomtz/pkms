package index

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Location struct {
	File  string
	Token Token
}
type Indexer struct {
	root  string
	index map[string][]Location
}

func NewIndexer(root string) *Indexer {
	return &Indexer{
		index: make(map[string][]Location),
		root:  root,
	}
}

func (i *Indexer) Index() {
	files, err := filepath.Glob(filepath.Join(i.root, "*.md"))

	fmt.Printf("Indexing %d files of dir: %s\n\n", i.root, len(files))
	if err != nil {
		fmt.Println("Error scanning dir: "+i.root, err)
		return
	}
	fmt.Printf("%d files indexed\n", len(files))
	for _, file := range files {
		i.indexFile(file)
	}

	fmt.Printf("Final index %d\n", len(i.index))
}
func (i *Indexer) Search(query string) {
	matches, found := i.index[strings.ToLower(query)]

	fmt.Printf("%d results found for: %s\n\n", len(matches), query)

	if !found {
		fmt.Println("Not matches found")
		return
	}
	for _, match := range matches {
		fmt.Printf("File: %s\n", match.File)
		fmt.Printf("   -> %d %s\n", match.Token.Start, match.Token.Value)
		fmt.Println()
	}
}

func (i *Indexer) indexFile(filepath string) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", filepath, err)
		return
	}
	tokens := Tokenize(string(content))

	for _, token := range tokens {
		word := strings.ToLower(token.Value)
		if _, exist := i.index[word]; !exist {
			i.index[word] = []Location{}
		}
		i.index[word] = append(i.index[word], Location{File: filepath, Token: token})
	}
}

type Token struct {
	Value string
	Start int
	End   int
}

func Tokenize(text string) []Token {
	tokens := []Token{}

	re := regexp.MustCompile(`[a-zA-Z0-9áéíóúÁÉÍÓÚñÑ'-]+`)
	indexes := re.FindAllStringIndex(strings.ToLower(text), -1)
	for _, index := range indexes {
		s, e := index[0], index[1]
		word := text[s:e]
		tokens = append(tokens, Token{Value: word, Start: s, End: e})
	}
	return tokens
}
