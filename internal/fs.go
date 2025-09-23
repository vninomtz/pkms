package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

const OUTPUT_DIR = "output"

func WriteHtml(name string, content []byte) error {
	out := filepath.Join(HomePath(), OUTPUT_DIR)
	filename := fmt.Sprintf("%s.html", name)
	_, err := write(filename, content, out)
	return err
}
func WriteNote(name string, content []byte, out string) (string, error) {
	filename := fmt.Sprintf("%s.md", name)
	return write(filename, content, out)
}

func write(name string, content []byte, dir_out string) (string, error) {
	err := os.MkdirAll(dir_out, 0755)
	if err != nil {
		return "", err
	}
	out := filepath.Join(dir_out, name)
	err = os.WriteFile(out, content, 0644)
	return out, err
}
