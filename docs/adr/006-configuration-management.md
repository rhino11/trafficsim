\
# ADR-006: Configuration Management

**Status:** Proposed
**Date:** 2025-06-03

## Context

The simulator requires various parameters to be configurable, such as the number and types of platforms, simulation duration, CoT output endpoint, map settings, and behavioral parameters for platforms. A consistent and user-friendly way to manage this configuration is needed.

## Decision

*   **Configuration File:** The primary method for configuration will be a **single configuration file**.
*   **Format:** **YAML** is chosen as the configuration file format due to its human-readability and support for complex hierarchical structures. JSON will be considered as a secondary supported format if strong demand exists, but YAML is preferred.
*   **Loading:** The `internal/config` module will be responsible for loading and parsing the YAML file at startup. A library like `gopkg.in/yaml.v3` will be used.
*   **Structure:** Configuration values will be unmarshalled into Go structs defined within `internal/config` or `internal/models`.
*   **Defaults:** Sensible default values will be provided for most configuration options, allowing users to specify only what they need to change.
*   **Environment Variables (Optional Override):** For certain key parameters (e.g., CoT endpoint, log level), overrides via environment variables could be supported for flexibility in containerized or CI/CD environments. This would be a secondary mechanism.
*   **Command-Line Flags (Minimal):** Command-line flags will be used sparingly, primarily for specifying the path to the configuration file or for very common overrides like a "verbose" or "debug" mode.

## Rationale

*   **YAML:**
    *   **Human-Readable:** Easier for users to read and edit compared to JSON (especially with comments) or INI files for complex configurations.
    *   **Hierarchical:** Naturally supports nested configuration structures, which will be needed for different platform types and modules.
    *   **Comments:** Supports comments, which is crucial for documenting configuration options within the file itself.
    *   **Good Go Libraries:** Well-supported by Go libraries for parsing.
*   **Single File:** Simplifies management for users, especially for a standalone application.
*   **Centralized Logic (`internal/config`):** Encapsulates configuration loading and validation, making it easy for other modules to access configuration values.
*   **Defaults:** Improves user experience by reducing the amount of configuration needed for basic scenarios.

## Alternatives Considered

*   **JSON:** Less human-friendly for complex, commented configurations. Often used for machine-to-machine communication.
*   **INI Files:** Less suitable for hierarchical data.
*   **Environment Variables Only:** Can become unwieldy for a large number of configuration options. Better for a few key overrides.
*   **Command-Line Flags Only:** Not practical for complex or numerous settings.
*   **Database for Configuration:** Overkill for this application's needs.

## Consequences

*   A YAML parsing library will be added as a dependency.
*   Clear documentation for all configuration options and their structure in the YAML file will be essential. A sample configuration file should be provided.
*   The `internal/config` module will define Go structs that map to the YAML structure.
*   Changes to configuration structure will require updates to these Go structs and potentially migration guidance if backward compatibility is a concern.
*   Validation logic within `internal/config` will be needed to ensure that provided configuration values are sensible (e.g., ports are within range, required fields are present).
