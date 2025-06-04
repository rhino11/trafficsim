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
}

// Global functions for button clicks
function exportScenario() {
    scenarioBuilder.exportScenario();
}

function loadScenario() {
    scenarioBuilder.loadScenario();
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
