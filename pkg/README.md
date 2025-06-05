# Public Packages

This directory contains packages that are intended for external use and provide reusable functionality that other projects can import and utilize.

## Package Overview

### geospatial/
Geospatial calculations and utilities for geographic coordinate systems.

**Purpose**: 
- Provide geographic coordinate calculations
- Handle coordinate system transformations
- Implement geospatial distance and bearing calculations
- Support various geographic projections and datums

**Key Features**:
- Latitude/longitude coordinate handling
- Distance calculations (great circle, rhumb line)
- Bearing and heading calculations
- Coordinate system conversions
- Geographic bounds and region calculations

## Package Philosophy

Packages in `pkg/` follow these principles:

### Public API Design
- **Stable Interfaces**: Public APIs are designed for long-term stability
- **Clear Documentation**: All exported functions and types are well-documented
- **Semantic Versioning**: Changes follow semantic versioning principles
- **Minimal Dependencies**: Reduce external dependencies where possible

### Reusability
- **Generic Implementation**: Designed for use beyond TrafficSim
- **Standard Patterns**: Follow Go community standards and idioms
- **Composable Design**: Functions and types can be combined effectively
- **Performance Optimized**: Efficient implementations for real-time use

## geospatial Package

### Core Types
```go
// Coordinate represents a geographic position
type Coordinate struct {
    Latitude  float64
    Longitude float64
    Altitude  float64  // Optional, in meters
}

// Distance represents a calculated distance
type Distance struct {
    Meters      float64
    Kilometers  float64
    NauticalMiles float64
}

// Bearing represents a calculated bearing
type Bearing struct {
    Degrees float64
    Radians float64
}
```

### Key Functions
```go
// Distance calculations
func DistanceHaversine(from, to Coordinate) Distance
func DistanceVincenty(from, to Coordinate) Distance

// Bearing calculations  
func BearingTo(from, to Coordinate) Bearing
func FinalBearing(from, to Coordinate) Bearing

// Coordinate transformations
func ProjectToUTM(coord Coordinate) UTMCoordinate
func ProjectFromUTM(utm UTMCoordinate) Coordinate

// Movement calculations
func MoveCoordinate(start Coordinate, bearing Bearing, distance Distance) Coordinate
func IntermediatePoint(from, to Coordinate, fraction float64) Coordinate
```

### Usage Examples

#### Basic Distance Calculation
```go
import "github.com/rhino11/trafficsim/pkg/geospatial"

start := geospatial.Coordinate{
    Latitude:  40.7128,
    Longitude: -74.0060,
}

end := geospatial.Coordinate{
    Latitude:  34.0522,
    Longitude: -118.2437,
}

distance := geospatial.DistanceHaversine(start, end)
fmt.Printf("Distance: %.2f km\n", distance.Kilometers)
```

#### Navigation Calculations
```go
// Calculate bearing from New York to Los Angeles
bearing := geospatial.BearingTo(start, end)
fmt.Printf("Initial bearing: %.1f°\n", bearing.Degrees)

// Move 100km in that direction
newPos := geospatial.MoveCoordinate(start, bearing, geospatial.Distance{
    Kilometers: 100,
})
```

#### Route Planning
```go
// Calculate intermediate waypoints
waypoints := make([]geospatial.Coordinate, 10)
for i := 0; i < 10; i++ {
    fraction := float64(i) / 9.0
    waypoints[i] = geospatial.IntermediatePoint(start, end, fraction)
}
```

## Performance Characteristics

### geospatial Package Benchmarks
```
BenchmarkDistanceHaversine-8     1000000    1052 ns/op    0 allocs/op
BenchmarkDistanceVincenty-8       500000    2847 ns/op    0 allocs/op
BenchmarkBearingTo-8             2000000     654 ns/op    0 allocs/op
BenchmarkMoveCoordinate-8        1000000    1123 ns/op    0 allocs/op
```

### Optimization Features
- **Zero Allocations**: Core calculations avoid memory allocations
- **Vectorized Math**: Uses efficient mathematical operations
- **Caching**: Expensive calculations cached where appropriate
- **Precision Control**: Configurable precision vs performance trade-offs

## Testing and Quality

### Test Coverage
- Comprehensive unit tests for all public functions
- Property-based testing for mathematical invariants
- Benchmark tests for performance validation
- Integration tests with known geographic datasets

### Validation
- Cross-validation with established geographic libraries
- Testing against known geographic control points
- Precision validation for various coordinate ranges
- Edge case testing (poles, antimeridian, etc.)

## Dependencies

### Current Dependencies
- **Standard Library Only**: No external dependencies for core functionality
- **Math Package**: Uses Go's math package for trigonometric functions
- **Testing**: Uses standard testing and benchmark frameworks

### Future Considerations
- **PROJ Library**: Potential integration for advanced projections
- **Geographic Datasets**: Optional integration with geographic reference data
- **Precision Math**: Higher precision coordinate calculations if needed

## Usage in TrafficSim

The geospatial package is used throughout TrafficSim for:
- **Platform Movement**: Calculating new positions based on speed and heading
- **Distance Calculations**: Determining proximity between platforms
- **Route Planning**: Calculating waypoints and navigation paths
- **Coordinate Conversion**: Converting between different coordinate systems
- **Visualization**: Supporting map display and platform positioning

## External Usage

This package is designed for use in other projects requiring geospatial calculations:
- **Navigation Systems**: GPS and mapping applications
- **Simulation Software**: Other traffic and movement simulations
- **Geographic Analysis**: Spatial analysis and GIS applications
- **Logistics Software**: Route planning and optimization
- **Gaming**: Location-based games and simulations

## API Stability

### Current Status
- **Version**: 1.0.0
- **Stability**: Stable API, no breaking changes planned
- **Documentation**: Complete API documentation available
- **Testing**: Comprehensive test coverage

### Versioning Policy
- **Major Version**: Breaking API changes
- **Minor Version**: New features, backward compatible
- **Patch Version**: Bug fixes, performance improvements

## Contributing

### Adding New Features
1. Design public API with stability in mind
2. Implement with comprehensive tests
3. Add benchmarks for performance validation
4. Update documentation and examples
5. Ensure backward compatibility

### Performance Requirements
- All functions should complete in microseconds for typical inputs
- Zero or minimal memory allocations in hot paths
- Benchmark validation required for new functions
- Performance regression testing in CI

## Related Documentation

- [Geospatial API Documentation](geospatial/doc.go)
- [Architecture Overview](../docs/ARCHITECTURAL_DESCRIPTION.md)
- [Performance Benchmarks](../docs/PERFORMANCE.md)
- [Main README](../README.md)
