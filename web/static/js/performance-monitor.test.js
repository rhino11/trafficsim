/**
 * Tests for PerformanceMonitor
 * Ensures basic functionality and key methods work correctly
 */

// Import the PerformanceMonitor class
require('./performance-monitor.js');

// Mock performance API if not available
global.performance = global.performance || {
    now: () => Date.now(),
    memory: {
        usedJSHeapSize: 50 * 1024 * 1024, // 50MB
        totalJSHeapSize: 100 * 1024 * 1024, // 100MB
        jsHeapSizeLimit: 200 * 1024 * 1024 // 200MB
    }
};

// Mock requestAnimationFrame
global.requestAnimationFrame = global.requestAnimationFrame || ((cb) => setTimeout(cb, 16));

describe('PerformanceMonitor', () => {
    let performanceMonitor;
    let mockElements;
    let originalAppendChild;
    let originalRemoveChild;
    let originalCreateElement;

    beforeEach(() => {
        // Create mock DOM elements
        mockElements = {
            fps: { textContent: '', style: {} },
            platformCount: { textContent: '', style: {} },
            visibleCount: { textContent: '', style: {} },
            memoryUsage: { textContent: '', style: {} },
            renderTime: { textContent: '', style: {} },
            updateRate: { textContent: '', style: {} },
            dataRate: { textContent: '', style: {} }
        };

        // Mock document.getElementById
        document.getElementById = jest.fn((id) => mockElements[id] || null);

        // Store original methods
        originalAppendChild = document.body.appendChild;
        originalRemoveChild = document.body.removeChild;
        originalCreateElement = document.createElement;

        // Mock DOM manipulation methods
        document.body.appendChild = jest.fn();
        document.body.removeChild = jest.fn();
        document.createElement = jest.fn(() => ({
            href: '',
            download: '',
            click: jest.fn(),
            style: { cssText: '' },
            innerHTML: '',
            id: '',
            remove: jest.fn()
        }));

        // Mock console methods to avoid spam during tests
        jest.spyOn(console, 'log').mockImplementation(() => {});
        jest.spyOn(console, 'warn').mockImplementation(() => {});
        jest.spyOn(console, 'error').mockImplementation(() => {});

        // Mock URL and Blob for export functionality
        global.URL = {
            createObjectURL: jest.fn(() => 'mock-url'),
            revokeObjectURL: jest.fn()
        };
        global.Blob = jest.fn();

        performanceMonitor = new PerformanceMonitor();
    });

    afterEach(() => {
        if (performanceMonitor) {
            performanceMonitor.destroy();
        }
        
        // Restore original methods
        document.body.appendChild = originalAppendChild;
        document.body.removeChild = originalRemoveChild;
        document.createElement = originalCreateElement;
        
        jest.restoreAllMocks();
        jest.clearAllTimers();
    });

    describe('Initialization', () => {
        test('should create instance with default metrics', () => {
            expect(performanceMonitor).toBeDefined();
            expect(performanceMonitor.isActive).toBe(true);
            expect(performanceMonitor.frameCount).toBe(0);
            expect(performanceMonitor.metrics).toEqual(expect.objectContaining({
                fps: 0,
                renderTime: 0,
                memoryUsage: 0,
                platformCount: 0,
                visibleCount: 0,
                updateRate: 0,
                dataRate: 0,
                latency: 0
            }));
        });

        test('should bind DOM elements correctly', () => {
            expect(document.getElementById).toHaveBeenCalledWith('fps');
            expect(document.getElementById).toHaveBeenCalledWith('platformCount');
            expect(document.getElementById).toHaveBeenCalledWith('visibleCount');
            expect(document.getElementById).toHaveBeenCalledWith('memoryUsage');
            expect(performanceMonitor.elements).toBeDefined();
        });
    });

    describe('FPS Tracking', () => {
        test('should update FPS metrics', () => {
            // Simulate frame updates
            performanceMonitor.frameCount = 60;
            performanceMonitor.lastFPSUpdate = performance.now() - 1000; // 1 second ago
            
            // This would normally be called by requestAnimationFrame
            performanceMonitor.currentFPS = 60;
            performanceMonitor.metrics.fps = 60;
            performanceMonitor.updateFPSDisplay();

            expect(performanceMonitor.metrics.fps).toBe(60);
            expect(mockElements.fps.textContent).toBe(60);
        });

        test('should color-code FPS display based on performance', () => {
            // Test high FPS (green)
            performanceMonitor.currentFPS = 60;
            performanceMonitor.updateFPSDisplay();
            expect(mockElements.fps.style.color).toBe('#4CAF50');

            // Test medium FPS (orange)
            performanceMonitor.currentFPS = 40;
            performanceMonitor.updateFPSDisplay();
            expect(mockElements.fps.style.color).toBe('#ff9800');

            // Test low FPS (red)
            performanceMonitor.currentFPS = 20;
            performanceMonitor.updateFPSDisplay();
            expect(mockElements.fps.style.color).toBe('#f44336');
        });
    });

    describe('Render Time Tracking', () => {
        test('should track render times correctly', () => {
            // Start with a fresh render timer
            performanceMonitor.renderStartTime = null;
            performanceMonitor.renderTimes = [];
            
            // Mock performance.now for consistent timing
            const mockNow = jest.fn()
                .mockReturnValueOnce(100) // start time
                .mockReturnValueOnce(116.5); // end time (16.5ms later)
            
            const originalNow = performance.now;
            performance.now = mockNow;

            performanceMonitor.startRenderTimer();
            expect(performanceMonitor.renderStartTime).toBe(100);

            performanceMonitor.endRenderTimer();

            expect(performanceMonitor.renderTimes).toContain(16.5);
            expect(performanceMonitor.metrics.renderTime).toBe(16.5);
            expect(performanceMonitor.renderStartTime).toBeNull();

            performance.now = originalNow;
        });

        test('should maintain render time history within limits', () => {
            // Clear existing render times first
            performanceMonitor.renderTimes = [];
            
            // Add exactly the max number of render times
            for (let i = 0; i < performanceMonitor.maxRenderTimeHistory; i++) {
                performanceMonitor.renderTimes.push(16.67);
            }

            // Mock performance.now for the new render time
            const originalNow = performance.now;
            performance.now = jest.fn()
                .mockReturnValueOnce(100) // start time
                .mockReturnValueOnce(120); // end time (20ms later)

            performanceMonitor.startRenderTimer();
            performanceMonitor.endRenderTimer();

            expect(performanceMonitor.renderTimes.length).toBeLessThanOrEqual(performanceMonitor.maxRenderTimeHistory);

            performance.now = originalNow;
        });
    });

    describe('Memory Tracking', () => {
        test('should update memory usage with performance.memory', () => {
            // Ensure performance.memory is available for this test
            const originalMemory = performance.memory;
            performance.memory = {
                usedJSHeapSize: 50 * 1024 * 1024, // 50MB
                totalJSHeapSize: 100 * 1024 * 1024, // 100MB
                jsHeapSizeLimit: 200 * 1024 * 1024 // 200MB
            };
            
            performanceMonitor.updateMemoryUsage();
            
            // Should use actual memory if available
            const expectedMemoryMB = Math.round(50 * 1024 * 1024 / (1024 * 1024)); // 50MB
            expect(performanceMonitor.metrics.memoryUsage).toBe(expectedMemoryMB);
            expect(mockElements.memoryUsage.textContent).toBe(`${expectedMemoryMB} MB`);
            
            // Restore original memory object
            performance.memory = originalMemory;
        });

        test('should fallback to estimation when performance.memory unavailable', () => {
            // Temporarily remove performance.memory
            const originalMemory = performance.memory;
            delete performance.memory;

            performanceMonitor.metrics.platformCount = 1000;
            performanceMonitor.metrics.visibleCount = 500;
            performanceMonitor.updateMemoryUsage();

            expect(performanceMonitor.metrics.memoryUsage).toBeGreaterThan(0);
            expect(mockElements.memoryUsage.textContent).toContain('~');

            performance.memory = originalMemory;
        });
    });

    describe('Platform Count Updates', () => {
        test('should update platform counts correctly', () => {
            performanceMonitor.updatePlatformCount(1500, 750);

            expect(performanceMonitor.metrics.platformCount).toBe(1500);
            expect(performanceMonitor.metrics.visibleCount).toBe(750);
            expect(mockElements.platformCount.textContent).toBe(1500);
            expect(mockElements.visibleCount.textContent).toBe(750);
        });

        test('should color-code visible count based on ratio', () => {
            // High visibility ratio (red)
            performanceMonitor.updatePlatformCount(100, 90);
            expect(mockElements.visibleCount.style.color).toBe('#f44336');

            // Medium visibility ratio (orange)
            performanceMonitor.updatePlatformCount(100, 60);
            expect(mockElements.visibleCount.style.color).toBe('#ff9800');

            // Low visibility ratio (green)
            performanceMonitor.updatePlatformCount(100, 30);
            expect(mockElements.visibleCount.style.color).toBe('#4CAF50');
        });
    });

    describe('Data Metrics', () => {
        test('should update data streaming metrics', () => {
            performanceMonitor.updateDataMetrics(30, 51200); // 30 updates/s, 50KB/s

            expect(performanceMonitor.metrics.updateRate).toBe(30);
            expect(performanceMonitor.metrics.dataRate).toBe(51200);
            expect(mockElements.updateRate.textContent).toBe(30);
            expect(mockElements.dataRate.textContent).toBe('50 KB/s');
        });
    });

    describe('Performance Analysis', () => {
        test('should generate performance report', () => {
            // Set up some test data
            performanceMonitor.metrics.fps = 45;
            performanceMonitor.renderTimes = [16.5, 17.2, 15.8];

            const report = performanceMonitor.getPerformanceReport();

            expect(report).toEqual(expect.objectContaining({
                fps: 45,
                averageRenderTime: expect.any(Number),
                maxRenderTime: expect.any(Number),
                minRenderTime: expect.any(Number),
                renderTimeStdDev: expect.any(Number),
                timestamp: expect.any(Number)
            }));
        });

        test('should calculate standard deviation correctly', () => {
            const values = [10, 12, 14, 16, 18];
            const stdDev = performanceMonitor.calculateStandardDeviation(values);
            
            expect(stdDev).toBeCloseTo(2.83, 1); // Expected standard deviation
        });

        test('should return 0 for empty array standard deviation', () => {
            const stdDev = performanceMonitor.calculateStandardDeviation([]);
            expect(stdDev).toBe(0);
        });
    });

    describe('Performance Warnings', () => {
        test('should detect performance issues', () => {
            performanceMonitor.metrics.fps = 25; // Low FPS
            performanceMonitor.metrics.renderTime = 40; // High render time
            performanceMonitor.metrics.memoryUsage = 600; // High memory
            performanceMonitor.metrics.visibleCount = 6000; // Too many objects

            const warnings = performanceMonitor.checkPerformanceThresholds();

            expect(warnings).toContain('Low FPS: 25');
            expect(warnings).toContain('High render time: 40ms');
            expect(warnings).toContain('High memory usage: 600MB');
            expect(warnings).toContain('Too many visible objects: 6000');
        });

        test('should return no warnings for good performance', () => {
            performanceMonitor.metrics.fps = 60;
            performanceMonitor.metrics.renderTime = 15;
            performanceMonitor.metrics.memoryUsage = 100;
            performanceMonitor.metrics.visibleCount = 1000;

            const warnings = performanceMonitor.checkPerformanceThresholds();

            expect(warnings).toHaveLength(0);
        });
    });

    describe('Optimization Suggestions', () => {
        test('should provide optimization suggestions', () => {
            performanceMonitor.metrics.visibleCount = 3000;
            performanceMonitor.metrics.renderTime = 20;
            performanceMonitor.metrics.fps = 40;
            performanceMonitor.metrics.memoryUsage = 400;

            const suggestions = performanceMonitor.getOptimizationSuggestions();

            expect(suggestions.length).toBeGreaterThan(0);
            expect(suggestions.some(s => s.includes('clustering'))).toBe(true);
            expect(suggestions.some(s => s.includes('trail length'))).toBe(true);
            expect(suggestions.some(s => s.includes('zooming out'))).toBe(true);
            expect(suggestions.some(s => s.includes('object pooling'))).toBe(true);
        });
    });

    describe('Public API', () => {
        test('should start and stop monitoring', () => {
            performanceMonitor.stop();
            expect(performanceMonitor.isActive).toBe(false);

            performanceMonitor.start();
            expect(performanceMonitor.isActive).toBe(true);
        });

        test('should reset metrics', () => {
            // Set some values
            performanceMonitor.frameCount = 60;
            performanceMonitor.renderTimes = [16, 17, 18];
            performanceMonitor.metrics.fps = 60;

            performanceMonitor.reset();

            expect(performanceMonitor.frameCount).toBe(0);
            expect(performanceMonitor.renderTimes).toHaveLength(0);
            expect(performanceMonitor.metrics.fps).toBe(0);
        });

        test('should record platform count', () => {
            performanceMonitor.recordPlatformCount(500);
            
            expect(performanceMonitor.metrics.platformCount).toBe(500);
            expect(mockElements.platformCount.textContent).toBe(500);
        });

        test('should record errors', () => {
            const errorSpy = jest.spyOn(console, 'error').mockImplementation(() => {});
            
            performanceMonitor.recordError('Test error');
            
            expect(errorSpy).toHaveBeenCalledWith('Performance Monitor - Error recorded:', 'Test error');
        });

        test('should set up update callback', () => {
            jest.useFakeTimers();
            const mockCallback = jest.fn();
            
            performanceMonitor.metrics.fps = 60;
            performanceMonitor.metrics.visibleCount = 100;
            performanceMonitor.metrics.memoryUsage = 50;
            performanceMonitor.metrics.dataRate = 1024;
            
            performanceMonitor.onUpdate(mockCallback);
            
            // Fast-forward time to trigger the callback
            jest.advanceTimersByTime(1000);
            
            expect(mockCallback).toHaveBeenCalledWith({
                fps: 60,
                visibleCount: 100,
                memoryMB: 50,
                dataRate: 1024
            });
            
            jest.useRealTimers();
        });
    });

    describe('Export Functionality', () => {
        test('should export performance data', () => {
            performanceMonitor.exportPerformanceData();

            expect(global.Blob).toHaveBeenCalled();
            expect(global.URL.createObjectURL).toHaveBeenCalled();
            expect(document.createElement).toHaveBeenCalledWith('a');
        });
    });

    describe('Debug Overlay', () => {
        test('should create performance debug overlay', () => {
            jest.useFakeTimers();
            
            const overlay = performanceMonitor.createPerformanceOverlay();
            
            expect(overlay).toBeDefined();
            expect(overlay.id).toBe('performance-debug-overlay');
            expect(document.body.appendChild).toHaveBeenCalledWith(overlay);
            
            jest.useRealTimers();
        });
    });

    describe('Cleanup', () => {
        test('should destroy monitor correctly', () => {
            // Create a mock overlay element
            const mockOverlay = { remove: jest.fn() };
            document.getElementById = jest.fn((id) => {
                if (id === 'performance-debug-overlay') return mockOverlay;
                return mockElements[id] || null;
            });

            performanceMonitor.destroy();

            expect(performanceMonitor.isActive).toBe(false);
        });
    });
});
