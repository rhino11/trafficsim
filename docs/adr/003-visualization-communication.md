\
# ADR-003: Communication Method for Visualization

**Status:** Accepted
**Date:** 2025-06-03

## Context

The web-based UI needs to display simulated tracks in near real-time. This requires a mechanism for the backend (Go server) to push updates to the frontend (browser) efficiently.

## Decision

**Server-Sent Events (SSE)** will be the primary method for streaming track updates from the `internal/server` to the `web/` UI. **WebSockets** will be considered as a fallback or future enhancement if bidirectional communication becomes a strong requirement.

## Rationale

*   **Server-Sent Events (SSE):**
    *   **Simplicity:** SSE is a simpler protocol than WebSockets, built on standard HTTP. It's easier to implement on both client (EventSource API) and server.
    *   **Unidirectional:** Well-suited for server-to-client streaming, which is the primary need for track visualization.
    *   **Automatic Reconnection:** Browsers automatically handle reconnection if the connection drops.
    *   **Text-Based:** Easy to debug.
    *   **Sufficient for Prototype:** Meets the immediate needs of pushing track updates.

*   **WebSockets (as an alternative/future option):**
    *   **Bidirectional:** Allows two-way communication, which could be useful for future features like sending commands from the UI to the simulator.
    *   **Lower Latency (potentially):** Can offer slightly lower latency once the connection is established due to less HTTP overhead per message.
    *   **More Complex:** Requires more setup and handling on both client and server.

## Alternatives Considered

*   **Short Polling:** Client repeatedly asks the server for updates. Inefficient, high latency, and high server load. Not suitable for real-time updates.
*   **Long Polling:** Client makes a request, and the server holds it open until an update is available. Better than short polling but still has overhead and complexity compared to SSE or WebSockets.

## Consequences

*   **Server Implementation:** The `internal/server` will need an HTTP handler dedicated to SSE streaming. This handler will subscribe to updates from the `internal/sim` engine (likely via channels) and forward them to connected clients.
*   **Client Implementation:** The JavaScript in `web/static/js/` will use the `EventSource` API to connect to the SSE endpoint and update the map display.
*   **Unidirectional Focus:** If significant client-to-server real-time communication is needed later (beyond simple HTTP requests for initial data or configuration), migrating to or adding WebSockets might be necessary.
*   **Proxy/Firewall Issues:** SSE is generally well-supported by proxies as it's plain HTTP. WebSockets can sometimes have issues with older proxies.
*   **Connection Limits:** Browsers have a limit on the number of concurrent HTTP connections per domain (typically around 6), which includes SSE connections. This is unlikely to be an issue for a single-user local visualization.
