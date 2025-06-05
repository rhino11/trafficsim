# ADR-008: Military Symbology Rendering with mil-sym-ts

**Date:** 2025-06-05
**Status:** Accepted
**Deciders:** Development Team

## Context

TrafficSim currently generates CoT (Cursor on Target) messages with MIL-STD-2525 type codes for military symbology, but the web visualization interface lacks proper military symbol rendering. Users see generic markers instead of standardized military symbols, which reduces situational awareness and professional presentation for military/defense use cases.

The current visualization system uses:
- Simple colored circles/markers for all platforms
- Basic text labels for identification
- No standardized military symbology
- Limited visual distinction between platform types and affiliations

For a military traffic simulation system, proper symbology is essential for:
- Immediate visual recognition of platform types
- Clear affiliation identification (friendly, hostile, neutral, unknown)
- Professional military presentation standards
- Integration with existing military C2 systems

## Decision

We will integrate the **mil-sym-ts** library (https://github.com/missioncommand/mil-sym-ts) for rendering MIL-STD-2525D+ compliant military symbols in the browser interface.

### Key reasons for choosing mil-sym-ts:

1. **Standards Compliance**: Implements MIL-STD-2525D+, the latest military symbology standard
2. **Perfect Technology Match**: Written in TypeScript (96.9%), aligns with our frontend stack
3. **Open Source**: Apache 2.0 license, compatible with our project
4. **Active Development**: Latest release v2.1.8 (May 2025), actively maintained
5. **Browser Optimized**: Designed specifically for web/JavaScript environments
6. **Comprehensive API**: Includes all necessary classes (MilStdSymbol, WebRenderer, SymbolUtilities)
7. **Existing Integration**: We already generate compatible MIL-STD type codes in our CoT messages

## Implementation Strategy

### Phase 1: Library Integration (Day 1)
- Install mil-sym-ts via npm
- Initialize renderer in map engine
- Create symbol utility wrapper class
- Update platform renderer to use military symbols

### Phase 2: Symbol Enhancement (Day 2)
- Implement dynamic symbol generation based on platform data
- Add affiliation color coding
- Integrate modifiers (speed vectors, status indicators)
- Add symbol caching for performance

### Phase 3: Advanced Features (Future)
- Symbol animation for movement
- Tactical graphics support
- Enhanced modifier display
- Symbol selection/interaction

## Technical Implementation

### Integration Points:
- **Platform Renderer**: Replace generic markers with mil-sym-ts symbols
- **Map Engine**: Initialize C5Ren renderer and symbol utilities
- **Data Pipeline**: Map our platform data to MIL-STD attributes
- **Performance**: Implement symbol caching and lazy loading

### Symbol Generation Flow:
```
Platform Data → MIL-STD Type Code → mil-sym-ts Symbol → SVG Rendering
```

### Dependencies:
- mil-sym-ts library (latest version)
- Update build process for library assets
- Enhance platform-renderer.js with symbol utilities

## Consequences

### Positive:
- **Professional Presentation**: Military-grade symbol rendering
- **Standards Compliance**: Aligns with defense industry expectations
- **Enhanced UX**: Immediate visual platform identification
- **Integration Ready**: Compatible with military C2 systems
- **Future Proof**: Based on latest MIL-STD-2525D+ standard
- **Maintainable**: Active open-source project with community support

### Considerations:
- **Bundle Size**: Additional library increases frontend payload
- **Learning Curve**: Team needs to understand MIL-STD symbology concepts
- **Asset Management**: Symbol assets need proper caching strategy
- **Performance**: SVG rendering performance with many symbols
- **Dependency**: External dependency for critical visualization feature

### Risks Mitigated:
- **Compatibility**: Apache 2.0 license ensures no legal issues
- **Maintenance**: Active project with recent updates reduces abandonment risk
- **Standards**: Using official MIL-STD ensures long-term compatibility

## Alternatives Considered

1. **Custom Symbol Implementation**: Too complex, reinventing the wheel
2. **Other Military Symbol Libraries**: mil-sym-ts is the most current and TypeScript-native
3. **Static Symbol Assets**: Limited flexibility, maintenance overhead
4. **Third-party Services**: Introduces external dependencies and potential costs

## Success Metrics

- Military symbols render correctly for all platform types
- Visual distinction between affiliations is clear
- Symbol rendering performance remains acceptable (< 100ms per symbol)
- Integration with existing CoT pipeline is seamless
- User feedback indicates improved situational awareness

## Related ADRs

- ADR-004: CoT Output (provides the MIL-STD type codes this library will render)
- ADR-003: Visualization Communication (defines the data flow this enhances)
- ADR-005: Geospatial Data Handling (coordinates with map display)
