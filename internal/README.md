# Internal Packages

This directory contains the core internal packages that implement the traffic simulation engine. These packages are not intended for external use and provide the fundamental building blocks of the simulation system.

## Package Overview

### config/
Configuration management and YAML parsing.

**Purpose**: 
- Load and validate YAML configuration files
- Provide typed configuration structures
- Handle platform factory initialization
- Support multiple configuration scenarios

**Key Components**:
- `Config` struct for main simulation parameters
- Platform factory creation and validation
- Configuration loading with error handling
- Default configuration values

### models/
Platform definitions and physics implementations.

**Purpose**:
- Define platform types (airborne, maritime, land, space)
- Implement platform-specific physics and behavior
- Provide universal platform interface
- Handle platform state management

**Key Components**:
- `Platform` interface for universal platform operations
- Specialized platform types (aircraft, ships, vehicles, satellites)
- Physics calculations for movement and navigation
- Platform factory for creation and initialization

### server/
Web server and API endpoints.

**Purpose**:
- HTTP/WebSocket server for web interface
- REST API for platform management
- Real-time communication with frontend
- Static file serving for web assets

**Key Components**:
- HTTP router setup and middleware
- WebSocket connection management
- API endpoints for simulation control
- Template rendering for web interface

### sim/
Physics engine and simulation core.

**Purpose**:
- Core simulation loop and timing
- Physics calculations and movement updates
- Platform interaction and collision detection
- Simulation state management

**Key Components**:
- Physics engine for realistic movement calculations
- Simulation timestep management
- Platform position and velocity updates
- Environmental factors (gravity, air density)

### output/
External data output and transmission.

**Purpose**:
- CoT (Cursor on Target) message generation
- Multicast UDP transmission
- Platform data serialization
- Real-time data streaming

**Key Components**:
- CoT XML message formatting
- UDP multicast transmission
- Platform data conversion
- Output configuration management

### testutil/
Testing utilities and helpers.

**Purpose**:
- Common testing infrastructure
- Mock objects and test fixtures
- Logging configuration for tests
- Test data generation utilities

**Key Components**:
- Logger setup for consistent test output
- Helper functions for test data creation
- Mock implementations for testing
- Test configuration utilities

## Package Dependencies

```
config/ ─────┐
             ▼
models/ ──▶ sim/ ──▶ server/
   ▲         ▲         ▲
   │         │         │
   └─────────┼─────────┘
             │
           output/
             ▲
             │
         testutil/
```

## Development Guidelines

### Internal Package Rules
- Packages in `internal/` are not importable by external projects
- All public APIs within internal packages should be well documented
- Each package should have a clear, single responsibility
- Avoid circular dependencies between packages

### Testing Standards
- Each package must have comprehensive unit tests
- Use `testutil/` for common testing infrastructure
- Aim for >80% test coverage across all packages
- Include integration tests where appropriate

### Code Quality
- Follow Go best practices and idioms
- Use consistent error handling patterns
- Implement proper logging throughout
- Document all exported functions and types

## Current Test Coverage

| Package | Coverage | Status | Priority |
|---------|----------|--------|----------|
| models/ | 84.7% | Good | Medium |
| config/ | 29.0% | Needs Work | High |
| server/ | 34.8% | Needs Work | High |
| sim/ | 32.5% | Needs Work | High |
| output/ | 72.1% | Good | Low |
| testutil/ | N/A | Support | Low |

## Adding New Packages

1. Create package directory under `internal/`
2. Implement core functionality with proper interfaces
3. Add comprehensive unit tests
4. Update package documentation
5. Add to dependency diagram if applicable
6. Update this README

## Security Considerations

- All input validation happens in these packages
- Configuration parsing includes security checks
- HTTP handlers implement proper sanitization
- Platform data is validated before processing

## Performance Notes

- Physics calculations are optimized for real-time performance
- Platform updates are batched for efficiency
- Memory allocation is minimized in hot paths
- Concurrent processing where appropriate

## Related Documentation

- [Architecture Overview](../docs/ARCHITECTURAL_DESCRIPTION.md)
- [ADR Directory](../docs/adr/) - Architecture decisions
- [Testability Guide](../docs/TESTABILITY_GUIDE.md)
- [Main README](../README.md)
