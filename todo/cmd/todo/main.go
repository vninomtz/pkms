package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/vninomtz/swe-notes/todo"
)

// File name

var todoFileName = ".todo.json"

const FILE_NAME_ENV = "TODO_FILENAME"

func main()  {
  // check if the user defined the ENV VAR for a custom file name
  if os.Getenv(FILE_NAME_ENV) != "" {
    todoFileName = os.Getenv(FILE_NAME_ENV)
  }
  flag.Usage = func() {
    fmt.Fprintf(flag.CommandLine.Output(), "%s tool. Developed for vninomtz\n", os.Args[0])
    fmt.Fprintf(flag.CommandLine.Output(), "Copyright 2023\n")
    fmt.Fprintf(flag.CommandLine.Output(), "Usage information:\n")
    flag.PrintDefaults()
  }
  // Parsing command line flags
  task := flag.String("task", "", "Task to be included in the ToDo list")
  list := flag.Bool("list", false, "List all tasks")
  complete := flag.Int("complete", 0, "Item to be completed")

  flag.Parse()

  l := &todo.Tasks{}

  // Read items from file
  if err := l.Get(todoFileName); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }

  switch {
  case *list:
    // List current to do items
    fmt.Print(l)

  case *complete > 0:
    // complete the given item
    if err := l.Complete(*complete); err != nil {
      fmt.Fprintln(os.Stderr, err)
      os.Exit(1)
    }
    if err := l.Save(todoFileName); err != nil {
      fmt.Fprintln(os.Stderr, err)
      os.Exit(1)
    }
  
  case *task != "":
    l.Add(*task)

    if err := l.Save(todoFileName); err != nil {
      fmt.Fprintln(os.Stderr, err)
      os.Exit(1)
    }
  default:
    fmt.Fprintln(os.Stderr, "Invalid option")
    os.Exit(1)

  }
}
