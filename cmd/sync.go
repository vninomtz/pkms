package cmd

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/vninomtz/pkms/internal/config"
	"github.com/vninomtz/pkms/internal/notes"
)

var syncFiles = map[string]string{
	"exercises.yml": "exercises/exercises.yml",
	"routines.yml":  "routines/routines.yml",
	"schedule.yml":  "schedule/schedule.yml",
	"resources.yml": "resources/resources.yml",
}

func SyncCommand(args []string) {
	startTime := time.Now()

	fs := flag.NewFlagSet("sync", flag.ExitOnError)
	dryRun := fs.Bool("dry-run", false, "Preview changes without writing files")
	fs.Parse(args)

	cfg := config.New()
	cfg.Load()

	if cfg.BycuriosityDir == "" {
		log.Fatal("PKMS_BYCURIOSITY_DIR is not set. Example:\n  export PKMS_BYCURIOSITY_DIR=/path/to/bycuriosity/website/src/content")
	}

	fmt.Printf("Syncing content → bycuriosity\n")
	fmt.Printf("  Source : %s\n", cfg.NotesDir)
	fmt.Printf("  Target : %s\n\n", cfg.BycuriosityDir)

	// Sync YAML files
	fmt.Println("YAML Files:")
	syncedYAML := 0
	for file, subpath := range syncFiles {
		src := filepath.Join(cfg.NotesDir, file)
		dst := filepath.Join(cfg.BycuriosityDir, subpath)

		if _, err := os.Stat(src); os.IsNotExist(err) {
			fmt.Printf("  skip  %s (not found in notes)\n", file)
			continue
		}

		if *dryRun {
			fmt.Printf("  would copy  %s → %s\n", src, dst)
			syncedYAML++
			continue
		}

		if err := copyFile(src, dst); err != nil {
			fmt.Printf("  error  %s: %v\n", file, err)
			continue
		}

		fmt.Printf("  copied  %s\n", file)
		syncedYAML++
	}

	// Sync Notes
	fmt.Println("\nNotes:")
	srv := notes.New(cfg.NotesDir)

	// Get public notes count for preview
	publicNotes, err := srv.GetPublic()
	if err != nil {
		fmt.Printf("  error loading notes: %v\n", err)
	} else {
		fmt.Printf("  Found %d public notes to sync\n", len(publicNotes))
	}

	notesStats, err := syncNotes(srv, filepath.Join(cfg.BycuriosityDir, "notes"), *dryRun)
	if err != nil {
		fmt.Printf("  error syncing notes: %v\n", err)
		notesStats = SyncStats{}
	}

	// Summary
	duration := time.Since(startTime).Seconds()
	fmt.Printf("\n=== Summary ===\n")
	fmt.Printf("YAML: %d/%d files synced\n", syncedYAML, len(syncFiles))

	totalNoteChanges := notesStats.Added + notesStats.Updated
	if notesStats.Deleted == 0 {
		fmt.Printf("Notes: %d public notes synced to ByCuriosity\n", totalNoteChanges)
	} else {
		fmt.Printf("Notes synced: %d to add/update, %d to delete (made private or removed)\n",
			totalNoteChanges, notesStats.Deleted)
	}
	fmt.Printf("Total time: %.2fs\n", duration)

	if *dryRun {
		fmt.Printf("\n📋 Dry-run mode — no files were written\n")
	} else {
		fmt.Printf("\n✓ Sync complete\n")
	}
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open source: %w", err)
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("create destination dir: %w", err)
	}

	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("create destination: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("copy: %w", err)
	}

	return out.Sync()
}

// SyncStats tracks the results of syncing notes
type SyncStats struct {
	Added   int
	Updated int
	Deleted int
}

// slugify converts a title to a URL-safe filename
func slugify(title string) string {
	// Convert to lowercase
	s := strings.ToLower(title)
	// Replace spaces with hyphens
	s = strings.ReplaceAll(s, " ", "-")
	// Remove special characters, keep only alphanumeric and hyphens
	reg := regexp.MustCompile("[^a-z0-9-]+")
	s = reg.ReplaceAllString(s, "")
	// Replace multiple hyphens with single hyphen
	s = regexp.MustCompile("-+").ReplaceAllString(s, "-")
	// Trim hyphens from edges
	s = strings.Trim(s, "-")
	return s
}

// getOutputFilename determines the output filename for a note
func getOutputFilename(note notes.Note) string {
	if note.Title != "" {
		return slugify(note.Title) + ".md"
	}
	return note.Entry.Filename
}

// listNotesInDirectory returns all markdown files in a directory
func listNotesInDirectory(dir string) (map[string]bool, error) {
	files := make(map[string]bool)

	// Create directory if it doesn't exist
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return files, nil
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return files, err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			files[entry.Name()] = true
		}
	}

	return files, nil
}

// syncNotes synchronizes public notes from PKMS to ByCuriosity
func syncNotes(srv notes.NoteService, targetDir string, dryRun bool) (SyncStats, error) {
	stats := SyncStats{}

	// Create target directory if it doesn't exist
	if !dryRun {
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return stats, fmt.Errorf("create target dir: %w", err)
		}
	}

	// Get all public notes
	publicNotes, err := srv.GetPublic()
	if err != nil {
		return stats, err
	}

	// Track which files we're syncing
	syncedFiles := make(map[string]bool)

	// Sync each public note
	for _, note := range publicNotes {
		outputFilename := getOutputFilename(note)
		outputPath := filepath.Join(targetDir, outputFilename)

		// Check if file already exists (for update detection)
		_, exists := os.Stat(outputPath)
		isNew := os.IsNotExist(exists)

		syncedFiles[outputFilename] = true

		if dryRun {
			if isNew {
				fmt.Printf("  would add    %s\n", outputFilename)
				stats.Added++
			} else {
				fmt.Printf("  would update %s\n", outputFilename)
				stats.Updated++
			}
			continue
		}

		// Write note file with preserved frontmatter
		if err := os.WriteFile(outputPath, note.Entry.Content, 0644); err != nil {
			fmt.Printf("  error writing %s: %v\n", outputFilename, err)
			continue
		}

		if isNew {
			fmt.Printf("  added    %s\n", outputFilename)
			stats.Added++
		} else {
			fmt.Printf("  updated  %s\n", outputFilename)
			stats.Updated++
		}
	}

	// Find and delete files that are no longer public
	existingFiles, err := listNotesInDirectory(targetDir)
	if err != nil {
		return stats, err
	}

	for filename := range existingFiles {
		if !syncedFiles[filename] {
			filePath := filepath.Join(targetDir, filename)

			if dryRun {
				fmt.Printf("  would delete %s\n", filename)
				stats.Deleted++
				continue
			}

			if err := os.Remove(filePath); err != nil {
				fmt.Printf("  error deleting %s: %v\n", filename, err)
				continue
			}

			fmt.Printf("  deleted  %s\n", filename)
			stats.Deleted++
		}
	}

	return stats, nil
}
