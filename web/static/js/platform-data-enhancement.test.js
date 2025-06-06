/**
 * Test suite for MIL-STD-2525D platform data enhancement
 * Tests enhanced platform data display and editing capabilities
 */

describe('Platform Data Enhancement for MIL-STD-2525D', () => {
    let scenarioBuilder;
    let mockDOM;

    beforeEach(() => {
        // Create mock DOM environment
        mockDOM = {
            createElement: jest.fn((tag) => ({
                tagName: tag.toUpperCase(),
                className: '',
                innerHTML: '',
                textContent: '',
                style: {},
                addEventListener: jest.fn(),
                appendChild: jest.fn(),
                setAttribute: jest.fn(),
                getAttribute: jest.fn(),
                querySelectorAll: jest.fn(() => []),
                querySelector: jest.fn()
            })),
            getElementById: jest.fn(),
            querySelectorAll: jest.fn(() => []),
            querySelector: jest.fn()
        };

        global.document = mockDOM;

        // Mock scenario builder with platform data
        scenarioBuilder = {
            platforms: [
                {
                    id: 'f16_fighter',
                    name: 'F-16 Fighting Falcon',
                    class: 'F-16C Block 50',
                    category: 'fighter_aircraft',
                    domain: 'airborne',
                    affiliation: 'military',
                    description: 'Multirole fighter aircraft'
                },
                {
                    id: 'civilian_airliner',
                    name: 'Boeing 737',
                    class: 'Boeing 737-800',
                    category: 'commercial_aircraft',
                    domain: 'airborne',
                    affiliation: 'commercial',
                    description: 'Commercial passenger aircraft'
                }
            ],

            // Enhanced platform data methods
            enhancePlatformData: function (platform) {
                const enhanced = { ...platform };

                // Generate symbol code based on platform
                enhanced.symbolCode = this.generateSymbolCode(platform);
                enhanced.unitDesignation = this.generateUnitDesignation(platform);
                enhanced.higherFormation = this.generateHigherFormation(platform);
                enhanced.specifications = this.generateSpecifications(platform);
                enhanced.operational = this.generateOperationalData(platform);

                return enhanced;
            },

            generateSymbolCode: function (platform) {
                const codes = {
                    'f16_fighter': 'a-f-A-M-F',
                    'civilian_airliner': 'a-n-A-C-A'
                };
                return codes[platform.id] || 'a-u-G-U-U';
            },

            generateUnitDesignation: function (platform) {
                const designations = {
                    'f16_fighter': '34th FS',
                    'civilian_airliner': 'UAL1234'
                };
                return designations[platform.id] || 'UNIT-001';
            },

            generateHigherFormation: function (platform) {
                const formations = {
                    'f16_fighter': '388th FW',
                    'civilian_airliner': 'United Airlines'
                };
                return formations[platform.id] || 'Formation';
            },

            generateSpecifications: function (platform) {
                const specs = {
                    'f16_fighter': {
                        length: '15.06 m',
                        wingspan: '9.96 m',
                        height: '5.09 m',
                        maxSpeed: '2,414 km/h',
                        serviceceiling: '15,240 m',
                        range: '4,220 km'
                    },
                    'civilian_airliner': {
                        length: '39.5 m',
                        wingspan: '35.8 m',
                        height: '12.5 m',
                        maxSpeed: '842 km/h',
                        serviceceiling: '12,500 m',
                        range: '5,765 km'
                    }
                };
                return specs[platform.id] || {};
            },

            generateOperationalData: function (platform) {
                const data = {
                    'f16_fighter': {
                        crewSize: 1,
                        fuelLevel: '85%',
                        missionReadiness: 'ready',
                        maintenanceStatus: 'operational'
                    },
                    'civilian_airliner': {
                        crewSize: 6,
                        fuelLevel: '92%',
                        missionReadiness: 'ready',
                        maintenanceStatus: 'operational'
                    }
                };
                return data[platform.id] || {};
            },

            createEnhancedPlatformDisplay: function (platform) {
                const enhanced = this.enhancePlatformData(platform);
                const element = mockDOM.createElement('div');
                element.className = 'platform-enhanced';
                element.innerHTML = `
                    <div class="platform-header">
                        <h4 class="platform-name">${enhanced.name}</h4>
                        <button class="platform-edit-btn">Edit</button>
                    </div>
                    <div class="platform-sections">
                        <div class="platform-section mil-std">
                            <div class="platform-section-title">MIL-STD-2525D</div>
                            <div class="platform-field">
                                <span class="platform-field-label">Symbol Code:</span>
                                <span class="platform-field-value symbol-code">${enhanced.symbolCode}</span>
                            </div>
                            <div class="platform-field">
                                <span class="platform-field-label">Unit:</span>
                                <span class="platform-field-value unit-designation">${enhanced.unitDesignation}</span>
                            </div>
                            <div class="platform-field">
                                <span class="platform-field-label">Formation:</span>
                                <span class="platform-field-value formation">${enhanced.higherFormation}</span>
                            </div>
                        </div>
                    </div>
                `;
                return element;
            }
        };
    });

    describe('Platform Data Enhancement', () => {
        test('should enhance basic platform data with MIL-STD-2525D fields', () => {
            const platform = scenarioBuilder.platforms[0]; // F-16
            const enhanced = scenarioBuilder.enhancePlatformData(platform);

            expect(enhanced.symbolCode).toBe('a-f-A-M-F');
            expect(enhanced.unitDesignation).toBe('34th FS');
            expect(enhanced.higherFormation).toBe('388th FW');
            expect(enhanced.specifications).toBeDefined();
            expect(enhanced.operational).toBeDefined();
        });

        test('should generate correct symbol codes for different platform types', () => {
            const fighter = scenarioBuilder.platforms[0];
            const airliner = scenarioBuilder.platforms[1];

            expect(scenarioBuilder.generateSymbolCode(fighter)).toBe('a-f-A-M-F');
            expect(scenarioBuilder.generateSymbolCode(airliner)).toBe('a-n-A-C-A');
        });

        test('should generate appropriate unit designations', () => {
            const fighter = scenarioBuilder.platforms[0];
            const airliner = scenarioBuilder.platforms[1];

            expect(scenarioBuilder.generateUnitDesignation(fighter)).toBe('34th FS');
            expect(scenarioBuilder.generateUnitDesignation(airliner)).toBe('UAL1234');
        });

        test('should generate specifications for platforms', () => {
            const fighter = scenarioBuilder.platforms[0];
            const specs = scenarioBuilder.generateSpecifications(fighter);

            expect(specs.length).toBe('15.06 m');
            expect(specs.maxSpeed).toBe('2,414 km/h');
            expect(specs.serviceceiling).toBe('15,240 m');
        });

        test('should generate operational data', () => {
            const fighter = scenarioBuilder.platforms[0];
            const operational = scenarioBuilder.generateOperationalData(fighter);

            expect(operational.crewSize).toBe(1);
            expect(operational.fuelLevel).toBe('85%');
            expect(operational.missionReadiness).toBe('ready');
        });
    });

    describe('Enhanced Platform Display', () => {
        test('should create enhanced platform display element', () => {
            const platform = scenarioBuilder.platforms[0];
            const element = scenarioBuilder.createEnhancedPlatformDisplay(platform);

            expect(element.className).toBe('platform-enhanced');
            expect(element.innerHTML).toContain('F-16 Fighting Falcon');
            expect(element.innerHTML).toContain('a-f-A-M-F');
            expect(element.innerHTML).toContain('34th FS');
            expect(element.innerHTML).toContain('388th FW');
            expect(element.innerHTML).toContain('platform-edit-btn');
        });

        test('should include all required CSS classes', () => {
            const platform = scenarioBuilder.platforms[0];
            const element = scenarioBuilder.createEnhancedPlatformDisplay(platform);

            expect(element.innerHTML).toContain('platform-header');
            expect(element.innerHTML).toContain('platform-sections');
            expect(element.innerHTML).toContain('platform-section mil-std');
            expect(element.innerHTML).toContain('symbol-code');
            expect(element.innerHTML).toContain('unit-designation');
            expect(element.innerHTML).toContain('formation');
        });

        test('should handle different platform types correctly', () => {
            const airliner = scenarioBuilder.platforms[1];
            const element = scenarioBuilder.createEnhancedPlatformDisplay(airliner);

            expect(element.innerHTML).toContain('Boeing 737');
            expect(element.innerHTML).toContain('a-n-A-C-A');
            expect(element.innerHTML).toContain('UAL1234');
            expect(element.innerHTML).toContain('United Airlines');
        });
    });

    describe('Platform Data Validation', () => {
        test('should provide fallback values for unknown platforms', () => {
            const unknownPlatform = {
                id: 'unknown_platform',
                name: 'Unknown Platform',
                domain: 'unknown',
                affiliation: 'unknown'
            };

            const symbolCode = scenarioBuilder.generateSymbolCode(unknownPlatform);
            const unitDesignation = scenarioBuilder.generateUnitDesignation(unknownPlatform);
            const formation = scenarioBuilder.generateHigherFormation(unknownPlatform);

            expect(symbolCode).toBe('a-u-G-U-U');
            expect(unitDesignation).toBe('UNIT-001');
            expect(formation).toBe('Formation');
        });

        test('should handle missing or incomplete platform data', () => {
            const incompletePlatform = {
                name: 'Incomplete Platform'
            };

            const enhanced = scenarioBuilder.enhancePlatformData(incompletePlatform);

            expect(enhanced.symbolCode).toBeDefined();
            expect(enhanced.unitDesignation).toBeDefined();
            expect(enhanced.higherFormation).toBeDefined();
            expect(enhanced.specifications).toBeDefined();
            expect(enhanced.operational).toBeDefined();
        });
    });
});
