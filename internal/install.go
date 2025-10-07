package internal

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/vninomtz/pkms/internal/store"
)

const CLI_NAME = "pkms"
const PKMS_HOME_DIR = "PKMS_HOME_DIR"
const PKMS_NOTES_DIR = "PKMS_NOTES_DIR"
const DB_FILENAME = "pkms.db"

func Install() error {
	home := HomePath()
	if _, err := os.Stat(home); os.IsNotExist(err) {
		if err := os.MkdirAll(home, 0755); err != nil {
			return fmt.Errorf("error to create directory %s: %w", home, err)
		}
	}
	db_path := filepath.Join(home, DB_FILENAME)

	file, err := os.Create(db_path)
	if err != nil {
		return fmt.Errorf("Error creating DB at %s: %w", home, err)
	}
	defer file.Close()

	st, err := store.New(db_path)
	if err != nil {
		return err
	}
	defer st.Close()

	err = st.Setup()
	if err != nil {
		return fmt.Errorf("Error to Setup db schema: %w", err)
	}

	return nil
}

func HomePath() string {
	base_dir := baseDir()
	return filepath.Join(base_dir, fmt.Sprintf(".%s", CLI_NAME))
}

func DatabasePath() string {
	return filepath.Join(HomePath(), DB_FILENAME)
}
func NotesPath() string {
	return os.Getenv(PKMS_NOTES_DIR)
}

func setup(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("error to create directory %s: %w", dir, err)
		}
	}
	return nil
}
func baseDir() string {
	// Check env variable or use $HOME as default
	dir := os.Getenv(PKMS_HOME_DIR)
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			panic("Error to get $HOME directory")
		}
		dir = home
	}
	return dir
}
