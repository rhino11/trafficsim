\
# ADR-002: Simulation State Management

**Status:** Accepted
**Date:** 2025-06-03

## Context

The simulation engine needs to manage the state of numerous platforms, each with its own attributes (position, velocity, fuel, etc.) and behaviors. This state will be updated frequently. We need an approach that is performant, allows for concurrent updates, and is manageable in terms of complexity.

## Decision

The primary simulation state (e.g., platform objects, current routes) will be managed **in-memory** within the Simulation Engine (`internal/sim`). Concurrency will be handled using **goroutines and channels** for updates and communication between different parts of the simulation.

## Rationale

*   **Performance:** In-memory access is the fastest way to read and update platform states, crucial for a real-time or faster-than-real-time simulation.
*   **Concurrency:**
    *   Goroutines allow individual platforms or groups of platforms to be processed concurrently.
    *   Channels can be used for safe communication of state changes or events between goroutines, minimizing the need for complex locking mechanisms if designed carefully.
    *   Where direct shared access to mutable state is unavoidable, Go's `sync` package (e.g., `sync.Mutex`, `sync.RWMutex`) will be used.
*   **Simplicity (for initial prototype):** Avoids the overhead and complexity of setting up and managing an external database or complex persistence layer for the core simulation loop.
*   **Go Idiomatic:** This approach aligns well with Go's strengths in concurrency and memory management.

## Alternatives Considered

*   **External Database (e.g., Redis, PostgreSQL):**
    *   **Pros:** Offers persistence, potentially easier querying, and could scale to larger-than-memory datasets.
    *   **Cons:** Adds significant latency for frequent updates, increases deployment complexity, and might be overkill for the initial prototype's requirements.
*   **Embedded Database (e.g., SQLite, BadgerDB):**
    *   **Pros:** Offers persistence without an external service.
    *   **Cons:** Still introduces I/O overhead compared to pure in-memory, and adds a dependency. Might be considered for future features like saving/loading simulation scenarios.
*   **Actor Model (e.g., using a library or custom implementation):**
    *   **Pros:** Provides strong encapsulation and concurrency management.
    *   **Cons:** Can add a layer of abstraction and complexity. Go's goroutines and channels provide similar benefits with less boilerplate for many use cases.

## Consequences

*   **Memory Limitation:** The maximum number of simulated platforms will be limited by available system memory. This is acceptable for the initial target scale.
*   **No Automatic Persistence:** Simulation state is lost when the application stops. Saving/loading scenarios would be a separate feature to be implemented if needed (potentially using serialization to disk or an embedded DB).
*   **Careful Concurrency Design:** Requires careful design to avoid race conditions and deadlocks. The Go race detector will be used during testing.
*   The `internal/models` package will define the structures for these in-memory objects.
*   The `internal/sim` package will contain the logic for updating these objects and managing their lifecycle.
