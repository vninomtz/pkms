# Documentation


## Collector

- Provide an interface to get data from files or a database.
- Support multi nodes but for the moment will only support *FileNode*

```golang
type collector interface{
    Collect() ([]FileNode, error)
}
```

## Searcher

- Search, filter, and sort FileNodes
- In the future I will implement a inner Tree data structure but for now only a simple List is enough

```golang
type Searcher interface {
    File(filename string) (FileNode, error)
}
```

## Features

- search by name
- search by tag
- open note by name
- create a note in line
- create note using vim as editor


## CLI Commands

Find all notes with the tag or tags
```
 cmd -find -tags "tag1,tag2"
```

## Concerns

Options for an easier way to get a note.

Problem: sometimes the note has a long name and to get the note the full name should
be used in the CLI.

- Create a hash for the name and use it to get the note (like docker or git)

Options to open a note.

- Output in stdout
- Create tmp pdf file
- Use a markdown preview tool

## TODO

- Improve the cli with [cobra](https://github.com/spf13/cobra) and [viper](https://github.com/spf13/viper)
- Create a new method to display the Notes in console with better format
