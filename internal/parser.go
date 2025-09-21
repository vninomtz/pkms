package internal

import (
	"bufio"
	"bytes"
	"html/template"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
)

type content struct {
	Title string
	Body  template.HTML
}
type parser struct {
	filename string
	tmpl     *template.Template
}

func NewTemplateParser(pathTmp, templateName string) *parser {
	lp := filepath.Join(pathTmp, "*.html")
	t, _ := template.ParseGlob(lp)
	return &parser{
		filename: templateName,
		tmpl:     t,
	}
}

func (p *parser) Parse(node Document) ([]byte, error) {
	html, err := p.MDToHTML(node.Content)
	if err != nil {
		return nil, err
	}
	c := content{
		Title: node.Name(),
		Body:  template.HTML(html),
	}
	var bufferHtml bytes.Buffer
	if err := p.tmpl.ExecuteTemplate(&bufferHtml, "layout", c); err != nil {
		return nil, err
	}

	return bufferHtml.Bytes(), nil
}

func (p *parser) MDToHTML(input []byte) (string, error) {
	var buf bytes.Buffer

	md := goldmark.New(
		goldmark.WithExtensions(highlighting.NewHighlighting(
			highlighting.WithStyle("dracula"),
		)),
	)
	if err := md.Convert(input, &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}
func MDToHTML(input []byte) (string, error) {
	var buf bytes.Buffer

	md := goldmark.New(
		goldmark.WithExtensions(highlighting.NewHighlighting(
			highlighting.WithStyle("dracula"),
		)),
	)
	if err := md.Convert(input, &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func ParseNodeToHTML(node FileNode, templateFile string) ([]byte, error) {
	var buf bytes.Buffer

	md := goldmark.New(
		goldmark.WithExtensions(highlighting.NewHighlighting(
			highlighting.WithStyle("dracula"),
		)),
	)
	if err := md.Convert(node.Content, &buf); err != nil {
		return nil, err
	}

	t, err := template.ParseFiles(templateFile)
	if err != nil {
		return nil, err
	}

	//Instantiate the content type, adding the title and body
	c := content{
		Title: node.Name(),
		Body:  template.HTML(buf.String()),
	}

	//Create a buffer of bytes to write to file
	var bufferHtml bytes.Buffer

	// Write html to bytes buffer
	if err := t.Execute(&bufferHtml, c); err != nil {
		return nil, err
	}

	return bufferHtml.Bytes(), nil
}

func ParseDocument(doc Document) (Note, error) {
	name := strings.TrimSuffix(doc.Filename, doc.Ext)
	content, meta, err := ExtractMetadata(doc.Content)
	if err != nil {
		return Note{}, err
	}
	links, err := ExtractLinks(doc.Content)
	if err != nil {
		return Note{}, err
	}
	n := Note{
		Title:    name,
		Content:  string(content),
		IsPublic: meta.IsPublic,
		Tags:     meta.Tags,
		Type:     meta.Type,
		Links:    links,
	}
	if meta.Title != "" {
		n.Title = meta.Title
	}
	return n, nil
}
func ExtractLinks(content []byte) ([]string, error) {
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
