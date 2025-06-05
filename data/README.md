# Configuration and Data

This directory contains all configuration files, platform definitions, and sample data used by the TrafficSim system.

## Directory Structure

### config.yaml
Main simulation configuration file.

**Purpose**: Primary configuration for simulation parameters, physics settings, and default platform definitions.

**Key Sections**:
- `simulation`: Timestep, duration, real-time settings
- `physics`: Earth radius, gravity, air density constants
- `platforms`: Default platform configurations
- `output`: CoT message and multicast settings
- `web`: Web server and visualization settings

### configs/
Alternative configuration scenarios.

```
configs/
├── emergency_response.yaml   # Emergency services simulation
├── global_traffic_demo.yaml  # Global traffic demonstration
└── military_exercise.yaml    # Military training scenario
```

**Purpose**: Predefined scenarios for different simulation use cases.

### platforms/
Platform type definitions organized by domain.

```
platforms/
├── airborne/
│   ├── commercial/           # Commercial aircraft
│   ├── military/            # Military aircraft
│   ├── general_aviation/    # Private aircraft
│   └── drones/              # UAV platforms
├── land/
│   ├── vehicles/            # Road vehicles
│   ├── military/            # Military ground vehicles
│   └── rail/                # Train systems
├── maritime/
│   ├── commercial/          # Cargo and passenger ships
│   ├── military/            # Naval vessels
│   ├── fishing/             # Fishing vessels
│   └── recreational/        # Pleasure craft
└── space/
    ├── satellites/          # Satellite platforms
    ├── space_stations/      # Crewed space platforms
    └── debris/              # Space debris objects
```

### sample_routes/
Predefined routes and flight paths.

**Purpose**: Sample route definitions for realistic platform movement patterns.

## Configuration Format

### Main Configuration Structure
```yaml
simulation:
  timestep: 1s              # Simulation update interval
  duration: 3600s           # Total simulation time
  realtime: true            # Real-time vs fast execution

physics:
  earth_radius: 6371000.0   # Earth radius in meters
  gravity: 9.81             # Gravitational acceleration
  air_density: 1.225        # Standard air density

platforms:
  - type: "airborne"
    count: 100
    routes: "data/sample_routes/commercial_flights.yaml"

output:
  cot:
    enabled: true
    endpoint: "http://localhost:8080/cot"
  multicast:
    enabled: false
    address: "239.2.3.1"
    port: 6969

web:
  port: 8080
  static_dir: "web/static"
  template_dir: "web/templates"
```

### Platform Definitions
```yaml
# Example: Commercial aircraft platform
platform_type: "airborne"
category: "commercial"
specifications:
  max_speed: 250.0          # m/s (cruise speed)
  max_altitude: 12000.0     # meters
  turn_rate: 3.0            # degrees per second
  climb_rate: 5.0           # m/s
  fuel_capacity: 50000.0    # kg

physics:
  mass: 75000.0             # kg
  drag_coefficient: 0.02
  wing_area: 120.0          # m²

navigation:
  gps_accuracy: 5.0         # meters
  update_rate: 1.0          # Hz

visualization:
  symbol: "aircraft"
  color: "#0066CC"
  size: 16
```

### Route Definitions
```yaml
# Example: Transcontinental route
route_name: "JFK_to_LAX"
description: "New York to Los Angeles flight path"
waypoints:
  - position: {lat: 40.6413, lon: -73.7781, alt: 0}      # JFK Airport
    action: "takeoff"
  - position: {lat: 40.7, lon: -74.0, alt: 10000}       # Departure climb
    action: "climb"
  - position: {lat: 39.0, lon: -95.0, alt: 11000}       # Cruise altitude
    action: "cruise"
  - position: {lat: 34.0, lon: -118.0, alt: 11000}      # Approach
    action: "descend"
  - position: {lat: 33.9425, lon: -118.4081, alt: 0}    # LAX Airport
    action: "landing"

flight_parameters:
  cruise_speed: 240.0       # m/s
  climb_rate: 8.0           # m/s
  descent_rate: -6.0        # m/s
```

## Configuration Usage

### Loading Configurations
```bash
# Use default configuration
./trafficsim

# Use specific configuration file
./trafficsim -config data/configs/military_exercise.yaml

# Validate configuration before use
./validate-yaml data/config.yaml
```

### Environment-Specific Configurations
- **Development**: `config.yaml` with debug settings
- **Testing**: Minimal platform counts for fast tests
- **Demo**: `global_traffic_demo.yaml` for presentations
- **Military**: `military_exercise.yaml` for defense scenarios

### Configuration Validation
```bash
# Validate single file
./validate-yaml data/config.yaml

# Validate all configurations
./validate-yaml data/configs/*.yaml

# Batch validation with summary
./validate-yaml data/configs/ -summary
```

## Platform Categories

### Airborne Platforms
- **Commercial**: Airlines, cargo aircraft
- **Military**: Fighters, transports, reconnaissance
- **General Aviation**: Private aircraft, training planes
- **UAV/Drones**: Military and civilian unmanned aircraft

### Maritime Platforms
- **Commercial**: Container ships, tankers, cruise ships
- **Military**: Destroyers, submarines, aircraft carriers
- **Fishing**: Commercial fishing vessels
- **Recreational**: Yachts, pleasure boats

### Land Platforms
- **Vehicles**: Cars, trucks, buses
- **Military**: Tanks, APCs, support vehicles
- **Rail**: Trains, light rail, subway systems

### Space Platforms
- **Satellites**: Communication, weather, GPS satellites
- **Space Stations**: ISS, commercial stations
- **Debris**: Tracked space debris objects

## Scenario Examples

### Emergency Response Scenario
```yaml
# emergency_response.yaml
scenario_name: "Emergency Response Coordination"
description: "Multi-agency emergency response simulation"

platforms:
  - type: "airborne"
    category: "emergency"
    count: 5
    aircraft_types: ["helicopter", "fixed_wing"]

  - type: "land"
    category: "emergency"
    count: 20
    vehicle_types: ["ambulance", "fire_truck", "police"]

areas_of_interest:
  - name: "incident_site"
    center: {lat: 40.7128, lon: -74.0060}
    radius: 5000  # meters
```

### Military Exercise Scenario
```yaml
# military_exercise.yaml
scenario_name: "Joint Training Exercise"
description: "Multi-domain military training simulation"

platforms:
  - type: "airborne"
    category: "military"
    count: 50
    mix: ["fighter", "transport", "reconnaissance"]

  - type: "maritime"
    category: "military"
    count: 10
    mix: ["destroyer", "frigate", "support"]

exercise_parameters:
  duration: 7200s  # 2 hours
  real_time: false
  acceleration: 4.0
```

## Data Management

### File Organization
- Use descriptive names for all configuration files
- Group related configurations in subdirectories
- Include version information in configuration comments
- Maintain backward compatibility when possible

### Version Control
- All configuration files are version controlled
- Use meaningful commit messages for configuration changes
- Tag stable configuration versions
- Document breaking changes in configuration format

### Backup and Recovery
- Regular backups of critical configuration files
- Configuration file validation in CI/CD pipeline
- Rollback procedures for configuration issues
- Default fallback configurations available

## Development Guidelines

### Adding New Configurations
1. Create configuration file with descriptive name
2. Include comprehensive documentation comments
3. Validate configuration format
4. Add to version control
5. Update this README

### Configuration Best Practices
- Use clear, descriptive parameter names
- Include units for all physical quantities
- Provide reasonable default values
- Document all configuration options
- Validate configuration on load

### Testing Configurations
- Unit tests for configuration loading
- Integration tests with different scenarios
- Performance testing with various platform counts
- Validation testing for malformed configurations

## Related Documentation

- [Configuration Management ADR](../docs/adr/006-configuration-management.md)
- [Platform Architecture](../docs/ARCHITECTURAL_DESCRIPTION.md)
- [Validation Tool Documentation](../cmd/README.md#validate-yaml)
- [Main README](../README.md)
