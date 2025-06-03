\
# ADR-007: Modularity of Platform Types

**Status:** Proposed
**Date:** 2025-06-03

## Context

The simulator needs to support various platform types (maritime, airborne, land, space), each with potentially unique behaviors, attributes, and movement models. The design should allow for easy addition of new platform types or modification of existing ones without extensive changes to the core simulation engine.

## Decision

*   **Interface-Based Design:** A common Go **interface** (e.g., `Platform`) will be defined in `internal/models` or `internal/sim`. This interface will specify the essential methods that all platform types must implement.
    ```go
    package models // or sim

    type Platform interface {
        ID() string
        Type() string
        Update(deltaTime float64, spatialContext SpatialContext) error // deltaTime in seconds
        Position() (lat, lon, alt float64)
        CoTEvent() (CoTMessage, error) // Or similar method to get CoT data
        // Other common methods: Speed(), Heading(), FuelRemaining(), etc.
    }
    ```
*   **Concrete Implementations:** Each specific platform type (e.g., `Aircraft`, `Vessel`, `GroundVehicle`, `Spacecraft`) will be a concrete struct in `internal/models` that implements the `Platform` interface.
    ```go
    package models

    type Aircraft struct {
        // Common fields (e.g., from an embedded struct)
        // Aircraft-specific fields (e.g., flightLevel, callsign)
    }
    func (a *Aircraft) Update(deltaTime float64, spatialContext SpatialContext) error { /* ... */ }
    // ... other interface methods
    ```
*   **Simulation Engine Interaction:** The core simulation loop in `internal/sim` will operate on a collection of `Platform` interfaces (`[]Platform`). This decouples the engine from the concrete types.
*   **Factory Pattern (Optional):** A factory pattern or registration mechanism might be used in `internal/sim` or `internal/config` to create instances of specific platform types based on configuration.
*   **Behavioral Composition:** Common behaviors (e.g., basic physics, fuel consumption model) might be implemented in shared utility functions or embedded structs to promote code reuse across different platform types. Platform-specific behaviors will reside within their respective `Update` methods.

## Rationale

*   **Extensibility:** Adding a new platform type primarily involves creating a new struct that implements the `Platform` interface and potentially registering it with a factory. The core simulation engine requires minimal or no changes.
*   **Maintainability:** Code for each platform type is encapsulated within its own implementation, making it easier to understand, modify, and test.
*   **Polymorphism:** The simulation engine can treat all platforms uniformly through the `Platform` interface, simplifying its logic.
*   **Testability:** Individual platform types can be tested in isolation by mocking their dependencies (like `SpatialContext`).

## Alternatives Considered

*   **Large Switch Statements / Type Assertions:** The simulation engine could use type assertions or large switch statements to handle different platform types. This becomes unwieldy and error-prone as the number of types grows, violating the Open/Closed Principle.
*   **Inheritance (via struct embedding for shared fields):** Go doesn't have classical inheritance. While struct embedding can be used for sharing common fields and some methods, an interface-based approach is more robust for defining contracts and achieving polymorphism for behavior.
*   **Generic Programming (Go 1.18+):** While Go generics could be used to define collections of specific platform types, the interface approach remains fundamental for defining the behavioral contract that the simulation engine relies upon. Generics might be useful within the implementations or for utility functions operating on platforms.

## Consequences

*   The `Platform` interface needs to be carefully designed to include all necessary methods for the simulation engine and CoT output module to function correctly. It might evolve as new requirements emerge.
*   Each new platform type requires a new struct definition and implementation of all interface methods.
*   The `SpatialContext` interface (or a similar mechanism) will be crucial for providing platforms with information about their environment (e.g., terrain, roads, airspaces) in a decoupled way.
*   Initial platform implementations will focus on the core requirements, with more nuanced behaviors added iteratively.
