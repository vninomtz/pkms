package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/vninomtz/swe-notes/api"
	"github.com/vninomtz/swe-notes/walker"
)

func main()  {
  root := flag.String("root", ".", "Root directory to start")
  ext := flag.String("ext", ".md", "File extension to filter out")
  size := flag.Int64("size", 0, "Minimum file size")
  cmd := flag.String("cmd", "", "Command to execute")
  port := flag.String("port", "8000", "Port for http server")
  host := flag.String("host", "", "Server host")

  flag.Parse()

  c := walker.Config{
    Root: *root,
    Ext: *ext,
    Size: *size,
    Cmd: *cmd,
  }
  notes, err := walker.Run(*root, os.Stdout, c);
  if err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }

  handler := api.NewNoteHandler(notes)

  http.Handle("/api/notes", handler)

  addr := fmt.Sprintf("%s:%s",*host,*port)

  fmt.Printf("Server runing at: http://localhost:%s", *port)
  err = http.ListenAndServe(addr, nil)
  if err != nil {
    fmt.Println(err)
    return
  }
}
