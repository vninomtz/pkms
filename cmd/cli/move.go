package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/vninomtz/pkms/internal"
)

func MoveNotes(nodes []internal.Node, output string) {
	saved := 0
	for _, n := range nodes {
		path := filepath.Join(output, fmt.Sprintf("%s.md", n.Title))
		err := os.WriteFile(path, n.Bytes, 0644)
		if err != nil {
			fmt.Printf("Error writing to %s\n", path)
		} else {
			saved++
		}
	}
	fmt.Printf("%d saved of %d \n", saved, len(nodes))
}
