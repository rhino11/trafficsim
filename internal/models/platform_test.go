package models

import (
	"math"
	"testing"
	"time"
)

func TestPosition(t *testing.T) {
	pos := Position{
		Latitude:  40.7128,
		Longitude: -74.0060,
		Altitude:  100.0,
	}

	if pos.Latitude != 40.7128 {
		t.Errorf("Expected latitude 40.7128, got %f", pos.Latitude)
	}
	if pos.Longitude != -74.0060 {
		t.Errorf("Expected longitude -74.0060, got %f", pos.Longitude)
	}
	if pos.Altitude != 100.0 {
		t.Errorf("Expected altitude 100.0, got %f", pos.Altitude)
	}
}

func TestVelocity(t *testing.T) {
	vel := Velocity{
		North: 100.0,
		East:  50.0,
		Up:    10.0,
	}

	if vel.North != 100.0 {
		t.Errorf("Expected north velocity 100.0, got %f", vel.North)
	}
	if vel.East != 50.0 {
		t.Errorf("Expected east velocity 50.0, got %f", vel.East)
	}
	if vel.Up != 10.0 {
		t.Errorf("Expected up velocity 10.0, got %f", vel.Up)
	}
}

func TestAcceleration(t *testing.T) {
	accel := Acceleration{
		North: 2.0,
		East:  1.0,
		Up:    0.5,
	}

	if accel.North != 2.0 {
		t.Errorf("Expected north acceleration 2.0, got %f", accel.North)
	}
	if accel.East != 1.0 {
		t.Errorf("Expected east acceleration 1.0, got %f", accel.East)
	}
	if accel.Up != 0.5 {
		t.Errorf("Expected up acceleration 0.5, got %f", accel.Up)
	}
}

func TestAttitude(t *testing.T) {
	attitude := Attitude{
		Roll:  15.0,
		Pitch: 10.0,
		Yaw:   270.0,
	}

	if attitude.Roll != 15.0 {
		t.Errorf("Expected roll 15.0, got %f", attitude.Roll)
	}
	if attitude.Pitch != 10.0 {
		t.Errorf("Expected pitch 10.0, got %f", attitude.Pitch)
	}
	if attitude.Yaw != 270.0 {
		t.Errorf("Expected yaw 270.0, got %f", attitude.Yaw)
	}
}

func TestAngularVelocity(t *testing.T) {
	angVel := AngularVelocity{
		RollRate:  5.0,
		PitchRate: 3.0,
		YawRate:   2.0,
	}

	if angVel.RollRate != 5.0 {
		t.Errorf("Expected roll rate 5.0, got %f", angVel.RollRate)
	}
	if angVel.PitchRate != 3.0 {
		t.Errorf("Expected pitch rate 3.0, got %f", angVel.PitchRate)
	}
	if angVel.YawRate != 2.0 {
		t.Errorf("Expected yaw rate 2.0, got %f", angVel.YawRate)
	}
}

func TestForces(t *testing.T) {
	forces := Forces{
		Thrust: 50000.0,
		Drag:   5000.0,
		Lift:   800000.0,
		Weight: 780000.0,
		Normal: 0.0,
	}

	if forces.Thrust != 50000.0 {
		t.Errorf("Expected thrust 50000.0, got %f", forces.Thrust)
	}
	if forces.Drag != 5000.0 {
		t.Errorf("Expected drag 5000.0, got %f", forces.Drag)
	}
	if forces.Lift != 800000.0 {
		t.Errorf("Expected lift 800000.0, got %f", forces.Lift)
	}
	if forces.Weight != 780000.0 {
		t.Errorf("Expected weight 780000.0, got %f", forces.Weight)
	}
}

func TestMomentOfInertia(t *testing.T) {
	moi := MomentOfInertia{
		Ixx: 1000000.0,
		Iyy: 2000000.0,
		Izz: 1500000.0,
	}

	if moi.Ixx != 1000000.0 {
		t.Errorf("Expected Ixx 1000000.0, got %f", moi.Ixx)
	}
	if moi.Iyy != 2000000.0 {
		t.Errorf("Expected Iyy 2000000.0, got %f", moi.Iyy)
	}
	if moi.Izz != 1500000.0 {
		t.Errorf("Expected Izz 1500000.0, got %f", moi.Izz)
	}
}

func TestUniversalPlatformBasicOperations(t *testing.T) {
	pos := Position{Latitude: 40.0, Longitude: -74.0, Altitude: 0}

	platform := &UniversalPlatform{
		ID:           "TEST-001",
		PlatformType: PlatformTypeAirborne,
		TypeDef: &PlatformTypeDefinition{
			Class:    "Test Aircraft",
			Category: "test",
			Performance: PerformanceCharacteristics{
				MaxSpeed:    100.0,
				CruiseSpeed: 80.0,
				MaxAltitude: 10000.0,
			},
			Physical: PhysicalCharacteristics{
				Length: 20.0,
				Width:  15.0,
				Height: 5.0,
				Mass:   5000.0,
			},
		},
		Config: &PlatformConfiguration{
			ID:   "TEST-001",
			Name: "Test Platform",
		},
		State: PlatformState{
			ID:       "TEST-001",
			Position: pos,
			Velocity: Velocity{},
			Heading:  0.0,
			Speed:    0.0,
		},
		CallSign: "TEST001",
	}

	// Test basic getters
	if platform.GetID() != "TEST-001" {
		t.Errorf("Expected ID TEST-001, got %s", platform.GetID())
	}
	if platform.GetType() != PlatformTypeAirborne {
		t.Errorf("Expected type airborne, got %s", platform.GetType())
	}
	if platform.GetClass() != "Test Aircraft" {
		t.Errorf("Expected class Test Aircraft, got %s", platform.GetClass())
	}
	if platform.GetMaxSpeed() != 100.0 {
		t.Errorf("Expected max speed 100.0, got %f", platform.GetMaxSpeed())
	}
	if platform.GetMaxAltitude() != 10000.0 {
		t.Errorf("Expected max altitude 10000.0, got %f", platform.GetMaxAltitude())
	}

	// Test destination setting
	dest := Position{Latitude: 41.0, Longitude: -73.0, Altitude: 1000.0}
	if err := platform.SetDestination(dest); err != nil {
		t.Errorf("Unexpected error setting destination: %v", err)
	}

	if platform.Destination == nil {
		t.Error("Destination should not be nil after setting")
	}
	if platform.Destination.Latitude != 41.0 {
		t.Errorf("Expected destination latitude 41.0, got %f", platform.Destination.Latitude)
	}
}

func TestUniversalPlatformCalculations(t *testing.T) {
	platform := &UniversalPlatform{
		State: PlatformState{
			Position: Position{Latitude: 40.0, Longitude: -74.0, Altitude: 0},
		},
	}

	// Test great circle distance calculation
	target := Position{Latitude: 41.0, Longitude: -74.0, Altitude: 0}
	distance := platform.calculateGreatCircleDistance(target)

	// Should be approximately 111 km (1 degree of latitude)
	expectedDistance := 111320.0 // meters
	tolerance := 1000.0          // 1km tolerance

	if math.Abs(distance-expectedDistance) > tolerance {
		t.Errorf("Expected distance ~%f, got %f", expectedDistance, distance)
	}

	// Test bearing calculation
	bearing := platform.calculateBearing(target)
	expectedBearing := 0.0 // Due north

	if math.Abs(bearing-expectedBearing) > 1.0 {
		t.Errorf("Expected bearing ~%f, got %f", expectedBearing, bearing)
	}

	// Test horizontal distance calculation
	horizontalDist := platform.calculateHorizontalDistance(target)
	if math.Abs(horizontalDist-expectedDistance) > tolerance {
		t.Errorf("Expected horizontal distance ~%f, got %f", expectedDistance, horizontalDist)
	}
}

func TestUniversalPlatformAccelerationConstraints(t *testing.T) {
	platform := &UniversalPlatform{
		TypeDef: &PlatformTypeDefinition{
			Performance: PerformanceCharacteristics{
				MaxSpeed:     100.0,
				Acceleration: 2.0, // 2 m/s²
			},
		},
		State: PlatformState{
			Speed: 50.0,
		},
	}

	// Test acceleration within limits
	targetSpeed := 55.0
	deltaTime := 2.0 // 2 seconds
	platform.applyAccelerationConstraints(targetSpeed, deltaTime)

	expectedSpeed := 54.0 // 50 + (2 * 2) = 54, which is less than target of 55
	if platform.State.Speed != expectedSpeed {
		t.Errorf("Expected speed %f, got %f", expectedSpeed, platform.State.Speed)
	}

	// Test acceleration exceeding max speed
	platform.State.Speed = 98.0
	targetSpeed = 105.0
	platform.applyAccelerationConstraints(targetSpeed, deltaTime)

	if platform.State.Speed != 100.0 { // Should be capped at max speed
		t.Errorf("Expected speed capped at 100.0, got %f", platform.State.Speed)
	}

	// Test deceleration
	platform.State.Speed = 60.0
	targetSpeed = 50.0
	platform.applyAccelerationConstraints(targetSpeed, deltaTime)

	expectedSpeed = 56.0 // 60 - (2 * 2) = 56, which is greater than target of 50
	if platform.State.Speed != expectedSpeed {
		t.Errorf("Expected speed %f, got %f", expectedSpeed, platform.State.Speed)
	}
}

func TestUniversalPlatformTurningConstraints(t *testing.T) {
	platform := &UniversalPlatform{
		TypeDef: &PlatformTypeDefinition{
			Performance: PerformanceCharacteristics{
				TurningRadius: 100.0, // 100 meter turning radius
			},
		},
		State: PlatformState{
			Heading: 0.0,
			Speed:   20.0, // 20 m/s
		},
	}

	// Test turning constraints
	desiredHeading := 90.0 // Turn to east
	deltaTime := 1.0       // 1 second

	platform.applyTurningConstraints(desiredHeading, deltaTime)

	// Calculate expected turn rate: speed / turning_radius * 180/π
	expectedTurnRate := (20.0 / 100.0) * 180.0 / math.Pi // ~11.46 degrees/second
	expectedHeading := 0.0 + expectedTurnRate*deltaTime

	tolerance := 0.1
	if math.Abs(platform.State.Heading-expectedHeading) > tolerance {
		t.Errorf("Expected heading ~%f, got %f", expectedHeading, platform.State.Heading)
	}
}

func TestUniversalPlatformUpdate(t *testing.T) {
	platform := &UniversalPlatform{
		ID:           "TEST-002",
		PlatformType: PlatformTypeAirborne,
		TypeDef: &PlatformTypeDefinition{
			Performance: PerformanceCharacteristics{
				MaxSpeed:        100.0,
				CruiseSpeed:     80.0,
				Acceleration:    2.0,
				FuelConsumption: 0.1, // kg/s
			},
		},
		State: PlatformState{
			Position: Position{Latitude: 40.0, Longitude: -74.0, Altitude: 1000.0},
			Speed:    0.0,
		},
		FuelRemaining: 1000.0, // kg
	}

	// Set a destination
	dest := Position{Latitude: 40.01, Longitude: -74.0, Altitude: 1000.0}
	if err := platform.SetDestination(dest); err != nil {
		t.Errorf("Unexpected error setting destination: %v", err)
	}

	// Update for 10 seconds
	deltaTime := 10 * time.Second
	if err := platform.Update(deltaTime); err != nil {
		t.Errorf("Unexpected error during update: %v", err)
	}

	// Check that fuel was consumed
	expectedFuelRemaining := 1000.0 - (0.1 * 10.0) // 999.0 kg
	if platform.FuelRemaining != expectedFuelRemaining {
		t.Errorf("Expected fuel remaining %f, got %f", expectedFuelRemaining, platform.FuelRemaining)
	}

	// Check that mission time was updated
	expectedMissionTime := deltaTime
	if platform.MissionTime != expectedMissionTime {
		t.Errorf("Expected mission time %v, got %v", expectedMissionTime, platform.MissionTime)
	}

	// Check that platform moved towards destination
	if platform.State.Position.Latitude <= 40.0 {
		t.Error("Platform should have moved north towards destination")
	}
}

func TestPerformanceCharacteristicAccess(t *testing.T) {
	platform := &UniversalPlatform{
		TypeDef: &PlatformTypeDefinition{
			Performance: PerformanceCharacteristics{
				CruiseSpeed:     80.0,
				FuelConsumption: 0.5,
				TurningRadius:   150.0,
				Acceleration:    1.5,
				ClimbRate:       12.0,
				OrbitalPeriod:   5400.0,
			},
		},
	}

	// Test valid characteristics
	testCases := []struct {
		name     string
		expected float64
	}{
		{"cruise_speed", 80.0},
		{"fuel_consumption", 0.5},
		{"turning_radius", 150.0},
		{"acceleration", 1.5},
		{"climb_rate", 12.0},
		{"orbital_period", 5400.0},
	}

	for _, tc := range testCases {
		value, err := platform.GetPerformanceCharacteristic(tc.name)
		if err != nil {
			t.Errorf("Unexpected error for %s: %v", tc.name, err)
		}
		if value != tc.expected {
			t.Errorf("For %s, expected %f, got %f", tc.name, tc.expected, value)
		}
	}

	// Test invalid characteristic
	_, err := platform.GetPerformanceCharacteristic("invalid_param")
	if err == nil {
		t.Error("Expected error for invalid performance characteristic")
	}
}

// Benchmark tests for performance-critical operations
func BenchmarkCalculateGreatCircleDistance(b *testing.B) {
	platform := &UniversalPlatform{
		State: PlatformState{
			Position: Position{Latitude: 40.0, Longitude: -74.0, Altitude: 0},
		},
	}
	target := Position{Latitude: 41.0, Longitude: -73.0, Altitude: 1000.0}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		platform.calculateGreatCircleDistance(target)
	}
}

func BenchmarkCalculateBearing(b *testing.B) {
	platform := &UniversalPlatform{
		State: PlatformState{
			Position: Position{Latitude: 40.0, Longitude: -74.0, Altitude: 0},
		},
	}
	target := Position{Latitude: 41.0, Longitude: -73.0, Altitude: 1000.0}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		platform.calculateBearing(target)
	}
}

func BenchmarkUniversalPlatformUpdate(b *testing.B) {
	platform := &UniversalPlatform{
		ID:           "BENCH-001",
		PlatformType: PlatformTypeAirborne,
		TypeDef: &PlatformTypeDefinition{
			Performance: PerformanceCharacteristics{
				MaxSpeed:        100.0,
				CruiseSpeed:     80.0,
				Acceleration:    2.0,
				FuelConsumption: 0.1,
			},
		},
		State: PlatformState{
			Position: Position{Latitude: 40.0, Longitude: -74.0, Altitude: 1000.0},
			Speed:    50.0,
		},
		FuelRemaining: 1000.0,
	}

	dest := Position{Latitude: 40.1, Longitude: -74.0, Altitude: 1000.0}
	if err := platform.SetDestination(dest); err != nil {
		b.Fatalf("Error setting destination: %v", err)
	}

	deltaTime := 1 * time.Second
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := platform.Update(deltaTime); err != nil {
			b.Errorf("Update failed: %v", err)
		}
	}
}
