package models

import (
	"fmt"
	"math"
	"time"
)

// MaritimePlatform represents naval and commercial vessels
type MaritimePlatform struct {
	// Base identification
	ID       string
	Class    string // Ship class/model
	Name     string // Ship name
	CallSign string

	// Current state
	State PlatformState

	// Maritime characteristics
	MaxSpeed     float64 // m/s
	CruiseSpeed  float64 // m/s
	Draft        float64 // meters (depth below waterline)
	Displacement float64 // tonnes
	Range        float64 // nautical miles (required by tests)

	// Physical characteristics
	Length float64 // meters
	Width  float64 // beam in meters
	Height float64 // meters (above waterline)
	Mass   float64 // kg

	// Navigation
	Destination *Position
	Route       []Position
}

// Core Platform interface implementation
func (m *MaritimePlatform) GetID() string           { return m.ID }
func (m *MaritimePlatform) GetType() PlatformType   { return PlatformTypeMaritime }
func (m *MaritimePlatform) GetClass() string        { return m.Class }
func (m *MaritimePlatform) GetName() string         { return m.Name }
func (m *MaritimePlatform) GetCallSign() string     { return m.CallSign }
func (m *MaritimePlatform) GetState() PlatformState { return m.State }
func (m *MaritimePlatform) GetMaxSpeed() float64    { return m.MaxSpeed }
func (m *MaritimePlatform) GetMaxAltitude() float64 { return 0 } // Ships don't fly
func (m *MaritimePlatform) GetLength() float64      { return m.Length }
func (m *MaritimePlatform) GetWidth() float64       { return m.Width }
func (m *MaritimePlatform) GetHeight() float64      { return m.Height }
func (m *MaritimePlatform) GetMass() float64        { return m.Mass }

func (m *MaritimePlatform) UpdateState(state PlatformState) {
	m.State = state
	// Keep ships at sea level
	m.State.Position.Altitude = 0
}

func (m *MaritimePlatform) SetDestination(pos Position) error {
	// Force maritime platforms to sea level
	pos.Altitude = 0
	m.Destination = &pos
	return nil
}

func (m *MaritimePlatform) Update(deltaTime time.Duration) error {
	if m.Destination == nil {
		return nil
	}

	// Simple movement towards destination
	dt := deltaTime.Seconds()

	// Calculate distance and bearing to destination
	deltaLat := m.Destination.Latitude - m.State.Position.Latitude
	deltaLon := m.Destination.Longitude - m.State.Position.Longitude

	distance := math.Sqrt(deltaLat*deltaLat + deltaLon*deltaLon)

	if distance < 0.0001 { // Close enough (smaller threshold for ships)
		m.Destination = nil
		return nil
	}

	// Move towards destination at cruise speed
	speed := math.Min(m.CruiseSpeed, m.MaxSpeed)

	// Update position
	factor := (speed * dt) / (distance * 111320) // rough meters per degree
	m.State.Position.Latitude += deltaLat * factor
	m.State.Position.Longitude += deltaLon * factor
	m.State.Position.Altitude = 0 // Always at sea level

	// Update heading and speed
	m.State.Heading = math.Atan2(deltaLon, deltaLat) * 180 / math.Pi
	if m.State.Heading < 0 {
		m.State.Heading += 360
	}
	m.State.Speed = speed
	m.State.LastUpdated = time.Now()

	return nil
}

// Ship factory functions for real-world platforms

// NewArleighBurkeDestroyer creates an Arleigh Burke-class destroyer (US Navy)
func NewArleighBurkeDestroyer(id, shipName string, startPos Position) *MaritimePlatform {
	startPos.Altitude = 0 // Ensure at sea level
	return &MaritimePlatform{
		ID:       id,
		Class:    "Arleigh Burke-class",
		Name:     fmt.Sprintf("USS %s", shipName),
		CallSign: fmt.Sprintf("NAVY%s", id[len(id)-3:]),
		State: PlatformState{
			ID:          id,
			Position:    startPos,
			Velocity:    Velocity{},
			Heading:     0,
			Speed:       0,
			LastUpdated: time.Now(),
		},
		MaxSpeed:     15.4,    // m/s (30+ knots)
		CruiseSpeed:  10.3,    // m/s (20 knots)
		Draft:        6.3,     // meters
		Displacement: 9200,    // tonnes
		Length:       155,     // meters
		Width:        20,      // meters
		Height:       18,      // meters
		Mass:         9200000, // kg
		Range:        4000,    // nautical miles
	}
}

// NewTiconderogaCruiser creates a Ticonderoga-class cruiser (US Navy)
func NewTiconderogaCruiser(id, shipName string, startPos Position) *MaritimePlatform {
	startPos.Altitude = 0
	return &MaritimePlatform{
		ID:       id,
		Class:    "Ticonderoga-class",
		Name:     fmt.Sprintf("USS %s", shipName),
		CallSign: fmt.Sprintf("NAVY%s", id[len(id)-3:]),
		State: PlatformState{
			ID:          id,
			Position:    startPos,
			Velocity:    Velocity{},
			Heading:     0,
			Speed:       0,
			LastUpdated: time.Now(),
		},
		MaxSpeed:     15.4,    // m/s (30+ knots)
		CruiseSpeed:  10.3,    // m/s (20 knots)
		Draft:        10.2,    // meters
		Displacement: 9800,    // tonnes
		Length:       173,     // meters
		Width:        16.8,    // meters
		Height:       20,      // meters
		Mass:         9800000, // kg
		Range:        5000,    // nautical miles
	}
}

// NewContainerShip creates a large container vessel (commercial)
func NewContainerShip(id, shipName string, startPos Position) *MaritimePlatform {
	startPos.Altitude = 0
	return &MaritimePlatform{
		ID:       id,
		Class:    "Ultra Large Container Vessel",
		Name:     shipName,
		CallSign: fmt.Sprintf("CARGO%s", id[len(id)-3:]),
		State: PlatformState{
			ID:          id,
			Position:    startPos,
			Velocity:    Velocity{},
			Heading:     0,
			Speed:       0,
			LastUpdated: time.Now(),
		},
		MaxSpeed:     12.9,      // m/s (25 knots)
		CruiseSpeed:  10.3,      // m/s (20 knots)
		Draft:        16,        // meters
		Displacement: 200000,    // tonnes
		Length:       400,       // meters
		Width:        59,        // meters
		Height:       73,        // meters
		Mass:         200000000, // kg
		Range:        10000,     // nautical miles
	}
}

// NewOilTanker creates a large crude oil tanker (commercial)
func NewOilTanker(id, shipName string, startPos Position) *MaritimePlatform {
	startPos.Altitude = 0
	return &MaritimePlatform{
		ID:       id,
		Class:    "Very Large Crude Carrier",
		Name:     shipName,
		CallSign: fmt.Sprintf("TANKER%s", id[len(id)-2:]),
		State: PlatformState{
			ID:          id,
			Position:    startPos,
			Velocity:    Velocity{},
			Heading:     0,
			Speed:       0,
			LastUpdated: time.Now(),
		},
		MaxSpeed:     8.2,       // m/s (16 knots)
		CruiseSpeed:  6.7,       // m/s (13 knots)
		Draft:        22,        // meters
		Displacement: 320000,    // tonnes
		Length:       330,       // meters
		Width:        60,        // meters
		Height:       35,        // meters
		Mass:         320000000, // kg
		Range:        12000,     // nautical miles
	}
}

// NewCoastGuardCutter creates a Coast Guard cutter
func NewCoastGuardCutter(id, shipName string, startPos Position) *MaritimePlatform {
	startPos.Altitude = 0
	return &MaritimePlatform{
		ID:       id,
		Class:    "Legend-class Cutter",
		Name:     fmt.Sprintf("USCGC %s", shipName),
		CallSign: fmt.Sprintf("COASTGUARD%s", id[len(id)-2:]),
		State: PlatformState{
			ID:          id,
			Position:    startPos,
			Velocity:    Velocity{},
			Heading:     0,
			Speed:       0,
			LastUpdated: time.Now(),
		},
		MaxSpeed:     14.4,    // m/s (28 knots)
		CruiseSpeed:  10.3,    // m/s (20 knots)
		Draft:        6.7,     // meters
		Displacement: 4500,    // tonnes
		Length:       127,     // meters
		Width:        16.4,    // meters
		Height:       15,      // meters
		Mass:         4500000, // kg
		Range:        3000,    // nautical miles
	}
}
