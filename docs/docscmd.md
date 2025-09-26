# Docs Command

Goal: Manage my plain text files located at my personal directory.

- Read and load all the files of the provided directory.
- Index files to perform faster searches using a sqlite database.

## Add

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

