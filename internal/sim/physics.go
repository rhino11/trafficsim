package sim

import (
	"math"
	"time"

	"github.com/rhino11/trafficsim/internal/models"
)

// Platform type constants
const (
	PlatformTypeAircraft = "aircraft"
	PlatformTypeLand     = "land"
	PlatformTypeMaritime = "maritime"
	PlatformTypeSpace    = "space"
)

// PhysicsEngine handles realistic movement calculations for all platform types
type PhysicsEngine struct {
	// Environmental constants
	EarthRadius   float64 // meters
	GravityAccel  float64 // m/s²
	AirDensity    float64 // kg/m³ (sea level)
	SeaLevelPress float64 // Pa

	// Simulation parameters
	TimeStep      time.Duration
	EnableWeather bool
	EnableTerrain bool
}

// NewPhysicsEngine creates a new physics engine with realistic constants
func NewPhysicsEngine() *PhysicsEngine {
	return &PhysicsEngine{
		EarthRadius:   6371000.0, // meters
		GravityAccel:  9.81,      // m/s²
		AirDensity:    1.225,     // kg/m³
		SeaLevelPress: 101325.0,  // Pa
		TimeStep:      time.Second,
		EnableWeather: false, // Start simple
		EnableTerrain: false, // Start simple
	}
}

// CalculateMovement performs physics-based movement calculation for a platform
func (pe *PhysicsEngine) CalculateMovement(platform models.Platform, deltaTime time.Duration) error {
	// Try to cast to UniversalPlatform for enhanced physics
	if universalPlatform, ok := platform.(*models.UniversalPlatform); ok {
		return pe.updateUniversalPlatform(universalPlatform, deltaTime)
	}

	// Fallback to platform's own Update method for all platforms
	return platform.Update(deltaTime)
}

// updateUniversalPlatform handles movement for the new universal platform system
func (pe *PhysicsEngine) updateUniversalPlatform(platform *models.UniversalPlatform, deltaTime time.Duration) error {
	deltaSeconds := deltaTime.Seconds()

	// Skip movement if no destination
	if platform.Destination == nil {
		return nil
	}

	// Calculate distance and bearing to destination
	distance := pe.CalculateGreatCircleDistance(
		platform.State.Position,
		*platform.Destination,
	)

	// Check if we've reached the destination
	arrivalThreshold := pe.getArrivalThreshold(platform.PlatformType)
	if distance < arrivalThreshold {
		platform.State.Position = *platform.Destination
		platform.Destination = nil
		platform.State.Speed = 0
		return nil
	}

	// Calculate desired movement vector
	bearing := pe.CalculateBearing(platform.State.Position, *platform.Destination)

	// Apply platform-specific physics
	switch platform.PlatformType {
	case models.PlatformTypeAirborne:
		return pe.updateAircraftPhysics(platform, bearing, distance, deltaSeconds)
	case models.PlatformTypeMaritime:
		return pe.updateMaritimePhysics(platform, bearing, distance, deltaSeconds)
	case models.PlatformTypeLand:
		return pe.updateLandPhysics(platform, bearing, distance, deltaSeconds)
	case models.PlatformTypeSpace:
		return pe.updateSpacePhysics(platform, deltaSeconds)
	default:
		return pe.updateGenericPhysics(platform, bearing, distance, deltaSeconds)
	}
}

// updateAircraftPhysics implements realistic aircraft movement
func (pe *PhysicsEngine) updateAircraftPhysics(platform *models.UniversalPlatform, bearing, _ /* distance */, deltaSeconds float64) error {
	// Get performance characteristics
	maxSpeed := platform.TypeDef.Performance.MaxSpeed
	cruiseSpeed := platform.TypeDef.Performance.CruiseSpeed
	climbRate := platform.TypeDef.Performance.ClimbRate
	if climbRate == 0 {
		climbRate = 10.0 // Default climb rate
	}

	// Calculate turning constraints
	turningRadius := platform.TypeDef.Performance.TurningRadius
	if turningRadius == 0 {
		// Calculate based on speed and standard bank angle (30°)
		bankAngle := 30.0 * math.Pi / 180.0
		turningRadius = (cruiseSpeed * cruiseSpeed) / (pe.GravityAccel * math.Tan(bankAngle))
	}

	// Apply heading change with turning constraints
	newHeading := pe.applyTurningConstraints(
		platform.State.Heading,
		bearing,
		platform.State.Speed,
		turningRadius,
		deltaSeconds,
	)
	platform.State.Heading = newHeading

	// Handle altitude changes
	altitudeDiff := platform.Destination.Altitude - platform.State.Position.Altitude
	if math.Abs(altitudeDiff) > 10 {
		maxAltChange := climbRate * deltaSeconds
		if math.Abs(altitudeDiff) <= maxAltChange {
			platform.State.Position.Altitude = platform.Destination.Altitude
		} else if altitudeDiff > 0 {
			platform.State.Position.Altitude += maxAltChange
		} else {
			platform.State.Position.Altitude -= maxAltChange
		}
	}

	// Apply speed control with acceleration limits
	targetSpeed := math.Min(cruiseSpeed, maxSpeed)
	platform.State.Speed = pe.applyAcceleration(
		platform.State.Speed,
		targetSpeed,
		platform.TypeDef.Performance.Acceleration,
		deltaSeconds,
	)

	// Update position
	pe.updatePosition(&platform.State, deltaSeconds)
	platform.State.LastUpdated = time.Now()

	return nil
}

// updateMaritimePhysics implements realistic ship movement
func (pe *PhysicsEngine) updateMaritimePhysics(platform *models.UniversalPlatform, bearing, _ /* distance */, deltaSeconds float64) error {
	// Ships have different characteristics
	cruiseSpeed := platform.TypeDef.Performance.CruiseSpeed

	// Ships have large turning radii
	turningRadius := platform.TypeDef.Performance.TurningRadius
	if turningRadius == 0 {
		turningRadius = platform.TypeDef.Physical.Length * 6 // 6x ship length
	}

	// Apply heading change (ships turn slowly)
	newHeading := pe.applyTurningConstraints(
		platform.State.Heading,
		bearing,
		platform.State.Speed,
		turningRadius,
		deltaSeconds,
	)
	platform.State.Heading = newHeading

	// Ships stay at sea level
	platform.State.Position.Altitude = 0

	// Apply acceleration (ships accelerate slowly)
	acceleration := platform.TypeDef.Performance.Acceleration
	if acceleration == 0 {
		acceleration = 0.3 // Default slow acceleration for ships
	}

	platform.State.Speed = pe.applyAcceleration(
		platform.State.Speed,
		cruiseSpeed,
		acceleration,
		deltaSeconds,
	)

	// Update position
	pe.updatePosition(&platform.State, deltaSeconds)
	platform.State.LastUpdated = time.Now()

	return nil
}

// updateLandPhysics implements realistic land vehicle movement
func (pe *PhysicsEngine) updateLandPhysics(platform *models.UniversalPlatform, bearing, distance, deltaSeconds float64) error {
	// Land vehicles have terrain constraints
	cruiseSpeed := platform.TypeDef.Performance.CruiseSpeed

	// Apply turning constraints
	newHeading := pe.applyTurningConstraints(
		platform.State.Heading,
		bearing,
		platform.State.Speed,
		platform.TypeDef.Performance.TurningRadius,
		deltaSeconds,
	)
	platform.State.Heading = newHeading

	// Apply gradient constraints
	if pe.EnableTerrain {
		pe.applyGradientConstraints(platform, distance)
	}

	// Apply acceleration
	platform.State.Speed = pe.applyAcceleration(
		platform.State.Speed,
		cruiseSpeed,
		platform.TypeDef.Performance.Acceleration,
		deltaSeconds,
	)

	// Update position
	pe.updatePosition(&platform.State, deltaSeconds)
	platform.State.LastUpdated = time.Now()

	return nil
}

// updateSpacePhysics implements orbital mechanics for satellites
func (pe *PhysicsEngine) updateSpacePhysics(platform *models.UniversalPlatform, deltaSeconds float64) error {
	// Space platforms have orbital characteristics
	orbitalVelocity := platform.TypeDef.Performance.OrbitalVelocity
	if orbitalVelocity == 0 {
		// Calculate orbital velocity based on altitude
		altitude := platform.State.Position.Altitude
		orbitalVelocity = math.Sqrt(pe.GravityAccel * pe.EarthRadius * pe.EarthRadius / (pe.EarthRadius + altitude))
	}

	// Maintain orbital velocity
	platform.State.Speed = orbitalVelocity

	// Update position in orbit
	pe.updateOrbitalPosition(&platform.State, deltaSeconds)
	platform.State.LastUpdated = time.Now()

	return nil
}

// updateGenericPhysics provides basic movement for unknown platform types
func (pe *PhysicsEngine) updateGenericPhysics(platform *models.UniversalPlatform, bearing, _ /* distance */, deltaSeconds float64) error {
	// Basic movement
	platform.State.Heading = bearing
	platform.State.Speed = platform.TypeDef.Performance.CruiseSpeed

	// Update position
	pe.updatePosition(&platform.State, deltaSeconds)
	platform.State.LastUpdated = time.Now()

	return nil
}

// Helper methods

func (pe *PhysicsEngine) getArrivalThreshold(platformType models.PlatformType) float64 {
	switch platformType {
	case models.PlatformTypeAirborne:
		return 100.0 // 100 meters
	case models.PlatformTypeMaritime:
		return 50.0 // 50 meters
	case models.PlatformTypeLand:
		return 10.0 // 10 meters
	case models.PlatformTypeSpace:
		return 1000.0 // 1 km for satellites
	default:
		return 50.0 // Default threshold
	}
}

func (pe *PhysicsEngine) applyTurningConstraints(currentHeading, targetHeading, speed, turningRadius, deltaTime float64) float64 {
	if turningRadius <= 0 || speed <= 0 {
		return targetHeading
	}

	// Calculate heading difference
	headingDiff := targetHeading - currentHeading
	for headingDiff > 180 {
		headingDiff -= 360
	}
	for headingDiff < -180 {
		headingDiff += 360
	}

	// Calculate maximum turn rate
	maxTurnRate := speed / turningRadius * 180.0 / math.Pi // degrees per second
	maxTurnChange := maxTurnRate * deltaTime

	// Apply turn constraint
	if math.Abs(headingDiff) <= maxTurnChange {
		return targetHeading
	} else if headingDiff > 0 {
		return math.Mod(currentHeading+maxTurnChange, 360)
	} else {
		return math.Mod(currentHeading-maxTurnChange+360, 360)
	}
}

func (pe *PhysicsEngine) applyAcceleration(currentSpeed, targetSpeed, acceleration, deltaTime float64) float64 {
	if acceleration <= 0 {
		return targetSpeed
	}

	speedDiff := targetSpeed - currentSpeed
	maxSpeedChange := acceleration * deltaTime

	if math.Abs(speedDiff) <= maxSpeedChange {
		return targetSpeed
	} else if speedDiff > 0 {
		return currentSpeed + maxSpeedChange
	} else {
		return currentSpeed - maxSpeedChange
	}
}

func (pe *PhysicsEngine) applyGradientConstraints(platform *models.UniversalPlatform, _ /* distance */ float64) {
	// Apply terrain gradient constraints for land vehicles
	maxGradient := platform.TypeDef.Performance.MaxGradient
	if maxGradient > 0 {
		// Implement gradient checking logic here
		// For now, this is a placeholder
	}
}

func (pe *PhysicsEngine) updatePosition(state *models.PlatformState, deltaTime float64) {
	// Convert heading to radians
	headingRad := state.Heading * math.Pi / 180.0

	// Calculate distance moved
	distance := state.Speed * deltaTime

	// Calculate new position
	deltaLat := distance * math.Cos(headingRad) / pe.EarthRadius * 180.0 / math.Pi
	deltaLon := distance * math.Sin(headingRad) / pe.EarthRadius * 180.0 / math.Pi / math.Cos(state.Position.Latitude*math.Pi/180.0)

	state.Position.Latitude += deltaLat
	state.Position.Longitude += deltaLon

	// Update velocity components
	state.Velocity.North = state.Speed * math.Cos(headingRad)
	state.Velocity.East = state.Speed * math.Sin(headingRad)
}

func (pe *PhysicsEngine) updateOrbitalPosition(state *models.PlatformState, deltaTime float64) {
	// Simplified orbital mechanics
	orbitalPeriod := 90.0 * 60.0 // 90 minutes in seconds
	angularVelocity := 2.0 * math.Pi / orbitalPeriod

	// Update orbital position
	currentAngle := angularVelocity * deltaTime

	// Move eastward along the equator (simplified)
	deltaLon := currentAngle * 180.0 / math.Pi
	state.Position.Longitude += deltaLon

	// Normalize longitude
	for state.Position.Longitude > 180 {
		state.Position.Longitude -= 360
	}
	for state.Position.Longitude < -180 {
		state.Position.Longitude += 360
	}
}

// CalculateGreatCircleDistance calculates the distance between two positions
func (pe *PhysicsEngine) CalculateGreatCircleDistance(pos1, pos2 models.Position) float64 {
	// Haversine formula
	lat1 := pos1.Latitude * math.Pi / 180.0
	lat2 := pos2.Latitude * math.Pi / 180.0
	deltaLat := (pos2.Latitude - pos1.Latitude) * math.Pi / 180.0
	deltaLon := (pos2.Longitude - pos1.Longitude) * math.Pi / 180.0

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := pe.EarthRadius * c

	// Add altitude difference
	altDiff := pos2.Altitude - pos1.Altitude
	return math.Sqrt(distance*distance + altDiff*altDiff)
}

// CalculateBearing calculates the bearing from one position to another
func (pe *PhysicsEngine) CalculateBearing(from, to models.Position) float64 {
	lat1 := from.Latitude * math.Pi / 180.0
	lat2 := to.Latitude * math.Pi / 180.0
	deltaLon := (to.Longitude - from.Longitude) * math.Pi / 180.0

	y := math.Sin(deltaLon) * math.Cos(lat2)
	x := math.Cos(lat1)*math.Sin(lat2) - math.Sin(lat1)*math.Cos(lat2)*math.Cos(deltaLon)

	bearing := math.Atan2(y, x) * 180.0 / math.Pi
	return math.Mod(bearing+360, 360)
}
