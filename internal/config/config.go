package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const CLI_NAME = "pkms"
const PKMS_HOME_DIR = "PKMS_HOME_DIR"
const PKMS_NOTES_DIR = "PKMS_NOTES_DIR"
const DB_FILENAME = "pkms.db"

type config struct {
	HomeDir    string
	NotesDir   string
	SQLiteFile string
	CLIDir     string
}

func New() *config {
	return &config{}
}

func (c *config) Load() {

	home, err := os.UserHomeDir()
	if err != nil {
		panic("Error to get $HOME directory")
	}
	c.HomeDir = home
	c.CLIDir = filepath.Join(home, fmt.Sprintf(".%s", CLI_NAME))
	c.SQLiteFile = filepath.Join(c.CLIDir, DB_FILENAME)
	c.NotesDir = os.Getenv(PKMS_NOTES_DIR)
	if c.NotesDir == "" {
		panic("Error to get $PKMS_HOME_DIR env variable")
	}
}
