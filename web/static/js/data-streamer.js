/**
 * High-Performance Data Streamer
 * Handles real-time platform updates via WebSocket/SSE
 */
class DataStreamer {
    constructor(options = {}) {
        this.options = {
            // Connection settings
            reconnectInterval: 5000,
            maxReconnectAttempts: 10,
            heartbeatInterval: 30000,
            batchSize: 100,
            updateThrottle: 16, // ~60 FPS
            // Data compression
            enableCompression: true,
            enableDeltaCompression: true,
            ...options
        };

        // Connection management
        this.websocket = null;
        this.eventSource = null;
        this.connectionType = 'websocket'; // 'websocket' or 'sse'
        this.isConnected = false;
        this.reconnectAttempts = 0;
        this.heartbeatTimer = null;

        // Data management
        this.updateQueue = [];
        this.lastPlatformStates = new Map(); // For delta compression
        this.pendingUpdates = new Map();
        this.updateThrottleTimer = null;

        // Callback handlers
        this.onPlatformUpdateCallback = null;
        this.onSimulationStatusCallback = null;
        this.onConnectionStatusCallback = null;
        this.onStatsUpdateCallback = null;

        // Performance tracking
        this.stats = {
            messagesReceived: 0,
            bytesReceived: 0,
            updateRate: 0,
            lastUpdateTime: 0,
            averageLatency: 0,
            compressionRatio: 1.0
        };

        // Rate limiting
        this.updateRateTracker = [];
        this.lastStatsUpdate = Date.now();

        this.init();
    }

    init() {
        console.log('Initializing DataStreamer...');
        this.detectBestConnectionType();
        this.setupPerformanceTracking();
        console.log('DataStreamer initialized');
    }

    // Set platform update callback
    onPlatformUpdate(callback) {
        this.onPlatformUpdateCallback = callback;
    }

    // Set simulation status callback
    onSimulationStatus(callback) {
        this.onSimulationStatusCallback = callback;
    }

    // Set connection status callback
    onConnectionStatus(callback) {
        this.onConnectionStatusCallback = callback;
    }

    detectBestConnectionType() {
        // Prefer WebSocket for bidirectional communication
        if (typeof WebSocket !== 'undefined') {
            this.connectionType = 'websocket';
        } else if (typeof EventSource !== 'undefined') {
            this.connectionType = 'sse';
        } else {
            console.error('No supported real-time connection type available');
            return false;
        }

        console.log(`Using connection type: ${this.connectionType}`);
        return true;
    }

    connect() {
        if (this.isConnected) {
            console.log('Already connected');
            return;
        }

        console.log(`Connecting via ${this.connectionType}...`);

        if (this.connectionType === 'websocket') {
            this.connectWebSocket();
        } else if (this.connectionType === 'sse') {
            this.connectSSE();
        }
    }

    connectWebSocket() {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/ws`;

        try {
            this.websocket = new WebSocket(wsUrl);

            this.websocket.onopen = (event) => {
                console.log('WebSocket connected');
                this.onConnectionOpen(event);
            };

            this.websocket.onmessage = (event) => {
                this.onDataReceived(event.data);
            };

            this.websocket.onclose = (event) => {
                console.log('WebSocket closed:', event.code, event.reason);
                this.onConnectionClosed(event);
            };

            this.websocket.onerror = (event) => {
                console.error('WebSocket error:', event);
                this.onConnectionError(event);
            };

        } catch (error) {
            console.error('Failed to create WebSocket:', error);
            this.fallbackToSSE();
        }
    }

    connectSSE() {
        const sseUrl = `/api/stream/platforms`;

        try {
            this.eventSource = new EventSource(sseUrl);

            this.eventSource.onopen = (event) => {
                console.log('SSE connected');
                this.onConnectionOpen(event);
            };

            this.eventSource.onmessage = (event) => {
                this.onDataReceived(event.data);
            };

            this.eventSource.onerror = (event) => {
                console.error('SSE error:', event);
                this.onConnectionError(event);
            };

            // Handle specific event types
            this.eventSource.addEventListener('platform_update', (event) => {
                this.onDataReceived(event.data);
            });

            this.eventSource.addEventListener('simulation_status', (event) => {
                this.onSimulationStatusUpdate(JSON.parse(event.data));
            });

        } catch (error) {
            console.error('Failed to create EventSource:', error);
            this.onConnectionError(error);
        }
    }

    fallbackToSSE() {
        if (this.connectionType === 'websocket') {
            console.log('Falling back to SSE');
            this.connectionType = 'sse';
            setTimeout(() => this.connectSSE(), 1000);
        }
    }

    onConnectionOpen(event) {
        this.isConnected = true;
        this.reconnectAttempts = 0;
        this.updateConnectionStatus('connected');

        // Start heartbeat for WebSocket
        if (this.connectionType === 'websocket') {
            this.startHeartbeat();
        }

        // Request initial platform data
        this.requestInitialData();
    }

    onConnectionClosed(event) {
        this.isConnected = false;
        this.updateConnectionStatus('disconnected');
        this.stopHeartbeat();

        // Attempt reconnection
        if (this.reconnectAttempts < this.options.maxReconnectAttempts) {
            this.scheduleReconnect();
        } else {
            console.error('Max reconnection attempts reached');
            this.updateConnectionStatus('failed');
        }
    }

    onConnectionError(event) {
        console.error('Connection error:', event);
        this.isConnected = false;
        this.updateConnectionStatus('error');
    }

    scheduleReconnect() {
        this.reconnectAttempts++;
        const delay = Math.min(this.options.reconnectInterval * this.reconnectAttempts, 30000);

        console.log(`Reconnecting in ${delay}ms (attempt ${this.reconnectAttempts})`);
        this.updateConnectionStatus('connecting');

        setTimeout(() => {
            this.connect();
        }, delay);
    }

    startHeartbeat() {
        this.heartbeatTimer = setInterval(() => {
            if (this.websocket && this.websocket.readyState === WebSocket.OPEN) {
                this.websocket.send(JSON.stringify({ type: 'ping', timestamp: Date.now() }));
            }
        }, this.options.heartbeatInterval);
    }

    stopHeartbeat() {
        if (this.heartbeatTimer) {
            clearInterval(this.heartbeatTimer);
            this.heartbeatTimer = null;
        }
    }

    requestInitialData() {
        if (this.connectionType === 'websocket' && this.websocket) {
            this.websocket.send(JSON.stringify({
                type: 'request_initial_data',
                timestamp: Date.now()
            }));
        } else {
            // For SSE, make an HTTP request for initial data
            this.fetchInitialData();
        }
    }

    async fetchInitialData() {
        try {
            const response = await fetch('/api/platforms');
            if (response.ok) {
                const platforms = await response.json();
                this.processPlatformUpdate({ platforms: platforms });
            }
        } catch (error) {
            console.error('Failed to fetch initial data:', error);
        }
    }

    onDataReceived(data) {
        const startTime = performance.now();

        try {
            // Track statistics
            this.stats.messagesReceived++;
            this.stats.bytesReceived += data.length;

            // Parse message
            let message;
            if (typeof data === 'string') {
                message = JSON.parse(data);
            } else {
                message = data;
            }

            // Handle different message types
            this.handleMessage(message);

            // Track processing time
            const processingTime = performance.now() - startTime;
            this.updateLatencyStats(processingTime);

        } catch (error) {
            console.error('Error processing received data:', error, data);
        }
    }

    handleMessage(message) {
        switch (message.type) {
            case 'platform_update':
            case 'platforms':
                this.processPlatformUpdate(message);
                break;

            case 'platform_batch':
                this.processPlatformBatch(message);
                break;

            case 'simulation_status':
                this.onSimulationStatusUpdate(message.data);
                break;

            case 'pong':
                this.handlePong(message);
                break;

            case 'error':
                console.error('Server error:', message.error);
                break;

            default:
                console.warn('Unknown message type:', message.type);
        }
    }

    processPlatformUpdate(message) {
        let platforms = message.platforms || message.data || [];

        // Handle single platform update
        if (message.platform) {
            platforms = [message.platform];
        }

        // Apply delta compression if enabled
        if (this.options.enableDeltaCompression) {
            platforms = this.applyDeltaDecompression(platforms);
        }

        // Queue updates for throttled processing
        this.queuePlatformUpdates(platforms);
    }

    processPlatformBatch(message) {
        const batch = message.batch || [];
        this.queuePlatformUpdates(batch);
    }

    queuePlatformUpdates(platforms) {
        // Add to update queue
        this.updateQueue.push(...platforms);

        // Process queue with throttling
        if (!this.updateThrottleTimer) {
            this.updateThrottleTimer = setTimeout(() => {
                this.processUpdateQueue();
                this.updateThrottleTimer = null;
            }, this.options.updateThrottle);
        }
    }

    processUpdateQueue() {
        if (this.updateQueue.length === 0) return;

        const startTime = performance.now();

        // Process updates in batches for better performance
        const batchSize = Math.min(this.options.batchSize, this.updateQueue.length);
        const batch = this.updateQueue.splice(0, batchSize);

        // Send updates to callback if available
        if (this.onPlatformUpdateCallback) {
            this.onPlatformUpdateCallback(batch);
        }

        // Track update rate
        this.updateRateTracker.push(Date.now());
        this.cleanupRateTracker();

        // Update statistics
        this.stats.lastUpdateTime = Date.now();
        this.updateStats();

        // If there are more updates, schedule next batch
        if (this.updateQueue.length > 0) {
            setTimeout(() => this.processUpdateQueue(), this.options.updateThrottle);
        }

        // Track performance
        const processingTime = performance.now() - startTime;
        if (this.onPerformanceUpdateCallback) {
            this.onPerformanceUpdateCallback('batchProcessTime', processingTime);
        }
    }

    applyDeltaDecompression(platforms) {
        // Decompress delta updates by merging with last known state
        return platforms.map(platform => {
            const lastState = this.lastPlatformStates.get(platform.id);

            if (lastState && platform.delta) {
                // Merge delta with last state
                const fullPlatform = { ...lastState };

                // Apply delta changes
                Object.keys(platform.delta).forEach(key => {
                    if (key === 'position' && lastState.position) {
                        fullPlatform.position = { ...lastState.position, ...platform.delta.position };
                    } else if (key === 'velocity' && lastState.velocity) {
                        fullPlatform.velocity = { ...lastState.velocity, ...platform.delta.velocity };
                    } else {
                        fullPlatform[key] = platform.delta[key];
                    }
                });

                // Store updated state
                this.lastPlatformStates.set(platform.id, fullPlatform);
                return fullPlatform;
            } else {
                // Full platform update
                this.lastPlatformStates.set(platform.id, platform);
                return platform;
            }
        });
    }

    handlePong(message) {
        const latency = Date.now() - message.timestamp;
        this.updateLatencyStats(latency);
    }

    updateLatencyStats(latency) {
        // Simple moving average for latency
        this.stats.averageLatency = (this.stats.averageLatency * 0.9) + (latency * 0.1);
    }

    cleanupRateTracker() {
        const now = Date.now();
        const cutoff = now - 1000; // Keep only last second
        this.updateRateTracker = this.updateRateTracker.filter(time => time > cutoff);
    }

    updateStats() {
        const now = Date.now();
        if (now - this.lastStatsUpdate > 1000) {
            // Calculate update rate
            this.stats.updateRate = this.updateRateTracker.length;

            // Calculate compression ratio
            if (this.stats.bytesReceived > 0) {
                const estimatedUncompressedSize = this.stats.messagesReceived * 500; // Estimate
                this.stats.compressionRatio = estimatedUncompressedSize / this.stats.bytesReceived;
            }

            this.lastStatsUpdate = now;

            // Notify performance monitor
            if (this.onStatsUpdateCallback) {
                this.onStatsUpdateCallback(this.getStats());
            }
        }
    }

    // Simulation control methods
    async startSimulation() {
        return this.sendControlMessage('start_simulation');
    }

    async stopSimulation() {
        return this.sendControlMessage('stop_simulation');
    }

    async resetSimulation() {
        return this.sendControlMessage('reset_simulation');
    }

    sendControlMessage(action, data = {}) {
        const message = {
            type: 'control',
            action: action,
            data: data,
            timestamp: Date.now()
        };

        if (this.connectionType === 'websocket' && this.websocket && this.websocket.readyState === WebSocket.OPEN) {
            this.websocket.send(JSON.stringify(message));
            return Promise.resolve();
        } else {
            // Fallback to HTTP for control messages
            return this.sendHTTPControlMessage(action, data);
        }
    }

    async sendHTTPControlMessage(action, data) {
        try {
            const response = await fetch(`/api/simulation/${action}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(data)
            });

            if (!response.ok) {
                throw new Error(`HTTP ${response.status}: ${response.statusText}`);
            }

            return await response.json();
        } catch (error) {
            console.error(`Failed to send ${action} command:`, error);
            throw error;
        }
    }

    // Event handlers
    onSimulationStatusUpdate(status) {
        if (this.onSimulationStatusCallback) {
            this.onSimulationStatusCallback(status);
        }
    }

    updateConnectionStatus(status) {
        if (this.onConnectionStatusCallback) {
            this.onConnectionStatusCallback(status);
        }
    }

    // Public API methods
    getStats() {
        return {
            ...this.stats,
            isConnected: this.isConnected,
            connectionType: this.connectionType,
            queueSize: this.updateQueue.length,
            platformStates: this.lastPlatformStates.size
        };
    }

    // Update viewport for server-side filtering
    updateViewport(bounds) {
        if (this.connectionType === 'websocket' && this.websocket && this.websocket.readyState === WebSocket.OPEN) {
            this.websocket.send(JSON.stringify({
                type: 'viewport_update',
                viewport: bounds,
                timestamp: Date.now()
            }));
        }
    }

    // Update filters for server-side filtering
    updateFilters(filters) {
        if (this.connectionType === 'websocket' && this.websocket && this.websocket.readyState === WebSocket.OPEN) {
            this.websocket.send(JSON.stringify({
                type: 'filter_update',
                filters: filters,
                timestamp: Date.now()
            }));
        }
    }

    // Cleanup
    disconnect() {
        this.isConnected = false;
        this.stopHeartbeat();

        if (this.websocket) {
            this.websocket.close();
            this.websocket = null;
        }

        if (this.eventSource) {
            this.eventSource.close();
            this.eventSource = null;
        }

        if (this.updateThrottleTimer) {
            clearTimeout(this.updateThrottleTimer);
            this.updateThrottleTimer = null;
        }

        console.log('DataStreamer disconnected');
    }

    setupPerformanceTracking() {
        // Set up performance monitoring
        setInterval(() => {
            this.updateStats();
        }, 1000);
    }
}

// Export for use in other modules
window.DataStreamer = DataStreamer;
