package models

import (
	"time"
)

// PlatformType represents the category of platform
type PlatformType string

const (
	PlatformTypeAirborne PlatformType = "airborne"
	PlatformTypeMaritime PlatformType = "maritime"
	PlatformTypeLand     PlatformType = "land"
	PlatformTypeSpace    PlatformType = "space"
)

// Position represents a 3D position in space
type Position struct {
	Latitude  float64 `json:"latitude"`  // degrees
	Longitude float64 `json:"longitude"` // degrees
	Altitude  float64 `json:"altitude"`  // meters above sea level
}

// Velocity represents 3D velocity vector
type Velocity struct {
	North float64 `json:"north"` // m/s
	East  float64 `json:"east"`  // m/s
	Up    float64 `json:"up"`    // m/s
}

// PlatformState represents the current state of a platform
type PlatformState struct {
	ID          string    `json:"id"`
	Position    Position  `json:"position"`
	Velocity    Velocity  `json:"velocity"`
	Heading     float64   `json:"heading"`     // degrees, 0-360
	Speed       float64   `json:"speed"`       // m/s
	LastUpdated time.Time `json:"lastUpdated"`
}

// Platform interface defines the contract for all platform types
type Platform interface {
	// Core identification
	GetID() string
	GetType() PlatformType
	GetClass() string // e.g., "Boeing 737-800", "Arleigh Burke-class"
	GetName() string  // e.g., "United 1234", "USS Cole"
	
	// State management
	GetState() PlatformState
	UpdateState(state PlatformState)
	
	// Behavior
	Update(deltaTime time.Duration) error
	SetDestination(pos Position) error
	
	// Properties
	GetMaxSpeed() float64    // m/s
	GetMaxAltitude() float64 // meters (for applicable platforms)
	GetCallSign() string
	
	// Real-world characteristics
	GetLength() float64  // meters
	GetWidth() float64   // meters
	GetHeight() float64  // meters
	GetMass() float64    // kg
}