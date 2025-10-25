package cmd

import (
	"flag"
	"fmt"
	"log"
	"os/exec"

	"github.com/vninomtz/pkms/internal/config"
	"github.com/vninomtz/pkms/internal/notes"
)

func PublishCommand(args []string) {
	cfg := config.New()
	cfg.Load()

	fs := flag.NewFlagSet("publish", flag.ExitOnError)
	out := fs.String("o", "", "Directory to move the notes")
	fs.Parse(args)

	if *out == "" {
		log.Fatal("Missing out directory")
	}

	srv := notes.New(cfg.NotesDir)

	res, err := srv.GetPublic()
	if err != nil {
		log.Fatal(err)
	}

	copied := 0
	for _, n := range res {
		cmd := exec.Command("cp", n.Entry.Path, *out)
		_, err := cmd.Output()
		if err != nil {
			fmt.Println(err)
		} else {
			copied++
		}
	}

	fmt.Printf("%d copied to %s of %d notes\n", copied, *out, len(res))
}
