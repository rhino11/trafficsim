# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Node parameters
NPMCMD=npm
JEST=npx jest

# Project parameters
BINARY_NAME=trafficsim
BINARY_UNIX=$(BINARY_NAME)_unix
MAIN_PATH=./cmd/simrunner

# Build targets
.PHONY: all build clean test test-go test-js test-all test-coverage test-coverage-go test-coverage-js test-verbose deps deps-js fmt vet lint run help

# Default target
all: test-all build

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

# Run Go tests
test-go:
	@echo "Running Go tests..."
	$(GOTEST) -v ./...

# Run JavaScript tests
test-js:
	@echo "Running JavaScript tests..."
	@which npm > /dev/null || (echo "Error: npm is not installed" && exit 1)
	$(NPMCMD) test

# Run all tests (Go + JavaScript)
test-all: test-go test-js

# Legacy test target (runs all tests for backward compatibility)
test: test-all

# Run Go tests with coverage
test-coverage-go:
	@echo "Running Go tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Go coverage report generated: coverage.html"

# Run JavaScript tests with coverage
test-coverage-js:
	@echo "Running JavaScript tests with coverage..."
	@which npm > /dev/null || (echo "Error: npm is not installed" && exit 1)
	$(NPMCMD) run test:coverage

# Run all tests with coverage
test-coverage: test-coverage-go test-coverage-js

# Run tests with race detection (Go only)
test-race:
	@echo "Running Go tests with race detection..."
	$(GOTEST) -v -race ./...

# Run tests in verbose mode
test-verbose:
	@echo "Running tests in verbose mode..."
	$(GOTEST) -v -count=1 ./...
	@which npm > /dev/null && $(NPMCMD) run test:debug || echo "Skipping JS verbose tests (npm not available)"

# Run specific test
test-run:
	@echo "Running specific test (use TEST=TestName)..."
	$(GOTEST) -v -run $(TEST) ./...

# Watch JavaScript tests
test-js-watch:
	@echo "Running JavaScript tests in watch mode..."
	@which npm > /dev/null || (echo "Error: npm is not installed" && exit 1)
	$(NPMCMD) run test:watch

# Benchmark tests
benchmark:
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

# Download Go dependencies
deps:
	@echo "Downloading Go dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Install JavaScript dependencies
deps-js:
	@echo "Installing JavaScript dependencies..."
	@which npm > /dev/null || (echo "Error: npm is not installed" && exit 1)
	$(NPMCMD) ci

# Install all dependencies
deps-all: deps deps-js

# Update dependencies
deps-update:
	@echo "Updating Go dependencies..."
	$(GOGET) -u ./...
	$(GOMOD) tidy
	@which npm > /dev/null && (echo "Updating JavaScript dependencies..." && $(NPMCMD) update) || echo "Skipping JS dependency updates (npm not available)"

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
	@echo "Running Go linter..."
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	$(shell go env GOPATH)/bin/golangci-lint run --issues-exit-code=0

# Run JavaScript linting
lint-js:
	@echo "Running JavaScript linting..."
	@which npm > /dev/null || (echo "Error: npm is not installed" && exit 1)
	$(NPMCMD) run lint:js

# Run all linting
lint-all: lint lint-js

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
	@which npm > /dev/null && (echo "Installing JavaScript tools..." && $(NPMCMD) ci) || echo "Skipping JS tools installation (npm not available)"

# Generate mocks (if using mockgen)
mocks:
	@echo "Generating mocks..."
	@which mockgen > /dev/null || (echo "Installing mockgen..." && go install go.uber.org/mock/mockgen@latest)
	go generate ./...

# Security check
security:
	@echo "Running security checks..."
	@which gosec > /dev/null || (echo "Installing gosec..." && go install github.com/securego/gosec/v2/cmd/gosec@latest)
	$(shell go env GOPATH)/bin/gosec ./...

# Security check with detailed output
security-verbose:
	@echo "Running security checks with detailed output..."
	@which gosec > /dev/null || (echo "Installing gosec..." && go install github.com/securego/gosec/v2/cmd/gosec@latest)
	$(shell go env GOPATH)/bin/gosec -fmt=json -out=gosec-report.json ./...
	$(shell go env GOPATH)/bin/gosec -fmt=text ./...

# Docker build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME) .

# Docker run
docker-run: docker-build
	@echo "Running Docker container..."
	docker run -p 8080:8080 $(BINARY_NAME)

# Quick quality check
check: fmt vet test-all

# Full quality check
check-all: fmt vet lint-all test-race test-coverage security

# CI check (for CI environments)
ci: deps-all check-all

# Display help
help:
	@echo "Available targets:"
	@echo "  build            - Build the application"
	@echo "  build-linux      - Build for Linux"
	@echo "  clean            - Clean build artifacts"
	@echo "  test             - Run all tests (Go + JavaScript)"
	@echo "  test-go          - Run Go tests only"
	@echo "  test-js          - Run JavaScript tests only"
	@echo "  test-all         - Run all tests (Go + JavaScript)"
	@echo "  test-coverage    - Run all tests with coverage"
	@echo "  test-coverage-go - Run Go tests with coverage"
	@echo "  test-coverage-js - Run JavaScript tests with coverage"
	@echo "  test-race        - Run Go tests with race detection"
	@echo "  test-verbose     - Run tests in verbose mode"
	@echo "  test-run         - Run specific test (use TEST=TestName)"
	@echo "  test-js-watch    - Run JavaScript tests in watch mode"
	@echo "  benchmark        - Run benchmark tests"
	@echo "  deps             - Download Go dependencies"
	@echo "  deps-js          - Install JavaScript dependencies"
	@echo "  deps-all         - Install all dependencies"
	@echo "  deps-update      - Update dependencies"
	@echo "  fmt              - Format code"
	@echo "  vet              - Vet code"
	@echo "  lint             - Run Go linter"
	@echo "  lint-js          - Run JavaScript linter"
	@echo "  lint-all         - Run all linters"
	@echo "  run              - Build and run the application"
	@echo "  dev              - Start development server with hot reload"
	@echo "  install-tools    - Install development tools"
	@echo "  mocks            - Generate mocks"
	@echo "  security         - Run security checks"
	@echo "  docker-build     - Build Docker image"
	@echo "  docker-run       - Build and run Docker container"
	@echo "  check            - Quick quality check (fmt, vet, test-all)"
	@echo "  check-all        - Full quality check"
	@echo "  ci               - CI environment check (deps + check-all)"
	@echo "  help             - Display this help"
