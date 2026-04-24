package cmd

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/vninomtz/pkms/internal/config"
	"github.com/vninomtz/pkms/internal/notes"
)

func BookmarksCommand(args []string) {
	startTime := time.Now()

	fs := flag.NewFlagSet("bookmarks", flag.ExitOnError)
	dryRun := fs.Bool("dry-run", false, "Preview changes without writing files")
	retryFailed := fs.Bool("retry-failed", false, "Re-fetch bookmarks that previously failed")
	fs.Parse(args)

	cfg := config.New()
	cfg.Load()

	fmt.Printf("Syncing bookmarks\n")
	fmt.Printf("  Source : %s\n", cfg.NotesDir)
	if cfg.BycuriosityDir != "" {
		fmt.Printf("  Target : %s\n", filepath.Join(cfg.BycuriosityDir, "bookmarks/bookmarks.yml"))
	}
	fmt.Println()

	srv := notes.New(cfg.NotesDir)
	bookmarksPath := filepath.Join(cfg.NotesDir, "bookmarks.yml")

	stats, err := syncBookmarks(srv, bookmarksPath, *dryRun, *retryFailed)
	if err != nil {
		log.Fatalf("error syncing bookmarks: %v", err)
	}

	// Copy to ByCuriosity if configured
	if cfg.BycuriosityDir != "" && !*dryRun {
		dst := filepath.Join(cfg.BycuriosityDir, "bookmarks/bookmarks.yml")
		if err := copyFile(bookmarksPath, dst); err != nil {
			fmt.Printf("  warning: could not copy to bycuriosity: %v\n", err)
		} else {
			fmt.Printf("  copied bookmarks.yml → bycuriosity\n")
		}
	}

	duration := time.Since(startTime).Seconds()
	fmt.Printf("\n=== Summary ===\n")
	if stats.Discovered > 0 {
		fmt.Printf("Bookmarks: %d discovered, %d fetched (new), %d cached (reused), %d updated (refreshed)",
			stats.Discovered, stats.Fetched, stats.Cached, stats.Updated)
		if stats.Failed > 0 {
			fmt.Printf(", %d failed", stats.Failed)
		}
		fmt.Println()
	}
	fmt.Printf("Total time: %.2fs\n", duration)

	if *dryRun {
		fmt.Printf("\n📋 Dry-run mode — no files were written\n")
	} else {
		fmt.Printf("\n✓ Bookmarks sync complete\n")
	}
}
