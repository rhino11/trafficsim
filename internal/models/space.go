package models

import (
	"fmt"
	"math"
	"time"
)

// SpacePlatform represents satellites and spacecraft
type SpacePlatform struct {
	// Embed UniversalPlatform for base functionality
	UniversalPlatform

	// Space-specific characteristics
	OrbitalPeriod float64 // seconds (time for one orbit)
	Apogee        float64 // meters (highest point)
	Perigee       float64 // meters (lowest point)
	Inclination   float64 // degrees (orbital plane angle)
	Eccentricity  float64 // orbital eccentricity (0 = circular)

	// Mission characteristics
	MissionType     string // communication, navigation, observation, etc.
	LaunchDate      time.Time
	MissionDuration time.Duration // expected mission lifetime

	// Space environment
	SolarPanelArea  float64 // square meters
	PowerGeneration float64 // watts
	FuelRemaining   float64 // kg
	RadiationShield bool
}

// Core Platform interface implementation - delegate to embedded UniversalPlatform
func (s *SpacePlatform) GetID() string           { return s.UniversalPlatform.GetID() }
func (s *SpacePlatform) GetType() PlatformType   { return s.UniversalPlatform.GetType() }
func (s *SpacePlatform) GetClass() string        { return s.UniversalPlatform.GetClass() }
func (s *SpacePlatform) GetName() string         { return s.UniversalPlatform.GetName() }
func (s *SpacePlatform) GetCallSign() string     { return s.UniversalPlatform.GetCallSign() }
func (s *SpacePlatform) GetState() PlatformState { return s.UniversalPlatform.GetState() }
func (s *SpacePlatform) GetMaxSpeed() float64    { return s.UniversalPlatform.GetMaxSpeed() }
func (s *SpacePlatform) GetMaxAltitude() float64 { return s.Apogee }
func (s *SpacePlatform) GetLength() float64      { return s.UniversalPlatform.GetLength() }
func (s *SpacePlatform) GetWidth() float64       { return s.UniversalPlatform.GetWidth() }
func (s *SpacePlatform) GetHeight() float64      { return s.UniversalPlatform.GetHeight() }
func (s *SpacePlatform) GetMass() float64        { return s.UniversalPlatform.GetMass() }

func (s *SpacePlatform) UpdateState(state PlatformState) {
	s.UniversalPlatform.UpdateState(state)
}

func (s *SpacePlatform) SetDestination(pos Position) error {
	return s.UniversalPlatform.SetDestination(pos)
}

// Enhanced 3D physics methods
func (s *SpacePlatform) Initialize3DPhysics() {
	s.UniversalPlatform.Initialize3DPhysics()
}

func (s *SpacePlatform) Update3DPhysics(deltaTime time.Duration) error {
	return s.UniversalPlatform.Update3DPhysics(deltaTime)
}

func (s *SpacePlatform) GetPhysicsState() PhysicsState {
	return s.UniversalPlatform.GetPhysicsState()
}

func (s *SpacePlatform) SetPhysicsState(physics PhysicsState) {
	s.UniversalPlatform.SetPhysicsState(physics)
}

// Update implements simplified orbital mechanics
func (s *SpacePlatform) Update(deltaTime time.Duration) error {
	// Use enhanced orbital mechanics from UniversalPlatform base
	// but can override for specialized space platform behavior

	// Simplified orbital mechanics - circular orbit approximation
	dt := deltaTime.Seconds()

	// Calculate orbital velocity based on altitude
	earthRadius := 6371000.0 // meters
	altitude := s.UniversalPlatform.State.Position.Altitude
	orbitalRadius := earthRadius + altitude

	// Simplified orbital velocity: v = sqrt(GM/r)
	// Using approximation for Earth: GM ≈ 3.986e14 m³/s²
	GM := 3.986e14
	orbitalVelocity := math.Sqrt(GM / orbitalRadius)

	// Angular velocity (radians per second)
	angularVelocity := orbitalVelocity / orbitalRadius

	// Update position in simplified circular orbit
	// Convert lat/lon to radians for calculation
	lonRad := s.UniversalPlatform.State.Position.Longitude * math.Pi / 180

	// Simple eastward progression (simplified)
	lonRad += angularVelocity * dt

	// Wrap longitude
	for lonRad > math.Pi {
		lonRad -= 2 * math.Pi
	}
	for lonRad < -math.Pi {
		lonRad += 2 * math.Pi
	}

	// Convert back to degrees
	s.UniversalPlatform.State.Position.Longitude = lonRad * 180 / math.Pi

	// Update speed and heading
	s.UniversalPlatform.State.Speed = orbitalVelocity
	s.UniversalPlatform.State.Heading = 90 // Eastward movement
	s.UniversalPlatform.State.LastUpdated = time.Now()

	return nil
}

// Space platform factory functions using UniversalPlatform base

// NewISSModule creates an International Space Station module
func NewISSModule(id, moduleName string, startPos Position) *SpacePlatform {
	// ISS orbits at approximately 408 km altitude
	startPos.Altitude = 408000

	typeDef := &PlatformTypeDefinition{
		Class: "ISS Module",
		Performance: PerformanceCharacteristics{
			MaxSpeed:      7660, // m/s (ISS orbital velocity)
			CruiseSpeed:   7660, // m/s
			Acceleration:  0,    // orbital mechanics
			TurningRadius: 0,    // not applicable
			Range:         0,    // unlimited in orbit
		},
		Physical: PhysicalCharacteristics{
			Length: 73,     // meters (full ISS)
			Width:  109,    // meters (solar array span)
			Height: 20,     // meters
			Mass:   420000, // kg
		},
	}

	config := &PlatformConfiguration{
		ID:   id,
		Name: fmt.Sprintf("ISS %s", moduleName),
		Type: "ISS Module",
	}

	universalPlatform := UniversalPlatform{
		ID:           id,
		PlatformType: PlatformTypeSpace,
		TypeDef:      typeDef,
		Config:       config,
		CallSign:     fmt.Sprintf("ISS%s", id[len(id)-2:]),
		State: PlatformState{
			ID:          id,
			Position:    startPos,
			Velocity:    Velocity{},
			Heading:     90,   // Eastward
			Speed:       7660, // m/s
			LastUpdated: time.Now(),
		},
	}

	return &SpacePlatform{
		UniversalPlatform: universalPlatform,
		OrbitalPeriod:     5520,   // seconds (92 minutes)
		Apogee:            420000, // meters
		Perigee:           408000, // meters
		Inclination:       51.6,   // degrees
		Eccentricity:      0.0001, // nearly circular
		MissionType:       "space station",
		LaunchDate:        time.Date(1998, 11, 20, 0, 0, 0, 0, time.UTC),
		MissionDuration:   time.Hour * 24 * 365 * 30, // 30 years planned
		SolarPanelArea:    2500,                      // square meters
		PowerGeneration:   84000,                     // watts
		FuelRemaining:     1000,                      // kg (for reboost)
		RadiationShield:   true,
	}
}

// NewStarlinkSatellite creates a Starlink communication satellite
func NewStarlinkSatellite(id, satelliteNumber string, startPos Position) *SpacePlatform {
	// Starlink operates at ~550 km altitude
	startPos.Altitude = 550000

	typeDef := &PlatformTypeDefinition{
		Class: "Starlink Satellite",
		Performance: PerformanceCharacteristics{
			MaxSpeed:     7590, // m/s
			CruiseSpeed:  7590,
			Acceleration: 0,
			Range:        0,
		},
		Physical: PhysicalCharacteristics{
			Length: 2.8,  // meters
			Width:  1.9,  // meters
			Height: 0.32, // meters
			Mass:   260,  // kg
		},
	}

	config := &PlatformConfiguration{
		ID:   id,
		Name: fmt.Sprintf("Starlink-%s", satelliteNumber),
		Type: "Starlink Satellite",
	}

	universalPlatform := UniversalPlatform{
		ID:           id,
		PlatformType: PlatformTypeSpace,
		TypeDef:      typeDef,
		Config:       config,
		CallSign:     fmt.Sprintf("STARLINK%s", id[len(id)-3:]),
		State: PlatformState{
			ID:          id,
			Position:    startPos,
			Velocity:    Velocity{},
			Heading:     90,
			Speed:       7590,
			LastUpdated: time.Now(),
		},
	}

	return &SpacePlatform{
		UniversalPlatform: universalPlatform,
		OrbitalPeriod:     5760,   // seconds (96 minutes)
		Apogee:            550000, // meters
		Perigee:           550000, // meters
		Inclination:       53.0,   // degrees
		Eccentricity:      0.0,    // circular orbit
		MissionType:       "communication",
		LaunchDate:        time.Now(),
		MissionDuration:   time.Hour * 24 * 365 * 5, // 5 years
		SolarPanelArea:    8.0,                      // square meters
		PowerGeneration:   4000,                     // watts
		FuelRemaining:     50,                       // kg
		RadiationShield:   false,
	}
}

// NewGPSSatellite creates a GPS navigation satellite
func NewGPSSatellite(id, satelliteNumber string, startPos Position) *SpacePlatform {
	// GPS operates at ~20,200 km altitude
	startPos.Altitude = 20200000

	typeDef := &PlatformTypeDefinition{
		Class: "GPS Block III",
		Performance: PerformanceCharacteristics{
			MaxSpeed:     3870, // m/s
			CruiseSpeed:  3870,
			Acceleration: 0,
			Range:        0,
		},
		Physical: PhysicalCharacteristics{
			Length: 3.0,  // meters
			Width:  2.0,  // meters
			Height: 1.7,  // meters
			Mass:   2000, // kg
		},
	}

	config := &PlatformConfiguration{
		ID:   id,
		Name: fmt.Sprintf("GPS III-%s", satelliteNumber),
		Type: "GPS Block III",
	}

	universalPlatform := UniversalPlatform{
		ID:           id,
		PlatformType: PlatformTypeSpace,
		TypeDef:      typeDef,
		Config:       config,
		CallSign:     fmt.Sprintf("GPS%s", satelliteNumber),
		State: PlatformState{
			ID:          id,
			Position:    startPos,
			Velocity:    Velocity{},
			Heading:     90,
			Speed:       3870,
			LastUpdated: time.Now(),
		},
	}

	return &SpacePlatform{
		UniversalPlatform: universalPlatform,
		OrbitalPeriod:     43200,    // seconds (12 hours)
		Apogee:            20200000, // meters
		Perigee:           20200000, // meters
		Inclination:       55.0,     // degrees
		Eccentricity:      0.0,      // circular orbit
		MissionType:       "navigation",
		LaunchDate:        time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC),
		MissionDuration:   time.Hour * 24 * 365 * 15, // 15 years
		SolarPanelArea:    22.0,                      // square meters
		PowerGeneration:   2900,                      // watts
		FuelRemaining:     200,                       // kg
		RadiationShield:   true,
	}
}

// NewHubbleTelescope creates the Hubble Space Telescope
func NewHubbleTelescope(id string, startPos Position) *SpacePlatform {
	// Hubble orbits at ~547 km altitude
	startPos.Altitude = 547000

	typeDef := &PlatformTypeDefinition{
		Class: "Space Telescope",
		Performance: PerformanceCharacteristics{
			MaxSpeed:     7590, // m/s
			CruiseSpeed:  7590,
			Acceleration: 0,
			Range:        0,
		},
		Physical: PhysicalCharacteristics{
			Length: 13.3,  // meters
			Width:  4.3,   // meters (diameter)
			Height: 4.3,   // meters
			Mass:   11110, // kg
		},
	}

	config := &PlatformConfiguration{
		ID:   id,
		Name: "Hubble Space Telescope",
		Type: "Space Telescope",
	}

	universalPlatform := UniversalPlatform{
		ID:           id,
		PlatformType: PlatformTypeSpace,
		TypeDef:      typeDef,
		Config:       config,
		CallSign:     "HUBBLE",
		State: PlatformState{
			ID:          id,
			Position:    startPos,
			Velocity:    Velocity{},
			Heading:     90,
			Speed:       7590,
			LastUpdated: time.Now(),
		},
	}

	return &SpacePlatform{
		UniversalPlatform: universalPlatform,
		OrbitalPeriod:     5730,   // seconds (95.5 minutes)
		Apogee:            547000, // meters
		Perigee:           547000, // meters
		Inclination:       28.5,   // degrees
		Eccentricity:      0.0003, // nearly circular
		MissionType:       "observation",
		LaunchDate:        time.Date(1990, 4, 24, 0, 0, 0, 0, time.UTC),
		MissionDuration:   time.Hour * 24 * 365 * 35, // extended mission
		SolarPanelArea:    50.0,                      // square meters
		PowerGeneration:   2800,                      // watts
		FuelRemaining:     0,                         // no fuel, uses reaction wheels
		RadiationShield:   true,
	}
}

// NewDragonCapsule creates a SpaceX Dragon spacecraft
func NewDragonCapsule(id, missionName string, startPos Position) *SpacePlatform {
	// Dragon typically operates at ISS altitude
	startPos.Altitude = 408000

	typeDef := &PlatformTypeDefinition{
		Class: "Dragon 2 Capsule",
		Performance: PerformanceCharacteristics{
			MaxSpeed:     7660, // m/s
			CruiseSpeed:  7660,
			Acceleration: 0,
			Range:        0,
		},
		Physical: PhysicalCharacteristics{
			Length: 8.1,   // meters (with trunk)
			Width:  3.7,   // meters (diameter)
			Height: 3.7,   // meters
			Mass:   12055, // kg (fully loaded)
		},
	}

	config := &PlatformConfiguration{
		ID:   id,
		Name: fmt.Sprintf("Dragon %s", missionName),
		Type: "Dragon 2 Capsule",
	}

	universalPlatform := UniversalPlatform{
		ID:           id,
		PlatformType: PlatformTypeSpace,
		TypeDef:      typeDef,
		Config:       config,
		CallSign:     fmt.Sprintf("DRAGON%s", id[len(id)-2:]),
		State: PlatformState{
			ID:          id,
			Position:    startPos,
			Velocity:    Velocity{},
			Heading:     90,
			Speed:       7660,
			LastUpdated: time.Now(),
		},
	}

	return &SpacePlatform{
		UniversalPlatform: universalPlatform,
		OrbitalPeriod:     5520,   // seconds
		Apogee:            420000, // meters
		Perigee:           408000, // meters
		Inclination:       51.6,   // degrees
		Eccentricity:      0.001,  // nearly circular
		MissionType:       "crew transport",
		LaunchDate:        time.Now(),
		MissionDuration:   time.Hour * 24 * 180, // 6 months typical
		SolarPanelArea:    15.0,                 // square meters
		PowerGeneration:   3000,                 // watts
		FuelRemaining:     400,                  // kg
		RadiationShield:   true,
	}
}
