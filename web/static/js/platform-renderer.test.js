/**
 * Unit tests for PlatformRenderer class
 */

// Import the PlatformRenderer class
require('./platform-renderer.js');

describe('PlatformRenderer', () => {
    let platformRenderer;
    let mockMapEngine;
    let mockMap;
    let mockPlatformLayer;
    let mockTrailLayer;
    let mockCanvasRenderer;
    let mockMarkerCluster;

    beforeEach(() => {
        // Mock Leaflet objects
        global.L = {
            Canvas: jest.fn(),
            canvas: jest.fn(() => mockCanvasRenderer),
            circleMarker: jest.fn(),
            marker: jest.fn(),
            divIcon: jest.fn(),
            polyline: jest.fn(),
            markerClusterGroup: jest.fn(() => mockMarkerCluster),
            DivIcon: jest.fn(),
            Point: jest.fn()
        };

        // Mock canvas renderer
        mockCanvasRenderer = {
            addTo: jest.fn()
        };

        // Mock marker cluster
        mockMarkerCluster = {
            addLayer: jest.fn(),
            removeLayer: jest.fn(),
            clearLayers: jest.fn(),
            getChildCount: jest.fn(() => 5)
        };

        // Mock Leaflet layers
        mockPlatformLayer = {
            addLayer: jest.fn(),
            removeLayer: jest.fn(),
            clearLayers: jest.fn()
        };

        mockTrailLayer = {
            addLayer: jest.fn(),
            removeLayer: jest.fn(),
            clearLayers: jest.fn()
        };

        // Mock map
        mockMap = {
            addLayer: jest.fn(),
            removeLayer: jest.fn(),
            getZoom: jest.fn(() => 10),
            on: jest.fn(),
            off: jest.fn()
        };

        // Mock map engine
        mockMapEngine = {
            getMap: jest.fn(() => mockMap),
            getPlatformLayer: jest.fn(() => mockPlatformLayer),
            getTrailLayer: jest.fn(() => mockTrailLayer),
            isInViewport: jest.fn(() => true),
            centerOn: jest.fn(),
            onViewportChangeCallback: null
        };

        // Reset console mocks
        console.log.mockClear();
        console.error.mockClear();
        console.warn.mockClear();

        // Mock performance
        global.performance = {
            now: jest.fn(() => Date.now())
        };
    });

    afterEach(() => {
        if (platformRenderer) {
            platformRenderer.clearAllPlatforms();
        }
    });

    describe('Constructor and Initialization', () => {
        it('should initialize with default settings', () => {
            platformRenderer = new PlatformRenderer(mockMapEngine);

            expect(platformRenderer.mapEngine).toBe(mockMapEngine);
            expect(platformRenderer.map).toBe(mockMap);
            expect(platformRenderer.platformLayer).toBe(mockPlatformLayer);
            expect(platformRenderer.trailLayer).toBe(mockTrailLayer);
            expect(platformRenderer.platforms).toBeInstanceOf(Map);
            expect(platformRenderer.markers).toBeInstanceOf(Map);
            expect(platformRenderer.trails).toBeInstanceOf(Map);
            expect(platformRenderer.useCanvas).toBe(true);
            expect(platformRenderer.maxTrailLength).toBe(20);
            expect(platformRenderer.showTrails).toBe(true);
        });

        it('should set up canvas renderer when available', () => {
            platformRenderer = new PlatformRenderer(mockMapEngine);

            expect(global.L.canvas).toHaveBeenCalledWith({
                padding: 0.5,
                tolerance: 10,
                updateWhenIdle: false,
                updateWhenZooming: true
            });
            expect(mockMap.addLayer).toHaveBeenCalledWith(mockCanvasRenderer);
        });

        it('should set up marker clustering when available', () => {
            platformRenderer = new PlatformRenderer(mockMapEngine);

            expect(global.L.markerClusterGroup).toHaveBeenCalledWith(
                expect.objectContaining({
                    maxClusterRadius: 50,
                    spiderfyOnMaxZoom: false,
                    showCoverageOnHover: false,
                    zoomToBoundsOnClick: false,
                    disableClusteringAtZoom: 10
                })
            );
        });

        it('should set up event handlers', () => {
            platformRenderer = new PlatformRenderer(mockMapEngine);

            expect(mockMap.on).toHaveBeenCalledWith('zoomend', expect.any(Function));
            expect(mockMapEngine.onViewportChangeCallback).toBeDefined();
        });
    });

    describe('Platform Marker Creation', () => {
        beforeEach(() => {
            platformRenderer = new PlatformRenderer(mockMapEngine);
        });

        it('should create canvas-based marker when canvas is enabled', () => {
            const platform = {
                id: 'test-1',
                platform_type: 'airborne',
                position: { latitude: 40.0, longitude: -75.0 }
            };

            const mockMarker = {
                platformData: null,
                bindPopup: jest.fn(),
                on: jest.fn()
            };
            global.L.circleMarker.mockReturnValue(mockMarker);

            const marker = platformRenderer.createPlatformMarker(platform);

            expect(global.L.circleMarker).toHaveBeenCalledWith(
                [40.0, -75.0],
                expect.objectContaining({
                    radius: expect.any(Number),
                    fillColor: '#2196F3', // Blue for airborne
                    renderer: mockCanvasRenderer
                })
            );
            expect(marker.platformData).toBe(platform);
        });

        it('should create div icon marker when canvas is disabled', () => {
            platformRenderer.useCanvas = false;
            platformRenderer.canvasRenderer = null;

            const platform = {
                id: 'test-1',
                platform_type: 'maritime',
                position: { latitude: 40.0, longitude: -75.0 }
            };

            const mockIcon = {};
            const mockMarker = {
                platformData: null,
                bindPopup: jest.fn(),
                on: jest.fn()
            };

            global.L.divIcon.mockReturnValue(mockIcon);
            global.L.marker.mockReturnValue(mockMarker);

            const marker = platformRenderer.createPlatformMarker(platform);

            expect(global.L.divIcon).toHaveBeenCalled();
            expect(global.L.marker).toHaveBeenCalledWith([40.0, -75.0], { icon: mockIcon, bubblingMouseEvents: false });
        });

        it('should calculate correct marker size based on platform type', () => {
            const testCases = [
                { type: 'space', expectedMultiplier: 0.7 },
                { type: 'airborne', expectedMultiplier: 1.0 },
                { type: 'maritime', expectedMultiplier: 1.2 },
                { type: 'land', expectedMultiplier: 0.9 }
            ];

            testCases.forEach(({ type, expectedMultiplier }) => {
                const platform = { platform_type: type };
                const size = platformRenderer.calculateMarkerSize(platform);
                expect(size).toBe(platformRenderer.markerSize * expectedMultiplier);
            });
        });

        it('should scale marker size based on zoom level', () => {
            platformRenderer.zoomBasedSizing = true;
            mockMap.getZoom.mockReturnValue(5);

            const platform = { platform_type: 'airborne' };
            const size = platformRenderer.calculateMarkerSize(platform);

            expect(size).toBe(Math.max(3, Math.min(20, platformRenderer.markerSize * (5 / 10))));
        });

        it('should return correct colors for platform types', () => {
            const expectedColors = {
                airborne: '#2196F3',
                maritime: '#00BCD4',
                land: '#4CAF50',
                space: '#9C27B0'
            };

            Object.entries(expectedColors).forEach(([type, color]) => {
                expect(platformRenderer.getPlatformColor(type)).toBe(color);
            });

            expect(platformRenderer.getPlatformColor('unknown')).toBe('#757575');
        });
    });

    describe('Platform Updates', () => {
        beforeEach(() => {
            platformRenderer = new PlatformRenderer(mockMapEngine);
        });

        it('should add new platform to map', () => {
            const platform = {
                id: 'test-1',
                platform_type: 'airborne',
                position: { latitude: 40.0, longitude: -75.0 },
                velocity: { north: 10, east: 5, up: 0 },
                speed: 15,
                heading: 90,
                lastUpdated: Date.now()
            };

            const mockMarker = {
                platformData: null,
                bindPopup: jest.fn(),
                on: jest.fn()
            };
            global.L.circleMarker.mockReturnValue(mockMarker);

            platformRenderer.updatePlatform(platform);

            expect(platformRenderer.platforms.has('test-1')).toBe(true);
            expect(platformRenderer.markers.has('test-1')).toBe(true);
            expect(mockPlatformLayer.addLayer).toHaveBeenCalledWith(mockMarker);
        });

        it('should update existing platform position', () => {
            const platform = {
                id: 'test-1',
                platform_type: 'airborne',
                position: { latitude: 40.0, longitude: -75.0 },
                velocity: { north: 10, east: 5, up: 0 },
                speed: 15,
                heading: 90,
                lastUpdated: Date.now()
            };

            const mockMarker = {
                platformData: null,
                bindPopup: jest.fn(),
                on: jest.fn(),
                setLatLng: jest.fn(),
                isPopupOpen: jest.fn(() => false)
            };
            global.L.circleMarker.mockReturnValue(mockMarker);

            // Add platform first
            platformRenderer.updatePlatform(platform);

            // Update platform position
            platform.position.latitude = 41.0;
            platformRenderer.updatePlatform(platform);

            expect(mockMarker.setLatLng).toHaveBeenCalledWith([41.0, -75.0]);
            expect(mockMarker.platformData).toBe(platform);
        });

        it('should remove platform when not visible', () => {
            const platform = {
                id: 'test-1',
                platform_type: 'airborne',
                position: { latitude: 40.0, longitude: -75.0 }
            };

            // Add platform first
            platformRenderer.updatePlatform(platform);

            // Make platform invisible
            platformRenderer.visibilityFilters.airborne = false;
            platformRenderer.updatePlatform(platform);

            expect(platformRenderer.platforms.has('test-1')).toBe(false);
            expect(platformRenderer.markers.has('test-1')).toBe(false);
        });

        it('should handle viewport culling', () => {
            mockMapEngine.isInViewport.mockReturnValue(false);

            const platform = {
                id: 'test-1',
                platform_type: 'airborne',
                position: { latitude: 40.0, longitude: -75.0 }
            };

            const mockMarker = {
                platformData: null,
                bindPopup: jest.fn(),
                on: jest.fn()
            };
            global.L.circleMarker.mockReturnValue(mockMarker);

            // Add platform when in viewport
            mockMapEngine.isInViewport.mockReturnValue(true);
            platformRenderer.updatePlatform(platform);
            expect(mockPlatformLayer.addLayer).toHaveBeenCalled();

            // Remove from rendering when out of viewport
            mockMapEngine.isInViewport.mockReturnValue(false);
            platformRenderer.updatePlatform(platform);
            expect(mockPlatformLayer.removeLayer).toHaveBeenCalled();
        });
    });

    describe('Platform Trails', () => {
        beforeEach(() => {
            platformRenderer = new PlatformRenderer(mockMapEngine);
            platformRenderer.showTrails = true;
        });

        it('should create trail for platform', () => {
            const platform = {
                id: 'test-1',
                platform_type: 'airborne',
                position: { latitude: 40.0, longitude: -75.0 }
            };

            const mockTrail = {
                setLatLngs: jest.fn()
            };
            global.L.polyline.mockReturnValue(mockTrail);

            // Add two positions to create a trail
            platformRenderer.updatePlatformTrail(platform);
            platform.position.latitude = 40.1;
            platformRenderer.updatePlatformTrail(platform);

            expect(global.L.polyline).toHaveBeenCalledWith(
                [[40.0, -75.0], [40.1, -75.0]],
                expect.objectContaining({
                    color: '#2196F3',
                    weight: 2,
                    opacity: 0.6
                })
            );
            expect(mockTrailLayer.addLayer).toHaveBeenCalledWith(mockTrail);
        });

        it('should limit trail length', () => {
            platformRenderer.maxTrailLength = 3;

            const platform = {
                id: 'test-1',
                platform_type: 'airborne',
                position: { latitude: 40.0, longitude: -75.0 }
            };

            // Add more points than max trail length
            for (let i = 0; i < 5; i++) {
                platform.position.latitude = 40.0 + i * 0.1;
                platformRenderer.updatePlatformTrail(platform);
            }

            const trailPoints = platformRenderer.trailPoints.get('test-1');
            expect(trailPoints.length).toBe(3);
            expect(trailPoints[0]).toEqual([40.2, -75.0]); // Should keep last 3 points
        });

        it('should update existing trail', () => {
            const platform = {
                id: 'test-1',
                platform_type: 'airborne',
                position: { latitude: 40.0, longitude: -75.0 }
            };

            const mockTrail = {
                setLatLngs: jest.fn()
            };
            global.L.polyline.mockReturnValue(mockTrail);

            // Create trail
            platformRenderer.updatePlatformTrail(platform);
            platform.position.latitude = 40.1;
            platformRenderer.updatePlatformTrail(platform);

            // Update trail
            platform.position.latitude = 40.2;
            platformRenderer.updatePlatformTrail(platform);

            expect(mockTrail.setLatLngs).toHaveBeenCalledWith([
                [40.0, -75.0],
                [40.1, -75.0],
                [40.2, -75.0]
            ]);
        });
    });

    describe('Batch Operations', () => {
        beforeEach(() => {
            platformRenderer = new PlatformRenderer(mockMapEngine);
        });

        it('should update multiple platforms efficiently', () => {
            const platforms = [
                {
                    id: 'test-1',
                    platform_type: 'airborne',
                    position: { latitude: 40.0, longitude: -75.0 }
                },
                {
                    id: 'test-2',
                    platform_type: 'maritime',
                    position: { latitude: 41.0, longitude: -76.0 }
                }
            ];

            const mockMarker = {
                platformData: null,
                bindPopup: jest.fn(),
                on: jest.fn()
            };
            global.L.circleMarker.mockReturnValue(mockMarker);

            platformRenderer.updatePlatforms(platforms);

            expect(platformRenderer.platforms.size).toBe(2);
            expect(platformRenderer.markers.size).toBe(2);
            expect(platformRenderer.renderStats.renderTime).toBeGreaterThan(0);
        });

        it('should remove platforms not in update batch', () => {
            // Add initial platforms
            const initialPlatforms = [
                { id: 'test-1', platform_type: 'airborne', position: { latitude: 40.0, longitude: -75.0 } },
                { id: 'test-2', platform_type: 'maritime', position: { latitude: 41.0, longitude: -76.0 } }
            ];

            const mockMarker = {
                platformData: null,
                bindPopup: jest.fn(),
                on: jest.fn()
            };
            global.L.circleMarker.mockReturnValue(mockMarker);

            platformRenderer.updatePlatforms(initialPlatforms);
            expect(platformRenderer.platforms.size).toBe(2);

            // Update with only one platform
            const updatedPlatforms = [
                { id: 'test-1', platform_type: 'airborne', position: { latitude: 40.1, longitude: -75.1 } }
            ];

            platformRenderer.updatePlatforms(updatedPlatforms);
            expect(platformRenderer.platforms.size).toBe(1);
            expect(platformRenderer.platforms.has('test-1')).toBe(true);
            expect(platformRenderer.platforms.has('test-2')).toBe(false);
        });
    });

    describe('Clustering', () => {
        beforeEach(() => {
            platformRenderer = new PlatformRenderer(mockMapEngine);
        });

        it('should enable clustering', () => {
            const mockMarker = {
                platformData: null,
                bindPopup: jest.fn(),
                on: jest.fn()
            };

            // Add some markers first
            platformRenderer.markers.set('test-1', mockMarker);

            platformRenderer.enableClustering();

            expect(platformRenderer.clusteringEnabled).toBe(true);
            expect(mockPlatformLayer.removeLayer).toHaveBeenCalledWith(mockMarker);
            expect(mockMarkerCluster.addLayer).toHaveBeenCalledWith(mockMarker);
            expect(mockMap.addLayer).toHaveBeenCalledWith(mockMarkerCluster);
        });

        it('should disable clustering', () => {
            const mockMarker = {
                platformData: null,
                bindPopup: jest.fn(),
                on: jest.fn()
            };

            // Enable clustering first
            platformRenderer.clusteringEnabled = true;
            platformRenderer.markers.set('test-1', mockMarker);

            platformRenderer.disableClustering();

            expect(platformRenderer.clusteringEnabled).toBe(false);
            expect(mockMarkerCluster.removeLayer).toHaveBeenCalledWith(mockMarker);
            expect(mockPlatformLayer.addLayer).toHaveBeenCalledWith(mockMarker);
            expect(mockMap.removeLayer).toHaveBeenCalledWith(mockMarkerCluster);
        });

        it('should auto-enable clustering for high platform counts', () => {
            const bounds = { north: 45, south: 35, east: -70, west: -80 };

            // Mock high visible platform count
            platformRenderer.countVisiblePlatforms = jest.fn(() => 1500);

            platformRenderer.onViewportChange(bounds);

            expect(platformRenderer.clusteringEnabled).toBe(true);
        });

        it('should auto-disable clustering for low platform counts', () => {
            platformRenderer.clusteringEnabled = true;
            const bounds = { north: 45, south: 35, east: -70, west: -80 };

            // Mock low visible platform count
            platformRenderer.countVisiblePlatforms = jest.fn(() => 300);

            platformRenderer.onViewportChange(bounds);

            expect(platformRenderer.clusteringEnabled).toBe(false);
        });
    });

    describe('Visibility Filters', () => {
        beforeEach(() => {
            platformRenderer = new PlatformRenderer(mockMapEngine);
        });

        it('should filter platforms by type', () => {
            const airborne = { id: 'air-1', platform_type: 'airborne', position: { latitude: 40.0, longitude: -75.0 } };
            const maritime = { id: 'sea-1', platform_type: 'maritime', position: { latitude: 41.0, longitude: -76.0 } };

            expect(platformRenderer.isPlatformVisible(airborne)).toBe(true);
            expect(platformRenderer.isPlatformVisible(maritime)).toBe(true);

            platformRenderer.setPlatformFilter('airborne', false);

            expect(platformRenderer.isPlatformVisible(airborne)).toBe(false);
            expect(platformRenderer.isPlatformVisible(maritime)).toBe(true);
        });

        it('should update existing platforms when filter changes', () => {
            const platform = {
                id: 'test-1',
                platform_type: 'airborne',
                position: { latitude: 40.0, longitude: -75.0 }
            };

            const mockMarker = {
                platformData: null,
                bindPopup: jest.fn(),
                on: jest.fn()
            };
            global.L.circleMarker.mockReturnValue(mockMarker);

            // Add platform
            platformRenderer.updatePlatform(platform);
            expect(platformRenderer.platforms.has('test-1')).toBe(true);

            // Hide airborne platforms
            platformRenderer.setPlatformFilter('airborne', false);
            expect(platformRenderer.platforms.has('test-1')).toBe(false);

            // Show airborne platforms again
            platformRenderer.setPlatformFilter('airborne', true);
            expect(platformRenderer.platforms.has('test-1')).toBe(true);
        });
    });

    describe('Trail Management', () => {
        beforeEach(() => {
            platformRenderer = new PlatformRenderer(mockMapEngine);
        });

        it('should show/hide trails', () => {
            platformRenderer.setTrailsVisible(true);
            expect(mockMap.addLayer).toHaveBeenCalledWith(mockTrailLayer);

            platformRenderer.setTrailsVisible(false);
            expect(mockMap.removeLayer).toHaveBeenCalledWith(mockTrailLayer);
        });

        it('should set trail length', () => {
            const platform = {
                id: 'test-1',
                platform_type: 'airborne',
                position: { latitude: 40.0, longitude: -75.0 }
            };

            const mockTrail = {
                setLatLngs: jest.fn()
            };

            // Create long trail
            const longTrail = [];
            for (let i = 0; i < 10; i++) {
                longTrail.push([40.0 + i * 0.1, -75.0]);
            }
            platformRenderer.trailPoints.set('test-1', longTrail);
            platformRenderer.trails.set('test-1', mockTrail);

            // Set shorter trail length
            platformRenderer.setTrailLength(5);

            expect(platformRenderer.maxTrailLength).toBe(5);
            expect(mockTrail.setLatLngs).toHaveBeenCalledWith(
                longTrail.slice(-5)
            );
        });
    });

    describe('Platform Focus and Navigation', () => {
        beforeEach(() => {
            platformRenderer = new PlatformRenderer(mockMapEngine);
        });

        it('should focus on specific platform', () => {
            const platform = {
                id: 'test-1',
                platform_type: 'airborne',
                position: { latitude: 40.0, longitude: -75.0 }
            };

            const mockMarker = {
                platformData: null,
                bindPopup: jest.fn(),
                on: jest.fn(),
                openPopup: jest.fn()
            };

            // Add platform
            platformRenderer.platforms.set('test-1', platform);
            platformRenderer.markers.set('test-1', mockMarker);

            platformRenderer.focusOnPlatform('test-1');

            expect(mockMapEngine.centerOn).toHaveBeenCalledWith(40.0, -75.0, 12);
            expect(mockMarker.openPopup).toHaveBeenCalled();
        });

        it('should handle focus on non-existent platform', () => {
            platformRenderer.focusOnPlatform('non-existent');

            expect(mockMapEngine.centerOn).not.toHaveBeenCalled();
        });
    });

    describe('Performance and Statistics', () => {
        beforeEach(() => {
            platformRenderer = new PlatformRenderer(mockMapEngine);
        });

        it('should track render statistics', () => {
            const platforms = [
                { id: 'test-1', platform_type: 'airborne', position: { latitude: 40.0, longitude: -75.0 } }
            ];

            const mockMarker = {
                platformData: null,
                bindPopup: jest.fn(),
                on: jest.fn()
            };
            global.L.circleMarker.mockReturnValue(mockMarker);

            platformRenderer.updatePlatforms(platforms);

            const stats = platformRenderer.getRenderStats();
            expect(stats.totalPlatforms).toBe(1);
            expect(stats.renderTime).toBeGreaterThanOrEqual(0);
            expect(stats.lastUpdate).toBeGreaterThan(0);
        });

        it('should count visible platforms correctly', () => {
            const platforms = [
                { id: 'test-1', platform_type: 'airborne', position: { latitude: 40.0, longitude: -75.0 } },
                { id: 'test-2', platform_type: 'maritime', position: { latitude: 41.0, longitude: -76.0 } }
            ];

            platforms.forEach(platform => {
                platformRenderer.platforms.set(platform.id, platform);
            });

            const count = platformRenderer.countVisiblePlatforms();
            expect(count).toBe(2);

            // Hide one type
            platformRenderer.visibilityFilters.airborne = false;
            const filteredCount = platformRenderer.countVisiblePlatforms();
            expect(filteredCount).toBe(1);
        });

        it('should update marker sizes on zoom change', () => {
            const platform = {
                id: 'test-1',
                platform_type: 'airborne',
                position: { latitude: 40.0, longitude: -75.0 }
            };

            const mockMarker = {
                platformData: platform,
                bindPopup: jest.fn(),
                on: jest.fn(),
                setRadius: jest.fn()
            };

            platformRenderer.platforms.set('test-1', platform);
            platformRenderer.markers.set('test-1', mockMarker);

            platformRenderer.updateMarkerSizes();

            expect(mockMarker.setRadius).toHaveBeenCalled();
        });
    });

    describe('Cleanup', () => {
        beforeEach(() => {
            platformRenderer = new PlatformRenderer(mockMapEngine);
        });

        it('should clear all platforms', () => {
            // Add some test data
            platformRenderer.platforms.set('test-1', {});
            platformRenderer.markers.set('test-1', {});
            platformRenderer.trails.set('test-1', {});
            platformRenderer.trailPoints.set('test-1', []);

            platformRenderer.clearAllPlatforms();

            expect(platformRenderer.platforms.size).toBe(0);
            expect(platformRenderer.markers.size).toBe(0);
            expect(platformRenderer.trails.size).toBe(0);
            expect(platformRenderer.trailPoints.size).toBe(0);
            expect(mockPlatformLayer.clearLayers).toHaveBeenCalled();
            expect(mockTrailLayer.clearLayers).toHaveBeenCalled();
            expect(mockMarkerCluster.clearLayers).toHaveBeenCalled();
        });
    });

    describe('Popup Content Generation', () => {
        beforeEach(() => {
            platformRenderer = new PlatformRenderer(mockMapEngine);
        });

        it('should create detailed popup content', () => {
            const platform = {
                id: 'test-aircraft-1',
                platform_type: 'airborne',
                position: { latitude: 40.1234, longitude: -75.5678, altitude: 10000 },
                velocity: { north: 100, east: 50, up: 10 },
                speed: 150.5,
                heading: 90,
                lastUpdated: 1623456789000
            };

            const content = platformRenderer.createPopupContent(platform);

            expect(content).toContain('test-aircraft-1');
            expect(content).toContain('airborne');
            expect(content).toContain('40.1234, -75.5678');
            expect(content).toContain('10000 m');
            expect(content).toContain('150.5 m/s');
            expect(content).toContain('90°');
            expect(content).toContain('N:100.0 E:50.0 U:10.0');
        });
    });
});
