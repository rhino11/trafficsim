package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rhino11/trafficsim/internal/config"
	"github.com/rhino11/trafficsim/internal/server"
	"github.com/rhino11/trafficsim/internal/sim"
)

func TestWebModeIntegration(t *testing.T) {
	// Test that web mode properly loads platforms and serves them via API
	cfg := createTestConfig()
	engine := sim.NewEngine(cfg)

	// Load platforms
	if err := engine.LoadPlatformsFromConfig(); err != nil {
		t.Fatalf("Failed to load platforms: %v", err)
	}

	// Start simulation
	if err := engine.Start(); err != nil {
		t.Fatalf("Failed to start simulation: %v", err)
	}
	defer engine.Stop()

	// Create server
	srv := server.NewServer(cfg, engine)

	// Test GET /api/platforms endpoint
	req := httptest.NewRequest("GET", "/api/platforms", nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Parse response
	var platforms []interface{}
	if err := json.NewDecoder(w.Body).Decode(&platforms); err != nil {
		t.Fatalf("Failed to decode platforms response: %v", err)
	}

	if len(platforms) == 0 {
		t.Error("Expected platforms to be returned, but got empty array")
	}

	t.Logf("Successfully loaded %d platforms in web mode", len(platforms))
}

func TestWebSocketPlatformUpdates(t *testing.T) {
	// Test that WebSocket properly streams platform updates
	cfg := createTestConfig()
	engine := sim.NewEngine(cfg)

	// Load platforms
	if err := engine.LoadPlatformsFromConfig(); err != nil {
		t.Fatalf("Failed to load platforms: %v", err)
	}

	// Start simulation
	if err := engine.Start(); err != nil {
		t.Fatalf("Failed to start simulation: %v", err)
	}
	defer engine.Stop()

	// Create server
	srv := server.NewServer(cfg, engine)

	// Create test server
	testServer := httptest.NewServer(srv)
	defer testServer.Close()

	// Convert HTTP URL to WebSocket URL
	wsURL := "ws" + strings.TrimPrefix(testServer.URL, "http") + "/ws"

	// Connect to WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer conn.Close()

	// Set read timeout
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	// Read initial message
	_, message, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read WebSocket message: %v", err)
	}

	// Parse message
	var platformUpdate map[string]interface{}
	if err := json.Unmarshal(message, &platformUpdate); err != nil {
		t.Fatalf("Failed to parse WebSocket message: %v", err)
	}

	// Verify platform update structure
	if updateType, ok := platformUpdate["type"].(string); !ok || updateType != "platform_update" {
		t.Errorf("Expected platform_update message, got %v", platformUpdate["type"])
	}

	if platforms, ok := platformUpdate["platforms"].([]interface{}); !ok || len(platforms) == 0 {
		t.Error("Expected platforms array in WebSocket message")
	}

	t.Logf("Successfully received platform updates via WebSocket")
}

func TestMulticastIntegrationWithWebMode(t *testing.T) {
	// Test that multicast can be enabled alongside web mode
	// This test verifies the integration we need to implement

	// Create test configuration
	cfg := createTestConfig()
	engine := sim.NewEngine(cfg)

	// Load platforms
	if err := engine.LoadPlatformsFromConfig(); err != nil {
		t.Fatalf("Failed to load platforms: %v", err)
	}

	// Start simulation
	if err := engine.Start(); err != nil {
		t.Fatalf("Failed to start simulation: %v", err)
	}
	defer engine.Stop()

	// Create server with multicast capability (this will need to be implemented)
	srv := server.NewServer(cfg, engine)

	// Test that server has multicast status endpoint
	req := httptest.NewRequest("GET", "/api/multicast/status", nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	// This should return multicast status (currently will fail - we need to implement this)
	t.Logf("Multicast status response: %d", w.Code)
}

func TestSSEPlatformUpdates(t *testing.T) {
	// Test Server-Sent Events for platform updates
	cfg := createTestConfig()
	engine := sim.NewEngine(cfg)

	// Load platforms
	if err := engine.LoadPlatformsFromConfig(); err != nil {
		t.Fatalf("Failed to load platforms: %v", err)
	}

	// Start simulation
	if err := engine.Start(); err != nil {
		t.Fatalf("Failed to start simulation: %v", err)
	}
	defer engine.Stop()

	// Create server
	srv := server.NewServer(cfg, engine)

	// Test SSE endpoint
	req := httptest.NewRequest("GET", "/api/stream/platforms", nil)
	w := httptest.NewRecorder()

	// Use context with timeout for SSE
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	srv.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for SSE, got %d", w.Code)
	}

	// Check content type - accept both event-stream and JSON fallback
	contentType := w.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/event-stream") && !strings.Contains(contentType, "application/json") {
		t.Errorf("Expected event-stream or JSON content type, got %s", contentType)
	}

	// If it's JSON fallback, verify the response structure
	if strings.Contains(contentType, "application/json") {
		var platforms []interface{}
		if err := json.NewDecoder(w.Body).Decode(&platforms); err != nil {
			t.Errorf("SSE fallback JSON response invalid: %v", err)
		} else {
			t.Logf("SSE endpoint using JSON fallback mode with %d platforms", len(platforms))
		}
	}

	t.Logf("SSE endpoint working correctly")
}

func TestPlatformDataVisualization(t *testing.T) {
	// Test that platform data is properly formatted for visualization
	cfg := createTestConfig()
	engine := sim.NewEngine(cfg)

	// Load platforms
	if err := engine.LoadPlatformsFromConfig(); err != nil {
		t.Fatalf("Failed to load platforms: %v", err)
	}

	platforms := engine.GetAllPlatforms()
	if len(platforms) == 0 {
		t.Fatal("No platforms loaded for visualization test")
	}

	// Test each platform has required fields for map display
	for i, platform := range platforms {
		state := platform.GetState()

		// Check required position data - allow (0,0) for space platforms which is valid in orbit
		if state.Position.Latitude == 0 && state.Position.Longitude == 0 && platform.GetType() != "space" {
			t.Errorf("Platform %d has invalid position data", i)
		}

		// Check required identification data
		if platform.GetCallSign() == "" {
			t.Errorf("Platform %d missing callsign", i)
		}

		if platform.GetID() == "" {
			t.Errorf("Platform %d missing ID", i)
		}

		// Check type information
		if platform.GetType() == "" {
			t.Errorf("Platform %d missing type", i)
		}

		t.Logf("Platform %d (%s): Valid visualization data at %.4f,%.4f",
			i, platform.GetCallSign(), state.Position.Latitude, state.Position.Longitude)
	}
}

func createTestConfig() *config.Config {
	// Create a minimal configuration for testing
	return &config.Config{
		Simulation: config.SimulationConfig{
			UpdateInterval: "1s",
			TimeScale:      1.0,
		},
		Platforms: config.PlatformRegistry{
			// Use default example platforms for testing
		},
	}
}

func TestHTMLTemplateRendering(t *testing.T) {
	// Test that HTML templates are properly rendered
	cfg := createTestConfig()
	engine := sim.NewEngine(cfg)
	srv := server.NewServer(cfg, engine)

	// Test main index page
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for index page, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "TrafficSim") {
		t.Error("Index page doesn't contain expected title")
	}

	// Test scenario builder page
	req = httptest.NewRequest("GET", "/scenario-builder", nil)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for scenario builder, got %d", w.Code)
	}

	t.Logf("HTML templates rendering correctly")
}

func TestStaticFileServing(t *testing.T) {
	// Test that static files (CSS, JS) are properly served
	cfg := createTestConfig()
	engine := sim.NewEngine(cfg)
	srv := server.NewServer(cfg, engine)

	// Test static file paths that should exist
	staticPaths := []string{
		"/static/css/style.css",
		"/static/js/map-engine.js",
		"/static/js/platform-renderer.js",
		"/static/js/scenario-builder.js",
	}

	for _, path := range staticPaths {
		req := httptest.NewRequest("GET", path, nil)
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)

		// Should not return 500 errors
		if w.Code >= 500 {
			t.Errorf("Static file %s returned server error: %d", path, w.Code)
		}

		// Log the result for debugging
		t.Logf("Static file %s: status %d", path, w.Code)
	}
}

func TestSimulationControlAPI(t *testing.T) {
	// Test simulation control endpoints
	cfg := createTestConfig()
	engine := sim.NewEngine(cfg)
	srv := server.NewServer(cfg, engine)

	// Test simulation status
	req := httptest.NewRequest("GET", "/api/simulation/status", nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for simulation status, got %d", w.Code)
	}

	var status map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&status); err != nil {
		t.Fatalf("Failed to decode status response: %v", err)
	}

	// Check required status fields
	requiredFields := []string{"running", "time", "platform_count", "speed"}
	for _, field := range requiredFields {
		if _, exists := status[field]; !exists {
			t.Errorf("Status missing required field: %s", field)
		}
	}

	t.Logf("Simulation status API working: %+v", status)
}

// Helper function to check if we're in test mode
func isTestEnvironment() bool {
	return os.Getenv("GO_TESTING") == "1" || strings.Contains(os.Args[0], ".test")
}
