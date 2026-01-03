# Release recipes
binary_name := "lazysql"
main_path := "."
bin_dir := "bin"
dist_dir := "dist"

# Show available release commands
[private]
default:
    @just --list --unsorted --justfile {{source_file()}}

# Build for macOS Intel
[private, no-cd]
build-darwin:
    #!/usr/bin/env bash
    set -euo pipefail
    mkdir -p {{bin_dir}}
    GOOS=darwin GOARCH=amd64 go build -ldflags '-s -w' -o {{bin_dir}}/{{binary_name}}-darwin-amd64 {{main_path}}

# Build for macOS Apple Silicon
[private, no-cd]
build-darwin-arm64:
    #!/usr/bin/env bash
    set -euo pipefail
    mkdir -p {{bin_dir}}
    GOOS=darwin GOARCH=arm64 go build -ldflags '-s -w' -o {{bin_dir}}/{{binary_name}}-darwin-arm64 {{main_path}}

# Build for Linux AMD64
[private, no-cd]
build-linux:
    #!/usr/bin/env bash
    set -euo pipefail
    mkdir -p {{bin_dir}}
    GOOS=linux GOARCH=amd64 go build -ldflags '-s -w' -o {{bin_dir}}/{{binary_name}}-linux-amd64 {{main_path}}

# Build for Windows AMD64
[private, no-cd]
build-windows:
    #!/usr/bin/env bash
    set -euo pipefail
    mkdir -p {{bin_dir}}
    GOOS=windows GOARCH=amd64 go build -ldflags '-s -w' -o {{bin_dir}}/{{binary_name}}-windows-amd64.exe {{main_path}}

# Create macOS release packages
[no-cd]
darwin: build-darwin build-darwin-arm64
    #!/usr/bin/env bash
    set -euo pipefail
    mkdir -p {{dist_dir}}
    tar -czf {{dist_dir}}/{{binary_name}}-darwin-amd64.tar.gz -C {{bin_dir}} {{binary_name}}-darwin-amd64
    tar -czf {{dist_dir}}/{{binary_name}}-darwin-arm64.tar.gz -C {{bin_dir}} {{binary_name}}-darwin-arm64

# Create Linux release package
[no-cd]
linux: build-linux
    #!/usr/bin/env bash
    set -euo pipefail
    mkdir -p {{dist_dir}}
    tar -czf {{dist_dir}}/{{binary_name}}-linux-amd64.tar.gz -C {{bin_dir}} {{binary_name}}-linux-amd64

# Create Windows release package
[no-cd]
windows: build-windows
    #!/usr/bin/env bash
    set -euo pipefail
    mkdir -p {{dist_dir}}
    cd {{bin_dir}} && zip -q ../{{dist_dir}}/{{binary_name}}-windows-amd64.zip {{binary_name}}-windows-amd64.exe

# Create all release packages
[no-cd]
release: darwin linux windows
    #!/usr/bin/env bash
    set -euo pipefail
    echo "All release packages created in {{dist_dir}}"
    ls -la {{dist_dir}}
