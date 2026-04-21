package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/vninomtz/pkms/internal/config"
	"github.com/vninomtz/pkms/internal/reading"
)

func AddBookCommand(args []string) {
	fs := flag.NewFlagSet("add-book", flag.ExitOnError)
	fs.Parse(args)

	cfg := config.New()
	cfg.Load()
	filePath := filepath.Join(cfg.NotesDir, "resources.yml")
	srv := reading.New(filePath)

	fmt.Println("\n📚 Agregar un nuevo libro a tus lecturas\n")

	title := reading.PromptString("Título del libro: ")
	if title == "" {
		fmt.Println("❌ El título es requerido")
		os.Exit(1)
	}

	author := reading.PromptString("Autor: ")
	if author == "" {
		fmt.Println("❌ El autor es requerido")
		os.Exit(1)
	}

	existingBook := srv.FindBook(title, author)
	if existingBook != nil {
		fmt.Printf("\n⚠️  El libro '%s' de %s ya existe en tu lista\n", title, author)
		fmt.Printf("ID: %s | Estado: %s\n\n", existingBook.ID, existingBook.Status)

		addReread := reading.PromptBool("¿Quieres agregar una relectura de este libro?")
		if !addReread {
			fmt.Println("Cancelado.")
			return
		}

		fmt.Println("\nAñadiendo relectura...")
		newBook := &reading.Book{
			Title:         title,
			Author:        author,
			Status:        "completed",
			Language:      existingBook.Language,
			Categories:    existingBook.Categories,
			CompletedDate: time.Now().Format("2006-01-02"),
			Notes:         "Relectura",
		}

		if err := srv.AddBook(newBook); err != nil {
			fmt.Printf("❌ Error al agregar relectura: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("\n✓ Relectura agregada exitosamente\n")
		fmt.Printf("📖 UUID: %s\n", newBook.ID)
		fmt.Printf("📅 Fecha: %s\n\n", newBook.CompletedDate)
		return
	}

	completed := reading.PromptBool("¿Ya lo completaste?")

	language := reading.PromptStringWithDefault("Idioma", "es")

	categories := reading.PromptCategories("Categorías")

	url := reading.PromptString("URL (opcional): ")

	notes := reading.PromptString("Notas (opcional): ")

	completedDate := ""
	if completed {
		completedDate = time.Now().Format("2006-01-02")
	}

	status := "pending"
	if completed {
		status = "completed"
	}

	book := &reading.Book{
		Title:         title,
		Author:        author,
		Status:        status,
		Language:      language,
		Categories:    categories,
		CompletedDate: completedDate,
		URL:           url,
		Notes:         notes,
	}

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("📋 PREVIEW")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Título: %s\n", book.Title)
	fmt.Printf("Autor: %s\n", book.Author)
	fmt.Printf("Estado: %s\n", statusToSpanish(book.Status))

	if book.CompletedDate != "" {
		fmt.Printf("Fecha completada: %s\n", book.CompletedDate)
	}

	fmt.Printf("Idioma: %s\n", book.Language)

	if len(book.Categories) > 0 {
		fmt.Printf("Categorías: %s\n", strings.Join(book.Categories, ", "))
	}

	if book.URL != "" {
		fmt.Printf("URL: %s\n", book.URL)
	}

	if book.Notes != "" {
		fmt.Printf("Notas: %s\n", book.Notes)
	}

	fmt.Println(strings.Repeat("=", 50) + "\n")

	if !reading.PromptBool("¿Confirmar?") {
		fmt.Println("Cancelado.")
		return
	}

	if err := srv.AddBook(book); err != nil {
		fmt.Printf("❌ Error al agregar libro: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n✅ Libro agregado exitosamente\n")
	fmt.Printf("📚 Título: %s\n", book.Title)
	fmt.Printf("✍️  Autor: %s\n", book.Author)
	fmt.Printf("📖 UUID: %s\n", book.ID)
	fmt.Printf("🔖 Estado: %s\n", statusToSpanish(book.Status))

	if book.CompletedDate != "" {
		fmt.Printf("📅 Fecha completada: %s\n", book.CompletedDate)
	}

	if len(book.Categories) > 0 {
		fmt.Printf("📌 Categorías: %s\n", strings.Join(book.Categories, ", "))
	}

	fmt.Println()
}

func statusToSpanish(status string) string {
	if status == "completed" {
		return "Completado"
	}
	return "Pendiente"
}

