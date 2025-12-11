.PHONY: build build-all clean test lint install

# Build variables
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -X 'github.com/ethantroy/oscal-cli/internal/cli.Version=$(VERSION)' \
           -X 'github.com/ethantroy/oscal-cli/internal/cli.GitCommit=$(COMMIT)' \
           -X 'github.com/ethantroy/oscal-cli/internal/cli.BuildDate=$(DATE)'

# Binary name
BINARY := oscal

# Build for current platform
build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) ./cmd/oscal

# Build for all platforms
build-all: clean
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY)-darwin-amd64 ./cmd/oscal
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY)-darwin-arm64 ./cmd/oscal
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY)-linux-amd64 ./cmd/oscal
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY)-linux-arm64 ./cmd/oscal
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY)-windows-amd64.exe ./cmd/oscal

# Clean build artifacts
clean:
	rm -f $(BINARY)
	rm -rf dist/

# Run tests
test:
	go test -v ./...

# Run linter
lint:
	golangci-lint run

# Install to GOPATH/bin
install:
	go install -ldflags "$(LDFLAGS)" ./cmd/oscal

# Run the binary
run: build
	./$(BINARY)
