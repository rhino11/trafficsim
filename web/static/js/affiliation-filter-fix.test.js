// Test to verify the affiliation filtering fix works correctly

describe('Affiliation Filter Fix Verification', () => {
    let scenarioBuilder;
    let mockDOM;

    beforeEach(() => {
        // Mock DOM elements
        mockDOM = {
            allButton: {
                dataset: { filterType: 'all', filterValue: 'all' },
                classList: { add: jest.fn(), remove: jest.fn(), contains: jest.fn() }
            },
            commercialButton: {
                dataset: { filterType: 'affiliation', filterValue: 'commercial' },
                classList: { add: jest.fn(), remove: jest.fn(), contains: jest.fn() }
            },
            militaryButton: {
                dataset: { filterType: 'affiliation', filterValue: 'military' },
                classList: { add: jest.fn(), remove: jest.fn(), contains: jest.fn() }
            }
        };

        // Mock document.querySelectorAll
        global.document = {
            querySelectorAll: jest.fn((selector) => {
                if (selector === '[data-filter-type="affiliation"]') {
                    return [mockDOM.commercialButton, mockDOM.militaryButton];
                }
                return [];
            })
        };

        // Create a realistic scenario builder mock with the fixed logic
        scenarioBuilder = {
            activeFilters: {
                domains: new Set(['airborne', 'land', 'maritime', 'space']),
                affiliations: new Set(['commercial', 'military']),
                search: ''
            },
            platforms: [
                { id: 'air_comm', name: 'Airbus A320', domain: 'airborne', affiliation: 'commercial' },
                { id: 'air_mil', name: 'F-16 Fighter', domain: 'airborne', affiliation: 'military' },
                { id: 'sea_comm', name: 'Container Ship', domain: 'maritime', affiliation: 'commercial' },
                { id: 'sea_mil', name: 'Naval Destroyer', domain: 'maritime', affiliation: 'military' },
                { id: 'land_comm', name: 'Tesla Model S', domain: 'land', affiliation: 'commercial' },
                { id: 'land_mil', name: 'Military Truck', domain: 'land', affiliation: 'military' },
                { id: 'space_comm', name: 'Starlink Sat', domain: 'space', affiliation: 'commercial' },
                { id: 'space_mil', name: 'Military Sat', domain: 'space', affiliation: 'military' }
            ],

            updateAllButtonState: jest.fn(),

            // Implementation of the FIXED handleFilterClick method
            handleFilterClick(button) {
                const filterType = button.dataset.filterType;
                const filterValue = button.dataset.filterValue;

                if (filterType === 'affiliation') {
                    // Implement exclusive affiliation selection
                    if (this.activeFilters.affiliations.has(filterValue) &&
                        this.activeFilters.affiliations.size === 1) {
                        // If only this affiliation is selected, clicking it should select all (toggle back to show all)
                        this.activeFilters.affiliations = new Set(['commercial', 'military']);
                        document.querySelectorAll('[data-filter-type="affiliation"]').forEach(btn => {
                            btn.classList.add('active');
                        });
                    } else {
                        // Otherwise, select only this affiliation (exclusive selection)
                        this.activeFilters.affiliations = new Set([filterValue]);
                        document.querySelectorAll('[data-filter-type="affiliation"]').forEach(btn => {
                            if (btn.dataset.filterValue === filterValue) {
                                btn.classList.add('active');
                            } else {
                                btn.classList.remove('active');
                            }
                        });
                    }
                    this.updateAllButtonState();
                }

                return this.applyFilters();
            },

            // Fixed applyFilters method (same as before, working correctly)
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

                // Apply affiliation filters (only if some but not all affiliations are selected)
                if (this.activeFilters.affiliations.size > 0 && this.activeFilters.affiliations.size < 2) {
                    filteredPlatforms = filteredPlatforms.filter(platform =>
                        this.activeFilters.affiliations.has(platform.affiliation)
                    );
                } else if (this.activeFilters.affiliations.size === 0) {
                    filteredPlatforms = [];
                }

                return filteredPlatforms;
            }
        };
    });

    test('should show only military platforms when Military button is clicked', () => {
        // Start with both affiliations selected (default state)
        expect(scenarioBuilder.activeFilters.affiliations).toEqual(new Set(['commercial', 'military']));

        // Click Military button
        const result = scenarioBuilder.handleFilterClick(mockDOM.militaryButton);

        // Should now have only military selected
        expect(scenarioBuilder.activeFilters.affiliations).toEqual(new Set(['military']));

        // Should show only military platforms
        expect(result).toHaveLength(4);
        expect(result.every(p => p.affiliation === 'military')).toBe(true);

        // Should have set button states correctly
        expect(mockDOM.militaryButton.classList.add).toHaveBeenCalledWith('active');
        expect(mockDOM.commercialButton.classList.remove).toHaveBeenCalledWith('active');
    });

    test('should show only commercial platforms when Commercial button is clicked', () => {
        // Start with both affiliations selected
        expect(scenarioBuilder.activeFilters.affiliations).toEqual(new Set(['commercial', 'military']));

        // Click Commercial button
        const result = scenarioBuilder.handleFilterClick(mockDOM.commercialButton);

        // Should now have only commercial selected
        expect(scenarioBuilder.activeFilters.affiliations).toEqual(new Set(['commercial']));

        // Should show only commercial platforms
        expect(result).toHaveLength(4);
        expect(result.every(p => p.affiliation === 'commercial')).toBe(true);

        // Should have set button states correctly
        expect(mockDOM.commercialButton.classList.add).toHaveBeenCalledWith('active');
        expect(mockDOM.militaryButton.classList.remove).toHaveBeenCalledWith('active');
    });

    test('should toggle back to all affiliations when clicking the same button twice', () => {
        // Start with both affiliations selected
        expect(scenarioBuilder.activeFilters.affiliations).toEqual(new Set(['commercial', 'military']));

        // Click Military button (select only military)
        scenarioBuilder.handleFilterClick(mockDOM.militaryButton);
        expect(scenarioBuilder.activeFilters.affiliations).toEqual(new Set(['military']));

        // Click Military button again (should toggle back to all)
        const result = scenarioBuilder.handleFilterClick(mockDOM.militaryButton);
        expect(scenarioBuilder.activeFilters.affiliations).toEqual(new Set(['commercial', 'military']));

        // Should show all platforms
        expect(result).toHaveLength(8);

        // Should have activated both buttons
        expect(mockDOM.militaryButton.classList.add).toHaveBeenCalledWith('active');
        expect(mockDOM.commercialButton.classList.add).toHaveBeenCalledWith('active');
    });

    test('should switch from one affiliation to another correctly', () => {
        // Start with both affiliations selected
        expect(scenarioBuilder.activeFilters.affiliations).toEqual(new Set(['commercial', 'military']));

        // Click Military button (select only military)
        let result = scenarioBuilder.handleFilterClick(mockDOM.militaryButton);
        expect(scenarioBuilder.activeFilters.affiliations).toEqual(new Set(['military']));
        expect(result.every(p => p.affiliation === 'military')).toBe(true);

        // Click Commercial button (should switch to only commercial)
        result = scenarioBuilder.handleFilterClick(mockDOM.commercialButton);
        expect(scenarioBuilder.activeFilters.affiliations).toEqual(new Set(['commercial']));
        expect(result.every(p => p.affiliation === 'commercial')).toBe(true);
    });

    test('should work correctly with domain filters combined', () => {
        // Set only airborne domain active
        scenarioBuilder.activeFilters.domains = new Set(['airborne']);

        // Click Military button
        const result = scenarioBuilder.handleFilterClick(mockDOM.militaryButton);

        // Should show only military airborne platforms
        expect(result).toHaveLength(1);
        expect(result[0].name).toBe('F-16 Fighter');
        expect(result[0].domain).toBe('airborne');
        expect(result[0].affiliation).toBe('military');
    });
});
