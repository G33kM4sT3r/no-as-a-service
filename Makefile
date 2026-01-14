# Build configuration
BINARY_NAME=NoAAS
BUILD_DIR=dist

# Version info (can be overridden)
VERSION?=1.0.0
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags
LDFLAGS=-ldflags "-s -w -X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildTime=$(BUILD_TIME)"

.PHONY: all build clean run test deps lint help \
        build-linux build-windows build-darwin-arm64 build-darwin-amd64 build-all

# Default target
all: build

# Build for current platform
build:
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .

# Run the application
run:
	go run .

# Run tests
test:
	go test -v ./...

# Install dependencies
deps:
	go mod download
	go mod tidy

# Lint the code
lint:
	golangci-lint run ./...

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)
	go clean

# Build for Linux (amd64)
build-linux:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .

# Build for Windows (amd64)
build-windows:
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .

# Build for macOS ARM64 (Apple Silicon)
build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .

# Build for macOS Intel
build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .

# Build for all platforms
build-all: clean
	@mkdir -p $(BUILD_DIR)
	@echo "Building for Linux (amd64)..."
	@$(MAKE) build-linux
	@echo "Building for Windows (amd64)..."
	@$(MAKE) build-windows
	@echo "Building for macOS ARM64..."
	@$(MAKE) build-darwin-arm64
	@echo "Building for macOS Intel..."
	@$(MAKE) build-darwin-amd64
	@echo "Build complete! Binaries are in $(BUILD_DIR)/"
	@ls -la $(BUILD_DIR)/

# Help
help:
	@echo "Available targets:"
	@echo "  build              - Build for current platform"
	@echo "  build-all          - Build for all platforms"
	@echo "  build-linux        - Build for Linux (amd64)"
	@echo "  build-windows      - Build for Windows (amd64)"
	@echo "  build-darwin-arm64 - Build for macOS ARM64 (Apple Silicon)"
	@echo "  build-darwin-amd64 - Build for macOS Intel"
	@echo "  run                - Run the application"
	@echo "  test               - Run tests"
	@echo "  deps               - Download and tidy dependencies"
	@echo "  lint               - Run linter"
	@echo "  clean              - Remove build artifacts"
	@echo "  help               - Show this help"
