#!/bin/bash
set -e

# Set workdir to the first argument, or use the current working directory if not provided
workdir=${1:-$(pwd)}

# Change to the specified workdir
cd "$workdir" || { echo "Failed to change directory to $workdir"; exit 1; }

# Run gofmt
run_go_fmt() {
    echo "Running gofmt..."
    gofmt -w .
}

# Run go vet
run_go_vet() {
    echo "Running go vet..."
    go vet ./...
}

# Run go mod tidy
run_go_mod_tidy() {
    echo "Running go mod tidy..."
    go mod tidy
}

# Run go imports
run_go_imports() {
    echo "Running go imports..."
    goimports -w .
}

# Run go lint
run_go_lint() {
    echo "Running go lint..."
    golangci-lint run --timeout=5m -D errcheck ./...
}

# Run go test
run_go_test() {
    echo "Running go test..."
    go test ./...
}

# Execute the functions
run_go_fmt
run_go_vet
run_go_mod_tidy
run_go_imports
run_go_lint
run_go_test

echo "All tasks completed successfully."
