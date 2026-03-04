# AGENTS.md

## Project Overview

This is a small Go command‑line program (`dailynotes`) that creates a daily note in an Obsidian vault. It targets macOS by default and is intended to be run manually or via a cron/launchd job.

## Architecture

- **Entry point**: `/main.go` builds a single binary.
- **Core logic**: resides in `internal/notes/notes.go`; it determines the vault path, ensures a `days` subdirectory exists, and creates an empty markdown file named `YYYY-MM-DD.md` for the current date.
- **Configuration**: a hard‑coded `DefaultVaultPath` constant is used by default; tests override it with `SetVaultPath`.
- **Dependencies**: only the Go standard library; `go.mod` defines the module `dailynotes`.
- There are no external services or libraries (aside from optional lint/test tooling).
- **Binary**: The executable stays in the `/bin` folder
## Project Conventions

- Keep logic in `internal/` so it’s not accidentally imported by other modules.
- Tests live alongside production files (`notes_test.go`) and exercise unexported helpers by staying in the same package.
- Use `t.TempDir()` for filesystem isolation; override global state with `SetVaultPath` and defer resetting it.
- File names use `snake_case` for `.go` files and `PascalCase` for types/functions.
- Add table-driven tests for new scenarios and check both successful paths and errors.
- Avoid adding any external configuration parsing until the need arises; the codebase currently has no CLI flags or environment variable support.

## Build, Lint, and Test Commands

### Build
```bash
go build -o dailynotes .
```

### Run
```bash
go run .
```

### Test
Run all tests:
```bash
go test ./...
```

Run a single test:
```bash
go test -run TestFunctionName ./...
go test -run TestFunctionName -v ./...
```

### Linting
Install golangci-lint if not present:
```bash
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
```

Run linter:
```bash
golangci-lint run ./...
```

### Format
```bash
go fmt ./...
```

### Vet
```bash
go vet ./...
```

### Dependencies
```bash
go mod tidy
go get <package>
```

---

## Code Style Guidelines

### General Principles
- Write clear, readable, idiomatic Go code
- Keep functions small and focused (single responsibility)
- Use meaningful variable and function names
- Comment non-trivial code; document exported functions

### Naming Conventions
- **Variables**: `camelCase` (e.g., `notePath`, `isCreated`)
- **Constants**: `PascalCase` or `ALL_CAPS` for grouped constants (e.g., `DefaultFolder`, `MaxRetries`)
- **Functions**: `PascalCase` for exported, `camelCase` for unexported (e.g., `CreateNote`, `generateFilename`)
- **Files**: `snake_case.go` (e.g., `note_service.go`, `main.go`)
- **Packages**: short, lowercase, no underscores (e.g., `notes`, `utils`)

### Imports
- Use the standard Go import organization:
  ```go
  import (
      "fmt"
      "os"

      "dailynotes/internal/notes"

      "github.com/pkg/errors"
  )
  ```
- Group: standard library → external packages → internal packages
- Use blank imports only when necessary for side effects (`import _ "..."`)

### Formatting
- Run `go fmt` before committing
- Use `gofumpt` for stricter formatting if desired
- Keep lines under 100 characters when reasonable
- No trailing whitespace
- Add final newline to files

### Types
- Use explicit types for function parameters and return values
- Prefer concrete types over interfaces unless polymorphism is needed
- Use `time.Time` for timestamps, not `int` or `string`
- Define custom types for domain concepts (e.g., `type NoteID string`)

### Error Handling
- Always handle errors explicitly; never ignore with `_`
- Return errors with context using `fmt.Errorf("failed to X: %w", err)` or `errors.Wrap()`
- Check for errors immediately after the call that can fail
- Avoid generic error messages; be specific
- Use sentinel errors for known failure cases when appropriate

### Concurrency
- Use goroutines and channels only when needed
- Always handle race conditions with `-race` flag during testing
- Use `sync.WaitGroup` for coordinating multiple goroutines
- Pass context (`context.Context`) for cancellation and timeouts

### Testing
- Test files: `filename_test.go` in the same package
- Test functions: `TestFunctionName(t *testing.T)`
- Use table-driven tests for multiple test cases
- Use `t.Run()` to structure sub-tests
- Assert with `require` or `assert` from `testify` if available
- Test both success and failure paths

### Logging
- Use structured logging (e.g., `log/slog` or `zap`)
- Include relevant context in log messages
- Use appropriate log levels (debug, info, warn, error)

### Configuration
- Use environment variables or config files for settings
- Never hardcode sensitive data (use env vars or secrets management)
- Provide sensible defaults

### File Organization
```
dailynotes/
├── cmd/
│   └── dailynotes/
│       └── main.go
├── internal/
│   └── notes/
│       └── notes.go
├── pkg/
│   └── utils/
├── go.mod
├── go.sum
└── Makefile
```

---

## Cursor / Copilot Rules

No existing rules found in `.cursor/rules/`, `.cursorrules`, or `.github/copilot-instructions.md`.

---

## Recommended VSCode Extensions

- Go (Google)
- Error Lens
- Go Test Explorer
