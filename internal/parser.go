package internal

import (
	"bytes"
	"html/template"
	"path/filepath"
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

func (p *parser) Parse(node FileNode) ([]byte, error) {
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
	name := strings.TrimSuffix(doc.Name, doc.Ext)
	content, meta, err := ExtractMetadata(doc.Content)
	if err != nil {
		return Note{}, err
	}
	n := Note{
		Title:    name,
		Content:  string(content),
		IsPublic: meta.IsPublic,
		Tags:     meta.Tags,
		Type:     meta.Type,
	}
	if meta.Title != "" {
		n.Title = meta.Title
	}
	return n, nil
}
