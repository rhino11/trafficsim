package sim

import (
	"testing"
	"time"

	"github.com/rhino11/trafficsim/internal/models"
)

func TestNewPhysicsEngine(t *testing.T) {
	pe := NewPhysicsEngine()

	if pe == nil {
		t.Fatal("NewPhysicsEngine returned nil")
	}

	// Check default values
	if pe.EarthRadius != 6371000.0 {
		t.Errorf("Expected EarthRadius 6371000.0, got %f", pe.EarthRadius)
	}

	if pe.GravityAccel != 9.81 {
		t.Errorf("Expected GravityAccel 9.81, got %f", pe.GravityAccel)
	}

	if pe.AirDensity != 1.225 {
		t.Errorf("Expected AirDensity 1.225, got %f", pe.AirDensity)
	}

	if pe.TimeStep != time.Second {
		t.Errorf("Expected TimeStep 1s, got %v", pe.TimeStep)
	}
}

func TestCalculateMovement(t *testing.T) {
	pe := NewPhysicsEngine()

	// Test with UniversalPlatform
	platform := &models.UniversalPlatform{
		ID:           "test-platform",
		PlatformType: models.PlatformTypeAirborne,
		State: models.PlatformState{
			Position: models.Position{
				Latitude:  40.7128,
				Longitude: -74.0060,
				Altitude:  1000,
			},
			Speed:   100,
			Heading: 90,
		},
		TypeDef: &models.PlatformTypeDefinition{
			Performance: models.PerformanceCharacteristics{
				MaxSpeed:      300,
				CruiseSpeed:   250,
				ClimbRate:     10,
				TurningRadius: 500,
				Acceleration:  2.0,
			},
		},
		Destination: &models.Position{
			Latitude:  40.7580,
			Longitude: -73.9855,
			Altitude:  2000,
		},
	}

	err := pe.CalculateMovement(platform, time.Second)
	if err != nil {
		t.Errorf("CalculateMovement failed: %v", err)
	}
}

func TestCalculateGreatCircleDistance(t *testing.T) {
	pe := NewPhysicsEngine()

	pos1 := models.Position{Latitude: 40.7128, Longitude: -74.0060, Altitude: 0}
	pos2 := models.Position{Latitude: 40.7580, Longitude: -73.9855, Altitude: 0}

	distance := pe.CalculateGreatCircleDistance(pos1, pos2)

	// Distance between these NYC coordinates should be roughly 5.5km
	if distance < 5000 || distance > 6000 {
		t.Errorf("Expected distance around 5500m, got %f", distance)
	}
}

func TestCalculateBearing(t *testing.T) {
	pe := NewPhysicsEngine()

	pos1 := models.Position{Latitude: 40.7128, Longitude: -74.0060, Altitude: 0}
	pos2 := models.Position{Latitude: 40.7580, Longitude: -73.9855, Altitude: 0}

	bearing := pe.CalculateBearing(pos1, pos2)

	// Should be roughly northeast
	if bearing < 0 || bearing > 360 {
		t.Errorf("Bearing should be between 0-360 degrees, got %f", bearing)
	}
}

func TestUpdateAircraftPhysics(t *testing.T) {
	pe := NewPhysicsEngine()

	platform := &models.UniversalPlatform{
		ID:           "aircraft-test",
		PlatformType: models.PlatformTypeAirborne,
		State: models.PlatformState{
			Position: models.Position{
				Latitude:  40.7128,
				Longitude: -74.0060,
				Altitude:  1000,
			},
			Speed:   100,
			Heading: 90,
		},
		TypeDef: &models.PlatformTypeDefinition{
			Performance: models.PerformanceCharacteristics{
				MaxSpeed:      300,
				CruiseSpeed:   250,
				ClimbRate:     10,
				TurningRadius: 500,
				Acceleration:  2.0,
			},
		},
		Destination: &models.Position{
			Latitude:  40.7580,
			Longitude: -73.9855,
			Altitude:  2000,
		},
	}

	initialAltitude := platform.State.Position.Altitude
	err := pe.CalculateMovement(platform, time.Second)
	if err != nil {
		t.Errorf("Aircraft physics update failed: %v", err)
	}

	// Should have gained altitude
	if platform.State.Position.Altitude <= initialAltitude {
		t.Errorf("Expected altitude increase, got %f -> %f",
			initialAltitude, platform.State.Position.Altitude)
	}
}

func TestUpdateMaritimePhysics(t *testing.T) {
	pe := NewPhysicsEngine()

	platform := &models.UniversalPlatform{
		ID:           "ship-test",
		PlatformType: models.PlatformTypeMaritime,
		State: models.PlatformState{
			Position: models.Position{
				Latitude:  40.7128,
				Longitude: -74.0060,
				Altitude:  10, // Should be corrected to 0
			},
			Speed:   20,
			Heading: 180,
		},
		TypeDef: &models.PlatformTypeDefinition{
			Performance: models.PerformanceCharacteristics{
				CruiseSpeed:   30,
				TurningRadius: 200,
				Acceleration:  0.5,
			},
			Physical: models.PhysicalCharacteristics{
				Length: 100,
			},
		},
		Destination: &models.Position{
			Latitude:  40.6892,
			Longitude: -74.0445,
			Altitude:  0,
		},
	}

	err := pe.CalculateMovement(platform, time.Second)
	if err != nil {
		t.Errorf("Maritime physics update failed: %v", err)
	}

	// Should be at sea level
	if platform.State.Position.Altitude != 0 {
		t.Errorf("Maritime platform should be at sea level, got %f",
			platform.State.Position.Altitude)
	}
}

func TestUpdateLandPhysics(t *testing.T) {
	pe := NewPhysicsEngine()

	platform := &models.UniversalPlatform{
		ID:           "vehicle-test",
		PlatformType: models.PlatformTypeLand,
		State: models.PlatformState{
			Position: models.Position{
				Latitude:  40.7128,
				Longitude: -74.0060,
				Altitude:  100,
			},
			Speed:   50,
			Heading: 270,
		},
		TypeDef: &models.PlatformTypeDefinition{
			Performance: models.PerformanceCharacteristics{
				CruiseSpeed:   80,
				TurningRadius: 50,
				Acceleration:  3.0,
			},
		},
		Destination: &models.Position{
			Latitude:  40.7300,
			Longitude: -74.0200,
			Altitude:  150,
		},
	}

	err := pe.CalculateMovement(platform, time.Second)
	if err != nil {
		t.Errorf("Land physics update failed: %v", err)
	}
}

func TestUpdateSpacePhysics(t *testing.T) {
	pe := NewPhysicsEngine()

	platform := &models.UniversalPlatform{
		ID:           "satellite-test",
		PlatformType: models.PlatformTypeSpace,
		State: models.PlatformState{
			Position: models.Position{
				Latitude:  0,
				Longitude: 0,
				Altitude:  400000, // 400km orbit
			},
			Speed:   7800, // Orbital velocity
			Heading: 0,
		},
		TypeDef: &models.PlatformTypeDefinition{
			Performance: models.PerformanceCharacteristics{
				OrbitalVelocity: 7800,
			},
		},
	}

	err := pe.CalculateMovement(platform, time.Second)
	if err != nil {
		t.Errorf("Space physics update failed: %v", err)
	}

	// Should maintain orbital velocity
	if platform.State.Speed != 7800 {
		t.Errorf("Expected orbital velocity 7800, got %f", platform.State.Speed)
	}
}

func TestPlatformWithoutDestination(t *testing.T) {
	pe := NewPhysicsEngine()

	platform := &models.UniversalPlatform{
		ID:           "stationary-test",
		PlatformType: models.PlatformTypeAirborne,
		State: models.PlatformState{
			Position: models.Position{
				Latitude:  40.7128,
				Longitude: -74.0060,
				Altitude:  1000,
			},
			Speed:   100,
			Heading: 90,
		},
		TypeDef: &models.PlatformTypeDefinition{
			Performance: models.PerformanceCharacteristics{
				MaxSpeed:    300,
				CruiseSpeed: 250,
			},
		},
		Destination: nil, // No destination
	}

	originalPosition := platform.State.Position
	err := pe.CalculateMovement(platform, time.Second)
	if err != nil {
		t.Errorf("Movement without destination failed: %v", err)
	}

	// Position should remain unchanged
	if platform.State.Position != originalPosition {
		t.Errorf("Position changed without destination")
	}
}

func TestPhysicsEngineConfiguration(t *testing.T) {
	pe := NewPhysicsEngine()

	// Test configuration changes
	pe.EnableWeather = true
	pe.EnableTerrain = true
	pe.TimeStep = 500 * time.Millisecond

	if !pe.EnableWeather {
		t.Error("Weather should be enabled")
	}

	if !pe.EnableTerrain {
		t.Error("Terrain should be enabled")
	}

	if pe.TimeStep != 500*time.Millisecond {
		t.Errorf("Expected TimeStep 500ms, got %v", pe.TimeStep)
	}
}
