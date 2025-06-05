# Command Line Applications

This directory contains the executable applications for the TrafficSim project.

## Applications

### simrunner
The main traffic simulation application with multiple execution modes.

**Location**: `cmd/simrunner/`

**Purpose**:
- Primary entry point for running traffic simulations
- Supports both command-line and web server modes
- Handles configuration loading and platform management

**Usage**:
```bash
# CLI mode (default)
./trafficsim -config data/config.yaml

# Web server mode
./trafficsim -web -port 8080

# Headless mode
./trafficsim -headless

# With multicast output
./trafficsim -multicast -multicast-addr 239.2.3.1 -multicast-port 6969
```

**Key Features**:
- Configuration-driven simulation setup
- Real-time web visualization interface
- CoT (Cursor on Target) message output
- Multicast UDP transmission support
- Graceful shutdown handling
- WebSocket-based real-time updates

### validate-yaml
YAML configuration validation utility.

**Location**: `cmd/validate-yaml/`

**Purpose**:
- Validates TrafficSim configuration files
- Ensures YAML syntax and semantic correctness
- Batch validation of multiple configuration files
- Integration with CI/CD pipelines

**Usage**:
```bash
# Validate single file
./validate-yaml data/config.yaml

# Validate multiple files
./validate-yaml data/configs/*.yaml

# Validate with detailed output
./validate-yaml -verbose data/config.yaml
```

**Validation Checks**:
- YAML syntax validation
- Configuration schema compliance
- Platform definition validation
- Route and scenario validation
- Physics parameter validation

## Building Applications

### Individual Applications
```bash
# Build simrunner
go build -o trafficsim ./cmd/simrunner

# Build validate-yaml
go build -o validate-yaml ./cmd/validate-yaml
```

### Using Make
```bash
# Build all applications
make build

# Build and run main application
make run

# Install applications to $GOPATH/bin
make install
```

## Testing

### Run Tests
```bash
# Test all command applications
go test ./cmd/...

# Test specific application
go test ./cmd/simrunner
go test ./cmd/validate-yaml
```

### Integration Tests
```bash
# Run full integration test suite
make test

# Test with coverage
make test-coverage
```

## Configuration

Both applications rely on YAML configuration files located in the `data/` directory:

- `data/config.yaml` - Main simulation configuration
- `data/configs/` - Alternative configuration scenarios
- `data/platforms/` - Platform type definitions
- `data/sample_routes/` - Route definitions

See the main [README.md](../README.md) for detailed configuration documentation.

## Development

### Adding New Commands

1. Create new directory under `cmd/`
2. Implement `main.go` with proper flag handling
3. Add tests in `main_test.go`
4. Update this README
5. Add to Makefile build targets

### Code Structure

Each command application should follow this structure:
```
cmd/
  new-command/
    main.go          # Application entry point
    main_test.go     # Application tests
    README.md        # Command-specific documentation (optional)
```

### Best Practices

- Use `flag` package for command-line arguments
- Implement proper error handling and logging
- Support graceful shutdown with signal handling
- Include comprehensive help text
- Write unit tests for core functionality
- Document all command-line flags

## Related Documentation

- [Main README](../README.md) - Project overview and usage
- [Architecture](../docs/ARCHITECTURAL_DESCRIPTION.md) - System design
- [Configuration Guide](../docs/) - Configuration details
- [Development Guide](../docs/) - Development practices
