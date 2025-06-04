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

// Acceleration represents 3D acceleration vector
type Acceleration struct {
	North float64 `json:"north"` // m/s²
	East  float64 `json:"east"`  // m/s²
	Up    float64 `json:"up"`    // m/s²
}

// Attitude represents rotational orientation in 3D space
type Attitude struct {
	Roll  float64 `json:"roll"`  // degrees, rotation around longitudinal axis
	Pitch float64 `json:"pitch"` // degrees, rotation around lateral axis
	Yaw   float64 `json:"yaw"`   // degrees, rotation around vertical axis (same as heading)
}

// AngularVelocity represents rotational rates
type AngularVelocity struct {
	RollRate  float64 `json:"roll_rate"`  // degrees/second
	PitchRate float64 `json:"pitch_rate"` // degrees/second
	YawRate   float64 `json:"yaw_rate"`   // degrees/second
}

// AngularAcceleration represents rotational accelerations
type AngularAcceleration struct {
	RollAccel  float64 `json:"roll_accel"`  // degrees/second²
	PitchAccel float64 `json:"pitch_accel"` // degrees/second²
	YawAccel   float64 `json:"yaw_accel"`   // degrees/second²
}

// PhysicsState represents comprehensive 3D physics state
type PhysicsState struct {
	Position            Position            `json:"position"`
	Velocity            Velocity            `json:"velocity"`
	Acceleration        Acceleration        `json:"acceleration"`
	Attitude            Attitude            `json:"attitude"`
	AngularVelocity     AngularVelocity     `json:"angular_velocity"`
	AngularAcceleration AngularAcceleration `json:"angular_acceleration"`
	Mass                float64             `json:"mass"`              // kg
	MomentOfInertia     MomentOfInertia     `json:"moment_of_inertia"` // kg⋅m²
	Forces              Forces              `json:"forces"`            // N
	Torques             Torques             `json:"torques"`           // N⋅m
}

// MomentOfInertia represents rotational inertia in 3 axes
type MomentOfInertia struct {
	Ixx float64 `json:"ixx"` // kg⋅m² around x-axis (roll)
	Iyy float64 `json:"iyy"` // kg⋅m² around y-axis (pitch)
	Izz float64 `json:"izz"` // kg⋅m² around z-axis (yaw)
}

// Forces represents applied forces in 3D
type Forces struct {
	Thrust float64 `json:"thrust"` // N, forward force
	Drag   float64 `json:"drag"`   // N, opposing motion
	Lift   float64 `json:"lift"`   // N, upward force (aircraft)
	Weight float64 `json:"weight"` // N, gravitational force
	Normal float64 `json:"normal"` // N, surface reaction force
}

// Torques represents applied torques in 3D
type Torques struct {
	Roll  float64 `json:"roll"`  // N⋅m around longitudinal axis
	Pitch float64 `json:"pitch"` // N⋅m around lateral axis
	Yaw   float64 `json:"yaw"`   // N⋅m around vertical axis
}

// PlatformState represents the current state of a platform (backwards compatible)
type PlatformState struct {
	ID          string    `json:"id"`
	Position    Position  `json:"position"`
	Velocity    Velocity  `json:"velocity"`
	Heading     float64   `json:"heading"` // degrees, 0-360 (maintained for compatibility)
	Speed       float64   `json:"speed"`   // m/s (maintained for compatibility)
	Roll        float64   `json:"roll"`    // degrees, banking angle for aircraft
	LastUpdated time.Time `json:"lastUpdated"`

	// Enhanced physics state
	Physics PhysicsState `json:"physics"`
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
	Range           float64 `yaml:"range"`

	// Enhanced performance characteristics
	MaxAcceleration   float64 `yaml:"max_acceleration,omitempty"`     // m/s²
	MaxDeceleration   float64 `yaml:"max_deceleration,omitempty"`     // m/s²
	MaxRollRate       float64 `yaml:"max_roll_rate,omitempty"`        // degrees/second
	MaxPitchRate      float64 `yaml:"max_pitch_rate,omitempty"`       // degrees/second
	MaxYawRate        float64 `yaml:"max_yaw_rate,omitempty"`         // degrees/second
	MaxBankAngle      float64 `yaml:"max_bank_angle,omitempty"`       // degrees
	MaxPitchAngle     float64 `yaml:"max_pitch_angle,omitempty"`      // degrees
	MaxThrustToWeight float64 `yaml:"max_thrust_to_weight,omitempty"` // ratio
	MaxLoadFactor     float64 `yaml:"max_load_factor,omitempty"`      // g-force

	// Control surface effectiveness
	ElevatorAuthority float64 `yaml:"elevator_authority,omitempty"` // degrees
	RudderAuthority   float64 `yaml:"rudder_authority,omitempty"`   // degrees
	AileronAuthority  float64 `yaml:"aileron_authority,omitempty"`  // degrees

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

	// Enhanced physical characteristics for realistic physics
	CenterOfGravity CenterOfGravity `yaml:"center_of_gravity,omitempty"`
	EmptyWeight     float64         `yaml:"empty_weight,omitempty"` // kg
	MaxWeight       float64         `yaml:"max_weight,omitempty"`   // kg
	WingArea        float64         `yaml:"wing_area,omitempty"`    // m² (aircraft)
	WetArea         float64         `yaml:"wet_area,omitempty"`     // m² (ships)
	FrontalArea     float64         `yaml:"frontal_area,omitempty"` // m² (land vehicles)
}

// CenterOfGravity represents the center of mass location
type CenterOfGravity struct {
	X float64 `yaml:"x"` // meters from nose/bow/front
	Y float64 `yaml:"y"` // meters from centerline (positive = starboard/right)
	Z float64 `yaml:"z"` // meters from keel/ground (positive = up)
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
	Sensors      SensorCharacteristics      `yaml:"sensors"`
	CallsignConf CallsignConfiguration      `yaml:"callsign_config"`
}

// SensorCharacteristics defines sensor capabilities
type SensorCharacteristics struct {
	HasGPS        bool    `yaml:"has_gps"`
	HasRadar      bool    `yaml:"has_radar"`
	HasCompass    bool    `yaml:"has_compass"`
	RadarRange    float64 `yaml:"radar_range,omitempty"`    // meters
	SonarRange    float64 `yaml:"sonar_range,omitempty"`    // meters
	OpticalRange  float64 `yaml:"optical_range,omitempty"`  // meters
	InfraredRange float64 `yaml:"infrared_range,omitempty"` // meters
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
	ID           string                   `json:"id"`
	PlatformType PlatformType             `json:"platform_type"`
	TypeDef      *PlatformTypeDefinition  `json:"type_def"`
	Config       *PlatformConfiguration   `json:"config"`
	State        PlatformState            `json:"state"`
	CallSign     string                   `json:"call_sign"`

	// Navigation
	Destination *Position `json:"destination,omitempty"`
	Route       []Position `json:"route,omitempty"`

	// Runtime state
	FuelRemaining float64       `json:"fuel_remaining"`
	MissionTime   time.Duration `json:"mission_time"`
	SystemStatus  SystemStatus  `json:"system_status"`

	// Physics state
	lastPosition  Position `json:"-"` // Internal field, don't serialize
	acceleration  float64  `json:"-"` // Internal field, don't serialize
	Mass          float64  `json:"mass"` // Mass for physics calculations
	AngularForces struct { // Angular forces for attitude control
		X float64 `json:"x"`
		Y float64 `json:"y"`
		Z float64 `json:"z"`
	} `json:"angular_forces"`
}

// SystemStatus represents the operational status of various platform systems
type SystemStatus struct {
	PowerSystem         SystemState `json:"power_system"`
	PropulsionSystem    SystemState `json:"propulsion_system"`
	NavigationSystem    SystemState `json:"navigation_system"`
	CommunicationSystem SystemState `json:"communication_system"`
	SensorSystem        SystemState `json:"sensor_system"`
	WeaponSystem        SystemState `json:"weapon_system,omitempty"`
	LifeSupport         SystemState `json:"life_support,omitempty"`
	FuelSystem          SystemState `json:"fuel_system"`
	WeaponStatus        string      `json:"weapon_status"` // "ARMED", "SAFE", "N/A"
}

// SystemState represents the operational state of a system
type SystemState struct {
	Operational bool      `json:"operational"`
	Efficiency  float64   `json:"efficiency"` // 0.0 to 1.0
	LastCheck   time.Time `json:"last_check"`
	Notes       string    `json:"notes,omitempty"`
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

// Update updates the platform's state for the given time step
func (up *UniversalPlatform) Update(deltaTime time.Duration) error {
	deltaSeconds := deltaTime.Seconds()

	// Update mission time
	up.MissionTime += deltaTime

	// Consume fuel if engines are running
	if up.State.Speed > 0 || up.Destination != nil {
		fuelRate := up.TypeDef.Performance.FuelConsumption
		if fuelRate > 0 {
			fuelConsumed := fuelRate * deltaSeconds
			up.FuelRemaining = math.Max(0, up.FuelRemaining-fuelConsumed)
		}
	}

	// If no destination but has velocity, update position from current velocity
	if up.Destination == nil {
		if up.State.Velocity.North != 0 || up.State.Velocity.East != 0 || up.State.Velocity.Up != 0 {
			up.updatePositionFromVelocity(deltaSeconds)
			up.State.LastUpdated = time.Now()
		}
		return nil
	}

	// Update based on platform type
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
		return up.updateGenericMovement(deltaSeconds)
	}
}

func (up *UniversalPlatform) updateAirborneMovement(deltaSeconds float64) error {
	if up.Destination == nil {
		// Even without destination, update position if there's velocity
		if up.State.Velocity.North != 0 || up.State.Velocity.East != 0 || up.State.Velocity.Up != 0 {
			up.updatePositionFromVelocity(deltaSeconds)
		}
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
			up.State.Velocity.Up = 0
		} else if altitudeDiff > 0 {
			up.State.Position.Altitude += maxAltChange
			up.State.Velocity.Up = climbRate
		} else {
			up.State.Position.Altitude -= maxAltChange
			up.State.Velocity.Up = -climbRate
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

	// Maritime movement with ship dynamics
	distance := up.calculateGreatCircleDistance(*up.Destination)
	if distance < 10 { // 10 meter threshold for ships
		up.State.Position = *up.Destination
		up.Destination = nil
		up.State.Speed = 0
		return nil
	}

	// Ships have turning constraints based on their size and speed
	desiredHeading := up.calculateBearing(*up.Destination)

	// Use turning radius for maritime movement calculations
	if up.TypeDef.Performance.TurningRadius == 0 {
		turningRadius := up.TypeDef.Physical.Length * 5 // Default: 5x ship length
		up.TypeDef.Performance.TurningRadius = turningRadius
	}

	up.applyTurningConstraints(desiredHeading, deltaSeconds)

	// Calculate movement with maritime characteristics
	targetSpeed := up.TypeDef.Performance.MaxSpeed * 0.9 // Ships can maintain near max speed in open water

	// Apply acceleration with maritime characteristics
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
	if up.Destination == nil {
		return nil
	}

	// Space movement with orbital mechanics
	distance := up.calculateGreatCircleDistance(*up.Destination)
	if distance < 100 { // 100 meter threshold for satellites
		up.State.Position = *up.Destination
		up.Destination = nil
		up.State.Speed = 0
		return nil
	}

	// Satellites move in predictable orbital paths
	desiredHeading := up.calculateBearing(*up.Destination)

	// Space vehicles have different movement characteristics
	up.applyTurningConstraints(desiredHeading, deltaSeconds)

	// Calculate movement with orbital velocity
	targetSpeed := up.TypeDef.Performance.MaxSpeed * 0.8 // Most satellites maintain consistent orbital speed

	// Maintain orbital altitude
	up.State.Position.Altitude = up.TypeDef.Performance.OrbitalAltitude
	up.State.Speed = targetSpeed
	up.State.Heading = 90 // Generally eastward

	up.State.LastUpdated = time.Now()
	return nil
}

func (up *UniversalPlatform) updateGenericMovement(deltaSeconds float64) error {
	if up.Destination == nil {
		return nil
	}

	// Generic movement for unknown platform types
	distance := up.calculateGreatCircleDistance(*up.Destination)
	if distance < 50 { // 50 meter threshold
		up.State.Position = *up.Destination
		up.Destination = nil
		up.State.Speed = 0
		return nil
	}

	// Calculate desired heading
	desiredHeading := up.calculateBearing(*up.Destination)
	up.applyTurningConstraints(desiredHeading, deltaSeconds)

	// Apply generic acceleration
	targetSpeed := up.TypeDef.Performance.CruiseSpeed
	up.applyAccelerationConstraints(targetSpeed, deltaSeconds)

	// Update position based on current heading and speed
	up.updatePositionFromHeadingAndSpeed(deltaSeconds)

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
	// Navigation heading: 0° = North, 90° = East, 180° = South, 270° = West
	headingRad := up.State.Heading * math.Pi / 180.0

	// Calculate distance moved
	distance := up.State.Speed * deltaSeconds

	// Earth radius in meters
	earthRadius := 6371000.0

	// Calculate new position using proper navigation calculations
	// North component: distance * cos(heading)
	// East component: distance * sin(heading)
	deltaLat := distance * math.Cos(headingRad) / earthRadius * 180.0 / math.Pi
	deltaLon := distance * math.Sin(headingRad) / earthRadius * 180.0 / math.Pi / math.Cos(up.State.Position.Latitude*math.Pi/180.0)

	up.State.Position.Latitude += deltaLat
	up.State.Position.Longitude += deltaLon

	// Calculate velocity components for state tracking
	up.State.Velocity.North = up.State.Speed * math.Cos(headingRad)
	up.State.Velocity.East = up.State.Speed * math.Sin(headingRad)
	up.State.Velocity.Up = 0 // Will be calculated separately for altitude changes
}

// updatePositionFromVelocity updates position based on current velocity components
func (up *UniversalPlatform) updatePositionFromVelocity(deltaSeconds float64) {
	// Earth radius in meters
	earthRadius := 6371000.0

	// Convert North/East velocity to lat/lon changes
	deltaLat := up.State.Velocity.North * deltaSeconds / earthRadius * 180.0 / math.Pi
	deltaLon := up.State.Velocity.East * deltaSeconds / (earthRadius * math.Cos(up.State.Position.Latitude*math.Pi/180.0)) * 180.0 / math.Pi
	deltaAlt := up.State.Velocity.Up * deltaSeconds

	// Update position
	up.State.Position.Latitude += deltaLat
	up.State.Position.Longitude += deltaLon
	up.State.Position.Altitude += deltaAlt

	// Sync physics state
	up.State.Physics.Position = up.State.Position

	// Update heading and speed from velocity
	if up.State.Velocity.East != 0 || up.State.Velocity.North != 0 {
		up.State.Heading = math.Atan2(up.State.Velocity.East, up.State.Velocity.North) * 180.0 / math.Pi
		if up.State.Heading < 0 {
			up.State.Heading += 360
		}
	}

	up.State.Speed = math.Sqrt(up.State.Velocity.North*up.State.Velocity.North + up.State.Velocity.East*up.State.Velocity.East)
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

	// Enhanced 3D physics methods
	Initialize3DPhysics()
	Update3DPhysics(deltaTime time.Duration) error
	GetPhysicsState() PhysicsState
	SetPhysicsState(physics PhysicsState)
}

// Enhanced 3D physics methods for UniversalPlatform

// Initialize3DPhysics sets up realistic physics parameters based on platform type
func (up *UniversalPlatform) Initialize3DPhysics() {
	// Initialize physics state with platform characteristics
	up.State.Physics = PhysicsState{
		Position:            up.State.Position,
		Velocity:            up.State.Velocity,
		Acceleration:        Acceleration{},
		Attitude:            Attitude{Yaw: up.State.Heading},
		AngularVelocity:     AngularVelocity{},
		AngularAcceleration: AngularAcceleration{},
		Mass:                up.TypeDef.Physical.Mass,
		MomentOfInertia:     up.calculateMomentOfInertia(),
		Forces:              Forces{Weight: up.TypeDef.Physical.Mass * 9.81},
		Torques:             Torques{},
	}
}

// calculateMomentOfInertia estimates moment of inertia based on platform geometry
func (up *UniversalPlatform) calculateMomentOfInertia() MomentOfInertia {
	mass := up.TypeDef.Physical.Mass
	length := up.TypeDef.Physical.Length
	width := up.TypeDef.Physical.Width
	height := up.TypeDef.Physical.Height

	// Simplified calculations for different platform types
	switch up.PlatformType {
	case PlatformTypeAirborne:
		// Aircraft approximated as elongated ellipsoid
		return MomentOfInertia{
			Ixx: mass * (width*width + height*height) / 20,   // Roll axis
			Iyy: mass * (length*length + height*height) / 20, // Pitch axis
			Izz: mass * (length*length + width*width) / 20,   // Yaw axis
		}
	case PlatformTypeMaritime:
		// Ship approximated as rectangular prism
		return MomentOfInertia{
			Ixx: mass * (width*width + height*height) / 12,
			Iyy: mass * (length*length + height*height) / 12,
			Izz: mass * (length*length + width*width) / 12,
		}
	case PlatformTypeLand:
		// Vehicle approximated as box
		return MomentOfInertia{
			Ixx: mass * (width*width + height*height) / 12,
			Iyy: mass * (length*length + height*height) / 12,
			Izz: mass * (length*length + width*width) / 12,
		}
	case PlatformTypeSpace:
		// Satellite approximated as uniform distribution
		avgDim := (length + width + height) / 3
		return MomentOfInertia{
			Ixx: mass * avgDim * avgDim / 6,
			Iyy: mass * avgDim * avgDim / 6,
			Izz: mass * avgDim * avgDim / 6,
		}
	default:
		// Default calculation
		return MomentOfInertia{
			Ixx: mass * (width*width + height*height) / 12,
			Iyy: mass * (length*length + height*height) / 12,
			Izz: mass * (length*length + width*width) / 12,
		}
	}
}

// Update3DPhysics applies 3D physics simulation to the platform
func (up *UniversalPlatform) Update3DPhysics(deltaTime time.Duration) error {
	deltaSeconds := deltaTime.Seconds()

	// Apply platform-specific forces
	switch up.PlatformType {
	case PlatformTypeAirborne:
		up.applyAerodynamicForces()
		up.applyThrustForces()
		up.applyGravity()
	case PlatformTypeLand:
		up.applyGroundForces()
		up.applyEngineForces()
		up.applyGravity()
	case PlatformTypeMaritime:
		up.applyHydrodynamicForces()
		up.applyPropulsionForces()
		up.applyBuoyancyAndGravity()
	case PlatformTypeSpace:
		up.applyOrbitalForces()
		up.applyThrusterForces()
	default:
		up.applyBasicForces()
	}

	// Integrate forces to acceleration and velocity
	up.integrateForces(deltaSeconds)

	// Integrate velocity to position
	up.integrateVelocity(deltaSeconds)

	// Update attitude from angular forces
	up.updateAttitudeFromForces(deltaSeconds)

	// Update timestamp
	up.State.LastUpdated = time.Now()

	return nil
}

// updateAircraft3DPhysics applies aircraft-specific physics
func (up *UniversalPlatform) updateAircraft3DPhysics(deltaTime float64) {
	// For aircraft, update position based on speed and heading
	up.updatePositionFromHeadingAndSpeed(deltaTime)
}

// updateLand3DPhysics applies land vehicle physics
func (up *UniversalPlatform) updateLand3DPhysics(deltaTime float64) {
	// For land vehicles, update position based on speed and heading
	up.updatePositionFromHeadingAndSpeed(deltaTime)
}

// updateMaritime3DPhysics applies maritime vehicle physics
func (up *UniversalPlatform) updateMaritime3DPhysics(deltaTime float64) {
	// For maritime vehicles, update position based on speed and heading
	up.updatePositionFromHeadingAndSpeed(deltaTime)
}

// updateSpace3DPhysics applies space vehicle physics
func (up *UniversalPlatform) updateSpace3DPhysics(deltaTime float64) {
	// For space vehicles, update position based on speed and heading
	up.updatePositionFromHeadingAndSpeed(deltaTime)
}

// Force represents a 3D force vector
type Force struct {
	X float64 `json:"x"` // Forward/backward force (N)
	Y float64 `json:"y"` // Left/right force (N)
	Z float64 `json:"z"` // Up/down force (N)
}

// Force application methods
func (up *UniversalPlatform) applyAerodynamicForces() {
	// Basic aerodynamic forces - can be expanded later
	// For now, apply basic drag opposing motion
	speed := math.Sqrt(up.State.Velocity.North*up.State.Velocity.North + up.State.Velocity.East*up.State.Velocity.East + up.State.Velocity.Up*up.State.Velocity.Up)
	if speed > 0 {
		dragCoeff := 0.1
		up.State.Physics.Forces.Drag -= dragCoeff * speed * speed // Apply drag force
	}
}

func (up *UniversalPlatform) applyThrustForces() {
	// Convert speed and heading to thrust forces
	if up.State.Speed > 0 {
		headingRad := up.State.Heading * math.Pi / 180.0
		thrustMagnitude := up.State.Speed * 0.1 // Simple thrust model
		up.State.Physics.Forces.Thrust += thrustMagnitude * math.Sin(headingRad)
		up.State.Physics.Forces.Normal += thrustMagnitude * math.Cos(headingRad)
	}
}

func (up *UniversalPlatform) applyGravity() {
	up.State.Physics.Forces.Weight = 9.81 * up.State.Physics.Mass // Standard gravity
}

func (up *UniversalPlatform) applyHydrodynamicForces() {
	// Similar to aerodynamic but for water
	up.applyAerodynamicForces() // Simplified for now
}

func (up *UniversalPlatform) applyPropulsionForces() {
	up.applyThrustForces() // Simplified for now
}

func (up *UniversalPlatform) applyBuoyancyAndGravity() {
	up.applyGravity()
	// Add buoyancy force to counteract gravity for surface vessels
	up.State.Physics.Forces.Lift += 9.81 * up.State.Physics.Mass // Simplified buoyancy
}

func (up *UniversalPlatform) applyGroundForces() {
	// Ground friction and resistance
	up.applyAerodynamicForces() // Simplified friction model
}

func (up *UniversalPlatform) applyEngineForces() {
	up.applyThrustForces() // Simplified for now
}

func (up *UniversalPlatform) applyOrbitalForces() {
	// Basic orbital mechanics - simplified
	up.applyGravity()
}

func (up *UniversalPlatform) applyThrusterForces() {
	up.applyThrustForces() // Simplified for now
}

func (up *UniversalPlatform) applyBasicForces() {
	up.applyThrustForces()
	up.applyGravity()
}

// Force and velocity integration methods
func (up *UniversalPlatform) integrateForces(deltaTime float64) {
	if up.State.Physics.Mass > 0 {
		// F = ma, so a = F/m
		up.State.Physics.Acceleration.North = up.State.Physics.Forces.Thrust / up.State.Physics.Mass
		up.State.Physics.Acceleration.East = up.State.Physics.Forces.Normal / up.State.Physics.Mass
		up.State.Physics.Acceleration.Up = up.State.Physics.Forces.Lift / up.State.Physics.Mass

		// Integrate acceleration to velocity
		up.State.Physics.Velocity.North += up.State.Physics.Acceleration.North * deltaTime
		up.State.Physics.Velocity.East += up.State.Physics.Acceleration.East * deltaTime
		up.State.Physics.Velocity.Up += up.State.Physics.Acceleration.Up * deltaTime

		// Reset forces for next iteration
		up.State.Physics.Forces.Thrust = 0
		up.State.Physics.Forces.Normal = 0
		up.State.Physics.Forces.Lift = 0
	}
}

func (up *UniversalPlatform) integrateVelocity(deltaTime float64) {
	// Convert velocity to lat/lon/alt changes
	earthRadius := 6371000.0 // Earth radius in meters

	// Convert North/East velocity to lat/lon changes
	deltaLat := up.State.Velocity.North * deltaTime / earthRadius * 180.0 / math.Pi
	deltaLon := up.State.Velocity.East * deltaTime / (earthRadius * math.Cos(up.State.Position.Latitude*math.Pi/180.0)) * 180.0 / math.Pi
	deltaAlt := up.State.Velocity.Up * deltaTime

	// Update position
	up.State.Position.Latitude += deltaLat
	up.State.Position.Longitude += deltaLon
	up.State.Position.Altitude += deltaAlt

	// Update speed from velocity magnitude
	up.State.Speed = math.Sqrt(up.State.Velocity.North*up.State.Velocity.North + up.State.Velocity.East*up.State.Velocity.East + up.State.Velocity.Up*up.State.Velocity.Up)

	// Update heading from velocity direction
	if up.State.Velocity.East != 0 || up.State.Velocity.North != 0 {
		up.State.Heading = math.Atan2(up.State.Velocity.East, up.State.Velocity.North) * 180.0 / math.Pi
		if up.State.Heading < 0 {
			up.State.Heading += 360
		}
	}
}

func (up *UniversalPlatform) updateAttitudeFromForces(deltaTime float64) {
	// Apply angular velocity to attitude changes
	up.State.Physics.Attitude.Roll += up.State.Physics.AngularVelocity.RollRate * deltaTime
	up.State.Physics.Attitude.Pitch += up.State.Physics.AngularVelocity.PitchRate * deltaTime
	up.State.Physics.Attitude.Yaw += up.State.Physics.AngularVelocity.YawRate * deltaTime

	// Update legacy roll field for compatibility
	up.State.Roll = up.State.Physics.Attitude.Roll

	// Basic attitude control from angular forces
	if up.AngularForces.X != 0 || up.AngularForces.Y != 0 || up.AngularForces.Z != 0 {
		// Integrate angular forces to angular velocity
		if up.State.Physics.Mass > 0 {
			up.State.Physics.AngularVelocity.RollRate += up.AngularForces.X * deltaTime / up.State.Physics.Mass
			up.State.Physics.AngularVelocity.PitchRate += up.AngularForces.Y * deltaTime / up.State.Physics.Mass
			up.State.Physics.AngularVelocity.YawRate += up.AngularForces.Z * deltaTime / up.State.Physics.Mass
		}

		// Reset angular forces
		up.AngularForces.X = 0
		up.AngularForces.Y = 0
		up.AngularForces.Z = 0
	}
}

// GetPhysicsState returns the current physics state
func (up *UniversalPlatform) GetPhysicsState() PhysicsState {
	return up.State.Physics
}

// SetPhysicsState sets the physics state
func (up *UniversalPlatform) SetPhysicsState(physics PhysicsState) {
	up.State.Physics = physics
	// Sync legacy state
	up.State.Position = physics.Position
	up.State.Velocity = physics.Velocity
	up.State.Heading = physics.Attitude.Yaw
	up.State.Speed = math.Sqrt(physics.Velocity.North*physics.Velocity.North + physics.Velocity.East*physics.Velocity.East)
}

// CalculateDistanceTo calculates distance to another position
func (up *UniversalPlatform) CalculateDistanceTo(target Position) float64 {
	return up.calculateGreatCircleDistance(target)
}

// GetStatus returns a status summary of the platform
func (up *UniversalPlatform) GetStatus() PlatformStatus {
	return PlatformStatus{
		ID:           up.ID,
		PlatformType: up.PlatformType,
		Position:     up.State.Position,
		Velocity:     up.State.Velocity,
		Heading:      up.State.Heading,
		Speed:        up.State.Speed,
		LastUpdated:  up.State.LastUpdated,
	}
}

// ApplyForce applies a 3D force to the platform
func (up *UniversalPlatform) ApplyForce(force Force) {
	// Add the applied force to existing forces
	up.State.Physics.Forces.Thrust += force.X
	up.State.Physics.Forces.Normal += force.Y
	up.State.Physics.Forces.Lift += force.Z
}

// calculateTurningRate calculates the turning rate for the platform
func (p *UniversalPlatform) calculateTurningRate(deltaTime time.Duration) float64 {
	// Get turning radius from performance characteristics
	turningRadius := p.TypeDef.Performance.TurningRadius
	if turningRadius == 0 {
		turningRadius = 500.0 // Default turning radius in meters
	}

	// Calculate turning rate based on current speed and turning radius
	// Angular velocity = linear velocity / radius
	if p.State.Speed > 0 && turningRadius > 0 {
		return p.State.Speed / turningRadius // radians per second
	}
	return 0.0
}

// PlatformStatus represents a summary of platform state
type PlatformStatus struct {
	ID           string       `json:"id"`
	PlatformType PlatformType `json:"platform_type"`
	Position     Position     `json:"position"`
	Velocity     Velocity     `json:"velocity"`
	Heading      float64      `json:"heading"`
	Speed        float64      `json:"speed"`
	LastUpdated  time.Time    `json:"last_updated"`
}
