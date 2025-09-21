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
	return write(filename, content, out)
}

func write(name string, content []byte, dir_out string) error {
	err := os.Mkdir(dir_out, 0755)
	if err != nil {
		return err
	}
	out := filepath.Join(dir_out, name)
	return os.WriteFile(out, content, 0644)
}
