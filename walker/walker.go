package walker

import (
	"io"
	"io/fs"
	"path/filepath"

	"github.com/vninomtz/swe-notes/mdreader"
	"github.com/vninomtz/swe-notes/model"
)

type Config struct {
  // root directory to search
  Root string
  // extenstion to filter out: .md
  Ext string
  // min file size
  Size int64
  // command to execute
  Cmd string
}

func Run(root string, out io.Writer, cfg Config) ([]*model.Note, error)  {
  notes := []*model.Note{}
  
  err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
    if err != nil {
      return err
    }

    if cfg.Ext != "" && filepath.Ext(path) == cfg.Ext {
      dir := filepath.Dir(path)
      parent := filepath.Base(dir)
      if parent == "." {
        parent = ""
      }
      md, err := mdreader.Read(path)
      if err != nil {
        return err
      }
      nt := &model.Note{
        Name: info.Name(),
        Content: md.Content,
        Meta: md.Meta,
        Size: info.Size(),
        Parent: parent,
      }
      notes = append(notes, nt)
    }

    return nil
  })

  return notes, err
}
