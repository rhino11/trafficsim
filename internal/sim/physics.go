package sim

import (
	"math"
	"time"

	"github.com/rhino11/trafficsim/internal/models"
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
	switch p := platform.(type) {
	case *models.UniversalPlatform:
		return pe.updateUniversalPlatform(p, deltaTime)
	default:
		// Fallback to platform's own Update method
		return platform.Update(deltaTime)
	}
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
func (pe *PhysicsEngine) updateAircraftPhysics(platform *models.UniversalPlatform, bearing, distance, deltaSeconds float64) error {
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
func (pe *PhysicsEngine) updateMaritimePhysics(platform *models.UniversalPlatform, bearing, distance, deltaSeconds float64) error {
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

// updateLandPhysics implements realistic ground vehicle movement
func (pe *PhysicsEngine) updateLandPhysics(platform *models.UniversalPlatform, bearing, distance, deltaSeconds float64) error {
	cruiseSpeed := platform.TypeDef.Performance.CruiseSpeed

	// Land vehicles have small turning radii
	turningRadius := platform.TypeDef.Performance.TurningRadius
	if turningRadius == 0 {
		turningRadius = platform.TypeDef.Physical.Length * 1.5 // 1.5x vehicle length
	}

	// Apply heading change
	newHeading := pe.applyTurningConstraints(
		platform.State.Heading,
		bearing,
		platform.State.Speed,
		turningRadius,
		deltaSeconds,
	)
	platform.State.Heading = newHeading

	// Handle terrain constraints
	targetSpeed := cruiseSpeed
	if pe.EnableTerrain {
		targetSpeed = pe.applyTerrainConstraints(platform, targetSpeed)
	}

	// Apply acceleration
	acceleration := platform.TypeDef.Performance.Acceleration
	if acceleration == 0 {
		acceleration = 2.5 // Default acceleration for land vehicles
	}

	platform.State.Speed = pe.applyAcceleration(
		platform.State.Speed,
		targetSpeed,
		acceleration,
		deltaSeconds,
	)

	// Update position with terrain following
	pe.updatePosition(&platform.State, deltaSeconds)
	platform.State.LastUpdated = time.Now()

	return nil
}

// updateSpacePhysics implements orbital mechanics
func (pe *PhysicsEngine) updateSpacePhysics(platform *models.UniversalPlatform, deltaSeconds float64) error {
	orbitalPeriod := platform.TypeDef.Performance.OrbitalPeriod
	if orbitalPeriod <= 0 {
		return nil // No orbital motion defined
	}

	// Calculate angular velocity (radians per second)
	angularVelocity := 2 * math.Pi / orbitalPeriod
	deltaAngle := angularVelocity * deltaSeconds

	// Update longitude (eastward motion)
	currentLon := platform.State.Position.Longitude * math.Pi / 180.0
	newLon := currentLon + deltaAngle

	// Normalize longitude
	for newLon > math.Pi {
		newLon -= 2 * math.Pi
	}
	for newLon < -math.Pi {
		newLon += 2 * math.Pi
	}

	platform.State.Position.Longitude = newLon * 180.0 / math.Pi

	// Calculate latitude oscillation based on orbital inclination
	inclination := platform.TypeDef.Performance.Inclination * math.Pi / 180.0
	platform.MissionTime += time.Duration(deltaSeconds * float64(time.Second))
	timeInOrbit := platform.MissionTime.Seconds()
	latitudePhase := 2 * math.Pi * timeInOrbit / orbitalPeriod

	maxLatitude := inclination
	platform.State.Position.Latitude = maxLatitude * math.Sin(latitudePhase) * 180.0 / math.Pi

	// Maintain orbital altitude
	if platform.TypeDef.Performance.OrbitalAltitude > 0 {
		platform.State.Position.Altitude = platform.TypeDef.Performance.OrbitalAltitude
	}

	// Set orbital velocity
	platform.State.Speed = platform.TypeDef.Performance.OrbitalVelocity
	platform.State.Heading = 90 // Generally eastward

	platform.State.LastUpdated = time.Now()
	return nil
}

// updateGenericPhysics provides fallback movement calculations
func (pe *PhysicsEngine) updateGenericPhysics(platform *models.UniversalPlatform, bearing, distance, deltaSeconds float64) error {
	// Simple point-to-point movement
	platform.State.Heading = bearing
	platform.State.Speed = platform.TypeDef.Performance.CruiseSpeed

	pe.updatePosition(&platform.State, deltaSeconds)
	platform.State.LastUpdated = time.Now()

	return nil
}

// Physics calculation helper functions

// CalculateGreatCircleDistance calculates the distance between two positions
func (pe *PhysicsEngine) CalculateGreatCircleDistance(pos1, pos2 models.Position) float64 {
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

// CalculateBearing calculates the bearing from pos1 to pos2
func (pe *PhysicsEngine) CalculateBearing(pos1, pos2 models.Position) float64 {
	lat1 := pos1.Latitude * math.Pi / 180.0
	lat2 := pos2.Latitude * math.Pi / 180.0
	deltaLon := (pos2.Longitude - pos1.Longitude) * math.Pi / 180.0

	y := math.Sin(deltaLon) * math.Cos(lat2)
	x := math.Cos(lat1)*math.Sin(lat2) - math.Sin(lat1)*math.Cos(lat2)*math.Cos(deltaLon)

	bearing := math.Atan2(y, x) * 180.0 / math.Pi
	return math.Mod(bearing+360, 360) // Normalize to 0-360
}

// applyTurningConstraints applies realistic turning limitations
func (pe *PhysicsEngine) applyTurningConstraints(currentHeading, desiredHeading, speed, turningRadius, deltaSeconds float64) float64 {
	if turningRadius <= 0 || speed <= 0 {
		return desiredHeading
	}

	// Calculate heading difference
	headingDiff := desiredHeading - currentHeading

	// Normalize to -180 to 180
	for headingDiff > 180 {
		headingDiff -= 360
	}
	for headingDiff < -180 {
		headingDiff += 360
	}

	// Calculate maximum turn rate based on turning radius and speed
	maxTurnRate := (speed / turningRadius) * 180.0 / math.Pi // degrees per second
	maxTurnChange := maxTurnRate * deltaSeconds

	// Apply turn rate limitation
	if math.Abs(headingDiff) <= maxTurnChange {
		return desiredHeading
	} else if headingDiff > 0 {
		return math.Mod(currentHeading+maxTurnChange, 360)
	} else {
		return math.Mod(currentHeading-maxTurnChange+360, 360)
	}
}

// applyAcceleration applies realistic acceleration constraints
func (pe *PhysicsEngine) applyAcceleration(currentSpeed, targetSpeed, acceleration, deltaSeconds float64) float64 {
	if acceleration <= 0 {
		return targetSpeed
	}

	speedDiff := targetSpeed - currentSpeed
	maxSpeedChange := acceleration * deltaSeconds

	if math.Abs(speedDiff) <= maxSpeedChange {
		return targetSpeed
	} else if speedDiff > 0 {
		return currentSpeed + maxSpeedChange
	} else {
		return math.Max(0, currentSpeed-maxSpeedChange)
	}
}

// applyTerrainConstraints modifies speed based on terrain conditions
func (pe *PhysicsEngine) applyTerrainConstraints(platform *models.UniversalPlatform, targetSpeed float64) float64 {
	if platform.Destination == nil {
		return targetSpeed
	}

	// Calculate gradient
	altDiff := platform.Destination.Altitude - platform.State.Position.Altitude
	horizontalDist := pe.CalculateGreatCircleDistance(
		platform.State.Position,
		models.Position{
			Latitude:  platform.Destination.Latitude,
			Longitude: platform.Destination.Longitude,
			Altitude:  platform.State.Position.Altitude,
		},
	)

	if horizontalDist > 0 {
		gradient := math.Atan(altDiff/horizontalDist) * 180.0 / math.Pi
		maxGradient := platform.TypeDef.Performance.MaxGradient
		if maxGradient == 0 {
			maxGradient = 30.0 // Default max gradient
		}

		if math.Abs(gradient) > maxGradient {
			// Reduce speed on steep gradients
			speedReduction := math.Abs(gradient) / maxGradient
			return targetSpeed / speedReduction
		}
	}

	return targetSpeed
}

// updatePosition updates the platform position based on heading and speed
func (pe *PhysicsEngine) updatePosition(state *models.PlatformState, deltaSeconds float64) {
	// Convert heading to radians (0° = North, 90° = East)
	headingRad := (90 - state.Heading) * math.Pi / 180.0

	// Calculate distance moved
	distance := state.Speed * deltaSeconds

	// Calculate position change
	deltaLat := distance * math.Cos(headingRad) / pe.EarthRadius * 180.0 / math.Pi
	deltaLon := distance * math.Sin(headingRad) / pe.EarthRadius * 180.0 / math.Pi / math.Cos(state.Position.Latitude*math.Pi/180.0)

	// Update position
	state.Position.Latitude += deltaLat
	state.Position.Longitude += deltaLon

	// Update velocity components
	state.Velocity.North = state.Speed * math.Cos(headingRad)
	state.Velocity.East = state.Speed * math.Sin(headingRad)
	state.Velocity.Up = 0 // Simplified - altitude changes handled separately
}

// getArrivalThreshold returns the distance threshold for considering arrival
func (pe *PhysicsEngine) getArrivalThreshold(platformType models.PlatformType) float64 {
	switch platformType {
	case models.PlatformTypeAirborne:
		return 100.0 // 100 meters for aircraft
	case models.PlatformTypeMaritime:
		return 50.0 // 50 meters for ships
	case models.PlatformTypeLand:
		return 10.0 // 10 meters for land vehicles
	case models.PlatformTypeSpace:
		return 1000.0 // 1km for space platforms
	default:
		return 50.0
	}
}
