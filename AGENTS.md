# AGENTS.md - Development Guide for AI Coding Agents

This guide provides essential information for AI coding agents working on the PKMS (Personal Knowledge Management System) codebase.

## Project Overview

- **Language**: Go 1.23.2
- **Module**: `github.com/vninomtz/pkms`
- **Architecture**: Standard Go project layout with `cmd/` for CLI commands, `internal/` for business logic
- **Type**: Monorepo containing multiple tools (pkms, todo, pomo)

## Build, Test, and Run Commands

### Main CLI Build
```bash
# Build main pkm binary
go build -o pkm

# Build to specific location
go build -o ./bin/pkm

# Install to /usr/local/bin
./install.sh

# Run without installing
go run main.go <command> [options]

# Common development workflow
./bin/pkm install      # First time setup
./bin/pkm index        # Index notes after changes
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run specific package tests
go test ./internal/notes/
go test ./internal/index/

# Run single test function
go test -v -run TestExtractMetadata ./internal/notes/
go test -v -run TestTokenize ./internal/index/

# Run tests in specific file (use package path)
go test -v ./internal/security/
```

### Development Server (with hot reload)
```bash
# Start server with auto-reload (uses .air.toml)
air

# Server runs at ./tmp/server and watches .go, .html, .tmpl files
```

### Todo CLI (subproject)
```bash
cd todo/cmd/todo
make build  # Builds to ./bin/todo
make test   # Runs tests with verbose output
```

## Code Style Guidelines

### Import Organization

Use three groups separated by blank lines:
```go
import (
    // 1. Standard library (alphabetically)
    "context"
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    
    // 2. Third-party packages (alphabetically)
    "github.com/adrg/frontmatter"
    "github.com/google/uuid"
    
    // 3. Internal packages (alphabetically)
    "github.com/vninomtz/pkms/internal/config"
    "github.com/vninomtz/pkms/internal/notes"
)
```

**Special cases:**
- Blank imports for drivers: `_ "github.com/mattn/go-sqlite3"`
- Named imports for conflicts: `rand2 "math/rand"`

### Naming Conventions

**Variables:**
- `camelCase` for local variables: `cfg`, `srv`, `notesDir`
- Short names in small scopes: `i`, `n`, `t`, `err`
- Common abbreviations: `res` (result), `cfg` (config), `srv` (service)

**Constants:**
- `UPPER_SNAKE_CASE` for package-level constants:
  ```go
  const PKM_VERSION = "0.1.0"
  const DB_FILENAME = "pkms.db"
  ```
- `PascalCase` for exported enum-like constants:
  ```go
  const (
      CategoryPomodoro   = "Pomodoro"
      CategoryShortBreak = "ShortBreak"
  )
  ```

**Functions:**
- `PascalCase` for exported: `AddCommand`, `NewIndexer`
- `camelCase` for private: `newTimeId`, `indexFile`

**Types:**
- `PascalCase` for exported: `Note`, `Entry`, `Store`
- `camelCase` for private: `config`, `parser`

**Files:**
- `snake_case`: `note_service.go`, `html_parser.go`
- Test files: `*_test.go`

**Method Receivers:**
- Single letter or short abbreviation: `(n *noteService)`, `(s *Store)`

### Error Handling

**Error Declaration:**
```go
var (
    ErrNoIntervals        = errors.New("No intervals")
    ErrIntervalNotRunning = errors.New("Interval not running")
    ErrInvalidState       = errors.New("Invalid State")
)
```

**Error Wrapping (always use %w):**
```go
return fmt.Errorf("error trying to open DB: %w", err)
return nil, fmt.Errorf("error reading file %s: %w", filename, err)
```

**Error Checking Patterns:**
```go
// 1. Early return pattern (preferred)
if err != nil {
    return "", err
}

// 2. Log and continue (for non-critical errors)
if err != nil {
    log.Printf("Error to index %s: %w\n", filename, err)
    continue
}

// 3. Fatal for CLI commands
if err != nil {
    log.Fatal(err)
}

// 4. Specific error checking
if errors.Is(err, os.ErrNotExist) {
    return nil
}
```

### Function and Method Design

**Command Pattern (cmd/ directory):**
```go
func CommandName(args []string) {
    fs := flag.NewFlagSet("commandname", flag.ExitOnError)
    flagVar := fs.String("flag", "default", "description")
    fs.Parse(args)
    
    cfg := config.New()
    cfg.Load()
    
    srv := notes.New(cfg.NotesDir)
    // Implementation
}
```

**Service Pattern (internal/ packages):**
```go
type ServiceName interface {
    Method1() error
    Method2(param string) (Result, error)
}

type serviceName struct {
    field string
}

func New(param string) ServiceName {
    return &serviceName{field: param}
}
```

**Factory Functions:**
```go
func New(path string) (*Store, error) {
    // Initialization
    return &Store{db: db}, nil
}
```

### Struct Design

**Use clear, descriptive fields:**
```go
type Note struct {
    Title   string
    Content string
    Type    string
    Created time.Time
    Updated time.Time
    Public  bool
    Tags    []string
    Links   []string
}
```

### Concurrency

**Use WaitGroup and Mutex for parallel operations:**
```go
var wg sync.WaitGroup
var mu sync.Mutex
results := []*Result{}

for _, item := range items {
    wg.Add(1)
    go func(i Item) {
        defer wg.Done()
        result := process(i)
        
        mu.Lock()
        results = append(results, result)
        mu.Unlock()
    }(item)
}

wg.Wait()
```

**Use context for cancellation:**
```go
func Process(ctx context.Context) error {
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            // Continue processing
        }
    }
}
```

### Resource Management

**Always use defer for cleanup:**
```go
file, err := os.CreateTemp(os.TempDir(), "pkms")
if err != nil {
    return err
}
defer os.Remove(file.Name())
defer file.Close()
```

## Testing Guidelines

**Test function structure:**
```go
func TestFunctionName(t *testing.T) {
    // Arrange
    input := "test data"
    expected := "expected result"
    
    // Act
    result, err := FunctionUnderTest(input)
    
    // Assert
    if err != nil {
        t.Errorf("unexpected error: %s", err)
    }
    if result != expected {
        t.Errorf("expected %s, got %s", expected, result)
    }
}
```

**Use testdata/ directory for fixtures:**
- Place test files in `testdata/` directories
- Use `os.CreateTemp()` for temporary test files

## Dependencies

**Current key dependencies:**
- `github.com/adrg/frontmatter` - YAML frontmatter parsing
- `github.com/google/uuid` - UUID generation
- `github.com/mattn/go-sqlite3` - SQLite database
- `github.com/yuin/goldmark` - Markdown parsing
- `golang.org/x/net/html` - HTML parsing

**Preference:** Use standard library where possible before adding dependencies

## File Organization

- **Commands**: `cmd/*.go` - CLI command implementations
- **Business logic**: `internal/` - Private application packages
- **Tests**: `*_test.go` - Colocated with source files
- **Templates**: `templates/` - HTML templates
- **Documentation**: `docs/` - Project documentation
- **Test data**: `testdata/` - Test fixtures

## Common Patterns to Follow

1. Use interfaces for service boundaries
2. Return errors, don't panic (except in main/init)
3. Prefer explicit over implicit
4. Keep functions small and focused
5. Use `flag.FlagSet` for command-line parsing
6. Configuration via `internal/config` package
7. Database operations via `internal/store` package

## Additional Commands (Not Yet Integrated)

### share (WIP - not in main.go router)
Share encrypted notes via external service.

**Implementation exists** in `cmd/share.go` but not wired to main command router.

**Options:**
- `--path <directory>`: Path to documents directory
- `--name <filename>`: Name of the note to share

**Environment Variables:**
- `PKMS_SHARE_URL`: Base URL for sharing service

**Usage (when integrated):**
```bash
export PKMS_SHARE_URL="https://share.example.com"
pkm share --name "mynote.md" --path "/path/to/notes"
```

**Features:**
- Encrypts note content with random key
- Uploads to external service
- Returns shareable URL with decryption key as fragment

## Historical Command Patterns

The project previously used a `docs` subcommand pattern:
```bash
pkm docs -add        # Old pattern (now: pkm add)
pkm docs -index      # Old pattern (now: pkm index)
pkm docs -get <name> # Old pattern (not currently implemented)
pkm docs -public     # Old pattern (now: pkm search --public)
```

Some older commands not currently implemented:
- `ls` / `docs -ls`: List all notes
- `get <id>`: Get note by ID or filename
- `find -n <name>`: Find notes by name
- `find -t <tag>`: Find notes by tag
