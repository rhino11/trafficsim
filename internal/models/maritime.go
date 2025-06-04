package models

import (
	"fmt"
	"time"
)

// MaritimePlatform represents naval and commercial vessels
type MaritimePlatform struct {
	// Embed UniversalPlatform for base functionality
	UniversalPlatform

	// Maritime-specific characteristics
	Draft         float64 // meters (depth below waterline)
	Displacement  float64 // tonnes
	Beam          float64 // width at widest point
	CrewSize      int
	CargoCapacity float64 // tonnes

	// Naval capabilities
	SonarRange     float64 // meters
	RadarRange     float64 // meters
	WeaponSystems  []string
	ArmorThickness float64 // mm

	// Propulsion
	PropulsionType string // diesel, gas turbine, nuclear, etc.
	Screws         int    // number of propellers
	Rudders        int    // number of rudders
}

// Core Platform interface implementation - delegate to embedded UniversalPlatform
func (m *MaritimePlatform) GetID() string           { return m.UniversalPlatform.GetID() }
func (m *MaritimePlatform) GetType() PlatformType   { return m.UniversalPlatform.GetType() }
func (m *MaritimePlatform) GetClass() string        { return m.UniversalPlatform.GetClass() }
func (m *MaritimePlatform) GetName() string         { return m.UniversalPlatform.GetName() }
func (m *MaritimePlatform) GetCallSign() string     { return m.UniversalPlatform.GetCallSign() }
func (m *MaritimePlatform) GetState() PlatformState { return m.UniversalPlatform.GetState() }
func (m *MaritimePlatform) GetMaxSpeed() float64    { return m.UniversalPlatform.GetMaxSpeed() }
func (m *MaritimePlatform) GetMaxAltitude() float64 { return 0 } // Ships don't fly
func (m *MaritimePlatform) GetLength() float64      { return m.UniversalPlatform.GetLength() }
func (m *MaritimePlatform) GetWidth() float64       { return m.UniversalPlatform.GetWidth() }
func (m *MaritimePlatform) GetHeight() float64      { return m.UniversalPlatform.GetHeight() }
func (m *MaritimePlatform) GetMass() float64        { return m.UniversalPlatform.GetMass() }

func (m *MaritimePlatform) UpdateState(state PlatformState) {
	// Keep ships at sea level
	state.Position.Altitude = 0
	m.UniversalPlatform.UpdateState(state)
}

func (m *MaritimePlatform) SetDestination(pos Position) error {
	// Force maritime platforms to sea level
	pos.Altitude = 0
	return m.UniversalPlatform.SetDestination(pos)
}

// Enhanced 3D physics methods
func (m *MaritimePlatform) Initialize3DPhysics() {
	m.UniversalPlatform.Initialize3DPhysics()
	// Ensure ship stays at sea level
	m.UniversalPlatform.State.Position.Altitude = 0
	m.UniversalPlatform.State.Physics.Position.Altitude = 0
}

func (m *MaritimePlatform) Update3DPhysics(deltaTime time.Duration) error {
	err := m.UniversalPlatform.Update3DPhysics(deltaTime)
	// Ensure ship stays at sea level
	m.UniversalPlatform.State.Position.Altitude = 0
	m.UniversalPlatform.State.Physics.Position.Altitude = 0
	return err
}

func (m *MaritimePlatform) GetPhysicsState() PhysicsState {
	return m.UniversalPlatform.GetPhysicsState()
}

func (m *MaritimePlatform) SetPhysicsState(physics PhysicsState) {
	// Ensure ship stays at sea level
	physics.Position.Altitude = 0
	m.UniversalPlatform.SetPhysicsState(physics)
}

// Update uses the base UniversalPlatform movement with maritime constraints
func (m *MaritimePlatform) Update(deltaTime time.Duration) error {
	err := m.UniversalPlatform.Update(deltaTime)
	// Ensure ship stays at sea level
	m.UniversalPlatform.State.Position.Altitude = 0
	return err
}

// Ship factory functions using UniversalPlatform base

// NewArleighBurkeDestroyer creates an Arleigh Burke-class destroyer (US Navy)
func NewArleighBurkeDestroyer(id, shipName string, startPos Position) *MaritimePlatform {
	startPos.Altitude = 0 // Ensure at sea level

	typeDef := &PlatformTypeDefinition{
		Class: "Arleigh Burke-class",
		Performance: PerformanceCharacteristics{
			MaxSpeed:      15.4,    // m/s (30+ knots)
			CruiseSpeed:   10.3,    // m/s (20 knots)
			Acceleration:  0.5,     // m/s²
			TurningRadius: 800,     // meters
			Range:         7400000, // meters (4000 nautical miles)
		},
		Physical: PhysicalCharacteristics{
			Length: 155,     // meters
			Width:  20,      // meters
			Height: 18,      // meters
			Mass:   9200000, // kg
			Draft:  6.3,     // meters
		},
	}

	config := &PlatformConfiguration{
		ID:   id,
		Name: fmt.Sprintf("USS %s", shipName),
		Type: "Arleigh Burke-class",
	}

	universalPlatform := UniversalPlatform{
		ID:           id,
		PlatformType: PlatformTypeMaritime,
		TypeDef:      typeDef,
		Config:       config,
		CallSign:     fmt.Sprintf("NAVY%s", id[len(id)-3:]),
		State: PlatformState{
			ID:          id,
			Position:    startPos,
			Velocity:    Velocity{},
			Heading:     0,
			Speed:       0,
			LastUpdated: time.Now(),
		},
	}

	return &MaritimePlatform{
		UniversalPlatform: universalPlatform,
		Draft:             6.3,  // meters
		Displacement:      9200, // tonnes
		Beam:              20,   // meters
		CrewSize:          330,
		CargoCapacity:     0,      // combat vessel
		SonarRange:        50000,  // meters
		RadarRange:        180000, // meters
		WeaponSystems:     []string{"Aegis Combat System", "VLS Missiles", "5-inch Gun"},
		ArmorThickness:    25, // mm
		PropulsionType:    "Gas Turbine",
		Screws:            2,
		Rudders:           1,
	}
}

// NewTiconderogaCruiser creates a Ticonderoga-class cruiser (US Navy)
func NewTiconderogaCruiser(id, shipName string, startPos Position) *MaritimePlatform {
	startPos.Altitude = 0

	typeDef := &PlatformTypeDefinition{
		Class: "Ticonderoga-class",
		Performance: PerformanceCharacteristics{
			MaxSpeed:      15.4,    // m/s (30+ knots)
			CruiseSpeed:   10.3,    // m/s (20 knots)
			Acceleration:  0.5,     // m/s²
			TurningRadius: 900,     // meters
			Range:         9250000, // meters (5000 nautical miles)
		},
		Physical: PhysicalCharacteristics{
			Length: 173,     // meters
			Width:  16.8,    // meters
			Height: 20,      // meters
			Mass:   9800000, // kg
			Draft:  10.2,    // meters
		},
	}

	config := &PlatformConfiguration{
		ID:   id,
		Name: fmt.Sprintf("USS %s", shipName),
		Type: "Ticonderoga-class",
	}

	universalPlatform := UniversalPlatform{
		ID:           id,
		PlatformType: PlatformTypeMaritime,
		TypeDef:      typeDef,
		Config:       config,
		CallSign:     fmt.Sprintf("NAVY%s", id[len(id)-3:]),
		State: PlatformState{
			ID:          id,
			Position:    startPos,
			Velocity:    Velocity{},
			Heading:     0,
			Speed:       0,
			LastUpdated: time.Now(),
		},
	}

	return &MaritimePlatform{
		UniversalPlatform: universalPlatform,
		Draft:             10.2, // meters
		Displacement:      9800, // tonnes
		Beam:              16.8, // meters
		CrewSize:          400,
		CargoCapacity:     0,      // combat vessel
		SonarRange:        50000,  // meters
		RadarRange:        200000, // meters
		WeaponSystems:     []string{"Aegis Combat System", "VLS Missiles", "5-inch Gun", "Tomahawk Missiles"},
		ArmorThickness:    30, // mm
		PropulsionType:    "Gas Turbine",
		Screws:            2,
		Rudders:           1,
	}
}

// NewContainerShip creates a large container vessel (commercial)
func NewContainerShip(id, shipName string, startPos Position) *MaritimePlatform {
	startPos.Altitude = 0

	typeDef := &PlatformTypeDefinition{
		Class: "Ultra Large Container Vessel",
		Performance: PerformanceCharacteristics{
			MaxSpeed:      12.9,     // m/s (25 knots)
			CruiseSpeed:   10.3,     // m/s (20 knots)
			Acceleration:  0.1,      // m/s²
			TurningRadius: 2000,     // meters
			Range:         18500000, // meters (10000 nautical miles)
		},
		Physical: PhysicalCharacteristics{
			Length: 400,       // meters
			Width:  59,        // meters
			Height: 73,        // meters
			Mass:   200000000, // kg
			Draft:  16,        // meters
		},
	}

	config := &PlatformConfiguration{
		ID:   id,
		Name: shipName,
		Type: "Ultra Large Container Vessel",
	}

	universalPlatform := UniversalPlatform{
		ID:           id,
		PlatformType: PlatformTypeMaritime,
		TypeDef:      typeDef,
		Config:       config,
		CallSign:     fmt.Sprintf("CARGO%s", id[len(id)-3:]),
		State: PlatformState{
			ID:          id,
			Position:    startPos,
			Velocity:    Velocity{},
			Heading:     0,
			Speed:       0,
			LastUpdated: time.Now(),
		},
	}

	return &MaritimePlatform{
		UniversalPlatform: universalPlatform,
		Draft:             16,     // meters
		Displacement:      200000, // tonnes
		Beam:              59,     // meters
		CrewSize:          25,
		CargoCapacity:     200000,     // tonnes
		SonarRange:        5000,       // meters
		RadarRange:        50000,      // meters
		WeaponSystems:     []string{}, // unarmed commercial vessel
		ArmorThickness:    0,          // mm
		PropulsionType:    "Diesel",
		Screws:            1,
		Rudders:           1,
	}
}

// NewOilTanker creates a large crude oil tanker (commercial)
func NewOilTanker(id, shipName string, startPos Position) *MaritimePlatform {
	startPos.Altitude = 0

	typeDef := &PlatformTypeDefinition{
		Class: "Very Large Crude Carrier",
		Performance: PerformanceCharacteristics{
			MaxSpeed:      8.2,      // m/s (16 knots)
			CruiseSpeed:   6.7,      // m/s (13 knots)
			Acceleration:  0.05,     // m/s²
			TurningRadius: 3000,     // meters
			Range:         22200000, // meters (12000 nautical miles)
		},
		Physical: PhysicalCharacteristics{
			Length: 330,       // meters
			Width:  60,        // meters
			Height: 35,        // meters
			Mass:   320000000, // kg
			Draft:  22,        // meters
		},
	}

	config := &PlatformConfiguration{
		ID:   id,
		Name: shipName,
		Type: "Very Large Crude Carrier",
	}

	universalPlatform := UniversalPlatform{
		ID:           id,
		PlatformType: PlatformTypeMaritime,
		TypeDef:      typeDef,
		Config:       config,
		CallSign:     fmt.Sprintf("TANKER%s", id[len(id)-2:]),
		State: PlatformState{
			ID:          id,
			Position:    startPos,
			Velocity:    Velocity{},
			Heading:     0,
			Speed:       0,
			LastUpdated: time.Now(),
		},
	}

	return &MaritimePlatform{
		UniversalPlatform: universalPlatform,
		Draft:             22,     // meters
		Displacement:      320000, // tonnes
		Beam:              60,     // meters
		CrewSize:          30,
		CargoCapacity:     320000,     // tonnes
		SonarRange:        2000,       // meters
		RadarRange:        25000,      // meters
		WeaponSystems:     []string{}, // unarmed commercial vessel
		ArmorThickness:    0,          // mm
		PropulsionType:    "Diesel",
		Screws:            1,
		Rudders:           1,
	}
}

// NewCoastGuardCutter creates a Coast Guard cutter
func NewCoastGuardCutter(id, shipName string, startPos Position) *MaritimePlatform {
	startPos.Altitude = 0

	typeDef := &PlatformTypeDefinition{
		Class: "Legend-class Cutter",
		Performance: PerformanceCharacteristics{
			MaxSpeed:      14.4,    // m/s (28 knots)
			CruiseSpeed:   10.3,    // m/s (20 knots)
			Acceleration:  0.8,     // m/s²
			TurningRadius: 600,     // meters
			Range:         5550000, // meters (3000 nautical miles)
		},
		Physical: PhysicalCharacteristics{
			Length: 127,     // meters
			Width:  16.4,    // meters
			Height: 15,      // meters
			Mass:   4500000, // kg
			Draft:  6.7,     // meters
		},
	}

	config := &PlatformConfiguration{
		ID:   id,
		Name: fmt.Sprintf("USCGC %s", shipName),
		Type: "Legend-class Cutter",
	}

	universalPlatform := UniversalPlatform{
		ID:           id,
		PlatformType: PlatformTypeMaritime,
		TypeDef:      typeDef,
		Config:       config,
		CallSign:     fmt.Sprintf("COASTGUARD%s", id[len(id)-2:]),
		State: PlatformState{
			ID:          id,
			Position:    startPos,
			Velocity:    Velocity{},
			Heading:     0,
			Speed:       0,
			LastUpdated: time.Now(),
		},
	}

	return &MaritimePlatform{
		UniversalPlatform: universalPlatform,
		Draft:             6.7,  // meters
		Displacement:      4500, // tonnes
		Beam:              16.4, // meters
		CrewSize:          150,
		CargoCapacity:     500,    // tonnes
		SonarRange:        20000,  // meters
		RadarRange:        100000, // meters
		WeaponSystems:     []string{"57mm Gun", "Close-in Weapons System"},
		ArmorThickness:    15, // mm
		PropulsionType:    "Diesel",
		Screws:            2,
		Rudders:           1,
	}
}
