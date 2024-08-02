# syntax = docker/dockerfile:1
FROM golang:1.22.4-alpine3.19 as compiler

ENV GO111MODULE=on

# Set the working directory inside the container
WORKDIR /gopath/src/

# Copy go.mod and go.sum files to the working directory
# COPY ./go.mod ./go.sum ./
COPY go.mod go.sum ./

# Download the dependencies
RUN --mount=type=cache,target=/go/pkg go mod download

# Copy the entire project directory to the working directory
# COPY ./ ./
COPY . ./

# Build the Go binary
RUN --mount=type=cache,target=/go/pkg \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o ./realtime ./cmd/realtime/server.go

# Use a minimal Alpine image to run the Go binary
FROM alpine:3.19

# Copy the binary from the build stage
COPY --from=compiler /gopath/src/realtime /go/bin/realtime
