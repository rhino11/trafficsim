package models

import (
	"math"
	"testing"
	"time"
)

func TestAirbornePlatformCreation(t *testing.T) {
	startPos := Position{Latitude: 40.7128, Longitude: -74.0060, Altitude: 0}

	// Test Boeing 737-800 creation
	boeing737 := NewBoeing737_800("B737-001", "AA1234", startPos)

	if boeing737.GetID() != "B737-001" {
		t.Errorf("Expected ID B737-001, got %s", boeing737.GetID())
	}
	if boeing737.GetClass() != "Boeing 737-800" {
		t.Errorf("Expected class Boeing 737-800, got %s", boeing737.GetClass())
	}
	if boeing737.GetName() != "AA1234" {
		t.Errorf("Expected name AA1234, got %s", boeing737.GetName())
	}
	if boeing737.FlightPhase != FlightPhaseParked {
		t.Errorf("Expected initial flight phase parked, got %s", boeing737.FlightPhase)
	}
	if boeing737.GetMaxSpeed() != 257 {
		t.Errorf("Expected max speed 257 m/s, got %f", boeing737.GetMaxSpeed())
	}
	if boeing737.MaxRollRate != 15 {
		t.Errorf("Expected max roll rate 15 deg/s, got %f", boeing737.MaxRollRate)
	}
	if boeing737.State.Physics.Mass != 79010.0 {
		t.Errorf("Expected mass 79010.0 kg, got %f", boeing737.State.Physics.Mass)
	}
}

func TestAirbusA320Creation(t *testing.T) {
	startPos := Position{Latitude: 40.7128, Longitude: -74.0060, Altitude: 0}

	airbus := NewAirbusA320("A320-001", "DL5678", startPos)

	if airbus.GetClass() != "Airbus A320" {
		t.Errorf("Expected class Airbus A320, got %s", airbus.GetClass())
	}
	if airbus.GetMaxSpeed() != 257 {
		t.Errorf("Expected max speed 257 m/s, got %f", airbus.GetMaxSpeed())
	}
	if airbus.GetLength() != 37.6 {
		t.Errorf("Expected length 37.6 m, got %f", airbus.GetLength())
	}
	if airbus.WingArea != 122.6 {
		t.Errorf("Expected wing area 122.6 m², got %f", airbus.WingArea)
	}
}

func TestF16Creation(t *testing.T) {
	startPos := Position{Latitude: 35.0, Longitude: -118.0, Altitude: 0}

	f16 := NewF16FightingFalcon("F16-001", "87-0001", startPos)

	if f16.GetClass() != "F-16 Fighting Falcon" {
		t.Errorf("Expected class F-16 Fighting Falcon, got %s", f16.GetClass())
	}
	if f16.MaxLoadFactor != 9.0 {
		t.Errorf("Expected max load factor 9.0 g, got %f", f16.MaxLoadFactor)
	}
	if f16.MaxRollRate != 720 {
		t.Errorf("Expected max roll rate 720 deg/s, got %f", f16.MaxRollRate)
	}
	if f16.MaxBankAngle != 90 {
		t.Errorf("Expected max bank angle 90 degrees, got %f", f16.MaxBankAngle)
	}
	if f16.CallSign != "VIPER001" {
		t.Errorf("Expected callsign VIPER001, got %s", f16.CallSign)
	}
}

func TestC130Creation(t *testing.T) {
	startPos := Position{Latitude: 32.0, Longitude: -106.0, Altitude: 0}

	c130 := NewC130Hercules("C130-001", "12-5678", startPos)

	if c130.GetClass() != "C-130 Hercules" {
		t.Errorf("Expected class C-130 Hercules, got %s", c130.GetClass())
	}
	if c130.MaxRollRate != 10 {
		t.Errorf("Expected max roll rate 10 deg/s, got %f", c130.MaxRollRate)
	}
	if c130.MaxBankAngle != 25 {
		t.Errorf("Expected max bank angle 25 degrees, got %f", c130.MaxBankAngle)
	}
	if c130.CallSign != "HERKY78" {
		t.Errorf("Expected callsign HERKY78, got %s", c130.CallSign)
	}
}

func TestFlightPhaseTransitions(t *testing.T) {
	startPos := Position{Latitude: 40.7128, Longitude: -74.0060, Altitude: 0}
	aircraft := NewBoeing737_800("TEST-001", "TEST123", startPos)

	// Test parked phase
	aircraft.updateFlightPhase()
	if aircraft.FlightPhase != FlightPhaseParked {
		t.Errorf("Expected parked phase, got %s", aircraft.FlightPhase)
	}

	// Test takeoff phase
	aircraft.State.Position.Altitude = 200
	aircraft.State.Velocity.Up = 5
	aircraft.State.Speed = 80
	aircraft.updateFlightPhase()
	if aircraft.FlightPhase != FlightPhaseTakeoff {
		t.Errorf("Expected takeoff phase, got %s", aircraft.FlightPhase)
	}

	// Test climb phase
	aircraft.State.Position.Altitude = 1000
	aircraft.State.Velocity.Up = 10
	aircraft.updateFlightPhase()
	if aircraft.FlightPhase != FlightPhaseClimb {
		t.Errorf("Expected climb phase, got %s", aircraft.FlightPhase)
	}

	// Test cruise phase
	aircraft.State.Position.Altitude = 11000
	aircraft.State.Velocity.Up = 0
	aircraft.updateFlightPhase()
	if aircraft.FlightPhase != FlightPhaseCruise {
		t.Errorf("Expected cruise phase, got %s", aircraft.FlightPhase)
	}

	// Test descent phase
	aircraft.State.Velocity.Up = -8
	aircraft.updateFlightPhase()
	if aircraft.FlightPhase != FlightPhaseDescent {
		t.Errorf("Expected descent phase, got %s", aircraft.FlightPhase)
	}
}

func TestAerodynamicForceCalculation(t *testing.T) {
	startPos := Position{Latitude: 40.0, Longitude: -74.0, Altitude: 10000}
	aircraft := NewBoeing737_800("TEST-002", "TEST456", startPos)

	// Set cruise conditions
	aircraft.State.Speed = 230               // cruise speed
	aircraft.State.Position.Altitude = 11000 // cruise altitude

	aircraft.calculateAerodynamicForces()

	// Check that forces are calculated
	if aircraft.State.Physics.Forces.Lift <= 0 {
		t.Error("Lift force should be positive")
	}
	if aircraft.State.Physics.Forces.Drag <= 0 {
		t.Error("Drag force should be positive")
	}
	if aircraft.State.Physics.Forces.Weight <= 0 {
		t.Error("Weight force should be positive")
	}
	if aircraft.State.Physics.Forces.Thrust <= 0 {
		t.Error("Thrust force should be positive")
	}

	// At cruise, lift should approximately equal weight
	liftToWeightRatio := aircraft.State.Physics.Forces.Lift / aircraft.State.Physics.Forces.Weight
	if math.Abs(liftToWeightRatio-1.0) > 0.2 { // Allow 20% tolerance
		t.Errorf("Lift/Weight ratio should be close to 1.0, got %f", liftToWeightRatio)
	}
}

func TestFlightDynamicsUpdate(t *testing.T) {
	startPos := Position{Latitude: 40.0, Longitude: -74.0, Altitude: 0}
	aircraft := NewBoeing737_800("TEST-003", "TEST789", startPos)

	// Test takeoff acceleration
	aircraft.FlightPhase = FlightPhaseTakeoff
	aircraft.State.Speed = 50

	aircraft.updateFlightDynamics(1.0) // 1 second

	// Speed should increase during takeoff
	if aircraft.State.Speed <= 50 {
		t.Error("Speed should increase during takeoff")
	}

	// Should not exceed max acceleration limits
	if aircraft.State.Speed > 50+aircraft.MaxAcceleration {
		t.Errorf("Speed increase should not exceed max acceleration: %f", aircraft.State.Speed)
	}

	// Test stall speed protection
	aircraft.State.Position.Altitude = 1000 // In flight
	aircraft.State.Speed = 50               // Below stall speed
	aircraft.updateFlightDynamics(1.0)

	if aircraft.State.Speed < aircraft.StallSpeed {
		t.Errorf("Speed should not fall below stall speed %f, got %f", aircraft.StallSpeed, aircraft.State.Speed)
	}
}

func TestAttitudeControl(t *testing.T) {
	startPos := Position{Latitude: 40.0, Longitude: -74.0, Altitude: 10000}
	aircraft := NewBoeing737_800("TEST-004", "TEST101", startPos)

	// Set up for a turn - aircraft initially pointing east, destination north
	dest := Position{Latitude: 40.1, Longitude: -74.0, Altitude: 10000}
	aircraft.Destination = &dest
	aircraft.State.Speed = 200
	aircraft.State.Heading = 90 // Initially pointing east, need to turn left to north

	initialRoll := aircraft.State.Physics.Attitude.Roll
	aircraft.updateAttitude(1.0) // 1 second

	// Aircraft should bank into the turn (left turn, so negative roll)
	if math.Abs(aircraft.State.Physics.Attitude.Roll) <= math.Abs(initialRoll) {
		t.Errorf("Aircraft should bank for the turn. Initial roll: %.3f, After: %.3f",
			initialRoll, aircraft.State.Physics.Attitude.Roll)
	}

	// Bank angle should not exceed limits
	if math.Abs(aircraft.State.Physics.Attitude.Roll) > aircraft.MaxBankAngle {
		t.Errorf("Bank angle %f should not exceed max %f",
			math.Abs(aircraft.State.Physics.Attitude.Roll), aircraft.MaxBankAngle)
	}

	// Test pitch angle for climb
	aircraft.State.Velocity.Up = 10 // Climbing
	aircraft.updateAttitude(1.0)

	if aircraft.State.Physics.Attitude.Pitch <= 0 {
		t.Error("Pitch should be positive when climbing")
	}

	if aircraft.State.Physics.Attitude.Pitch > aircraft.MaxPitchAngle {
		t.Errorf("Pitch angle %f should not exceed max %f",
			aircraft.State.Physics.Attitude.Pitch, aircraft.MaxPitchAngle)
	}
}

func TestBearingCalculation(t *testing.T) {
	startPos := Position{Latitude: 40.0, Longitude: -74.0, Altitude: 0}
	aircraft := NewBoeing737_800("BEARING-TEST", "TEST123", startPos)

	// Test bearing to point directly north
	northTarget := Position{Latitude: 41.0, Longitude: -74.0, Altitude: 0}
	bearing := aircraft.calculateBearing(aircraft.State.Position, northTarget)

	if math.Abs(bearing-0.0) > 1.0 { // 1 degree tolerance
		t.Errorf("Expected bearing ~0 degrees (north), got %f", bearing)
	}

	// Test bearing to point directly east
	eastTarget := Position{Latitude: 40.0, Longitude: -73.0, Altitude: 0}
	bearing = aircraft.calculateBearing(aircraft.State.Position, eastTarget)

	if math.Abs(bearing-90.0) > 1.0 {
		t.Errorf("Expected bearing ~90 degrees (east), got %f", bearing)
	}

	// Test bearing to point directly south
	southTarget := Position{Latitude: 39.0, Longitude: -74.0, Altitude: 0}
	bearing = aircraft.calculateBearing(aircraft.State.Position, southTarget)

	if math.Abs(bearing-180.0) > 1.0 {
		t.Errorf("Expected bearing ~180 degrees (south), got %f", bearing)
	}
}

func TestPositionUpdate(t *testing.T) {
	startPos := Position{Latitude: 40.0, Longitude: -74.0, Altitude: 10000}
	aircraft := NewBoeing737_800("TEST-005", "TEST202", startPos)

	// Set up flight conditions
	aircraft.State.Speed = 200  // m/s
	aircraft.State.Heading = 90 // Due east
	aircraft.State.Velocity.Up = 0

	initialLat := aircraft.State.Position.Latitude
	initialLon := aircraft.State.Position.Longitude

	aircraft.updatePosition(1.0) // 1 second

	// Should move east (longitude should increase)
	if aircraft.State.Position.Longitude <= initialLon {
		t.Error("Aircraft should move east (longitude should increase)")
	}

	// Latitude should remain approximately the same
	if math.Abs(aircraft.State.Position.Latitude-initialLat) > 0.001 {
		t.Error("Latitude should remain approximately constant when flying due east")
	}

	// Test altitude change
	aircraft.State.Velocity.Up = 10 // 10 m/s climb rate
	initialAlt := aircraft.State.Position.Altitude
	aircraft.updatePosition(1.0)

	expectedAlt := initialAlt + 10
	if math.Abs(aircraft.State.Position.Altitude-expectedAlt) > 0.1 {
		t.Errorf("Expected altitude %f, got %f", expectedAlt, aircraft.State.Position.Altitude)
	}
}

func TestAircraftUpdate(t *testing.T) {
	startPos := Position{Latitude: 40.0, Longitude: -74.0, Altitude: 0}
	aircraft := NewBoeing737_800("TEST-006", "TEST303", startPos)

	// Set destination
	dest := Position{Latitude: 40.1, Longitude: -74.0, Altitude: 11000}
	aircraft.SetDestination(dest)

	initialTime := aircraft.State.LastUpdated

	// Update for 10 seconds
	deltaTime := 10 * time.Second
	err := aircraft.Update(deltaTime)

	if err != nil {
		t.Errorf("Unexpected error during update: %v", err)
	}

	// Check that time was updated
	if !aircraft.State.LastUpdated.After(initialTime) {
		t.Error("LastUpdated should be updated after Update call")
	}

	// Check that aircraft moved towards destination
	if aircraft.State.Position.Latitude <= 40.0 {
		t.Error("Aircraft should move towards destination (north)")
	}

	// Check that physics state is updated
	if aircraft.State.Physics.Position.Latitude != aircraft.State.Position.Latitude {
		t.Error("Physics position should match state position")
	}
}

func TestPlatformInterfaceCompliance(t *testing.T) {
	startPos := Position{Latitude: 40.0, Longitude: -74.0, Altitude: 0}
	aircraft := NewBoeing737_800("TEST-007", "TEST404", startPos)

	// Test that AirbornePlatform implements Platform interface
	var platform Platform = aircraft

	if platform.GetID() != "TEST-007" {
		t.Errorf("Expected ID TEST-007, got %s", platform.GetID())
	}

	if platform.GetType() != PlatformTypeAirborne {
		t.Errorf("Expected type airborne, got %s", platform.GetType())
	}

	if platform.GetClass() != "Boeing 737-800" {
		t.Errorf("Expected class Boeing 737-800, got %s", platform.GetClass())
	}

	// Test state update
	newState := platform.GetState()
	newState.Speed = 250
	platform.UpdateState(newState)

	if platform.GetState().Speed != 250 {
		t.Errorf("Expected speed 250, got %f", platform.GetState().Speed)
	}
}

// Benchmark tests for airborne platform operations
func BenchmarkAirbornePlatformUpdate(b *testing.B) {
	startPos := Position{Latitude: 40.0, Longitude: -74.0, Altitude: 10000}
	aircraft := NewBoeing737_800("BENCH-001", "BENCH123", startPos)

	dest := Position{Latitude: 41.0, Longitude: -73.0, Altitude: 11000}
	aircraft.SetDestination(dest)

	deltaTime := 1 * time.Second

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		aircraft.Update(deltaTime)
	}
}

func BenchmarkAerodynamicCalculation(b *testing.B) {
	startPos := Position{Latitude: 40.0, Longitude: -74.0, Altitude: 10000}
	aircraft := NewBoeing737_800("BENCH-002", "BENCH456", startPos)
	aircraft.State.Speed = 230

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		aircraft.calculateAerodynamicForces()
	}
}

func BenchmarkAttitudeUpdate(b *testing.B) {
	startPos := Position{Latitude: 40.0, Longitude: -74.0, Altitude: 10000}
	aircraft := NewBoeing737_800("BENCH-003", "BENCH789", startPos)

	dest := Position{Latitude: 40.1, Longitude: -74.0, Altitude: 10000}
	aircraft.Destination = &dest
	aircraft.State.Speed = 200

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		aircraft.updateAttitude(1.0)
	}
}

func TestAirbornePlatformUpdate(t *testing.T) {
	startPos := Position{Latitude: 40.7128, Longitude: -74.0060, Altitude: 10000}
	aircraft := NewBoeing737_800("TEST001", "UA123", startPos)

	// Set a destination
	destination := Position{Latitude: 41.0, Longitude: -74.0, Altitude: 11000}
	aircraft.SetDestination(destination)

	// Convert float64 to time.Duration for Update call
	dt := time.Duration(1.0 * float64(time.Second))
	err := aircraft.Update(dt)

	if err != nil {
		t.Errorf("Update failed: %v", err)
	}

	// Check that the aircraft moved towards destination
	if aircraft.State.Position.Latitude == startPos.Latitude {
		t.Error("Aircraft should have moved")
	}

	// Check that flight phase is appropriate
	if aircraft.FlightPhase == FlightPhaseParked {
		t.Error("Aircraft should not be parked when flying")
	}
}
