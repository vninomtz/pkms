package internal

import (
	"bytes"
	"html/template"

	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
)

type content struct {
	Title string
	Body  template.HTML
}

func ParseNodeToHTML(input []byte, templateFile string) ([]byte, error) {
	var buf bytes.Buffer

	md := goldmark.New(
		goldmark.WithExtensions(highlighting.NewHighlighting(
			highlighting.WithStyle("dracula"),
		)),
	)
	if err := md.Convert(input, &buf); err != nil {
		return nil, err
	}

	t, err := template.ParseFiles(templateFile)
	if err != nil {
		return nil, err
	}

	//Instantiate the content type, adding the title and body
	c := content{
		Title: "Markdown Preview Tool",
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
