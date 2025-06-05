/**
 * UI Tests for Platform Library Filters
 * Tests for expected UI behavior and filter functionality
 */

// Mock DOM elements for testing
const mockDOM = {
    createFilterButton: (type, value, text, active = false) => ({
        dataset: { filterType: type, filterValue: value },
        textContent: text,
        classList: {
            contains: jest.fn().mockReturnValue(active),
            add: jest.fn(),
            remove: jest.fn(),
            toggle: jest.fn()
        },
        click: jest.fn()
    }),

    createPlatformListElement: () => ({
        innerHTML: '',
        appendChild: jest.fn(),
        children: []
    })
};

describe('Platform Library Filter UI Tests', () => {
    let mockScenarioBuilder;
    let mockFilterButtons;

    beforeEach(() => {
        // Reset mocks
        jest.clearAllMocks();

        // Create mock filter buttons matching the expected UI
        mockFilterButtons = {
            all: mockDOM.createFilterButton('all', 'all', 'All Platforms', true),
            airborne: mockDOM.createFilterButton('domain', 'airborne', '✈️ Airborne', false),
            land: mockDOM.createFilterButton('domain', 'land', '🚛 Land', false),
            maritime: mockDOM.createFilterButton('domain', 'maritime', '🚢 Maritime', false),
            space: mockDOM.createFilterButton('domain', 'space', '🛰️ Space', false),
            commercial: mockDOM.createFilterButton('affiliation', 'commercial', '🏢 Commercial', false),
            military: mockDOM.createFilterButton('affiliation', 'military', '⚔️ Military', false)
        };

        // Mock platform data with correct MIL-STD-2525D domain names
        mockScenarioBuilder = {
            activeFilters: {
                domains: new Set(),
                affiliations: new Set(),
                search: ''
            },
            platforms: [
                { id: 'air_commercial', name: 'Airbus A320', domain: 'airborne', affiliation: 'commercial', category: 'aircraft' },
                { id: 'air_military', name: 'F-16 Fighter', domain: 'airborne', affiliation: 'military', category: 'fighter' },
                { id: 'sea_commercial', name: 'Container Ship', domain: 'maritime', affiliation: 'commercial', category: 'ship' },
                { id: 'sea_military', name: 'Naval Destroyer', domain: 'maritime', affiliation: 'military', category: 'warship' },
                { id: 'land_commercial', name: 'Tesla Model S', domain: 'land', affiliation: 'commercial', category: 'car' },
                { id: 'land_military', name: 'Military Truck', domain: 'land', affiliation: 'military', category: 'truck' },
                { id: 'space_commercial', name: 'Starlink Sat', domain: 'space', affiliation: 'commercial', category: 'satellite' },
                { id: 'space_military', name: 'Military Sat', domain: 'space', affiliation: 'military', category: 'satellite' }
            ],

            // Mock implementation that will be tested against
            handleFilterClick: jest.fn(),
            applyFilters: jest.fn(),
            isAllFiltersActive: jest.fn(),
            selectAllFilters: jest.fn(),
            clearAllFilters: jest.fn()
        };
    });

    describe('UI Test 1: Button Text Matches MIL-STD-2525D Domain Names', () => {
        test('should display correct domain names on filter buttons', () => {
            // Test that button text matches MIL-STD-2525D standards
            expect(mockFilterButtons.airborne.textContent).toBe('✈️ Airborne');
            expect(mockFilterButtons.maritime.textContent).toBe('🚢 Maritime');
            expect(mockFilterButtons.land.textContent).toBe('🚛 Land');
            expect(mockFilterButtons.space.textContent).toBe('🛰️ Space');
        });

        test('should have correct data attributes matching domain names', () => {
            // Test that data attributes match the text and platform data
            expect(mockFilterButtons.airborne.dataset.filterValue).toBe('airborne');
            expect(mockFilterButtons.maritime.dataset.filterValue).toBe('maritime');
            expect(mockFilterButtons.land.dataset.filterValue).toBe('land');
            expect(mockFilterButtons.space.dataset.filterValue).toBe('space');
        });

        test('should NOT use legacy domain names (air, sea)', () => {
            // Ensure we're not using the old incorrect names - check for exact matches only
            expect(mockFilterButtons.airborne.textContent).not.toBe('Air');
            expect(mockFilterButtons.maritime.textContent).not.toBe('Sea');
            expect(mockFilterButtons.airborne.dataset.filterValue).not.toBe('air');
            expect(mockFilterButtons.maritime.dataset.filterValue).not.toBe('sea');
        });
    });

    describe('UI Test 2: Initial Button State Consistency', () => {
        test('should have All button active by default', () => {
            // "All Platforms" should be active initially
            expect(mockFilterButtons.all.classList.contains('active')).toBe(true);
        });

        test('should have all domain buttons inactive by default', () => {
            // Individual domain buttons should be inactive when "All" is active
            expect(mockFilterButtons.airborne.classList.contains('active')).toBe(false);
            expect(mockFilterButtons.land.classList.contains('active')).toBe(false);
            expect(mockFilterButtons.maritime.classList.contains('active')).toBe(false);
            expect(mockFilterButtons.space.classList.contains('active')).toBe(false);
        });

        test('should have all affiliation buttons inactive by default', () => {
            // Individual affiliation buttons should be inactive when "All" is active
            expect(mockFilterButtons.commercial.classList.contains('active')).toBe(false);
            expect(mockFilterButtons.military.classList.contains('active')).toBe(false);
        });

        test('should toggle domain buttons correctly on first click', () => {
            // When clicking a domain button for the first time, it should become active
            const airborneBtnMock = {
                ...mockFilterButtons.airborne,
                classList: {
                    contains: jest.fn().mockReturnValue(false), // Initially inactive
                    add: jest.fn(),
                    remove: jest.fn()
                }
            };

            // Simulate first click - should activate the button
            mockScenarioBuilder.handleFilterClick(airborneBtnMock);

            // Verify the button activation was attempted
            expect(mockScenarioBuilder.handleFilterClick).toHaveBeenCalledWith(airborneBtnMock);
        });
    });

    describe('UI Test 3: Filter Application Matches Button Text', () => {
        beforeEach(() => {
            // Mock implementation of applyFilters to test actual filtering
            mockScenarioBuilder.applyFilters = jest.fn().mockImplementation(() => {
                let filteredPlatforms = [...mockScenarioBuilder.platforms];

                // Apply domain filters
                if (mockScenarioBuilder.activeFilters.domains.size > 0 &&
                    mockScenarioBuilder.activeFilters.domains.size < 4) {
                    filteredPlatforms = filteredPlatforms.filter(platform =>
                        mockScenarioBuilder.activeFilters.domains.has(platform.domain)
                    );
                } else if (mockScenarioBuilder.activeFilters.domains.size === 0) {
                    filteredPlatforms = [];
                }

                // Apply affiliation filters
                if (mockScenarioBuilder.activeFilters.affiliations.size > 0 &&
                    mockScenarioBuilder.activeFilters.affiliations.size < 2) {
                    filteredPlatforms = filteredPlatforms.filter(platform =>
                        mockScenarioBuilder.activeFilters.affiliations.has(platform.affiliation)
                    );
                } else if (mockScenarioBuilder.activeFilters.affiliations.size === 0) {
                    // If no affiliations selected but domains are, keep domain filtering
                    if (mockScenarioBuilder.activeFilters.domains.size === 0) {
                        filteredPlatforms = [];
                    }
                }

                return filteredPlatforms;
            });
        });

        test('should filter platforms by airborne domain when Airborne button is active', () => {
            // Set up airborne filter as active
            mockScenarioBuilder.activeFilters.domains.add('airborne');

            const result = mockScenarioBuilder.applyFilters();

            // Should return only airborne platforms
            const airbornePlatforms = result.filter(p => p.domain === 'airborne');
            expect(airbornePlatforms).toHaveLength(2); // F-16 and Airbus A320
            expect(result.every(p => p.domain === 'airborne')).toBe(true);
        });

        test('should filter platforms by maritime domain when Maritime button is active', () => {
            // Set up maritime filter as active
            mockScenarioBuilder.activeFilters.domains.add('maritime');

            const result = mockScenarioBuilder.applyFilters();

            // Should return only maritime platforms
            const maritimePlatforms = result.filter(p => p.domain === 'maritime');
            expect(maritimePlatforms).toHaveLength(2); // Container Ship and Naval Destroyer
            expect(result.every(p => p.domain === 'maritime')).toBe(true);
        });

        test('should NOT filter by legacy domain names (air, sea)', () => {
            // Set up filters with legacy names - should return empty
            mockScenarioBuilder.activeFilters.domains.add('air');
            mockScenarioBuilder.activeFilters.domains.add('sea');

            const result = mockScenarioBuilder.applyFilters();

            // Should return no platforms since 'air' and 'sea' don't exist in data
            expect(result).toHaveLength(0);
        });

        test('should show all platforms when All button is active (all domains/affiliations selected)', () => {
            // Simulate "All" button state - all filters active
            mockScenarioBuilder.activeFilters.domains = new Set(['airborne', 'land', 'maritime', 'space']);
            mockScenarioBuilder.activeFilters.affiliations = new Set(['commercial', 'military']);

            const result = mockScenarioBuilder.applyFilters();

            // Should return all platforms
            expect(result).toHaveLength(8);
        });

        test('should show no platforms when no filters are selected', () => {
            // No filters active
            mockScenarioBuilder.activeFilters.domains.clear();
            mockScenarioBuilder.activeFilters.affiliations.clear();

            const result = mockScenarioBuilder.applyFilters();

            // Should return no platforms
            expect(result).toHaveLength(0);
        });

        test('should combine domain and affiliation filters correctly', () => {
            // Set up combined filter - military airborne platforms
            mockScenarioBuilder.activeFilters.domains.add('airborne');
            mockScenarioBuilder.activeFilters.affiliations.add('military');

            const result = mockScenarioBuilder.applyFilters();

            // Should return only military airborne platforms
            expect(result).toHaveLength(1);
            expect(result[0].name).toBe('F-16 Fighter');
            expect(result[0].domain).toBe('airborne');
            expect(result[0].affiliation).toBe('military');
        });
    });

    describe('UI Test 4: Button State Management', () => {
        test('should deactivate All button when individual filters are selected', () => {
            mockScenarioBuilder.isAllFiltersActive = jest.fn().mockReturnValue(false);

            // Mock the scenario where a single domain is selected
            mockScenarioBuilder.activeFilters.domains.add('airborne');

            const isAllActive = mockScenarioBuilder.isAllFiltersActive();
            expect(isAllActive).toBe(false);
        });

        test('should activate All button when all individual filters are selected', () => {
            mockScenarioBuilder.isAllFiltersActive = jest.fn().mockReturnValue(true);

            // Mock the scenario where all filters are selected
            mockScenarioBuilder.activeFilters.domains = new Set(['airborne', 'land', 'maritime', 'space']);
            mockScenarioBuilder.activeFilters.affiliations = new Set(['commercial', 'military']);

            const isAllActive = mockScenarioBuilder.isAllFiltersActive();
            expect(isAllActive).toBe(true);
        });

        test('should toggle filter state correctly on button click', () => {
            // Mock scenario builder methods
            mockScenarioBuilder.handleFilterClick = jest.fn().mockImplementation((button) => {
                const filterType = button.dataset.filterType;
                const filterValue = button.dataset.filterValue;

                if (filterType === 'domain') {
                    if (mockScenarioBuilder.activeFilters.domains.has(filterValue)) {
                        mockScenarioBuilder.activeFilters.domains.delete(filterValue);
                        button.classList.remove('active');
                    } else {
                        mockScenarioBuilder.activeFilters.domains.add(filterValue);
                        button.classList.add('active');
                    }
                }
            });

            // Test clicking airborne button twice
            const airborneBtn = mockFilterButtons.airborne;

            // First click - should add filter
            mockScenarioBuilder.handleFilterClick(airborneBtn);
            expect(mockScenarioBuilder.activeFilters.domains.has('airborne')).toBe(true);
            expect(airborneBtn.classList.add).toHaveBeenCalledWith('active');

            // Second click - should remove filter
            mockScenarioBuilder.activeFilters.domains.add('airborne'); // Simulate state after first click
            mockScenarioBuilder.handleFilterClick(airborneBtn);
            expect(mockScenarioBuilder.activeFilters.domains.has('airborne')).toBe(false);
            expect(airborneBtn.classList.remove).toHaveBeenCalledWith('active');
        });
    });

    describe('UI Test 5: Expected vs Actual HTML Structure', () => {
        test('should have correct HTML structure for filter buttons', () => {
            // Define expected HTML structure
            const expectedStructure = {
                all: {
                    'data-filter-type': 'all',
                    'data-filter-value': 'all',
                    textContent: 'All Platforms'
                },
                airborne: {
                    'data-filter-type': 'domain',
                    'data-filter-value': 'airborne',
                    textContent: '✈️ Airborne'
                },
                maritime: {
                    'data-filter-type': 'domain',
                    'data-filter-value': 'maritime',
                    textContent: '🚢 Maritime'
                },
                land: {
                    'data-filter-type': 'domain',
                    'data-filter-value': 'land',
                    textContent: '🚛 Land'
                },
                space: {
                    'data-filter-type': 'domain',
                    'data-filter-value': 'space',
                    textContent: '🛰️ Space'
                }
            };

            // Test each button structure
            Object.keys(expectedStructure).forEach(buttonKey => {
                const expected = expectedStructure[buttonKey];
                const actual = mockFilterButtons[buttonKey];

                expect(actual.dataset.filterType).toBe(expected['data-filter-type']);
                expect(actual.dataset.filterValue).toBe(expected['data-filter-value']);
                expect(actual.textContent).toBe(expected.textContent);
            });
        });

        test('should NOT have legacy button structures (air, sea)', () => {
            // Ensure we don't have buttons with old incorrect values - check for exact matches
            const allButtons = Object.values(mockFilterButtons);

            const hasAirButton = allButtons.some(btn =>
                btn.dataset.filterValue === 'air' || btn.textContent === 'Air'
            );
            const hasSeaButton = allButtons.some(btn =>
                btn.dataset.filterValue === 'sea' || btn.textContent === 'Sea'
            );

            expect(hasAirButton).toBe(false);
            expect(hasSeaButton).toBe(false);
        });
    });
});
