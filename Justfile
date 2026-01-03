#!/usr/bin/env just --justfile
# lazysql CLI build and development tasks

set dotenv-load := true

# Configuration
binary_name := "lazysql"
bin_dir := "bin"
dist_dir := "dist"
gobin := env_var_or_default("GOBIN", `go env GOPATH` + "/bin")

# Modules
[doc('Cross-platform builds')]
mod xbuild '.justfiles/build.just'

[doc('Testing (unit, coverage)')]
mod test '.justfiles/test.just'

[doc('Release packaging')]
mod release '.justfiles/release.just'

# Show available commands
[private]
default:
    #!/usr/bin/env bash
    echo "lazysql - TUI SQL client"
    echo ""
    just --list --unsorted

# Build lazysql binary (local)
build:
    #!/usr/bin/env bash
    set -euo pipefail
    mkdir -p {{bin_dir}}
    go build -o {{bin_dir}}/{{binary_name}} .

# Run lazysql from source (pass connection args as needed)
run *args:
    go run . {{args}}

# Format Go code
fmt:
    go fmt ./...

# Run go vet
vet:
    go vet ./...

# Run formatting and vetting
lint: fmt vet
    @echo "Linting complete"

# Run all tests
test:
    @just test::all

# Clean build artifacts
clean:
    #!/usr/bin/env bash
    set -euo pipefail
    rm -rf {{bin_dir}} {{dist_dir}} coverage.out coverage.html

# Update and tidy dependencies
deps:
    go get -u ./...
    go mod tidy

# Install lazysql to $GOBIN
install:
    #!/usr/bin/env bash
    set -euo pipefail
    just build
    mkdir -p {{gobin}}
    cp {{bin_dir}}/{{binary_name}} {{gobin}}/{{binary_name}}
    echo "lazysql installed to {{gobin}}/{{binary_name}}"

# Uninstall lazysql from $GOBIN
uninstall:
    #!/usr/bin/env bash
    set -euo pipefail
    if [ -f {{gobin}}/{{binary_name}} ]; then
        rm {{gobin}}/{{binary_name}}
        echo "lazysql uninstalled from {{gobin}}"
    else
        echo "lazysql not found in {{gobin}}"
    fi
