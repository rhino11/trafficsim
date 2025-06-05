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
                description: 'Short to medium-range commercial airliner',
                performance: { max_speed: 257.0, cruise_speed: 230.0, max_altitude: 12000 }
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
                id: 'tesla_model_s',
                name: 'Tesla Model S',
                class: 'Tesla Model S',
                category: 'passenger_vehicle',
                domain: 'land',
                description: 'Electric luxury sedan',
                performance: { max_speed: 69.4, cruise_speed: 33.3 }
            },
            {
                id: 'starlink_satellite',
                name: 'Starlink Satellite',
                class: 'Starlink Satellite',
                category: 'communications_satellite',
                domain: 'space',
                description: 'Low Earth orbit communications satellite',
                performance: { max_speed: 7660.0, cruise_speed: 7660.0, max_altitude: 550000 }
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
        // Remove selected class from all platform items
        document.querySelectorAll('.platform-item').forEach(p => {
            p.classList.remove('selected');
        });

        // Add selected class to clicked platform
        element.classList.add('selected');

        // Set the selected platform
        this.selectedPlatform = platform;

        // Update UI
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
            route: []
        };

        this.scenarioPlatforms.push(platformData);

        // Add marker to map
        const marker = L.marker([latlng.lat, latlng.lng])
            .bindPopup(`${platform.name}<br>ID: ${platformData.id}`)
            .addTo(this.map);

        this.mapMarkers.push(marker);
        this.updateScenarioPlatformsList();
        this.updateStatus(`Placed ${platform.name} at ${latlng.lat.toFixed(4)}, ${latlng.lng.toFixed(4)}`);
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
        const waypointMarker = L.marker([latlng.lat, latlng.lng], {
            icon: L.divIcon({
                className: 'waypoint-marker',
                html: `<div style="background-color: red; width: 8px; height: 8px; border-radius: 50%; border: 2px solid white;"></div>`,
                iconSize: [12, 12],
                iconAnchor: [6, 6]
            })
        }).addTo(this.map);

        // Draw polyline if we have more than one waypoint
        if (this.currentRoute.length > 1) {
            const polyline = L.polyline(this.currentRoute, {
                color: 'red',
                weight: 3,
                opacity: 0.7
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
        a.download = `${document.getElementById('scenarioName').value || 'scenario'}.yaml`;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);

        this.updateStatus('Scenario exported successfully');
    }

    loadScenario(file) {
        const reader = new FileReader();
        reader.onload = (e) => {
            try {
                // Simple YAML parsing for basic scenarios
                const content = e.target.result;
                const lines = content.split('\n');

                // This is a simplified parser - in production you'd use a proper YAML library
                let currentPlatform = null;
                const platforms = [];

                lines.forEach(line => {
                    const trimmed = line.trim();
                    if (trimmed.startsWith('- id:')) {
                        if (currentPlatform) platforms.push(currentPlatform);
                        currentPlatform = { id: trimmed.split('"')[1] };
                    } else if (currentPlatform) {
                        if (trimmed.startsWith('type:')) currentPlatform.type = trimmed.split('"')[1];
                        else if (trimmed.startsWith('name:')) currentPlatform.name = trimmed.split('"')[1];
                        else if (trimmed.startsWith('latitude:')) currentPlatform.latitude = parseFloat(trimmed.split(':')[1]);
                        else if (trimmed.startsWith('longitude:')) currentPlatform.longitude = parseFloat(trimmed.split(':')[1]);
                        else if (trimmed.startsWith('altitude:')) currentPlatform.altitude = parseInt(trimmed.split(':')[1]);
                    }
                });

                if (currentPlatform) platforms.push(currentPlatform);

                // Load platforms into scenario
                this.scenarioPlatforms = platforms.map(p => ({
                    id: p.id,
                    type: p.type,
                    name: p.name,
                    start_position: {
                        latitude: p.latitude,
                        longitude: p.longitude,
                        altitude: p.altitude
                    },
                    mission: { type: 'patrol' }
                }));

                this.renderScenarioPlatforms();
                this.updateStatus(`Loaded ${platforms.length} platforms from scenario`);

            } catch (error) {
                console.error('Error loading scenario:', error);
                alert('Error loading scenario file');
            }
        };
        reader.readAsText(file);
    }

    clearScenario() {
        if (confirm('Are you sure you want to clear the current scenario?')) {
            this.scenarioPlatforms = [];
            this.mapMarkers.forEach(markerData => {
                this.map.removeLayer(markerData.marker);
            });
            this.mapMarkers = [];
            this.renderScenarioPlatforms();
            this.updateStatus('Scenario cleared');
        }
    }

    previewScenario() {
        if (this.scenarioPlatforms.length === 0) {
            alert('Please add platforms to preview');
            return;
        }

        // Fit map to show all platforms
        const group = new L.featureGroup(this.mapMarkers.map(m => m.marker));
        this.map.fitBounds(group.getBounds().pad(0.1));

        this.updateStatus('Scenario preview updated');
    }

    getDomainStats() {
        const stats = {
            airborne: 0,
            maritime: 0,
            land: 0,
            space: 0
        };

        this.scenarioPlatforms.forEach(platform => {
            if (stats.hasOwnProperty(platform.domain)) {
                stats[platform.domain]++;
            }
        });

        return stats;
    }

    validateScenario() {
        const issues = [];

        if (this.scenarioPlatforms.length === 0) {
            issues.push('No platforms added to scenario');
        }

        // Check for overlapping platforms
        for (let i = 0; i < this.scenarioPlatforms.length; i++) {
            for (let j = i + 1; j < this.scenarioPlatforms.length; j++) {
                const p1 = this.scenarioPlatforms[i].start_position;
                const p2 = this.scenarioPlatforms[j].start_position;

                const distance = Math.sqrt(
                    Math.pow(p1.latitude - p2.latitude, 2) +
                    Math.pow(p1.longitude - p2.longitude, 2)
                );

                if (distance < 0.001) { // Very close positions
                    issues.push(`Platforms ${this.scenarioPlatforms[i].name} and ${this.scenarioPlatforms[j].name} are very close`);
                }
            }
        }

        return issues;
    }

    loadPlatformLibrary() {
        // Mock platform library with platforms for all domains
        this.platformLibrary = {
            airborne: [
                { name: 'Boeing 747', type: 'commercial', speed: 500 },
                { name: 'F-16 Fighter', type: 'military', speed: 1200 }
            ],
            maritime: [
                { name: 'Container Ship', type: 'commercial', speed: 25 },
                { name: 'Navy Destroyer', type: 'military', speed: 35 }
            ],
            land: [
                { name: 'Delivery Truck', type: 'commercial', speed: 80 },
                { name: 'Military Tank', type: 'military', speed: 60 }
            ],
            space: [
                { name: 'Commercial Satellite', type: 'commercial', speed: 7800 },
                { name: 'Military Satellite', type: 'military', speed: 7800 }
            ]
        };
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
