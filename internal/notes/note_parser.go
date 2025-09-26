package notes

import (
	"bufio"
	"bytes"
	"net/url"
	"regexp"

	"github.com/adrg/frontmatter"
)

func ParseMarkdown(content []byte) (Note, error) {
	note := Note{}
	var data map[string]interface{}
	content, err := frontmatter.Parse(bytes.NewReader(content), &data)
	if err != nil {
		return note, err
	}
	note.Content = string(content)

	isPublic, ok := data["isPublic"]
	if ok {
		if b, ok := isPublic.(bool); ok {
			note.IsPublic = b
		} else {
			note.IsPublic = false
		}
	}
	noteType, ok := data["type"]
	if ok {
		note.Type = noteType.(string)
	}
	title, ok := data["title"]
	if ok {
		note.Title = title.(string)
	}
	note.Tags = parse_tags(data)
	note.Links, err = parse_links(content)
	if err != nil {
		return note, err
	}
	return note, nil
}
func parse_tags(data map[string]interface{}) []string {
	res := []string{}
	tags, ok := data["tags"]
	if !ok {
		return res
	}
	switch v := tags.(type) {
	case string:
		res = append(res, v)
	case []interface{}:
		for _, item := range v {
			str, ok := item.(string)
			if ok {
				res = append(res, str)
			}
		}
	}
	return res
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

func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
