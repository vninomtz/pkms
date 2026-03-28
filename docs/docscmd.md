# PKM Commands Documentation

Goal: Manage plain text files located in your personal directory.

- Read and load all files from the configured directory
- Index files to perform faster searches using a SQLite database
- Search, publish, and inspect content

## Commands Reference

### add

Creates a new note by opening vim editor.

**Usage:**
```bash
pkm add
```

**Behavior:**
- Opens vim editor for content input
- Saves content on exit (`:q`)
- Generates unique identifier in format *YYYYMMDDHHmmss*
- Stores note in configured notes directory

**Implementation:** `cmd/add.go:14`

---

### search

Search for notes in your collection.

**Usage:**
```bash
pkm search [options]
```

**Options:**
- `--filename <name>`: Search note by exact filename match
- `--public`: List all notes marked as public

**Examples:**
```bash
# Find specific note by filename
pkm search --filename "20240315120000.md"

# List all public notes
pkm search --public
```

**Implementation:** `cmd/search.go:12`

---

### inspect

Inspect URLs and fetch their content. Useful for bookmarking and content archival.

**Usage:**
```bash
pkm inspect --url <urls>
```

**Options:**
- `--url <urls>`: Comma-separated list of URLs to inspect

**Behavior:**
- Fetches content from each URL
- Displays status code and HTML length
- Exports all results to `pages.json` file

**Example:**
```bash
pkm inspect --url "https://example.com,https://github.com"
```

**Implementation:** `cmd/inspect.go:14`

---

### index

Index all notes into SQLite database for faster searches and queries.

**Usage:**
```bash
pkm index
```

**Behavior:**
- Scans all notes in configured directory
- Parses note metadata and content
- Stores in SQLite database
- Reports number of successfully indexed documents

**Implementation:** `cmd/index.go:13`

---

### publish

Copy all public notes to a specified output directory.

**Usage:**
```bash
pkm publish -o <directory>
```

**Options:**
- `-o <directory>`: Required. Output directory path where notes will be copied

**Example:**
```bash
pkm publish -o /path/to/blog/content
```

**Behavior:**
- Searches for all notes marked as public
- Copies each public note to output directory
- Reports number of successfully copied notes

**Implementation:** `cmd/publish.go:13`

---

### install

Initialize PKMS in your home directory.

**Usage:**
```bash
pkm install
```

**Behavior:**
- Creates configuration directory at `$HOME/.pkms`
- Sets up default configuration
- Creates necessary subdirectories

**Implementation:** `cmd/install.go:11`

---

### version

Display the current version of PKM.

**Usage:**
```bash
pkm version
```

**Output:** Current PKM version number

**Implementation:** `main.go:35`

---

## Configuration

After running `pkm install`, configuration is stored in your home directory. The config includes:
- **NotesDir**: Directory where notes are stored
- **SQLiteFile**: Path to SQLite database for indexed notes

## Future Considerations

### Easier Note Retrieval

**Problem:** Long note names require typing the full filename.

**Possible solutions:**
- Create a hash for note names (similar to Docker/Git short hashes)
- Allow fuzzy search by partial filename
- Add note aliases or tags

### Note Display Options

**Current:** Output to stdout only

**Possible improvements:**
- Generate temporary PDF file
- Use markdown preview tool
- Open in default markdown viewer
- Web-based preview interface

