package config

import (
	"os"
	"testing"

	"github.com/rhino11/trafficsim/internal/testutil"
)

func TestLoadConfig(t *testing.T) {
	logger := testutil.SetupTestLogging(t)
	logger.Debug("Starting TestLoadConfig")

	// Create temporary config file in current directory (relative path)
	configFile := "test_config.yaml"
	configContent := `
simulation:
  update_interval: "16ms"
  time_scale: 1.5
server:
  host: "localhost"
  port: 8080
  web_root: "./web"
platforms:
  airborne_types:
    boeing_737:
      class: "Boeing 737-800"
      category: "commercial"
  maritime_types:
    destroyer:
      class: "Arleigh Burke"
      category: "military"
  land_types:
    tank:
      class: "M1A2 Abrams"
      category: "military"
  space_types:
    satellite:
      class: "GPS"
      category: "navigation"
  scenarios:
    test_scenario:
      instances:
        - type_id: "boeing_737"
          id: "test1"
          name: "Test Boeing"
          start_position:
            latitude: 40.0
            longitude: -74.0
            altitude: 10000
        - type_id: "destroyer"
          id: "test2"
          name: "Test Destroyer"
          start_position:
            latitude: 40.0
            longitude: -74.0
            altitude: 0
`

	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(configFile)

	// Test loading valid config
	config, err := LoadConfig(configFile)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Verify simulation config
	if config.Simulation.UpdateInterval != "16ms" {
		t.Errorf("Expected update_interval '16ms', got '%s'", config.Simulation.UpdateInterval)
	}
	if config.Simulation.TimeScale != 1.5 {
		t.Errorf("Expected time_scale 1.5, got %f", config.Simulation.TimeScale)
	}

	// Verify server config
	if config.Server.Host != "localhost" {
		t.Errorf("Expected host 'localhost', got '%s'", config.Server.Host)
	}
	if config.Server.Port != 8080 {
		t.Errorf("Expected port 8080, got %d", config.Server.Port)
	}
	if config.Server.WebRoot != "./web" {
		t.Errorf("Expected web_root './web', got '%s'", config.Server.WebRoot)
	}

	// Verify platforms
	if len(config.Platforms.AirborneTypes) != 1 {
		t.Errorf("Expected 1 airborne type, got %d", len(config.Platforms.AirborneTypes))
	}
	if len(config.Platforms.MaritimeTypes) != 1 {
		t.Errorf("Expected 1 maritime type, got %d", len(config.Platforms.MaritimeTypes))
	}
	if len(config.Platforms.LandTypes) != 1 {
		t.Errorf("Expected 1 land type, got %d", len(config.Platforms.LandTypes))
	}
	if len(config.Platforms.SpaceTypes) != 1 {
		t.Errorf("Expected 1 space type, got %d", len(config.Platforms.SpaceTypes))
	}
	if len(config.Platforms.Scenarios) != 1 {
		t.Errorf("Expected 1 scenario, got %d", len(config.Platforms.Scenarios))
	}

	// Test boeing_737 platform
	boeing737, exists := config.Platforms.AirborneTypes["boeing_737"]
	if !exists {
		t.Error("Expected boeing_737 platform to exist")
	} else {
		if boeing737.Class != "Boeing 737-800" {
			t.Errorf("Expected class 'Boeing 737-800', got '%s'", boeing737.Class)
		}
		if boeing737.Category != "commercial" {
			t.Errorf("Expected category 'commercial', got '%s'", boeing737.Category)
		}
	}

	// Test scenario
	scenario, exists := config.Platforms.Scenarios["test_scenario"]
	if !exists {
		t.Error("Expected test_scenario to exist")
	} else {
		if len(scenario.Instances) != 2 {
			t.Errorf("Expected 2 instances, got %d", len(scenario.Instances))
		}
	}

	// Test loading non-existent file
	_, err = LoadConfig("non_existent_file.yaml")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}

	// Test loading invalid YAML
	invalidConfigFile := "invalid_test_config.yaml"
	invalidContent := "invalid: [unclosed"
	if err := os.WriteFile(invalidConfigFile, []byte(invalidContent), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(invalidConfigFile)
	_, err = LoadConfig(invalidConfigFile)
	if err == nil {
		t.Error("Expected error for invalid YAML")
	}
}

func TestConfigDefaults(t *testing.T) {
	logger := testutil.SetupTestLogging(t)
	logger.Debug("Starting TestConfigDefaults")

	// Create minimal config file in current directory (relative path)
	configFile := "minimal_test_config.yaml"
	minimalContent := `
platforms:
  airborne_types:
    test_aircraft:
      class: "Test Aircraft"
`

	if err := os.WriteFile(configFile, []byte(minimalContent), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(configFile)

	config, err := LoadConfig(configFile)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Check defaults are applied
	if config.Simulation.UpdateInterval == "" {
		t.Error("Expected default update_interval to be set")
	}
	if config.Simulation.TimeScale == 0 {
		t.Error("Expected default time_scale to be set")
	}
	if config.Server.Host == "" {
		t.Error("Expected default host to be set")
	}
	if config.Server.Port == 0 {
		t.Error("Expected default port to be set")
	}
}

func TestPlatformsHasType(t *testing.T) {
	logger := testutil.SetupTestLogging(t)
	logger.Debug("Starting TestPlatformsHasType")

	platforms := PlatformRegistry{
		AirborneTypes: map[string]PlatformTypeDefinition{
			"boeing_737": {Class: "Boeing 737"},
		},
		MaritimeTypes: map[string]PlatformTypeDefinition{
			"destroyer": {Class: "Destroyer"},
		},
		LandTypes: map[string]PlatformTypeDefinition{
			"tank": {Class: "Tank"},
		},
		SpaceTypes: map[string]PlatformTypeDefinition{
			"satellite": {Class: "Satellite"},
		},
	}

	// Test existing types
	if !platforms.HasType("boeing_737") {
		t.Error("Expected HasType to return true for boeing_737")
	}
	if !platforms.HasType("destroyer") {
		t.Error("Expected HasType to return true for destroyer")
	}
	if !platforms.HasType("tank") {
		t.Error("Expected HasType to return true for tank")
	}
	if !platforms.HasType("satellite") {
		t.Error("Expected HasType to return true for satellite")
	}

	// Test non-existing type
	if platforms.HasType("non_existent") {
		t.Error("Expected HasType to return false for non_existent")
	}
	if platforms.HasType("") {
		t.Error("Expected HasType to return false for empty string")
	}
}

func TestConfigValidation(t *testing.T) {
	logger := testutil.SetupTestLogging(t)
	logger.Debug("Starting TestConfigValidation")

	// Config with invalid port - using relative path
	configFile := "invalid_port_test_config.yaml"
	configContent := `
simulation:
  update_interval: "16ms"
  time_scale: 1.0
server:
  port: 99999
platforms:
  airborne_types:
    test:
      class: "Test"
`

	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(configFile)

	// Config should fail validation due to invalid port
	_, err := LoadConfig(configFile)
	if err == nil {
		t.Error("Expected error for invalid port")
	}
}
