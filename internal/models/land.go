package models

import (
	"fmt"
	"time"
)

// Constants to reduce string duplication
const (
	DefaultHeading       = 0
	DefaultSpeed         = 0
	ArmorCallsignPrefix  = "ARMOR"
	HumveeCallsignPrefix = "HUMVEE"
	TruckCallsignPrefix  = "TRUCK"
	PatrolCallsignPrefix = "PATROL"
)

// LandPlatform represents ground vehicles and installations
type LandPlatform struct {
	// Embed UniversalPlatform for base functionality
	UniversalPlatform

	// Land-specific characteristics
	MaxGradient     float64 // degrees (max slope)
	GroundClearance float64 // meters
	TurningRadius   float64 // meters
	FuelConsumption float64 // liters per km
	CargoCapacity   float64 // kg
	CrewCapacity    int

	// Terrain capabilities
	OffRoadCapable bool
	WaterFording   float64 // max depth in meters
	ClimbAngle     float64 // degrees
}

// Core Platform interface implementation - delegate to embedded UniversalPlatform
func (l *LandPlatform) GetID() string           { return l.UniversalPlatform.GetID() }
func (l *LandPlatform) GetType() PlatformType   { return l.UniversalPlatform.GetType() }
func (l *LandPlatform) GetClass() string        { return l.UniversalPlatform.GetClass() }
func (l *LandPlatform) GetName() string         { return l.UniversalPlatform.GetName() }
func (l *LandPlatform) GetCallSign() string     { return l.UniversalPlatform.GetCallSign() }
func (l *LandPlatform) GetState() PlatformState { return l.UniversalPlatform.GetState() }
func (l *LandPlatform) GetMaxSpeed() float64    { return l.UniversalPlatform.GetMaxSpeed() }
func (l *LandPlatform) GetMaxAltitude() float64 { return l.UniversalPlatform.GetMaxAltitude() }
func (l *LandPlatform) GetLength() float64      { return l.UniversalPlatform.GetLength() }
func (l *LandPlatform) GetWidth() float64       { return l.UniversalPlatform.GetWidth() }
func (l *LandPlatform) GetHeight() float64      { return l.UniversalPlatform.GetHeight() }
func (l *LandPlatform) GetMass() float64        { return l.UniversalPlatform.GetMass() }

func (l *LandPlatform) UpdateState(state PlatformState) {
	l.UniversalPlatform.UpdateState(state)
}

func (l *LandPlatform) SetDestination(pos Position) error {
	return l.UniversalPlatform.SetDestination(pos)
}

// Enhanced 3D physics methods
func (l *LandPlatform) Initialize3DPhysics() {
	l.UniversalPlatform.Initialize3DPhysics()
}

func (l *LandPlatform) Update3DPhysics(deltaTime time.Duration) error {
	return l.UniversalPlatform.Update3DPhysics(deltaTime)
}

func (l *LandPlatform) GetPhysicsState() PhysicsState {
	return l.UniversalPlatform.GetPhysicsState()
}

func (l *LandPlatform) SetPhysicsState(physics PhysicsState) {
	l.UniversalPlatform.SetPhysicsState(physics)
}

// Update uses the base UniversalPlatform movement but can be enhanced with land-specific logic
func (l *LandPlatform) Update(deltaTime time.Duration) error {
	// Use enhanced land movement from UniversalPlatform
	return l.UniversalPlatform.Update(deltaTime)
}

// createLandPlatformBase creates a base UniversalPlatform for land vehicles
func createLandPlatformBase(id, name, platformType, callsign string, startPos Position, typeDef *PlatformTypeDefinition) UniversalPlatform {
	return UniversalPlatform{
		ID:           id,
		PlatformType: PlatformTypeLand,
		TypeDef:      typeDef,
		Config: &PlatformConfiguration{
			ID:            id,
			Type:          platformType,
			Name:          name,
			StartPosition: startPos,
		},
		State: PlatformState{
			ID:          id,
			Position:    startPos,
			Velocity:    Velocity{},
			Heading:     0,
			Speed:       0,
			LastUpdated: time.Now(),
			Physics: PhysicsState{
				Position:        startPos,
				Mass:            typeDef.Physical.Mass,
				MomentOfInertia: calculateMomentOfInertia(typeDef.Physical.Mass, PlatformTypeLand),
			},
		},
		CallSign:      callsign,
		FuelRemaining: typeDef.Physical.FuelCapacity,
		MissionTime:   0,
		SystemStatus: SystemStatus{
			PropulsionSystem:    SystemState{Operational: true, Efficiency: 1.0, LastCheck: time.Now()},
			NavigationSystem:    SystemState{Operational: true, Efficiency: 1.0, LastCheck: time.Now()},
			CommunicationSystem: SystemState{Operational: true, Efficiency: 1.0, LastCheck: time.Now()},
			SensorSystem:        SystemState{Operational: true, Efficiency: 1.0, LastCheck: time.Now()},
			FuelSystem:          SystemState{Operational: true, Efficiency: 1.0, LastCheck: time.Now()},
			WeaponStatus:        WeaponStatusSafe,
		},
		lastPosition: startPos,
		acceleration: 0,
	}
}

// CreateM1A2Tank creates an M1A2 Abrams main battle tank
func CreateM1A2Tank(id string, callsign string, startPos Position) *LandPlatform {
	if callsign == "" {
		callsign = fmt.Sprintf("%s-%s", ArmorCallsignPrefix, id)
	}

	typeDef := &PlatformTypeDefinition{
		Class:    "Main Battle Tank",
		Category: "military",
		Performance: PerformanceCharacteristics{
			MaxSpeed:      18.6, // m/s (67 km/h)
			CruiseSpeed:   11.1, // m/s (40 km/h)
			MaxAltitude:   4267, // meters
			TurningRadius: 8.4,  // meters
		},
		Physical: PhysicalCharacteristics{
			Length: 9.77,  // meters
			Width:  3.66,  // meters
			Height: 2.44,  // meters
			Mass:   62000, // kg
		},
	}

	base := createLandPlatformBase(id, "M1A2 Abrams", "Main Battle Tank", callsign, startPos, typeDef)

	return &LandPlatform{
		UniversalPlatform: base,
		MaxGradient:       60,    // degrees
		GroundClearance:   0.432, // meters
		TurningRadius:     8.4,   // meters
		FuelConsumption:   2.6,   // liters per km
		CargoCapacity:     1000,  // kg
		CrewCapacity:      4,
		OffRoadCapable:    true,
		WaterFording:      1.2, // meters
		ClimbAngle:        60,  // degrees
	}
}

// CreateM2Bradley creates an M2 Bradley Infantry Fighting Vehicle
func CreateM2Bradley(id string, callsign string, startPos Position) *LandPlatform {
	if callsign == "" {
		callsign = fmt.Sprintf("%s-%s", ArmorCallsignPrefix, id)
	}

	typeDef := &PlatformTypeDefinition{
		Class:    "Infantry Fighting Vehicle",
		Category: "military",
		Performance: PerformanceCharacteristics{
			MaxSpeed:      18.3, // m/s (66 km/h)
			CruiseSpeed:   11.1, // m/s (40 km/h)
			MaxAltitude:   4267, // meters
			TurningRadius: 6.0,  // meters
		},
		Physical: PhysicalCharacteristics{
			Length: 6.55,  // meters
			Width:  3.6,   // meters
			Height: 2.98,  // meters
			Mass:   27600, // kg
		},
	}

	base := createLandPlatformBase(id, "M2 Bradley", "Infantry Fighting Vehicle", callsign, startPos, typeDef)

	return &LandPlatform{
		UniversalPlatform: base,
		MaxGradient:       60,    // degrees
		GroundClearance:   0.432, // meters
		TurningRadius:     6.0,   // meters
		FuelConsumption:   1.5,   // liters per km
		CargoCapacity:     2000,  // kg
		CrewCapacity:      9,     // 3 crew + 6 infantry
		OffRoadCapable:    true,
		WaterFording:      1.0, // meters
		ClimbAngle:        60,  // degrees
	}
}

// CreateHumvee creates a High Mobility Multipurpose Wheeled Vehicle
func CreateHumvee(id string, callsign string, startPos Position) *LandPlatform {
	if callsign == "" {
		callsign = fmt.Sprintf("%s-%s", HumveeCallsignPrefix, id)
	}

	typeDef := &PlatformTypeDefinition{
		Class:    "Utility Vehicle",
		Category: "military",
		Performance: PerformanceCharacteristics{
			MaxSpeed:      31.4, // m/s (113 km/h)
			CruiseSpeed:   19.4, // m/s (70 km/h)
			MaxAltitude:   4267, // meters
			TurningRadius: 7.6,  // meters
		},
		Physical: PhysicalCharacteristics{
			Length: 4.57, // meters
			Width:  2.16, // meters
			Height: 1.75, // meters
			Mass:   2359, // kg
		},
	}

	base := createLandPlatformBase(id, "HMMWV", "Utility Vehicle", callsign, startPos, typeDef)

	return &LandPlatform{
		UniversalPlatform: base,
		MaxGradient:       60,    // degrees
		GroundClearance:   0.406, // meters
		TurningRadius:     7.6,   // meters
		FuelConsumption:   0.8,   // liters per km
		CargoCapacity:     1200,  // kg
		CrewCapacity:      4,
		OffRoadCapable:    true,
		WaterFording:      0.76, // meters
		ClimbAngle:        60,   // degrees
	}
}

// CreateLAV25 creates a Light Armored Vehicle-25
func CreateLAV25(id string, callsign string, startPos Position) *LandPlatform {
	if callsign == "" {
		callsign = fmt.Sprintf("%s-%s", ArmorCallsignPrefix, id)
	}

	typeDef := &PlatformTypeDefinition{
		Class:    "Light Armored Vehicle",
		Category: "military",
		Performance: PerformanceCharacteristics{
			MaxSpeed:      27.8, // m/s (100 km/h)
			CruiseSpeed:   16.7, // m/s (60 km/h)
			MaxAltitude:   4267, // meters
			TurningRadius: 5.5,  // meters
		},
		Physical: PhysicalCharacteristics{
			Length: 6.39,  // meters
			Width:  2.5,   // meters
			Height: 2.69,  // meters
			Mass:   12900, // kg
		},
	}

	base := createLandPlatformBase(id, "LAV-25", "Light Armored Vehicle", callsign, startPos, typeDef)

	return &LandPlatform{
		UniversalPlatform: base,
		MaxGradient:       60,   // degrees
		GroundClearance:   0.5,  // meters
		TurningRadius:     5.5,  // meters
		FuelConsumption:   1.2,  // liters per km
		CargoCapacity:     1500, // kg
		CrewCapacity:      6,    // 3 crew + 3 marines
		OffRoadCapable:    true,
		WaterFording:      1.5, // meters (amphibious)
		ClimbAngle:        60,  // degrees
	}
}

// CreateM35Truck creates an M35 2.5-ton cargo truck
func CreateM35Truck(id string, callsign string, startPos Position) *LandPlatform {
	if callsign == "" {
		callsign = fmt.Sprintf("%s-%s", TruckCallsignPrefix, id)
	}

	typeDef := &PlatformTypeDefinition{
		Class:    "Cargo Truck",
		Category: "military",
		Performance: PerformanceCharacteristics{
			MaxSpeed:      25.0, // m/s (90 km/h)
			CruiseSpeed:   16.7, // m/s (60 km/h)
			MaxAltitude:   4267, // meters
			TurningRadius: 8.7,  // meters
		},
		Physical: PhysicalCharacteristics{
			Length: 6.71, // meters
			Width:  2.44, // meters
			Height: 2.94, // meters
			Mass:   6350, // kg empty
		},
	}

	base := createLandPlatformBase(id, "M35 Truck", "Cargo Truck", callsign, startPos, typeDef)

	return &LandPlatform{
		UniversalPlatform: base,
		MaxGradient:       30,    // degrees
		GroundClearance:   0.279, // meters
		TurningRadius:     8.7,   // meters
		FuelConsumption:   0.6,   // liters per km
		CargoCapacity:     2268,  // kg (5000 lbs)
		CrewCapacity:      3,     // driver + 2 passengers
		OffRoadCapable:    true,
		WaterFording:      0.76, // meters
		ClimbAngle:        30,   // degrees
	}
}

// CreateM1126Stryker creates an M1126 Stryker Infantry Carrier Vehicle
func CreateM1126Stryker(id string, callsign string, startPos Position) *LandPlatform {
	if callsign == "" {
		callsign = fmt.Sprintf("%s-%s", ArmorCallsignPrefix, id)
	}

	typeDef := &PlatformTypeDefinition{
		Class:    "Infantry Carrier Vehicle",
		Category: "military",
		Performance: PerformanceCharacteristics{
			MaxSpeed:      27.8, // m/s (100 km/h)
			CruiseSpeed:   19.4, // m/s (70 km/h)
			MaxAltitude:   4267, // meters
			TurningRadius: 7.0,  // meters
		},
		Physical: PhysicalCharacteristics{
			Length: 6.95,  // meters
			Width:  2.72,  // meters
			Height: 2.64,  // meters
			Mass:   16470, // kg
		},
	}

	base := createLandPlatformBase(id, "M1126 Stryker", "Infantry Carrier Vehicle", callsign, startPos, typeDef)

	return &LandPlatform{
		UniversalPlatform: base,
		MaxGradient:       60,    // degrees
		GroundClearance:   0.533, // meters
		TurningRadius:     7.0,   // meters
		FuelConsumption:   1.1,   // liters per km
		CargoCapacity:     1800,  // kg
		CrewCapacity:      11,    // 2 crew + 9 infantry
		OffRoadCapable:    true,
		WaterFording:      1.0, // meters
		ClimbAngle:        60,  // degrees
	}
}

// CreateMRAP creates a Mine-Resistant Ambush Protected vehicle
func CreateMRAP(id string, callsign string, startPos Position) *LandPlatform {
	if callsign == "" {
		callsign = fmt.Sprintf("%s-%s", PatrolCallsignPrefix, id)
	}

	typeDef := &PlatformTypeDefinition{
		Class:    "Mine-Resistant Vehicle",
		Category: "military",
		Performance: PerformanceCharacteristics{
			MaxSpeed:      29.2, // m/s (105 km/h)
			CruiseSpeed:   19.4, // m/s (70 km/h)
			MaxAltitude:   4267, // meters
			TurningRadius: 8.0,  // meters
		},
		Physical: PhysicalCharacteristics{
			Length: 6.7,   // meters
			Width:  2.7,   // meters
			Height: 2.7,   // meters
			Mass:   14500, // kg
		},
	}

	base := createLandPlatformBase(id, "MRAP", "Mine-Resistant Vehicle", callsign, startPos, typeDef)

	return &LandPlatform{
		UniversalPlatform: base,
		MaxGradient:       30,    // degrees
		GroundClearance:   0.406, // meters
		TurningRadius:     8.0,   // meters
		FuelConsumption:   1.0,   // liters per km
		CargoCapacity:     1000,  // kg
		CrewCapacity:      6,     // crew + passengers
		OffRoadCapable:    true,
		WaterFording:      0.6, // meters
		ClimbAngle:        30,  // degrees
	}
}
