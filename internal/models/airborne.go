package models

import (
	"fmt"
	"math"
	"time"
)

// AirbornePlatform represents aircraft platforms
type AirbornePlatform struct {
	// Base identification
	ID       string
	Class    string // Aircraft model
	Name     string // Flight number/tail number
	CallSign string

	// Current state
	State PlatformState

	// Flight characteristics
	MaxSpeed    float64 // m/s
	MaxAltitude float64 // meters
	CruiseSpeed float64 // m/s
	CruiseAlt   float64 // meters
	ServiceCeil float64 // meters

	// Physical characteristics
	Length float64 // meters
	Width  float64 // wingspan in meters
	Height float64 // meters
	Mass   float64 // kg

	// Navigation
	Destination *Position
	Route       []Position
}

// Core Platform interface implementation
func (a *AirbornePlatform) GetID() string           { return a.ID }
func (a *AirbornePlatform) GetType() PlatformType   { return PlatformTypeAirborne }
func (a *AirbornePlatform) GetClass() string        { return a.Class }
func (a *AirbornePlatform) GetName() string         { return a.Name }
func (a *AirbornePlatform) GetCallSign() string     { return a.CallSign }
func (a *AirbornePlatform) GetState() PlatformState { return a.State }
func (a *AirbornePlatform) GetMaxSpeed() float64    { return a.MaxSpeed }
func (a *AirbornePlatform) GetMaxAltitude() float64 { return a.MaxAltitude }
func (a *AirbornePlatform) GetLength() float64      { return a.Length }
func (a *AirbornePlatform) GetWidth() float64       { return a.Width }
func (a *AirbornePlatform) GetHeight() float64      { return a.Height }
func (a *AirbornePlatform) GetMass() float64        { return a.Mass }

func (a *AirbornePlatform) UpdateState(state PlatformState) {
	a.State = state
}

func (a *AirbornePlatform) SetDestination(pos Position) error {
	a.Destination = &pos
	return nil
}

func (a *AirbornePlatform) Update(deltaTime time.Duration) error {
	if a.Destination == nil {
		return nil
	}

	// Simple movement towards destination
	dt := deltaTime.Seconds()

	// Calculate distance and bearing to destination
	deltaLat := a.Destination.Latitude - a.State.Position.Latitude
	deltaLon := a.Destination.Longitude - a.State.Position.Longitude
	deltaAlt := a.Destination.Altitude - a.State.Position.Altitude

	distance := math.Sqrt(deltaLat*deltaLat + deltaLon*deltaLon)

	if distance < 0.001 { // Close enough
		a.Destination = nil
		return nil
	}

	// Move towards destination at cruise speed
	speed := math.Min(a.CruiseSpeed, a.MaxSpeed)

	// Update position
	factor := (speed * dt) / (distance * 111320) // rough meters per degree
	a.State.Position.Latitude += deltaLat * factor
	a.State.Position.Longitude += deltaLon * factor

	// Altitude change (simpler vertical movement)
	if math.Abs(deltaAlt) > 10 { // 10 meter threshold
		altChangeRate := 5.0 // m/s climb/descent rate
		altChange := math.Min(math.Abs(deltaAlt), altChangeRate*dt)
		if deltaAlt > 0 {
			a.State.Position.Altitude += altChange
		} else {
			a.State.Position.Altitude -= altChange
		}
	}

	// Update heading and speed
	a.State.Heading = math.Atan2(deltaLon, deltaLat) * 180 / math.Pi
	if a.State.Heading < 0 {
		a.State.Heading += 360
	}
	a.State.Speed = speed
	a.State.LastUpdated = time.Now()

	return nil
}

// Aircraft factory functions for real-world platforms

// NewBoeing737_800 creates a Boeing 737-800 (commercial airliner)
func NewBoeing737_800(id, flightNumber string, startPos Position) *AirbornePlatform {
	return &AirbornePlatform{
		ID:       id,
		Class:    "Boeing 737-800",
		Name:     flightNumber,
		CallSign: flightNumber,
		State: PlatformState{
			ID:          id,
			Position:    startPos,
			Velocity:    Velocity{},
			Heading:     0,
			Speed:       0,
			LastUpdated: time.Now(),
		},
		MaxSpeed:    257,   // m/s (500 kts)
		MaxAltitude: 12500, // meters (41,000 ft)
		CruiseSpeed: 230,   // m/s (447 kts)
		CruiseAlt:   11000, // meters (36,000 ft)
		ServiceCeil: 12500,
		Length:      39.5,  // meters
		Width:       35.8,  // wingspan
		Height:      12.5,  // meters
		Mass:        79010, // kg (max takeoff weight)
	}
}

// NewAirbusA320 creates an Airbus A320 (commercial airliner)
func NewAirbusA320(id, flightNumber string, startPos Position) *AirbornePlatform {
	return &AirbornePlatform{
		ID:       id,
		Class:    "Airbus A320",
		Name:     flightNumber,
		CallSign: flightNumber,
		State: PlatformState{
			ID:          id,
			Position:    startPos,
			Velocity:    Velocity{},
			Heading:     0,
			Speed:       0,
			LastUpdated: time.Now(),
		},
		MaxSpeed:    257,   // m/s (500 kts)
		MaxAltitude: 12000, // meters (39,370 ft)
		CruiseSpeed: 230,   // m/s (447 kts)
		CruiseAlt:   11000, // meters (36,000 ft)
		ServiceCeil: 12000,
		Length:      37.6,  // meters
		Width:       36.0,  // wingspan
		Height:      11.8,  // meters
		Mass:        78000, // kg (max takeoff weight)
	}
}

// NewF16FightingFalcon creates an F-16 Fighting Falcon (military fighter)
func NewF16FightingFalcon(id, tailNumber string, startPos Position) *AirbornePlatform {
	return &AirbornePlatform{
		ID:       id,
		Class:    "F-16 Fighting Falcon",
		Name:     tailNumber,
		CallSign: fmt.Sprintf("VIPER%s", tailNumber[len(tailNumber)-3:]),
		State: PlatformState{
			ID:          id,
			Position:    startPos,
			Velocity:    Velocity{},
			Heading:     0,
			Speed:       0,
			LastUpdated: time.Now(),
		},
		MaxSpeed:    617,   // m/s (Mach 2.0+ at altitude)
		MaxAltitude: 15240, // meters (50,000 ft)
		CruiseSpeed: 257,   // m/s (500 kts)
		CruiseAlt:   9000,  // meters (30,000 ft)
		ServiceCeil: 15240,
		Length:      15.0,  // meters
		Width:       10.0,  // wingspan
		Height:      5.1,   // meters
		Mass:        19187, // kg (max takeoff weight)
	}
}

// NewC130Hercules creates a C-130 Hercules (military transport)
func NewC130Hercules(id, tailNumber string, startPos Position) *AirbornePlatform {
	return &AirbornePlatform{
		ID:       id,
		Class:    "C-130 Hercules",
		Name:     tailNumber,
		CallSign: fmt.Sprintf("HERKY%s", tailNumber[len(tailNumber)-2:]),
		State: PlatformState{
			ID:          id,
			Position:    startPos,
			Velocity:    Velocity{},
			Heading:     0,
			Speed:       0,
			LastUpdated: time.Now(),
		},
		MaxSpeed:    190,   // m/s (370 kts)
		MaxAltitude: 10060, // meters (33,000 ft)
		CruiseSpeed: 160,   // m/s (310 kts)
		CruiseAlt:   7000,  // meters (23,000 ft)
		ServiceCeil: 10060,
		Length:      29.8,  // meters
		Width:       40.4,  // wingspan
		Height:      11.7,  // meters
		Mass:        70300, // kg (max takeoff weight)
	}
}
