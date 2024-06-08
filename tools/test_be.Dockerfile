# Start from the official Go base image
FROM golang:1.22.4-alpine3.19

ENV GO111MODULE=on

# cache module
WORKDIR /gopath/src/
COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg go mod download

# Install golangci-lint(v1.59.0)
RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.59.0 && \
    mv ./bin/golangci-lint /usr/local/bin/

# Install golint and goimports
RUN go install golang.org/x/lint/golint@v0.0.0-20210508222113-6edffad5e616 && \
    go install golang.org/x/tools/cmd/goimports@v0.22.0

# Copy the source from the current directory to the Working Directory inside the container
WORKDIR /gopath/src/
COPY . ./

CMD ["/gopath/src/tools/test_be.sh /gopath/src"]
