package internal

import (
	"bytes"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"html/template"
	"path/filepath"
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

func (p *parser) Parse(name string, _content []byte) ([]byte, error) {
	html, err := p.MDToHTML(_content)
	if err != nil {
		return nil, err
	}
	c := content{
		Title: name,
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
