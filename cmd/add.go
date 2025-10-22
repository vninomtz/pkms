package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/vninomtz/pkms/internal/config"
	"github.com/vninomtz/pkms/internal/notes"
)

func AddCommand(args []string) {
	fs := flag.NewFlagSet("add", flag.ExitOnError)
	fs.Parse(args)

	cfg := config.New()
	cfg.Load()

	srv := notes.New(cfg.NotesDir)

	input, err := ReadInputFromEditor()
	if err != nil {
		log.Fatal(err)
	}
	path, err := srv.New(input)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Document created %s\n", path)
}

func ReadInputFromEditor() ([]byte, error) {
	file, err := os.CreateTemp(os.TempDir(), "pkms")
	if err != nil {
		return []byte{}, err
	}
	filename := file.Name()
	defer os.Remove(filename)

	if err = file.Close(); err != nil {
		return []byte{}, err
	}

	if err = OpenFileInEditor(filename); err != nil {
		return []byte{}, err
	}

	bytes, err := os.ReadFile(filename)
	if err != nil {
		return []byte{}, err
	}

	return bytes, nil
}

func OpenFileInEditor(filename string) error {
	editor := "vim"
	path, err := exec.LookPath(editor)
	if err != nil {
		return err
	}
	cmd := exec.Command(path, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
