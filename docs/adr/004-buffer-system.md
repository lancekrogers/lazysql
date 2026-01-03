# ADR 004: Buffer System

## Status

Proposed

## Context

Lazysql uses a `TabbedPane` component to manage multiple result tabs, but there
is no formal buffer abstraction for SQL queries, file-backed content, or dirty
state. To support Vim-like workflows and the `\b` namespace, we need a buffer
manager that tracks query buffers, persists file paths, and integrates with the
existing tab UI.

## Decision Drivers

- Enable multiple concurrent SQL query buffers with predictable navigation.
- Track dirty state and prevent data loss.
- Reuse existing `TabbedPane` UI patterns where possible.
- Provide a clean interface for command-line commands (`:w`, `:e`, `:q`).

## Decision

Introduce a buffer manager that owns buffer lifecycle and exposes navigation
operations. Each buffer is rendered as a tab in the editor area using the
existing `TabbedPane` infrastructure.

### Buffer Structure

```go
type BufferID uint64

type CursorPos struct {
	Row int
	Col int
}

type Buffer struct {
	ID        BufferID
	Name      string     // Display name (filename or "Untitled-N")
	Path      string     // Absolute file path (empty if unsaved)
	Content   string
	Dirty     bool
	Cursor    CursorPos
	LastRunAt time.Time
	CreatedAt time.Time
}
```

### Buffer Manager Interface (Go)

```go
type BufferManager interface {
	Create(ctx context.Context) (*Buffer, error)
	Open(ctx context.Context, path string) (*Buffer, error)
	Current() *Buffer
	Switch(ctx context.Context, id BufferID) (*Buffer, error)
	Next(ctx context.Context) *Buffer
	Previous(ctx context.Context) *Buffer
	Alternate(ctx context.Context) *Buffer
	Close(ctx context.Context, id BufferID) error
	CloseForce(ctx context.Context, id BufferID) error
	List() []*Buffer
	MarkDirty(id BufferID, dirty bool)
	UpdateContent(id BufferID, content string, cursor CursorPos)
}
```

### UI Integration

- Each buffer is backed by a tab (existing `TabbedPane` headers).
- Tab title: buffer name + `*` when dirty.
- The active buffer corresponds to the active tab.
- Tab overflow is handled by `TabbedPane` layout; future enhancement can add
  scroll or dropdown.

### Navigation Commands (`\b` namespace)

- `\bn` Next buffer
- `\bp` Previous buffer
- `\bb` Alternate buffer (last visited)
- `\bd` Delete/close buffer (respect dirty state)
- `\bl` List buffers (picker)

### Lifecycle Rules

- **Create**: New buffer with `Untitled-N`, empty content.
- **Open**: Load file content into a new buffer; path is stored.
- **Modify**: Any content change marks buffer dirty.
- **Save**: Writes content to path; clears dirty flag.
- **Close**: If dirty, return error unless forced.
- **Alternate**: Track last active buffer for quick switching.

## Consequences

- Adds a consistent buffer abstraction for Vim-like workflows.
- Requires minimal updates to `TabbedPane` to reflect dirty state.
- Enables `:w/:e/:q` commands to operate on buffers via a stable API.
