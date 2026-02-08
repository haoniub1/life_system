.PHONY: build run clean test help

# Default target
help:
	@echo "Life System Backend - Available targets:"
	@echo "  make build      - Build the binary"
	@echo "  make run        - Build and run the application"
	@echo "  make clean      - Remove build artifacts and database"
	@echo "  make test       - Run tests"
	@echo "  make deps       - Download dependencies"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

# Build the application
build: deps
	@echo "Building..."
	@mkdir -p data
	@go build -o bin/life-system-backend main.go
	@echo "Build complete: bin/life-system-backend"

# Run the application
run: build
	@echo "Starting server..."
	@./bin/life-system-backend -f etc/config.yaml

# Clean build artifacts and database
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f data/life-system.db
	@echo "Clean complete"

# Run tests
test:
	@echo "Running tests..."
	@go test ./... -v

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run ./...
