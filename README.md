# Personal Knowledge Management Tools

## Goals

Develop tools for a personal knowledge management such as:

- Search, manipulation and aggregation of information
- Information visualization, classification and publication
- Use of plain text files with Markdown syntax
- Information backup and conversion to different types of storage (SQLite, CSV, etc)

## Installation

```bash
# Build the project
go build -o pkm

# Install PKMS in your home directory
./pkm install
```

## Available Commands

### `add`
Creates a new note by opening vim editor.

```bash
pkm add
```

The note is saved with a unique identifier in the format *YYYYMMDDHHmmss*.

### `search`
Search for notes in your collection.

**Options:**
- `--filename <name>`: Search note by filename
- `--public`: List all public notes

```bash
# Search by filename
pkm search --filename "myfile.md"

# List public notes
pkm search --public
```

### `inspect`
Inspect URLs and fetch their content. Exports results to `pages.json`.

**Options:**
- `--url <urls>`: Comma-separated list of URLs to inspect

```bash
pkm inspect --url "https://example.com,https://another.com"
```

### `index`
Index all notes into SQLite database for faster searches.

```bash
pkm index
```

### `publish`
Copy all public notes to a specified directory.

**Options:**
- `-o <directory>`: Output directory where notes will be copied

```bash
pkm publish -o /path/to/output
```

### `install`
Install PKMS configuration in your home directory.

```bash
pkm install
```

### `version`
Display the current version of PKM.

```bash
pkm version
```

## Configuration

PKMS uses a configuration file stored in your home directory after running `pkm install`. The configuration includes:
- Notes directory path
- SQLite database file location

## TODO

- [ ] Implement notes subcommand
- [ ] Add more search filters
- [ ] Improve markdown preview options
