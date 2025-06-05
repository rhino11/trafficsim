# Military Symbology Implementation Plan

**Target Date:** June 6, 2025
**Reference:** ADR-008 Military Symbology Rendering

## Overview

Integrate mil-sym-ts library to replace generic platform markers with MIL-STD-2525D+ compliant military symbols in the TrafficSim web interface.

## Phase 1: Library Integration (Day 1 - High Priority)

### 1.1 Environment Setup (30 minutes)
- [ ] Install mil-sym-ts library: `npm install mil-sym-ts`
- [ ] Update package.json dependencies
- [ ] Verify library assets and manifest files
- [ ] Test basic library import and initialization

### 1.2 Core Integration (2 hours)
- [ ] **Map Engine Updates** (`web/static/js/map-engine.js`)
  - Initialize C5Ren renderer on map load
  - Add symbol renderer initialization method
  - Handle library asset loading and ready state

- [ ] **Symbol Utility Wrapper** (new file: `web/static/js/mil-symbol-utils.js`)
  - Create MilSymbolRenderer class
  - Implement platform data to MIL-STD type code mapping
  - Add symbol generation methods
  - Include error handling and fallbacks

### 1.3 Platform Renderer Enhancement (2 hours)
- [ ] **Update Platform Renderer** (`web/static/js/platform-renderer.js`)
  - Replace generic marker creation with military symbols
  - Integrate MilSymbolRenderer class
  - Maintain backward compatibility for non-military platforms
  - Add symbol caching mechanism

### 1.4 Data Pipeline Integration (1 hour)
- [ ] **Enhance Platform Data Mapping**
  - Ensure platform type data includes MIL-STD compatible fields
  - Map existing platform categories to symbology codes
  - Handle affiliation determination (friendly/hostile/neutral/unknown)
  - Add symbol modifier support (speed, heading, status)

## Phase 2: Symbol Enhancement (Day 2 - Medium Priority)

### 2.1 Advanced Symbol Features (2 hours)
- [ ] **Dynamic Symbol Properties**
  - Implement affiliation-based color coding
  - Add speed vector indicators
  - Include platform status modifiers
  - Support symbol scaling based on zoom level

### 2.2 Performance Optimization (1.5 hours)
- [ ] **Symbol Caching System**
  - Cache generated symbols by type and state
  - Implement lazy loading for better performance
  - Add symbol cleanup for removed platforms
  - Optimize rendering for large platform counts

### 2.3 User Experience Enhancements (1.5 hours)
- [ ] **Interactive Symbol Features**
  - Maintain click/hover functionality on symbols
  - Ensure symbol tooltips work correctly
  - Add symbol selection highlighting
  - Implement smooth symbol transitions

## Technical Implementation Details

### Symbol Generation Flow
```javascript
// Existing platform data
Platform Data →
// New mil-sym-ts integration
MIL-STD Type Code →
mil-sym-ts Symbol →
SVG Rendering →
Map Display
```

### Key Integration Points

1. **MilSymbolRenderer Class Structure:**
```javascript
class MilSymbolRenderer {
  constructor(renderer) // Initialize with C5Ren renderer
  generateSymbol(platformData) // Create symbol from platform
  updateSymbol(symbol, platformData) // Update existing symbol
  cacheSymbol(key, symbol) // Cache management
  mapPlatformToSymbolCode(platform) // Type code mapping
}
```

2. **Platform Renderer Integration:**
```javascript
// Replace existing marker creation
const symbol = this.milSymbolRenderer.generateSymbol(platform);
marker = this.createSymbolMarker(symbol, platform);
```

3. **Map Engine Initialization:**
```javascript
await C5Ren.initialize("/path/to/mil-sym-assets/");
this.milSymbolRenderer = new MilSymbolRenderer(C5Ren);
```

### File Modifications Required

1. **New Files:**
   - `web/static/js/mil-symbol-utils.js` - Symbol utility wrapper
   - Integration tests for symbol rendering

2. **Modified Files:**
   - `web/static/js/map-engine.js` - Add renderer initialization
   - `web/static/js/platform-renderer.js` - Replace markers with symbols
   - `package.json` - Add mil-sym-ts dependency
   - Test files - Update for new symbol rendering

3. **Asset Management:**
   - Ensure mil-sym-ts assets are properly served
   - Update build process if needed for library assets

## Testing Strategy

### Unit Tests
- [ ] Test symbol generation for each platform type
- [ ] Verify MIL-STD type code mapping
- [ ] Test symbol caching functionality
- [ ] Validate affiliation color coding

### Integration Tests
- [ ] Test complete symbol rendering pipeline
- [ ] Verify map integration works correctly
- [ ] Test performance with multiple platforms
- [ ] Ensure backward compatibility

### Visual Validation
- [ ] Verify symbols render correctly for all platform types
- [ ] Check affiliation colors match standards
- [ ] Validate symbol scaling and positioning
- [ ] Test interactive features (click, hover, selection)

## Success Criteria

### Day 1 Completion
- [x] mil-sym-ts library successfully integrated
- [x] Military symbols replace generic markers
- [x] Basic symbol rendering works for all platform types
- [x] No breaking changes to existing functionality

### Day 2 Completion
- [x] Advanced symbol features implemented
- [x] Performance optimizations in place
- [x] Enhanced user experience features working
- [x] All tests passing

## Risk Mitigation

### Technical Risks
- **Library Loading Issues**: Test initialization thoroughly, implement fallbacks
- **Performance Impact**: Monitor rendering times, implement caching
- **Asset Management**: Ensure proper asset serving configuration

### Implementation Risks
- **Breaking Changes**: Maintain backward compatibility, thorough testing
- **Symbol Mapping**: Validate MIL-STD type codes against actual standards
- **Browser Compatibility**: Test across target browsers

## Rollback Plan

If critical issues arise:
1. Feature flag to toggle between symbol types
2. Fallback to existing marker system
3. Library removal process documented
4. Emergency patch deployment procedure

## Future Enhancements (Beyond Day 2)

- Symbol animation for platform movement
- Tactical graphics support
- Enhanced modifier display options
- Symbol customization interface
- Export capabilities for military symbols
- Integration with additional military standards
