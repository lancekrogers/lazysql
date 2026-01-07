# Keybinding Coverage Report

## Scope

- Leader key sequences (namespaces and top-level bindings).
- Vim mode keybindings (Normal and Insert handlers).
- Which-key hint descriptions and conflict checks.

## Automated verification

- `internal/vim/leader/keybinding_coverage_test.go`: validates leader registry coverage, which-key descriptions, group hints, and exact binding set.
- `internal/vim/namespace/*_test.go`: validates namespace handler execution paths and error handling.
- `components/home_leader_test.go`: validates top-level `\e` registration.
- `internal/vim/modes/normal_handler_test.go` and `internal/vim/modes/insert_handler_test.go`: validates vim mode keybindings.

## Leader key sequences

### Top-level

- `\e` Toggle tree

### Buffer namespace (`\b`)

- `\bn` Next buffer
- `\bp` Previous buffer
- `\bb` Alternate buffer
- `\bd` Delete buffer
- `\bD` Force delete buffer
- `\bl` List buffers
- `\bN` New buffer
- `\b1`..`\b9` Jump to buffer by index

### Shell namespace (`\s`)

- `\ss` Toggle shell pane
- `\sf` Focus shell input
- `\sq` Close shell pane
- `\sc` Clear shell output
- `\sr` Resize shell pane
- `\sh` Shell history

### Run namespace (`\r`)

- `\rr` Run current query
- `\rl` Run selected lines
- `\ra` Run all queries
- `\rp` Run with parameters
- `\rh` Query history

### Describe namespace (`\d`)

- `\dt` Describe table
- `\dv` Describe view
- `\di` Describe indexes
- `\ds` Describe sequences
- `\df` Describe functions
- `\dd` Describe database
- `\d<Enter>` Describe under cursor

### Find namespace (`\f`)

- `\ft` Find table
- `\fc` Find column
- `\ff` Find function
- `\fv` Find view
- `\fs` Find schema
- `\fq` Find in queries
- `\fr` Find recent

### Tree namespace (`\t`)

- `\tt` Toggle tree
- `\te` Expand node
- `\tc` Collapse node
- `\tE` Expand all
- `\tC` Collapse all
- `\tf` Focus tree
- `\tr` Refresh tree
- `\ts` Sync with buffer

### Workspace namespace (`\w`)

- `\wc` Connect
- `\wd` Disconnect
- `\wl` List connections
- `\ws` Switch connection

### Explain namespace (`\x`)

- `\xx` Explain current query
- `\xX` Explain analyze

### Cancel namespace (`\c`)

- `\cc` Cancel running query
- `\cl` Clear results

## Vim mode keybindings

### Normal mode

- Movement: `h`, `j`, `k`, `l`, `w`, `b`, `gg`, `G`
- Insert transitions: `i`, `I`, `a`, `A`, `o`, `O`
- Mode actions: `\` (leader), `:` (command line)

### Insert mode

- Exit: `Esc`, `Ctrl+C`
- Editing: `Enter`, `Backspace`, `Delete`, `Tab`, rune insert

## Results

- Leader bindings: all sequences registered with expected descriptions and which-key hints.
- Vim mode bindings: all key paths covered by unit tests.
- Conflicts: none detected; exact binding set matches expected coverage.
- Manual UI validation: pending in the next festival validation task.
