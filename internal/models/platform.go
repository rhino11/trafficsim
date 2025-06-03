package models

import (
	"fmt"
	"math"
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
	Heading     float64   `json:"heading"` // degrees, 0-360
	Speed       float64   `json:"speed"`   // m/s
	LastUpdated time.Time `json:"lastUpdated"`
}

// PerformanceCharacteristics holds configurable performance data
type PerformanceCharacteristics struct {
	MaxSpeed        float64 `yaml:"max_speed"`
	CruiseSpeed     float64 `yaml:"cruise_speed"`
	MaxAltitude     float64 `yaml:"max_altitude,omitempty"`
	FuelConsumption float64 `yaml:"fuel_consumption"`
	TurningRadius   float64 `yaml:"turning_radius,omitempty"`
	Acceleration    float64 `yaml:"acceleration"`
	MaxGradient     float64 `yaml:"max_gradient,omitempty"`
	ClimbRate       float64 `yaml:"climb_rate,omitempty"`
	StallSpeed      float64 `yaml:"stall_speed,omitempty"`

	// Orbital characteristics for space platforms
	OrbitalVelocity float64 `yaml:"orbital_velocity,omitempty"`
	OrbitalPeriod   float64 `yaml:"orbital_period,omitempty"`
	OrbitalAltitude float64 `yaml:"orbital_altitude,omitempty"`
	Inclination     float64 `yaml:"inclination,omitempty"`
	Eccentricity    float64 `yaml:"eccentricity,omitempty"`
}

// PhysicalCharacteristics holds physical properties
type PhysicalCharacteristics struct {
	Length          float64 `yaml:"length"`
	Width           float64 `yaml:"width"`
	Height          float64 `yaml:"height"`
	Mass            float64 `yaml:"mass"`
	FuelCapacity    float64 `yaml:"fuel_capacity,omitempty"`
	Draft           float64 `yaml:"draft,omitempty"`            // Maritime
	GroundClearance float64 `yaml:"ground_clearance,omitempty"` // Land
	SolarPanelArea  float64 `yaml:"solar_panel_area,omitempty"` // Space
}

// OperationalCharacteristics holds mission and operational data
type OperationalCharacteristics struct {
	Range             float64  `yaml:"range"`
	CrewCapacity      int      `yaml:"crew_capacity,omitempty"`
	PassengerCapacity int      `yaml:"passenger_capacity,omitempty"`
	CargoCapacity     float64  `yaml:"cargo_capacity,omitempty"`
	MissionLife       float64  `yaml:"mission_life,omitempty"`    // Space platforms
	SensorRange       float64  `yaml:"sensor_range,omitempty"`    // Naval platforms
	WeaponSystems     []string `yaml:"weapon_systems,omitempty"`  // Military platforms
	FrequencyBands    []string `yaml:"frequency_bands,omitempty"` // Satellites
}

// CallsignConfiguration defines how callsigns are generated
type CallsignConfiguration struct {
	Prefix      string   `yaml:"prefix"`
	Format      string   `yaml:"format"`
	NumberRange [2]int   `yaml:"number_range,omitempty"`
	Names       []string `yaml:"names,omitempty"`
	Modules     []string `yaml:"modules,omitempty"`
}

// PlatformTypeDefinition defines the characteristics of a platform type
type PlatformTypeDefinition struct {
	Class        string                     `yaml:"class"`
	Category     string                     `yaml:"category"`
	Performance  PerformanceCharacteristics `yaml:"performance"`
	Physical     PhysicalCharacteristics    `yaml:"physical"`
	Operational  OperationalCharacteristics `yaml:"operational"`
	CallsignConf CallsignConfiguration      `yaml:"callsign_config"`
}

// MissionConfiguration defines platform-specific mission parameters
type MissionConfiguration struct {
	Type       string                 `yaml:"type"`
	Parameters map[string]interface{} `yaml:",inline"`
}

// PlatformConfiguration represents a configured platform instance
type PlatformConfiguration struct {
	ID            string               `yaml:"id"`
	Type          string               `yaml:"type"`
	Name          string               `yaml:"name"`
	StartPosition Position             `yaml:"start_position"`
	Mission       MissionConfiguration `yaml:"mission"`
}

// UniversalPlatform implements the Platform interface using configuration data
type UniversalPlatform struct {
	ID           string
	PlatformType PlatformType
	TypeDef      *PlatformTypeDefinition
	Config       *PlatformConfiguration
	State        PlatformState
	CallSign     string

	// Navigation
	Destination *Position
	Route       []Position

	// Runtime state
	FuelRemaining float64
	MissionTime   time.Duration

	// Physics state
	lastPosition Position
	acceleration float64
}

// Platform interface implementation
func (up *UniversalPlatform) GetID() string {
	return up.ID
}

func (up *UniversalPlatform) GetType() PlatformType {
	return up.PlatformType
}

func (up *UniversalPlatform) GetClass() string {
	return up.TypeDef.Class
}

func (up *UniversalPlatform) GetName() string {
	return up.Config.Name
}

func (up *UniversalPlatform) GetCallSign() string {
	return up.CallSign
}

func (up *UniversalPlatform) GetState() PlatformState {
	return up.State
}

func (up *UniversalPlatform) UpdateState(state PlatformState) {
	up.State = state
}

func (up *UniversalPlatform) GetMaxSpeed() float64 {
	return up.TypeDef.Performance.MaxSpeed
}

func (up *UniversalPlatform) GetMaxAltitude() float64 {
	if up.TypeDef.Performance.OrbitalAltitude > 0 {
		return up.TypeDef.Performance.OrbitalAltitude
	}
	return up.TypeDef.Performance.MaxAltitude
}

func (up *UniversalPlatform) GetLength() float64 {
	return up.TypeDef.Physical.Length
}

func (up *UniversalPlatform) GetWidth() float64 {
	return up.TypeDef.Physical.Width
}

func (up *UniversalPlatform) GetHeight() float64 {
	return up.TypeDef.Physical.Height
}

func (up *UniversalPlatform) GetMass() float64 {
	return up.TypeDef.Physical.Mass
}

func (up *UniversalPlatform) SetDestination(pos Position) error {
	up.Destination = &pos
	return nil
}

// GetPerformanceCharacteristic allows access to any performance parameter
func (up *UniversalPlatform) GetPerformanceCharacteristic(name string) (float64, error) {
	switch name {
	case "cruise_speed":
		return up.TypeDef.Performance.CruiseSpeed, nil
	case "fuel_consumption":
		return up.TypeDef.Performance.FuelConsumption, nil
	case "turning_radius":
		return up.TypeDef.Performance.TurningRadius, nil
	case "acceleration":
		return up.TypeDef.Performance.Acceleration, nil
	case "climb_rate":
		return up.TypeDef.Performance.ClimbRate, nil
	case "orbital_period":
		return up.TypeDef.Performance.OrbitalPeriod, nil
	default:
		return 0, fmt.Errorf("unknown performance characteristic: %s", name)
	}
}

// Update implements platform-specific movement and behavior with improved physics
func (up *UniversalPlatform) Update(deltaTime time.Duration) error {
	up.MissionTime += deltaTime
	deltaSeconds := deltaTime.Seconds()

	// Store previous position for velocity calculation
	up.lastPosition = up.State.Position

	// Update fuel consumption
	fuelConsumed := up.TypeDef.Performance.FuelConsumption * deltaSeconds
	up.FuelRemaining = max(0, up.FuelRemaining-fuelConsumed)

	// Platform-specific movement logic based on type
	switch up.PlatformType {
	case PlatformTypeAirborne:
		return up.updateAirborneMovement(deltaSeconds)
	case PlatformTypeMaritime:
		return up.updateMaritimeMovement(deltaSeconds)
	case PlatformTypeLand:
		return up.updateLandMovement(deltaSeconds)
	case PlatformTypeSpace:
		return up.updateSpaceMovement(deltaSeconds)
	default:
		return fmt.Errorf("unknown platform type: %s", up.PlatformType)
	}
}

func (up *UniversalPlatform) updateAirborneMovement(deltaSeconds float64) error {
	if up.Destination == nil {
		return nil
	}

	// Enhanced aircraft movement with realistic climb rates and banking
	distance := up.calculateGreatCircleDistance(*up.Destination)
	if distance < 100 { // 100 meter threshold
		up.State.Position = *up.Destination
		up.Destination = nil
		up.State.Speed = 0
		return nil
	}

	// Calculate desired heading
	desiredHeading := up.calculateBearing(*up.Destination)

	// Apply turning constraints
	if up.TypeDef.Performance.TurningRadius > 0 {
		up.applyTurningConstraints(desiredHeading, deltaSeconds)
	} else {
		up.State.Heading = desiredHeading
	}

	// Calculate altitude change with climb rate constraints
	altitudeDiff := up.Destination.Altitude - up.State.Position.Altitude
	if math.Abs(altitudeDiff) > 10 {
		climbRate := up.TypeDef.Performance.ClimbRate
		if climbRate == 0 {
			climbRate = 10.0 // Default 10 m/s
		}

		maxAltChange := climbRate * deltaSeconds
		if math.Abs(altitudeDiff) <= maxAltChange {
			up.State.Position.Altitude = up.Destination.Altitude
		} else if altitudeDiff > 0 {
			up.State.Position.Altitude += maxAltChange
		} else {
			up.State.Position.Altitude -= maxAltChange
		}
	}

	// Apply acceleration constraints
	targetSpeed := up.TypeDef.Performance.CruiseSpeed
	if up.TypeDef.Performance.Acceleration > 0 {
		up.applyAccelerationConstraints(targetSpeed, deltaSeconds)
	} else {
		up.State.Speed = targetSpeed
	}

	// Update position based on current heading and speed
	up.updatePositionFromHeadingAndSpeed(deltaSeconds)

	up.State.LastUpdated = time.Now()
	return nil
}

func (up *UniversalPlatform) updateMaritimeMovement(deltaSeconds float64) error {
	if up.Destination == nil {
		return nil
	}

	// Maritime movement with inertia and turning constraints
	distance := up.calculateGreatCircleDistance(*up.Destination)
	if distance < 50 { // 50 meter threshold for ships
		up.State.Position = *up.Destination
		up.Destination = nil
		up.State.Speed = 0
		return nil
	}

	// Ships have larger turning radii and slower acceleration
	desiredHeading := up.calculateBearing(*up.Destination)

	// Apply maritime turning constraints (ships turn slower)
	turningRadius := up.TypeDef.Performance.TurningRadius
	if turningRadius == 0 {
		turningRadius = up.TypeDef.Physical.Length * 5 // Default: 5x ship length
	}
	up.applyTurningConstraints(desiredHeading, deltaSeconds)

	// Maritime platforms stay at sea level
	up.State.Position.Altitude = 0

	// Apply acceleration with maritime characteristics
	targetSpeed := up.TypeDef.Performance.CruiseSpeed
	acceleration := up.TypeDef.Performance.Acceleration
	if acceleration == 0 {
		acceleration = 0.5 // Default slow acceleration for ships
	}
	up.applyAccelerationConstraints(targetSpeed, deltaSeconds)

	up.updatePositionFromHeadingAndSpeed(deltaSeconds)
	up.State.LastUpdated = time.Now()
	return nil
}

func (up *UniversalPlatform) updateLandMovement(deltaSeconds float64) error {
	if up.Destination == nil {
		return nil
	}

	// Land movement with terrain and gradient constraints
	distance := up.calculateGreatCircleDistance(*up.Destination)
	if distance < 10 { // 10 meter threshold for land vehicles
		up.State.Position = *up.Destination
		up.Destination = nil
		up.State.Speed = 0
		return nil
	}

	desiredHeading := up.calculateBearing(*up.Destination)
	up.applyTurningConstraints(desiredHeading, deltaSeconds)

	// Apply gradient constraints for altitude changes
	altitudeDiff := up.Destination.Altitude - up.State.Position.Altitude
	horizontalDist := up.calculateHorizontalDistance(*up.Destination)

	if horizontalDist > 0 {
		gradient := math.Atan(altitudeDiff/horizontalDist) * 180.0 / math.Pi
		maxGradient := up.TypeDef.Performance.MaxGradient
		if maxGradient == 0 {
			maxGradient = 30.0 // Default 30 degree max gradient
		}

		if math.Abs(gradient) > maxGradient {
			// Reduce speed when climbing steep grades
			up.State.Speed *= (maxGradient / math.Abs(gradient))
		}
	}

	targetSpeed := up.TypeDef.Performance.CruiseSpeed
	up.applyAccelerationConstraints(targetSpeed, deltaSeconds)

	up.updatePositionFromHeadingAndSpeed(deltaSeconds)
	up.State.LastUpdated = time.Now()
	return nil
}

func (up *UniversalPlatform) updateSpaceMovement(deltaSeconds float64) error {
	// Enhanced orbital mechanics
	if up.TypeDef.Performance.OrbitalPeriod > 0 {
		// Calculate angular velocity
		angularVelocity := 2 * math.Pi / up.TypeDef.Performance.OrbitalPeriod
		deltaAngle := angularVelocity * deltaSeconds

		// Apply orbital inclination
		inclination := up.TypeDef.Performance.Inclination * math.Pi / 180.0

		// Simplified orbital motion (circular orbit approximation)
		currentLongitude := up.State.Position.Longitude * math.Pi / 180.0
		newLongitude := currentLongitude + deltaAngle

		// Apply inclination effect on latitude (simplified)
		maxLatitude := inclination * 180.0 / math.Pi
		timeInOrbit := up.MissionTime.Seconds()
		latitudePhase := 2 * math.Pi * timeInOrbit / up.TypeDef.Performance.OrbitalPeriod
		up.State.Position.Latitude = maxLatitude * math.Sin(latitudePhase)

		up.State.Position.Longitude = newLongitude * 180.0 / math.Pi
		if up.State.Position.Longitude > 180 {
			up.State.Position.Longitude -= 360
		} else if up.State.Position.Longitude < -180 {
			up.State.Position.Longitude += 360
		}

		// Maintain orbital altitude
		up.State.Position.Altitude = up.TypeDef.Performance.OrbitalAltitude
		up.State.Speed = up.TypeDef.Performance.OrbitalVelocity
		up.State.Heading = 90 // Generally eastward
	}

	up.State.LastUpdated = time.Now()
	return nil
}

// Physics helper functions

func (up *UniversalPlatform) calculateGreatCircleDistance(target Position) float64 {
	// Haversine formula for great circle distance
	lat1 := up.State.Position.Latitude * math.Pi / 180.0
	lat2 := target.Latitude * math.Pi / 180.0
	deltaLat := (target.Latitude - up.State.Position.Latitude) * math.Pi / 180.0
	deltaLon := (target.Longitude - up.State.Position.Longitude) * math.Pi / 180.0

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	earthRadius := 6371000.0 // meters
	distance := earthRadius * c

	// Add altitude difference
	altDiff := target.Altitude - up.State.Position.Altitude
	return math.Sqrt(distance*distance + altDiff*altDiff)
}

func (up *UniversalPlatform) calculateHorizontalDistance(target Position) float64 {
	lat1 := up.State.Position.Latitude * math.Pi / 180.0
	lat2 := target.Latitude * math.Pi / 180.0
	deltaLat := (target.Latitude - up.State.Position.Latitude) * math.Pi / 180.0
	deltaLon := (target.Longitude - up.State.Position.Longitude) * math.Pi / 180.0

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	earthRadius := 6371000.0 // meters
	return earthRadius * c
}

func (up *UniversalPlatform) calculateBearing(target Position) float64 {
	lat1 := up.State.Position.Latitude * math.Pi / 180.0
	lat2 := target.Latitude * math.Pi / 180.0
	deltaLon := (target.Longitude - up.State.Position.Longitude) * math.Pi / 180.0

	y := math.Sin(deltaLon) * math.Cos(lat2)
	x := math.Cos(lat1)*math.Sin(lat2) - math.Sin(lat1)*math.Cos(lat2)*math.Cos(deltaLon)

	bearing := math.Atan2(y, x) * 180.0 / math.Pi
	return math.Mod(bearing+360, 360) // Normalize to 0-360
}

func (up *UniversalPlatform) applyTurningConstraints(desiredHeading, deltaSeconds float64) {
	currentHeading := up.State.Heading
	headingDiff := desiredHeading - currentHeading

	// Normalize heading difference to -180 to 180
	for headingDiff > 180 {
		headingDiff -= 360
	}
	for headingDiff < -180 {
		headingDiff += 360
	}

	// Calculate maximum turn rate based on turning radius and speed
	if up.TypeDef.Performance.TurningRadius > 0 && up.State.Speed > 0 {
		maxTurnRate := (up.State.Speed / up.TypeDef.Performance.TurningRadius) * 180.0 / math.Pi // deg/s
		maxTurnChange := maxTurnRate * deltaSeconds

		if math.Abs(headingDiff) <= maxTurnChange {
			up.State.Heading = desiredHeading
		} else if headingDiff > 0 {
			up.State.Heading = math.Mod(currentHeading+maxTurnChange, 360)
		} else {
			up.State.Heading = math.Mod(currentHeading-maxTurnChange+360, 360)
		}
	} else {
		up.State.Heading = desiredHeading
	}
}

func (up *UniversalPlatform) applyAccelerationConstraints(targetSpeed, deltaSeconds float64) {
	speedDiff := targetSpeed - up.State.Speed

	if up.TypeDef.Performance.Acceleration > 0 {
		maxSpeedChange := up.TypeDef.Performance.Acceleration * deltaSeconds

		if math.Abs(speedDiff) <= maxSpeedChange {
			up.State.Speed = targetSpeed
		} else if speedDiff > 0 {
			up.State.Speed += maxSpeedChange
		} else {
			up.State.Speed -= maxSpeedChange
		}
	} else {
		up.State.Speed = targetSpeed
	}

	// Ensure speed doesn't exceed maximum
	if up.State.Speed > up.TypeDef.Performance.MaxSpeed {
		up.State.Speed = up.TypeDef.Performance.MaxSpeed
	}

	// Ensure speed isn't negative
	if up.State.Speed < 0 {
		up.State.Speed = 0
	}
}

func (up *UniversalPlatform) updatePositionFromHeadingAndSpeed(deltaSeconds float64) {
	// Convert heading to radians (0 degrees = North, 90 degrees = East)
	headingRad := (90 - up.State.Heading) * math.Pi / 180.0

	// Calculate distance moved
	distance := up.State.Speed * deltaSeconds

	// Earth radius in meters
	earthRadius := 6371000.0

	// Calculate new position
	deltaLat := distance * math.Cos(headingRad) / earthRadius * 180.0 / math.Pi
	deltaLon := distance * math.Sin(headingRad) / earthRadius * 180.0 / math.Pi / math.Cos(up.State.Position.Latitude*math.Pi/180.0)

	up.State.Position.Latitude += deltaLat
	up.State.Position.Longitude += deltaLon

	// Calculate velocity components for state tracking
	up.State.Velocity.North = up.State.Speed * math.Cos(headingRad)
	up.State.Velocity.East = up.State.Speed * math.Sin(headingRad)
	up.State.Velocity.Up = 0 // Will be calculated separately for altitude changes
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
	GetLength() float64 // meters
	GetWidth() float64  // meters
	GetHeight() float64 // meters
	GetMass() float64   // kg
}
