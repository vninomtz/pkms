package main

import (
	"fmt"
	"os"

	"github.com/vninomtz/pkms/cmd"
)

const PKM_VERSION = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Use: pkm <command> [options]")
		os.Exit(1)
	}

	// First arg is the binary, second is the subcommand
	subcommand := os.Args[1]
	args := os.Args[2:] // rest go to subcommand

	switch subcommand {
	case "add":
		cmd.AddCommand(args)
	case "search":
		cmd.SearchCommand(args)
	case "inspect":
		cmd.InspectCommand(args)
	case "install":
		cmd.InstallCommand(args)
	case "publish":
		cmd.PublishCommand(args)
	case "index":
		cmd.IndexCommand(args)
	case "version":
		fmt.Printf("PKM version %s\n", PKM_VERSION)
	default:
		fmt.Printf("Unknown command: %s\n", subcommand)
		os.Exit(1)
	}
}
