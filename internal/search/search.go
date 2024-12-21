package search

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type searcher struct {
	index   map[string][]int
	indexed []Doc
	dir     string
}

type Doc struct {
	file  string
	title string
}

func NewSercher(dir string) *searcher {
	return &searcher{
		index:   make(map[string][]int),
		indexed: []Doc{},
		dir:     dir,
	}
}

func (s *searcher) Index() error {
	df, err := os.Open(s.dir)
	if err != nil {
		return err
	}
	fis, err := df.Readdir(-1)
	if err != nil {
		return err
	}
	if len(fis) == 0 {
		return errors.New(fmt.Sprintf("No files in %s", s.dir))
	}
	indexed := 0

	for _, fi := range fis {
		if !fi.IsDir() {
			if s.indexFile(s.dir + "/" + fi.Name()) {
				indexed++
			}
		}
	}
	return nil
}

func (s *searcher) indexFile(fn string) bool {
	if filepath.Ext(fn) != ".md" {
		//fmt.Printf("Just index .md files %s\n", fn)
		return false
	}
	f, err := os.Open(fn)
	if err != nil {
		fmt.Println(err)
		return false
	}
	// register new file
	x := len(s.indexed)
	s.indexed = append(s.indexed, Doc{fn, fn})
	//pdoc := &s.indexed[x]

	// scan lines
	r := bufio.NewReader(f)
	//lines := 0
	for {
		b, isPrefix, err := r.ReadLine()
		switch {
		case err == io.EOF:
			return true
		case err != nil:
			fmt.Println(err)
			return true
		case isPrefix:
			fmt.Printf("%s, unexpected long line\n", fn)
			return true
		}
		// index line of text in b
		// again, in a real program you would write a much
		// nicer word splitter
	wordLoop:
		for _, bword := range bytes.Fields(b) {
			bword := bytes.Trim(bword, ".,-~?!\"'`;:()<>[]{}\\|/=_+*&^%$#@")
			if len(bword) > 0 {
				word := string(bword)
				dl := s.index[word]
				for _, d := range dl {
					if d == x {
						continue wordLoop
					}
				}
				s.index[word] = append(dl, x)
			}
		}
	}
}

func (s *searcher) Search(str string) []string {
	files := []string{}
	switch dl := s.index[str]; len(dl) {
	case 0:
		fmt.Println("no match")
	case 1:
		doc := s.indexed[dl[0]]
		files = append(files, doc.title)
	default:
		for _, d := range dl {
			doc := s.indexed[d]
			fmt.Println("    ", doc.title)
			files = append(files, doc.title)
		}
	}
	return files
}
