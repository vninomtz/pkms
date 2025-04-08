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

### Settings

- PKMS_STORE_TYPE: Enviroment variable to define the type of Store. 0 for File System, 1 for SQLite.
- PKMS_STORE_PATH: Enviroment variable to define the directory or file to use.

### Add

Create a new note

Flags:
- c: Include the content inline. By the fault open a vim terminal.
- title: Define the name of the note. By the fault use a unique identifier in the format *YYYYMMDDHHmmss*

```bash
 cmd add -c "Test content" -title "custom-title"
```

### Ls

List all the notes

Flags:
- t: List and count all the tags in the notes

```bash
 cmd ls
```

### Get

Get a note by name

```bash
 cmd get "note-title"
```

### Find

Search a note by title or by tags

Flags:
- n: Search by note title
- t: Search by tags. Allow multiple separeted by comma


```bash
 cmd find -t "tag,tag1"
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
