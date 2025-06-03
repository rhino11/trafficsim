\
# ADR-005: Geospatial Data Handling

**Status:** Proposed
**Date:** 2025-06-03

## Context

The simulator needs to model realistic platform movement, which involves interacting with geospatial data such as roads, shipping lanes, airspace boundaries, airports, and seaports. The complexity of this data and its processing can vary significantly.

## Decision

*   **Initial Prototype:** Geospatial data will be represented using **simplified, embedded, or file-based formats** (e.g., JSON, CSV, or hardcoded structures for basic scenarios). This includes:
    *   Defining points of interest (airports, seaports, military bases) with coordinates and types.
    *   Simple linear paths or waypoints for roads and shipping lanes.
    *   Basic polygonal or circular regions for airspace boundaries (if implemented in the early stages).
*   **Pathfinding:**
    *   For land vehicles on roads: A simple graph traversal (e.g., Dijkstra or A* on a predefined, small graph) or direct waypoint following.
    *   For maritime and airborne platforms: Great-circle navigation between waypoints or direct flight paths.
    *   For spacecraft: Simplified Keplerian orbital mechanics or predefined elliptical paths.
*   **Geospatial Calculations:** Basic calculations (distance, bearing) will be implemented within a `pkg/geospatial` utility package or use a lightweight, well-vetted Go library if one meets the needs without adding excessive complexity.
*   **Future Iterations:** Plan for the potential integration of more standard geospatial data formats (e.g., GeoJSON, Shapefiles) and libraries (e.g., PostGIS-like capabilities or Go-native spatial libraries) if the complexity and realism requirements increase significantly.

## Rationale

*   **Rapid Prototyping:** Starting with simplified data allows for quicker development of the core simulation engine and platform behaviors.
*   **Reduced Dependencies:** Avoids introducing heavy geospatial libraries or databases in the early stages.
*   **Focus on Core Logic:** Allows the team to focus on simulation mechanics and CoT output first.
*   **Scalability of Approach:** The design allows for replacing the simplified geospatial module with a more sophisticated one later without rewriting the entire simulation engine. The `internal/sim` will interact with geospatial data through defined interfaces.

## Alternatives Considered

*   **Full GIS Integration from Start (e.g., PostGIS, QGIS libraries):**
    *   **Pros:** Highly accurate and powerful.
    *   **Cons:** Significant upfront complexity in setup, data loading, and integration. Steep learning curve. Overkill for initial prototype.
*   **Using External Routing APIs (e.g., Google Maps, OSRM):**
    *   **Pros:** Offloads complex routing.
    *   **Cons:** Violates the requirement for the simulator to not draw on external streams and generate all data internally. Introduces external dependencies and potential costs.

## Consequences

*   **Initial Realism:** The realism of path adherence and environmental interaction will be limited in the first prototype. For example, vehicles might not follow complex road curvatures perfectly, or airspace interactions might be simplified.
*   **Data Preparation:** Some effort will be required to create or convert geospatial data into the simplified formats, even for basic scenarios. This data will reside in the `data/` directory.
*   **Modularity:** The `internal/sim` component will need to be designed to abstract its interaction with geospatial data, allowing the underlying implementation to be swapped out or enhanced later. An interface like `RouteProvider` or `SpatialContext` could be defined.
*   The `pkg/geospatial` package will house common, reusable geospatial calculations to avoid code duplication.
