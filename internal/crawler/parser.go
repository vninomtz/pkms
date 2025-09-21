package crawler

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

var f func(*html.Node)

func ParseHtml(htmlBytes []byte) (map[string]string, error) {
	result := make(map[string]string)

	doc, err := html.Parse(bytes.NewReader(htmlBytes))
	if err != nil {
		return nil, fmt.Errorf("error parseando HTML: %w", err)
	}

	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			// <title>
			if n.Data == "title" && n.FirstChild != nil {
				result["title"] = n.FirstChild.Data
			}

			// <meta ...>
			if n.Data == "meta" {
				var name, content, property string
				for _, attr := range n.Attr {
					switch strings.ToLower(attr.Key) {
					case "name":
						name = attr.Val
					case "content":
						content = attr.Val
					case "property":
						property = attr.Val
					}
				}

				if name != "" && content != "" {
					result[strings.ToLower(name)] = content
				}

				if property != "" && content != "" {
					result[strings.ToLower(property)] = content
				}
			}
		}

		// Recursivo por hijos
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)
	return result, nil
}
