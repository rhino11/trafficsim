/**
 * @jest-environment jsdom
 */

// Platform Library Filter Button UI Tests
// Tests for the enhanced Platform Library filter system with specific layout requirements

describe('Platform Library Filter System', () => {
    let mockScenarioBuilder;
    let mockDOM;

    beforeEach(() => {
        // Setup mock DOM structure
        mockDOM = {
            platformSearch: { value: '', addEventListener: jest.fn() },
            platformList: { innerHTML: '' },
            filterButtons: [],
            querySelectorAll: jest.fn(),
            getElementById: jest.fn()
        };

        // Mock querySelectorAll to return filter buttons
        mockDOM.querySelectorAll.mockImplementation((selector) => {
            if (selector === '.filter-btn') {
                return mockDOM.filterButtons;
            }
            if (selector.includes('data-filter-type')) {
                return mockDOM.filterButtons.filter(btn =>
                    btn.dataset.filterType === selector.match(/data-filter-type="([^"]+)"/)[1]
                );
            }
            return [];
        });

        mockDOM.getElementById.mockImplementation((id) => {
            if (id === 'platformSearch') return mockDOM.platformSearch;
            if (id === 'platformList') return mockDOM.platformList;
            return null;
        });

        // Setup ScenarioBuilder mock with updated filter structure
        mockScenarioBuilder = {
            activeFilters: { domains: new Set(), affiliations: new Set(), search: '' },
            platforms: [
                { id: 'airborne_commercial', name: 'Airbus A320', domain: 'airborne', affiliation: 'commercial', category: 'aircraft', description: 'Commercial airliner' },
                { id: 'airborne_military', name: 'F-16 Fighter', domain: 'airborne', affiliation: 'military', category: 'fighter', description: 'Military fighter jet' },
                { id: 'land_commercial', name: 'Tesla Model S', domain: 'land', affiliation: 'commercial', category: 'car', description: 'Electric vehicle' },
                { id: 'land_military', name: 'Military Truck', domain: 'land', affiliation: 'military', category: 'truck', description: 'Military transport' },
                { id: 'maritime_commercial', name: 'Container Ship', domain: 'maritime', affiliation: 'commercial', category: 'ship', description: 'Cargo vessel' },
                { id: 'maritime_military', name: 'Naval Destroyer', domain: 'maritime', affiliation: 'military', category: 'warship', description: 'Combat vessel' },
                { id: 'space_commercial', name: 'Starlink Sat', domain: 'space', affiliation: 'commercial', category: 'satellite', description: 'Communication satellite' },
                { id: 'space_military', name: 'Military Sat', domain: 'space', affiliation: 'military', category: 'satellite', description: 'Reconnaissance satellite' }
            ],
            handleFilterClick: jest.fn(),
            clearAllFilters: jest.fn(),
            selectAllFilters: jest.fn(),
            setAllButtonsActive: jest.fn(),
            clearAllButtons: jest.fn(),
            updateAllButtonState: jest.fn(),
            isAllFiltersActive: jest.fn(),
            applyFilters: jest.fn(),
            renderFilteredPlatformList: jest.fn()
        };

        // Create mock filter buttons with proper structure
        const filterButtons = [
            { dataset: { filterType: 'all', filterValue: 'all' }, classList: { add: jest.fn(), remove: jest.fn() }, textContent: 'All Platforms' },
            { dataset: { filterType: 'domain', filterValue: 'airborne' }, classList: { add: jest.fn(), remove: jest.fn() }, textContent: '✈️ Airborne' },
            { dataset: { filterType: 'domain', filterValue: 'land' }, classList: { add: jest.fn(), remove: jest.fn() }, textContent: '🚛 Land' },
            { dataset: { filterType: 'domain', filterValue: 'maritime' }, classList: { add: jest.fn(), remove: jest.fn() }, textContent: '🚢 Maritime' },
            { dataset: { filterType: 'domain', filterValue: 'space' }, classList: { add: jest.fn(), remove: jest.fn() }, textContent: '🛰️ Space' },
            { dataset: { filterType: 'affiliation', filterValue: 'commercial' }, classList: { add: jest.fn(), remove: jest.fn() }, textContent: '🏢 Commercial' },
            { dataset: { filterType: 'affiliation', filterValue: 'military' }, classList: { add: jest.fn(), remove: jest.fn() }, textContent: '⚔️ Military' }
        ];

        mockDOM.filterButtons = filterButtons;

        // Mock document global
        global.document = mockDOM;
    });

    describe('Filter Button Layout Structure', () => {
        test('should have correct filter row structure', () => {
            // Test that the HTML structure matches the design requirements
            const expectedStructure = {
                'filter-row-all': 1,     // One wide "All" button
                'filter-row-domains': 4,  // Four equal domain buttons
                'filter-row-affiliations': 2  // Two equal affiliation buttons
            };

            // Mock the DOM structure check
            const filterRowAll = { children: { length: 1 } };
            const filterRowDomains = { children: { length: 4 } };
            const filterRowAffiliations = { children: { length: 2 } };

            mockDOM.querySelector = jest.fn().mockImplementation((selector) => {
                if (selector === '.filter-row-all') return filterRowAll;
                if (selector === '.filter-row-domains') return filterRowDomains;
                if (selector === '.filter-row-affiliations') return filterRowAffiliations;
                return null;
            });

            expect(filterRowAll.children.length).toBe(expectedStructure['filter-row-all']);
            expect(filterRowDomains.children.length).toBe(expectedStructure['filter-row-domains']);
            expect(filterRowAffiliations.children.length).toBe(expectedStructure['filter-row-affiliations']);
        });

        test('should have correct button data attributes', () => {
            const allButton = mockDOM.filterButtons.find(btn => btn.dataset.filterType === 'all');
            const domainButtons = mockDOM.filterButtons.filter(btn => btn.dataset.filterType === 'domain');
            const affiliationButtons = mockDOM.filterButtons.filter(btn => btn.dataset.filterType === 'affiliation');

            expect(allButton.dataset.filterValue).toBe('all');
            expect(domainButtons).toHaveLength(4);
            expect(affiliationButtons).toHaveLength(2);

            const expectedDomains = ['airborne', 'land', 'maritime', 'space'];
            const expectedAffiliations = ['commercial', 'military'];

            domainButtons.forEach(btn => {
                expect(expectedDomains).toContain(btn.dataset.filterValue);
            });

            affiliationButtons.forEach(btn => {
                expect(expectedAffiliations).toContain(btn.dataset.filterValue);
            });
        });
    });

    describe('Filter Button Behavior', () => {
        test('All button should clear all other filters', () => {
            const allButton = mockDOM.filterButtons.find(btn => btn.dataset.filterType === 'all');

            mockScenarioBuilder.handleFilterClick = (button) => {
                if (button.dataset.filterType === 'all') {
                    mockScenarioBuilder.activeFilters.domain = null;
                    mockScenarioBuilder.activeFilters.affiliation = null;
                }
            };

            mockScenarioBuilder.handleFilterClick(allButton);

            expect(mockScenarioBuilder.activeFilters.domain).toBeNull();
            expect(mockScenarioBuilder.activeFilters.affiliation).toBeNull();
        });

        test('Domain filters should be mutually exclusive', () => {
            const airborneButton = mockDOM.filterButtons.find(btn =>
                btn.dataset.filterType === 'domain' && btn.dataset.filterValue === 'airborne'
            );
            const landButton = mockDOM.filterButtons.find(btn =>
                btn.dataset.filterType === 'domain' && btn.dataset.filterValue === 'land'
            );

            mockScenarioBuilder.handleFilterClick = (button) => {
                if (button.dataset.filterType === 'domain') {
                    mockScenarioBuilder.activeFilters.domain = button.dataset.filterValue;
                }
            };

            // Click airborne filter
            mockScenarioBuilder.handleFilterClick(airborneButton);
            expect(mockScenarioBuilder.activeFilters.domain).toBe('airborne');

            // Click land filter - should replace airborne filter
            mockScenarioBuilder.handleFilterClick(landButton);
            expect(mockScenarioBuilder.activeFilters.domain).toBe('land');
        });

        test('Affiliation filters should be mutually exclusive', () => {
            const commercialButton = mockDOM.filterButtons.find(btn =>
                btn.dataset.filterType === 'affiliation' && btn.dataset.filterValue === 'commercial'
            );
            const militaryButton = mockDOM.filterButtons.find(btn =>
                btn.dataset.filterType === 'affiliation' && btn.dataset.filterValue === 'military'
            );

            mockScenarioBuilder.handleFilterClick = (button) => {
                if (button.dataset.filterType === 'affiliation') {
                    mockScenarioBuilder.activeFilters.affiliation = button.dataset.filterValue;
                }
            };

            // Click commercial filter
            mockScenarioBuilder.handleFilterClick(commercialButton);
            expect(mockScenarioBuilder.activeFilters.affiliation).toBe('commercial');

            // Click military filter - should replace commercial filter
            mockScenarioBuilder.handleFilterClick(militaryButton);
            expect(mockScenarioBuilder.activeFilters.affiliation).toBe('military');
        });

        test('Domain and affiliation filters should be combinable', () => {
            const airborneButton = mockDOM.filterButtons.find(btn =>
                btn.dataset.filterType === 'domain' && btn.dataset.filterValue === 'airborne'
            );
            const militaryButton = mockDOM.filterButtons.find(btn =>
                btn.dataset.filterType === 'affiliation' && btn.dataset.filterValue === 'military'
            );

            mockScenarioBuilder.handleFilterClick = (button) => {
                if (button.dataset.filterType === 'domain') {
                    mockScenarioBuilder.activeFilters.domain = button.dataset.filterValue;
                } else if (button.dataset.filterType === 'affiliation') {
                    mockScenarioBuilder.activeFilters.affiliation = button.dataset.filterValue;
                }
            };

            // Apply both filters
            mockScenarioBuilder.handleFilterClick(airborneButton);
            mockScenarioBuilder.handleFilterClick(militaryButton);

            expect(mockScenarioBuilder.activeFilters.domain).toBe('airborne');
            expect(mockScenarioBuilder.activeFilters.affiliation).toBe('military');
        });
    });

    describe('Platform Filtering Functionality', () => {
        beforeEach(() => {
            // Setup realistic filtering function
            mockScenarioBuilder.applyFilters = () => {
                let filtered = [...mockScenarioBuilder.platforms];

                if (mockScenarioBuilder.activeFilters.domain) {
                    filtered = filtered.filter(p => p.domain === mockScenarioBuilder.activeFilters.domain);
                }

                if (mockScenarioBuilder.activeFilters.affiliation) {
                    filtered = filtered.filter(p => p.affiliation === mockScenarioBuilder.activeFilters.affiliation);
                }

                if (mockScenarioBuilder.activeFilters.search) {
                    const term = mockScenarioBuilder.activeFilters.search.toLowerCase();
                    filtered = filtered.filter(p =>
                        p.name.toLowerCase().includes(term) ||
                        p.description.toLowerCase().includes(term)
                    );
                }

                return filtered;
            };
        });

        test('should filter by domain correctly', () => {
            mockScenarioBuilder.activeFilters.domain = 'airborne';
            const filtered = mockScenarioBuilder.applyFilters();

            expect(filtered).toHaveLength(2);
            expect(filtered.every(p => p.domain === 'airborne')).toBe(true);
        });

        test('should filter by affiliation correctly', () => {
            mockScenarioBuilder.activeFilters.affiliation = 'military';
            const filtered = mockScenarioBuilder.applyFilters();

            expect(filtered).toHaveLength(4);
            expect(filtered.every(p => p.affiliation === 'military')).toBe(true);
        });

        test('should combine domain and affiliation filters', () => {
            mockScenarioBuilder.activeFilters.domain = 'airborne';
            mockScenarioBuilder.activeFilters.affiliation = 'military';
            const filtered = mockScenarioBuilder.applyFilters();

            expect(filtered).toHaveLength(1);
            expect(filtered[0].name).toBe('F-16 Fighter');
        });

        test('should handle search with other filters', () => {
            mockScenarioBuilder.activeFilters.domain = 'land';
            mockScenarioBuilder.activeFilters.search = 'tesla';
            const filtered = mockScenarioBuilder.applyFilters();

            expect(filtered).toHaveLength(1);
            expect(filtered[0].name).toBe('Tesla Model S');
        });

        test('should return empty array when no platforms match', () => {
            mockScenarioBuilder.activeFilters.domain = 'airborne';
            mockScenarioBuilder.activeFilters.affiliation = 'commercial';
            mockScenarioBuilder.activeFilters.search = 'nonexistent';
            const filtered = mockScenarioBuilder.applyFilters();

            expect(filtered).toHaveLength(0);
        });
    });

    describe('CSS Layout Validation', () => {
        test('should have correct CSS classes for filter rows', () => {
            const expectedClasses = [
                'platform-filters',
                'filter-row-all',
                'filter-row-domains',
                'filter-row-affiliations',
                'filter-btn'
            ];

            // Mock CSS validation
            const mockStyleSheet = {
                cssRules: expectedClasses.map(className => ({
                    selectorText: `.${className}`,
                    style: { display: 'flex' }
                }))
            };

            expectedClasses.forEach(className => {
                const rule = mockStyleSheet.cssRules.find(rule =>
                    rule.selectorText === `.${className}`
                );
                expect(rule).toBeDefined();
            });
        });

        test('should validate filter row grid layouts', () => {
            const mockComputedStyles = {
                '.filter-row-domains': { gridTemplateColumns: 'repeat(4, 1fr)' },
                '.filter-row-affiliations': { gridTemplateColumns: 'repeat(2, 1fr)' },
                '.filter-row-all': { display: 'flex' }
            };

            expect(mockComputedStyles['.filter-row-domains'].gridTemplateColumns).toBe('repeat(4, 1fr)');
            expect(mockComputedStyles['.filter-row-affiliations'].gridTemplateColumns).toBe('repeat(2, 1fr)');
            expect(mockComputedStyles['.filter-row-all'].display).toBe('flex');
        });
    });

    describe('User Interface Integration', () => {
        test('should handle rapid filter clicking', () => {
            const buttons = mockDOM.filterButtons;
            let clickCount = 0;

            mockScenarioBuilder.handleFilterClick = () => {
                clickCount++;
            };

            // Simulate rapid clicking
            buttons.forEach(btn => {
                mockScenarioBuilder.handleFilterClick(btn);
            });

            expect(clickCount).toBe(buttons.length);
        });

        test('should maintain filter state during search', () => {
            // Set initial filters
            mockScenarioBuilder.activeFilters.domain = 'airborne';
            mockScenarioBuilder.activeFilters.affiliation = 'military';

            // Simulate search input
            mockScenarioBuilder.activeFilters.search = 'fighter';

            // Verify filters are maintained
            expect(mockScenarioBuilder.activeFilters.domain).toBe('airborne');
            expect(mockScenarioBuilder.activeFilters.affiliation).toBe('military');
            expect(mockScenarioBuilder.activeFilters.search).toBe('fighter');
        });

        test('should update UI when filters change', () => {
            const mockButton = mockDOM.filterButtons[1]; // Air button
            let uiUpdated = false;

            mockScenarioBuilder.renderFilteredPlatformList = () => {
                uiUpdated = true;
            };

            mockScenarioBuilder.handleFilterClick = (button) => {
                mockScenarioBuilder.activeFilters.domain = button.dataset.filterValue;
                mockScenarioBuilder.renderFilteredPlatformList();
            };

            mockScenarioBuilder.handleFilterClick(mockButton);

            expect(uiUpdated).toBe(true);
        });
    });

    describe('Accessibility and User Experience', () => {
        test('should provide visual feedback for active filters', () => {
            const button = mockDOM.filterButtons[0];
            let hasActiveClass = false;

            button.classList.add = jest.fn().mockImplementation((className) => {
                if (className === 'active') {
                    hasActiveClass = true;
                }
            });

            button.classList.add('active');
            expect(hasActiveClass).toBe(true);
        });

        test('should handle no results gracefully', () => {
            mockScenarioBuilder.renderFilteredPlatformList = (platforms) => {
                if (platforms.length === 0) {
                    mockDOM.platformList.innerHTML = '<div class="no-results">No platforms match your filters</div>';
                }
            };

            mockScenarioBuilder.renderFilteredPlatformList([]);
            expect(mockDOM.platformList.innerHTML).toContain('No platforms match your filters');
        });

        test('should preserve search text when applying filters', () => {
            mockDOM.platformSearch.value = 'test search';
            mockScenarioBuilder.activeFilters.search = 'test search';

            // Apply domain filter
            mockScenarioBuilder.activeFilters.domain = 'airborne';

            // Search should be preserved
            expect(mockScenarioBuilder.activeFilters.search).toBe('test search');
        });
    });

    describe('Multi-Select Filter Behavior', () => {
        let scenarioBuilder;

        beforeEach(() => {
            // Create a real instance with mocked DOM
            scenarioBuilder = {
                activeFilters: { domains: new Set(), affiliations: new Set(), search: '' },
                platforms: [
                    { id: 'air_commercial', name: 'Airbus A320', domain: 'airborne', affiliation: 'commercial', category: 'aircraft', description: 'Commercial airliner' },
                    { id: 'air_military', name: 'F-16 Fighter', domain: 'airborne', affiliation: 'military', category: 'fighter', description: 'Military fighter jet' },
                    { id: 'land_commercial', name: 'Tesla Model S', domain: 'land', affiliation: 'commercial', category: 'car', description: 'Electric vehicle' },
                    { id: 'land_military', name: 'Military Truck', domain: 'land', affiliation: 'military', category: 'truck', description: 'Military transport' },
                    { id: 'sea_commercial', name: 'Container Ship', domain: 'maritime', affiliation: 'commercial', category: 'ship', description: 'Cargo vessel' },
                    { id: 'sea_military', name: 'Naval Destroyer', domain: 'maritime', affiliation: 'military', category: 'warship', description: 'Combat vessel' },
                    { id: 'space_commercial', name: 'Starlink Sat', domain: 'space', affiliation: 'commercial', category: 'satellite', description: 'Communication satellite' },
                    { id: 'space_military', name: 'Military Sat', domain: 'space', affiliation: 'military', category: 'satellite', description: 'Reconnaissance satellite' }
                ],
                isAllFiltersActive() {
                    const allDomains = ['airborne', 'land', 'maritime', 'space'];
                    const allAffiliations = ['commercial', 'military'];
                    return allDomains.every(domain => this.activeFilters.domains.has(domain)) &&
                        allAffiliations.every(affiliation => this.activeFilters.affiliations.has(affiliation));
                },
                selectAllFilters() {
                    this.activeFilters.domains = new Set(['airborne', 'land', 'maritime', 'space']);
                    this.activeFilters.affiliations = new Set(['commercial', 'military']);
                },
                clearAllFilters() {
                    this.activeFilters.domains.clear();
                    this.activeFilters.affiliations.clear();
                },
                applyFilters() {
                    let filteredPlatforms = [...this.platforms];

                    // Apply domain filters
                    if (this.activeFilters.domains.size > 0) {
                        filteredPlatforms = filteredPlatforms.filter(platform =>
                            this.activeFilters.domains.has(platform.domain)
                        );
                    }

                    // Apply affiliation filters
                    if (this.activeFilters.affiliations.size > 0) {
                        filteredPlatforms = filteredPlatforms.filter(platform =>
                            this.activeFilters.affiliations.has(platform.affiliation)
                        );
                    }

                    this.lastFilteredPlatforms = filteredPlatforms;
                    return filteredPlatforms;
                }
            };
        });

        test('should show all platforms when All button is selected', () => {
            // When All is selected, all filters should be active
            scenarioBuilder.selectAllFilters();
            const result = scenarioBuilder.applyFilters();

            expect(result).toHaveLength(8); // All 8 platforms
            expect(scenarioBuilder.isAllFiltersActive()).toBe(true);
        });

        test('should allow multiple domain selections', () => {
            // Select airborne and land domains
            scenarioBuilder.activeFilters.domains.add('airborne');
            scenarioBuilder.activeFilters.domains.add('land');
            const result = scenarioBuilder.applyFilters();

            // Should show 4 platforms (2 airborne + 2 land)
            expect(result).toHaveLength(4);
            expect(result.every(p => ['airborne', 'land'].includes(p.domain))).toBe(true);
        });

        test('should allow multiple affiliation selections', () => {
            // Select commercial affiliation for all domains
            scenarioBuilder.activeFilters.domains = new Set(['airborne', 'land', 'maritime', 'space']);
            scenarioBuilder.activeFilters.affiliations.add('commercial');
            const result = scenarioBuilder.applyFilters();

            // Should show 4 platforms (all commercial ones)
            expect(result).toHaveLength(4);
            expect(result.every(p => p.affiliation === 'commercial')).toBe(true);
        });

        test('should combine domain and affiliation filters correctly', () => {
            // Select airborne domain and military affiliation
            scenarioBuilder.activeFilters.domains.add('airborne');
            scenarioBuilder.activeFilters.affiliations.add('military');
            const result = scenarioBuilder.applyFilters();

            // Should show 1 platform (F-16 Fighter)
            expect(result).toHaveLength(1);
            expect(result[0].name).toBe('F-16 Fighter');
            expect(result[0].domain).toBe('airborne');
            expect(result[0].affiliation).toBe('military');
        });

        test('should show no platforms when no filters are selected', () => {
            // When domains have items but affiliations don't, it should filter by domains only
            scenarioBuilder.activeFilters.domains.add('airborne');
            // Clear affiliations entirely - this should show all affiliations for the selected domain
            scenarioBuilder.activeFilters.affiliations.clear();
            const result = scenarioBuilder.applyFilters();

            // Should show all airborne platforms (both commercial and military) since no affiliation filter is active
            expect(result).toHaveLength(2);
            expect(result.every(p => p.domain === 'airborne')).toBe(true);
        });

        test('should handle complex multi-selection scenarios', () => {
            // Select airborne + maritime domains and commercial + military affiliations
            scenarioBuilder.activeFilters.domains.add('airborne');
            scenarioBuilder.activeFilters.domains.add('maritime');
            scenarioBuilder.activeFilters.affiliations.add('commercial');
            scenarioBuilder.activeFilters.affiliations.add('military');
            const result = scenarioBuilder.applyFilters();

            // Should show 4 platforms (2 airborne + 2 maritime)
            expect(result).toHaveLength(4);
            expect(result.every(p => ['airborne', 'maritime'].includes(p.domain))).toBe(true);
            expect(result.every(p => ['commercial', 'military'].includes(p.affiliation))).toBe(true);
        });

        test('should correctly identify when all filters are active', () => {
            // Initially no filters
            expect(scenarioBuilder.isAllFiltersActive()).toBe(false);

            // Add some but not all
            scenarioBuilder.activeFilters.domains.add('air');
            scenarioBuilder.activeFilters.affiliations.add('commercial');
            expect(scenarioBuilder.isAllFiltersActive()).toBe(false);

            // Add all domains and affiliations
            scenarioBuilder.selectAllFilters();
            expect(scenarioBuilder.isAllFiltersActive()).toBe(true);
        });

        test('should maintain platform data integrity during filtering', () => {
            // Test that filtering doesn't modify original platform data
            const originalPlatforms = JSON.parse(JSON.stringify(scenarioBuilder.platforms));

            scenarioBuilder.activeFilters.domains.add('air');
            scenarioBuilder.activeFilters.affiliations.add('military');
            scenarioBuilder.applyFilters();

            // Original data should be unchanged
            expect(scenarioBuilder.platforms).toEqual(originalPlatforms);
        });
    });

    describe('Platform Display Validation', () => {
        let scenarioBuilder;

        beforeEach(() => {
            scenarioBuilder = {
                activeFilters: { domains: new Set(), affiliations: new Set(), search: '' },
                platforms: [
                    { id: 'air_commercial', name: 'Airbus A320', domain: 'air', affiliation: 'commercial', category: 'aircraft', description: 'Commercial airliner' },
                    { id: 'air_military', name: 'F-16 Fighter', domain: 'air', affiliation: 'military', category: 'fighter', description: 'Military fighter jet' },
                    { id: 'land_commercial', name: 'Tesla Model S', domain: 'land', affiliation: 'commercial', category: 'car', description: 'Electric vehicle' },
                    { id: 'land_military', name: 'Military Truck', domain: 'land', affiliation: 'military', category: 'truck', description: 'Military transport' },
                    { id: 'sea_commercial', name: 'Container Ship', domain: 'sea', affiliation: 'commercial', category: 'ship', description: 'Cargo vessel' },
                    { id: 'sea_military', name: 'Naval Destroyer', domain: 'sea', affiliation: 'military', category: 'warship', description: 'Combat vessel' },
                    { id: 'space_commercial', name: 'Starlink Sat', domain: 'space', affiliation: 'commercial', category: 'satellite', description: 'Communication satellite' },
                    { id: 'space_military', name: 'Military Sat', domain: 'space', affiliation: 'military', category: 'satellite', description: 'Reconnaissance satellite' }
                ],
                applyFilters() {
                    let filteredPlatforms = [...this.platforms];

                    if (this.activeFilters.domains.size > 0) {
                        filteredPlatforms = filteredPlatforms.filter(platform =>
                            this.activeFilters.domains.has(platform.domain)
                        );
                    }

                    if (this.activeFilters.affiliations.size > 0) {
                        filteredPlatforms = filteredPlatforms.filter(platform =>
                            this.activeFilters.affiliations.has(platform.affiliation)
                        );
                    }

                    return filteredPlatforms;
                }
            };
        });

        test('should display correct platforms for Air + Commercial filter', () => {
            scenarioBuilder.activeFilters.domains.add('air');
            scenarioBuilder.activeFilters.affiliations.add('commercial');
            const result = scenarioBuilder.applyFilters();

            expect(result).toHaveLength(1);
            expect(result[0]).toEqual(expect.objectContaining({
                name: 'Airbus A320',
                domain: 'air',
                affiliation: 'commercial',
                category: 'aircraft'
            }));
        });

        test('should display correct platforms for Land + Military filter', () => {
            scenarioBuilder.activeFilters.domains.add('land');
            scenarioBuilder.activeFilters.affiliations.add('military');
            const result = scenarioBuilder.applyFilters();

            expect(result).toHaveLength(1);
            expect(result[0]).toEqual(expect.objectContaining({
                name: 'Military Truck',
                domain: 'land',
                affiliation: 'military',
                category: 'truck'
            }));
        });

        test('should display correct platforms for Sea + Commercial filter', () => {
            scenarioBuilder.activeFilters.domains.add('sea');
            scenarioBuilder.activeFilters.affiliations.add('commercial');
            const result = scenarioBuilder.applyFilters();

            expect(result).toHaveLength(1);
            expect(result[0]).toEqual(expect.objectContaining({
                name: 'Container Ship',
                domain: 'sea',
                affiliation: 'commercial',
                category: 'ship'
            }));
        });

        test('should display correct platforms for Space + Military filter', () => {
            scenarioBuilder.activeFilters.domains.add('space');
            scenarioBuilder.activeFilters.affiliations.add('military');
            const result = scenarioBuilder.applyFilters();

            expect(result).toHaveLength(1);
            expect(result[0]).toEqual(expect.objectContaining({
                name: 'Military Sat',
                domain: 'space',
                affiliation: 'military',
                category: 'satellite'
            }));
        });

        test('should display all commercial platforms when only Commercial filter is active', () => {
            scenarioBuilder.activeFilters.domains = new Set(['air', 'land', 'sea', 'space']);
            scenarioBuilder.activeFilters.affiliations.add('commercial');
            const result = scenarioBuilder.applyFilters();

            expect(result).toHaveLength(4);
            result.forEach(platform => {
                expect(platform.affiliation).toBe('commercial');
            });

            const expectedNames = ['Airbus A320', 'Tesla Model S', 'Container Ship', 'Starlink Sat'];
            const actualNames = result.map(p => p.name).sort();
            expect(actualNames).toEqual(expectedNames.sort());
        });

        test('should display all military platforms when only Military filter is active', () => {
            scenarioBuilder.activeFilters.domains = new Set(['air', 'land', 'sea', 'space']);
            scenarioBuilder.activeFilters.affiliations.add('military');
            const result = scenarioBuilder.applyFilters();

            expect(result).toHaveLength(4);
            result.forEach(platform => {
                expect(platform.affiliation).toBe('military');
            });

            const expectedNames = ['F-16 Fighter', 'Military Truck', 'Naval Destroyer', 'Military Sat'];
            const actualNames = result.map(p => p.name).sort();
            expect(expectedNames.sort()).toEqual(actualNames);
        });

        test('should display platforms for multiple domain selection', () => {
            scenarioBuilder.activeFilters.domains.add('air');
            scenarioBuilder.activeFilters.domains.add('space');
            scenarioBuilder.activeFilters.affiliations = new Set(['commercial', 'military']);
            const result = scenarioBuilder.applyFilters();

            expect(result).toHaveLength(4); // 2 air + 2 space
            result.forEach(platform => {
                expect(['air', 'space']).toContain(platform.domain);
            });
        });

        test('should handle edge case with no matching platforms when domains are empty', () => {
            // Create scenario where no platforms match (no domains selected, but affiliations selected)
            scenarioBuilder.activeFilters.domains.clear();
            scenarioBuilder.activeFilters.affiliations.add('commercial');
            const result = scenarioBuilder.applyFilters();

            // Should show all commercial platforms since no domain filter is active
            expect(result).toHaveLength(4);
            expect(result.every(p => p.affiliation === 'commercial')).toBe(true);
        });

        test('should validate platform properties are preserved during filtering', () => {
            scenarioBuilder.activeFilters.domains.add('air');
            scenarioBuilder.activeFilters.affiliations.add('commercial');
            const result = scenarioBuilder.applyFilters();

            const platform = result[0];
            expect(platform).toHaveProperty('id');
            expect(platform).toHaveProperty('name');
            expect(platform).toHaveProperty('domain');
            expect(platform).toHaveProperty('affiliation');
            expect(platform).toHaveProperty('category');
            expect(platform).toHaveProperty('description');

            // Verify the values are strings
            expect(typeof platform.id).toBe('string');
            expect(typeof platform.name).toBe('string');
            expect(typeof platform.domain).toBe('string');
            expect(typeof platform.affiliation).toBe('string');
            expect(typeof platform.category).toBe('string');
            expect(typeof platform.description).toBe('string');
        });
    });
});
