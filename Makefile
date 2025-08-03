# Makefile for log_delta

# Variables
BINARY_NAME=log_delta
MAIN_FILE=main.go

# Default target
all: build

# Build the application
build:
	go build -o $(BINARY_NAME) $(MAIN_FILE)

# Run tests
test:
	go test -v

# Run tests with coverage
test-coverage:
	go test -cover

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)

# Run the application (requires LOG_FILE environment variable)
run:
	./$(BINARY_NAME) $(LOG_FILE)

# Install dependencies (if any are added later)
deps:
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Lint code (requires golint)
lint:
	golint ./...

# Show help
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  clean         - Remove build artifacts"
	@echo "  run           - Run the application (set LOG_FILE=path/to/file)"
	@echo "  deps          - Install dependencies"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code"
	@echo "  help          - Show this help"

.PHONY: all build test test-coverage clean run deps fmt lint help
