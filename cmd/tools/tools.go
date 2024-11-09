package main

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/vninomtz/pkms/internal"
)

func main() {
	path := flag.String("p", "", "Path name")
	flag.Parse()

	if err := Run(*path); err != nil {
		log.Println(err)
	}
}

func Run(path string) error {
	if path == "" {
		return errors.New("Path is required")
	}
	collector := internal.NewCollector(path, "")
	nodes, err := collector.Collect()

	if err != nil {
		log.Println(err)
		return err
	}

	for i, n := range nodes {
		fmt.Printf("%d.- %d  %s, %s %v\n", i, n.Size, n.Name, n.Parent, n.Meta)
	}
	return nil
}
