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
VALIDATOR_PATH=./cmd/validate-yaml

# Runtime parameters
DEFAULT_PORT=8080
DEFAULT_MULTICAST_ADDR=239.2.3.1
DEFAULT_MULTICAST_PORT=6969

# Package building parameters
PACKAGE_VERSION ?= 1.0.0
PACKAGE_ARCH ?= amd64
DIST_DIR = dist

# Build targets
.PHONY: all build clean test test-go test-js test-all test-package test-package-coverage test-package-race test-coverage test-coverage-go test-coverage-js test-verbose deps deps-js fmt vet lint run run-headless run-web run-multicast run-web-multicast validate-yaml help

# Default target - now includes YAML validation
all: validate-yaml test-all build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH)

# Build for Linux
build-linux:
	@echo "Building for Linux..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v $(MAIN_PATH)

# Validate YAML configuration files
validate-yaml:
	@echo "Validating YAML configuration files..."
	@$(GOBUILD) -o validate-yaml $(VALIDATOR_PATH)
	@./validate-yaml data/
	@rm -f validate-yaml

# Validate specific YAML file or directory
validate-yaml-file:
	@echo "Validating specific YAML file/directory: $(FILE)"
	@if [ -z "$(FILE)" ]; then echo "Error: FILE variable not set. Use: make validate-yaml-file FILE=path/to/file.yaml"; exit 1; fi
	@$(GOBUILD) -o validate-yaml $(VALIDATOR_PATH)
	@./validate-yaml $(FILE)
	@rm -f validate-yaml

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f validate-yaml

# Run Go tests
test-go:
	@echo "Running Go tests..."
	$(GOTEST) -v ./...

# Run tests for a specific package
test-package:
	@echo "Running tests for package: $(PKG)"
	@if [ -z "$(PKG)" ]; then echo "Error: PKG variable not set. Use: make test-package PKG=./internal/server"; exit 1; fi
	$(GOTEST) -v $(PKG)

# Run tests for a specific package with coverage
test-package-coverage:
	@echo "Running tests with coverage for package: $(PKG)"
	@if [ -z "$(PKG)" ]; then echo "Error: PKG variable not set. Use: make test-package-coverage PKG=./internal/server"; exit 1; fi
	$(GOTEST) -v -coverprofile=coverage-$(shell basename $(PKG)).out $(PKG)
	$(GOCMD) tool cover -html=coverage-$(shell basename $(PKG)).out -o coverage-$(shell basename $(PKG)).html
	@echo "Coverage report generated: coverage-$(shell basename $(PKG)).html"

# Run tests for a specific package with race detection
test-package-race:
	@echo "Running tests with race detection for package: $(PKG)"
	@if [ -z "$(PKG)" ]; then echo "Error: PKG variable not set. Use: make test-package-race PKG=./internal/server"; exit 1; fi
	$(GOTEST) -v -race $(PKG)

# Run JavaScript tests
test-js:
	@echo "Running JavaScript tests..."
	@which npm > /dev/null || (echo "Error: npm is not installed" && exit 1)
	$(NPMCMD) test

# Run all tests (Go + JavaScript)
test-all: test-go test-js

# Legacy test target (runs all tests for backward compatibility)
test: test-all

# Run Go tests with coverage (quiet mode - coverage summary only)
test-coverage-go-quiet:
	@echo "Running Go tests with coverage (quiet mode)..."
	$(GOTEST) -coverprofile=coverage.out ./...
	@echo ""
	@echo "=== Coverage Summary ==="
	$(GOCMD) tool cover -func=coverage.out
	@echo ""
	@echo "HTML coverage report generated: coverage.html"
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Run Go tests with coverage
test-coverage-go:
	@echo "Running Go tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Go coverage report generated: coverage.html"

# Run Go tests with coverage and detailed function-level breakdown
test-coverage-go-detailed:
	@echo "Running Go tests with detailed coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	@echo ""
	@echo "=== Detailed Coverage Breakdown ==="
	$(GOCMD) tool cover -func=coverage.out
	@echo ""
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

# Run the application (backward compatibility - defaults to CLI mode)
run:
	@echo "Running $(BINARY_NAME) in CLI mode..."
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH) && ./$(BINARY_NAME)

# Run in headless mode (no frontend, command-line only)
run-headless:
	@echo "Running $(BINARY_NAME) in headless mode (no frontend)..."
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH) && ./$(BINARY_NAME) -headless

# Run with web frontend
run-web:
	@echo "Running $(BINARY_NAME) with web frontend on port $(DEFAULT_PORT)..."
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH) && ./$(BINARY_NAME) -web -port=$(DEFAULT_PORT)

# Run with web frontend on custom port
run-web-port:
	@echo "Running $(BINARY_NAME) with web frontend on port $(PORT)..."
	@if [ -z "$(PORT)" ]; then echo "Error: PORT variable not set. Use: make run-web-port PORT=8081"; exit 1; fi
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH) && ./$(BINARY_NAME) -web -port=$(PORT)

# Run in headless mode with multicast transmission
run-multicast:
	@echo "Running $(BINARY_NAME) in headless mode with multicast transmission..."
	@echo "Multicast address: $(DEFAULT_MULTICAST_ADDR):$(DEFAULT_MULTICAST_PORT)"
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH) && ./$(BINARY_NAME) -headless -multicast -multicast-addr=$(DEFAULT_MULTICAST_ADDR) -multicast-port=$(DEFAULT_MULTICAST_PORT)

# Run with web frontend and multicast transmission
run-web-multicast:
	@echo "Running $(BINARY_NAME) with web frontend and multicast transmission..."
	@echo "Web server on port $(DEFAULT_PORT), multicast on $(DEFAULT_MULTICAST_ADDR):$(DEFAULT_MULTICAST_PORT)"
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH) && ./$(BINARY_NAME) -web -port=$(DEFAULT_PORT) -multicast -multicast-addr=$(DEFAULT_MULTICAST_ADDR) -multicast-port=$(DEFAULT_MULTICAST_PORT)

# Run with custom multicast settings
run-multicast-custom:
	@echo "Running $(BINARY_NAME) with custom multicast settings..."
	@if [ -z "$(ADDR)" ]; then echo "Error: ADDR variable not set. Use: make run-multicast-custom ADDR=239.255.42.100 PORT=9998"; exit 1; fi
	@if [ -z "$(PORT)" ]; then echo "Error: PORT variable not set. Use: make run-multicast-custom ADDR=239.255.42.100 PORT=9998"; exit 1; fi
	@echo "Multicast address: $(ADDR):$(PORT)"
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH) && ./$(BINARY_NAME) -headless -multicast -multicast-addr=$(ADDR) -multicast-port=$(PORT)

# Run all modes for testing (in background with different ports/addresses)
run-all-modes:
	@echo "Starting all running modes for testing..."
	@echo "Starting headless mode in background..."
	@$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH) && ./$(BINARY_NAME) -headless &
	@sleep 2
	@echo "Starting web mode on port 8081 in background..."
	@./$(BINARY_NAME) -web -port=8081 &
	@sleep 2
	@echo "Starting multicast mode on different address in background..."
	@./$(BINARY_NAME) -headless -multicast -multicast-addr=239.255.42.100 -multicast-port=9998 &
	@echo "All modes started. Use 'make stop-all' to stop them."
	@echo "Web interface available at: http://localhost:8081"
	@echo "Multicast streams on: 239.255.42.99:9999 and 239.255.42.100:9998"

# Stop all background processes
stop-all:
	@echo "Stopping all $(BINARY_NAME) processes..."
	@pkill -f $(BINARY_NAME) || echo "No $(BINARY_NAME) processes found"

# Download and install air (if not already installed)
install-air:
	@echo "Installing air for hot reloading..."
	@which air > /dev/null || (go install github.com/air-verse/air@latest)

# Development mode with hot reloading (requires air)
dev:
	@echo "Starting development server with hot reload..."
	@$(MAKE) install-air
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

# Prepare artifacts for CI pipeline
prepare-artifacts: test-coverage
	@echo "Preparing artifacts for CI pipeline..."
	@mkdir -p coverage test-results
	@cp coverage.out coverage/ 2>/dev/null || true
	@cp coverage.html coverage/ 2>/dev/null || true
	@cp coverage-*.out coverage/ 2>/dev/null || true
	@cp coverage-*.html coverage/ 2>/dev/null || true
	@cp -r coverage-js/* coverage/ 2>/dev/null || true
	@echo "Artifacts prepared in coverage/ directory"

# Clean artifacts directories
clean-artifacts:
	@echo "Cleaning artifact directories..."
	@rm -rf coverage test-results

# Quick quality check - now includes YAML validation
check: validate-yaml fmt vet test-all

# Full quality check - now includes YAML validation
check-all: validate-yaml check lint-all test-race test-coverage security

# CI check (for CI environments) - now includes YAML validation
ci: deps-all validate-yaml check-all

# Set up package directories
setup-package-dirs:
	@echo "Setting up package directories..."

# Package building targets
# Package building targets
# These are simplified versions without heredocs to avoid syntax issues

# Basic package directory setup
setup-package-dirs:
	@echo "Setting up package directories..."
	@mkdir -p $(DIST_DIR)

# Linux RPM package (simplified)
build-package-rpm: setup-package-dirs
	@echo "Building RPM package for $(PACKAGE_ARCH)..."
	@echo "Note: RPM building requires rpmbuild tool and proper setup"
	@mkdir -p rpm/BUILD
	@$(GOBUILD) -ldflags="-s -w" -o $(DIST_DIR)/$(BINARY_NAME)-$(PACKAGE_VERSION).linux-$(PACKAGE_ARCH) $(MAIN_PATH)

# Linux DEB package (simplified)
build-package-deb: setup-package-dirs
	@echo "Building DEB package for $(PACKAGE_ARCH)..."
	@echo "Note: DEB building requires dpkg-deb tool"
	@mkdir -p deb/BUILD
	@$(GOBUILD) -ldflags="-s -w" -o $(DIST_DIR)/$(BINARY_NAME)-$(PACKAGE_VERSION).linux-$(PACKAGE_ARCH).deb $(MAIN_PATH)

# Windows EXE
build-package-exe: setup-package-dirs
	@echo "Building Windows EXE for $(PACKAGE_ARCH)..."
	@GOOS=windows GOARCH=$(PACKAGE_ARCH) $(GOBUILD) -ldflags="-s -w" -o $(DIST_DIR)/$(BINARY_NAME)-$(PACKAGE_VERSION).windows-$(PACKAGE_ARCH).exe $(MAIN_PATH)

# macOS binary
build-package-macos: setup-package-dirs
	@echo "Building macOS binary for $(PACKAGE_ARCH)..."
	@GOOS=darwin GOARCH=$(PACKAGE_ARCH) $(GOBUILD) -ldflags="-s -w" -o $(DIST_DIR)/$(BINARY_NAME)-$(PACKAGE_VERSION).darwin-$(PACKAGE_ARCH) $(MAIN_PATH)

# Linux binary
build-package-linux: setup-package-dirs
	@echo "Building Linux binary for $(PACKAGE_ARCH)..."
	@GOOS=linux GOARCH=$(PACKAGE_ARCH) $(GOBUILD) -ldflags="-s -w" -o $(DIST_DIR)/$(BINARY_NAME)-$(PACKAGE_VERSION).linux-$(PACKAGE_ARCH) $(MAIN_PATH)

# Test package building (simplified)
test-package-builds: build-package-linux build-package-macos build-package-exe
	@echo "Successfully built basic packages"
	@ls -la $(DIST_DIR)/

# Clean package artifacts
clean-packages:
	@echo "Cleaning package artifacts..."
	@rm -rf $(DIST_DIR) rpm deb appimage wix nsis dmg android ios
	@echo "Package artifacts cleaned"

.PHONY: setup-package-dirs build-package-rpm build-package-deb build-package-exe build-package-macos build-package-linux test-package-builds clean-packages
