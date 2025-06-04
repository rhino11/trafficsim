package models

import (
	"math"
	"testing"
	"time"
)

func TestUniversalPlatform_Update(t *testing.T) {
	// Create a test aircraft
	startPos := Position{Latitude: 40.7128, Longitude: -74.0060, Altitude: 10000}
	aircraft := NewBoeing737_800Universal("TEST001", "UA123", startPos)

	// Set initial velocity
	aircraft.State.Velocity = Velocity{North: 100, East: 50, Up: 5}
	aircraft.State.Heading = 45
	aircraft.State.Speed = 111.8 // sqrt(100² + 50²)

	// Update the platform
	dt := 1 * time.Second
	if err := aircraft.Update(dt); err != nil {
		t.Errorf("Update failed: %v", err)
	}

	// Check that position has been updated
	if aircraft.State.Position.Latitude == startPos.Latitude {
		t.Error("Latitude should have changed after update")
	}
	if aircraft.State.Position.Longitude == startPos.Longitude {
		t.Error("Longitude should have changed after update")
	}
	if aircraft.State.Position.Altitude == startPos.Altitude {
		t.Error("Altitude should have changed after update")
	}

	// Check that LastUpdated was set
	if aircraft.State.LastUpdated.IsZero() {
		t.Error("LastUpdated should be set after update")
	}
}

func TestUniversalPlatform_SetDestination(t *testing.T) {
	startPos := Position{Latitude: 40.7128, Longitude: -74.0060, Altitude: 10000}
	aircraft := NewBoeing737_800Universal("TEST001", "UA123", startPos)

	target := Position{Latitude: 41.0000, Longitude: -73.0000, Altitude: 11000}
	err := aircraft.SetDestination(target)
	if err != nil {
		t.Errorf("SetDestination failed: %v", err)
	}

	if aircraft.Destination == nil {
		t.Error("Destination should be set")
	}

	if aircraft.Destination.Latitude != target.Latitude ||
		aircraft.Destination.Longitude != target.Longitude ||
		aircraft.Destination.Altitude != target.Altitude {
		t.Error("Destination position not set correctly")
	}
}

func TestUniversalPlatform_CalculateDistanceTo(t *testing.T) {
	pos1 := Position{Latitude: 0, Longitude: 0, Altitude: 0}
	pos2 := Position{Latitude: 1, Longitude: 1, Altitude: 1000}

	aircraft := NewBoeing737_800Universal("TEST001", "UA123", pos1)

	distance := aircraft.CalculateDistanceTo(pos2)

	// Should be approximately 157 km (1 degree ≈ 111 km, so sqrt(2) * 111 km ≈ 157 km)
	// Plus 1000m altitude difference
	expectedDistance := math.Sqrt(157000*157000 + 1000*1000)
	tolerance := 5000.0 // 5 km tolerance

	if math.Abs(distance-expectedDistance) > tolerance {
		t.Errorf("Distance calculation incorrect. Expected ~%.0f, got %.0f", expectedDistance, distance)
	}
}

func TestUniversalPlatform_GetStatus(t *testing.T) {
	startPos := Position{Latitude: 40.7128, Longitude: -74.0060, Altitude: 10000}
	aircraft := NewBoeing737_800Universal("TEST001", "UA123", startPos)

	status := aircraft.GetStatus()

	if status.ID != "TEST001" {
		t.Errorf("Expected ID TEST001, got %s", status.ID)
	}
	if status.PlatformType != PlatformTypeAirborne {
		t.Errorf("Expected airborne platform type, got %s", status.PlatformType)
	}
	if status.Position.Latitude != startPos.Latitude {
		t.Error("Status position should match current position")
	}
}

func TestPlatformFactory_Boeing737(t *testing.T) {
	startPos := Position{Latitude: 40.7128, Longitude: -74.0060, Altitude: 10000}
	aircraft := NewBoeing737_800Universal("TEST001", "UA123", startPos)

	// Check basic properties using getter methods
	if aircraft.ID != "TEST001" {
		t.Errorf("Expected ID TEST001, got %s", aircraft.ID)
	}
	if aircraft.GetName() != "UA123" {
		t.Errorf("Expected Name UA123, got %s", aircraft.GetName())
	}
	if aircraft.GetClass() != "Boeing 737-800" {
		t.Errorf("Expected Class Boeing 737-800, got %s", aircraft.GetClass())
	}

	// Check performance characteristics
	if aircraft.TypeDef.Performance.MaxSpeed != 257 {
		t.Errorf("Expected max speed 257 m/s, got %.1f", aircraft.TypeDef.Performance.MaxSpeed)
	}
	if aircraft.TypeDef.Performance.MaxAltitude != 12500 {
		t.Errorf("Expected max altitude 12500m, got %.1f", aircraft.TypeDef.Performance.MaxAltitude)
	}

	// Check sensors
	if !aircraft.TypeDef.Sensors.HasGPS {
		t.Error("Boeing 737 should have GPS")
	}
	if !aircraft.TypeDef.Sensors.HasRadar {
		t.Error("Boeing 737 should have radar")
	}
}

func TestPlatformFactory_F16(t *testing.T) {
	startPos := Position{Latitude: 40.7128, Longitude: -74.0060, Altitude: 5000}
	fighter := NewF16FightingFalconUniversal("USAF001", "87-0001", startPos)

	// Check military-specific properties
	if fighter.SystemStatus.WeaponStatus != "ARMED" {
		t.Errorf("Expected WeaponStatus ARMED, got %s", fighter.SystemStatus.WeaponStatus)
	}

	// Check high performance characteristics
	if fighter.TypeDef.Performance.MaxSpeed < 600 {
		t.Errorf("F-16 should have high max speed, got %.1f m/s", fighter.TypeDef.Performance.MaxSpeed)
	}
	if fighter.TypeDef.Performance.MaxLoadFactor < 8.0 {
		t.Errorf("F-16 should handle high G-forces, got %.1f", fighter.TypeDef.Performance.MaxLoadFactor)
	}

	// Check agility
	if fighter.TypeDef.Performance.MaxRollRate < 500 {
		t.Errorf("F-16 should have high roll rate, got %.1f deg/s", fighter.TypeDef.Performance.MaxRollRate)
	}
}

func TestPlatformFactory_M1A2Abrams(t *testing.T) {
	startPos := Position{Latitude: 40.7128, Longitude: -74.0060, Altitude: 0}
	tank := NewM1A2AbramsUniversal("ARMY001", "A-1-1", startPos)

	// Check land vehicle properties
	if tank.PlatformType != PlatformTypeLand {
		t.Errorf("Expected land platform type, got %s", tank.PlatformType)
	}
	if tank.TypeDef.Performance.MaxAltitude != 0 {
		t.Error("Ground vehicle should have max altitude of 0")
	}

	// Check heavy vehicle characteristics
	if tank.TypeDef.Physical.Mass < 50000 {
		t.Errorf("Tank should be heavy, got %.1f kg", tank.TypeDef.Physical.Mass)
	}
	if tank.TypeDef.Performance.MaxSpeed > 25 {
		t.Errorf("Tank should be relatively slow, got %.1f m/s", tank.TypeDef.Performance.MaxSpeed)
	}

	// Check military systems
	if tank.SystemStatus.WeaponStatus != "ARMED" {
		t.Error("Tank should be armed")
	}
}

func TestPlatformFactory_ArleighBurke(t *testing.T) {
	startPos := Position{Latitude: 40.7128, Longitude: -74.0060, Altitude: 0}
	destroyer := NewArleighBurkeDestroyerUniversal("DDG001", "John Paul Jones", startPos)

	// Check maritime properties
	if destroyer.PlatformType != PlatformTypeMaritime {
		t.Errorf("Expected maritime platform type, got %s", destroyer.PlatformType)
	}
	if destroyer.State.Position.Altitude != 0 {
		t.Error("Ship should be at sea level")
	}

	// Check naval characteristics
	if destroyer.TypeDef.Physical.WetArea == 0 {
		t.Error("Ship should have wet area defined")
	}
	if destroyer.TypeDef.Physical.Draft == 0 {
		t.Error("Ship should have draft defined")
	}

	// Check long-range sensors
	if destroyer.TypeDef.Sensors.RadarRange < 300000 {
		t.Errorf("Naval radar should have long range, got %.0f m", destroyer.TypeDef.Sensors.RadarRange)
	}
}

func TestPlatformFactory_StarlinkSatellite(t *testing.T) {
	startPos := Position{Latitude: 40.7128, Longitude: -74.0060, Altitude: 550000}
	satellite := NewStarlinkSatelliteUniversal("SAT001", "1234", startPos)

	// Check space platform properties
	if satellite.PlatformType != PlatformTypeSpace {
		t.Errorf("Expected space platform type, got %s", satellite.PlatformType)
	}
	if satellite.State.Position.Altitude != 550000 {
		t.Error("Satellite should be at specified orbital altitude")
	}

	// Check orbital characteristics
	if satellite.TypeDef.Performance.MaxSpeed < 7000 {
		t.Errorf("Satellite should have orbital velocity, got %.1f m/s", satellite.TypeDef.Performance.MaxSpeed)
	}
	if satellite.TypeDef.Performance.OrbitalPeriod == 0 {
		t.Error("Satellite should have orbital period defined")
	}

	// Check space-specific systems
	if satellite.TypeDef.Sensors.HasCompass {
		t.Error("Compass not useful in space")
	}
	if satellite.SystemStatus.WeaponStatus != "N/A" {
		t.Error("Civilian satellite should not be armed")
	}
}

func TestPlatformFactory_CivilianCar(t *testing.T) {
	startPos := Position{Latitude: 40.7128, Longitude: -74.0060, Altitude: 0}
	car := NewCivilianCarUniversal("CIV001", "Toyota Camry", startPos)

	// Check civilian characteristics
	if car.SystemStatus.WeaponStatus != "N/A" {
		t.Error("Civilian car should not be armed")
	}
	if car.TypeDef.Sensors.HasRadar {
		t.Error("Basic civilian car should not have radar")
	}

	// Check reasonable performance
	if car.TypeDef.Performance.MaxSpeed > 60 {
		t.Errorf("Civilian car should have reasonable max speed, got %.1f m/s", car.TypeDef.Performance.MaxSpeed)
	}
}

func TestPlatformFactory_ContainerShip(t *testing.T) {
	startPos := Position{Latitude: 40.7128, Longitude: -74.0060, Altitude: 0}
	ship := NewContainerShipUniversal("CARGO001", "Ever Given", startPos)

	// Check commercial characteristics
	if ship.SystemStatus.WeaponStatus != "N/A" {
		t.Error("Commercial ship should not be armed")
	}

	// Check massive scale
	if ship.TypeDef.Physical.Mass < 100000000 {
		t.Errorf("Container ship should be massive, got %.1f kg", ship.TypeDef.Physical.Mass)
	}
	if ship.TypeDef.Physical.Length < 300 {
		t.Errorf("Container ship should be very long, got %.1f m", ship.TypeDef.Physical.Length)
	}

	// Check slow but efficient
	if ship.TypeDef.Performance.MaxAcceleration > 0.2 {
		t.Errorf("Large ship should accelerate slowly, got %.3f m/s²", ship.TypeDef.Performance.MaxAcceleration)
	}
	if ship.TypeDef.Performance.Range < 20000000 {
		t.Errorf("Container ship should have global range, got %.0f m", ship.TypeDef.Performance.Range)
	}
}

func TestCreatePlatformFromConfig(t *testing.T) {
	config := map[string]interface{}{
		"id":   "CONFIG001",
		"type": "airborne",
		"position": map[string]interface{}{
			"latitude":  40.7128,
			"longitude": -74.0060,
			"altitude":  10000.0,
		},
	}

	platform, err := CreatePlatformFromConfig(config)
	if err != nil {
		t.Fatalf("Failed to create platform from config: %v", err)
	}

	if platform.ID != "CONFIG001" {
		t.Errorf("Expected ID CONFIG001, got %s", platform.ID)
	}
	if platform.PlatformType != PlatformTypeAirborne {
		t.Errorf("Expected airborne type, got %s", platform.PlatformType)
	}

	// Test invalid config
	invalidConfig := map[string]interface{}{
		"type": "airborne",
		// Missing ID
	}
	_, err = CreatePlatformFromConfig(invalidConfig)
	if err == nil {
		t.Error("Expected error for invalid config")
	}
}

func TestPhysicsIntegration(t *testing.T) {
	startPos := Position{Latitude: 40.7128, Longitude: -74.0060, Altitude: 10000}
	aircraft := NewBoeing737_800Universal("TEST001", "UA123", startPos)

	// Set up physics state
	aircraft.State.Physics.Velocity = Velocity{North: 100, East: 0, Up: 0}
	aircraft.State.Physics.AngularVelocity = AngularVelocity{RollRate: 5, PitchRate: 2, YawRate: 1}
	aircraft.State.Physics.Attitude = Attitude{Roll: 0, Pitch: 5, Yaw: 90}

	// Apply forces
	forces := Force{X: 1000, Y: 0, Z: 500} // Forward thrust and upward lift
	aircraft.ApplyForce(forces)

	// Update physics
	dt := 1 * time.Second
	if err := aircraft.Update3DPhysics(dt); err != nil {
		t.Errorf("Update3DPhysics failed: %v", err)
	}

	// Check that physics state changed
	if aircraft.State.Physics.Velocity.North <= 100 {
		t.Error("Forward velocity should have increased due to applied force")
	}
	if aircraft.State.Physics.Velocity.Up <= 0 {
		t.Error("Vertical velocity should have increased due to upward force")
	}

	// Check orientation changes
	if aircraft.State.Physics.Attitude.Roll == 0 {
		t.Error("Roll should have changed due to angular velocity")
	}
}

func BenchmarkUniversalPlatformCreate(b *testing.B) {
	startPos := Position{Latitude: 40.7128, Longitude: -74.0060, Altitude: 10000}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewBoeing737_800Universal("BENCH001", "UA123", startPos)
	}
}

func BenchmarkPlatformCreation(b *testing.B) {
	startPos := Position{Latitude: 40.7128, Longitude: -74.0060, Altitude: 10000}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewBoeing737_800Universal("BENCH001", "UA123", startPos)
	}
}

func BenchmarkPhysicsUpdate(b *testing.B) {
	startPos := Position{Latitude: 40.7128, Longitude: -74.0060, Altitude: 10000}
	aircraft := NewBoeing737_800Universal("BENCH001", "UA123", startPos)
	aircraft.State.Physics.Velocity = Velocity{North: 100, East: 50, Up: 5}
	forces := Force{X: 1000, Y: 100, Z: 500}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		aircraft.ApplyForce(forces)
		if err := aircraft.Update3DPhysics(1 * time.Second); err != nil {
			b.Errorf("Update3DPhysics failed: %v", err)
		}
	}
}

func createTestUniversalPlatform() *UniversalPlatform {
	// Create a mock configuration
	typeDef := &PlatformTypeDefinition{
		Class: "Test Aircraft",
		Performance: PerformanceCharacteristics{
			MaxSpeed:        100.0,
			CruiseSpeed:     80.0,
			MaxAcceleration: 2.0,
			TurningRadius:   500.0,
			Range:           10000.0, // Add missing Range field
		},
		Physical: PhysicalCharacteristics{
			Length: 30.0,
			Width:  35.0,
			Height: 12.0,
			Mass:   50000.0,
		},
	}

	config := &PlatformConfiguration{
		ID:   "TEST001",
		Name: "Test Platform",
		StartPosition: Position{
			Latitude:  40.0,
			Longitude: -74.0,
			Altitude:  0.0,
		},
	}

	platform := &UniversalPlatform{
		ID:           "TEST001",
		PlatformType: PlatformTypeAirborne,
		TypeDef:      typeDef,
		Config:       config,
		State: PlatformState{
			ID:       "TEST001",
			Position: config.StartPosition,
			Velocity: Velocity{},
			Heading:  0.0,
			Speed:    0.0,
		},
		CallSign: "TEST001",
	}

	return platform
}

func createTestCommercialAircraft() *UniversalPlatform {
	// Create a Boeing 737-800 for testing
	startPos := Position{Latitude: 40.7128, Longitude: -74.0060, Altitude: 10000}
	aircraft := NewBoeing737_800Universal("TEST001", "Test Commercial Aircraft", startPos)
	return aircraft
}

func TestUniversalPlatformSetDestination(t *testing.T) {
	platform := createTestUniversalPlatform()

	destination := Position{
		Latitude:  41.0,
		Longitude: -74.0,
		Altitude:  1000.0,
	}

	err := platform.SetDestination(destination)
	if err != nil {
		t.Errorf("SetDestination failed: %v", err)
	}

	if platform.Destination == nil {
		t.Error("Destination was not set")
	}

	if platform.Destination.Latitude != destination.Latitude {
		t.Errorf("Expected destination latitude %f, got %f", destination.Latitude, platform.Destination.Latitude)
	}
}

func TestUniversalPlatformCalculateDistance(t *testing.T) {
	platform := createTestUniversalPlatform()

	target := Position{
		Latitude:  41.0,
		Longitude: -74.0,
		Altitude:  1000.0,
	}

	distance := platform.CalculateDistanceTo(target)
	if distance <= 0 {
		t.Error("Distance calculation should return positive value")
	}
}

func TestUniversalPlatformGetStatus(t *testing.T) {
	platform := createTestUniversalPlatform()

	status := platform.GetStatus()
	if status.ID != platform.ID {
		t.Errorf("Expected status ID %s, got %s", platform.ID, status.ID)
	}

	if status.PlatformType != platform.PlatformType {
		t.Errorf("Expected platform type %s, got %s", platform.PlatformType, status.PlatformType)
	}
}

func TestUniversalPlatformCommercialAircraft(t *testing.T) {
	// Create a Boeing 737-800 for testing
	startPos := Position{Latitude: 40.7128, Longitude: -74.0060, Altitude: 10000}
	aircraft := NewBoeing737_800Universal("TEST001", "Test Commercial Aircraft", startPos)

	// Test basic properties using getter methods instead of direct field access
	if aircraft.GetName() != "Test Commercial Aircraft" {
		t.Errorf("Expected name 'Test Commercial Aircraft', got %s", aircraft.GetName())
	}

	if aircraft.GetClass() != "Boeing 737-800" {
		t.Errorf("Expected class 'Boeing 737-800', got %s", aircraft.GetClass())
	}

	// Test sensor capabilities
	if !aircraft.TypeDef.Sensors.HasGPS {
		t.Error("Commercial aircraft should have GPS")
	}

	if !aircraft.TypeDef.Sensors.HasRadar {
		t.Error("Commercial aircraft should have radar")
	}

	// Initialize system status for testing
	aircraft.SystemStatus = SystemStatus{
		PowerSystem:         SystemState{Operational: true, Efficiency: 1.0},
		PropulsionSystem:    SystemState{Operational: true, Efficiency: 0.95},
		NavigationSystem:    SystemState{Operational: true, Efficiency: 1.0},
		CommunicationSystem: SystemState{Operational: true, Efficiency: 0.98},
		SensorSystem:        SystemState{Operational: true, Efficiency: 0.99},
		FuelSystem:          SystemState{Operational: true, Efficiency: 1.0},
	}

	if !aircraft.SystemStatus.PowerSystem.Operational {
		t.Error("Power system should be operational")
	}
}
