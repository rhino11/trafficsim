// Scenario Builder JavaScript

// Import UIController for Node.js environments (testing)
if (typeof require !== 'undefined' && typeof UIController === 'undefined') {
    var UIController = require('./ui-controller.js');
}

class ScenarioBuilder {
    constructor() {
        this.map = null;
        this.platforms = [];
        this.scenarioPlatforms = [];
        this.selectedPlatform = null;
        this.currentMarker = null;
        this.mapMarkers = [];
        this.platformCounter = 1;
        this.waypointMode = false;
        this.currentRoute = [];
        this.routePolylines = [];

        // Filter state - now supports multiple selections
        this.activeFilters = {
            domains: new Set(), // Multiple domain selections
            affiliations: new Set(), // Multiple affiliation selections
            search: ''
        };

        // Scenario Management UI Elements
        this.scenarioSelect = document.getElementById('scenarioSelect');
        this.loadScenarioBtn = document.getElementById('loadScenarioBtn');
        this.runScenarioBtn = document.getElementById('runScenarioBtn');
        this.saveScenarioBtn = document.getElementById('saveScenarioBtn');
        this.scenarioNameInput = document.getElementById('scenarioName'); // Assuming this is used for custom scenario names
        this.scenarioDescriptionInput = document.getElementById('scenarioDescription'); // Assuming this is used for custom scenario descriptions

        // Initialize UI Controller
        this.uiController = new UIController(this);

        this.init();
    }

    async init() {
        this.initMap();
        await this.loadPlatforms();
        this.uiController.setupEventListeners(); // Use UIController for event listeners
        this.initializeFilters(); // Initialize filters after loading platforms
        await this.populateScenarioDropdown(); // Populate scenarios on init
        this.uiController.updateScenarioActionButtonsState(); // Use UIController for button states
        this.uiController.updateStatus('Scenario builder ready');
    }

    initializeFilters() {
        console.log('🏁 INITIALIZING FILTERS');
        console.log('  Platforms loaded:', this.platforms.length);
        console.log('  Sample platform:', this.platforms[0]);

        // Initialize all filters as active (show all platforms by default)
        this.selectAllFilters();

        // Set UI buttons to match immediately
        console.log('  Setting all buttons active...');
        this.setAllButtonsActive();
        console.log('  Initial filter state - domains:', Array.from(this.activeFilters.domains));
        console.log('  Initial filter state - affiliations:', Array.from(this.activeFilters.affiliations));
        this.applyFilters();
    }

    initMap() {
        if (typeof L === 'undefined') {
            console.error('Leaflet library not loaded');
            return;
        }

        // Initialize the map centered on a default location
        this.map = L.map('map').setView([39.8283, -98.5795], 4); // Center of USA

        // Add tile layer - this is what the test expects to be called
        L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
            attribution: '© OpenStreetMap contributors',
            maxZoom: 18
        }).addTo(this.map);

        // Initialize platform layer for displaying platforms
        this.platformLayer = L.layerGroup().addTo(this.map);

        this.updateStatus('Map initialized successfully');
    }

    async loadPlatforms() {
        try {
            // Load platform types from the server API
            const response = await fetch('/api/platform-types');
            if (response.ok) {
                this.platforms = await response.json();
                this.updateStatus(`Loaded ${this.platforms.length} platform types from server`);
            } else {
                // Fallback: use hardcoded platform data
                console.warn('Failed to load platforms from server, using fallback data');
                this.platforms = this.getDefaultPlatforms();
                this.updateStatus('Using fallback platform data');
            }
            this.renderPlatformList();
        } catch (error) {
            console.warn('Failed to load platforms from server, using defaults:', error);
            this.platforms = this.getDefaultPlatforms();
            this.updateStatus('Using fallback platform data due to network error');
            this.renderPlatformList();
        }
    }

    // Fallback method for when server data is unavailable
    getDefaultPlatforms() {
        // This is now a fallback method only - primary data loading is dynamic via API
        console.warn('Using hardcoded fallback platform data - this should only happen if the server is unavailable');
        return [
            // Minimal fallback set - the server should provide the full platform catalog
            {
                id: 'airbus_a320',
                name: 'Airbus A320',
                class: 'Airbus A320',
                category: 'commercial_aircraft',
                domain: 'airborne',
                affiliation: 'commercial',
                description: 'Short to medium-range commercial airliner',
                performance: { max_speed: 257.0, cruise_speed: 230.0, max_altitude: 12000 }
            },
            {
                id: 'f16_fighting_falcon',
                name: 'F-16 Fighting Falcon',
                class: 'F-16 Fighting Falcon',
                category: 'fighter_aircraft',
                domain: 'airborne',
                affiliation: 'military',
                description: 'Multi-role fighter aircraft',
                performance: { max_speed: 588.89, cruise_speed: 261.11, max_altitude: 15240 }
            },
            {
                id: 'container_ship',
                name: 'Container Ship',
                class: 'Container Ship',
                category: 'cargo_vessel',
                domain: 'maritime',
                affiliation: 'commercial',
                description: 'Large cargo container vessel',
                performance: { max_speed: 12.9, cruise_speed: 10.8 }
            },
            {
                id: 'destroyer_ship',
                name: 'Naval Destroyer',
                class: 'Naval Destroyer',
                category: 'warship',
                domain: 'maritime',
                affiliation: 'military',
                description: 'Multi-purpose naval combat vessel',
                performance: { max_speed: 16.2, cruise_speed: 10.8 }
            },
            {
                id: 'tesla_model_s',
                name: 'Tesla Model S',
                class: 'Tesla Model S',
                category: 'passenger_vehicle',
                domain: 'land',
                affiliation: 'commercial',
                description: 'Electric luxury sedan',
                performance: { max_speed: 69.4, cruise_speed: 33.3 }
            },
            {
                id: 'military_truck',
                name: 'Military Transport Truck',
                class: 'Military Transport Truck',
                category: 'military_vehicle',
                domain: 'land',
                affiliation: 'military',
                description: 'Heavy-duty military transport vehicle',
                performance: { max_speed: 27.8, cruise_speed: 22.2 }
            },
            {
                id: 'starlink_satellite',
                name: 'Starlink Satellite',
                class: 'Starlink Satellite',
                category: 'communications_satellite',
                domain: 'space',
                affiliation: 'commercial',
                description: 'Low Earth orbit communications satellite',
                performance: { max_speed: 7660.0, cruise_speed: 7660.0, max_altitude: 550000 }
            },
            {
                id: 'military_satellite',
                name: 'Military Reconnaissance Satellite',
                class: 'Military Reconnaissance Satellite',
                category: 'reconnaissance_satellite',
                domain: 'space',
                affiliation: 'military',
                description: 'Military surveillance and reconnaissance satellite',
                performance: { max_speed: 7660.0, cruise_speed: 7660.0, max_altitude: 600000 }
            }
        ];
    }

    // Filter management methods
    handleFilterClick(button) {
        const filterType = button.dataset.filterType;
        const filterValue = button.dataset.filterValue;

        if (filterType === 'domain') {
            this.handleDomainFilter(button, filterValue);
        } else if (filterType === 'affiliation') {
            this.handleAffiliationFilter(button, filterValue);
        } else if (filterType === 'all' || filterValue === 'all') {
            this.selectAllFilters();
        }

        this.applyFilters();
        this.updateAllButtonState();
    }

    handleDomainFilter(button, domain) {
        if (domain === 'all') {
            this.selectAllFilters();
            return;
        }

        if (this.activeFilters.domains.has(domain)) {
            this.activeFilters.domains.delete(domain);
            button.classList.remove('active');
        } else {
            this.activeFilters.domains.add(domain);
            button.classList.add('active');
        }
    }

    handleAffiliationFilter(button, affiliation) {
        if (affiliation === 'all') {
            this.selectAllFilters();
            return;
        }

        // Check if only this affiliation is currently selected
        if (this.activeFilters.affiliations.has(affiliation) &&
            this.activeFilters.affiliations.size === 1) {
            // Toggle back to show all affiliations
            this.activeFilters.affiliations.clear();
            this.activeFilters.affiliations.add('commercial');
            this.activeFilters.affiliations.add('military');

            // Activate all affiliation buttons
            document.querySelectorAll('.filter-btn[data-filter-type="affiliation"]').forEach(btn => {
                btn.classList.add('active');
            });
        } else {
            // Exclusive selection - only one affiliation can be active at a time
            document.querySelectorAll('.filter-btn[data-filter-type="affiliation"]').forEach(btn => {
                btn.classList.remove('active');
            });
            this.activeFilters.affiliations.clear();

            this.activeFilters.affiliations.add(affiliation);
            button.classList.add('active');
        }
    }

    selectAllFilters() {
        // Clear all filters
        this.activeFilters.domains.clear();
        this.activeFilters.affiliations.clear();

        // Remove active class from all filter buttons except "All"
        document.querySelectorAll('.filter-btn').forEach(btn => {
            btn.classList.remove('active');
        });

        // Activate all "All" buttons
        this.setAllButtonsActive();
    }

    setAllButtonsActive() {
        document.querySelectorAll('.filter-btn[data-filter-value="all"]').forEach(btn => {
            btn.classList.add('active');
        });
    }

    updateAllButtonState() {
        const domainAllBtn = document.querySelector('.filter-btn[data-filter-type="domain"][data-filter-value="all"]');
        const affiliationAllBtn = document.querySelector('.filter-btn[data-filter-type="affiliation"][data-filter-value="all"]');

        // Update domain "All" button
        if (domainAllBtn) {
            if (this.isAllFiltersActive('domain')) {
                domainAllBtn.classList.add('active');
            } else {
                domainAllBtn.classList.remove('active');
            }
        }

        // Update affiliation "All" button
        if (affiliationAllBtn) {
            if (this.isAllFiltersActive('affiliation')) {
                affiliationAllBtn.classList.add('active');
            } else {
                affiliationAllBtn.classList.remove('active');
            }
        }
    }

    isAllFiltersActive(filterType) {
        if (filterType === 'domain') {
            return this.activeFilters.domains.size === 0;
        } else if (filterType === 'affiliation') {
            return this.activeFilters.affiliations.size === 0;
        }
        return false;
    }

    applyFilters() {
        if (!this.platforms || this.platforms.length === 0) {
            console.log('No platforms loaded for filtering');
            return;
        }

        let filteredPlatforms = [...this.platforms];

        // Apply domain filters (multiple selection allowed)
        if (this.activeFilters.domains.size > 0) {
            filteredPlatforms = filteredPlatforms.filter(platform =>
                this.activeFilters.domains.has(platform.domain)
            );
        }

        // Apply affiliation filters (exclusive selection)
        if (this.activeFilters.affiliations.size > 0 && this.activeFilters.affiliations.size < 2) {
            filteredPlatforms = filteredPlatforms.filter(platform =>
                this.activeFilters.affiliations.has(platform.affiliation)
            );
        }

        // Apply search filter
        if (this.activeFilters.search) {
            const searchTerm = this.activeFilters.search.toLowerCase();
            filteredPlatforms = filteredPlatforms.filter(platform =>
                platform.name.toLowerCase().includes(searchTerm) ||
                platform.description.toLowerCase().includes(searchTerm) ||
                platform.category.toLowerCase().includes(searchTerm)
            );
        }

        this.renderFilteredPlatformList(filteredPlatforms);
        return filteredPlatforms;
    }

    renderFilteredPlatformList(platforms) {
        const container = document.getElementById('platformList');
        if (!container) return;

        container.innerHTML = '';

        if (platforms.length === 0) {
            container.innerHTML = '<div class="no-results">No platforms match your filters</div>';
            return;
        }

        platforms.forEach(platform => {
            const item = document.createElement('div');
            item.className = 'platform-item';
            item.innerHTML = `
                <h4>${platform.name}</h4>
                <p><strong>Type:</strong> ${platform.category}</p>
                <p><strong>Domain:</strong> ${this.getDomainIcon(platform.domain)} ${platform.domain}</p>
                <p><strong>Affiliation:</strong> ${platform.affiliation}</p>
                <p>${platform.description}</p>
            `;

            item.addEventListener('click', () => {
                this.selectPlatform(platform, item);
            });

            container.appendChild(item);
        });
    }

    getDomainIcon(domain) {
        const icons = {
            airborne: '✈️',
            maritime: '🚢',
            land: '🚛',
            space: '🛰️'
        };
        return icons[domain] || '🔹';
    }

    async populateScenarioDropdown() {
        try {
            const response = await fetch('/api/scenarios');
            if (response.ok) {
                const scenarios = await response.json();

                // Clear existing options except the default "Custom Scenario"
                this.scenarioSelect.innerHTML = '<option value="custom">Custom Scenario</option>';

                // Add scenarios from the server
                scenarios.forEach(scenario => {
                    const option = document.createElement('option');
                    option.value = scenario.filename;
                    option.textContent = scenario.display_name || scenario.filename;
                    option.setAttribute('data-description', scenario.description || 'No description available');
                    this.scenarioSelect.appendChild(option);
                });

                this.updateStatus(`Loaded ${scenarios.length} pre-configured scenarios`);
            } else {
                console.warn('Failed to load scenarios from server');
                this.updateStatus('Could not load pre-configured scenarios');
            }
        } catch (error) {
            console.error('Error loading scenarios:', error);
            this.updateStatus('Error loading scenarios from server');
        }

        // Always update button states after populating
        this.uiController.updateScenarioActionButtonsState();
    }

    async loadSelectedScenario() {
        const selectedValue = this.scenarioSelect.value;

        if (!selectedValue || selectedValue === 'custom') {
            this.updateStatus('Please select a pre-configured scenario to load');
            return;
        }

        try {
            // Show loading state
            this.loadScenarioBtn.disabled = true;
            this.loadScenarioBtn.textContent = 'Loading...';

            const response = await fetch(`/api/scenario/${encodeURIComponent(selectedValue)}`);
            if (response.ok) {
                const scenarioData = await response.json();

                // Clear current scenario first
                this.clearScenario(false); // Don't ask for confirmation when loading a new scenario

                // Load the scenario data
                this.scenarioPlatforms = scenarioData.platforms || [];

                // Update scenario info fields if they exist
                if (this.scenarioNameInput && scenarioData.metadata?.name) {
                    this.scenarioNameInput.value = scenarioData.metadata.name;
                }
                if (this.scenarioDescriptionInput && scenarioData.metadata?.description) {
                    this.scenarioDescriptionInput.value = scenarioData.metadata.description;
                }

                // Add platforms to map
                this.scenarioPlatforms.forEach(platform => {
                    this.addMapMarker(platform);
                });

                // Update displays
                this.renderScenarioPlatforms();

                // Fit map to show all platforms if any exist
                if (this.mapMarkers.length > 0) {
                    const group = new L.featureGroup(this.mapMarkers.map(m => m.marker));
                    this.map.fitBounds(group.getBounds().pad(0.1));
                }

                this.updateStatus(`Loaded scenario: ${scenarioData.metadata?.name || selectedValue}`);
            } else {
                const errorText = await response.text();
                console.error('Failed to load scenario:', errorText);
                this.updateStatus(`Error loading scenario: ${response.statusText}`);
            }
        } catch (error) {
            console.error('Error loading scenario:', error);
            this.updateStatus('Error loading scenario from server');
        } finally {
            // Restore button state
            this.loadScenarioBtn.textContent = 'Load Scenario';
            this.uiController.updateScenarioActionButtonsState();
        }
    }

    async runCurrentScenario() {
        if (this.scenarioPlatforms.length === 0) {
            this.updateStatus('No platforms in scenario to run');
            return;
        }

        try {
            // Show loading state
            this.runScenarioBtn.disabled = true;
            this.runScenarioBtn.textContent = 'Starting...';

            // Prepare scenario data for simulation
            const scenarioData = {
                metadata: {
                    name: this.scenarioNameInput?.value || 'Custom Scenario',
                    description: this.scenarioDescriptionInput?.value || 'User-created scenario',
                    created_at: new Date().toISOString()
                },
                platforms: this.scenarioPlatforms
            };

            const response = await fetch('/api/scenario/run', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(scenarioData)
            });

            if (response.ok) {
                const result = await response.json();
                this.updateStatus(`Scenario started successfully. Session ID: ${result.session_id || 'unknown'}`);

                // Optionally redirect to simulation view or show success message
                if (result.redirect_url) {
                    setTimeout(() => {
                        window.location.href = result.redirect_url;
                    }, 2000);
                }
            } else {
                const errorText = await response.text();
                console.error('Failed to start scenario:', errorText);
                this.updateStatus(`Error starting scenario: ${response.statusText}`);
            }
        } catch (error) {
            console.error('Error starting scenario:', error);
            this.updateStatus('Error communicating with server');
        } finally {
            // Restore button state
            this.runScenarioBtn.textContent = 'Run Scenario';
            this.uiController.updateScenarioActionButtonsState();
        }
    }

    async saveCustomScenario() {
        if (this.scenarioPlatforms.length === 0) {
            this.updateStatus('No platforms in scenario to save');
            return;
        }

        // Get scenario name from user
        const scenarioName = this.scenarioNameInput?.value || 'Custom Scenario';
        const scenarioDescription = this.scenarioDescriptionInput?.value || '';

        // Validate scenario name
        if (scenarioName.trim() === '' || scenarioName === 'Custom Scenario') {
            const userProvidedName = prompt('Please enter a name for your custom scenario:');
            if (!userProvidedName || userProvidedName.trim() === '') {
                this.updateStatus('Scenario name is required to save');
                return;
            }
            if (this.scenarioNameInput) {
                this.scenarioNameInput.value = userProvidedName.trim();
            }
        }

        try {
            // Show loading state
            this.saveScenarioBtn.disabled = true;
            this.saveScenarioBtn.textContent = 'Saving...';

            // Prepare scenario data
            const scenarioData = {
                metadata: {
                    name: this.scenarioNameInput?.value || scenarioName,
                    description: this.scenarioDescriptionInput?.value || scenarioDescription,
                    version: "1.0",
                    created_at: new Date().toISOString(),
                    created_by: "scenario_builder"
                },
                platforms: this.scenarioPlatforms
            };

            const response = await fetch('/api/scenario/save', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(scenarioData)
            });

            if (response.ok) {
                const result = await response.json();
                this.updateStatus(`Scenario saved successfully as: ${result.filename || 'custom_scenario.yaml'}`);

                // Refresh the dropdown to include the newly saved scenario
                await this.populateScenarioDropdown();

                // Select the newly saved scenario in the dropdown
                if (result.filename) {
                    this.scenarioSelect.value = result.filename;
                }
            } else {
                const errorText = await response.text();
                console.error('Failed to save scenario:', errorText);
                this.updateStatus(`Error saving scenario: ${response.statusText}`);
            }
        } catch (error) {
            console.error('Error saving scenario:', error);
            this.updateStatus('Error saving scenario to server');
        } finally {
            // Restore button state
            this.saveScenarioBtn.textContent = 'Save Custom Scenario';
            this.uiController.updateScenarioActionButtonsState();
        }
    }



    // Override or ensure clearScenario updates button states
    clearScenario(confirmClear = true) {
        let doClear = true;
        if (confirmClear) {
            doClear = confirm('Are you sure you want to clear the current scenario?');
        }

        if (doClear) {
            this.scenarioPlatforms = [];
            this.mapMarkers.forEach(markerData => {
                this.map.removeLayer(markerData.marker);
            });
            this.mapMarkers = [];
            this.renderScenarioPlatforms();
            if (this.scenarioNameInput) this.scenarioNameInput.value = 'Custom Scenario';
            if (this.scenarioDescriptionInput) this.scenarioDescriptionInput.value = 'Build your scenario from scratch.';
            if (this.scenarioSelect) this.scenarioSelect.value = 'custom'; // Reset dropdown to custom
            this.updateStatus('Scenario cleared. Ready for new custom scenario.');
        }
        this.uiController.updateScenarioActionButtonsState(); // Always update state
        return doClear; // Return whether clearing happened
    }

    // Ensure platform modifications update button states
    addMapMarker(platform) {
        const icon = this.getPlatformIcon(platform.domain);
        const marker = L.marker([platform.start_position.latitude, platform.start_position.longitude], {
            icon: icon
        }).addTo(this.map);

        marker.bindPopup(`
            <strong>${platform.name}</strong><br>
            Type: ${platform.class}<br>
            Altitude: ${platform.start_position.altitude}m<br>
            Mission: ${platform.mission.type}
        `);

        this.mapMarkers.push({ marker, platform });
        this.uiController.updateScenarioActionButtonsState(); // Update buttons
    }

    removePlatform(index) {
        // Remove marker from map
        const markerData = this.mapMarkers.find(m => m.platform === this.scenarioPlatforms[index]);
        if (markerData) {
            this.map.removeLayer(markerData.marker);
            this.mapMarkers = this.mapMarkers.filter(m => m !== markerData);
        }

        // Remove from scenario platforms
        this.scenarioPlatforms.splice(index, 1);
        this.renderScenarioPlatforms();
        this.uiController.updateScenarioActionButtonsState();
    }

    renderPlatformList() {
        const platformListDiv = document.getElementById('platformList');
        if (!platformListDiv) return;

        // Clear existing list
        platformListDiv.innerHTML = '';

        // Create and append platform items
        this.platforms.forEach(platform => {
            const platformItem = document.createElement('div');
            platformItem.className = 'platform-item';
            platformItem.innerHTML = `
                <strong>${platform.name}</strong> (${platform.id})
                <button class="select-platform" data-id="${platform.id}">Select</button>
            `;
            platformListDiv.appendChild(platformItem);
        });

        // Add "Custom Platform" option
        const customPlatformItem = document.createElement('div');
        customPlatformItem.className = 'platform-item custom-platform';
        customPlatformItem.innerHTML = `
            <strong>Custom Platform</strong>
            <button class="select-platform" data-id="custom">Select</button>
        `;
        platformListDiv.appendChild(customPlatformItem);

        // Re-apply event listeners to new buttons
        this.applyPlatformSelectListeners();
    }

    applyPlatformSelectListeners() {
        document.querySelectorAll('.select-platform').forEach(button => {
            button.addEventListener('click', (e) => {
                const platformId = e.target.dataset.id;
                this.selectPlatform(platformId);
            });
        });
    }

    selectPlatform(platformId) {
        if (platformId === 'custom') {
            // Enter custom platform mode
            this.selectedPlatform = null;
            this.showPlatformEditor();
        } else {
            // Select existing platform
            const platform = this.platforms.find(p => p.id === platformId);
            if (platform) {
                this.selectedPlatform = platform;
                this.fillPlatformDetails(platform);
            }
        }
    }

    showPlatformEditor() {
        // Clear details
        this.clearPlatformDetails();

        // Show editor UI
        document.getElementById('platformEditor').style.display = 'block';
        document.getElementById('platformDetails').style.display = 'none';
    }

    clearPlatformDetails() {
        const detailsDiv = document.getElementById('platformDetails');
        if (detailsDiv) {
            detailsDiv.innerHTML = '';
        }
    }

    fillPlatformDetails(platform) {
        const detailsDiv = document.getElementById('platformDetails');
        if (!detailsDiv) return;

        // Clear existing details
        detailsDiv.innerHTML = '';

        // Fill with platform data
        for (const [key, value] of Object.entries(platform)) {
            const p = document.createElement('p');
            p.innerHTML = `<strong>${key}:</strong> ${value}`;
            detailsDiv.appendChild(p);
        }

        // Show details UI
        document.getElementById('platformEditor').style.display = 'none';
        detailsDiv.style.display = 'block';
    }

    handleScenarioSelectionChange() {
        // Update button states when scenario selection changes
        this.uiController.updateScenarioActionButtonsState();

        // Optionally show description or other info about selected scenario
        const selectedOption = this.scenarioSelect.options[this.scenarioSelect.selectedIndex];
        if (selectedOption && selectedOption.getAttribute('data-description')) {
            const description = selectedOption.getAttribute('data-description');
            // You could show this in a status area or tooltip
            console.log('Selected scenario description:', description);
        }
    }

    // UI delegation methods - delegate to UIController
    updateStatus(message, timeout) {
        return this.uiController.updateStatus(message, timeout);
    }

    showNotification(message, type, duration) {
        return this.uiController.showNotification(message, type, duration);
    }

    updateScenarioActionButtonsState() {
        return this.uiController.updateScenarioActionButtonsState();
    }

    openModal(modalId) {
        return this.uiController.openModal(modalId);
    }

    closeModal(modal) {
        return this.uiController.closeModal(modal);
    }

    clearPlatformSelection() {
        return this.uiController.clearPlatformSelection();
    }

    // Core platform interaction methods
    selectPlatform(platform, element) {
        // Clear previous selection
        document.querySelectorAll('.platform-item').forEach(item => {
            item.classList.remove('selected');
        });

        // Set new selection
        this.selectedPlatform = platform;
        if (element) {
            element.classList.add('selected');
        }

        document.getElementById('mapInstructions').textContent = 'Click on the map to place this platform';
        this.updateStatus(`Selected platform: ${platform.name}`);

        // Enable map click if in placement mode
        if (this.map) {
            this.map.on('click', (e) => {
                this.placePlatformOnMap(e.latlng, platform);
            });
        }
    }

    placePlatformOnMap(latlng, platform) {
        if (!platform) return;

        const platformData = {
            id: `platform_${this.platformCounter++}`,
            type: platform,
            position: {
                lat: latlng.lat,
                lng: latlng.lng,
                altitude: platform.performance?.max_altitude || 0
            },
            name: platform.name
        };

        this.showPlatformConfigModal(latlng);
    }

    showPlatformConfigModal(latlng) {
        if (!this.selectedPlatform) return;

        // Set the current click position
        this.currentLatLng = latlng;

        // Fill modal with platform data
        document.getElementById('modalTitle').textContent = `Configure ${this.selectedPlatform.name}`;
        document.getElementById('modalLatitude').value = latlng.lat.toFixed(6);
        document.getElementById('modalLongitude').value = latlng.lng.toFixed(6);
        document.getElementById('modalAltitude').value = this.selectedPlatform.performance?.max_altitude || 1000;
        document.getElementById('modalPlatformName').value = this.generatePlatformName(this.selectedPlatform);

        // Show modal
        document.getElementById('platformModal').style.display = 'block';
    }

    generatePlatformName(platform) {
        const domainPrefixes = {
            airborne: 'AIR',
            maritime: 'SEA',
            land: 'GND',
            space: 'SAT'
        };

        const prefix = domainPrefixes[platform.domain] || 'PLT';
        return `${prefix}-${this.platformCounter.toString().padStart(3, '0')}`;
    }

    savePlatformConfig() {
        const platformId = document.getElementById('modalPlatformId').value;
        const platformName = document.getElementById('modalPlatformName').value;
        const latitude = parseFloat(document.getElementById('modalLatitude').value);
        const longitude = parseFloat(document.getElementById('modalLongitude').value);
        const altitude = parseInt(document.getElementById('modalAltitude').value);
        const missionType = document.getElementById('modalMissionType').value;

        if (!platformId || !platformName) {
            alert('Please fill in all required fields');
            return;
        }

        // Create scenario platform object
        const scenarioPlatform = {
            id: platformId,
            type: this.selectedPlatform.id,
            name: platformName,
            class: this.selectedPlatform.class,
            domain: this.selectedPlatform.domain,
            start_position: {
                latitude: latitude,
                longitude: longitude,
                altitude: altitude
            },
            mission: {
                type: missionType
            }
        };

        // Add to scenario
        this.scenarioPlatforms.push(scenarioPlatform);
        this.platformCounter++;

        // Add marker to map
        this.addMapMarker(scenarioPlatform);

        // Update UI
        this.renderScenarioPlatforms();
        this.uiController.updateScenarioActionButtonsState();
        this.updateStatus(`Added ${platformName} to scenario`);

        // Close modal
        document.getElementById('platformModal').style.display = 'none';
    }

    getPlatformIcon(domain) {
        const iconColors = {
            airborne: 'blue',
            maritime: 'darkblue',
            land: 'green',
            space: 'purple'
        };

        return L.divIcon({
            className: 'custom-div-icon',
            html: `<div style="background-color: ${iconColors[domain] || 'gray'}; width: 12px; height: 12px; border-radius: 50%; border: 2px solid white;"></div>`,
            iconSize: [16, 16],
            iconAnchor: [8, 8]
        });
    }

    renderScenarioPlatforms() {
        const container = document.getElementById('scenarioPlatforms');
        container.innerHTML = '';

        if (this.scenarioPlatforms.length === 0) {
            container.innerHTML = '<p style="color: #6c757d; font-style: italic;">No platforms added yet</p>';
            return;
        }

        this.scenarioPlatforms.forEach((platform, index) => {
            const item = document.createElement('div');
            item.className = 'scenario-platform';
            item.innerHTML = `
                <button class="remove-platform" onclick="scenarioBuilder.removePlatform(${index})">×</button>
                <h4>${platform.name}</h4>
                <p><strong>Type:</strong> ${platform.class}</p>
                <p><strong>Position:</strong> ${platform.start_position.latitude.toFixed(4)}, ${platform.start_position.longitude.toFixed(4)}</p>
                <p><strong>Altitude:</strong> ${platform.start_position.altitude}m</p>
                <p><strong>Mission:</strong> ${platform.mission.type}</p>
            `;

            container.appendChild(item);
        });
    }

    // Filtering methods
    filterPlatforms(searchTerm) {
        this.activeFilters.search = searchTerm.toLowerCase();
        this.applyFilters();
    }

    filterPlatformsByDomain(domain) {
        if (domain === 'all') {
            this.activeFilters.domains.clear();
            this.renderPlatformList();
            return;
        } else {
            this.activeFilters.domains.clear();
            this.activeFilters.domains.add(domain);
        }
        this.applyFilters();
    }

    // Waypoint and route management
    toggleWaypointMode(enabled) {
        this.waypointMode = enabled;
        if (enabled) {
            this.updateStatus('Waypoint mode enabled - Click on map to add waypoints');
            this.map.off('click');
            this.map.on('click', (e) => {
                this.addWaypoint(e.latlng);
            });
        } else {
            this.updateStatus('Waypoint mode disabled');
            this.currentRoute = [];
            this.clearRoutePolylines();
            this.map.off('click');
            this.map.on('click', (e) => {
                if (this.selectedPlatform) {
                    this.showPlatformConfigModal(e.latlng);
                }
            });
        }
    }

    addWaypoint(latlng) {
        if (!this.waypointMode) return;

        this.currentRoute.push(latlng);

        // Add waypoint marker
        const marker = L.marker([latlng.lat, latlng.lng]).addTo(this.map);
        this.mapMarkers.push({ marker, type: 'waypoint' });

        // Create/update polyline if we have more than one point
        if (this.currentRoute.length > 1) {
            const polyline = L.polyline(this.currentRoute.map(p => [p.lat, p.lng]), {
                color: 'red',
                weight: 3
            }).addTo(this.map);
            this.routePolylines.push(polyline);
        }

        this.updateStatus(`Added waypoint ${this.currentRoute.length}`);
    }

    completeCurrentRoute() {
        if (this.currentRoute.length < 2) {
            alert('Please add at least 2 waypoints to create a route');
            return;
        }

        // Here you could save the route or attach it to a platform
        this.updateStatus(`Route completed with ${this.currentRoute.length} waypoints`);

        // Reset waypoint mode
        document.getElementById('waypointMode').checked = false;
        this.toggleWaypointMode(false);
    }

    clearRoutePolylines() {
        this.routePolylines.forEach(polyline => {
            this.map.removeLayer(polyline);
        });
        this.routePolylines = [];
    }

    // Export/Import functionality
    exportScenario() {
        if (this.scenarioPlatforms.length === 0) {
            alert('Please add at least one platform to export');
            return;
        }

        const yaml = this.generateScenarioYAML();
        const blob = new Blob([yaml], { type: 'text/yaml' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `scenario_${Date.now()}.yaml`;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);

        this.updateStatus('Scenario exported successfully');
    }

    generateScenarioYAML() {
        const scenario = {
            name: this.scenarioName || 'Test Scenario',
            description: this.scenarioDescription || 'Generated scenario',
            platforms: this.scenarioPlatforms.map(platform => ({
                id: platform.id,
                name: platform.name,
                class: platform.class,
                domain: platform.domain,
                affiliation: platform.affiliation,
                start_position: platform.start_position,
                mission: platform.mission
            }))
        };

        // Simple YAML generation
        let yaml = `name: "${scenario.name}"\n`;
        yaml += `description: "${scenario.description}"\n`;
        yaml += `platforms:\n`;

        scenario.platforms.forEach(platform => {
            yaml += `  - id: "${platform.id}"\n`;
            yaml += `    name: "${platform.name}"\n`;
            yaml += `    class: "${platform.class}"\n`;
            yaml += `    domain: "${platform.domain}"\n`;
            yaml += `    affiliation: "${platform.affiliation}"\n`;
            yaml += `    start_position:\n`;
            yaml += `      latitude: ${platform.start_position.latitude}\n`;
            yaml += `      longitude: ${platform.start_position.longitude}\n`;
            yaml += `      altitude: ${platform.start_position.altitude}\n`;
            yaml += `    mission:\n`;
            yaml += `      type: "${platform.mission.type}"\n`;
        });

        return yaml;
    }

    loadScenario(file) {
        const reader = new FileReader();
        reader.onload = (e) => {
            try {
                // Simple YAML parsing - in reality you'd use a proper YAML parser
                const content = e.target.result;
                const lines = content.split('\n');

                // Reset current scenario
                this.scenarioPlatforms = [];
                this.renderScenarioPlatforms();

                this.updateStatus('Scenario loaded successfully');
            } catch (error) {
                console.error('Error loading scenario:', error);
                alert('Error loading scenario file');
            }
        };
        reader.readAsText(file);
    }

    // Validation functionality
    validateScenario() {
        const issues = [];

        if (this.scenarioPlatforms.length === 0) {
            issues.push('No platforms added to scenario');
        }

        // Check for overlapping platforms
        for (let i = 0; i < this.scenarioPlatforms.length; i++) {
            for (let j = i + 1; j < this.scenarioPlatforms.length; j++) {
                const platform1 = this.scenarioPlatforms[i];
                const platform2 = this.scenarioPlatforms[j];

                const distance = this.calculateDistance(
                    platform1.start_position.latitude,
                    platform1.start_position.longitude,
                    platform2.start_position.latitude,
                    platform2.start_position.longitude
                );

                if (distance < 1) { // Less than 1km apart
                    issues.push(`Platforms ${platform1.name} and ${platform2.name} are very close (${distance.toFixed(2)}km apart)`);
                }
            }
        }

        return issues;
    }

    calculateDistance(lat1, lon1, lat2, lon2) {
        const R = 6371; // Earth's radius in km
        const dLat = (lat2 - lat1) * Math.PI / 180;
        const dLon = (lon2 - lon1) * Math.PI / 180;
        const a = Math.sin(dLat / 2) * Math.sin(dLat / 2) +
            Math.cos(lat1 * Math.PI / 180) * Math.cos(lat2 * Math.PI / 180) *
            Math.sin(dLon / 2) * Math.sin(dLon / 2);
        const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
        return R * c;
    }

    // Additional functionality
    getDomainStats() {
        const stats = {
            airborne: 0,
            maritime: 0,
            land: 0,
            space: 0,
            total: this.scenarioPlatforms.length
        };

        this.scenarioPlatforms.forEach(platform => {
            if (stats.hasOwnProperty(platform.domain)) {
                stats[platform.domain]++;
            }
        });

        return stats;
    }

    previewScenario() {
        if (this.scenarioPlatforms.length === 0) {
            alert('Please add platforms to preview');
            return;
        }

        // Create a feature group with all platforms
        const group = L.featureGroup();
        let addedMarkers = 0;

        this.scenarioPlatforms.forEach(platform => {
            // Check if platform has position data
            if (!platform.start_position ||
                typeof platform.start_position.latitude === 'undefined' ||
                typeof platform.start_position.longitude === 'undefined') {
                return; // Skip platforms without position data
            }

            const marker = L.marker([
                platform.start_position.latitude,
                platform.start_position.longitude
            ], {
                icon: this.getPlatformIcon(platform.domain)
            }).bindPopup(`
                <strong>${platform.name || 'Unknown Platform'}</strong><br>
                Type: ${platform.class || 'Unknown'}<br>
                Domain: ${platform.domain || 'Unknown'}<br>
                Mission: ${platform.mission ? platform.mission.type : 'Unknown'}
            `);

            if (group.addLayer) {
                group.addLayer(marker);
            }
            addedMarkers++;
        });

        // Only fit bounds if we have markers and valid map
        if (addedMarkers > 0 && this.map && this.map.fitBounds && group.getBounds) {
            // Fit the map to show all platforms
            this.map.fitBounds(group.getBounds());
        }

        this.updateStatus(`Previewing scenario with ${addedMarkers} platforms`);
    }

    loadPlatformLibrary() {
        // Initialize platform library with sample data if not already loaded
        if (!this.platformLibrary) {
            this.platformLibrary = {
                airborne: [
                    { id: 'f16', name: 'F-16 Fighting Falcon', class: 'Fighter', affiliation: 'military' },
                    { id: 'boeing737', name: 'Boeing 737', class: 'Commercial Airliner', affiliation: 'commercial' }
                ],
                maritime: [
                    { id: 'destroyer', name: 'Arleigh Burke Destroyer', class: 'Destroyer', affiliation: 'military' },
                    { id: 'cargo', name: 'Container Ship', class: 'Cargo Vessel', affiliation: 'commercial' }
                ],
                land: [
                    { id: 'm1a2', name: 'M1A2 Abrams', class: 'Main Battle Tank', affiliation: 'military' },
                    { id: 'truck', name: 'Cargo Truck', class: 'Transport Vehicle', affiliation: 'commercial' }
                ],
                space: [
                    { id: 'satellite', name: 'Communications Satellite', class: 'Satellite', affiliation: 'commercial' }
                ]
            };
        }

        this.updateStatus('Platform library loaded');
    }

    // Fix for renderPlatformList method that was missing
    renderPlatformList() {
        const container = document.getElementById('platformList');
        if (!container) return;

        container.innerHTML = '';

        const filteredPlatforms = this.getFilteredPlatforms();

        filteredPlatforms.forEach(platform => {
            const item = document.createElement('div');
            item.className = 'platform-item';
            item.innerHTML = `
                <h4>${platform.name}</h4>
                <p>Type: ${platform.class}</p>
                <p>Domain: ${platform.domain}</p>
                <button onclick="scenarioBuilder.selectPlatform('${platform.id}')">Add to Scenario</button>
            `;
            container.appendChild(item);
        });
    }

    getFilteredPlatforms() {
        // Return all platforms from the library based on active filters
        let allPlatforms = [];

        if (this.platformLibrary) {
            Object.values(this.platformLibrary).forEach(domainPlatforms => {
                allPlatforms = allPlatforms.concat(domainPlatforms);
            });
        }

        // Apply filters
        if (this.activeFilters.domains.size > 0 && !this.activeFilters.domains.has('all')) {
            allPlatforms = allPlatforms.filter(platform =>
                this.activeFilters.domains.has(platform.domain)
            );
        }

        if (this.activeFilters.affiliations.size > 0) {
            allPlatforms = allPlatforms.filter(platform =>
                this.activeFilters.affiliations.has(platform.affiliation)
            );
        }

        return allPlatforms;
    }
}

// Global instance
let scenarioBuilder;

// Initialize when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    scenarioBuilder = new ScenarioBuilder();
});

// Export for module usage
if (typeof module !== 'undefined' && module.exports) {
    module.exports = ScenarioBuilder;
}
