package main

import (
	"flag"
	"fmt"

	"github.com/vninomtz/swe-notes/mdreader"
)

func main()  {
  flag.Parse()
  if len(flag.Args()) == 0 {
    fmt.Printf("Usage: <file>")
    return
  }
  file := flag.Args()[0]
  err := mdreader.Read(file)
  if err != nil {
    fmt.Println(err)
  }
}
