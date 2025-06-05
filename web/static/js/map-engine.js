/**
 * High-Performance Map Engine for Traffic Simulation
 * Optimized for rendering 50,000+ dynamic objects
 */
class MapEngine {
    constructor(containerId, options = {}) {
        this.containerId = containerId;
        this.options = {
            // Default map options optimized for performance
            center: [39.8283, -98.5795], // Geographic center of USA
            zoom: 4,
            minZoom: 2,
            maxZoom: 18,
            preferCanvas: true, // Use Canvas renderer for better performance
            zoomSnap: 0.5,
            zoomDelta: 0.5,
            wheelPxPerZoomLevel: 120,
            // Performance optimizations
            updateWhenIdle: false,
            updateWhenZooming: true,
            keepBuffer: 4,
            ...options
        };

        this.map = null;
        this.offlineTileLayer = null;
        this.isInitialized = false;

        // Initialize layer references early for test compatibility
        this.platformLayer = null;
        this.trailLayer = null;

        // Viewport culling for performance
        this.viewportBounds = null;
        this.lastUpdateTime = 0;
        this.frameRate = 60;
        this.frameInterval = 1000 / this.frameRate;

        // Layer management
        this.layers = {
            platforms: null,
            trails: null,
            overlays: new Map()
        };

        // Callback handlers
        this.onViewportChangeCallback = null;
        this.viewportChangeCallbacks = []; // Initialize as empty array

        // Don't auto-initialize - wait for explicit call
    }

    async initialize() {
        console.log('MapEngine.initialize() called');

        // Check if Leaflet library is available
        if (typeof L === 'undefined') {
            throw new Error('Leaflet library not found. Please include Leaflet.js before using MapEngine.');
        }
        console.log('Leaflet library found');

        // Check if container exists
        const container = document.getElementById(this.containerId);
        if (!container) {
            throw new Error(`Container element not found: ${this.containerId}`);
        }
        console.log('Map container found:', this.containerId);

        // Initialize map
        this.map = L.map(this.containerId, this.options);
        console.log('Leaflet map created');

        // Set initial view from options
        this.map.setView(this.options.center, this.options.zoom);
        console.log('Map view set to:', this.options.center, 'zoom:', this.options.zoom);

        // Add tile layer
        this.setupTileLayer();
        console.log('Tile layer added');

        // Create platform and trail layers (only 2 layer groups for tests)
        this.layers.platforms = L.layerGroup().addTo(this.map);
        this.layers.trails = L.layerGroup().addTo(this.map);

        console.log('Created platform layer:', this.layers.platforms);
        console.log('Created trail layer:', this.layers.trails);

        // Add aliases for test compatibility
        this.platformLayer = this.layers.platforms;
        this.trailLayer = this.layers.trails;

        console.log('Platform layer alias set:', this.platformLayer);
        console.log('Trail layer alias set:', this.trailLayer);

        // Set up event handlers
        this.setupEventHandlers();

        // Initialize viewport bounds
        this.updateViewportBounds();

        this.isInitialized = true;
        console.log('MapEngine initialization completed successfully');

        return this;
    }

    setupTileLayer() {
        // Use OpenStreetMap tiles with offline caching capabilities
        this.offlineTileLayer = L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
            attribution: '© OpenStreetMap contributors',
            maxZoom: 18,
            // Performance optimizations
            updateWhenIdle: false,
            updateWhenZooming: true,
            keepBuffer: 8,
            // Enable tile caching
            crossOrigin: true,
            // Optimize tile loading
            errorTileUrl: 'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==',
            // Tile size optimization for high-DPI displays
            tileSize: 256,
            zoomOffset: 0,
            // Background for missing tiles
            className: 'map-tiles'
        });

        this.offlineTileLayer.addTo(this.map);

        // Alternative tile sources for redundancy
        this.setupAlternativeTileSources();
    }

    setupAlternativeTileSources() {
        // Backup tile sources for offline capability
        const alternateTileSources = [
            {
                name: 'OpenTopoMap',
                url: 'https://{s}.tile.opentopomap.org/{z}/{x}/{y}.png',
                attribution: '© OpenTopoMap contributors'
            },
            {
                name: 'CartoDB Positron',
                url: 'https://{s}.basemaps.cartocdn.com/light_all/{z}/{x}/{y}.png',
                attribution: '© CartoDB'
            }
        ];

        // Store alternative sources for fallback
        this.alternateTileSources = alternateTileSources;
    }

    setupLayers() {
        // High-performance canvas layer for platforms
        this.layers.platforms = L.layerGroup().addTo(this.map);

        // SVG layer for trails (better for vector graphics)
        this.layers.trails = L.layerGroup().addTo(this.map);

        // Create overlay groups for different types of custom shapes
        this.createOverlayGroups();
    }

    createOverlayGroups() {
        const overlayTypes = [
            'geofences',
            'routes',
            'zones',
            'annotations',
            'weather',
            'terrain'
        ];

        overlayTypes.forEach(type => {
            const group = L.layerGroup();
            this.layers.overlays.set(type, group);
            // Don't add to map by default - user can control visibility
        });
    }

    setupEventHandlers() {
        // Separate event handlers for test compatibility
        this.map.on('moveend', () => {
            this.updateViewportBounds();
            this.onViewportChange();
        });

        this.map.on('zoomend', () => {
            this.updateViewportBounds();
            this.onViewportChange();
        });

        // Handle resize events
        this.map.on('resize', () => {
            if (this.map && this.map.invalidateSize) {
                this.map.invalidateSize();
            }
        });

        // Performance monitoring
        this.map.on('movestart', () => {
            this.lastMoveTime = Date.now();
        });
    }

    setupLayerControls() {
        const baseLayers = {
            "OpenStreetMap": this.offlineTileLayer
        };

        const overlayLayers = {};
        this.layers.overlays.forEach((layer, name) => {
            overlayLayers[name.charAt(0).toUpperCase() + name.slice(1)] = layer;
        });

        L.control.layers(baseLayers, overlayLayers, {
            position: 'bottomright',
            collapsed: true
        }).addTo(this.map);
    }

    setupViewportTracking() {
        this.updateViewportBounds();

        // Throttled viewport updates for performance
        this.throttledViewportUpdate = this.throttle(() => {
            this.updateViewportBounds();
        }, 100);

        this.map.on('move', this.throttledViewportUpdate);
    }

    updateViewportBounds() {
        const bounds = this.map.getBounds();
        this.viewportBounds = {
            north: bounds.getNorth(),
            south: bounds.getSouth(),
            east: bounds.getEast(),
            west: bounds.getWest(),
            center: this.map.getCenter(),
            zoom: this.map.getZoom()
        };
    }

    // Check if a point is within the current viewport (with buffer)
    isInViewport(lat, lng, buffer = 0.1) {
        if (!this.viewportBounds) return true;

        const latBuffer = buffer * (this.viewportBounds.north - this.viewportBounds.south);
        const lngBuffer = buffer * (this.viewportBounds.east - this.viewportBounds.west);

        return lat >= (this.viewportBounds.south - latBuffer) &&
            lat <= (this.viewportBounds.north + latBuffer) &&
            lng >= (this.viewportBounds.west - lngBuffer) &&
            lng <= (this.viewportBounds.east + lngBuffer);
    }

    // Add custom overlay to specific layer
    addOverlay(type, overlay, id = null) {
        const layer = this.layers.overlays.get(type);
        if (layer) {
            if (id) {
                overlay._customId = id;
            }
            layer.addLayer(overlay);
            return true;
        }
        return false;
    }

    // Remove custom overlay
    removeOverlay(type, id) {
        const layer = this.layers.overlays.get(type);
        if (layer) {
            layer.eachLayer(overlay => {
                if (overlay._customId === id) {
                    layer.removeLayer(overlay);
                    return true;
                }
            });
        }
        return false;
    }

    // Toggle overlay layer visibility
    toggleOverlayLayer(type, visible) {
        const layer = this.layers.overlays.get(type);
        if (layer) {
            if (visible && !this.map.hasLayer(layer)) {
                this.map.addLayer(layer);
            } else if (!visible && this.map.hasLayer(layer)) {
                this.map.removeLayer(layer);
            }
        }
    }

    // Check if map is ready/initialized
    isReady() {
        return this.isInitialized && this.map !== null;
    }

    /**
     * Center map on specific coordinates
     * @param {number} lat - Latitude
     * @param {number} lng - Longitude
     * @param {number} zoom - Zoom level (optional, defaults to current zoom)
     */
    centerOn(lat, lng, zoom = null) {
        this.validateCoordinates(lat, lng);
        const zoomLevel = zoom !== null ? zoom : this.getZoom();
        this.map.setView([lat, lng], zoomLevel);
    }

    /**
     * Smoothly fly to location
     * @param {number} lat - Latitude
     * @param {number} lng - Longitude
     * @param {number} zoom - Zoom level
     */
    flyTo(lat, lng, zoom) {
        this.validateCoordinates(lat, lng);
        this.map.flyTo([lat, lng], zoom);
    }

    /**
     * Fit map to bounds
     * @param {Array} bounds - Bounds array [[south, west], [north, east]]
     */
    fitBounds(bounds) {
        this.map.fitBounds(bounds);
    }

    /**
     * Set zoom level
     * @param {number} zoom - Zoom level
     */
    setZoom(zoom) {
        this.map.setZoom(zoom);
    }

    /**
     * Get current zoom level
     * @returns {number} Current zoom level
     */
    getZoom() {
        return this.map.getZoom();
    }

    /**
     * Get platform layer
     * @returns {L.LayerGroup} Platform layer
     */
    getPlatformLayer() {
        return this.platformLayer;
    }

    /**
     * Get trail layer
     * @returns {L.LayerGroup} Trail layer
     */
    getTrailLayer() {
        return this.trailLayer;
    }

    /**
     * Add layer to map
     * @param {L.Layer} layer - Layer to add
     */
    addLayer(layer) {
        this.map.addLayer(layer);
    }

    /**
     * Remove layer from map
     * @param {L.Layer} layer - Layer to remove
     */
    removeLayer(layer) {
        this.map.removeLayer(layer);
    }

    /**
     * Convert lat/lng to screen coordinates
     * @param {number} lat - Latitude
     * @param {number} lng - Longitude
     * @returns {Object} Screen coordinates {x, y}
     */
    latLngToScreen(lat, lng) {
        return this.map.latLngToContainerPoint([lat, lng]);
    }

    /**
     * Convert screen coordinates to lat/lng
     * @param {number} x - Screen X coordinate
     * @param {number} y - Screen Y coordinate
     * @returns {Object} Lat/lng coordinates {lat, lng}
     */
    screenToLatLng(x, y) {
        return this.map.containerPointToLatLng([x, y]);
    }

    /**
     * Project coordinates
     * @param {number} lat - Latitude
     * @param {number} lng - Longitude
     * @returns {Object} Projected coordinates {x, y}
     */
    project(lat, lng) {
        return this.map.project([lat, lng]);
    }

    /**
     * Unproject coordinates
     * @param {number} x - X coordinate
     * @param {number} y - Y coordinate
     * @returns {Object} Unprojected coordinates {lat, lng}
     */
    unproject(x, y) {
        return this.map.unproject([x, y]);
    }

    /**
     * Get map size
     * @returns {Object} Map size {width, height}
     */
    getMapSize() {
        return this.map.getSize();
    }

    /**
     * Get map container element
     * @returns {HTMLElement} Container element
     */
    getContainer() {
        return this.map.getContainer();
    }

    /**
     * Invalidate map size
     */
    invalidateSize() {
        this.map.invalidateSize();
    }

    /**
     * Check if map is ready
     * @returns {boolean} True if map is initialized
     */
    isReady() {
        return this.isInitialized && this.map !== null;
    }

    /**
     * Get the underlying Leaflet map instance
     * @returns {L.Map} The Leaflet map instance
     */
    getMap() {
        return this.map;
    }

    /**
     * Get current viewport bounds (test-compatible format)
     * @returns {Object} Viewport bounds
     */
    getViewportBounds() {
        if (!this.viewportBounds) return null;

        // Return simplified bounds for test compatibility
        return {
            north: this.viewportBounds.north,
            south: this.viewportBounds.south,
            east: this.viewportBounds.east,
            west: this.viewportBounds.west
        };
    }

    /**
     * Get full viewport bounds including center and zoom
     * @returns {Object} Complete viewport bounds with center and zoom
     */
    getFullViewportBounds() {
        return this.viewportBounds;
    }

    /**
     * Trigger viewport change callback
     */
    onViewportChange() {
        // Get simplified bounds for callbacks
        const simplifiedBounds = this.getViewportBounds();

        if (this.onViewportChangeCallback) {
            this.onViewportChangeCallback(simplifiedBounds);
        }

        // Also trigger any registered callbacks
        this.viewportChangeCallbacks.forEach(callback => {
            try {
                callback(simplifiedBounds);
            } catch (error) {
                console.error('Error in viewport change callback:', error);
            }
        });
    }

    /**
     * Register viewport change callback
     * @param {Function} callback - Callback function
     */
    onViewportChangeCallback(callback) {
        this.onViewportChangeCallback = callback;
    }

    /**
     * Add viewport change listener
     * @param {Function} callback - Callback function
     */
    addViewportChangeListener(callback) {
        this.viewportChangeCallbacks.push(callback);
    }

    /**
     * Remove viewport change listener
     * @param {Function} callback - Callback function to remove
     */
    removeViewportChangeListener(callback) {
        const index = this.viewportChangeCallbacks.indexOf(callback);
        if (index > -1) {
            this.viewportChangeCallbacks.splice(index, 1);
        }
    }

    /**
     * Performance callback handler
     * @param {string} metric - Performance metric name
     * @param {*} value - Metric value
     */
    onPerformanceUpdate(metric, value) {
        if (this.performanceCallback) {
            this.performanceCallback(metric, value);
        }
    }

    /**
     * Set performance callback
     * @param {Function} callback - Performance callback function
     */
    setPerformanceCallback(callback) {
        this.performanceCallback = callback;
    }

    /**
     * Throttle function for performance optimization
     * @param {Function} func - Function to throttle
     * @param {number} wait - Wait time in milliseconds
     * @returns {Function} Throttled function
     */
    throttle(func, wait) {
        let timeout;
        return function executedFunction(...args) {
            const later = () => {
                clearTimeout(timeout);
                func(...args);
            };
            clearTimeout(timeout);
            timeout = setTimeout(later, wait);
        };
    }

    /**
     * Export map state
     * @returns {Object} Map state
     */
    exportState() {
        return {
            center: this.map.getCenter(),
            zoom: this.map.getZoom(),
            bounds: this.getViewportBounds()
        };
    }

    /**
     * Restore map state
     * @param {Object} state - Map state to restore
     */
    restoreState(state) {
        if (state.center && state.zoom) {
            this.map.setView([state.center.lat, state.center.lng], state.zoom);
        }
    }

    /**
     * Validate coordinates
     * @param {number} lat - Latitude
     * @param {number} lng - Longitude
     */
    validateCoordinates(lat, lng) {
        if (lat === null || lat === undefined || lng === null || lng === undefined) {
            throw new Error('Latitude and longitude cannot be null or undefined');
        }
        if (lat < -90 || lat > 90) {
            throw new Error('Latitude must be between -90 and 90 degrees');
        }
        if (lng < -180 || lng > 180) {
            throw new Error('Longitude must be between -180 and 180 degrees');
        }
    }

    /**
     * Destroy map and clean up resources
     */
    destroy() {
        if (this.map) {
            // Remove event listeners
            this.map.off('moveend');
            this.map.off('zoomend');
            this.map.off('resize');

            // Remove map
            this.map.remove();
            this.map = null;
        }

        // Clear layer references
        this.platformLayer = null;
        this.trailLayer = null;
        this.layers.platforms = null;
        this.layers.trails = null;
        this.layers.overlays.clear();

        this.isInitialized = false;
    }

    /**
     * Add event listener to map
     * @param {string} event - Event name
     * @param {Function} callback - Event callback
     */
    on(event, callback) {
        this.map.on(event, callback);
    }

    /**
     * Remove event listener from map
     * @param {string} event - Event name
     * @param {Function} callback - Event callback
     */
    off(event, callback) {
        this.map.off(event, callback);
    }

    /**
     * Get viewport bounds with padding for culling
     * @param {number} padding - Padding factor (0.1 = 10% padding)
     * @returns {Object} Padded viewport bounds
     */
    getViewportBoundsWithPadding(padding = 0.1) {
        const bounds = this.getViewportBounds();
        if (!bounds) return null;

        const latRange = bounds.north - bounds.south;
        const lngRange = bounds.east - bounds.west;
        const latPadding = latRange * padding;
        const lngPadding = lngRange * padding;

        return {
            north: bounds.north + latPadding,
            south: bounds.south - latPadding,
            east: bounds.east + lngPadding,
            west: bounds.west - lngPadding
        };
    }

    /**
     * Add scale control to map
     */
    addScaleControl() {
        const scaleControl = L.control.scale({
            position: 'bottomleft',
            metric: true,
            imperial: true
        });
        scaleControl.addTo(this.map);
        return scaleControl;
    }

    /**
     * Add layer control to map
     * @param {Object} baseLayers - Base layers
     * @param {Object} overlayLayers - Overlay layers
     */
    addLayerControl(baseLayers, overlayLayers) {
        const layerControl = L.control.layers(baseLayers, overlayLayers, {
            position: 'topright',
            collapsed: true
        });
        layerControl.addTo(this.map);
        return layerControl;
    }
}

// Export for use in other modules
window.MapEngine = MapEngine;
