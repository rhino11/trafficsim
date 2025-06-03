package models

import (
	"fmt"
	"math"
	"time"
)

// LandPlatform represents ground vehicles and installations
type LandPlatform struct {
	// Base identification
	ID       string
	Class    string // Vehicle/platform type
	Name     string // Unit designation
	CallSign string

	// Current state
	State PlatformState

	// Land characteristics
	MaxSpeed     float64 // m/s
	CruiseSpeed  float64 // m/s
	MaxGradient  float64 // degrees (max slope)
	FuelCapacity float64 // liters
	Range        float64 // meters

	// Physical characteristics
	Length float64 // meters
	Width  float64 // meters
	Height float64 // meters
	Mass   float64 // kg

	// Navigation
	Destination *Position
	Route       []Position
}

// Core Platform interface implementation
func (l *LandPlatform) GetID() string           { return l.ID }
func (l *LandPlatform) GetType() PlatformType   { return PlatformTypeLand }
func (l *LandPlatform) GetClass() string        { return l.Class }
func (l *LandPlatform) GetName() string         { return l.Name }
func (l *LandPlatform) GetCallSign() string     { return l.CallSign }
func (l *LandPlatform) GetState() PlatformState { return l.State }
func (l *LandPlatform) GetMaxSpeed() float64    { return l.MaxSpeed }
func (l *LandPlatform) GetMaxAltitude() float64 { return 0 } // Ground vehicles
func (l *LandPlatform) GetLength() float64      { return l.Length }
func (l *LandPlatform) GetWidth() float64       { return l.Width }
func (l *LandPlatform) GetHeight() float64      { return l.Height }
func (l *LandPlatform) GetMass() float64        { return l.Mass }

func (l *LandPlatform) UpdateState(state PlatformState) {
	l.State = state
}

func (l *LandPlatform) SetDestination(pos Position) error {
	l.Destination = &pos
	return nil
}

func (l *LandPlatform) Update(deltaTime time.Duration) error {
	if l.Destination == nil {
		return nil
	}

	// Simple movement towards destination
	dt := deltaTime.Seconds()

	// Calculate distance and bearing to destination
	deltaLat := l.Destination.Latitude - l.State.Position.Latitude
	deltaLon := l.Destination.Longitude - l.State.Position.Longitude

	distance := math.Sqrt(deltaLat*deltaLat + deltaLon*deltaLon)

	if distance < 0.00001 { // Close enough (very small threshold for ground vehicles)
		l.Destination = nil
		return nil
	}

	// Move towards destination at cruise speed
	speed := math.Min(l.CruiseSpeed, l.MaxSpeed)

	// Update position
	factor := (speed * dt) / (distance * 111320) // rough meters per degree
	l.State.Position.Latitude += deltaLat * factor
	l.State.Position.Longitude += deltaLon * factor
	// Note: Altitude changes would be based on terrain, keeping simple for now

	// Update heading and speed
	l.State.Heading = math.Atan2(deltaLon, deltaLat) * 180 / math.Pi
	if l.State.Heading < 0 {
		l.State.Heading += 360
	}
	l.State.Speed = speed
	l.State.LastUpdated = time.Now()

	return nil
}

// Land platform factory functions for real-world platforms

// NewM1A2Abrams creates an M1A2 Abrams main battle tank
func NewM1A2Abrams(id, unitDesignation string, startPos Position) *LandPlatform {
	return &LandPlatform{
		ID:       id,
		Class:    "M1A2 Abrams MBT",
		Name:     unitDesignation,
		CallSign: fmt.Sprintf("ARMOR%s", id[len(id)-2:]),
		State: PlatformState{
			ID:          id,
			Position:    startPos,
			Velocity:    Velocity{},
			Heading:     0,
			Speed:       0,
			LastUpdated: time.Now(),
		},
		MaxSpeed:     20,     // m/s (45 mph)
		CruiseSpeed:  13.4,   // m/s (30 mph)
		MaxGradient:  30,     // degrees
		FuelCapacity: 1900,   // liters
		Range:        426000, // meters (265 miles)
		Length:       9.8,    // meters
		Width:        3.7,    // meters
		Height:       2.4,    // meters
		Mass:         62000,  // kg
	}
}

// NewHMMWV creates a High Mobility Multipurpose Wheeled Vehicle (Humvee)
func NewHMMWV(id, unitDesignation string, startPos Position) *LandPlatform {
	return &LandPlatform{
		ID:       id,
		Class:    "HMMWV",
		Name:     unitDesignation,
		CallSign: fmt.Sprintf("HUMVEE%s", id[len(id)-2:]),
		State: PlatformState{
			ID:          id,
			Position:    startPos,
			Velocity:    Velocity{},
			Heading:     0,
			Speed:       0,
			LastUpdated: time.Now(),
		},
		MaxSpeed:     31,     // m/s (70 mph)
		CruiseSpeed:  22.4,   // m/s (50 mph)
		MaxGradient:  40,     // degrees
		FuelCapacity: 95,     // liters
		Range:        480000, // meters (300 miles)
		Length:       4.6,    // meters
		Width:        2.2,    // meters
		Height:       1.8,    // meters
		Mass:         5900,   // kg
	}
}

// NewFreightlinerCascadia creates a commercial freight truck
func NewFreightlinerCascadia(id, truckNumber string, startPos Position) *LandPlatform {
	return &LandPlatform{
		ID:       id,
		Class:    "Freightliner Cascadia",
		Name:     truckNumber,
		CallSign: fmt.Sprintf("TRUCK%s", id[len(id)-3:]),
		State: PlatformState{
			ID:          id,
			Position:    startPos,
			Velocity:    Velocity{},
			Heading:     0,
			Speed:       0,
			LastUpdated: time.Now(),
		},
		MaxSpeed:     33.5,    // m/s (75 mph)
		CruiseSpeed:  29.1,    // m/s (65 mph)
		MaxGradient:  15,      // degrees
		FuelCapacity: 1135,    // liters (300 gallons)
		Range:        1600000, // meters (1000 miles)
		Length:       6.1,     // meters (tractor only)
		Width:        2.6,     // meters
		Height:       4.0,     // meters
		Mass:         16000,   // kg (tractor only)
	}
}

// NewPolicePatrolCar creates a police patrol vehicle
func NewPolicePatrolCar(id, unitNumber string, startPos Position) *LandPlatform {
	return &LandPlatform{
		ID:       id,
		Class:    "Ford Police Interceptor",
		Name:     fmt.Sprintf("Unit %s", unitNumber),
		CallSign: fmt.Sprintf("PATROL%s", unitNumber),
		State: PlatformState{
			ID:          id,
			Position:    startPos,
			Velocity:    Velocity{},
			Heading:     0,
			Speed:       0,
			LastUpdated: time.Now(),
		},
		MaxSpeed:     50,     // m/s (112 mph)
		CruiseSpeed:  26.8,   // m/s (60 mph)
		MaxGradient:  25,     // degrees
		FuelCapacity: 68,     // liters (18 gallons)
		Range:        640000, // meters (400 miles)
		Length:       5.2,    // meters
		Width:        1.9,    // meters
		Height:       1.5,    // meters
		Mass:         2000,   // kg
	}
}
