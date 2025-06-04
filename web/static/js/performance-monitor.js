/**
 * Performance Monitor for High-Performance Map Rendering
 * Tracks FPS, memory usage, and rendering metrics
 */
class PerformanceMonitor {
    constructor() {
        this.isActive = true;
        this.frameCount = 0;
        this.lastFPSUpdate = performance.now();
        this.currentFPS = 0;

        // Performance metrics
        this.metrics = {
            fps: 0,
            renderTime: 0,
            memoryUsage: 0,
            platformCount: 0,
            visibleCount: 0,
            updateRate: 0,
            dataRate: 0,
            latency: 0
        };

        // DOM elements
        this.elements = {};

        // Performance tracking
        this.renderTimes = [];
        this.maxRenderTimeHistory = 60; // Keep 1 second at 60fps

        // Memory tracking
        this.lastMemoryCheck = 0;
        this.memoryCheckInterval = 5000; // Check every 5 seconds

        this.init();
    }

    init() {
        console.log('Initializing PerformanceMonitor...');
        this.bindDOMElements();
        this.startFPSTracking();
        this.setupMemoryTracking();
        console.log('PerformanceMonitor initialized');
    }

    bindDOMElements() {
        this.elements = {
            fps: document.getElementById('fps'),
            platformCount: document.getElementById('platformCount'),
            visibleCount: document.getElementById('visibleCount'),
            memoryUsage: document.getElementById('memoryUsage'),
            renderTime: document.getElementById('renderTime'),
            updateRate: document.getElementById('updateRate'),
            dataRate: document.getElementById('dataRate')
        };

        // Check if all elements are found
        Object.keys(this.elements).forEach(key => {
            if (!this.elements[key]) {
                console.warn(`Performance monitor element not found: ${key}`);
            }
        });
    }

    startFPSTracking() {
        const trackFPS = () => {
            if (!this.isActive) return;

            this.frameCount++;
            const now = performance.now();
            const delta = now - this.lastFPSUpdate;

            // Update FPS every second
            if (delta >= 1000) {
                this.currentFPS = Math.round((this.frameCount * 1000) / delta);
                this.metrics.fps = this.currentFPS;
                this.updateFPSDisplay();

                this.frameCount = 0;
                this.lastFPSUpdate = now;
            }

            requestAnimationFrame(trackFPS);
        };

        requestAnimationFrame(trackFPS);
    }

    setupMemoryTracking() {
        setInterval(() => {
            if (!this.isActive) return;
            this.updateMemoryUsage();
        }, this.memoryCheckInterval);
    }

    updateMemoryUsage() {
        if (performance.memory) {
            // Chrome/Chromium browsers
            const used = performance.memory.usedJSHeapSize;
            const total = performance.memory.totalJSHeapSize;

            this.metrics.memoryUsage = Math.round(used / (1024 * 1024)); // MB

            if (this.elements.memoryUsage) {
                this.elements.memoryUsage.textContent = `${this.metrics.memoryUsage} MB`;
            }

            // Warn if memory usage is high
            const usagePercent = (used / total) * 100;
            if (usagePercent > 80) {
                console.warn(`High memory usage: ${usagePercent.toFixed(1)}%`);
            }
        } else {
            // Fallback estimation based on number of objects
            const estimatedMemory = Math.round((this.metrics.platformCount * 2 + this.metrics.visibleCount * 5) / 1024);
            this.metrics.memoryUsage = estimatedMemory;

            if (this.elements.memoryUsage) {
                this.elements.memoryUsage.textContent = `~${estimatedMemory} MB`;
            }
        }
    }

    updateFPSDisplay() {
        if (this.elements.fps) {
            this.elements.fps.textContent = this.currentFPS;

            // Color code FPS for visual feedback
            if (this.currentFPS >= 55) {
                this.elements.fps.style.color = '#4CAF50'; // Green
            } else if (this.currentFPS >= 30) {
                this.elements.fps.style.color = '#ff9800'; // Orange
            } else {
                this.elements.fps.style.color = '#f44336'; // Red
            }
        }
    }

    // Track rendering performance
    startRenderTimer() {
        this.renderStartTime = performance.now();
    }

    endRenderTimer() {
        if (this.renderStartTime) {
            const renderTime = performance.now() - this.renderStartTime;
            this.renderTimes.push(renderTime);

            // Keep only recent render times
            if (this.renderTimes.length > this.maxRenderTimeHistory) {
                this.renderTimes.shift();
            }

            // Calculate average render time
            const avgRenderTime = this.renderTimes.reduce((a, b) => a + b, 0) / this.renderTimes.length;
            this.metrics.renderTime = Math.round(avgRenderTime * 100) / 100;

            this.updateRenderTimeDisplay();
            this.renderStartTime = null;
        }
    }

    updateRenderTimeDisplay() {
        if (this.elements.renderTime) {
            this.elements.renderTime.textContent = `${this.metrics.renderTime}`;

            // Color code render time
            if (this.metrics.renderTime <= 16.67) { // 60 FPS target
                this.elements.renderTime.style.color = '#4CAF50'; // Green
            } else if (this.metrics.renderTime <= 33.33) { // 30 FPS
                this.elements.renderTime.style.color = '#ff9800'; // Orange
            } else {
                this.elements.renderTime.style.color = '#f44336'; // Red
            }
        }
    }

    // Update platform counts
    updatePlatformCount(total, visible) {
        this.metrics.platformCount = total;
        this.metrics.visibleCount = visible;

        if (this.elements.platformCount) {
            this.elements.platformCount.textContent = total;
        }

        if (this.elements.visibleCount) {
            this.elements.visibleCount.textContent = visible;

            // Color code based on visibility ratio
            const visibilityRatio = visible / total;
            if (visibilityRatio > 0.8) {
                this.elements.visibleCount.style.color = '#f44336'; // Red (too many visible)
            } else if (visibilityRatio > 0.5) {
                this.elements.visibleCount.style.color = '#ff9800'; // Orange
            } else {
                this.elements.visibleCount.style.color = '#4CAF50'; // Green
            }
        }
    }

    // Update data streaming metrics
    updateDataMetrics(updateRate, dataRate) {
        this.metrics.updateRate = updateRate;
        this.metrics.dataRate = dataRate;

        if (this.elements.updateRate) {
            this.elements.updateRate.textContent = updateRate;
        }

        if (this.elements.dataRate) {
            this.elements.dataRate.textContent = `${Math.round(dataRate / 1024)} KB/s`;
        }
    }

    // Performance analysis
    getPerformanceReport() {
        const report = {
            ...this.metrics,
            timestamp: Date.now(),
            averageRenderTime: this.renderTimes.length > 0 ?
                this.renderTimes.reduce((a, b) => a + b, 0) / this.renderTimes.length : 0,
            maxRenderTime: this.renderTimes.length > 0 ? Math.max(...this.renderTimes) : 0,
            minRenderTime: this.renderTimes.length > 0 ? Math.min(...this.renderTimes) : 0,
            renderTimeStdDev: this.calculateStandardDeviation(this.renderTimes)
        };

        return report;
    }

    calculateStandardDeviation(values) {
        if (values.length === 0) return 0;

        const mean = values.reduce((a, b) => a + b, 0) / values.length;
        const squaredDifferences = values.map(value => Math.pow(value - mean, 2));
        const avgSquaredDiff = squaredDifferences.reduce((a, b) => a + b, 0) / values.length;

        return Math.sqrt(avgSquaredDiff);
    }

    // Performance warnings
    checkPerformanceThresholds() {
        const warnings = [];

        if (this.metrics.fps < 30) {
            warnings.push(`Low FPS: ${this.metrics.fps}`);
        }

        if (this.metrics.renderTime > 33) {
            warnings.push(`High render time: ${this.metrics.renderTime}ms`);
        }

        if (this.metrics.memoryUsage > 500) {
            warnings.push(`High memory usage: ${this.metrics.memoryUsage}MB`);
        }

        if (this.metrics.visibleCount > 5000) {
            warnings.push(`Too many visible objects: ${this.metrics.visibleCount}`);
        }

        return warnings;
    }

    // Optimization suggestions
    getOptimizationSuggestions() {
        const suggestions = [];

        if (this.metrics.visibleCount > 2000) {
            suggestions.push('Enable clustering for better performance with many objects');
        }

        if (this.metrics.renderTime > 16.67) {
            suggestions.push('Consider reducing trail length or marker complexity');
        }

        if (this.metrics.fps < 45) {
            suggestions.push('Try zooming out or enabling viewport culling');
        }

        if (this.metrics.memoryUsage > 300) {
            suggestions.push('Consider implementing object pooling for markers');
        }

        return suggestions;
    }

    // Log performance data (for debugging)
    logPerformanceData() {
        const report = this.getPerformanceReport();
        const warnings = this.checkPerformanceThresholds();
        const suggestions = this.getOptimizationSuggestions();

        console.group('Performance Report');
        console.table(report);

        if (warnings.length > 0) {
            console.warn('Performance Warnings:', warnings);
        }

        if (suggestions.length > 0) {
            console.info('Optimization Suggestions:', suggestions);
        }

        console.groupEnd();
    }

    // Export performance data
    exportPerformanceData() {
        const data = {
            timestamp: new Date().toISOString(),
            metrics: this.getPerformanceReport(),
            renderTimeHistory: [...this.renderTimes],
            warnings: this.checkPerformanceThresholds(),
            suggestions: this.getOptimizationSuggestions()
        };

        // Create downloadable JSON file
        const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' });
        const url = URL.createObjectURL(blob);

        const a = document.createElement('a');
        a.href = url;
        a.download = `performance-report-${Date.now()}.json`;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);
    }

    // Real-time performance overlay (for debugging)
    createPerformanceOverlay() {
        const overlay = document.createElement('div');
        overlay.id = 'performance-debug-overlay';
        overlay.style.cssText = `
            position: fixed;
            top: 100px;
            right: 10px;
            background: rgba(0, 0, 0, 0.9);
            color: white;
            padding: 10px;
            border-radius: 5px;
            font-family: monospace;
            font-size: 12px;
            z-index: 10000;
            min-width: 200px;
        `;

        document.body.appendChild(overlay);

        // Update overlay every second
        setInterval(() => {
            const report = this.getPerformanceReport();
            overlay.innerHTML = `
                <div><strong>Debug Performance</strong></div>
                <div>FPS: ${report.fps}</div>
                <div>Render: ${report.renderTime.toFixed(2)}ms</div>
                <div>Avg Render: ${report.averageRenderTime.toFixed(2)}ms</div>
                <div>Max Render: ${report.maxRenderTime.toFixed(2)}ms</div>
                <div>Memory: ${report.memoryUsage}MB</div>
                <div>Objects: ${report.platformCount}</div>
                <div>Visible: ${report.visibleCount}</div>
                <div>Update Rate: ${report.updateRate}/s</div>
            `;
        }, 1000);

        return overlay;
    }

    // Add missing methods that the HTML template expects
    recordPlatformCount(count) {
        this.metrics.platformCount = count;
        if (this.elements.platformCount) {
            this.elements.platformCount.textContent = count;
        }
    }

    recordError(error) {
        console.error('Performance Monitor - Error recorded:', error);
        // Could implement error tracking here
    }

    onUpdate(callback) {
        this.updateCallback = callback;

        // Set up interval to call the callback with current stats
        setInterval(() => {
            if (this.updateCallback && this.isActive) {
                const stats = {
                    fps: this.metrics.fps,
                    visibleCount: this.metrics.visibleCount,
                    memoryMB: this.metrics.memoryUsage,
                    dataRate: this.metrics.dataRate
                };
                this.updateCallback(stats);
            }
        }, 1000); // Update every second
    }

    // Public API methods
    start() {
        this.isActive = true;
    }

    stop() {
        this.isActive = false;
    }

    reset() {
        this.frameCount = 0;
        this.renderTimes = [];
        this.metrics = {
            fps: 0,
            renderTime: 0,
            memoryUsage: 0,
            platformCount: 0,
            visibleCount: 0,
            updateRate: 0,
            dataRate: 0,
            latency: 0
        };
    }

    // Cleanup
    destroy() {
        this.isActive = false;

        // Remove debug overlay if it exists
        const overlay = document.getElementById('performance-debug-overlay');
        if (overlay) {
            overlay.remove();
        }
    }
}

// Export for use in other modules
window.PerformanceMonitor = PerformanceMonitor;
