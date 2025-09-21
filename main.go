package main

import (
	"fmt"
	"os"

	"github.com/vninomtz/pkms/cmd"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Use: pkm <command> [options]")
		os.Exit(1)
	}

	// First arg is the binary, second is the subcommand
	subcommand := os.Args[1]
	args := os.Args[2:] // rest go to subcommand

	switch subcommand {
	case "docs":
		cmd.DocsCommand(args)
	case "inspect":
		cmd.InspectCommand(args)
	case "install":
		cmd.InstallCommand(args)
	default:
		fmt.Printf("Unknown command: %s\n", subcommand)
		os.Exit(1)
	}
}
