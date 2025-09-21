package cmd

import (
	"flag"
	"fmt"
	"log"

	"github.com/vninomtz/pkms/internal"
)

func InstallCommand(args []string) {
	fs := flag.NewFlagSet("install", flag.ExitOnError)
	fs.Parse(args)

	err := internal.Install()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("PKMS installled at %s\n", internal.HomePath())
}
