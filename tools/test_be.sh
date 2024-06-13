#!/bin/bash
set -e

# Set workdir to the first argument, or use the current working directory if not provided
workdir=${1:-$(pwd)}

# Change to the specified workdir
cd "$workdir" || { echo "Failed to change directory to $workdir"; exit 1; }

# TODO: refine
# Run gofmt
run_go_fmt() {
    echo "Running gofmt..."
    if ! gofmt -w .; then
        echo "gofmt failed"
        exit 1
    fi
}

# Run go vet
run_go_vet() {
    echo "Running go vet..."
    if ! go vet ./...; then
        echo "go vet failed"
        exit 1
    fi
}

# Run go mod tidy
run_go_mod_tidy() {
    echo "Running go mod tidy..."
    if ! go mod tidy; then
        echo "go mod tidy failed"
        exit 1
    fi
}

# Run go imports
run_go_imports() {
    echo "Running go imports..."
    if ! goimports -w .; then
        echo "goimports failed"
        exit 1
    fi
}

# Run go lint
run_go_lint() {
    echo "Running go lint..."
    if ! golangci-lint run --timeout=5m -D errcheck ./...; then
        echo "golangci-lint failed"
        exit 1
    fi
}

# Run go test
run_go_test() {
    echo "Running go test..."
    if ! go test ./...; then
        echo "go test failed"
        exit 1
    fi
}

# Execute the functions
run_go_fmt
run_go_vet
run_go_mod_tidy
run_go_imports
run_go_lint
run_go_test

echo "All tasks completed successfully."