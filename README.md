# Global Traffic Simulator (GTS)

A Go-based traffic simulation system that generates realistic movement patterns for maritime, airborne, land, and space platforms. The simulator produces Cursor on Target (CoT) XML output and provides real-time visualization through a web interface.

## 🎯 Features

- **Multi-Domain Simulation**: Support for airborne, maritime, land, and space platforms
- **Realistic Physics**: Platform-specific movement models with authentic behaviors
- **Configuration-Driven**: YAML-based platform definitions and scenario configuration
- **CoT Output**: Standards-compliant Cursor on Target XML message generation
- **Real-time Visualization**: Web-based map interface with live track updates
- **Scalable Architecture**: Handles hundreds to thousands of concurrent platforms
- **Extensible Design**: Easy addition of new platform types and behaviors

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or later
- Modern web browser for visualization

### Installation

```bash
# Clone the repository
git clone https://github.com/rhino11/trafficsim.git
cd trafficsim

# Build the application
go build -o bin/simrunner cmd/simrunner/main.go

# Run with default configuration
./bin/simrunner -config data/config.yaml
```

### Basic Usage

```bash
# Run a specific scenario
./bin/simrunner -config data/config.yaml -scenario east_coast_demo

# Override simulation parameters
./bin/simrunner -config data/config.yaml -duration 1h -timescale 2.0

# Enable debug logging
./bin/simrunner -config data/config.yaml -loglevel debug
```

Access the web interface at `http://localhost:8080` to view real-time platform movements.

## 📋 Architecture Overview

### Core Components

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  Configuration  │───▶│ Simulation Core │◀──▶│  Data Models    │
│ (config.yaml)   │    │ (internal/sim)  │    │(internal/models)│
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │
                                │ Outputs CoT
                                ▼
┌─────────────────┐    ┌─────────────────┐
│  CoT Output     │───▶│ Network Output  │
│(internal/output)│    │  (HTTP/UDP)     │
└─────────────────┘    └─────────────────┘
          ▲
          │ Track Updates
          │
┌─────────────────┐    ┌─────────────────┐
│   Web Server    │◀──▶│   Web UI        │
│(internal/server)│    │   (web/)        │
└─────────────────┘    └─────────────────┘
```

### Platform Types

#### Airborne Platforms
- **Commercial Aviation**: Boeing 737-800, Airbus A320
- **Military Aviation**: F-16 Fighting Falcon, C-130 Hercules
- **Realistic Flight Models**: Altitude constraints, cruise speeds, fuel consumption

#### Maritime Platforms
- **Naval Vessels**: Arleigh Burke-class Destroyer
- **Commercial Shipping**: Container ships, oil tankers
- **Maritime Physics**: Draft considerations, port operations

#### Land Platforms
- **Military Vehicles**: M1A2 Abrams MBT, HMMWV
- **Civilian Vehicles**: Police patrol cars
- **Terrain Awareness**: Gradient limits, road following

#### Space Platforms
- **Satellites**: Starlink, GPS constellation
- **Space Stations**: ISS modules
- **Orbital Mechanics**: Realistic orbital periods and velocities

## ⚙️ Configuration

### Platform Type Definitions

Platform types are defined in `data/config.yaml` with comprehensive characteristics:

```yaml
platforms:
  airborne_types:
    boeing_737_800:
      name: "Boeing 737-800"
      class: "Boeing 737-800"
      type: "airborne"
      category: "commercial"
      max_speed: 257.0      # m/s (500 kts)
      cruise_speed: 230.0   # m/s (447 kts)
      max_altitude: 12500.0 # meters
      length: 39.5          # meters
      width: 35.8           # wingspan
      height: 12.5          # meters
      mass: 79010.0         # kg
      fuel_capacity: 26020  # liters
      range: 5665000        # meters
```

### Scenario Configuration

Define simulation scenarios with specific platform instances:

```yaml
scenarios:
  east_coast_demo:
    name: "East Coast Traffic Demo"
    duration: "30m"
    instances:
      - id: "AA1234"
        type_id: "boeing_737_800"
        name: "American 1234"
        callsign: "AAL1234"
        start_position:
          latitude: 40.7128
          longitude: -74.0060
          altitude: 10000
        destination:
          latitude: 25.7617
          longitude: -80.1918
          altitude: 11000
```

## 🏗️ Project Structure

```
trafficsim/
├── cmd/simrunner/          # Main application entry point
├── internal/
│   ├── config/             # Configuration management
│   ├── models/             # Platform data models
│   ├── sim/                # Simulation engine
│   ├── server/             # Web server
│   └── output/             # CoT output generation
├── pkg/geospatial/         # Geospatial utilities
├── web/                    # Web UI assets
├── data/                   # Configuration and sample data
├── docs/                   # Documentation and ADRs
└── scripts/                # Build and utility scripts
```

## 🔧 Development

### Building from Source

```bash
# Install dependencies
go mod download

# Run tests
go test ./...

# Build for multiple platforms
make build-all

# Run linting
make lint
```

### Adding New Platform Types

1. **Define Platform Characteristics**: Add to `data/config.yaml` under appropriate type category
2. **Implement Physics Model**: Add movement logic in `internal/models/`
3. **Update Factory**: Modify `internal/config/factory.go` for platform creation
4. **Add Scenarios**: Create test scenarios in configuration

### Architecture Decision Records (ADRs)

Key architectural decisions are documented in `docs/adr/`:
- **ADR-001**: Language Choice (Go)
- **ADR-002**: Simulation State Management
- **ADR-003**: Visualization Communication
- **ADR-004**: CoT Output Format
- **ADR-005**: Geospatial Data Handling
- **ADR-006**: Configuration Management
- **ADR-007**: Platform Modularity

## 📡 CoT Output

The simulator generates standards-compliant Cursor on Target (CoT) XML messages:

```xml
<event version="2.0" uid="AA1234" type="a-f-A-C" time="2025-06-03T14:30:00Z">
  <point lat="40.7128" lon="-74.0060" hae="10000" ce="50" le="25"/>
  <detail>
    <track course="090" speed="230"/>
    <contact callsign="AAL1234"/>
  </detail>
</event>
```

### Output Configuration

```yaml
output:
  cot:
    enabled: true
    endpoint: "udp://239.2.3.1:6969"  # Multicast UDP
    update_rate: "5s"
```

## 🌐 Web Interface

The integrated web interface provides:
- **Real-time Map Display**: Live platform positions and tracks
- **Platform Information**: Detailed platform characteristics
- **Scenario Control**: Start, stop, and configure simulations
- **Performance Metrics**: Track counts and update rates

Access at `http://localhost:8080` when running the simulator.

## 🔍 Physics Engine

### Movement Models

Each platform type implements realistic physics:

#### Aircraft
- **Altitude Management**: Climb/descent rates, service ceilings
- **Speed Control**: Stall speeds, cruise optimization
- **Navigation**: Great circle routes, waypoint following

#### Ships
- **Hydrodynamics**: Draft limitations, port constraints
- **Weather Effects**: Basic sea state considerations
- **Traffic Separation**: Shipping lane adherence

#### Land Vehicles
- **Terrain Following**: Gradient limitations, road networks
- **Fuel Consumption**: Range-based movement constraints
- **Obstacle Avoidance**: Basic pathfinding

#### Spacecraft
- **Orbital Mechanics**: Kepler's laws, orbital periods
- **Station Keeping**: Altitude maintenance
- **Ground Track**: Realistic satellite ground coverage

## 📊 Performance

- **Platform Capacity**: 1000+ concurrent platforms on typical hardware
- **Update Rate**: Configurable 1-60 second intervals
- **Time Scaling**: Real-time to 10x acceleration
- **Memory Usage**: ~1MB per 100 platforms

## 🤝 Contributing

1. **Fork the Repository**
2. **Create Feature Branch**: `git checkout -b feature/amazing-feature`
3. **Follow Go Standards**: Use `gofmt`, `golangci-lint`
4. **Add Tests**: Maintain test coverage
5. **Update Documentation**: Include ADRs for architectural changes
6. **Submit Pull Request**

### Code Style

- Follow standard Go conventions
- Use meaningful variable names
- Add comprehensive documentation
- Include unit tests for new features

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by real-world military simulation systems
- Built with Go's excellent concurrency primitives
- Uses open-source mapping libraries for visualization

## 📞 Support

- **Issues**: Report bugs via GitHub Issues
- **Documentation**: See `docs/` directory
- **Examples**: Check `data/sample_routes/` for scenario examples

---

**Version**: 0.1.0  
**Status**: Active Development  
**Go Version**: 1.21+
