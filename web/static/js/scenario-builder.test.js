/**
 * @jest-environment jsdom
 */

// Mock Leaflet library
global.L = {
    map: jest.fn(() => ({
        setView: jest.fn().mockReturnThis(),
        on: jest.fn(),
        off: jest.fn(),
        removeLayer: jest.fn(),
        eachLayer: jest.fn()
    })),
    tileLayer: jest.fn(() => ({
        addTo: jest.fn()
    })),
    marker: jest.fn(() => ({
        addTo: jest.fn().mockReturnThis(),
        bindPopup: jest.fn().mockReturnThis()
    })),
    divIcon: jest.fn(() => ({})),
    circleMarker: jest.fn(() => ({
        addTo: jest.fn().mockReturnThis(),
        bindPopup: jest.fn()
    })),
    polyline: jest.fn(() => ({
        addTo: jest.fn().mockReturnThis()
    })),
    layerGroup: jest.fn(() => ({
        addTo: jest.fn().mockReturnThis()
    })),
    featureGroup: jest.fn(() => ({
        getBounds: jest.fn(() => ({
            pad: jest.fn(() => 'mock-bounds')
        }))
    }))
};

// Mock fetch
global.fetch = jest.fn();

// Mock URL and Blob for file export
global.URL = {
    createObjectURL: jest.fn(() => 'mock-url'),
    revokeObjectURL: jest.fn()
};
global.Blob = jest.fn();

// Mock FileReader
global.FileReader = jest.fn(() => ({
    readAsText: jest.fn(),
    onload: null,
    onerror: null,
    result: 'mock file content'
}));

// Mock setTimeout to avoid timing issues in tests
global.setTimeout = jest.fn((callback, delay) => {
    // Don't execute callback immediately in tests
    return 123;
});

// Mock alert and confirm
global.alert = jest.fn();
global.confirm = jest.fn(() => true);

// Mock DOM elements
const mockElements = {
    'map': { style: {} },
    'platformList': { innerHTML: '', appendChild: jest.fn() },
    'platformSearch': { addEventListener: jest.fn() },
    'scenarioPlatforms': { innerHTML: '', appendChild: jest.fn() },
    'statusBar': { textContent: '' },
    'mapInstructions': { textContent: '' },
    'platformModal': { style: { display: 'none' } },
    'modalTitle': { textContent: '' },
    'modalPlatformId': { value: '' },
    'modalPlatformName': { value: '' },
    'modalLatitude': { value: '' },
    'modalLongitude': { value: '' },
    'modalAltitude': { value: '' },
    'modalMissionType': { value: 'transport' },
    'waypointMode': { addEventListener: jest.fn(), checked: false },
    'completeRoute': { addEventListener: jest.fn(), disabled: false },
    'routeControls': { style: { display: 'none' } },
    'routeModal': { style: { display: 'none' } },
    'routeModalTitle': { textContent: '' },
    'routePlatformId': { value: '' },
    'routePlatformName': { value: '' },
    'routeSummary': { innerHTML: '' },
    'routeSpeed': { value: '10' },
    'routeAltitude': { value: '1000' },
    'routeMissionType': { value: 'transport' },
    'scenarioName': { value: 'Test Scenario' },
    'scenarioDescription': { value: 'Test Description' },
    'scenarioDuration': { value: '30' }
};

document.getElementById = jest.fn((id) => mockElements[id] || null);
document.querySelectorAll = jest.fn((selector) => {
    if (selector === '.domain-filter button') {
        return [
            { addEventListener: jest.fn(), classList: { remove: jest.fn(), add: jest.fn() }, dataset: { domain: 'all' } },
            { addEventListener: jest.fn(), classList: { remove: jest.fn(), add: jest.fn() }, dataset: { domain: 'airborne' } }
        ];
    }
    if (selector === '.close') {
        return [{ addEventListener: jest.fn(), closest: jest.fn(() => ({ style: { display: 'none' } })) }];
    }
    if (selector === '.platform-item') {
        return [
            { classList: { remove: jest.fn(), add: jest.fn() }, style: {}, textContent: 'Airbus A320 Commercial Aircraft' },
            { classList: { remove: jest.fn(), add: jest.fn() }, style: {}, textContent: 'F-16 Fighter Jet' }
        ];
    }
    return [];
});

document.createElement = jest.fn((tag) => ({
    className: '',
    innerHTML: '',
    addEventListener: jest.fn(),
    appendChild: jest.fn(),
    download: '',
    href: '',
    click: jest.fn(),
    style: {},
    classList: {
        add: jest.fn(),
        remove: jest.fn(),
        contains: jest.fn()
    }
}));

// Mock document body
Object.defineProperty(document, 'body', {
    value: {
        appendChild: jest.fn(),
        removeChild: jest.fn()
    },
    writable: true
});

// Mock window and document events
global.window = {
    addEventListener: jest.fn()
};

document.addEventListener = jest.fn();

// Import the actual ScenarioBuilder class
const ScenarioBuilder = require('./scenario-builder.js');

describe('ScenarioBuilder', () => {
    let scenarioBuilder;

    beforeEach(() => {
        jest.clearAllMocks();

        // Reset mock element values
        Object.keys(mockElements).forEach(key => {
            if (mockElements[key].value !== undefined) {
                mockElements[key].value = key === 'scenarioName' ? 'Test Scenario' :
                    key === 'scenarioDescription' ? 'Test Description' :
                        key === 'scenarioDuration' ? '30' : '';
            }
            if (mockElements[key].innerHTML !== undefined) {
                mockElements[key].innerHTML = '';
            }
            if (mockElements[key].textContent !== undefined) {
                mockElements[key].textContent = '';
            }
        });

        scenarioBuilder = new ScenarioBuilder();
        // Initialize routePolylines to prevent issues with clearRoutePolylines
        scenarioBuilder.routePolylines = [];
    });

    describe('Constructor and Initialization', () => {
        test('should initialize with default values', () => {
            // The constructor calls init() which loads platforms, so we need to account for that
            expect(scenarioBuilder.scenarioPlatforms).toEqual([]);
            expect(scenarioBuilder.selectedPlatform).toBeNull();
            expect(scenarioBuilder.platformCounter).toBe(1);
            expect(scenarioBuilder.waypointMode).toBe(false);
            expect(scenarioBuilder.currentRoute).toEqual([]);
            // platforms array will be populated after loadPlatforms() is called in init()
            expect(Array.isArray(scenarioBuilder.platforms)).toBe(true);
        });

        test('should initialize map correctly', () => {
            scenarioBuilder.initMap();
            expect(L.map).toHaveBeenCalledWith('map');
            expect(L.tileLayer).toHaveBeenCalled();
            expect(L.layerGroup).toHaveBeenCalled();
        });

        test('should handle missing Leaflet library gracefully', () => {
            const originalL = global.L;
            global.L = undefined;
            const consoleSpy = jest.spyOn(console, 'error').mockImplementation();

            scenarioBuilder.initMap();

            expect(consoleSpy).toHaveBeenCalledWith('Leaflet library not loaded');

            global.L = originalL;
            consoleSpy.mockRestore();
        });

        test('should load platforms from server successfully', async () => {
            const mockPlatforms = [{ id: 'test', name: 'Test Platform' }];
            global.fetch.mockResolvedValueOnce({
                ok: true,
                json: jest.fn().mockResolvedValue(mockPlatforms)
            });

            await scenarioBuilder.loadPlatforms();

            expect(scenarioBuilder.platforms).toEqual(mockPlatforms);
        });

        test('should fallback to default platforms when server fails', async () => {
            global.fetch.mockResolvedValueOnce({
                ok: false
            });

            await scenarioBuilder.loadPlatforms();

            expect(scenarioBuilder.platforms.length).toBeGreaterThan(0);
            expect(scenarioBuilder.platforms[0]).toHaveProperty('id');
        });

        test('should fallback to default platforms when fetch throws', async () => {
            global.fetch.mockRejectedValueOnce(new Error('Network error'));
            const consoleSpy = jest.spyOn(console, 'warn').mockImplementation();

            await scenarioBuilder.loadPlatforms();

            expect(scenarioBuilder.platforms.length).toBeGreaterThan(0);
            expect(consoleSpy).toHaveBeenCalled();

            consoleSpy.mockRestore();
        });
    });

    describe('Platform Management', () => {
        beforeEach(async () => {
            await scenarioBuilder.loadPlatforms();
        });

        test('should select platform correctly', () => {
            const platform = scenarioBuilder.platforms[0];
            const mockElement = { classList: { add: jest.fn(), remove: jest.fn() } };

            scenarioBuilder.selectPlatform(platform, mockElement);

            expect(scenarioBuilder.selectedPlatform).toBe(platform);
            expect(mockElement.classList.add).toHaveBeenCalledWith('selected');
        });

        test('should get correct domain icons', () => {
            expect(scenarioBuilder.getDomainIcon('airborne')).toBe('✈️');
            expect(scenarioBuilder.getDomainIcon('maritime')).toBe('🚢');
            expect(scenarioBuilder.getDomainIcon('land')).toBe('🚛');
            expect(scenarioBuilder.getDomainIcon('space')).toBe('🛰️');
            expect(scenarioBuilder.getDomainIcon('unknown')).toBe('🔹');
        });

        test('should render platform list', () => {
            scenarioBuilder.renderPlatformList();
            expect(document.createElement).toHaveBeenCalledWith('div');
            expect(mockElements.platformList.appendChild).toHaveBeenCalled();
        });

        test('should filter platforms by domain', () => {
            scenarioBuilder.filterPlatformsByDomain('airborne');
            expect(document.createElement).toHaveBeenCalled();
        });

        test('should reset platform list when filtering by all domains', () => {
            const renderSpy = jest.spyOn(scenarioBuilder, 'renderPlatformList');
            scenarioBuilder.filterPlatformsByDomain('all');
            expect(renderSpy).toHaveBeenCalled();
        });

        test('should not place platform when platform is null', () => {
            const latlng = { lat: 40.0, lng: -100.0 };
            const initialLength = scenarioBuilder.scenarioPlatforms.length;

            scenarioBuilder.placePlatformOnMap(latlng, null);

            expect(scenarioBuilder.scenarioPlatforms.length).toBe(initialLength);
        });

        test('should generate platform names based on domain', () => {
            const platforms = [
                { domain: 'airborne' },
                { domain: 'maritime' },
                { domain: 'land' },
                { domain: 'space' },
                { domain: 'unknown' }
            ];

            platforms.forEach(platform => {
                const name = scenarioBuilder.generatePlatformName(platform);
                expect(typeof name).toBe('string');
                expect(name.length).toBeGreaterThan(0);
            });
        });

        test('should clear platform selection', () => {
            scenarioBuilder.selectedPlatform = scenarioBuilder.platforms[0];
            scenarioBuilder.clearPlatformSelection();

            expect(scenarioBuilder.selectedPlatform).toBeNull();
        });

        test('should get platform icon for different domains', () => {
            const domains = ['airborne', 'maritime', 'land', 'space', 'unknown'];

            domains.forEach(domain => {
                const icon = scenarioBuilder.getPlatformIcon(domain);
                expect(L.divIcon).toHaveBeenCalled();
            });
        });
    });

    describe('Platform Configuration Modal', () => {
        beforeEach(async () => {
            await scenarioBuilder.loadPlatforms();
            scenarioBuilder.selectedPlatform = scenarioBuilder.platforms[0];
        });

        test('should show platform config modal', () => {
            const latlng = { lat: 40.0, lng: -100.0 };
            scenarioBuilder.showPlatformConfigModal(latlng);

            expect(mockElements.modalTitle.textContent).toContain(scenarioBuilder.selectedPlatform.name);
            expect(mockElements.modalLatitude.value).toBe('40.000000');
            expect(mockElements.modalLongitude.value).toBe('-100.000000');
            expect(mockElements.platformModal.style.display).toBe('block');
        });

        test('should not save platform without required fields', () => {
            mockElements.modalPlatformId.value = '';
            mockElements.modalPlatformName.value = '';

            scenarioBuilder.savePlatformConfig();

            expect(global.alert).toHaveBeenCalledWith('Please fill in all required fields');
        });

        test('should add map marker for platform', () => {
            const platform = {
                name: 'Test Platform',
                class: 'Test Class',
                domain: 'airborne',
                start_position: { latitude: 40.0, longitude: -100.0, altitude: 1000 },
                mission: { type: 'transport' }
            };

            scenarioBuilder.addMapMarker(platform);

            expect(L.marker).toHaveBeenCalledWith([40.0, -100.0], { icon: expect.anything() });
        });
    });

    describe('Scenario Platform Management', () => {
        test('should render empty message when no platforms', () => {
            scenarioBuilder.scenarioPlatforms = [];
            scenarioBuilder.renderScenarioPlatforms();
            expect(mockElements.scenarioPlatforms.innerHTML).toContain('No platforms added yet');
        });

        test('should remove platform from scenario', () => {
            const platform = {
                name: 'Test Platform',
                start_position: { latitude: 40.0, longitude: -100.0, altitude: 1000 }
            };
            scenarioBuilder.scenarioPlatforms = [platform];
            scenarioBuilder.mapMarkers = [{ platform, marker: { remove: jest.fn() } }];
            scenarioBuilder.map = { removeLayer: jest.fn() };

            scenarioBuilder.removePlatform(0);

            // The removePlatform method doesn't return a value, so just check the side effects
            expect(scenarioBuilder.scenarioPlatforms.length).toBe(0);
        });

        test('should clear entire scenario', () => {
            scenarioBuilder.currentRoute = [{ lat: 1, lng: 1 }];
            scenarioBuilder.mapMarkers = [{ marker: { remove: jest.fn() } }];
            scenarioBuilder.map = { removeLayer: jest.fn() };

            scenarioBuilder.clearScenario();

            expect(scenarioBuilder.scenarioPlatforms.length).toBe(0);
            expect(scenarioBuilder.mapMarkers.length).toBe(0);
        });
    });

    describe('Search and Filtering', () => {
        test('should filter platforms by search term', () => {
            scenarioBuilder.filterPlatforms('airbus');

            const mockItems = document.querySelectorAll('.platform-item');
            // Verify the filter logic is called
            expect(document.querySelectorAll).toHaveBeenCalledWith('.platform-item');
        });
    });

    describe('Waypoint and Route Management', () => {
        beforeEach(() => {
            scenarioBuilder.map = {
                on: jest.fn(),
                off: jest.fn(),
                removeLayer: jest.fn()
            };
        });

        test('should toggle waypoint mode on', () => {
            scenarioBuilder.toggleWaypointMode(true);

            expect(scenarioBuilder.waypointMode).toBe(true);
            expect(scenarioBuilder.map.off).toHaveBeenCalledWith('click');
            expect(scenarioBuilder.map.on).toHaveBeenCalledWith('click', expect.any(Function));
        });

        test('should toggle waypoint mode off', () => {
            scenarioBuilder.waypointMode = true;
            scenarioBuilder.currentRoute = [{ lat: 1, lng: 1 }];
            scenarioBuilder.routePolylines = [{ remove: jest.fn() }];
            scenarioBuilder.mapMarkers = [];
            
            // Ensure map is properly mocked
            scenarioBuilder.map = {
                on: jest.fn(),
                off: jest.fn(),
                removeLayer: jest.fn()
            };

            scenarioBuilder.toggleWaypointMode(false);

            expect(scenarioBuilder.waypointMode).toBe(false);
            expect(scenarioBuilder.currentRoute.length).toBe(0);
            expect(scenarioBuilder.map.off).toHaveBeenCalledWith('click');
            expect(scenarioBuilder.map.on).toHaveBeenCalledWith('click', expect.any(Function));
        });

        test('should add waypoint in waypoint mode', () => {
            scenarioBuilder.waypointMode = true;
            const latlng = { lat: 40.0, lng: -100.0 };

            scenarioBuilder.addWaypoint(latlng);

            expect(scenarioBuilder.currentRoute.length).toBe(1);
            expect(L.marker).toHaveBeenCalled();
        });

        test('should not add waypoint when not in waypoint mode', () => {
            scenarioBuilder.waypointMode = false;
            const latlng = { lat: 40.0, lng: -100.0 };
            const initialLength = scenarioBuilder.currentRoute.length;

            scenarioBuilder.addWaypoint(latlng);

            expect(scenarioBuilder.currentRoute.length).toBe(initialLength);
        });

        test('should create polyline with multiple waypoints', () => {
            scenarioBuilder.waypointMode = true;
            scenarioBuilder.currentRoute = [{ lat: 40.0, lng: -100.0 }]; // Pre-existing waypoint

            const latlng = { lat: 41.0, lng: -101.0 };
            scenarioBuilder.addWaypoint(latlng);

            expect(L.polyline).toHaveBeenCalled();
            expect(scenarioBuilder.routePolylines.length).toBe(1);
        });

        test('should complete route with sufficient waypoints', () => {
            scenarioBuilder.currentRoute = [
                { lat: 40.0, lng: -100.0 },
                { lat: 41.0, lng: -101.0 }
            ];

            scenarioBuilder.completeCurrentRoute();

            expect(mockElements.waypointMode.checked).toBe(false);
        });

        test('should not complete route with insufficient waypoints', () => {
            scenarioBuilder.currentRoute = [{ lat: 40.0, lng: -100.0 }];

            scenarioBuilder.completeCurrentRoute();

            expect(global.alert).toHaveBeenCalledWith('Please add at least 2 waypoints to create a route');
        });
    });

    describe('Scenario Export and Import', () => {
        test('should not export empty scenario', () => {
            scenarioBuilder.scenarioPlatforms = [];

            scenarioBuilder.exportScenario();

            expect(global.alert).toHaveBeenCalledWith('Please add at least one platform to export');
        });

        test('should generate YAML for scenario with platforms', () => {
            scenarioBuilder.scenarioPlatforms = [{
                id: 'test-1',
                type: 'test_type',
                name: 'Test Platform',
                start_position: { latitude: 40.0, longitude: -100.0, altitude: 1000 },
                mission: { type: 'patrol' }
            }];

            const yaml = scenarioBuilder.generateScenarioYAML();

            expect(yaml).toContain('Test Scenario');
            expect(yaml).toContain('test-1');
            expect(yaml).toContain('Test Platform');
        });

        test('should load scenario file', () => {
            const mockFile = new File(['test content'], 'test.yaml');

            scenarioBuilder.loadScenario(mockFile);

            expect(global.FileReader).toHaveBeenCalled();
        });
    });

    describe('Scenario Validation', () => {
        test('should validate scenario with issues', () => {
            scenarioBuilder.scenarioPlatforms = [];

            const issues = scenarioBuilder.validateScenario();

            expect(issues).toContain('No platforms added to scenario');
        });

        test('should detect overlapping platforms', () => {
            scenarioBuilder.scenarioPlatforms = [
                {
                    name: 'Platform 1',
                    start_position: { latitude: 40.0, longitude: -100.0, altitude: 1000 }
                },
                {
                    name: 'Platform 2',
                    start_position: { latitude: 40.0001, longitude: -100.0001, altitude: 1000 }
                }
            ];

            const issues = scenarioBuilder.validateScenario();

            expect(issues.some(issue => issue.includes('very close'))).toBe(true);
        });
    });

    describe('Additional Functionality', () => {
        test('should get domain statistics', () => {
            scenarioBuilder.scenarioPlatforms = [
                { domain: 'airborne' },
                { domain: 'airborne' },
                { domain: 'maritime' }
            ];

            const stats = scenarioBuilder.getDomainStats();

            expect(stats.airborne).toBe(2);
            expect(stats.maritime).toBe(1);
            expect(stats.land).toBe(0);
            expect(stats.space).toBe(0);
        });

        test('should preview scenario with platforms', () => {
            scenarioBuilder.scenarioPlatforms = [{ name: 'Test' }];
            scenarioBuilder.mapMarkers = [{ marker: {} }];
            scenarioBuilder.map = { fitBounds: jest.fn() };

            scenarioBuilder.previewScenario();

            expect(L.featureGroup).toHaveBeenCalled();
            expect(scenarioBuilder.map.fitBounds).toHaveBeenCalled();
        });

        test('should not preview empty scenario', () => {
            scenarioBuilder.scenarioPlatforms = [];

            scenarioBuilder.previewScenario();

            expect(global.alert).toHaveBeenCalledWith('Please add platforms to preview');
        });

        test('should load platform library', () => {
            scenarioBuilder.loadPlatformLibrary();

            expect(scenarioBuilder.platformLibrary).toBeDefined();
            expect(scenarioBuilder.platformLibrary.airborne).toBeDefined();
            expect(scenarioBuilder.platformLibrary.maritime).toBeDefined();
            expect(scenarioBuilder.platformLibrary.land).toBeDefined();
            expect(scenarioBuilder.platformLibrary.space).toBeDefined();
        });
    });

    describe('Default Platform Library', () => {
        test('should have platforms for all domains', () => {
            const platforms = scenarioBuilder.getDefaultPlatforms();
            const domains = ['airborne', 'maritime', 'land', 'space'];

            domains.forEach(domain => {
                const domainPlatforms = platforms.filter(p => p.domain === domain);
                expect(domainPlatforms.length).toBeGreaterThan(0);
            });
        });

        test('should have required platform properties', () => {
            const platforms = scenarioBuilder.getDefaultPlatforms();

            platforms.forEach(platform => {
                expect(platform).toHaveProperty('id');
                expect(platform).toHaveProperty('name');
                expect(platform).toHaveProperty('class');
                expect(platform).toHaveProperty('category');
                expect(platform).toHaveProperty('domain');
                expect(platform).toHaveProperty('description');
                expect(platform).toHaveProperty('performance');
            });
        });
    });
});

// Export for potential use in other tests
if (typeof module !== 'undefined' && module.exports) {
    module.exports = { ScenarioBuilder };
}
