/**
 * High-Performance Platform Renderer
 * Optimized for rendering 50,000+ dynamic moving objects
 */
class PlatformRenderer {
    constructor(mapEngine) {
        this.mapEngine = mapEngine;
        this.map = mapEngine.getMap();
        this.platformLayer = mapEngine.getPlatformLayer();
        this.trailLayer = mapEngine.getTrailLayer();

        // Platform management
        this.platforms = new Map(); // platformId -> platform data
        this.markers = new Map();   // platformId -> marker
        this.trails = new Map();    // platformId -> trail polyline
        this.trailPoints = new Map(); // platformId -> array of trail points

        // Layer management for different platform types
        this.layers = new Map();
        this.layers.set('airborne', L.layerGroup());
        this.layers.set('maritime', L.layerGroup());
        this.layers.set('land', L.layerGroup());
        this.layers.set('space', L.layerGroup());

        // Performance optimizations
        this.useCanvas = true;
        this.canvasRenderer = null;
        this.clusteringEnabled = false;
        this.markerCluster = null;
        this.visibilityFilters = {
            airborne: true,
            maritime: true,
            land: true,
            space: true
        };

        // Rendering settings
        this.maxTrailLength = 20;
        this.showTrails = true;
        this.markerSize = 8;
        this.defaultOpacity = 1.0;
        this.zoomBasedSizing = true;

        // Performance tracking
        this.renderStats = {
            totalPlatforms: 0,
            visiblePlatforms: 0,
            renderTime: 0,
            lastUpdate: 0
        };

        this.init();
    }

    async init() {
        console.log('Initializing PlatformRenderer...');

        // Set up high-performance canvas renderer
        this.setupCanvasRenderer();

        // Set up clustering if available
        this.setupClustering();

        // Set up event handlers
        this.setupEventHandlers();

        console.log('PlatformRenderer initialized');

        return this;
    }

    setupCanvasRenderer() {
        // Use Leaflet Canvas renderer for better performance with many markers
        if (L.Canvas && this.useCanvas) {
            this.canvasRenderer = L.canvas({
                padding: 0.5,
                tolerance: 10,
                updateWhenIdle: false,
                updateWhenZooming: true
            });
            this.map.addLayer(this.canvasRenderer);
        }
    }

    setupClustering() {
        // Initialize marker clustering for extreme zoom-out scenarios
        if (typeof L.markerClusterGroup !== 'undefined') {
            this.markerCluster = L.markerClusterGroup({
                maxClusterRadius: 50,
                spiderfyOnMaxZoom: false,
                showCoverageOnHover: false,
                zoomToBoundsOnClick: false,
                disableClusteringAtZoom: 10,
                // Performance optimizations
                animate: false,
                chunkedLoading: true,
                chunkInterval: 50,
                chunkDelay: 50,
                // Custom cluster icon
                iconCreateFunction: (cluster) => {
                    const count = cluster.getChildCount();
                    let className = 'marker-cluster-';
                    if (count < 10) {
                        className += 'small';
                    } else if (count < 100) {
                        className += 'medium';
                    } else {
                        className += 'large';
                    }

                    return new L.DivIcon({
                        html: `<div><span>${count}</span></div>`,
                        className: 'marker-cluster ' + className,
                        iconSize: new L.Point(40, 40)
                    });
                }
            });
        }
    }

    setupEventHandlers() {
        // Update marker sizes based on zoom level
        this.map.on('zoomend', () => {
            if (this.zoomBasedSizing) {
                this.updateMarkerSizes();
            }
        });

        // Viewport change optimization - fix the callback assignment
        if (this.mapEngine.addViewportChangeListener) {
            this.mapEngine.addViewportChangeListener((bounds) => {
                this.onViewportChange(bounds);
            });
        }
    }

    // Create optimized marker for platform
    createPlatformMarker(platform) {
        const { latitude, longitude } = platform.position;

        // Validate coordinates
        if (latitude < -90 || latitude > 90 || longitude < -180 || longitude > 180) {
            console.warn(`Invalid coordinates for platform ${platform.id}: ${latitude}, ${longitude}`);
            return null;
        }

        const position = [latitude, longitude];
        const icon = this.createPlatformIcon(platform);

        // Create marker with proper options
        const marker = L.marker(position, {
            icon: icon,
            title: `${platform.platform_type} - ${platform.id}`,
            alt: platform.id,
            opacity: this.defaultOpacity,
            riseOnHover: true,
            riseOffset: 250
        });

        // Add platform data to marker - ensure marker object exists first
        if (marker) {
            marker.platformData = platform;

            // Add popup with platform information
            this.attachPlatformPopup(marker, platform);

            // Store marker reference
            this.markers.set(platform.id, marker);

            // Add to appropriate layer
            const layer = this.layers.get(platform.platform_type);
            if (layer && this.map) {
                layer.addLayer(marker);
            }
        }

        return marker;
    }

    createPlatformIcon(platform) {
        const size = this.calculateMarkerSize(platform);
        const color = this.getPlatformColor(platform.platform_type);

        return L.divIcon({
            className: `platform-marker ${platform.platform_type}`,
            html: `<div style="
                width: ${size * 2}px;
                height: ${size * 2}px;
                background-color: ${color};
                border-radius: 50%;
                border: 2px solid rgba(255,255,255,0.8);
                box-shadow: 0 0 10px rgba(0,0,0,0.3);
            "></div>`,
            iconSize: [size * 2, size * 2],
            iconAnchor: [size, size]
        });
    }

    calculateMarkerSize(platform) {
        let baseSize = this.markerSize;

        if (this.zoomBasedSizing) {
            const zoom = this.map.getZoom();
            // Scale marker size based on zoom level
            baseSize = Math.max(3, Math.min(20, this.markerSize * (zoom / 10)));
        }

        // Adjust size based on platform type
        switch (platform.platform_type) {
            case 'space':
                return baseSize * 0.7; // Smaller for satellites
            case 'airborne':
                return baseSize;
            case 'maritime':
                return baseSize * 1.2; // Larger for ships
            case 'land':
                return baseSize * 0.9;
            default:
                return baseSize;
        }
    }

    getPlatformColor(platformType) {
        const colors = {
            airborne: '#2196F3',  // Blue
            maritime: '#00BCD4',  // Cyan
            land: '#4CAF50',      // Green
            space: '#9C27B0'      // Purple
        };
        return colors[platformType] || '#757575';
    }

    attachPlatformPopup(marker, platform) {
        const popupContent = this.createPopupContent(platform);
        marker.bindPopup(popupContent, {
            className: 'platform-popup',
            maxWidth: 300,
            keepInView: true,
            autoPan: false // Disable for performance
        });
    }

    createPopupContent(platform) {
        const pos = platform.position;
        const vel = platform.velocity || { north: 0, east: 0, up: 0 };
        const altitude = pos.altitude || 0;
        const speed = platform.speed || 0;
        const heading = platform.heading || 0;

        return `
            <div class="platform-popup">
                <h4>${platform.id} - ${platform.platform_type}</h4>
                <div class="popup-field">
                    <span class="popup-label">Position:</span>
                    <span class="popup-value">${pos.latitude.toFixed(4)}, ${pos.longitude.toFixed(4)}</span>
                </div>
                <div class="popup-field">
                    <span class="popup-label">Altitude:</span>
                    <span class="popup-value">${altitude.toFixed(0)} m</span>
                </div>
                <div class="popup-field">
                    <span class="popup-label">Speed:</span>
                    <span class="popup-value">${speed.toFixed(1)} m/s</span>
                </div>
                <div class="popup-field">
                    <span class="popup-label">Heading:</span>
                    <span class="popup-value">${heading.toFixed(0)}°</span>
                </div>
                <div class="popup-field">
                    <span class="popup-label">Velocity:</span>
                    <span class="popup-value">N:${vel.north.toFixed(1)} E:${vel.east.toFixed(1)} U:${vel.up.toFixed(1)}</span>
                </div>
                <div class="popup-field">
                    <span class="popup-label">Last Updated:</span>
                    <span class="popup-value">${new Date(platform.lastUpdated || Date.now()).toLocaleTimeString()}</span>
                </div>
            </div>
        `;
    }

    // Add or update platform on map
    updatePlatform(platform) {
        const platformId = platform.id;
        const position = [platform.position.latitude, platform.position.longitude];

        // Check if platform should be visible
        if (!this.isPlatformVisible(platform)) {
            this.removePlatform(platformId);
            return;
        }

        // Check viewport culling for performance
        if (!this.mapEngine.isInViewport(platform.position.latitude, platform.position.longitude, 0.2)) {
            // Platform is outside viewport, remove from rendering but keep in memory
            if (this.markers.has(platformId)) {
                const marker = this.markers.get(platformId);
                if (this.clusteringEnabled && this.markerCluster) {
                    this.markerCluster.removeLayer(marker);
                } else {
                    this.platformLayer.removeLayer(marker);
                }
            }
            return;
        }

        // Store platform data
        this.platforms.set(platformId, platform);

        let marker = this.markers.get(platformId);

        if (marker) {
            // Update existing marker position
            marker.setLatLng(position);
            marker.platformData = platform;

            // Update popup content if open
            if (marker.isPopupOpen()) {
                marker.setPopupContent(this.createPopupContent(platform));
            }
        } else {
            // Create new marker
            marker = this.createPlatformMarker(platform);
            this.markers.set(platformId, marker);

            // Add to appropriate layer
            if (this.clusteringEnabled && this.markerCluster) {
                this.markerCluster.addLayer(marker);
            } else {
                this.platformLayer.addLayer(marker);
            }
        }

        // Update trail
        if (this.showTrails) {
            this.updatePlatformTrail(platform);
        }

        // Update render stats
        this.renderStats.totalPlatforms = this.platforms.size;
    }

    updatePlatformTrail(platform) {
        const platformId = platform.id;
        const position = [platform.position.latitude, platform.position.longitude];

        // Get or create trail points array
        let trailPoints = this.trailPoints.get(platformId) || [];

        // Add new position to trail
        trailPoints.push(position);

        // Limit trail length for performance
        if (trailPoints.length > this.maxTrailLength) {
            trailPoints = trailPoints.slice(-this.maxTrailLength);
        }

        this.trailPoints.set(platformId, trailPoints);

        // Update trail polyline
        if (trailPoints.length > 1) {
            let trail = this.trails.get(platformId);

            if (trail) {
                trail.setLatLngs(trailPoints);
            } else {
                trail = L.polyline(trailPoints, {
                    color: this.getPlatformColor(platform.platform_type),
                    weight: 2,
                    opacity: 0.6,
                    className: `platform-trail ${platform.platform_type}`,
                    smoothFactor: 1,
                    interactive: false // Disable interaction for performance
                });

                this.trails.set(platformId, trail);
                this.trailLayer.addLayer(trail);
            }
        }
    }

    // Remove platform from map
    removePlatform(platformId) {
        // Remove marker
        if (this.markers.has(platformId)) {
            const marker = this.markers.get(platformId);
            if (this.clusteringEnabled && this.markerCluster) {
                this.markerCluster.removeLayer(marker);
            } else {
                this.platformLayer.removeLayer(marker);
            }
            this.markers.delete(platformId);
        }

        // Remove trail
        if (this.trails.has(platformId)) {
            const trail = this.trails.get(platformId);
            this.trailLayer.removeLayer(trail);
            this.trails.delete(platformId);
            this.trailPoints.delete(platformId);
        }

        // Remove from platforms
        this.platforms.delete(platformId);
    }

    // Batch update multiple platforms for better performance
    updatePlatforms(platforms) {
        const startTime = performance.now();

        // Track which platforms are in the update
        const updatedPlatformIds = new Set();

        platforms.forEach(platform => {
            this.updatePlatform(platform);
            updatedPlatformIds.add(platform.id);
        });

        // Remove platforms that are no longer in the update
        const platformsToRemove = [];
        this.platforms.forEach((platform, id) => {
            if (!updatedPlatformIds.has(id)) {
                platformsToRemove.push(id);
            }
        });

        platformsToRemove.forEach(id => {
            this.removePlatform(id);
        });

        // Update performance stats
        this.renderStats.renderTime = performance.now() - startTime;
        this.renderStats.lastUpdate = Date.now();
        this.renderStats.visiblePlatforms = this.countVisiblePlatforms();
    }

    // Check if platform should be visible based on filters
    isPlatformVisible(platform) {
        return this.visibilityFilters[platform.platform_type] !== false;
    }

    // Update marker sizes based on zoom level
    updateMarkerSizes() {
        this.markers.forEach((marker, platformId) => {
            const platform = this.platforms.get(platformId);
            if (platform && marker.platformData) {
                const newSize = this.calculateMarkerSize(platform);
                if (marker.setRadius) {
                    marker.setRadius(newSize);
                }
            }
        });
    }

    // Count visible platforms in viewport
    countVisiblePlatforms() {
        let count = 0;
        this.platforms.forEach(platform => {
            if (this.mapEngine.isInViewport(platform.position.latitude, platform.position.longitude) &&
                this.isPlatformVisible(platform)) {
                count++;
            }
        });
        return count;
    }

    // Viewport change handler for performance optimization
    onViewportChange(bounds) {
        // Re-evaluate which platforms should be visible
        const visibleCount = this.countVisiblePlatforms();
        this.renderStats.visiblePlatforms = visibleCount;

        // Enable clustering for high-density areas
        if (visibleCount > 1000 && !this.clusteringEnabled) {
            this.enableClustering();
        } else if (visibleCount < 500 && this.clusteringEnabled) {
            this.disableClustering();
        }
    }

    // Platform click handler
    onPlatformClick(platform, event) {
        // Override in external components
        if (this.onPlatformClickCallback) {
            this.onPlatformClickCallback(platform, event);
        }
    }

    // Public API methods

    // Set platform visibility filters
    setPlatformFilter(platformType, visible) {
        this.visibilityFilters[platformType] = visible;

        // Update existing platforms
        this.platforms.forEach(platform => {
            if (platform.platform_type === platformType) {
                if (visible) {
                    this.updatePlatform(platform);
                } else {
                    this.removePlatform(platform.id);
                }
            }
        });
    }

    // Enable/disable clustering
    enableClustering() {
        if (this.markerCluster && !this.clusteringEnabled) {
            this.clusteringEnabled = true;

            // Move markers to cluster group
            this.markers.forEach(marker => {
                this.platformLayer.removeLayer(marker);
                this.markerCluster.addLayer(marker);
            });

            this.map.addLayer(this.markerCluster);
            console.log('Clustering enabled');
        }
    }

    disableClustering() {
        if (this.markerCluster && this.clusteringEnabled) {
            this.clusteringEnabled = false;

            // Move markers back to platform layer
            this.markers.forEach(marker => {
                this.markerCluster.removeLayer(marker);
                this.platformLayer.addLayer(marker);
            });

            this.map.removeLayer(this.markerCluster);
            console.log('Clustering disabled');
        }
    }

    // Show/hide trails
    setTrailsVisible(visible) {
        this.showTrails = visible;

        if (visible) {
            this.map.addLayer(this.trailLayer);
        } else {
            this.map.removeLayer(this.trailLayer);
        }
    }

    // Set trail length
    setTrailLength(length) {
        this.maxTrailLength = length;

        // Update existing trails
        this.trailPoints.forEach((points, platformId) => {
            if (points.length > length) {
                const trimmedPoints = points.slice(-length);
                this.trailPoints.set(platformId, trimmedPoints);

                const trail = this.trails.get(platformId);
                if (trail) {
                    trail.setLatLngs(trimmedPoints);
                }
            }
        });
    }

    // Clear all platforms
    clearAllPlatforms() {
        this.platforms.clear();
        this.markers.clear();
        this.trails.clear();
        this.trailPoints.clear();
        this.platformLayer.clearLayers();
        this.trailLayer.clearLayers();

        if (this.markerCluster) {
            this.markerCluster.clearLayers();
        }
    }

    // Get render statistics
    getRenderStats() {
        return { ...this.renderStats };
    }

    // Focus on specific platform
    focusOnPlatform(platformId) {
        const platform = this.platforms.get(platformId);
        if (platform) {
            this.mapEngine.centerOn(
                platform.position.latitude,
                platform.position.longitude,
                12
            );

            // Open popup
            const marker = this.markers.get(platformId);
            if (marker) {
                marker.openPopup();
            }
        }
    }
}

// Export for use in other modules
window.PlatformRenderer = PlatformRenderer;
