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

	"net/url"

	"github.com/vninomtz/pkms/internal/bookmarks"
	"github.com/vninomtz/pkms/internal/config"
	"github.com/vninomtz/pkms/internal/crawler"
	"github.com/vninomtz/pkms/internal/notes"
)

var syncFiles = map[string]string{
	"exercises.yml":  "exercises/exercises.yml",
	"routines.yml":   "routines/routines.yml",
	"schedule.yml":   "schedule/schedule.yml",
	"resources.yml":  "resources/resources.yml",
	"bookmarks.yml":  "bookmarks/bookmarks.yml",
}

func SyncCommand(args []string) {
	startTime := time.Now()

	fs := flag.NewFlagSet("sync", flag.ExitOnError)
	dryRun := fs.Bool("dry-run", false, "Preview changes without writing files")
	retryFailed := fs.Bool("retry-failed", false, "Re-fetch bookmarks that previously failed")
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

	// Sync Bookmarks
	fmt.Println("\nBookmarks:")
	bookmarkStats, err := syncBookmarks(srv, filepath.Join(cfg.NotesDir, "bookmarks.yml"), *dryRun, *retryFailed)
	if err != nil {
		fmt.Printf("  error syncing bookmarks: %v\n", err)
		bookmarkStats = bookmarks.BookmarkStats{}
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

	if bookmarkStats.Discovered > 0 {
		fmt.Printf("Bookmarks: %d discovered, %d fetched (new), %d cached (reused), %d updated (refreshed)",
			bookmarkStats.Discovered, bookmarkStats.Fetched, bookmarkStats.Cached, bookmarkStats.Updated)
		if bookmarkStats.Failed > 0 {
			fmt.Printf(", %d failed", bookmarkStats.Failed)
		}
		fmt.Println()
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

// skipDomains is a list of domains that require auth or block scrapers
var skipDomains = []string{
	"claude.ai",
	"chat.openai.com",
	"chatgpt.com",
	"twitter.com",
	"x.com",
	"linkedin.com",
	"facebook.com",
	"instagram.com",
	"notion.so",
	"figma.com",
	"docs.google.com",
	"drive.google.com",
	"mail.google.com",
	"localhost",
	"127.0.0.1",
}

// isSkippedDomain returns true if the URL's domain should be skipped
func isSkippedDomain(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	host := strings.ToLower(strings.TrimPrefix(u.Hostname(), "www."))
	for _, skip := range skipDomains {
		if host == skip || strings.HasSuffix(host, "."+skip) {
			return true
		}
	}
	return false
}

// extractUniqueURLs collects all unique URLs from all notes, excluding skipped domains
func extractUniqueURLs(srv notes.NoteService) (map[string]bool, error) {
	urls := make(map[string]bool)

	allNotes, err := srv.GetAll()
	if err != nil {
		return urls, err
	}

	for _, note := range allNotes {
		for _, u := range note.Links {
			if !isSkippedDomain(u) {
				urls[u] = true
			}
		}
	}

	return urls, nil
}

// syncBookmarks discovers and syncs bookmarks from notes
func syncBookmarks(srv notes.NoteService, bookmarksPath string, dryRun bool, retryFailed bool) (bookmarks.BookmarkStats, error) {
	stats := bookmarks.BookmarkStats{}

	// Extract all unique URLs from notes
	discoveredURLs, err := extractUniqueURLs(srv)
	if err != nil {
		return stats, fmt.Errorf("error extracting URLs: %w", err)
	}

	stats.Discovered = len(discoveredURLs)
	if stats.Discovered == 0 {
		fmt.Println("  no URLs found in notes")
		return stats, nil
	}

	fmt.Printf("  found %d unique URLs in notes\n", stats.Discovered)

	// Load existing bookmarks
	bookmarksSrv := bookmarks.New(bookmarksPath)
	bookmarkList, err := bookmarksSrv.LoadBookmarks()
	if err != nil {
		return stats, fmt.Errorf("error loading bookmarks: %w", err)
	}

	// Build map of existing bookmarks by URL
	existingByURL := make(map[string]*bookmarks.Bookmark)
	for i, b := range bookmarkList.Bookmarks {
		existingByURL[b.URL] = &bookmarkList.Bookmarks[i]
	}

	// Categorize URLs: new, cached, need refresh
	var toFetch []string

	for url := range discoveredURLs {
		if existing, found := existingByURL[url]; found {
			if existing.Status == "failed" && retryFailed {
				stats.Fetched++
				toFetch = append(toFetch, url)
			} else if bookmarksSrv.IsCached(existing) {
				stats.Cached++
			} else {
				stats.Updated++
				toFetch = append(toFetch, url)
			}
		} else {
			stats.Fetched++
			toFetch = append(toFetch, url)
		}
	}

	if dryRun {
		label := "new/stale"
		if retryFailed {
			label = "new/stale/failed"
		}
		fmt.Printf("  would fetch %d %s URLs, reuse %d cached\n", len(toFetch), label, stats.Cached)
		fmt.Printf("  would save bookmarks.yml\n")
		return stats, nil
	}

	// Fetch metadata for new and stale URLs
	if len(toFetch) > 0 {
		fmt.Printf("  fetching metadata for %d URLs...\n", len(toFetch))
		pages, fetchErr := crawler.FetchMultiple(toFetch)
		if fetchErr != nil {
			return stats, fmt.Errorf("error fetching URLs: %w", fetchErr)
		}

		// Track which URLs were successfully fetched
		fetchedURLs := make(map[string]bool)

		// Process fetched pages
		for _, page := range pages {
			fetchedURLs[page.URL] = true

			metadata, err := crawler.ParseHtml(page.HTML)
			if err != nil {
				// Fetch succeeded but parse failed — mark as failed
				if existing, found := existingByURL[page.URL]; found {
					bookmarksSrv.MarkFailed(existing, err.Error())
				} else {
					b := bookmarksSrv.AddFailedBookmark(page.URL, err.Error())
					bookmarkList.Bookmarks = append(bookmarkList.Bookmarks, *b)
					existingByURL[page.URL] = b
				}
				stats.Failed++
				continue
			}

			if existing, found := existingByURL[page.URL]; found {
				bookmarksSrv.UpdateMetadata(existing, metadata)
			} else {
				b := bookmarksSrv.AddBookmark(page.URL, metadata)
				bookmarkList.Bookmarks = append(bookmarkList.Bookmarks, *b)
				existingByURL[page.URL] = b
			}
		}

		// URLs in toFetch that got no response — network error / timeout
		for _, u := range toFetch {
			if fetchedURLs[u] {
				continue
			}
			if existing, found := existingByURL[u]; found {
				bookmarksSrv.MarkFailed(existing, "fetch failed (network error or timeout)")
			} else {
				b := bookmarksSrv.AddFailedBookmark(u, "fetch failed (network error or timeout)")
				bookmarkList.Bookmarks = append(bookmarkList.Bookmarks, *b)
				existingByURL[u] = b
			}
			stats.Failed++
		}
	}

	// Remove bookmarks for URLs no longer in notes
	var updatedBookmarks []bookmarks.Bookmark
	for _, b := range bookmarkList.Bookmarks {
		if discoveredURLs[b.URL] {
			updatedBookmarks = append(updatedBookmarks, b)
		}
	}
	bookmarkList.Bookmarks = updatedBookmarks

	// Save bookmarks
	if err := bookmarksSrv.SaveBookmarks(bookmarkList); err != nil {
		return stats, fmt.Errorf("error saving bookmarks: %w", err)
	}

	return stats, nil
}
