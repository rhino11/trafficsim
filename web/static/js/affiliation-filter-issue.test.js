// Test to identify and fix the affiliation filtering issue
// This test captures the current broken behavior and defines the expected fix

describe('Affiliation Filter Issue Analysis', () => {
    let mockScenarioBuilder;
    let mockFilterButtons;

    beforeEach(() => {
        // Mock the DOM elements
        mockFilterButtons = {
            commercial: {
                dataset: { filterType: 'affiliation', filterValue: 'commercial' },
                classList: { add: jest.fn(), remove: jest.fn(), contains: jest.fn() },
                textContent: '🏢 Commercial'
            },
            military: {
                dataset: { filterType: 'affiliation', filterValue: 'military' },
                classList: { add: jest.fn(), remove: jest.fn(), contains: jest.fn() },
                textContent: '⚔️ Military'
            }
        };

        // Mock scenario builder with test platforms
        mockScenarioBuilder = {
            activeFilters: {
                domains: new Set(['airborne', 'land', 'maritime', 'space']), // All domains active
                affiliations: new Set(['commercial', 'military']), // All affiliations active initially
                search: ''
            },
            platforms: [
                { id: 'air_commercial', name: 'Airbus A320', domain: 'airborne', affiliation: 'commercial' },
                { id: 'air_military', name: 'F-16 Fighter', domain: 'airborne', affiliation: 'military' },
                { id: 'sea_commercial', name: 'Container Ship', domain: 'maritime', affiliation: 'commercial' },
                { id: 'sea_military', name: 'Naval Destroyer', domain: 'maritime', affiliation: 'military' },
                { id: 'land_commercial', name: 'Tesla Model S', domain: 'land', affiliation: 'commercial' },
                { id: 'land_military', name: 'Military Truck', domain: 'land', affiliation: 'military' },
                { id: 'space_commercial', name: 'Starlink Sat', domain: 'space', affiliation: 'commercial' },
                { id: 'space_military', name: 'Military Sat', domain: 'space', affiliation: 'military' }
            ],

            // Current broken applyFilters implementation (this represents the current buggy logic)
            applyFilters() {
                let filteredPlatforms = [...this.platforms];

                // Apply domain filters (only if some but not all domains are selected)
                if (this.activeFilters.domains.size > 0 && this.activeFilters.domains.size < 4) {
                    filteredPlatforms = filteredPlatforms.filter(platform =>
                        this.activeFilters.domains.has(platform.domain)
                    );
                } else if (this.activeFilters.domains.size === 0) {
                    filteredPlatforms = [];
                }

                // BROKEN LOGIC: Apply affiliation filters (only if some but not all affiliations are selected)
                if (this.activeFilters.affiliations.size > 0 && this.activeFilters.affiliations.size < 2) {
                    filteredPlatforms = filteredPlatforms.filter(platform =>
                        this.activeFilters.affiliations.has(platform.affiliation)
                    );
                } else if (this.activeFilters.affiliations.size === 0) {
                    filteredPlatforms = [];
                }

                return filteredPlatforms;
            },

            // Mock the filter click handling
            handleFilterClick(button) {
                const filterType = button.dataset.filterType;
                const filterValue = button.dataset.filterValue;

                if (filterType === 'affiliation') {
                    if (this.activeFilters.affiliations.has(filterValue)) {
                        this.activeFilters.affiliations.delete(filterValue);
                        button.classList.remove('active');
                    } else {
                        this.activeFilters.affiliations.add(filterValue);
                        button.classList.add('active');
                    }
                }
                return this.applyFilters();
            }
        };
    });

    describe('Current Broken Behavior', () => {
        test('should demonstrate the affiliation filtering bug', () => {
            // PROBLEM: When both affiliations are active (size === 2), no filtering is applied
            // This means clicking Military or Commercial buttons doesn't filter the list

            // Start with all affiliations active (this is the default state)
            expect(mockScenarioBuilder.activeFilters.affiliations.size).toBe(2);

            // When all affiliations are active, applyFilters should show all platforms
            let result = mockScenarioBuilder.applyFilters();
            expect(result).toHaveLength(8); // All platforms shown

            // Now simulate clicking the Military button to DESELECT commercial
            // This should leave only military active
            mockScenarioBuilder.activeFilters.affiliations.delete('commercial');
            expect(mockScenarioBuilder.activeFilters.affiliations.size).toBe(1);

            // This should work and show only military platforms
            result = mockScenarioBuilder.applyFilters();
            expect(result).toHaveLength(4); // Only military platforms
            expect(result.every(p => p.affiliation === 'military')).toBe(true);
        });

        test('should identify the specific issue: buttons do not behave as expected', () => {
            // ISSUE: The current logic treats "all selected" as "no filter"
            // But users expect "Military" button click to show ONLY military platforms

            // Start with both selected (All state)
            mockScenarioBuilder.activeFilters.affiliations = new Set(['commercial', 'military']);

            // User clicks "Military" button expecting to see ONLY military platforms
            // But the current implementation requires DESELECTING the other first

            // Current broken behavior: clicking Military when both are selected does nothing useful
            const militaryBtn = mockFilterButtons.military;

            // Simulate clicking Military button
            mockScenarioBuilder.handleFilterClick(militaryBtn);

            // This should result in ONLY military being selected
            expect(mockScenarioBuilder.activeFilters.affiliations.has('military')).toBe(true);
            expect(mockScenarioBuilder.activeFilters.affiliations.has('commercial')).toBe(false);

            // And should show only military platforms
            const result = mockScenarioBuilder.applyFilters();
            expect(result).toHaveLength(4);
            expect(result.every(p => p.affiliation === 'military')).toBe(true);
        });
    });

    describe('Expected Fixed Behavior', () => {
        test('should implement exclusive affiliation selection', () => {
            // SOLUTION: Affiliation buttons should be mutually exclusive
            // Clicking "Military" should select ONLY military, deselecting commercial
            // Clicking "Commercial" should select ONLY commercial, deselecting military

            const mockFixedScenarioBuilder = {
                ...mockScenarioBuilder,

                // FIXED handleFilterClick implementation
                handleFilterClick(button) {
                    const filterType = button.dataset.filterType;
                    const filterValue = button.dataset.filterValue;

                    if (filterType === 'affiliation') {
                        if (this.activeFilters.affiliations.has(filterValue) &&
                            this.activeFilters.affiliations.size === 1) {
                            // If only this affiliation is selected, clicking it should select all
                            this.activeFilters.affiliations = new Set(['commercial', 'military']);
                        } else {
                            // Otherwise, select only this affiliation (exclusive selection)
                            this.activeFilters.affiliations = new Set([filterValue]);
                        }

                        // Update button states
                        Object.values(mockFilterButtons).forEach(btn => {
                            if (btn.dataset.filterType === 'affiliation') {
                                if (this.activeFilters.affiliations.has(btn.dataset.filterValue)) {
                                    btn.classList.add('active');
                                } else {
                                    btn.classList.remove('active');
                                }
                            }
                        });
                    }
                    return this.applyFilters();
                },

                // FIXED applyFilters implementation
                applyFilters() {
                    let filteredPlatforms = [...this.platforms];

                    // Apply domain filters
                    if (this.activeFilters.domains.size > 0 && this.activeFilters.domains.size < 4) {
                        filteredPlatforms = filteredPlatforms.filter(platform =>
                            this.activeFilters.domains.has(platform.domain)
                        );
                    } else if (this.activeFilters.domains.size === 0) {
                        filteredPlatforms = [];
                    }

                    // FIXED: Apply affiliation filters - always filter if not all selected
                    if (this.activeFilters.affiliations.size > 0 && this.activeFilters.affiliations.size < 2) {
                        filteredPlatforms = filteredPlatforms.filter(platform =>
                            this.activeFilters.affiliations.has(platform.affiliation)
                        );
                    }
                    // Note: When size === 2 (all selected), no filtering is applied (show all)

                    return filteredPlatforms;
                }
            };

            // Test the fixed behavior
            mockFixedScenarioBuilder.activeFilters.affiliations = new Set(['commercial', 'military']);

            // Click Military button
            const militaryBtn = mockFilterButtons.military;
            let result = mockFixedScenarioBuilder.handleFilterClick(militaryBtn);

            // Should now show only military platforms
            expect(mockFixedScenarioBuilder.activeFilters.affiliations).toEqual(new Set(['military']));
            expect(result).toHaveLength(4);
            expect(result.every(p => p.affiliation === 'military')).toBe(true);

            // Click Commercial button
            const commercialBtn = mockFilterButtons.commercial;
            result = mockFixedScenarioBuilder.handleFilterClick(commercialBtn);

            // Should now show only commercial platforms
            expect(mockFixedScenarioBuilder.activeFilters.affiliations).toEqual(new Set(['commercial']));
            expect(result).toHaveLength(4);
            expect(result.every(p => p.affiliation === 'commercial')).toBe(true);

            // Click Military again (when only Military is selected)
            result = mockFixedScenarioBuilder.handleFilterClick(militaryBtn);

            // Should select all affiliations (toggle back to show all)
            expect(mockFixedScenarioBuilder.activeFilters.affiliations).toEqual(new Set(['commercial', 'military']));
            expect(result).toHaveLength(8);
        });
    });
});
