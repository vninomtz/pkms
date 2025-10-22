package notes

import (
	"bufio"
	"bytes"
	"net/url"
	"regexp"

	"github.com/adrg/frontmatter"
)

func Parse(content []byte) (Note, error) {
	var note Note
	res, err := frontmatter.Parse(bytes.NewReader(content), &note)
	if err != nil {
		return note, err
	}
	note.Content = string(res)
	note.Links, err = parse_links(content)
	if err != nil {
		return note, err
	}
	note.Notes, err = parse_wikilinks(res)
	if err != nil {
		return note, err
	}
	return note, nil
}

func parse_links(content []byte) ([]string, error) {
	links := []string{}
	scanner := bufio.NewScanner(bytes.NewBuffer(content))
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		word := scanner.Text()
		if IsUrl(word) {
			links = append(links, word)
		} else {
			re := regexp.MustCompile(`\((.*?)\)`)
			matches := re.FindAllStringSubmatch(word, -1)
			if len(matches) > 0 {
				for _, m := range matches {
					is := IsUrl(m[1])
					if is {
						links = append(links, m[1])
					}
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return links, err
	}
	return links, nil
}
func parse_wikilinks(content []byte) ([]string, error) {
	links := []string{}
	scanner := bufio.NewScanner(bytes.NewBuffer(content))
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		word := scanner.Text()
		re := regexp.MustCompile(`\[\[([^\]]+)\]\]`)
		matches := re.FindAllStringSubmatch(word, -1)
		if len(matches) > 0 {
			for _, m := range matches {
				links = append(links, m[1])
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return links, err
	}
	return links, nil

}

func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
