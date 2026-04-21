package cmd

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/vninomtz/pkms/internal/config"
)

var syncFiles = []string{
	"exercises.yml",
	"routines.yml",
	"schedule.yml",
	"resources.yml",
}

func SyncCommand(args []string) {
	fs := flag.NewFlagSet("sync", flag.ExitOnError)
	dryRun := fs.Bool("dry-run", false, "Preview changes without writing files")
	fs.Parse(args)

	cfg := config.New()
	cfg.Load()

	if cfg.BycuriosityDir == "" {
		log.Fatal("PKMS_BYCURIOSITY_DIR is not set. Example:\n  export PKMS_BYCURIOSITY_DIR=/path/to/bycuriosity/website/src/content")
	}

	fmt.Printf("Syncing notes → bycuriosity\n")
	fmt.Printf("  Source : %s\n", cfg.NotesDir)
	fmt.Printf("  Target : %s\n\n", cfg.BycuriosityDir)

	synced := 0
	for _, file := range syncFiles {
		src := filepath.Join(cfg.NotesDir, file)
		dst := filepath.Join(cfg.BycuriosityDir, file)

		if _, err := os.Stat(src); os.IsNotExist(err) {
			fmt.Printf("  skip  %s (not found in notes)\n", file)
			continue
		}

		if *dryRun {
			fmt.Printf("  would copy  %s → %s\n", src, dst)
			synced++
			continue
		}

		if err := copyFile(src, dst); err != nil {
			fmt.Printf("  error  %s: %v\n", file, err)
			continue
		}

		fmt.Printf("  copied  %s\n", file)
		synced++
	}

	fmt.Printf("\n%d/%d files synced", synced, len(syncFiles))
	if *dryRun {
		fmt.Print(" (dry-run, no files written)")
	}
	fmt.Println()
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
