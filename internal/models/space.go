package models

import (
	"fmt"
	"math"
	"time"
)

// SpacePlatform represents satellites and spacecraft
type SpacePlatform struct {
	// Base identification
	ID       string
	Class    string // Spacecraft/satellite type
	Name     string // Mission designation
	CallSign string

	// Current state
	State PlatformState

	// Orbital characteristics
	MaxSpeed      float64 // m/s (orbital velocity)
	OrbitalPeriod float64 // seconds (time for one orbit)
	Apogee        float64 // meters (highest point)
	Perigee       float64 // meters (lowest point)
	Inclination   float64 // degrees (orbital plane angle)

	// Physical characteristics
	Length float64 // meters
	Width  float64 // meters
	Height float64 // meters
	Mass   float64 // kg

	// Navigation (simplified orbital mechanics)
	Destination *Position
	OrbitCenter Position // Usually Earth center for LEO/GEO
}

// Core Platform interface implementation
func (s *SpacePlatform) GetID() string           { return s.ID }
func (s *SpacePlatform) GetType() PlatformType   { return PlatformTypeSpace }
func (s *SpacePlatform) GetClass() string        { return s.Class }
func (s *SpacePlatform) GetName() string         { return s.Name }
func (s *SpacePlatform) GetCallSign() string     { return s.CallSign }
func (s *SpacePlatform) GetState() PlatformState { return s.State }
func (s *SpacePlatform) GetMaxSpeed() float64    { return s.MaxSpeed }
func (s *SpacePlatform) GetMaxAltitude() float64 { return s.Apogee }
func (s *SpacePlatform) GetLength() float64      { return s.Length }
func (s *SpacePlatform) GetWidth() float64       { return s.Width }
func (s *SpacePlatform) GetHeight() float64      { return s.Height }
func (s *SpacePlatform) GetMass() float64        { return s.Mass }

func (s *SpacePlatform) UpdateState(state PlatformState) {
	s.State = state
}

func (s *SpacePlatform) SetDestination(pos Position) error {
	s.Destination = &pos
	return nil
}

func (s *SpacePlatform) Update(deltaTime time.Duration) error {
	// Simplified orbital mechanics - circular orbit approximation
	dt := deltaTime.Seconds()

	// Calculate orbital velocity based on altitude
	earthRadius := 6371000.0 // meters
	altitude := s.State.Position.Altitude
	orbitalRadius := earthRadius + altitude

	// Simplified orbital velocity: v = sqrt(GM/r)
	// Using approximation for Earth: GM ≈ 3.986e14 m³/s²
	GM := 3.986e14
	orbitalVelocity := math.Sqrt(GM / orbitalRadius)

	// Angular velocity (radians per second)
	angularVelocity := orbitalVelocity / orbitalRadius

	// Update position in simplified circular orbit
	// Convert lat/lon to radians for calculation
	lonRad := s.State.Position.Longitude * math.Pi / 180

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
	s.State.Position.Longitude = lonRad * 180 / math.Pi

	// Update speed and heading
	s.State.Speed = orbitalVelocity
	s.State.Heading = 90 // Eastward movement
	s.State.LastUpdated = time.Now()

	return nil
}

// Space platform factory functions for real-world platforms

// NewISSModule creates an International Space Station module
func NewISSModule(id, moduleName string, startPos Position) *SpacePlatform {
	// ISS orbits at approximately 408 km altitude
	startPos.Altitude = 408000
	return &SpacePlatform{
		ID:       id,
		Class:    "ISS Module",
		Name:     fmt.Sprintf("ISS %s", moduleName),
		CallSign: fmt.Sprintf("ISS%s", id[len(id)-2:]),
		State: PlatformState{
			ID:          id,
			Position:    startPos,
			Velocity:    Velocity{},
			Heading:     90,   // Eastward
			Speed:       7660, // m/s (ISS orbital velocity)
			LastUpdated: time.Now(),
		},
		MaxSpeed:      7660,   // m/s
		OrbitalPeriod: 5520,   // seconds (92 minutes)
		Apogee:        420000, // meters
		Perigee:       408000, // meters
		Inclination:   51.6,   // degrees
		Length:        73,     // meters (full ISS)
		Width:         109,    // meters (solar array span)
		Height:        20,     // meters
		Mass:          420000, // kg
	}
}

// NewStarlinkSatellite creates a Starlink communication satellite
func NewStarlinkSatellite(id, satelliteNumber string, startPos Position) *SpacePlatform {
	// Starlink operates at ~550 km altitude
	startPos.Altitude = 550000
	return &SpacePlatform{
		ID:       id,
		Class:    "Starlink Satellite",
		Name:     fmt.Sprintf("Starlink-%s", satelliteNumber),
		CallSign: fmt.Sprintf("STARLINK%s", id[len(id)-3:]),
		State: PlatformState{
			ID:          id,
			Position:    startPos,
			Velocity:    Velocity{},
			Heading:     90,
			Speed:       7590, // m/s
			LastUpdated: time.Now(),
		},
		MaxSpeed:      7590,   // m/s
		OrbitalPeriod: 5760,   // seconds (96 minutes)
		Apogee:        550000, // meters
		Perigee:       550000, // meters (circular orbit)
		Inclination:   53.0,   // degrees
		Length:        2.8,    // meters
		Width:         1.9,    // meters
		Height:        0.32,   // meters
		Mass:          260,    // kg
	}
}

// NewGPSSatellite creates a GPS navigation satellite
func NewGPSSatellite(id, satelliteNumber string, startPos Position) *SpacePlatform {
	// GPS operates at ~20,200 km altitude
	startPos.Altitude = 20200000
	return &SpacePlatform{
		ID:       id,
		Class:    "GPS Block III",
		Name:     fmt.Sprintf("GPS III-%s", satelliteNumber),
		CallSign: fmt.Sprintf("GPS%s", satelliteNumber),
		State: PlatformState{
			ID:          id,
			Position:    startPos,
			Velocity:    Velocity{},
			Heading:     90,
			Speed:       3870, // m/s
			LastUpdated: time.Now(),
		},
		MaxSpeed:      3870,     // m/s
		OrbitalPeriod: 43200,    // seconds (12 hours)
		Apogee:        20200000, // meters
		Perigee:       20200000, // meters
		Inclination:   55.0,     // degrees
		Length:        3.0,      // meters
		Width:         2.0,      // meters
		Height:        1.7,      // meters
		Mass:          2000,     // kg
	}
}

// NewHubbleTelescope creates the Hubble Space Telescope
func NewHubbleTelescope(id string, startPos Position) *SpacePlatform {
	// Hubble orbits at ~547 km altitude
	startPos.Altitude = 547000
	return &SpacePlatform{
		ID:       id,
		Class:    "Space Telescope",
		Name:     "Hubble Space Telescope",
		CallSign: "HUBBLE",
		State: PlatformState{
			ID:          id,
			Position:    startPos,
			Velocity:    Velocity{},
			Heading:     90,
			Speed:       7590, // m/s
			LastUpdated: time.Now(),
		},
		MaxSpeed:      7590,   // m/s
		OrbitalPeriod: 5730,   // seconds (95.5 minutes)
		Apogee:        547000, // meters
		Perigee:       547000, // meters
		Inclination:   28.5,   // degrees
		Length:        13.3,   // meters
		Width:         4.3,    // meters (diameter)
		Height:        4.3,    // meters
		Mass:          11110,  // kg
	}
}

// NewDragonCapsule creates a SpaceX Dragon spacecraft
func NewDragonCapsule(id, missionName string, startPos Position) *SpacePlatform {
	// Dragon typically operates at ISS altitude
	startPos.Altitude = 408000
	return &SpacePlatform{
		ID:       id,
		Class:    "Dragon 2 Capsule",
		Name:     fmt.Sprintf("Dragon %s", missionName),
		CallSign: fmt.Sprintf("DRAGON%s", id[len(id)-2:]),
		State: PlatformState{
			ID:          id,
			Position:    startPos,
			Velocity:    Velocity{},
			Heading:     90,
			Speed:       7660, // m/s
			LastUpdated: time.Now(),
		},
		MaxSpeed:      7660,   // m/s
		OrbitalPeriod: 5520,   // seconds
		Apogee:        420000, // meters
		Perigee:       408000, // meters
		Inclination:   51.6,   // degrees
		Length:        8.1,    // meters (with trunk)
		Width:         3.7,    // meters (diameter)
		Height:        3.7,    // meters
		Mass:          12055,  // kg (fully loaded)
	}
}
