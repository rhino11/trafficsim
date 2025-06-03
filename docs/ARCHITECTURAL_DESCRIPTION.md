\
# Architectural Description: Global Traffic Simulator (GTS)

**Version:** 0.1.0
**Date:** 2025-06-03
**Status:** Proposed

## 1. Introduction and Goals

This document outlines the architecture for the Global Traffic Simulator (GTS), a Go-based application designed to emulate a large volume of realistic maritime, airborne, land, and space tracks.

**Primary Goals:**

*   **Realistic Simulation:** Platforms exhibit behaviors authentic to their type (e.g., vehicles follow roads, aircraft adhere to airspace, spacecraft follow orbital mechanics).
*   **Scalability:** Capable of simulating a large number of concurrent tracks.
*   **Internal Data Generation:** Operates without reliance on external live data streams.
*   **Standardized Output:** Generates Cursor on Target (CoT) XML objects.
*   **Configurability:** Allows users to define simulation scenarios, including platform types, numbers, start/end points, and operational parameters.
*   **Local Visualization:** Provides a basic localhost map interface to display simulated tracks for rapid prototyping and validation.
*   **Extensibility:** Designed to easily incorporate new platform types, behaviors, and data sources in the future.

**Non-Goals (for initial prototype):**

*   Advanced 3D visualization beyond simple map markers.
*   Distributed simulation capabilities.
*   User authentication or complex UI for configuration.

## 2. System Overview

GTS is a modular system composed of several key components that interact to generate, manage, and output simulated track data.

```
+---------------------+      +---------------------+      +---------------------+
|   Configuration     |----->|  Simulation Engine  |<---->|     Data Models     |
|  (config.yaml/json) |      | (internal/sim)      |      | (internal/models)   |
+---------------------+      +---------------------+      +---------------------+
                                      |
                                      | Outputs CoT
                                      v
+---------------------+      +---------------------+
|   CoT Output Module |----->|   Network Endpoint  |
| (internal/output)   |      | (e.g., HTTP/UDP)    |
+---------------------+      +---------------------+
          ^
          | Updates
          |
+---------------------+      +---------------------+
|   Web Server        |<---->|   Web UI (Map)      |
| (internal/server)   |      | (web/)              |
+---------------------+      +---------------------+
```

## 3. Component Breakdown

### 3.1. `cmd/simrunner` - Main Application
*   **Responsibilities:**
    *   Entry point of the application.
    *   Parses command-line arguments.
    *   Initializes and wires together all other components (configuration, simulation engine, output, server).
    *   Manages the main application lifecycle.

### 3.2. `internal/config` - Configuration Management
*   **Responsibilities:**
    *   Loading simulation parameters from a configuration file (e.g., YAML, JSON) or environment variables.
    *   Validating configuration.
    *   Providing access to configuration values for other modules.
*   **Key Decisions:** ADR-006 (Configuration Management).

### 3.3. `internal/models` - Data Models
*   **Responsibilities:**
    *   Defining core data structures for platforms (e.g., `Aircraft`, `Vessel`, `Vehicle`, `Spacecraft`), their attributes (ID, type, position, velocity, heading, fuel), routes, waypoints, and geographic features (airports, ports, roads - simplified initially).
*   **Key Decisions:** ADR-007 (Modularity of Platform Types).

### 3.4. `internal/sim` - Simulation Engine
*   **Responsibilities:**
    *   Core logic for managing the state of all simulated platforms.
    *   Updating platform positions, speeds, and orientations based on their type, defined behaviors, and environmental constraints (roads, shipping lanes, airspaces, orbits).
    *   Implementing realistic behaviors:
        *   Pathfinding and adherence (e.g., A* for roads, great-circle for long-distance flights).
        *   Takeoff, landing, docking, and undocking procedures.
        *   Fuel consumption and range limitations.
        *   Distinction between military and commercial facility usage.
    *   Managing the simulation clock and event scheduling.
*   **Key Decisions:** ADR-002 (Simulation State Management), ADR-005 (Geospatial Data Handling).

### 3.5. `internal/output` - CoT Output Module
*   **Responsibilities:**
    *   Generating CoT XML messages from platform data provided by the Simulation Engine.
    *   Formatting CoT messages according to relevant standards.
    *   Sending CoT messages to a specified URL/endpoint or IP address (HTTP POST or UDP).
*   **Key Decisions:** ADR-004 (Cursor on Target Output).

### 3.6. `internal/server` - Web Server & Real-time Updates
*   **Responsibilities:**
    *   Serving the static files for the Web UI (HTML, CSS, JavaScript).
    *   Providing an API endpoint for the Web UI to fetch initial track data.
    *   Streaming real-time track updates to the Web UI (e.g., via WebSockets or Server-Sent Events).
*   **Key Decisions:** ADR-003 (Communication Method for Visualization).

### 3.7. `pkg/` - Shared Packages
*   **Responsibilities:**
    *   Contains reusable utility packages that are not specific to the internal workings of this application but could be useful in other projects.
    *   Example: `pkg/geospatial` for common geographic calculations (distance, bearing, coordinate conversion).

### 3.8. `web/` - Web User Interface
*   **Responsibilities:**
    *   `web/templates/`: HTML templates for the map display.
    *   `web/static/`: CSS and JavaScript files for the frontend.
    *   Displaying simulated tracks on a 2D map (e.g., using Leaflet.js or OpenLayers).
    *   Receiving and rendering real-time track updates from the `internal/server`.

### 3.9. `data/` - Simulation Data
*   **Responsibilities:**
    *   Storing static data required for the simulation, such as:
        *   Sample route definitions.
        *   Airport/port locations and types.
        *   Simplified road network data or shipping lane definitions.
        *   Platform characteristics (speed, range, fuel capacity).
    *   Initially, this might be simple CSV or JSON files.

### 3.10. `docs/` - Documentation
*   **Responsibilities:**
    *   `docs/adr/`: Architecture Decision Records.
    *   Other design documents, user guides, etc.

### 3.11. `scripts/` - Utility Scripts
*   **Responsibilities:**
    *   Helper scripts for building, running, testing, linting, generating data, etc.

## 4. Data Management

*   **Platform State:** Primarily managed in-memory within the Simulation Engine for performance. Goroutines and channels will be used for concurrent updates.
*   **Geospatial Data:** Initial versions will use simplified, embedded, or file-based representations of roads, airspaces, etc. Future iterations might integrate with more complex geospatial libraries or databases if needed.
*   **Configuration Data:** Loaded from files at startup.

## 5. Key Architectural Decisions (ADRs)

Refer to the `docs/adr/` directory for detailed ADRs. Summary:
*   **ADR-001:** Language Choice (Go)
*   **ADR-002:** Simulation State Management (In-memory, concurrent)
*   **ADR-003:** Communication for Visualization (SSE/WebSockets)
*   **ADR-004:** Cursor on Target Output (Configurable HTTP/UDP)
*   **ADR-005:** Geospatial Data Handling (Simplified initially, file-based)
*   **ADR-006:** Configuration Management (YAML/JSON file)
*   **ADR-007:** Modularity of Platform Types (Interface-based design)

## 6. Deployment View

The application will be a single executable binary. It can be run directly on a local machine.
For the web UI, a modern browser is required.

## 7. Non-Functional Requirements

*   **Performance:** The system should be able to simulate hundreds to thousands of tracks in real-time or faster-than-real-time on typical developer hardware for the prototype.
*   **Maintainability:** Code should be well-structured, documented, and testable.
*   **Extensibility:** The architecture should allow for adding new platform types, behaviors, and output formats with reasonable effort.

## 8. Future Considerations

*   Integration with more detailed geospatial datasets (e.g., OpenStreetMap, FAA airspace data).
*   Support for distributed simulation.
*   More sophisticated UI for scenario definition and control.
*   Plugin architecture for custom behaviors or data sources.
*   Persistence of simulation state for resume capabilities.
---

## Coding Conventions and Best Practices

1.  **Project Structure:** Adhere to the defined project layout (see Section 3).
2.  **Go Best Practices:**
    *   Follow standard Go formatting (`gofmt`/`goimports`).
    *   Use `golangci-lint` with a sensible configuration for static analysis.
    *   Effective Go principles should be followed.
3.  **Naming Conventions:**
    *   Packages: `lowercase`, short, and concise.
    *   Types, Functions, Variables: `PascalCase` for exported, `camelCase` for unexported.
    *   Interfaces: Often end with `er` (e.g., `Reader`, `Writer`, `Simulator`).
4.  **Error Handling:**
    *   Use `error` return values. Avoid panics for recoverable errors.
    *   Provide context to errors using `fmt.Errorf("module: operation failed: %w", err)` when wrapping.
5.  **Concurrency:**
    *   Use goroutines and channels appropriately.
    *   Protect shared data with mutexes or use channels for synchronization.
    *   Be mindful of potential deadlocks and race conditions. Use the race detector (`go test -race`).
6.  **Logging:**
    *   Use a structured logging library (e.g., Go's built-in `slog` (Go 1.21+), or a library like `zerolog` or `zap`).
    *   Provide different log levels (DEBUG, INFO, WARN, ERROR).
7.  **Testing:**
    *   Write unit tests for all critical components and logic. Place them in `_test.go` files.
    *   Aim for high test coverage.
    *   Use table-driven tests where appropriate.
    *   Consider integration tests for interactions between components.
8.  **Dependencies:**
    *   Manage dependencies using Go Modules (`go.mod`, `go.sum`).
    *   Minimize external dependencies. Justify additions.
9.  **Documentation:**
    *   Document all exported functions, types, and important internal logic.
    *   Maintain and update ADRs.
10. **Configuration:**
    *   Externalize configuration (e.g., into YAML or JSON files). Avoid hardcoding.
11. **Makefile/Scripts:**
    *   Provide a `Makefile` or scripts in `scripts/` for common development tasks (e.g., `build`, `test`, `run`, `lint`, `clean`).
12. **API Design:**
    *   Keep APIs simple and focused.
    *   Strive for backward compatibility where possible once an API is considered stable.
