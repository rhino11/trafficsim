/**
 * Comprehensive unit tests for MapEngine class
 */

// Import the MapEngine class
require('./map-engine.js');

describe('MapEngine', () => {
    let mapEngine;
    let mockMap;
    let mockPlatformLayer;
    let mockTrailLayer;
    let mockTileLayer;
    let containerElement;

    beforeEach(() => {
        // Create mock DOM container
        containerElement = document.createElement('div');
        containerElement.id = 'test-map';
        document.body.appendChild(containerElement);

        // Mock Leaflet objects
        mockTileLayer = {
            addTo: jest.fn()
        };

        mockPlatformLayer = {
            addTo: jest.fn(),
            addLayer: jest.fn(),
            removeLayer: jest.fn(),
            clearLayers: jest.fn()
        };

        mockTrailLayer = {
            addTo: jest.fn(),
            addLayer: jest.fn(),
            removeLayer: jest.fn(),
            clearLayers: jest.fn()
        };

        const mockBounds = {
            getNorth: () => 45,
            getSouth: () => 35,
            getEast: () => -70,
            getWest: () => -80
        };

        // Mock Leaflet map with all required methods
        mockMap = {
            addLayer: jest.fn(),
            removeLayer: jest.fn(),
            hasLayer: jest.fn(() => false),
            on: jest.fn(),
            off: jest.fn(),
            getBounds: jest.fn(() => mockBounds),
            getCenter: jest.fn(() => ({ lat: 39.8283, lng: -98.5795 })),
            getZoom: jest.fn(() => 10),
            setView: jest.fn(),
            setZoom: jest.fn(),
            panTo: jest.fn(),
            fitBounds: jest.fn(),
            remove: jest.fn(),
            invalidateSize: jest.fn(),
            openPopup: jest.fn(),
            closePopup: jest.fn(),
            eachLayer: jest.fn(),
            flyTo: jest.fn(),
            // Add missing coordinate conversion methods
            latLngToContainerPoint: jest.fn(() => ({ x: 400, y: 300 })),
            containerPointToLatLng: jest.fn(() => ({ lat: 40, lng: -75 })),
            project: jest.fn(() => ({ x: 100, y: 100 })),
            unproject: jest.fn(() => ({ lat: 40, lng: -75 })),
            getSize: jest.fn(() => ({ width: 800, height: 600 })),
            getContainer: jest.fn(() => containerElement)
        };

        // Mock Leaflet globally
        global.L = {
            map: jest.fn(() => mockMap),
            tileLayer: jest.fn(() => mockTileLayer),
            layerGroup: jest.fn(() => mockPlatformLayer),
            control: {
                layers: jest.fn(() => ({
                    addTo: jest.fn()
                }))
            }
        };

        // Mock console methods
        console.log = jest.fn();
        console.error = jest.fn();
        console.warn = jest.fn();
    });

    afterEach(() => {
        if (mapEngine) {
            mapEngine.destroy();
        }

        // Clean up DOM
        if (containerElement && containerElement.parentNode) {
            document.body.removeChild(containerElement);
        }

        jest.clearAllMocks();
    });

    describe('Constructor and Initialization', () => {
        it('should initialize with default options', async () => {
            mapEngine = new MapEngine('test-map');
            await mapEngine.initialize();

            expect(mapEngine.containerId).toBe('test-map');
            expect(mapEngine.map).toBe(mockMap);
            expect(mapEngine.platformLayer).toBe(mockPlatformLayer);
            expect(mapEngine.trailLayer).toBe(mockPlatformLayer); // Second call returns same mock
        });

        it('should initialize with custom options', async () => {
            const options = {
                center: [42.0, -76.0],
                zoom: 8,
                maxZoom: 20,
                attribution: 'Custom attribution'
            };

            mapEngine = new MapEngine('test-map', options);
            await mapEngine.initialize();

            expect(mockMap.setView).toHaveBeenCalledWith([42.0, -76.0], 8);
        });

        it('should create map with proper tile layer', async () => {
            mapEngine = new MapEngine('test-map');
            await mapEngine.initialize();

            expect(global.L.tileLayer).toHaveBeenCalledWith(
                expect.stringContaining('openstreetmap'),
                expect.objectContaining({
                    maxZoom: expect.any(Number),
                    attribution: expect.any(String)
                })
            );
            expect(mockTileLayer.addTo).toHaveBeenCalledWith(mockMap);
        });

        it('should create platform and trail layers', async () => {
            mapEngine = new MapEngine('test-map');
            await mapEngine.initialize();

            expect(global.L.layerGroup).toHaveBeenCalledTimes(2);
            expect(mockPlatformLayer.addTo).toHaveBeenCalledWith(mockMap);
        });

        it('should set up event handlers', async () => {
            mapEngine = new MapEngine('test-map');
            await mapEngine.initialize();

            expect(mockMap.on).toHaveBeenCalledWith('moveend', expect.any(Function));
            expect(mockMap.on).toHaveBeenCalledWith('zoomend', expect.any(Function));
            expect(mockMap.on).toHaveBeenCalledWith('resize', expect.any(Function));
        });

        it('should handle missing container element', async () => {
            const invalidMapEngine = new MapEngine('non-existent-container');

            await expect(invalidMapEngine.initialize()).rejects.toThrow('Container element not found: non-existent-container');
        });
    });

    describe('Map Controls and Navigation', () => {
        beforeEach(async () => {
            mapEngine = new MapEngine('test-map');
            await mapEngine.initialize();
        });

        it('should center map on coordinates', () => {
            mapEngine.centerOn(40.0, -75.0, 12);

            expect(mockMap.setView).toHaveBeenCalledWith([40.0, -75.0], 12);
        });

        it('should center map with default zoom if not provided', () => {
            mapEngine.centerOn(40.0, -75.0);

            expect(mockMap.setView).toHaveBeenCalledWith([40.0, -75.0], 10);
        });

        it('should smoothly fly to location', () => {
            mapEngine.flyTo(40.0, -75.0, 12);

            expect(mockMap.flyTo).toHaveBeenCalledWith([40.0, -75.0], 12);
        });

        it('should fit map to bounds', () => {
            const bounds = [[35, -80], [45, -70]];

            mapEngine.fitBounds(bounds);

            expect(mockMap.fitBounds).toHaveBeenCalledWith(bounds);
        });

        it('should set zoom level', () => {
            mapEngine.setZoom(15);

            expect(mockMap.setZoom).toHaveBeenCalledWith(15);
        });

        it('should get current zoom level', () => {
            const zoom = mapEngine.getZoom();

            expect(zoom).toBe(10);
            expect(mockMap.getZoom).toHaveBeenCalled();
        });
    });

    describe('Viewport and Bounds Management', () => {
        beforeEach(async () => {
            mapEngine = new MapEngine('test-map');
            await mapEngine.initialize();
        });

        it('should get current viewport bounds', () => {
            const bounds = mapEngine.getViewportBounds();

            expect(bounds).toEqual({
                north: 45,
                south: 35,
                east: -70,
                west: -80
            });
        });

        it('should check if coordinates are in viewport', () => {
            // Inside viewport
            expect(mapEngine.isInViewport(40, -75)).toBe(true);

            // Outside viewport
            expect(mapEngine.isInViewport(50, -75)).toBe(false);
            expect(mapEngine.isInViewport(40, -90)).toBe(false);
        });

        it('should handle viewport change events', () => {
            const callback = jest.fn();
            mapEngine.onViewportChangeCallback = callback;

            // Trigger moveend event
            const moveendHandler = mockMap.on.mock.calls.find(call => call[0] === 'moveend')[1];
            moveendHandler();

            expect(callback).toHaveBeenCalledWith({
                north: 45,
                south: 35,
                east: -70,
                west: -80
            });
        });
    });

    describe('Layer Management', () => {
        beforeEach(async () => {
            mapEngine = new MapEngine('test-map');
            await mapEngine.initialize();
        });

        it('should return platform layer', () => {
            const layer = mapEngine.getPlatformLayer();
            expect(layer).toBe(mockPlatformLayer);
        });

        it('should return trail layer', () => {
            const layer = mapEngine.getTrailLayer();
            expect(layer).toBe(mockTrailLayer);
        });

        it('should add custom layer to map', () => {
            const customLayer = { addTo: jest.fn() };

            mapEngine.addLayer(customLayer);

            expect(mockMap.addLayer).toHaveBeenCalledWith(customLayer);
        });

        it('should remove layer from map', () => {
            const customLayer = { removeFrom: jest.fn() };

            mapEngine.removeLayer(customLayer);

            expect(mockMap.removeLayer).toHaveBeenCalledWith(customLayer);
        });
    });

    describe('Coordinate Conversion', () => {
        beforeEach(async () => {
            mapEngine = new MapEngine('test-map');
            await mapEngine.initialize();
        });

        it('should convert lat/lng to screen coordinates', () => {
            const screenCoords = mapEngine.latLngToScreen(40, -75);

            expect(mockMap.latLngToContainerPoint).toHaveBeenCalledWith([40, -75]);
            expect(screenCoords).toEqual({ x: 400, y: 300 });
        });

        it('should convert screen coordinates to lat/lng', () => {
            const latLng = mapEngine.screenToLatLng(400, 300);

            expect(mockMap.containerPointToLatLng).toHaveBeenCalledWith([400, 300]);
            expect(latLng).toEqual({ lat: 40, lng: -75 });
        });

        it('should project coordinates', () => {
            const projected = mapEngine.project(40, -75);

            expect(mockMap.project).toHaveBeenCalledWith([40, -75]);
            expect(projected).toEqual({ x: 100, y: 100 });
        });

        it('should unproject coordinates', () => {
            const unprojected = mapEngine.unproject(100, 100);

            expect(mockMap.unproject).toHaveBeenCalledWith([100, 100]);
            expect(unprojected).toEqual({ lat: 40, lng: -75 });
        });
    });

    describe('Map State and Properties', () => {
        beforeEach(async () => {
            mapEngine = new MapEngine('test-map');
            await mapEngine.initialize();
        });

        it('should get map size', () => {
            const size = mapEngine.getMapSize();

            expect(size).toEqual({ width: 800, height: 600 });
            expect(mockMap.getSize).toHaveBeenCalled();
        });

        it('should get map container', () => {
            const container = mapEngine.getContainer();

            expect(container).toBe(containerElement);
            expect(mockMap.getContainer).toHaveBeenCalled();
        });

        it('should invalidate map size', () => {
            mapEngine.invalidateSize();

            expect(mockMap.invalidateSize).toHaveBeenCalled();
        });

        it('should check if map is ready', async () => {
            expect(mapEngine.isReady()).toBe(true);

            // Test before initialization
            const newMapEngine = new MapEngine('test-map');
            expect(newMapEngine.isReady()).toBe(false);
        });
    });

    describe('Event Handling', () => {
        beforeEach(async () => {
            mapEngine = new MapEngine('test-map');
            await mapEngine.initialize();
        });

        it('should add event listener', () => {
            const callback = jest.fn();

            mapEngine.on('click', callback);

            expect(mockMap.on).toHaveBeenCalledWith('click', callback);
        });

        it('should remove event listener', () => {
            const callback = jest.fn();

            mapEngine.off('click', callback);

            expect(mockMap.off).toHaveBeenCalledWith('click', callback);
        });

        it('should handle resize events', () => {
            // Trigger resize event
            const resizeHandler = mockMap.on.mock.calls.find(call => call[0] === 'resize')[1];
            resizeHandler();

            expect(mockMap.invalidateSize).toHaveBeenCalled();
        });
    });

    describe('Performance Optimization', () => {
        beforeEach(async () => {
            mapEngine = new MapEngine('test-map');
            await mapEngine.initialize();
        });

        it('should throttle viewport change events', () => {
            jest.useFakeTimers();
            const callback = jest.fn();
            mapEngine.onViewportChangeCallback = callback;

            // Trigger multiple moveend events quickly
            const moveendHandler = mockMap.on.mock.calls.find(call => call[0] === 'moveend')[1];
            moveendHandler();
            moveendHandler();
            moveendHandler();

            // Should only call once due to throttling
            jest.advanceTimersByTime(100);
            expect(callback).toHaveBeenCalledTimes(1);

            jest.useRealTimers();
        });

        it('should calculate viewport padding for culling', () => {
            const bounds = mapEngine.getViewportBounds();
            const paddedBounds = mapEngine.getViewportBoundsWithPadding(0.1);

            expect(paddedBounds.north).toBeGreaterThan(bounds.north);
            expect(paddedBounds.south).toBeLessThan(bounds.south);
            expect(paddedBounds.east).toBeGreaterThan(bounds.east);
            expect(paddedBounds.west).toBeLessThan(bounds.west);
        });
    });

    describe('Error Handling', () => {
        it('should handle map creation errors gracefully', async () => {
            global.L.map.mockImplementation(() => {
                throw new Error('Map creation failed');
            });

            const errorMapEngine = new MapEngine('test-map');

            await expect(errorMapEngine.initialize()).rejects.toThrow('Map creation failed');
        });

        it('should handle missing Leaflet library', async () => {
            const originalL = global.L;
            delete global.L;

            const errorMapEngine = new MapEngine('test-map');

            await expect(errorMapEngine.initialize()).rejects.toThrow('Leaflet library not found');

            global.L = originalL;
        });

        it('should validate coordinates before operations', () => {
            expect(() => mapEngine.centerOn(null, -75)).toThrow();
            expect(() => mapEngine.centerOn(40, null)).toThrow();
            expect(() => mapEngine.centerOn(91, -75)).toThrow(); // Invalid latitude
            expect(() => mapEngine.centerOn(40, 181)).toThrow(); // Invalid longitude
        });
    });

    describe('Configuration and Options', () => {
        it('should support different tile layer providers', async () => {
            const options = {
                tileLayer: 'https://custom-tiles.com/{z}/{x}/{y}.png',
                attribution: 'Custom tiles'
            };

            mapEngine = new MapEngine('test-map', options);
            await mapEngine.initialize();

            expect(global.L.tileLayer).toHaveBeenCalledWith(
                'https://custom-tiles.com/{z}/{x}/{y}.png',
                expect.objectContaining({
                    attribution: 'Custom tiles'
                })
            );
        });

        it('should support custom CRS', async () => {
            const customCRS = { code: 'EPSG:4326' };
            const options = {
                crs: customCRS
            };

            mapEngine = new MapEngine('test-map', options);
            await mapEngine.initialize();

            expect(global.L.map).toHaveBeenCalledWith('test-map', expect.objectContaining({
                crs: customCRS
            }));
        });

        it('should configure zoom constraints', async () => {
            const options = {
                minZoom: 2,
                maxZoom: 18
            };

            mapEngine = new MapEngine('test-map', options);
            await mapEngine.initialize();

            expect(global.L.map).toHaveBeenCalledWith('test-map', expect.objectContaining({
                minZoom: 2,
                maxZoom: 18
            }));
        });
    });

    describe('Cleanup and Destruction', () => {
        beforeEach(async () => {
            mapEngine = new MapEngine('test-map');
            await mapEngine.initialize();
        });

        it('should clean up resources on destroy', () => {
            const removeHandler = jest.fn();
            mockMap.remove = removeHandler;

            mapEngine.destroy();

            expect(removeHandler).toHaveBeenCalled();
            expect(mapEngine.map).toBeNull();
            expect(mapEngine.platformLayer).toBeNull();
            expect(mapEngine.trailLayer).toBeNull();
        });

        it('should remove event listeners on destroy', () => {
            mapEngine.destroy();

            expect(mockMap.off).toHaveBeenCalledWith('moveend');
            expect(mockMap.off).toHaveBeenCalledWith('zoomend');
            expect(mockMap.off).toHaveBeenCalledWith('resize');
        });
    });

    describe('Integration Features', () => {
        beforeEach(async () => {
            mapEngine = new MapEngine('test-map');
            await mapEngine.initialize();
        });

        it('should support custom map controls', async () => {
            const scaleControl = { addTo: jest.fn() };
            global.L.control.scale.mockReturnValue(scaleControl);

            mapEngine.addScaleControl();

            expect(global.L.control.scale).toHaveBeenCalled();
            expect(scaleControl.addTo).toHaveBeenCalledWith(mockMap);
        });

        it('should support layer control', async () => {
            const layerControl = { addTo: jest.fn() };
            global.L.control.layers.mockReturnValue(layerControl);

            const baseLayers = { 'OpenStreetMap': mockTileLayer };
            const overlayLayers = { 'Platforms': mockPlatformLayer };

            mapEngine.addLayerControl(baseLayers, overlayLayers);

            expect(global.L.control.layers).toHaveBeenCalledWith(baseLayers, overlayLayers);
            expect(layerControl.addTo).toHaveBeenCalledWith(mockMap);
        });

        it('should export map state', () => {
            const state = mapEngine.exportState();

            expect(state).toEqual({
                center: expect.any(Object),
                zoom: 10,
                bounds: expect.any(Object)
            });
        });

        it('should restore map state', () => {
            const state = {
                center: { lat: 42, lng: -76 },
                zoom: 12,
                bounds: [[35, -80], [45, -70]]
            };

            mapEngine.restoreState(state);

            expect(mockMap.setView).toHaveBeenCalledWith([42, -76], 12);
        });
    });
});
