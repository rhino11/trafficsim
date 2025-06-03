\
# ADR-004: Cursor on Target (CoT) Output

**Status:** Accepted
**Date:** 2025-06-03

## Context

A core requirement is for the simulator to output platform data as Cursor on Target (CoT) XML objects. This output needs to be sent to a specified URL/endpoint or IP address. We need to decide on the transport protocol and how the CoT generation will be integrated.

## Decision

*   **CoT Generation:** The `internal/output` module will be responsible for generating CoT XML strings from the platform data structures provided by the `internal/sim` module. Standard Go `encoding/xml` package will be used.
*   **Transport Protocol:** The simulator will support sending CoT messages via **HTTP POST** and **UDP unicast**. The specific endpoint (URL for HTTP, IP:Port for UDP) and protocol will be configurable.
*   **Configuration:** The output endpoint, protocol, and transmission frequency will be configurable via the application's configuration file (see ADR-006).

## Rationale

*   **CoT Standard:** Adherence to CoT XML schema is paramount for interoperability with public safety systems.
*   **HTTP POST:**
    *   **Reliability:** TCP-based, ensuring message delivery (if the endpoint is reachable and acknowledges).
    *   **Commonly Supported:** Widely accepted by web services and easy to integrate with many existing tools.
    *   **Overhead:** Higher overhead per message compared to UDP due to TCP handshakes and HTTP headers.
*   **UDP Unicast:**
    *   **Low Latency/Overhead:** "Fire and forget" nature makes it suitable for high-frequency updates where occasional packet loss is acceptable.
    *   **Simplicity:** Simpler protocol than TCP.
    *   **Unreliability:** No guarantee of delivery or ordering.
*   **Configurability:** Providing both HTTP and UDP options offers flexibility for different consumer systems and network environments.
*   **Dedicated Module (`internal/output`):** Separates the concern of CoT generation and transmission from the core simulation logic, improving modularity.

## Alternatives Considered

*   **TCP Socket (raw):** More complex to implement than HTTP and offers similar reliability benefits but without the standardized application layer protocol.
*   **Message Queues (e.g., RabbitMQ, Kafka):**
    *   **Pros:** Highly scalable, offers persistence, decoupling.
    *   **Cons:** Adds significant external dependencies and complexity, which is overkill for the initial requirement of direct CoT output to a specified endpoint. Could be a future consideration for more advanced scenarios.
*   **gRPC/Protobuf:** Efficient binary protocol. However, CoT is explicitly an XML-based standard, so this would only be relevant for internal communication, not the primary CoT output.

## Consequences

*   The `internal/output` module will need to implement:
    *   Logic to marshal platform data into valid CoT XML.
    *   HTTP client logic for sending POST requests.
    *   UDP client logic for sending datagrams.
*   The `internal/sim` module will provide platform state updates to the `internal/output` module at configurable intervals.
*   Error handling for network transmission (e.g., connection refused for HTTP, network unreachable) needs to be implemented and logged appropriately.
*   The choice between HTTP and UDP will impact reliability guarantees. Users must be aware of this when configuring the output.
