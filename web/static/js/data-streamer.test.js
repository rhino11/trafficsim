/**
 * Comprehensive unit tests for DataStreamer class
 */

// Import the DataStreamer class
require('./data-streamer.js');

describe('DataStreamer', () => {
    let dataStreamer;
    let mockWebSocket;
    let mockEventSource;
    let originalWebSocket;
    let originalEventSource;
    let originalFetch;

    beforeEach(() => {
        // Store original constructors
        originalWebSocket = global.WebSocket;
        originalEventSource = global.EventSource;
        originalFetch = global.fetch;

        // Mock WebSocket
        mockWebSocket = {
            readyState: 1, // OPEN
            send: jest.fn(),
            close: jest.fn(),
            addEventListener: jest.fn(),
            removeEventListener: jest.fn(),
            onopen: null,
            onmessage: null,
            onclose: null,
            onerror: null
        };

        // Mock EventSource
        mockEventSource = {
            readyState: 1, // OPEN
            close: jest.fn(),
            addEventListener: jest.fn(),
            removeEventListener: jest.fn(),
            onopen: null,
            onmessage: null,
            onerror: null
        };

        // Mock fetch
        global.fetch = jest.fn();

        // Mock WebSocket constructor
        global.WebSocket = jest.fn(() => mockWebSocket);
        global.WebSocket.OPEN = 1;
        global.WebSocket.CLOSED = 3;

        // Mock EventSource constructor
        global.EventSource = jest.fn(() => mockEventSource);

        // Mock performance
        global.performance = {
            now: jest.fn(() => Date.now())
        };

        // Mock console methods
        console.log = jest.fn();
        console.error = jest.fn();
        console.warn = jest.fn();
    });

    afterEach(() => {
        if (dataStreamer) {
            dataStreamer.disconnect();
        }

        // Restore original constructors
        global.WebSocket = originalWebSocket;
        global.EventSource = originalEventSource;
        global.fetch = originalFetch;

        // Clear timers
        jest.clearAllTimers();
        jest.useRealTimers();
    });

    describe('Constructor and Initialization', () => {
        it('should initialize with default options', () => {
            dataStreamer = new DataStreamer();

            expect(dataStreamer.options.reconnectInterval).toBe(5000);
            expect(dataStreamer.options.maxReconnectAttempts).toBe(10);
            expect(dataStreamer.options.batchSize).toBe(100);
            expect(dataStreamer.options.updateThrottle).toBe(16);
            expect(dataStreamer.connectionType).toBe('websocket');
            expect(dataStreamer.isConnected).toBe(false);
        });

        it('should merge custom options with defaults', () => {
            const customOptions = {
                reconnectInterval: 3000,
                batchSize: 50,
                customOption: 'test'
            };

            dataStreamer = new DataStreamer(customOptions);

            expect(dataStreamer.options.reconnectInterval).toBe(3000);
            expect(dataStreamer.options.batchSize).toBe(50);
            expect(dataStreamer.options.maxReconnectAttempts).toBe(10); // default
            expect(dataStreamer.options.customOption).toBe('test');
        });

        it('should detect best connection type', () => {
            dataStreamer = new DataStreamer();
            expect(dataStreamer.connectionType).toBe('websocket');

            // Test fallback to SSE when WebSocket is not available
            delete global.WebSocket;
            dataStreamer = new DataStreamer();
            expect(dataStreamer.connectionType).toBe('sse');
        });

        it('should initialize performance tracking', () => {
            jest.useFakeTimers();
            dataStreamer = new DataStreamer();

            expect(dataStreamer.stats).toEqual({
                messagesReceived: 0,
                bytesReceived: 0,
                updateRate: 0,
                lastUpdateTime: 0,
                averageLatency: 0,
                compressionRatio: 1.0
            });
        });
    });

    describe('WebSocket Connection', () => {
        beforeEach(() => {
            dataStreamer = new DataStreamer();
        });

        it('should create WebSocket connection', () => {
            dataStreamer.connect();

            expect(global.WebSocket).toHaveBeenCalledWith('ws://localhost/ws');
            expect(mockWebSocket.onopen).toBeDefined();
            expect(mockWebSocket.onmessage).toBeDefined();
            expect(mockWebSocket.onclose).toBeDefined();
            expect(mockWebSocket.onerror).toBeDefined();
        });

        it('should handle WebSocket connection open', () => {
            const statusCallback = jest.fn();
            dataStreamer.onConnectionStatus(statusCallback);

            dataStreamer.connect();
            mockWebSocket.onopen();

            expect(dataStreamer.isConnected).toBe(true);
            expect(dataStreamer.reconnectAttempts).toBe(0);
            expect(statusCallback).toHaveBeenCalledWith('connected');
        });

        it('should handle WebSocket messages', () => {
            const platformCallback = jest.fn();
            dataStreamer.onPlatformUpdate(platformCallback);

            dataStreamer.connect();

            const testMessage = {
                type: 'platform_update',
                platforms: [
                    { id: 'test-1', platform_type: 'airborne', position: { latitude: 40, longitude: -75 } }
                ]
            };

            mockWebSocket.onmessage({ data: JSON.stringify(testMessage) });

            // Process should be throttled, so we need to wait
            jest.useFakeTimers();
            jest.advanceTimersByTime(20);

            expect(platformCallback).toHaveBeenCalledWith([testMessage.platforms[0]]);
        });

        it('should handle WebSocket connection close and reconnect', () => {
            jest.useFakeTimers();
            const statusCallback = jest.fn();
            dataStreamer.onConnectionStatus(statusCallback);

            dataStreamer.connect();
            mockWebSocket.onopen();
            expect(dataStreamer.isConnected).toBe(true);

            // Simulate connection close
            mockWebSocket.onclose({ code: 1000, reason: 'Normal closure' });

            expect(dataStreamer.isConnected).toBe(false);
            expect(statusCallback).toHaveBeenCalledWith('disconnected');
            expect(statusCallback).toHaveBeenCalledWith('connecting');

            // Should schedule reconnection
            expect(dataStreamer.reconnectAttempts).toBe(1);

            // Advance timer to trigger reconnect
            jest.advanceTimersByTime(5000);
            expect(global.WebSocket).toHaveBeenCalledTimes(2);
        });

        it('should handle WebSocket errors and fallback to SSE', () => {
            const statusCallback = jest.fn();
            dataStreamer.onConnectionStatus(statusCallback);

            dataStreamer.connect();
            mockWebSocket.onerror({ error: 'Connection failed' });

            expect(dataStreamer.isConnected).toBe(false);
            expect(statusCallback).toHaveBeenCalledWith('error');
        });

        it('should send heartbeat messages', () => {
            jest.useFakeTimers();
            dataStreamer.connect();
            mockWebSocket.onopen();

            // Advance timer to trigger heartbeat
            jest.advanceTimersByTime(30000);

            expect(mockWebSocket.send).toHaveBeenCalledWith(
                expect.stringContaining('"type":"ping"')
            );
        });

        it('should handle pong responses and update latency', () => {
            const timestamp = Date.now() - 100; // 100ms ago
            dataStreamer.connect();

            const pongMessage = {
                type: 'pong',
                timestamp: timestamp
            };

            mockWebSocket.onmessage({ data: JSON.stringify(pongMessage) });

            expect(dataStreamer.stats.averageLatency).toBeGreaterThan(0);
        });
    });

    describe('Server-Sent Events (SSE) Connection', () => {
        beforeEach(() => {
            delete global.WebSocket; // Force SSE mode
            dataStreamer = new DataStreamer();
        });

        it('should create SSE connection when WebSocket is not available', () => {
            dataStreamer.connect();

            expect(global.EventSource).toHaveBeenCalledWith('/api/stream/platforms');
            expect(mockEventSource.onopen).toBeDefined();
            expect(mockEventSource.onmessage).toBeDefined();
            expect(mockEventSource.onerror).toBeDefined();
        });

        it('should handle SSE connection open', () => {
            const statusCallback = jest.fn();
            dataStreamer.onConnectionStatus(statusCallback);

            dataStreamer.connect();
            mockEventSource.onopen();

            expect(dataStreamer.isConnected).toBe(true);
            expect(statusCallback).toHaveBeenCalledWith('connected');
        });

        it('should handle SSE messages', () => {
            const platformCallback = jest.fn();
            dataStreamer.onPlatformUpdate(platformCallback);

            dataStreamer.connect();

            const testData = JSON.stringify([
                { id: 'test-1', platform_type: 'maritime', position: { latitude: 41, longitude: -76 } }
            ]);

            mockEventSource.onmessage({ data: testData });

            jest.useFakeTimers();
            jest.advanceTimersByTime(20);

            expect(platformCallback).toHaveBeenCalled();
        });

        it('should fetch initial data for SSE connections', async () => {
            const mockPlatforms = [
                { id: 'initial-1', platform_type: 'land', position: { latitude: 42, longitude: -77 } }
            ];

            global.fetch.mockResolvedValue({
                ok: true,
                json: () => Promise.resolve(mockPlatforms)
            });

            const platformCallback = jest.fn();
            dataStreamer.onPlatformUpdate(platformCallback);

            dataStreamer.connect();
            mockEventSource.onopen();

            // Wait for async operations
            await new Promise(resolve => setTimeout(resolve, 0));

            expect(global.fetch).toHaveBeenCalledWith('/api/platforms');
        });
    });

    describe('Message Processing', () => {
        beforeEach(() => {
            dataStreamer = new DataStreamer();
        });

        it('should process platform update messages', () => {
            const platformCallback = jest.fn();
            dataStreamer.onPlatformUpdate(platformCallback);

            const platforms = [
                { id: 'test-1', platform_type: 'airborne', position: { latitude: 40, longitude: -75 } },
                { id: 'test-2', platform_type: 'maritime', position: { latitude: 41, longitude: -76 } }
            ];

            const message = {
                type: 'platform_update',
                platforms: platforms
            };

            dataStreamer.handleMessage(message);

            jest.useFakeTimers();
            jest.advanceTimersByTime(20);

            expect(platformCallback).toHaveBeenCalledWith(platforms);
        });

        it('should process simulation status messages', () => {
            const statusCallback = jest.fn();
            dataStreamer.onSimulationStatus(statusCallback);

            const status = {
                running: true,
                time: 1234.5,
                platformCount: 100,
                speed: 2.0
            };

            const message = {
                type: 'simulation_status',
                data: status
            };

            dataStreamer.handleMessage(message);

            expect(statusCallback).toHaveBeenCalledWith(status);
        });

        it('should handle delta compression', () => {
            dataStreamer.options.enableDeltaCompression = true;
            const platformCallback = jest.fn();
            dataStreamer.onPlatformUpdate(platformCallback);

            // First, send full platform state
            const fullPlatform = {
                id: 'test-1',
                platform_type: 'airborne',
                position: { latitude: 40, longitude: -75, altitude: 1000 },
                velocity: { north: 10, east: 5, up: 1 },
                speed: 100
            };

            dataStreamer.handleMessage({
                type: 'platform_update',
                platforms: [fullPlatform]
            });

            jest.useFakeTimers();
            jest.advanceTimersByTime(20);

            // Now send delta update
            const deltaUpdate = {
                id: 'test-1',
                delta: {
                    position: { latitude: 40.001 },
                    speed: 105
                }
            };

            platformCallback.mockClear();
            dataStreamer.handleMessage({
                type: 'platform_update',
                platforms: [deltaUpdate]
            });

            jest.advanceTimersByTime(20);

            const updatedPlatform = platformCallback.mock.calls[0][0][0];
            expect(updatedPlatform.position.latitude).toBe(40.001);
            expect(updatedPlatform.position.longitude).toBe(-75); // Should retain from full state
            expect(updatedPlatform.speed).toBe(105);
        });
    });

    describe('Simulation Control', () => {
        beforeEach(() => {
            dataStreamer = new DataStreamer();
        });

        it('should send start simulation command via WebSocket', async () => {
            dataStreamer.connect();
            mockWebSocket.onopen();

            await dataStreamer.startSimulation();

            expect(mockWebSocket.send).toHaveBeenCalledWith(
                expect.stringContaining('"action":"start_simulation"')
            );
        });

        it('should send stop simulation command via WebSocket', async () => {
            dataStreamer.connect();
            mockWebSocket.onopen();

            await dataStreamer.stopSimulation();

            expect(mockWebSocket.send).toHaveBeenCalledWith(
                expect.stringContaining('"action":"stop_simulation"')
            );
        });

        it('should fallback to HTTP for control commands when WebSocket is not available', async () => {
            global.fetch.mockResolvedValue({
                ok: true,
                json: () => Promise.resolve({ status: 'success' })
            });

            await dataStreamer.startSimulation();

            expect(global.fetch).toHaveBeenCalledWith('/api/simulation/start_simulation', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({})
            });
        });

        it('should handle HTTP errors for control commands', async () => {
            global.fetch.mockResolvedValue({
                ok: false,
                status: 500,
                statusText: 'Internal Server Error'
            });

            await expect(dataStreamer.startSimulation()).rejects.toThrow('HTTP 500: Internal Server Error');
        });
    });

    describe('Performance and Statistics', () => {
        beforeEach(() => {
            dataStreamer = new DataStreamer();
        });

        it('should track message statistics', () => {
            const testData = JSON.stringify({ type: 'test', data: 'hello' });

            dataStreamer.onDataReceived(testData);

            expect(dataStreamer.stats.messagesReceived).toBe(1);
            expect(dataStreamer.stats.bytesReceived).toBe(testData.length);
        });

        it('should calculate update rate', () => {
            jest.useFakeTimers();

            // Simulate multiple updates
            for (let i = 0; i < 5; i++) {
                dataStreamer.updateRateTracker.push(Date.now());
                jest.advanceTimersByTime(100);
            }

            dataStreamer.updateStats();

            expect(dataStreamer.stats.updateRate).toBeGreaterThan(0);
        });

        it('should calculate compression ratio', () => {
            dataStreamer.stats.messagesReceived = 10;
            dataStreamer.stats.bytesReceived = 1000;

            dataStreamer.updateStats();

            expect(dataStreamer.stats.compressionRatio).toBeGreaterThan(1);
        });

        it('should provide comprehensive statistics', () => {
            const stats = dataStreamer.getStats();

            expect(stats).toHaveProperty('messagesReceived');
            expect(stats).toHaveProperty('bytesReceived');
            expect(stats).toHaveProperty('updateRate');
            expect(stats).toHaveProperty('averageLatency');
            expect(stats).toHaveProperty('isConnected');
            expect(stats).toHaveProperty('connectionType');
            expect(stats).toHaveProperty('queueSize');
            expect(stats).toHaveProperty('platformStates');
        });

        it('should clean up old rate tracking data', () => {
            const now = Date.now();

            // Add old and new entries
            dataStreamer.updateRateTracker = [
                now - 2000, // Old entry (should be removed)
                now - 500,  // Recent entry (should be kept)
                now         // Current entry (should be kept)
            ];

            dataStreamer.cleanupRateTracker();

            expect(dataStreamer.updateRateTracker).toHaveLength(2);
            expect(dataStreamer.updateRateTracker[0]).toBe(now - 500);
        });
    });

    describe('Advanced Features', () => {
        beforeEach(() => {
            dataStreamer = new DataStreamer();
        });

        it('should send viewport updates for server-side filtering', () => {
            dataStreamer.connect();
            mockWebSocket.onopen();

            const bounds = {
                north: 45,
                south: 35,
                east: -70,
                west: -80
            };

            dataStreamer.updateViewport(bounds);

            expect(mockWebSocket.send).toHaveBeenCalledWith(
                expect.stringContaining('"type":"viewport_update"')
            );
        });

        it('should send filter updates for server-side filtering', () => {
            dataStreamer.connect();
            mockWebSocket.onopen();

            const filters = {
                airborne: true,
                maritime: false,
                land: true,
                space: true
            };

            dataStreamer.updateFilters(filters);

            expect(mockWebSocket.send).toHaveBeenCalledWith(
                expect.stringContaining('"type":"filter_update"')
            );
        });

        it('should handle connection timeout during reconnection', () => {
            jest.useFakeTimers();
            dataStreamer.options.maxReconnectAttempts = 2;

            const statusCallback = jest.fn();
            dataStreamer.onConnectionStatus(statusCallback);

            // Simulate multiple failed reconnection attempts
            for (let i = 0; i < 3; i++) {
                dataStreamer.connect();
                mockWebSocket.onerror({ error: 'Connection failed' });
            }

            expect(statusCallback).toHaveBeenCalledWith('failed');
        });

        it('should handle malformed JSON messages gracefully', () => {
            const consoleSpy = jest.spyOn(console, 'error').mockImplementation();

            dataStreamer.onDataReceived('invalid json data');

            expect(consoleSpy).toHaveBeenCalledWith(
                'Error processing received data:',
                expect.any(Error),
                'invalid json data'
            );

            consoleSpy.mockRestore();
        });
    });

    describe('Cleanup and Disconnection', () => {
        beforeEach(() => {
            dataStreamer = new DataStreamer();
        });

        it('should clean up resources on disconnect', () => {
            jest.useFakeTimers();
            dataStreamer.connect();
            mockWebSocket.onopen();

            // Start heartbeat
            expect(dataStreamer.heartbeatTimer).toBeDefined();

            dataStreamer.disconnect();

            expect(dataStreamer.isConnected).toBe(false);
            expect(mockWebSocket.close).toHaveBeenCalled();
            expect(dataStreamer.heartbeatTimer).toBeNull();
        });

        it('should close EventSource on disconnect', () => {
            delete global.WebSocket;
            dataStreamer = new DataStreamer();
            dataStreamer.connect();

            dataStreamer.disconnect();

            expect(mockEventSource.close).toHaveBeenCalled();
        });

        it('should clear pending timers on disconnect', () => {
            jest.useFakeTimers();

            // Queue some updates to create timers
            dataStreamer.queuePlatformUpdates([{ id: 'test' }]);
            expect(dataStreamer.updateThrottleTimer).toBeDefined();

            dataStreamer.disconnect();

            expect(dataStreamer.updateThrottleTimer).toBeNull();
        });
    });

    describe('Callback Management', () => {
        beforeEach(() => {
            dataStreamer = new DataStreamer();
        });

        it('should register and call platform update callbacks', () => {
            const callback1 = jest.fn();
            const callback2 = jest.fn();

            dataStreamer.onPlatformUpdate(callback1);

            // Should replace previous callback
            dataStreamer.onPlatformUpdate(callback2);

            const platforms = [{ id: 'test', platform_type: 'airborne' }];
            dataStreamer.handleMessage({
                type: 'platform_update',
                platforms: platforms
            });

            jest.useFakeTimers();
            jest.advanceTimersByTime(20);

            expect(callback1).not.toHaveBeenCalled();
            expect(callback2).toHaveBeenCalledWith(platforms);
        });

        it('should handle missing callbacks gracefully', () => {
            // No callbacks registered
            const platforms = [{ id: 'test', platform_type: 'airborne' }];

            expect(() => {
                dataStreamer.handleMessage({
                    type: 'platform_update',
                    platforms: platforms
                });
            }).not.toThrow();
        });
    });
});
