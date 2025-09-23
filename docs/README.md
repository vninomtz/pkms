# Documentation


## Use cases

- Get public notes and move them to the bycuriosity project
- Create a new note from the terminal
- Get note by name
- Parse note to HTML
- Share note with custom link


## Docs Command

Manage text files in a directory

### Add
Creates a new note.

- Opens vim editor and save the content on close vim using `:q`.
- By default use a unique identifier in the format *YYYYMMDDHHmmss*

```bash
pkm docs -add
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
