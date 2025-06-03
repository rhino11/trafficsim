\
# ADR-001: Language Choice

**Status:** Accepted
**Date:** 2025-06-03

## Context

The project requires a language suitable for building a performant, concurrent simulation engine that can handle network I/O (for CoT output and web UI updates) and is relatively easy to develop and deploy. Key considerations include performance, concurrency support, ecosystem, learning curve, and tooling.

## Decision

**Go (Golang)** is selected as the primary programming language for this project.

## Rationale

*   **Concurrency:** Go's built-in support for goroutines and channels makes it well-suited for concurrent programming, which is essential for simulating many platforms simultaneously and handling network requests.
*   **Performance:** Go compiles to native machine code and offers performance characteristics generally sufficient for this type of simulation, approaching C/C++ in many I/O-bound and CPU-bound tasks relevant to this project.
*   **Networking Libraries:** Go has a strong standard library for networking, simplifying the implementation of CoT output (HTTP/UDP) and the web server for visualization.
*   **Simplicity and Readability:** Go's syntax is relatively simple and designed for readability and maintainability, which can speed up development.
*   **Tooling:** Excellent tooling for building, testing, and dependency management (`go build`, `go test`, Go Modules).
*   **Deployment:** Produces single, statically-linked binaries by default, simplifying deployment.
*   **Ecosystem:** A growing ecosystem with libraries for various needs, though for core simulation logic, we will likely build custom components.

## Alternatives Considered

*   **Python:** Strong for rapid prototyping and has many libraries. However, can suffer from performance issues (GIL) for CPU-bound concurrent tasks without resorting to multiprocessing, which adds complexity.
*   **Rust:** Offers excellent performance and memory safety without a garbage collector. However, it has a steeper learning curve than Go, which could slow down initial development.
*   **C++:** Provides maximum performance and control. However, it has manual memory management, a more complex build system, and generally longer development cycles.
*   **Java/Kotlin (JVM):** Mature ecosystem and good performance. Can be more resource-intensive than Go and might be overkill for this project's initial scope.

## Consequences

*   The development team will need to be proficient in Go or be willing to learn it.
*   We will leverage Go's standard library and select external libraries carefully to maintain a lean application.
*   Garbage collection pauses, while generally minimal in Go, could be a factor to monitor under very high load, though unlikely to be an issue for the planned scale.
