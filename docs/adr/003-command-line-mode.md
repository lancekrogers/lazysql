# ADR 003: Command Line Mode

## Status

Proposed

## Context

Lazysql does not currently provide a Vim-style `:` command line. Users perform
file operations and buffer control through direct keybindings or modal actions.
To match Vim/Neovim workflows, we need a command-line input that appears on `:`
in Normal mode, parses a small set of file commands, and executes them against
the buffer system with clear error feedback.

## Decision Drivers

- Provide familiar Vim-style `:` commands without breaking existing UI layout.
- Keep parsing deterministic and easy to extend.
- Ensure file operations are safe and report meaningful errors.
- Allow command execution to be tested independently of UI components.

## Decision

Introduce a command-line component (bottom bar) that accepts input prefixed by
`:` and uses a simple parser to map input to commands. Execution is delegated
to a command executor that interacts with the buffer system and filesystem.

### Command Input UI

- Appears at the bottom of the screen when `:` is pressed in Normal mode.
- Displays `:` as a fixed prefix; user edits the remainder.
- Supports basic editing: backspace, left/right movement, delete, home/end.
- `Enter` executes the parsed command; `Esc` cancels and hides the input.
- Command history is optional (future enhancement).

### Grammar

```
command      := ":" cmd [ws args]
cmd          := "w" | "w!" | "e" | "q" | "q!" | "wq" | "wq!"
args         := filepath
filepath     := token | quoted
token        := non-space characters
quoted       := '"' (any char except '"')* '"'
ws           := 1+ spaces
```

### Command Semantics

- `:w [path]`   Save current buffer. If `path` is omitted, use buffer path.
               Error if no buffer path exists.
- `:w! [path]`  Save and overwrite existing file.
- `:e <path>`   Open file into a new buffer (or reuse if already open).
- `:q`          Close current buffer. If dirty, warn and block.
- `:q!`         Force close current buffer, discard changes.
- `:wq [path]`  Save (path optional), then close buffer.
- `:wq! [path]` Force save (overwrite), then close buffer.

### Path Handling

- Relative paths are resolved against the process working directory.
- `~` is expanded to the user home directory.
- Paths are normalized with `filepath.Clean`.
- Use `filepath.IsAbs` and `filepath.Abs` to generate canonical paths.

### Error Handling

- `:w` with no path and no buffer path → "No file name".
- `:w` to existing file without `!` → "File exists (use :w!)".
- `:e` on missing file → "File not found".
- `:q` on dirty buffer → "No write since last change".
- Unknown command → "Not an editor command".

Errors are surfaced in a non-blocking status line or modal that does not
destroy command-line input history.

### Interfaces (Go)

```go
package cmdline

import (
	"context"
)

type CommandType uint8

const (
	CommandWrite CommandType = iota
	CommandEdit
	CommandQuit
	CommandWriteQuit
)

type Command struct {
	Type   CommandType
	Force  bool
	Path   string
	Raw    string
}

type Parser interface {
	Parse(input string) (Command, error)
}

type Executor interface {
	Execute(ctx context.Context, cmd Command) error
}

type Input interface {
	Show()
	Hide()
	SetText(text string)
	SetExecutor(exec Executor)
}
```

### Integration Notes

- The mode router triggers the command line only from Normal mode on `:`.
- Command execution should use `context.Context` and report errors with context.
- Buffer operations are delegated to the buffer manager; filesystem I/O lives
  in a dedicated service layer for testability.

## Consequences

- Adds a single, predictable entry point for file operations.
- Requires a buffer manager that exposes open/save/close behaviors.
- Introduces a small parsing layer that can be extended for future commands.
