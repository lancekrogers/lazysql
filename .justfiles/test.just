# Test recipes

# Show available test commands
[private]
default:
    @just --list --unsorted --justfile {{source_file()}}

# Run all tests
[no-cd]
all:
    go test ./...

# Run unit tests with short mode
[no-cd]
unit:
    go test -short ./...

# Run tests with verbose output
[no-cd]
verbose:
    go test -v ./...

# Run tests with race detector
[no-cd]
race:
    go test -race ./...

# Run tests with coverage
[no-cd]
coverage:
    go test -cover ./...

# Generate HTML coverage report
[no-cd]
coverage-html:
    go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html
