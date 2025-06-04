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

        // Don't auto-initialize - wait for explicit call
    }

    async initialize() {
        return new Promise((resolve, reject) => {
            try {
                console.log('Initializing MapEngine...');

                // Initialize Leaflet map with performance optimizations
                this.map = L.map(this.containerId, {
                    center: this.options.center,
                    zoom: this.options.zoom,
                    minZoom: this.options.minZoom,
                    maxZoom: this.options.maxZoom,
                    preferCanvas: this.options.preferCanvas,
                    zoomSnap: this.options.zoomSnap,
                    zoomDelta: this.options.zoomDelta,
                    wheelPxPerZoomLevel: this.options.wheelPxPerZoomLevel,
                    updateWhenIdle: this.options.updateWhenIdle,
                    updateWhenZooming: this.options.updateWhenZooming,
                    keepBuffer: this.options.keepBuffer,
                    // Disable default animations for better performance
                    zoomAnimation: true,
                    fadeAnimation: true,
                    markerZoomAnimation: false
                });

                this.setupOfflineTiles();
                this.setupLayers();
                this.setupEventHandlers();
                this.setupViewportTracking();

                this.isInitialized = true;
                console.log('MapEngine initialized successfully');
                resolve(this);
            } catch (error) {
                console.error('MapEngine initialization failed:', error);
                reject(error);
            }
        });
    }

    setupOfflineTiles() {
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
        // Viewport change events for culling
        this.map.on('moveend zoomend', () => {
            this.updateViewportBounds();
            this.onViewportChange();
        });

        // Performance monitoring
        this.map.on('movestart', () => {
            this.lastMoveTime = Date.now();
        });

        this.map.on('moveend', () => {
            const moveTime = Date.now() - this.lastMoveTime;
            this.onPerformanceUpdate('moveTime', moveTime);
        });

        // Layer visibility controls
        this.setupLayerControls();
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

    // Performance optimization utilities
    throttle(func, limit) {
        let inThrottle;
        return function () {
            const args = arguments;
            const context = this;
            if (!inThrottle) {
                func.apply(context, args);
                inThrottle = true;
                setTimeout(() => inThrottle = false, limit);
            }
        };
    }

    debounce(func, delay) {
        let timeoutId;
        return function (...args) {
            clearTimeout(timeoutId);
            timeoutId = setTimeout(() => func.apply(this, args), delay);
        };
    }

    // Event handlers for external components
    onViewportChange() {
        // Override in external components
        if (this.onViewportChangeCallback) {
            this.onViewportChangeCallback(this.viewportBounds);
        }
    }

    onPerformanceUpdate(metric, value) {
        // Override in external components
        if (this.onPerformanceUpdateCallback) {
            this.onPerformanceUpdateCallback(metric, value);
        }
    }

    // Public API methods
    getMap() {
        return this.map;
    }

    getViewportBounds() {
        return this.viewportBounds;
    }

    getPlatformLayer() {
        return this.layers.platforms;
    }

    getTrailLayer() {
        return this.layers.trails;
    }

    getOverlayLayer(type) {
        return this.layers.overlays.get(type);
    }

    // Center map on coordinates
    centerOn(lat, lng, zoom = null) {
        if (zoom !== null) {
            this.map.setView([lat, lng], zoom);
        } else {
            this.map.panTo([lat, lng]);
        }
    }

    // Fit map to bounds
    fitBounds(bounds, options = {}) {
        this.map.fitBounds(bounds, {
            padding: [20, 20],
            maxZoom: 15,
            ...options
        });
    }

    // Get current map state for persistence
    getMapState() {
        return {
            center: this.map.getCenter(),
            zoom: this.map.getZoom(),
            bounds: this.viewportBounds
        };
    }

    // Restore map state
    setMapState(state) {
        if (state.center && state.zoom) {
            this.map.setView([state.center.lat, state.center.lng], state.zoom);
        }
    }

    // Cleanup method
    destroy() {
        if (this.map) {
            this.map.remove();
            this.map = null;
        }
        this.isInitialized = false;
    }
}

// Export for use in other modules
window.MapEngine = MapEngine;
