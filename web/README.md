# Web Interface

This directory contains all web-related assets for the TrafficSim visualization interface, including static files, templates, and frontend JavaScript components.

## Directory Structure

### static/
Static web assets (CSS, JavaScript, images).

```
static/
├── css/           # Stylesheets
├── js/            # JavaScript modules
├── images/        # Graphics and icons
└── fonts/         # Web fonts (if any)
```

### templates/
HTML templates for server-side rendering.

```
templates/
├── index.html              # Main simulation interface
├── scenario-builder.html   # Scenario configuration interface
└── layouts/               # Shared template layouts
```

## Key JavaScript Modules

### Core Modules

#### map-engine.js
**Purpose**: Leaflet map integration and management
- Map initialization and configuration
- Layer management (base maps, overlays)
- **Upcoming**: Military symbology integration with mil-sym-ts library
- Camera controls and view management

#### platform-renderer.js
**Purpose**: Platform visualization and rendering
- Platform marker creation and updates
- Real-time position updates via WebSocket
- **Upcoming**: Military symbol rendering (MIL-STD-2525D+)
- Platform selection and interaction

#### data-streamer.js
**Purpose**: WebSocket communication and data flow
- Real-time data streaming from simulation engine
- Platform update message handling
- Connection management and reconnection logic
- Data parsing and validation

#### performance-monitor.js
**Purpose**: Real-time performance metrics display
- FPS monitoring and display
- Memory usage tracking
- Platform count statistics
- Simulation health indicators

#### scenario-builder.js
**Purpose**: Interactive scenario configuration
- Platform creation and editing interface
- Route definition tools
- Configuration export/import
- Visual scenario preview

## Web Interface Features

### Real-time Visualization
- **Live Map Display**: Interactive map showing all platforms
- **Platform Tracking**: Real-time position updates
- **Multi-Platform Types**: Support for air, land, maritime, and space platforms
- **Performance Metrics**: Live FPS and system statistics

### Upcoming Military Symbology (June 2025)
- **MIL-STD-2525D+ Symbols**: Professional military symbology
- **Affiliation Display**: Clear friendly/hostile/neutral/unknown indicators  
- **Platform Recognition**: Immediate visual platform type identification
- **Standards Compliance**: Full military symbology standards support

### Interactive Controls
- **Simulation Control**: Start, stop, pause, reset simulation
- **Platform Management**: Add, remove, modify platforms
- **View Controls**: Zoom, pan, layer management
- **Configuration**: Real-time parameter adjustment

## Technology Stack

### Frontend Libraries
- **Leaflet**: Interactive map library
- **WebSocket API**: Real-time communication
- **mil-sym-ts** (upcoming): Military symbology rendering
- **Vanilla JavaScript**: No heavy frameworks, optimized performance

### Backend Integration
- **Go HTTP Server**: Serves static files and templates
- **WebSocket Server**: Real-time data streaming
- **Template Engine**: Server-side HTML rendering
- **API Endpoints**: RESTful platform management

## Development Workflow

### Local Development
```bash
# Start web server mode
./trafficsim -web -port 8080

# Access interface
open http://localhost:8080

# Development tools
npm run test:watch    # JavaScript tests
npm run test:coverage # Coverage reports
```

### File Organization
- Keep JavaScript modules focused and single-purpose
- Use consistent naming conventions
- Implement proper error handling
- Add comprehensive JSDoc documentation

### Testing
```bash
# Run JavaScript tests
npm test

# Run with coverage
npm run test:coverage

# Run specific test file
npm run test -- data-streamer.test.js
```

## Current Test Coverage

JavaScript test coverage is comprehensive across all modules:
- **map-engine.js**: Full coverage
- **platform-renderer.js**: Full coverage  
- **data-streamer.js**: Full coverage
- **performance-monitor.js**: Full coverage
- **scenario-builder.js**: Full coverage

## API Integration

### WebSocket Events
```javascript
// Platform updates
{
  "type": "platform_update",
  "platform": {
    "id": "NAVY-89",
    "lat": 36.8485,
    "lon": -76.2951,
    "alt": 0,
    "heading": 197.8,
    "speed": 0.3
  }
}

// Performance metrics
{
  "type": "simulation_metrics",
  "metrics": {
    "fps": 60,
    "platforms": 4,
    "memory_mb": 45.2
  }
}
```

### REST Endpoints
```javascript
// Platform management
GET    /api/platforms          // List all platforms
POST   /api/platforms          // Create platform
PUT    /api/platforms/{id}     // Update platform
DELETE /api/platforms/{id}     // Remove platform

// Simulation control
GET    /api/simulation/status  // Get simulation state
POST   /api/simulation/start   // Start simulation
POST   /api/simulation/stop    // Stop simulation
```

## Styling Guidelines

### CSS Organization
- Use semantic class names
- Implement responsive design principles
- Follow BEM methodology where appropriate
- Maintain consistent color schemes and typography

### Color Scheme
- **Primary**: Blues and grays for professional appearance
- **Status Colors**: Green (active), Yellow (warning), Red (error)
- **Military Colors**: Standard military color codes for symbology

## Performance Considerations

### Optimization Strategies
- **Symbol Caching**: Cache military symbols for performance
- **Efficient Updates**: Only update changed platforms
- **Memory Management**: Clean up removed platforms
- **Batch Processing**: Group platform updates

### Monitoring
- Real-time FPS monitoring
- Memory usage tracking
- WebSocket connection health
- Platform rendering performance

## Future Enhancements

### Planned Features (Q3-Q4 2025)
- **Enhanced Military Symbology**: Advanced symbol modifiers
- **Tactical Graphics**: Support for military overlay graphics
- **Symbol Animation**: Smooth movement transitions
- **Export Capabilities**: Screenshot and data export
- **Mobile Optimization**: Responsive design improvements

## Related Documentation

- [Military Symbology Implementation Plan](../docs/MIL_SYMBOL_IMPLEMENTATION_PLAN.md)
- [ADR-008: Military Symbology Rendering](../docs/adr/008-military-symbology-rendering.md)
- [Architecture Overview](../docs/ARCHITECTURAL_DESCRIPTION.md)
- [Main README](../README.md)
