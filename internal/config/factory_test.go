package config

import (
	"testing"

	"github.com/rhino11/trafficsim/internal/models"
)

func TestNewPlatformFactory(t *testing.T) {
	registry := &PlatformRegistry{
		AirborneTypes: map[string]PlatformTypeDefinition{
			"boeing_737": {
				Class:       "Boeing 737-800",
				Category:    "commercial_aircraft",
				Type:        "airborne",
				MaxSpeed:    257.0,
				CruiseSpeed: 230.0,
				Length:      39.5,
				Width:       35.8,
				Mass:        79016.0,
			},
		},
	}

	factory := NewPlatformFactory(registry)
	if factory == nil {
		t.Error("NewPlatformFactory should not return nil")
	}

	if factory.registry != registry {
		t.Error("Factory should store the provided registry")
	}
}

func TestPlatformFactory_CreatePlatform(t *testing.T) {
	registry := createTestRegistry()
	factory := NewPlatformFactory(registry)

	instance := PlatformInstance{
		ID:     "test-1",
		TypeID: "f16_fighter",
		Name:   "Test Fighter",
		StartPos: Position{
			Latitude:  34.0522,
			Longitude: -118.2437,
			Altitude:  10668,
		},
		Route: []Position{
			{Latitude: 34.0522, Longitude: -118.2437, Altitude: 10668},
			{Latitude: 40.7128, Longitude: -74.0060, Altitude: 10668},
		},
	}

	platform, err := factory.CreatePlatform(instance)
	if err != nil {
		t.Fatalf("Failed to create platform: %v", err)
	}

	if platform.GetID() != "test-1" {
		t.Errorf("Expected ID 'test-1', got '%s'", platform.GetID())
	}

	// Test with invalid type ID
	invalidInstance := PlatformInstance{
		ID:     "invalid-test",
		TypeID: "nonexistent",
	}

	_, err = factory.CreatePlatform(invalidInstance)
	if err == nil {
		t.Error("Expected error for invalid type ID")
	}
}

func TestPlatformFactory_CreatePlatformInvalidType(t *testing.T) {
	registry := &PlatformRegistry{
		AirborneTypes: map[string]PlatformTypeDefinition{
			"boeing_737": {Class: "Boeing 737-800", Type: "airborne"},
		},
	}

	factory := NewPlatformFactory(registry)

	instance := PlatformInstance{
		ID:     "test-aircraft-1",
		TypeID: "boeing_737",
		Name:   "Test Aircraft",
		StartPos: Position{
			Latitude:  40.0,
			Longitude: -74.0,
			Altitude:  10000.0,
		},
	}

	platform, err := factory.CreatePlatform(instance)
	if err != nil {
		t.Fatalf("Failed to create platform: %v", err)
	}

	if platform == nil {
		t.Fatal("Expected non-nil platform")
	}

	if platform.GetID() != "test-aircraft-1" {
		t.Errorf("Expected ID 'test-aircraft-1', got '%s'", platform.GetID())
	}

	// Test with invalid type ID
	invalidInstance := PlatformInstance{
		ID:     "invalid-test",
		TypeID: "nonexistent",
	}

	_, err = factory.CreatePlatform(invalidInstance)
	if err == nil {
		t.Error("Expected error for invalid type ID")
	}
}

func TestPlatformFactory_CreateScenario(t *testing.T) {
	registry := createTestRegistry()

	// Add a test scenario to the registry with platforms that exist
	registry.Scenarios["test_scenario"] = ScenarioConfig{
		Name:        "Test Scenario",
		Description: "Test scenario for factory",
		Instances: []PlatformInstance{
			{
				ID:     "test-fighter-1",
				TypeID: "f16_fighter",
				Name:   "Test Fighter 1",
				StartPos: Position{
					Latitude:  34.0522,
					Longitude: -118.2437,
					Altitude:  10668,
				},
				Route: []Position{
					{Latitude: 34.0522, Longitude: -118.2437, Altitude: 10668},
					{Latitude: 40.7128, Longitude: -74.0060, Altitude: 10668},
				},
			},
			{
				ID:     "test-carrier-1",
				TypeID: "aircraft_carrier",
				Name:   "Test Carrier 1",
				StartPos: Position{
					Latitude:  35.0,
					Longitude: -120.0,
					Altitude:  0,
				},
			},
		},
	}

	factory := NewPlatformFactory(registry)

	// Test creating all platforms in a scenario
	platforms, err := factory.CreateScenario("test_scenario")
	if err != nil {
		t.Fatalf("Failed to create scenario: %v", err)
	}

	if len(platforms) != 2 {
		t.Errorf("Expected 2 platforms, got %d", len(platforms))
	}

	// Verify platform IDs
	expectedIDs := map[string]bool{"test-fighter-1": false, "test-carrier-1": false}
	for _, platform := range platforms {
		id := platform.GetID()
		if _, exists := expectedIDs[id]; exists {
			expectedIDs[id] = true
		} else {
			t.Errorf("Unexpected platform ID: %s", id)
		}
	}

	for id, found := range expectedIDs {
		if !found {
			t.Errorf("Expected platform ID not found: %s", id)
		}
	}

	// Test with nonexistent scenario
	_, err = factory.CreateScenario("nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent scenario")
	}
}

func TestPlatformFactory_GetAvailablePlatformTypes(t *testing.T) {
	registry := &PlatformRegistry{
		AirborneTypes: map[string]PlatformTypeDefinition{
			"boeing_737": {Class: "Boeing 737-800", Type: "airborne"},
			"f16":        {Class: "F-16 Fighting Falcon", Type: "airborne"},
		},
		MaritimeTypes: map[string]PlatformTypeDefinition{
			"container_ship": {Class: "Container Ship", Type: "maritime"},
		},
		LandTypes: map[string]PlatformTypeDefinition{
			"humvee": {Class: "HMMWV", Type: "land"},
		},
	}

	factory := NewPlatformFactory(registry)
	types := factory.GetAvailablePlatformTypes()

	// Check airborne types
	airborneTypes, exists := types["airborne"]
	if !exists {
		t.Error("Expected airborne types to be present")
	} else if len(airborneTypes) != 2 {
		t.Errorf("Expected 2 airborne types, got %d", len(airborneTypes))
	}

	// Check maritime types
	maritimeTypes, exists := types["maritime"]
	if !exists {
		t.Error("Expected maritime types to be present")
	} else if len(maritimeTypes) != 1 {
		t.Errorf("Expected 1 maritime type, got %d", len(maritimeTypes))
	}

	// Check land types
	landTypes, exists := types["land"]
	if !exists {
		t.Error("Expected land types to be present")
	} else if len(landTypes) != 1 {
		t.Errorf("Expected 1 land type, got %d", len(landTypes))
	}

	// Verify specific type IDs are present
	expectedAirborne := map[string]bool{"boeing_737": false, "f16": false}
	for _, typeID := range airborneTypes {
		if _, exists := expectedAirborne[typeID]; exists {
			expectedAirborne[typeID] = true
		}
	}
	for typeID, found := range expectedAirborne {
		if !found {
			t.Errorf("Expected airborne type %s not found", typeID)
		}
	}
}

func TestPlatformFactory_ValidateScenario(t *testing.T) {
	registry := &PlatformRegistry{
		AirborneTypes: map[string]PlatformTypeDefinition{
			"boeing_737": {Class: "Boeing 737-800", Type: "airborne"},
		},
		Scenarios: map[string]ScenarioConfig{
			"valid_scenario": {
				Name: "Valid Scenario",
				Instances: []PlatformInstance{
					{
						ID:     "test-1",
						TypeID: "boeing_737",
					},
				},
			},
			"invalid_scenario": {
				Name: "Invalid Scenario",
				Instances: []PlatformInstance{
					{
						ID:     "test-1",
						TypeID: "unknown_type",
					},
				},
			},
		},
	}

	factory := NewPlatformFactory(registry)

	// Valid scenario
	err := factory.ValidateScenario("valid_scenario")
	if err != nil {
		t.Errorf("Valid scenario should not produce error: %v", err)
	}

	// Invalid scenario - unknown type
	err = factory.ValidateScenario("invalid_scenario")
	if err == nil {
		t.Error("Invalid scenario should produce error")
	}

	// Nonexistent scenario
	err = factory.ValidateScenario("nonexistent_scenario")
	if err == nil {
		t.Error("Nonexistent scenario should produce error")
	}
}

func TestPlatformFactory_generateCallSign(t *testing.T) {
	factory := &PlatformFactory{}

	// Test with call sign prefix and format
	typeDef := PlatformTypeDefinition{
		CallSignPrefix: "UAL",
		CallSignFormat: "{prefix}-{id}",
	}
	callSign := factory.generateCallSign(&typeDef, "1234")
	expected := "UAL-1234"
	if callSign != expected {
		t.Errorf("Expected call sign '%s', got '%s'", expected, callSign)
	}

	// Test with call sign prefix but no format
	typeDef = PlatformTypeDefinition{
		CallSignPrefix: "UAL",
	}
	callSign = factory.generateCallSign(&typeDef, "1234567")
	expected = "UAL567" // Should use last 3 chars
	if callSign != expected {
		t.Errorf("Expected call sign '%s', got '%s'", expected, callSign)
	}

	// Test with short ID
	callSign = factory.generateCallSign(&typeDef, "12")
	expected = "UAL12"
	if callSign != expected {
		t.Errorf("Expected call sign '%s', got '%s'", expected, callSign)
	}

	// Test with no call sign prefix
	typeDef = PlatformTypeDefinition{}
	callSign = factory.generateCallSign(&typeDef, "test-id")
	expected = "test-id"
	if callSign != expected {
		t.Errorf("Expected call sign '%s', got '%s'", expected, callSign)
	}
}

func TestPlatformFactory_determinePlatformType(t *testing.T) {
	factory := &PlatformFactory{}

	testCases := []struct {
		input    string
		expected models.PlatformType
	}{
		{"airborne", models.PlatformTypeAirborne},
		{"maritime", models.PlatformTypeMaritime},
		{"land", models.PlatformTypeLand},
		{"space", models.PlatformTypeSpace},
		{"unknown", models.PlatformTypeAirborne}, // Default fallback
	}

	for _, tc := range testCases {
		result := factory.determinePlatformType(tc.input)
		if result != tc.expected {
			t.Errorf("Input %s: expected %v, got %v", tc.input, tc.expected, result)
		}
	}
}

// createTestRegistry creates a test registry with sample platform types
func createTestRegistry() *PlatformRegistry {
	return &PlatformRegistry{
		AirborneTypes: map[string]PlatformTypeDefinition{
			"f16_fighter": {
				Class:       "F-16 Fighting Falcon",
				Category:    "military_aircraft",
				Type:        "airborne",
				MaxSpeed:    686.0,
				CruiseSpeed: 577.0,
				MaxAltitude: 15240.0,
				Length:      14.8,
				Width:       9.8,
				Mass:        8573.0,
			},
		},
		MaritimeTypes: map[string]PlatformTypeDefinition{
			"aircraft_carrier": {
				Class:       "Nimitz-class",
				Category:    "military_ship",
				Type:        "maritime",
				MaxSpeed:    56.0,
				CruiseSpeed: 46.0,
				Length:      332.8,
				Width:       76.8,
				Mass:        104600000.0,
			},
		},
		LandTypes: map[string]PlatformTypeDefinition{
			"m1_abrams": {
				Class:       "M1 Abrams",
				Category:    "main_battle_tank",
				Type:        "land",
				MaxSpeed:    67.0,
				CruiseSpeed: 56.0,
				Length:      9.8,
				Width:       3.7,
				Mass:        62000.0,
			},
		},
		SpaceTypes: map[string]PlatformTypeDefinition{
			"satellite": {
				Class:       "Communications Satellite",
				Category:    "commercial_satellite",
				Type:        "space",
				MaxSpeed:    7800.0,
				CruiseSpeed: 7800.0,
				MaxAltitude: 35786000.0,
				Length:      4.0,
				Width:       2.0,
				Mass:        2000.0,
			},
		},
		Scenarios: make(map[string]ScenarioConfig),
	}
}

func TestPlatformFactory_CreateScenario_EmptyScenario(t *testing.T) {

	factory := &PlatformFactory{
		registry: &PlatformRegistry{
			Scenarios: make(map[string]ScenarioConfig),
		},
	}

	// Attempt to get platforms for a scenario that doesn't exist
	platforms, err := factory.CreateScenario("nonexistent_scenario")
	if err == nil {
		t.Error("Expected error for nonexistent scenario")
	} else if platforms != nil {
		t.Error("Expected nil platforms for nonexistent scenario")
	}

	// Add a valid scenario with no platforms
	factory.registry.Scenarios["empty_scenario"] = ScenarioConfig{
		Name: "Empty Scenario",
	}

	// Get platforms for the empty scenario
	platforms, err = factory.CreateScenario("empty_scenario")
	if err != nil {
		t.Fatalf("Unexpected error for empty scenario: %v", err)
	}

	if platforms != nil {
		t.Error("Expected nil platforms for empty scenario")
	}
}
