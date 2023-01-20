package mdreader

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"gopkg.in/yaml.v3"
  "fmt"
)

func Read(file string) (*MdFile, error) {
  var err error
  
  f, err := os.Open(file)
  if err != nil {
    return nil, err
  }
  defer f.Close()

  reader := newMdReader(f)

  err = reader.extract()
  if err != nil {
    return nil, err
  }

  meta := make(map[string]interface{})
  err = yaml.Unmarshal(reader.Meta(), &meta)
  if err != nil {
    fmt.Println(err)
  }
  md := &MdFile{
    Meta: meta,
    Content: string(reader.Content()) ,
  }
  return md, nil

}

type mdReader struct {
  reader *bufio.Reader
  output *bytes.Buffer

  current int
  start int
  end int

  init int

  meta bool
}

type MdFile struct{
  Meta map[string]interface{}
  Content string
}

func newMdReader(r io.Reader) *mdReader{
  return &mdReader{
    reader: bufio.NewReader(r),
    output: bytes.NewBuffer(nil),
  }
}


func (m *mdReader) Meta() ([]byte)  {
  if m.meta {
    return m.output.Bytes()[m.start:m.end]
  }
  return nil
}

func (m *mdReader) Content() ([]byte)  {
  return m.output.Bytes()[m.init:]
}

func (m *mdReader) extract() (error)  {
  var openMeta = false
  for {
    line, isEOF, err := m.readLine()

    if err != nil {
      return err
    }

    if isEOF {
      return nil
    }
    str := string(bytes.TrimSpace(line))
    if str == "---" {
      if !openMeta {
        openMeta = true 
        m.start = m.current
      } else {
        m.end = m.current - len(line)
        m.init = m.current
        m.meta = true
      }
    }
  }
}


func (m *mdReader) readLine() ([]byte, bool, error)  {
  line, err := m.reader.ReadBytes('\n')

  isEOF := err == io.EOF
  if err != nil && !isEOF {
    return nil, false, err
  }

  // save byte position
  m.current += len(line)
  _, err = m.output.Write(line)
  return line, isEOF, err
}
