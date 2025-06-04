# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Project parameters
BINARY_NAME=trafficsim
BINARY_UNIX=$(BINARY_NAME)_unix
MAIN_PATH=./cmd/simrunner

# Build targets
.PHONY: all build clean test test-coverage test-verbose deps fmt vet lint run help

# Default target
all: test build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH)

# Build for Linux
build-linux:
	@echo "Building for Linux..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v $(MAIN_PATH)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests with race detection
test-race:
	@echo "Running tests with race detection..."
	$(GOTEST) -v -race ./...

# Run tests in verbose mode
test-verbose:
	@echo "Running tests in verbose mode..."
	$(GOTEST) -v -count=1 ./...

# Run specific test
test-run:
	@echo "Running specific test (use TEST=TestName)..."
	$(GOTEST) -v -run $(TEST) ./...

# Benchmark tests
benchmark:
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Update dependencies
deps-update:
	@echo "Updating dependencies..."
	$(GOGET) -u ./...
	$(GOMOD) tidy

# Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

# Vet code
vet:
	@echo "Vetting code..."
	$(GOVET) ./...

# Install golangci-lint and run it
lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run

# Run the application
run:
	@echo "Running $(BINARY_NAME)..."
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH) && ./$(BINARY_NAME)

# Run with hot reload (requires air)
dev:
	@echo "Starting development server with hot reload..."
	@which air > /dev/null || (echo "Installing air..." && go install github.com/air-verse/air@latest)
	air

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/air-verse/air@latest
	go install golang.org/x/tools/cmd/goimports@latest

# Generate mocks (if using mockgen)
mocks:
	@echo "Generating mocks..."
	@which mockgen > /dev/null || (echo "Installing mockgen..." && go install go.uber.org/mock/mockgen@latest)
	go generate ./...

# Security check
security:
	@echo "Running security checks..."
	@which gosec > /dev/null || (echo "Installing gosec..." && go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest)
	gosec ./...

# Docker build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME) .

# Docker run
docker-run: docker-build
	@echo "Running Docker container..."
	docker run -p 8080:8080 $(BINARY_NAME)

# Quick quality check
check: fmt vet test

# Full quality check
check-all: fmt vet lint test-race test-coverage security

# Display help
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  build-linux   - Build for Linux"
	@echo "  clean         - Clean build artifacts"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  test-race     - Run tests with race detection"
	@echo "  test-verbose  - Run tests in verbose mode"
	@echo "  test-run      - Run specific test (use TEST=TestName)"
	@echo "  benchmark     - Run benchmark tests"
	@echo "  deps          - Download dependencies"
	@echo "  deps-update   - Update dependencies"
	@echo "  fmt           - Format code"
	@echo "  vet           - Vet code"
	@echo "  lint          - Run linter"
	@echo "  run           - Build and run the application"
	@echo "  dev           - Start development server with hot reload"
	@echo "  install-tools - Install development tools"
	@echo "  mocks         - Generate mocks"
	@echo "  security      - Run security checks"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Build and run Docker container"
	@echo "  check         - Quick quality check (fmt, vet, test)"
	@echo "  check-all     - Full quality check"
	@echo "  help          - Display this help"
