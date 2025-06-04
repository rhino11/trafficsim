package models

import "time"

// PlatformInterface defines the common contract for all platform types
type PlatformInterface interface {
	// Core identification
	GetID() string
	GetType() string
	GetCategory() string
	GetClass() string

	// Position and movement
	GetPosition() Position
	SetPosition(Position)
	GetVelocity() Velocity
	SetVelocity(Velocity)
	GetHeading() float64
	SetHeading(float64)

	// Physical properties
	GetMass() float64
	GetDimensions() Dimensions
	GetFuelLevel() float64
	GetMaxFuelCapacity() float64

	// Performance characteristics
	GetMaxSpeed() float64
	GetCruiseSpeed() float64
	GetTurningRadius() float64
	GetAcceleration() float64

	// Operational capabilities
	GetRange() float64
	GetCrewCapacity() int
	GetWeaponSystems() []string

	// Simulation methods
	Update(deltaTime time.Duration) error
	CalculateFuelConsumption(speed float64) float64
	IsOperational() bool
	GetStatus() PlatformStatus

	// Domain-specific methods (to be implemented by each domain)
	GetDomainSpecificData() interface{}
	ValidateConfiguration() error
}

// Common structs used across all platforms
type Position struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude"` // meters above sea level (0 for surface/land platforms)
}

type Velocity struct {
	North float64 `json:"north"` // m/s
	East  float64 `json:"east"`  // m/s
	Up    float64 `json:"up"`    // m/s (typically 0 for surface platforms)
}

// Dimensions struct for platform physical dimensions
type Dimensions struct {
	Length float64 `json:"length"` // meters
	Width  float64 `json:"width"`  // meters (beam for ships, wingspan for aircraft)
	Height float64 `json:"height"` // meters
}
