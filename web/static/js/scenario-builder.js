// Scenario Builder JavaScript
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

        this.init();
    }

    async init() {
        this.initMap();
        await this.loadPlatforms();
        this.setupEventListeners();
        this.updateStatus('Scenario builder ready');
    }

    initMap() {
        // Initialize Leaflet map centered on the US
        this.map = L.map('map').setView([39.8283, -98.5795], 4);

        // Add OpenStreetMap tiles
        L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
            attribution: '© OpenStreetMap contributors'
        }).addTo(this.map);

        // Add click event for placing platforms
        this.map.on('click', (e) => {
            if (this.selectedPlatform) {
                this.showPlatformConfigModal(e.latlng);
            }
        });
    }

    async loadPlatforms() {
        try {
            // Load platform types from the server
            const response = await fetch('/api/platform-types');
            if (response.ok) {
                this.platforms = await response.json();
            } else {
                // Fallback: use hardcoded platform data
                this.platforms = this.getDefaultPlatforms();
            }
            this.renderPlatformList();
        } catch (error) {
            console.warn('Failed to load platforms from server, using defaults:', error);
            this.platforms = this.getDefaultPlatforms();
            this.renderPlatformList();
        }
    }

    getDefaultPlatforms() {
        return [
            {
                id: 'airbus_a320',
                name: 'Airbus A320',
                class: 'Airbus A320',
                category: 'commercial_aircraft',
                domain: 'airborne',
                description: 'Short to medium-range commercial airliner',
                performance: { max_speed: 257.0, cruise_speed: 230.0, max_altitude: 12000 }
            },
            {
                id: 'boeing_777_300er',
                name: 'Boeing 777-300ER',
                class: 'Boeing 777-300ER',
                category: 'wide_body_airliner',
                domain: 'airborne',
                description: 'Long-range wide-body commercial airliner',
                performance: { max_speed: 290.0, cruise_speed: 257.0, max_altitude: 13100 }
            },
            {
                id: 'f16_fighting_falcon',
                name: 'F-16 Fighting Falcon',
                class: 'F-16 Fighting Falcon',
                category: 'fighter_aircraft',
                domain: 'airborne',
                description: 'Multi-role fighter aircraft',
                performance: { max_speed: 588.89, cruise_speed: 261.11, max_altitude: 15240 }
            },
            {
                id: 'container_ship',
                name: 'Container Ship',
                class: 'Container Ship',
                category: 'cargo_vessel',
                domain: 'maritime',
                description: 'Large cargo container vessel',
                performance: { max_speed: 12.9, cruise_speed: 10.8 }
            },
            {
                id: 'arleigh_burke_destroyer',
                name: 'Arleigh Burke Destroyer',
                class: 'Arleigh Burke-class Destroyer',
                category: 'guided_missile_destroyer',
                domain: 'maritime',
                description: 'US Navy guided missile destroyer',
                performance: { max_speed: 15.4, cruise_speed: 10.3 }
            }
        ];
    }

    setupEventListeners() {
        // Search functionality
        document.getElementById('platformSearch').addEventListener('input', (e) => {
            this.filterPlatforms(e.target.value);
        });

        // Domain filter buttons
        document.querySelectorAll('.domain-filter button').forEach(btn => {
            btn.addEventListener('click', (e) => {
                document.querySelectorAll('.domain-filter button').forEach(b => b.classList.remove('active'));
                e.target.classList.add('active');
                this.filterPlatformsByDomain(e.target.dataset.domain);
            });
        });

        // Waypoint mode toggle
        const waypointToggle = document.getElementById('waypointMode');
        if (waypointToggle) {
            waypointToggle.addEventListener('change', (e) => {
                this.toggleWaypointMode(e.target.checked);
            });
        }

        // Route completion button
        const completeRouteBtn = document.getElementById('completeRoute');
        if (completeRouteBtn) {
            completeRouteBtn.addEventListener('click', () => {
                this.completeCurrentRoute();
            });
        }

        // Modal close events
        document.querySelectorAll('.close').forEach(closeBtn => {
            closeBtn.addEventListener('click', (e) => {
                e.target.closest('.modal').style.display = 'none';
            });
        });

        // Click outside modal to close
        window.addEventListener('click', (e) => {
            if (e.target.classList.contains('modal')) {
                e.target.style.display = 'none';
            }
        });
    }

    renderPlatformList() {
        const container = document.getElementById('platformList');
        container.innerHTML = '';

        this.platforms.forEach(platform => {
            const item = document.createElement('div');
            item.className = 'platform-item';
            item.innerHTML = `
                <h4>${platform.name}</h4>
                <p><strong>Type:</strong> ${platform.category}</p>
                <p><strong>Domain:</strong> ${this.getDomainIcon(platform.domain)} ${platform.domain}</p>
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

    selectPlatform(platform, element) {
        // Clear previous selection
        document.querySelectorAll('.platform-item').forEach(item => {
            item.classList.remove('selected');
        });

        // Select new platform
        element.classList.add('selected');
        this.selectedPlatform = platform;
        this.updateStatus(`Selected ${platform.name} - Click on map to place`);

        // Update map instructions
        document.getElementById('mapInstructions').textContent =
            `Click on the map to place ${platform.name}`;
    }

    showPlatformConfigModal(latlng) {
        const modal = document.getElementById('platformModal');
        const platform = this.selectedPlatform;

        // Pre-fill modal with platform data
        document.getElementById('modalTitle').textContent = `Configure ${platform.name}`;
        document.getElementById('modalPlatformId').value =
            `${platform.id.toUpperCase()}_${String(this.platformCounter).padStart(3, '0')}`;
        document.getElementById('modalPlatformName').value =
            this.generatePlatformName(platform);
        document.getElementById('modalLatitude').value = latlng.lat.toFixed(6);
        document.getElementById('modalLongitude').value = latlng.lng.toFixed(6);

        // Set default altitude based on platform domain
        const defaultAltitudes = {
            airborne: 10000,
            maritime: 0,
            land: 100,
            space: 400000
        };
        document.getElementById('modalAltitude').value =
            defaultAltitudes[platform.domain] || 1000;

        // Store the coordinates for later use
        this.currentMarker = latlng;

        modal.style.display = 'block';
    }

    generatePlatformName(platform) {
        const nameTemplates = {
            airborne: ['United', 'Delta', 'American', 'Southwest', 'Air Force'],
            maritime: ['USS', 'USNS', 'MV', 'MS'],
            land: ['Convoy', 'Transport', 'Mobile'],
            space: ['ISS', 'Satellite', 'Station']
        };

        const templates = nameTemplates[platform.domain] || ['Vehicle'];
        const randomTemplate = templates[Math.floor(Math.random() * templates.length)];
        const randomNumber = Math.floor(Math.random() * 999) + 1;

        return `${randomTemplate} ${randomNumber}`;
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
        this.updateStatus(`Added ${platformName} to scenario`);

        // Close modal and clear selection
        document.getElementById('platformModal').style.display = 'none';
        this.clearPlatformSelection();
    }

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

    clearPlatformSelection() {
        document.querySelectorAll('.platform-item').forEach(item => {
            item.classList.remove('selected');
        });
        this.selectedPlatform = null;
        document.getElementById('mapInstructions').textContent = 'Click on the map to place platforms';
        this.updateStatus('Ready to build scenarios');
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

    removePlatform(index) {
        // Remove marker from map
        const markerData = this.mapMarkers.find(m => m.platform === this.scenarioPlatforms[index]);
        if (markerData) {
            this.map.removeLayer(markerData.marker);
            this.mapMarkers = this.mapMarkers.filter(m => m !== markerData);
        }

        // Remove from scenario
        this.scenarioPlatforms.splice(index, 1);
        this.renderScenarioPlatforms();
        this.updateStatus('Platform removed from scenario');
    }

    filterPlatforms(searchTerm) {
        const items = document.querySelectorAll('.platform-item');
        items.forEach(item => {
            const text = item.textContent.toLowerCase();
            if (text.includes(searchTerm.toLowerCase())) {
                item.style.display = 'block';
            } else {
                item.style.display = 'none';
            }
        });
    }

    filterPlatformsByDomain(domain) {
        if (domain === 'all') {
            this.renderPlatformList();
            return;
        }

        const filteredPlatforms = this.platforms.filter(p => p.domain === domain);
        const container = document.getElementById('platformList');
        container.innerHTML = '';

        filteredPlatforms.forEach(platform => {
            const item = document.createElement('div');
            item.className = 'platform-item';
            item.innerHTML = `
                <h4>${platform.name}</h4>
                <p><strong>Type:</strong> ${platform.category}</p>
                <p><strong>Domain:</strong> ${this.getDomainIcon(platform.domain)} ${platform.domain}</p>
                <p>${platform.description}</p>
            `;

            item.addEventListener('click', () => {
                this.selectPlatform(platform, item);
            });

            container.appendChild(item);
        });
    }

    generateScenarioYAML() {
        const scenarioName = document.getElementById('scenarioName').value || 'Custom Scenario';
        const scenarioDescription = document.getElementById('scenarioDescription').value || 'User-created scenario';
        const scenarioDuration = parseInt(document.getElementById('scenarioDuration').value) || 30;

        const yaml = `# Generated Scenario Configuration
metadata:
  name: "${scenarioName}"
  description: "${scenarioDescription}"
  duration: ${scenarioDuration * 60}  # ${scenarioDuration} minutes in seconds
  created: "${new Date().toISOString()}"
  author: "Scenario Builder"

platforms:
${this.scenarioPlatforms.map(platform => `  - id: "${platform.id}"
    type: "${platform.type}"
    name: "${platform.name}"
    start_position:
      latitude: ${platform.start_position.latitude}
      longitude: ${platform.start_position.longitude}
      altitude: ${platform.start_position.altitude}
    mission:
      type: "${platform.mission.type}"`).join('\n')}

# Additional scenario parameters
environment:
  weather: "clear"
  visibility: 10000  # meters
  wind_speed: 5      # m/s
  wind_direction: 270 # degrees

simulation:
  update_interval: "100ms"
  time_scale: 1.0
  physics_enabled: true
`;

        return yaml;
    }

    updateStatus(message) {
        document.getElementById('statusBar').textContent = message;
        setTimeout(() => {
            document.getElementById('statusBar').textContent = 'Ready to build scenarios';
        }, 3000);
    }

    toggleWaypointMode(enabled) {
        this.waypointMode = enabled;
        if (enabled) {
            this.updateStatus('Waypoint mode enabled - Click multiple points to create a route');
            document.getElementById('mapInstructions').textContent =
                'Click multiple points on the map to create a route, then click "Complete Route"';
        } else {
            this.clearCurrentRoute();
            this.updateStatus('Waypoint mode disabled');
            document.getElementById('mapInstructions').textContent =
                'Select a platform and click on the map to place';
        }
        this.updateRouteControls();
    }

    updateRouteControls() {
        const controls = document.getElementById('routeControls');
        if (controls) {
            controls.style.display = this.waypointMode ? 'block' : 'none';
        }

        const completeBtn = document.getElementById('completeRoute');
        if (completeBtn) {
            completeBtn.disabled = this.currentRoute.length < 2;
        }
    }

    addWaypoint(latlng) {
        this.currentRoute.push({
            latitude: latlng.lat,
            longitude: latlng.lng,
            timestamp: Date.now()
        });

        // Add waypoint marker
        const waypointMarker = L.circleMarker(latlng, {
            color: '#ff6b6b',
            fillColor: '#ff6b6b',
            fillOpacity: 0.7,
            radius: 6
        }).addTo(this.map);

        waypointMarker.bindPopup(`Waypoint ${this.currentRoute.length}`);

        // Update route line
        this.updateRouteLine();
        this.updateRouteControls();

        this.updateStatus(`Added waypoint ${this.currentRoute.length} - ${this.currentRoute.length < 2 ? 'Add more points' : 'Click Complete Route when finished'}`);
    }

    updateRouteLine() {
        // Remove existing route line
        if (this.currentRouteLine) {
            this.map.removeLayer(this.currentRouteLine);
        }

        if (this.currentRoute.length > 1) {
            const latLngs = this.currentRoute.map(wp => [wp.latitude, wp.longitude]);
            this.currentRouteLine = L.polyline(latLngs, {
                color: '#ff6b6b',
                weight: 3,
                opacity: 0.7,
                dashArray: '10, 10'
            }).addTo(this.map);
        }
    }

    completeCurrentRoute() {
        if (this.currentRoute.length < 2) {
            alert('Please add at least 2 waypoints to create a route');
            return;
        }

        if (!this.selectedPlatform) {
            alert('Please select a platform before creating a route');
            return;
        }

        this.showRouteConfigModal();
    }

    showRouteConfigModal() {
        const modal = document.getElementById('routeModal');
        const platform = this.selectedPlatform;

        // Pre-fill modal with platform and route data
        document.getElementById('routeModalTitle').textContent = `Configure Route for ${platform.name}`;
        document.getElementById('routePlatformId').value =
            `${platform.id.toUpperCase()}_${String(this.platformCounter).padStart(3, '0')}`;
        document.getElementById('routePlatformName').value =
            this.generatePlatformName(platform);

        // Set route summary
        document.getElementById('routeSummary').innerHTML = `
            <strong>Route Summary:</strong><br>
            Waypoints: ${this.currentRoute.length}<br>
            Start: ${this.currentRoute[0].latitude.toFixed(4)}, ${this.currentRoute[0].longitude.toFixed(4)}<br>
            End: ${this.currentRoute[this.currentRoute.length - 1].latitude.toFixed(4)}, ${this.currentRoute[this.currentRoute.length - 1].longitude.toFixed(4)}
        `;

        modal.style.display = 'block';
    }

    saveRouteConfig() {
        const platformId = document.getElementById('routePlatformId').value;
        const platformName = document.getElementById('routePlatformName').value;
        const routeSpeed = parseFloat(document.getElementById('routeSpeed').value);
        const routeAltitude = parseInt(document.getElementById('routeAltitude').value);
        const missionType = document.getElementById('routeMissionType').value;

        if (!platformId || !platformName || !routeSpeed) {
            alert('Please fill in all required fields');
            return;
        }

        // Create scenario platform with route
        const scenarioPlatform = {
            id: platformId,
            type: this.selectedPlatform.id,
            name: platformName,
            class: this.selectedPlatform.class,
            domain: this.selectedPlatform.domain,
            start_position: {
                latitude: this.currentRoute[0].latitude,
                longitude: this.currentRoute[0].longitude,
                altitude: routeAltitude
            },
            route: {
                waypoints: this.currentRoute,
                speed: routeSpeed,
                altitude: routeAltitude
            },
            mission: {
                type: missionType
            }
        };

        // Add to scenario
        this.scenarioPlatforms.push(scenarioPlatform);
        this.platformCounter++;

        // Add route visualization to map
        this.addRouteMarkers(scenarioPlatform);

        // Update UI
        this.renderScenarioPlatforms();
        this.updateStatus(`Added ${platformName} with ${this.currentRoute.length} waypoint route`);

        // Close modal and clear route
        document.getElementById('routeModal').style.display = 'none';
        this.clearCurrentRoute();
        this.clearPlatformSelection();
        this.toggleWaypointMode(false);
        document.getElementById('waypointMode').checked = false;
    }

    addRouteMarkers(platform) {
        const startIcon = this.getPlatformIcon(platform.domain);
        const startMarker = L.marker([platform.start_position.latitude, platform.start_position.longitude], {
            icon: startIcon
        }).addTo(this.map);

        startMarker.bindPopup(`
            <strong>${platform.name}</strong><br>
            Type: ${platform.class}<br>
            Route Speed: ${platform.route.speed} knots<br>
            Waypoints: ${platform.route.waypoints.length}<br>
            Mission: ${platform.mission.type}
        `);

        // Add route line
        const routeLatLngs = platform.route.waypoints.map(wp => [wp.latitude, wp.longitude]);
        const routeLine = L.polyline(routeLatLngs, {
            color: this.getDomainColor(platform.domain),
            weight: 3,
            opacity: 0.8
        }).addTo(this.map);

        // Add waypoint markers
        platform.route.waypoints.forEach((waypoint, index) => {
            if (index > 0) { // Skip start waypoint (already marked)
                const waypointMarker = L.circleMarker([waypoint.latitude, waypoint.longitude], {
                    color: this.getDomainColor(platform.domain),
                    fillColor: this.getDomainColor(platform.domain),
                    fillOpacity: 0.5,
                    radius: 4
                }).addTo(this.map);

                waypointMarker.bindPopup(`${platform.name} - Waypoint ${index + 1}`);
            }
        });

        this.mapMarkers.push({
            marker: startMarker,
            platform,
            routeLine,
            waypoints: platform.route.waypoints
        });
    }

    getDomainColor(domain) {
        const colors = {
            airborne: '#0066cc',
            maritime: '#004d99',
            land: '#009900',
            space: '#9933cc'
        };
        return colors[domain] || '#666666';
    }

    clearCurrentRoute() {
        // Remove waypoint markers
        this.map.eachLayer((layer) => {
            if (layer instanceof L.CircleMarker && layer.options.color === '#ff6b6b') {
                this.map.removeLayer(layer);
            }
        });

        // Remove route line
        if (this.currentRouteLine) {
            this.map.removeLayer(this.currentRouteLine);
            this.currentRouteLine = null;
        }

        this.currentRoute = [];
        this.updateRouteControls();
    }

    // Enhanced validation system
    validateScenario() {
        const errors = [];
        const warnings = [];

        // Check scenario basics
        if (!this.scenarioName.trim()) {
            errors.push('Scenario name is required');
        }

        if (this.platforms.length === 0) {
            warnings.push('No platforms added to scenario');
        }

        // Validate platforms
        this.platforms.forEach((platform, index) => {
            if (!platform.position) {
                errors.push(`Platform ${index + 1}: Position not set`);
            }

            if (platform.speed < 0) {
                errors.push(`Platform ${index + 1}: Speed cannot be negative`);
            }

            if (platform.waypoints && platform.waypoints.length > 1) {
                // Check for valid route
                for (let i = 0; i < platform.waypoints.length - 1; i++) {
                    const distance = this.calculateDistance(
                        platform.waypoints[i],
                        platform.waypoints[i + 1]
                    );
                    if (distance < 10) { // Less than 10 meters
                        warnings.push(`Platform ${index + 1}: Waypoints ${i + 1} and ${i + 2} are very close`);
                    }
                }
            }
        });

        // Check for platform collisions
        this.checkPlatformCollisions(warnings);

        return { errors, warnings };
    }

    checkPlatformCollisions(warnings) {
        for (let i = 0; i < this.platforms.length; i++) {
            for (let j = i + 1; j < this.platforms.length; j++) {
                const platform1 = this.platforms[i];
                const platform2 = this.platforms[j];

                if (platform1.position && platform2.position) {
                    const distance = this.calculateDistance(platform1.position, platform2.position);
                    if (distance < 100) { // Less than 100 meters
                        warnings.push(`Platforms ${i + 1} and ${j + 1} are very close (${Math.round(distance)}m apart)`);
                    }
                }
            }
        }
    }

    calculateDistance(pos1, pos2) {
        const R = 6371000; // Earth's radius in meters
        const dLat = (pos2.lat - pos1.lat) * Math.PI / 180;
        const dLon = (pos2.lng - pos1.lng) * Math.PI / 180;
        const a = Math.sin(dLat / 2) * Math.sin(dLat / 2) +
            Math.cos(pos1.lat * Math.PI / 180) * Math.cos(pos2.lat * Math.PI / 180) *
            Math.sin(dLon / 2) * Math.sin(dLon / 2);
        const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
        return R * c;
    }

    showValidationResults(errors, warnings) {
        const container = document.getElementById('validationResults') || this.createValidationContainer();
        container.innerHTML = '';

        if (errors.length > 0) {
            const errorDiv = document.createElement('div');
            errorDiv.className = 'validation-errors';
            errorDiv.innerHTML = `
                <h4>Errors:</h4>
                <ul>${errors.map(error => `<li>${error}</li>`).join('')}</ul>
            `;
            container.appendChild(errorDiv);
        }

        if (warnings.length > 0) {
            const warningDiv = document.createElement('div');
            warningDiv.className = 'validation-warnings';
            warningDiv.innerHTML = `
                <h4>Warnings:</h4>
                <ul>${warnings.map(warning => `<li>${warning}</li>`).join('')}</ul>
            `;
            container.appendChild(warningDiv);
        }

        container.style.display = (errors.length > 0 || warnings.length > 0) ? 'block' : 'none';
    }

    createValidationContainer() {
        const container = document.createElement('div');
        container.id = 'validationResults';
        container.className = 'validation-container';
        document.querySelector('.scenario-controls').appendChild(container);
        return container;
    }

    // Scenario statistics methods
    updateScenarioStats() {
        const stats = this.calculateScenarioStats();
        this.updateStatsDisplay(stats);
    }

    calculateScenarioStats() {
        const platforms = this.scenarioData.platforms;
        const stats = {
            totalPlatforms: platforms.length,
            byDomain: {},
            byType: {},
            withRoutes: 0,
            averageSpeed: 0,
            totalDistance: 0
        };

        let totalSpeed = 0;
        let speedCount = 0;

        platforms.forEach(platform => {
            // Count by domain
            stats.byDomain[platform.domain] = (stats.byDomain[platform.domain] || 0) + 1;

            // Count by type
            stats.byType[platform.type] = (stats.byType[platform.type] || 0) + 1;

            // Count platforms with routes
            if (platform.waypoints && platform.waypoints.length > 1) {
                stats.withRoutes++;
                stats.totalDistance += this.calculateRouteDistance(platform.waypoints);
            }

            // Calculate average speed
            if (platform.speed) {
                totalSpeed += platform.speed;
                speedCount++;
            }
        });

        if (speedCount > 0) {
            stats.averageSpeed = (totalSpeed / speedCount).toFixed(1);
        }

        return stats;
    }

    calculateRouteDistance(waypoints) {
        let distance = 0;
        for (let i = 1; i < waypoints.length; i++) {
            const prev = waypoints[i - 1];
            const curr = waypoints[i];
            distance += this.getDistance(prev.lat, prev.lng, curr.lat, curr.lng);
        }
        return distance;
    }

    updateStatsDisplay(stats) {
        const statsContainer = document.getElementById('scenario-stats');
        if (!statsContainer) return;

        const domainStats = Object.entries(stats.byDomain)
            .map(([domain, count]) => `${domain}: ${count}`)
            .join(', ');

        statsContainer.innerHTML = `
            <div class="stats-grid">
                <div class="stat-item">
                    <span class="stat-label">Total Platforms:</span>
                    <span class="stat-value">${stats.totalPlatforms}</span>
                </div>
                <div class="stat-item">
                    <span class="stat-label">By Domain:</span>
                    <span class="stat-value">${domainStats || 'None'}</span>
                </div>
                <div class="stat-item">
                    <span class="stat-label">With Routes:</span>
                    <span class="stat-value">${stats.withRoutes}</span>
                </div>
                <div class="stat-item">
                    <span class="stat-label">Avg Speed:</span>
                    <span class="stat-value">${stats.averageSpeed} kts</span>
                </div>
                <div class="stat-item">
                    <span class="stat-label">Total Distance:</span>
                    <span class="stat-value">${stats.totalDistance.toFixed(1)} nm</span>
                </div>
            </div>
        `;
    }

    // Enhanced bulk operations
    setupBulkOperations() {
        // Select all platforms of same type
        document.addEventListener('click', (e) => {
            if (e.target.classList.contains('select-by-type')) {
                const type = e.target.dataset.type;
                this.selectPlatformsByType(type);
            }
        });

        // Bulk delete selected platforms
        const bulkDeleteBtn = document.getElementById('bulk-delete-btn');
        if (bulkDeleteBtn) {
            bulkDeleteBtn.addEventListener('click', () => {
                this.bulkDeleteSelectedPlatforms();
            });
        }

        // Bulk update selected platforms
        const bulkUpdateBtn = document.getElementById('bulk-update-btn');
        if (bulkUpdateBtn) {
            bulkUpdateBtn.addEventListener('click', () => {
                this.showBulkUpdateModal();
            });
        }
    }

    selectPlatformsByType(type) {
        this.selectedPlatforms = this.scenarioData.platforms
            .filter(platform => platform.type === type)
            .map(platform => platform.id);
        this.updatePlatformSelection();
    }

    bulkDeleteSelectedPlatforms() {
        if (this.selectedPlatforms.length === 0) {
            this.showNotification('No platforms selected', 'warning');
            return;
        }

        if (confirm(`Delete ${this.selectedPlatforms.length} selected platforms?`)) {
            this.selectedPlatforms.forEach(id => {
                this.removePlatform(id);
            });
            this.selectedPlatforms = [];
            this.updatePlatformSelection();
            this.showNotification(`Deleted ${this.selectedPlatforms.length} platforms`, 'success');
        }
    }

    // Advanced scenario validation
    validateScenarioAdvanced() {
        const issues = [];
        const platforms = this.scenarioData.platforms;

        // Check for platform collisions
        const collisions = this.detectPlatformCollisions(platforms);
        if (collisions.length > 0) {
            issues.push(`${collisions.length} potential platform collisions detected`);
        }

        // Check for unrealistic speeds
        const speedIssues = platforms.filter(p => {
            const maxSpeed = this.getMaxSpeedForPlatform(p);
            return p.speed > maxSpeed;
        });
        if (speedIssues.length > 0) {
            issues.push(`${speedIssues.length} platforms have unrealistic speeds`);
        }

        // Check for platforms outside valid areas
        const outsideArea = platforms.filter(p => !this.isValidLocation(p.position));
        if (outsideArea.length > 0) {
            issues.push(`${outsideArea.length} platforms are outside valid simulation area`);
        }

        return issues;
    }

    detectPlatformCollisions(platforms) {
        const collisions = [];
        const COLLISION_THRESHOLD = 0.001; // degrees (roughly 100m)

        for (let i = 0; i < platforms.length; i++) {
            for (let j = i + 1; j < platforms.length; j++) {
                const distance = this.getDistance(
                    platforms[i].position.lat, platforms[i].position.lng,
                    platforms[j].position.lat, platforms[j].position.lng
                );

                if (distance < COLLISION_THRESHOLD) {
                    collisions.push([platforms[i], platforms[j]]);
                }
            }
        }

        return collisions;
    }

    getMaxSpeedForPlatform(platform) {
        const speedLimits = {
            land: { commercial: 100, military: 80 },
            maritime: { commercial: 30, military: 40 },
            airborne: { commercial: 500, military: 800 },
            space: { commercial: 28000, military: 28000 }
        };

        return speedLimits[platform.domain]?.[platform.category] || 100;
    }
}

// Global functions for button clicks
function exportScenario() {
    scenarioBuilder.exportScenario();
}

function loadScenario() {
    scenarioBuilder.loadScenario();
}

function loadExistingScenario(scenarioName) {
    scenarioBuilder.loadExistingScenario(scenarioName);
}

function previewScenario() {
    scenarioBuilder.previewScenario();
}

function clearScenario() {
    scenarioBuilder.clearScenario();
}

function savePlatformConfig() {
    scenarioBuilder.savePlatformConfig();
}

function downloadYaml() {
    scenarioBuilder.downloadYaml();
}

// Extend ScenarioBuilder with export/import functionality
ScenarioBuilder.prototype.exportScenario = function () {
    if (this.scenarioPlatforms.length === 0) {
        alert('No platforms in scenario to export');
        return;
    }

    const yaml = this.generateScenarioYAML();
    this.downloadFile(yaml, 'scenario.yaml', 'text/yaml');
    this.updateStatus('Scenario exported successfully');
};

ScenarioBuilder.prototype.previewScenario = function () {
    const yaml = this.generateScenarioYAML();
    document.getElementById('yamlPreview').textContent = yaml;
    document.getElementById('previewModal').style.display = 'block';
};

ScenarioBuilder.prototype.downloadYaml = function () {
    const yaml = document.getElementById('yamlPreview').textContent;
    this.downloadFile(yaml, 'scenario.yaml', 'text/yaml');
    document.getElementById('previewModal').style.display = 'none';
    this.updateStatus('YAML downloaded');
};

ScenarioBuilder.prototype.loadScenario = function () {
    const input = document.createElement('input');
    input.type = 'file';
    input.accept = '.yaml,.yml';
    input.onchange = (e) => {
        const file = e.target.files[0];
        if (file) {
            const reader = new FileReader();
            reader.onload = (e) => {
                try {
                    this.parseAndLoadYAML(e.target.result);
                } catch (error) {
                    alert('Error loading scenario: ' + error.message);
                }
            };
            reader.readAsText(file);
        }
    };
    input.click();
};

ScenarioBuilder.prototype.loadExistingScenario = function (scenarioName) {
    fetch(`/data/configs/${scenarioName}.yaml`)
        .then(response => {
            if (!response.ok) throw new Error('Network response was not ok');
            return response.text();
        })
        .then(yamlText => {
            this.parseAndLoadYAML(yamlText);
            this.updateStatus(`Loaded existing scenario: ${scenarioName}`);
        })
        .catch(error => {
            console.error('Error loading existing scenario:', error);
            alert('Failed to load scenario. Please check the console for details.');
        });
};

ScenarioBuilder.prototype.parseAndLoadYAML = function (yamlText) {
    // Simple YAML parser for basic scenario structure
    const lines = yamlText.split('\n');
    let currentPlatform = null;
    const platforms = [];

    // This is a simplified parser - in production, you'd use a proper YAML library
    lines.forEach(line => {
        const trimmed = line.trim();
        if (trimmed.startsWith('- id:')) {
            if (currentPlatform) platforms.push(currentPlatform);
            currentPlatform = { start_position: {}, mission: {} };
            currentPlatform.id = trimmed.split('"')[1];
        } else if (trimmed.startsWith('type:') && currentPlatform) {
            currentPlatform.type = trimmed.split('"')[1];
        } else if (trimmed.startsWith('name:') && currentPlatform) {
            currentPlatform.name = trimmed.split('"')[1];
        } else if (trimmed.startsWith('latitude:') && currentPlatform) {
            currentPlatform.start_position.latitude = parseFloat(trimmed.split(':')[1]);
        } else if (trimmed.startsWith('longitude:') && currentPlatform) {
            currentPlatform.start_position.longitude = parseFloat(trimmed.split(':')[1]);
        } else if (trimmed.startsWith('altitude:') && currentPlatform) {
            currentPlatform.start_position.altitude = parseInt(trimmed.split(':')[1]);
        }
    });

    if (currentPlatform) platforms.push(currentPlatform);

    // Load platforms into scenario
    this.clearScenario();
    platforms.forEach(platform => {
        const platformType = this.platforms.find(p => p.id === platform.type);
        if (platformType) {
            const scenarioPlatform = {
                ...platform,
                class: platformType.class,
                domain: platformType.domain,
                mission: platform.mission || { type: 'patrol' }
            };
            this.scenarioPlatforms.push(scenarioPlatform);
            this.addMapMarker(scenarioPlatform);
        }
    });

    this.renderScenarioPlatforms();
    this.updateStatus(`Loaded ${platforms.length} platforms from scenario`);
};

ScenarioBuilder.prototype.clearScenario = function () {
    // Clear scenario platforms
    this.scenarioPlatforms = [];

    // Clear map markers
    this.mapMarkers.forEach(markerData => {
        this.map.removeLayer(markerData.marker);
    });
    this.mapMarkers = [];

    // Reset counter
    this.platformCounter = 1;

    // Clear selection
    this.clearPlatformSelection();

    // Update UI
    this.renderScenarioPlatforms();
    this.updateStatus('Scenario cleared');
};

ScenarioBuilder.prototype.downloadFile = function (content, filename, contentType) {
    const a = document.createElement('a');
    const file = new Blob([content], { type: contentType });
    a.href = URL.createObjectURL(file);
    a.download = filename;
    a.click();
    URL.revokeObjectURL(a.href);
};

ScenarioBuilder.prototype.runSimulation = function () {
    if (this.scenarioPlatforms.length === 0) {
        alert('Please add platforms to the scenario before running the simulation');
        return;
    }

    const timestep = document.getElementById('timestep').value;
    const multicastGroup = document.getElementById('multicastGroup').value;
    const multicastPort = document.getElementById('multicastPort').value;

    // Validate inputs
    if (!timestep || timestep < 1 || timestep > 60) {
        alert('Please enter a valid timestep between 1 and 60 seconds');
        return;
    }

    if (!this.isValidMulticastAddress(multicastGroup)) {
        alert('Please enter a valid multicast IP address (e.g., 239.255.0.1)');
        return;
    }

    if (!multicastPort || multicastPort < 1024 || multicastPort > 65535) {
        alert('Please enter a valid multicast port between 1024 and 65535');
        return;
    }

    // Prepare simulation configuration
    const simConfig = {
        scenario: {
            name: document.getElementById('scenarioName').value || 'Custom Scenario',
            description: document.getElementById('scenarioDescription').value || 'User-created scenario',
            duration: parseInt(document.getElementById('scenarioDuration').value) * 60,
            platforms: this.scenarioPlatforms
        },
        simulation: {
            timestep: parseInt(timestep),
            multicast: {
                group: multicastGroup,
                port: parseInt(multicastPort)
            }
        }
    };

    this.updateStatus('Starting simulation...');
    this.startSimulation(simConfig);
};

ScenarioBuilder.prototype.isValidMulticastAddress = function (ip) {
    const parts = ip.split('.');
    if (parts.length !== 4) return false;

    const firstOctet = parseInt(parts[0]);
    return firstOctet >= 224 && firstOctet <= 239 &&
        parts.every(part => {
            const num = parseInt(part);
            return num >= 0 && num <= 255;
        });
};

ScenarioBuilder.prototype.startSimulation = function (config) {
    // Send simulation configuration to the server
    fetch('/api/simulation/start', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(config)
    })
        .then(response => {
            if (response.ok) {
                this.updateStatus('Simulation started successfully! Redirecting to live view...');

                // Wait a moment then redirect to the live map
                setTimeout(() => {
                    window.location.href = '/';
                }, 2000);
            } else {
                throw new Error('Failed to start simulation');
            }
        })
        .catch(error => {
            console.error('Error starting simulation:', error);
            this.updateStatus('Error: Failed to start simulation. Please check your configuration.');
            alert('Failed to start simulation. Please check the console for details.');
        });
};

// Global function for run button
function runSimulation() {
    scenarioBuilder.runSimulation();
}

// Initialize the scenario builder when the page loads
let scenarioBuilder;
document.addEventListener('DOMContentLoaded', () => {
    scenarioBuilder = new ScenarioBuilder();
});
