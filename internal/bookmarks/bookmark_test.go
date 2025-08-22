package bookmarks

import (
	"bytes"
	"testing"
)

const HTML_EXAMPLE = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Ejemplo de Página</title>
  <meta name="description" content="Esta es una página de prueba para test unitarios.">

  <!-- Open Graph -->
  <meta property="og:title" content="Título OG de Ejemplo">
  <meta property="og:description" content="Descripción OG de ejemplo">
  <meta property="og:image" content="https://example.com/imagen.jpg">
  <meta property="og:url" content="https://example.com/page">

  <!-- Twitter Card -->
  <meta name="twitter:card" content="summary_large_image">
  <meta name="twitter:title" content="Título Twitter de Ejemplo">
  <meta name="twitter:description" content="Descripción Twitter de ejemplo">
  <meta name="twitter:image" content="https://example.com/twitter-imagen.jpg">
</head>
<body>
  <h1>Hola Mundo</h1>
  <p>Este es un HTML de ejemplo para probar la extracción de metadatos.</p>
</body>
</html>
`

func TestExtractBookmark(t *testing.T) {
	expectedBK := Bookmark{
		Title:       "Ejemplo de Página",
		Description: "Esta es una página de prueba para test unitarios.",
		Image:       "https://example.com/imagen.jpg",
		Url:         "https://example.com/page",
	}
	r := bytes.NewReader([]byte(HTML_EXAMPLE))
	bk := ExtractMetadata(r)

	if expectedBK.Title != bk.Title {
		t.Errorf("Expected %s, got %s instead\n", expectedBK.Title, bk.Title)
	}
	if expectedBK.Description != bk.Description {
		t.Errorf("Expected %s, got %s instead\n", expectedBK.Description, bk.Description)
	}
	if expectedBK.Url != bk.Url {
		t.Errorf("Expected %s, got %s instead\n", expectedBK.Url, bk.Url)
	}
	if expectedBK.Image != bk.Image {
		t.Errorf("Expected %s, got %s instead\n", expectedBK.Image, bk.Image)
	}
}
