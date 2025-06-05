package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/rhino11/trafficsim/internal/config"
	"github.com/rhino11/trafficsim/internal/sim"
	"github.com/rhino11/trafficsim/internal/testutil"
)

// createTestConfig creates a basic config for testing
func createTestConfig() *config.Config {
	return &config.Config{
		Server: config.ServerConfig{
			Port: 8080,
			Host: "localhost",
		},
		Simulation: config.SimulationConfig{
			UpdateInterval: "1s",
		},
	}
}

// createTestEngine creates a basic simulation engine for testing
func createTestEngine() *sim.Engine {
	cfg := createTestConfig()
	engine := sim.NewEngine(cfg)
	return engine
}

func TestNewServer(t *testing.T) {
	logger := testutil.SetupTestLogging(t)
	logger.Info("Testing server creation")

	cfg := createTestConfig()
	engine := createTestEngine()

	server := NewServer(cfg, engine)

	if server == nil {
		t.Fatal("Expected server to be created")
	}

	if server.config != cfg {
		t.Error("Server config not set correctly")
	}

	if server.simulation != engine {
		t.Error("Simulation engine not set correctly")
	}

	if server.router == nil {
		t.Error("Router should be initialized")
	}

	if server.clients == nil {
		t.Error("Clients map should be initialized")
	}

	if server.broadcast == nil {
		t.Error("Broadcast channel should be initialized")
	}

	logger.Info("Server creation test completed successfully")
}

func TestSetupRoutes(t *testing.T) {
	cfg := createTestConfig()
	engine := createTestEngine()
	server := NewServer(cfg, engine)

	// Test that routes are properly registered by making test requests
	testRoutes := []struct {
		path           string
		method         string
		expectNotFound bool // Some routes may legitimately return 404 (like static files that don't exist)
	}{
		{"/", "GET", false},
		{"/scenario-builder", "GET", false},
		{"/api/platforms", "GET", false},
		{"/api/platform-types", "GET", false},
		{"/api/simulation/status", "GET", false},
		{"/api/simulation/start", "POST", false},
		{"/api/simulation/stop", "POST", false},
		{"/api/simulation/reset", "POST", false},
		{"/api/stream/platforms", "GET", false},
		{"/api/metrics", "GET", false},
		{"/api/scenarios", "POST", false},
		{"/static/test.css", "GET", true}, // This file doesn't exist, so 404 is expected
		{"/ws", "GET", false},
	}

	for _, route := range testRoutes {
		req := httptest.NewRequest(route.method, route.path, nil)
		rec := httptest.NewRecorder()
		server.router.ServeHTTP(rec, req)

		if route.expectNotFound {
			// For static files that don't exist, 404 is expected
			if rec.Code != 404 {
				t.Errorf("Route %s %s expected 404 for non-existent file, got %d", route.method, route.path, rec.Code)
			}
		} else {
			// Should not return 404 for registered routes (though they may return other errors)
			if rec.Code == 404 {
				t.Errorf("Route %s %s returned 404, expected route to be registered", route.method, route.path)
			}
		}
	}
}

func TestHandleIndex(t *testing.T) {
	cfg := createTestConfig()
	engine := createTestEngine()
	server := NewServer(cfg, engine)

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	server.router.ServeHTTP(rec, req)

	// Template may not exist in test environment, so we accept either success or template error
	if rec.Code != http.StatusOK && rec.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 200 or 500 (template error), got %d", rec.Code)
	}
}

func TestHandleScenarioBuilder(t *testing.T) {
	cfg := createTestConfig()
	engine := createTestEngine()
	server := NewServer(cfg, engine)

	req := httptest.NewRequest("GET", "/scenario-builder", nil)
	rec := httptest.NewRecorder()

	server.router.ServeHTTP(rec, req)

	// Template may not exist in test environment, so we accept either success or template error
	if rec.Code != http.StatusOK && rec.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 200 or 500 (template error), got %d", rec.Code)
	}
}

func TestHandleGetPlatforms(t *testing.T) {
	cfg := createTestConfig()
	engine := createTestEngine()
	server := NewServer(cfg, engine)

	req := httptest.NewRequest("GET", "/api/platforms", nil)
	rec := httptest.NewRecorder()

	server.router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	// Should return valid JSON array
	var platforms []interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &platforms); err != nil {
		t.Errorf("Expected valid JSON array, got error: %v", err)
	}
}

func TestHandleGetPlatformTypes(t *testing.T) {
	cfg := createTestConfig()
	engine := createTestEngine()
	server := NewServer(cfg, engine)

	req := httptest.NewRequest("GET", "/api/platform-types", nil)
	rec := httptest.NewRecorder()

	server.router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	// Should return valid JSON array
	var platformTypes []PlatformTypeInfo
	if err := json.Unmarshal(rec.Body.Bytes(), &platformTypes); err != nil {
		t.Errorf("Expected valid JSON array of platform types, got error: %v", err)
	}

	// Should have at least some platform types (fallback data)
	if len(platformTypes) == 0 {
		t.Error("Expected at least some platform types to be returned")
	}
}

func TestHandleSimulationStatus(t *testing.T) {
	cfg := createTestConfig()
	engine := createTestEngine()
	server := NewServer(cfg, engine)

	req := httptest.NewRequest("GET", "/api/simulation/status", nil)
	rec := httptest.NewRecorder()

	server.router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	var status SimulationStatus
	if err := json.Unmarshal(rec.Body.Bytes(), &status); err != nil {
		t.Errorf("Expected valid SimulationStatus JSON, got error: %v", err)
	}

	// Should have required fields
	if status.Speed == 0 {
		t.Error("Expected Speed to be set")
	}
}

func TestHandleStartSimulation(t *testing.T) {
	cfg := createTestConfig()
	engine := createTestEngine()
	server := NewServer(cfg, engine)

	req := httptest.NewRequest("POST", "/api/simulation/start", nil)
	rec := httptest.NewRecorder()

	server.router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	var response map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Errorf("Expected valid JSON response, got error: %v", err)
	}

	if response["status"] != "started" {
		t.Errorf("Expected status 'started', got '%s'", response["status"])
	}
}

func TestHandleStopSimulation(t *testing.T) {
	cfg := createTestConfig()
	engine := createTestEngine()
	server := NewServer(cfg, engine)

	req := httptest.NewRequest("POST", "/api/simulation/stop", nil)
	rec := httptest.NewRecorder()

	server.router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	var response map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Errorf("Expected valid JSON response, got error: %v", err)
	}

	if response["status"] != "stopped" {
		t.Errorf("Expected status 'stopped', got '%s'", response["status"])
	}
}

func TestHandleResetSimulation(t *testing.T) {
	cfg := createTestConfig()
	engine := createTestEngine()
	server := NewServer(cfg, engine)

	req := httptest.NewRequest("POST", "/api/simulation/reset", nil)
	rec := httptest.NewRecorder()

	server.router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	var response map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Errorf("Expected valid JSON response, got error: %v", err)
	}

	if response["status"] != "reset" {
		t.Errorf("Expected status 'reset', got '%s'", response["status"])
	}
}

func TestHandleCreateScenario_ValidJSON(t *testing.T) {
	cfg := createTestConfig()
	engine := createTestEngine()
	server := NewServer(cfg, engine)

	scenarioJSON := `{
		"name": "test-scenario",
		"description": "A test scenario",
		"duration": 300,
		"platforms": []
	}`

	req := httptest.NewRequest("POST", "/api/scenarios", strings.NewReader(scenarioJSON))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	server.router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}
}

func TestHandleCreateScenario_InvalidJSON(t *testing.T) {
	cfg := createTestConfig()
	engine := createTestEngine()
	server := NewServer(cfg, engine)

	req := httptest.NewRequest("POST", "/api/scenarios", strings.NewReader("{invalid json"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	server.router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid JSON, got %d", rec.Code)
	}
}

func TestHandleWebSocket_InvalidUpgrade(t *testing.T) {
	cfg := createTestConfig()
	engine := createTestEngine()
	server := NewServer(cfg, engine)

	// Regular HTTP request without WebSocket upgrade headers
	req := httptest.NewRequest("GET", "/ws", nil)
	rec := httptest.NewRecorder()

	server.router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid WebSocket upgrade, got %d", rec.Code)
	}
}

func TestHandleMetrics(t *testing.T) {
	cfg := createTestConfig()
	engine := createTestEngine()
	server := NewServer(cfg, engine)

	req := httptest.NewRequest("GET", "/api/metrics", nil)
	rec := httptest.NewRecorder()

	server.router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	var metrics map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &metrics); err != nil {
		t.Errorf("Expected valid JSON metrics, got error: %v", err)
	}

	// Should have required metric sections
	if _, exists := metrics["simulation"]; !exists {
		t.Error("Expected metrics to contain 'simulation' section")
	}
	if _, exists := metrics["server"]; !exists {
		t.Error("Expected metrics to contain 'server' section")
	}
	if _, exists := metrics["platforms"]; !exists {
		t.Error("Expected metrics to contain 'platforms' section")
	}
}

func TestHandleClientLog(t *testing.T) {
	cfg := createTestConfig()
	engine := createTestEngine()
	server := NewServer(cfg, engine)

	logData := map[string]interface{}{
		"type":      "client_log",
		"step":      "test",
		"message":   "test message",
		"timestamp": "2025-06-04T12:00:00Z",
		"userAgent": "test-agent",
	}

	jsonData, _ := json.Marshal(logData)
	req := httptest.NewRequest("POST", "/api/log", strings.NewReader(string(jsonData)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	server.router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", rec.Code)
	}
}

func TestHandleClientLog_InvalidJSON(t *testing.T) {
	cfg := createTestConfig()
	engine := createTestEngine()
	server := NewServer(cfg, engine)

	req := httptest.NewRequest("POST", "/api/log", strings.NewReader("{invalid json"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	server.router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid JSON, got %d", rec.Code)
	}
}

func TestServerLifecycle(t *testing.T) {
	cfg := createTestConfig()
	engine := createTestEngine()
	server := NewServer(cfg, engine)

	// Test Stop without Start (should not panic)
	server.Stop()

	// Test context cancellation
	if server.ctx.Err() != context.Canceled {
		t.Error("Expected context to be canceled after Stop()")
	}
}

func TestBroadcastSimulationStatus(t *testing.T) {
	cfg := createTestConfig()
	engine := createTestEngine()
	server := NewServer(cfg, engine)

	// Should not panic when no clients are connected
	server.broadcastSimulationStatus()

	// Test is primarily ensuring no panic occurs
}

func TestGenerateDescription(t *testing.T) {
	cfg := createTestConfig()
	engine := createTestEngine()
	server := NewServer(cfg, engine)

	tests := []struct {
		platformType string
		expected     string
	}{
		{"airbus_a320", "Short to medium-range commercial airliner"},
		{"boeing_777_300er", "Long-range wide-body commercial airliner"},
		{"f16_fighting_falcon", "Multi-role fighter aircraft"},
		{"container_ship", "Large cargo container vessel"},
		{"arleigh_burke_destroyer", "US Navy guided missile destroyer"},
		{"unknown_platform", "Unknown Platform platform"},
	}

	for _, test := range tests {
		result := server.generateDescription(test.platformType)
		if result != test.expected {
			t.Errorf("For platform type %s, expected %s, got %s", test.platformType, test.expected, result)
		}
	}
}

func TestDetermineDomain(t *testing.T) {
	cfg := createTestConfig()
	engine := createTestEngine()
	server := NewServer(cfg, engine)

	tests := []struct {
		platformType string
		expected     string
	}{
		{"airbus_a320", "airborne"},
		{"boeing_777_300er", "airborne"},
		{"f16_fighting_falcon", "airborne"},
		{"container_ship", "maritime"},
		{"arleigh_burke_destroyer", "maritime"},
		{"truck", "land"},
		{"tank", "land"},
		{"satellite", "space"},
		{"space_station", "space"},
		{"unknown_type", "unknown"},
	}

	for _, test := range tests {
		result := server.determineDomain(test.platformType)
		if result != test.expected {
			t.Errorf("For platform type %s, expected domain %s, got %s", test.platformType, test.expected, result)
		}
	}
}

func TestGetFallbackPlatformTypes(t *testing.T) {
	cfg := createTestConfig()
	engine := createTestEngine()
	server := NewServer(cfg, engine)

	fallbackTypes := server.getFallbackPlatformTypes()

	if len(fallbackTypes) == 0 {
		t.Error("Expected fallback platform types to be returned")
	}

	// Check that each platform type has required fields
	for _, platformType := range fallbackTypes {
		if platformType.ID == "" {
			t.Error("Platform type ID should not be empty")
		}
		if platformType.Name == "" {
			t.Error("Platform type Name should not be empty")
		}
		if platformType.Domain == "" {
			t.Error("Platform type Domain should not be empty")
		}
		if platformType.Performance == nil {
			t.Error("Platform type Performance should not be nil")
		}
	}
}

func TestLoggingMiddleware(t *testing.T) {
	cfg := createTestConfig()
	engine := createTestEngine()
	server := NewServer(cfg, engine)

	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("test response")); err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
	})

	// Wrap with logging middleware
	wrappedHandler := server.loggingMiddleware(testHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	if rec.Body.String() != "test response" {
		t.Errorf("Expected 'test response', got '%s'", rec.Body.String())
	}
}

func TestResponseWriter(t *testing.T) {
	rec := httptest.NewRecorder()
	wrapper := &responseWriter{ResponseWriter: rec, statusCode: 200}

	wrapper.WriteHeader(404)
	if wrapper.statusCode != 404 {
		t.Errorf("Expected status code 404, got %d", wrapper.statusCode)
	}

	if rec.Code != 404 {
		t.Errorf("Expected underlying recorder to have status 404, got %d", rec.Code)
	}
}

func TestLoadPlatformTypesFromFiles(t *testing.T) {
	cfg := createTestConfig()
	engine := createTestEngine()
	server := NewServer(cfg, engine)

	// Create temporary test platform files
	tempDir := t.TempDir()

	// Create test directory structure
	testDirs := []string{
		"airborne/commercial",
		"airborne/military",
		"maritime/commercial",
		"land/military",
		"space/commercial",
	}

	for _, dir := range testDirs {
		err := os.MkdirAll(fmt.Sprintf("%s/data/platforms/%s", tempDir, dir), 0755)
		if err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}
	}

	// Create test platform files
	testPlatforms := map[string]string{
		"airborne/commercial/boeing_737.yaml": `
platform_types:
  boeing_737_800:
    class: "Boeing 737-800"
    category: "commercial_aircraft"
    performance:
      max_speed: 257.0
      cruise_speed: 230.0
      max_altitude: 12000.0
`,
		"maritime/commercial/cargo_ship.yaml": `
platform_types:
  container_ship:
    class: "Container Ship"
    category: "cargo_vessel"
    performance:
      max_speed: 12.9
      cruise_speed: 10.8
`,
		"land/military/tank.yaml": `
platform_types:
  m1a2_abrams:
    class: "M1A2 Abrams"
    category: "main_battle_tank"
    performance:
      max_speed: 18.0
      cruise_speed: 12.0
`,
	}

	for filePath, content := range testPlatforms {
		fullPath := fmt.Sprintf("%s/data/platforms/%s", tempDir, filePath)
		err := os.WriteFile(fullPath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test platform file: %v", err)
		}
	}

	// Change working directory to temp dir for the test
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tempDir)

	// Test loading platforms from files
	platforms, err := server.loadPlatformTypesFromFiles()

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(platforms) == 0 {
		t.Error("Expected platforms to be loaded, got empty slice")
	}

	// Verify platform data
	platformMap := make(map[string]PlatformTypeInfo)
	for _, platform := range platforms {
		platformMap[platform.ID] = platform
	}

	// Check Boeing 737
	if boeing, exists := platformMap["boeing_737_800"]; exists {
		if boeing.Domain != "airborne" {
			t.Errorf("Expected domain 'airborne', got '%s'", boeing.Domain)
		}
		if boeing.Class != "Boeing 737-800" {
			t.Errorf("Expected class 'Boeing 737-800', got '%s'", boeing.Class)
		}
		if boeing.Performance["max_speed"] != 257.0 {
			t.Errorf("Expected max_speed 257.0, got %v", boeing.Performance["max_speed"])
		}
	} else {
		t.Error("Expected boeing_737_800 platform to be loaded")
	}

	// Check Container Ship
	if ship, exists := platformMap["container_ship"]; exists {
		if ship.Domain != "maritime" {
			t.Errorf("Expected domain 'maritime', got '%s'", ship.Domain)
		}
		if ship.Performance["cruise_speed"] != 10.8 {
			t.Errorf("Expected cruise_speed 10.8, got %v", ship.Performance["cruise_speed"])
		}
	} else {
		t.Error("Expected container_ship platform to be loaded")
	}
}

func TestLoadPlatformFromFile(t *testing.T) {
	cfg := createTestConfig()
	engine := createTestEngine()
	server := NewServer(cfg, engine)

	// Create temporary test directory with proper structure
	tempDir := t.TempDir()
	platformDir := fmt.Sprintf("%s/data/platforms/airborne/military", tempDir)
	err := os.MkdirAll(platformDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create test file in the proper location
	testFile := fmt.Sprintf("%s/test_platform.yaml", platformDir)
	testContent := `
platform_types:
  test_fighter:
    class: "Test Fighter"
    category: "fighter_aircraft"
    performance:
      max_speed: 600.0
      cruise_speed: 300.0
      max_altitude: 15000.0
`

	err = os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test content: %v", err)
	}

	// Change working directory to temp dir for the test
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tempDir)

	// Test loading platform from file with relative path
	platform, err := server.loadPlatformFromFile("data/platforms/airborne/military/test_platform.yaml", "airborne", "military")

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if platform == nil {
		t.Fatal("Expected platform to be loaded, got nil")
	}

	if platform.ID != "test_fighter" {
		t.Errorf("Expected ID 'test_fighter', got '%s'", platform.ID)
	}

	if platform.Domain != "airborne" {
		t.Errorf("Expected domain 'airborne', got '%s'", platform.Domain)
	}

	if platform.Class != "Test Fighter" {
		t.Errorf("Expected class 'Test Fighter', got '%s'", platform.Class)
	}

	if platform.Performance["max_altitude"] != 15000.0 {
		t.Errorf("Expected max_altitude 15000.0, got %v", platform.Performance["max_altitude"])
	}
}

func TestLoadPlatformFromFile_InvalidYAML(t *testing.T) {
	cfg := createTestConfig()
	engine := createTestEngine()
	server := NewServer(cfg, engine)

	// Create temporary test directory with proper structure
	tempDir := t.TempDir()
	platformDir := fmt.Sprintf("%s/data/platforms/airborne/commercial", tempDir)
	err := os.MkdirAll(platformDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create test file with invalid YAML in the proper location
	testFile := fmt.Sprintf("%s/invalid_platform.yaml", platformDir)
	invalidContent := `
platform_types:
  test_platform:
    class: "Test Platform"
    category: [invalid yaml structure
`

	err = os.WriteFile(testFile, []byte(invalidContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test content: %v", err)
	}

	// Change working directory to temp dir for the test
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tempDir)

	// Test loading platform from invalid file with relative path
	platform, err := server.loadPlatformFromFile("data/platforms/airborne/commercial/invalid_platform.yaml", "airborne", "commercial")

	if err == nil {
		t.Error("Expected error for invalid YAML, got nil")
	}

	if platform != nil {
		t.Error("Expected nil platform for invalid YAML, got non-nil")
	}
}

func TestLoadPlatformFromFile_NonExistentFile(t *testing.T) {
	cfg := createTestConfig()
	engine := createTestEngine()
	server := NewServer(cfg, engine)

	// Test loading from non-existent file in the expected directory structure
	platform, err := server.loadPlatformFromFile("data/platforms/airborne/commercial/nonexistent.yaml", "airborne", "commercial")

	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}

	if platform != nil {
		t.Error("Expected nil platform for non-existent file, got non-nil")
	}
}

func TestHandleGetPlatformTypes_FilesFirst(t *testing.T) {
	cfg := createTestConfig()
	engine := createTestEngine()
	server := NewServer(cfg, engine)

	// Create a test server that can load from files
	tempDir := t.TempDir()

	// Create test platform file
	platformDir := fmt.Sprintf("%s/data/platforms/airborne/commercial", tempDir)
	err := os.MkdirAll(platformDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	testPlatformFile := fmt.Sprintf("%s/test_aircraft.yaml", platformDir)
	testContent := `
platform_types:
  test_aircraft:
    class: "Test Aircraft"
    category: "test_aircraft"
    performance:
      max_speed: 500.0
      cruise_speed: 400.0
      max_altitude: 10000.0
`

	err = os.WriteFile(testPlatformFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test platform file: %v", err)
	}

	// Change working directory for the test
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tempDir)

	// Make HTTP request
	req := httptest.NewRequest("GET", "/api/platform-types", nil)
	rec := httptest.NewRecorder()

	server.handleGetPlatformTypes(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	// Parse response
	var platformTypes []PlatformTypeInfo
	if err := json.Unmarshal(rec.Body.Bytes(), &platformTypes); err != nil {
		t.Errorf("Expected valid JSON array of platform types, got error: %v", err)
	}

	// Should load from files (at least our test platform)
	found := false
	for _, platform := range platformTypes {
		if platform.ID == "test_aircraft" {
			found = true
			if platform.Class != "Test Aircraft" {
				t.Errorf("Expected class 'Test Aircraft', got '%s'", platform.Class)
			}
			break
		}
	}

	if !found {
		t.Error("Expected test_aircraft platform to be loaded from files")
	}
}

func TestHandleGetPlatformTypes_FallbackToConfig(t *testing.T) {
	cfg := createTestConfig()
	engine := createTestEngine()
	server := NewServer(cfg, engine)

	// Make request when no platform files exist
	req := httptest.NewRequest("GET", "/api/platform-types", nil)
	rec := httptest.NewRecorder()

	server.handleGetPlatformTypes(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	// Parse response
	var platformTypes []PlatformTypeInfo
	if err := json.Unmarshal(rec.Body.Bytes(), &platformTypes); err != nil {
		t.Errorf("Expected valid JSON array of platform types, got error: %v", err)
	}

	// Should have fallback platforms
	if len(platformTypes) == 0 {
		t.Error("Expected fallback platform types to be returned")
	}
}

func TestLoadPlatformTypesFromFiles_EmptyDirectories(t *testing.T) {
	cfg := createTestConfig()
	engine := createTestEngine()
	server := NewServer(cfg, engine)

	// Create temporary directory structure with empty directories
	tempDir := t.TempDir()

	testDirs := []string{
		"airborne/commercial",
		"maritime/military",
	}

	for _, dir := range testDirs {
		err := os.MkdirAll(fmt.Sprintf("%s/data/platforms/%s", tempDir, dir), 0755)
		if err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}
	}

	// Change working directory
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tempDir)

	// Test loading from empty directories
	platforms, err := server.loadPlatformTypesFromFiles()

	if err != nil {
		t.Errorf("Expected no error from empty directories, got: %v", err)
	}

	// Should fallback to default platforms when no files found
	if len(platforms) == 0 {
		t.Error("Expected fallback platforms when no files found")
	}
}

func TestLoadPlatformTypesFromFiles_MixedFileTypes(t *testing.T) {
	cfg := createTestConfig()
	engine := createTestEngine()
	server := NewServer(cfg, engine)

	// Create temporary test directory
	tempDir := t.TempDir()
	platformDir := fmt.Sprintf("%s/data/platforms/airborne/commercial", tempDir)
	err := os.MkdirAll(platformDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create YAML file (should be loaded)
	yamlFile := fmt.Sprintf("%s/valid_platform.yaml", platformDir)
	yamlContent := `
platform_types:
  valid_aircraft:
    class: "Valid Aircraft"
    category: "commercial_aircraft"
    performance:
      max_speed: 300.0
      cruise_speed: 250.0
      max_altitude: 11000.0
`
	err = os.WriteFile(yamlFile, []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create YAML file: %v", err)
	}

	// Create non-YAML file (should be ignored)
	txtFile := fmt.Sprintf("%s/ignored_file.txt", platformDir)
	err = os.WriteFile(txtFile, []byte("This should be ignored"), 0644)
	if err != nil {
		t.Fatalf("Failed to create text file: %v", err)
	}

	// Create JSON file (should be ignored)
	jsonFile := fmt.Sprintf("%s/ignored_file.json", platformDir)
	err = os.WriteFile(jsonFile, []byte(`{"ignored": true}`), 0644)
	if err != nil {
		t.Fatalf("Failed to create JSON file: %v", err)
	}

	// Change working directory
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tempDir)

	// Test loading platforms
	platforms, err := server.loadPlatformTypesFromFiles()

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Should only load from YAML files
	if len(platforms) == 0 {
		t.Error("Expected at least one platform to be loaded")
	}

	// Check that the YAML platform was loaded
	found := false
	for _, platform := range platforms {
		if platform.ID == "valid_aircraft" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected valid_aircraft platform to be loaded from YAML file")
	}
}
