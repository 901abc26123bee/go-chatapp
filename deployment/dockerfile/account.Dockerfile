
# syntax = docker/dockerfile:1
FROM golang:1.22.4-alpine3.19 as compiler

ENV GO111MODULE=on

# cache module
WORKDIR /gopath/src/
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg go mod download

# build release
COPY . ./
RUN --mount=type=cache,target=/go/pkg \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o ./account ./cmd/account/server.go

FROM alpine3.19
COPY --from=compiler /gopath/src/account /go/bin/account
