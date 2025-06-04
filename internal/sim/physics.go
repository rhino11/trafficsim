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
	// Try to cast directly to UniversalPlatform for enhanced physics
	if universalPlatform, isUniversal := platform.(*models.UniversalPlatform); isUniversal {
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
	state.Velocity.Up = 0 // Will be calculated separately for altitude changes
}

// getArrivalThreshold returns the appropriate arrival threshold for each platform type
func (pe *PhysicsEngine) getArrivalThreshold(platformType models.PlatformType) float64 {
	switch platformType {
	case models.PlatformTypeAirborne:
		return 100.0 // 100 meters for aircraft
	case models.PlatformTypeMaritime:
		return 50.0 // 50 meters for ships
	case models.PlatformTypeLand:
		return 10.0 // 10 meters for land vehicles
	case models.PlatformTypeSpace:
		return 1000.0 // 1 km for satellites
	default:
		return 50.0
	}
}

// Enhanced 3D Physics Methods

// UpdatePhysics3D performs comprehensive 3D physics update including rotational dynamics
func (pe *PhysicsEngine) UpdatePhysics3D(platform *models.UniversalPlatform, deltaTime time.Duration) error {
	deltaSeconds := deltaTime.Seconds()

	// Update translational motion
	if err := pe.updateTranslationalMotion(platform, deltaSeconds); err != nil {
		return err
	}

	// Update rotational motion
	pe.updateRotationalMotion(platform, deltaSeconds)

	// Apply environmental forces
	pe.applyEnvironmentalForces(platform, deltaSeconds)

	// Update derived parameters
	pe.updateDerivedParameters(platform)

	platform.State.LastUpdated = time.Now()
	return nil
}

// updateTranslationalMotion handles 3D position and velocity updates
func (pe *PhysicsEngine) updateTranslationalMotion(platform *models.UniversalPlatform, deltaSeconds float64) error {
	// Apply forces to calculate acceleration
	totalForce := pe.calculateTotalForces(platform)

	// F = ma, so a = F/m
	mass := platform.TypeDef.Physical.Mass
	if mass <= 0 {
		mass = 1000.0 // Default mass if not specified
	}

	// Update acceleration in physics state
	platform.State.Physics.Acceleration.North = totalForce.North / mass
	platform.State.Physics.Acceleration.East = totalForce.East / mass
	platform.State.Physics.Acceleration.Up = totalForce.Up / mass

	// Update velocity using acceleration
	platform.State.Physics.Velocity.North += platform.State.Physics.Acceleration.North * deltaSeconds
	platform.State.Physics.Velocity.East += platform.State.Physics.Acceleration.East * deltaSeconds
	platform.State.Physics.Velocity.Up += platform.State.Physics.Acceleration.Up * deltaSeconds

	// Apply velocity limits
	pe.applyVelocityLimits(platform)

	// Update position using velocity
	return pe.updatePositionFromVelocity(platform, deltaSeconds)
}

// updateRotationalMotion handles angular velocity and orientation updates
func (pe *PhysicsEngine) updateRotationalMotion(platform *models.UniversalPlatform, deltaSeconds float64) {
	// Apply torques to calculate angular acceleration
	totalTorque := pe.calculateTotalTorques(platform)

	// Apply moment of inertia (simplified as uniform values)
	momentOfInertia := pe.calculateMomentOfInertia(platform)

	// Update angular acceleration (τ = Iα, so α = τ/I)
	platform.State.Physics.AngularAcceleration.RollAccel = totalTorque.RollAccel / momentOfInertia.RollAccel
	platform.State.Physics.AngularAcceleration.PitchAccel = totalTorque.PitchAccel / momentOfInertia.PitchAccel
	platform.State.Physics.AngularAcceleration.YawAccel = totalTorque.YawAccel / momentOfInertia.YawAccel

	// Update angular velocity
	platform.State.Physics.AngularVelocity.RollRate += platform.State.Physics.AngularAcceleration.RollAccel * deltaSeconds
	platform.State.Physics.AngularVelocity.PitchRate += platform.State.Physics.AngularAcceleration.PitchAccel * deltaSeconds
	platform.State.Physics.AngularVelocity.YawRate += platform.State.Physics.AngularAcceleration.YawAccel * deltaSeconds

	// Apply angular velocity limits
	pe.applyAngularVelocityLimits(platform)

	// Update orientation using angular velocity
	platform.State.Physics.Attitude.Roll += platform.State.Physics.AngularVelocity.RollRate * deltaSeconds
	platform.State.Physics.Attitude.Pitch += platform.State.Physics.AngularVelocity.PitchRate * deltaSeconds
	platform.State.Physics.Attitude.Yaw += platform.State.Physics.AngularVelocity.YawRate * deltaSeconds

	// Normalize angles to [-π, π]
	platform.State.Physics.Attitude.Roll = pe.normalizeAngle(platform.State.Physics.Attitude.Roll*math.Pi/180.0) * 180.0 / math.Pi
	platform.State.Physics.Attitude.Pitch = pe.normalizeAngle(platform.State.Physics.Attitude.Pitch*math.Pi/180.0) * 180.0 / math.Pi
	platform.State.Physics.Attitude.Yaw = pe.normalizeAngle(platform.State.Physics.Attitude.Yaw*math.Pi/180.0) * 180.0 / math.Pi
}

// calculateTotalForces computes the net force acting on the platform
func (pe *PhysicsEngine) calculateTotalForces(platform *models.UniversalPlatform) models.Acceleration {
	var totalForce models.Acceleration

	// Thrust force (in body frame, need to convert to world frame)
	thrust := pe.calculateThrust(platform)
	worldThrust := pe.bodyToWorldForce(thrust, platform.State.Physics.Attitude)
	totalForce.North += worldThrust.North
	totalForce.East += worldThrust.East
	totalForce.Up += worldThrust.Up

	// Gravity
	mass := platform.TypeDef.Physical.Mass
	if mass <= 0 {
		mass = 1000.0
	}
	totalForce.Up -= mass * pe.GravityAccel

	// Drag forces
	drag := pe.calculateDrag(platform)
	totalForce.North -= drag.North
	totalForce.East -= drag.East
	totalForce.Up -= drag.Up

	// Environmental forces (wind, current, etc.)
	if pe.EnableWeather {
		envForce := pe.calculateEnvironmentalForces(platform)
		totalForce.North += envForce.North
		totalForce.East += envForce.East
		totalForce.Up += envForce.Up
	}

	return totalForce
}

// calculateTotalTorques computes the net torque acting on the platform
func (pe *PhysicsEngine) calculateTotalTorques(platform *models.UniversalPlatform) models.AngularAcceleration {
	var totalTorque models.AngularAcceleration

	// Control surface torques
	controlTorque := pe.calculateControlTorques(platform)
	totalTorque.RollAccel += controlTorque.RollAccel
	totalTorque.PitchAccel += controlTorque.PitchAccel
	totalTorque.YawAccel += controlTorque.YawAccel

	// Aerodynamic/hydrodynamic stability torques
	stabilityTorque := pe.calculateStabilityTorques(platform)
	totalTorque.RollAccel += stabilityTorque.RollAccel
	totalTorque.PitchAccel += stabilityTorque.PitchAccel
	totalTorque.YawAccel += stabilityTorque.YawAccel

	return totalTorque
}

// Helper methods for force and torque calculations

func (pe *PhysicsEngine) calculateThrust(platform *models.UniversalPlatform) models.Acceleration {
	// Simplified thrust calculation based on desired speed
	var thrust models.Acceleration

	// Forward thrust based on throttle setting (assumed to be related to desired speed)
	maxThrust := platform.TypeDef.Performance.MaxAcceleration
	if maxThrust <= 0 {
		maxThrust = 2.0 // Default 2 m/s² acceleration capability
	}

	// Calculate desired thrust based on speed error
	desiredSpeed := platform.TypeDef.Performance.CruiseSpeed
	currentSpeed := platform.State.Speed
	speedError := desiredSpeed - currentSpeed

	// Simple proportional control
	thrustPercent := math.Max(0, math.Min(1, speedError/desiredSpeed))
	thrust.North = maxThrust * thrustPercent // Assuming forward is North for simplicity

	return thrust
}

func (pe *PhysicsEngine) calculateDrag(platform *models.UniversalPlatform) models.Acceleration {
	// Simplified drag calculation: F_drag = 0.5 * ρ * v² * C_d * A
	var drag models.Acceleration

	dragCoeff := 0.3 // Default drag coefficient
	frontalArea := platform.TypeDef.Physical.Length * platform.TypeDef.Physical.Width

	// Get fluid density based on platform type and altitude
	density := pe.getFluidDensity(platform)

	// Calculate drag force in each direction
	velocityMag := math.Sqrt(
		platform.State.Velocity.North*platform.State.Velocity.North +
			platform.State.Velocity.East*platform.State.Velocity.East +
			platform.State.Velocity.Up*platform.State.Velocity.Up,
	)

	if velocityMag > 0 {
		dragMagnitude := 0.5 * density * velocityMag * velocityMag * dragCoeff * frontalArea

		// Apply drag opposite to velocity direction
		drag.North = dragMagnitude * (platform.State.Velocity.North / velocityMag) / platform.TypeDef.Physical.Mass
		drag.East = dragMagnitude * (platform.State.Velocity.East / velocityMag) / platform.TypeDef.Physical.Mass
		drag.Up = dragMagnitude * (platform.State.Velocity.Up / velocityMag) / platform.TypeDef.Physical.Mass
	}

	return drag
}

func (pe *PhysicsEngine) calculateControlTorques(platform *models.UniversalPlatform) models.AngularAcceleration {
	// Simplified control torque calculation
	var torque models.AngularAcceleration

	// Calculate desired orientation based on movement direction
	if platform.Destination != nil {
		desiredHeading := pe.CalculateBearing(platform.State.Position, *platform.Destination)
		desiredYaw := desiredHeading * math.Pi / 180.0

		yawError := pe.normalizeAngle(desiredYaw - platform.State.Physics.Attitude.Yaw*math.Pi/180.0)

		// Simple proportional control for yaw
		yawGain := 1.0
		torque.YawAccel = yawGain * yawError
	}

	return torque
}

func (pe *PhysicsEngine) calculateStabilityTorques(platform *models.UniversalPlatform) models.AngularAcceleration {
	// Stability torques tend to return the platform to level flight
	var torque models.AngularAcceleration

	rollDamping := 1.0
	pitchDamping := 1.0
	yawDamping := 1.0

	// Apply damping proportional to angular velocity
	torque.RollAccel = -rollDamping * platform.State.Physics.AngularVelocity.RollRate
	torque.PitchAccel = -pitchDamping * platform.State.Physics.AngularVelocity.PitchRate
	torque.YawAccel = -yawDamping * platform.State.Physics.AngularVelocity.YawRate

	// Add restoring torques for roll and pitch
	torque.RollAccel -= rollDamping * platform.State.Physics.Attitude.Roll * 0.5
	torque.PitchAccel -= pitchDamping * platform.State.Physics.Attitude.Pitch * 0.5

	return torque
}

// Utility methods

func (pe *PhysicsEngine) bodyToWorldForce(bodyForce models.Acceleration, orientation models.Attitude) models.Acceleration {
	// Simplified rotation transformation (assuming small angles)
	var worldForce models.Acceleration

	// For small angles, we can use simplified rotation
	yawRad := orientation.Yaw * math.Pi / 180.0
	worldForce.North = bodyForce.North*math.Cos(yawRad) - bodyForce.East*math.Sin(yawRad)
	worldForce.East = bodyForce.North*math.Sin(yawRad) + bodyForce.East*math.Cos(yawRad)
	worldForce.Up = bodyForce.Up

	return worldForce
}

func (pe *PhysicsEngine) getFluidDensity(platform *models.UniversalPlatform) float64 {
	switch platform.PlatformType {
	case models.PlatformTypeAirborne, models.PlatformTypeSpace:
		// Air density decreases with altitude
		altitude := platform.State.Position.Altitude
		if altitude < 11000 { // Troposphere
			return pe.AirDensity * math.Pow(1-0.0065*altitude/288.15, 4.255)
		}
		return pe.AirDensity * 0.1 // Simplified for higher altitudes
	case models.PlatformTypeMaritime:
		return 1025.0 // Seawater density (kg/m³)
	case models.PlatformTypeLand:
		return pe.AirDensity // Ground vehicles in air
	default:
		return pe.AirDensity
	}
}

func (pe *PhysicsEngine) calculateMomentOfInertia(platform *models.UniversalPlatform) models.AngularAcceleration {
	// Simplified moment of inertia calculation
	mass := platform.TypeDef.Physical.Mass
	length := platform.TypeDef.Physical.Length
	width := platform.TypeDef.Physical.Width
	height := platform.TypeDef.Physical.Height

	if mass <= 0 {
		mass = 1000.0
	}
	if length <= 0 {
		length = 10.0
	}
	if width <= 0 {
		width = 2.0
	}
	if height <= 0 {
		height = 2.0
	}

	// Approximate as rectangular solid
	return models.AngularAcceleration{
		RollAccel:  mass * (height*height + length*length) / 12.0,
		PitchAccel: mass * (width*width + height*height) / 12.0,
		YawAccel:   mass * (length*length + width*width) / 12.0,
	}
}

func (pe *PhysicsEngine) applyVelocityLimits(platform *models.UniversalPlatform) {
	maxSpeed := platform.TypeDef.Performance.MaxSpeed
	if maxSpeed <= 0 {
		return
	}

	// Calculate total velocity magnitude
	velMag := math.Sqrt(
		platform.State.Velocity.North*platform.State.Velocity.North +
			platform.State.Velocity.East*platform.State.Velocity.East +
			platform.State.Velocity.Up*platform.State.Velocity.Up,
	)

	// Scale down if exceeding max speed
	if velMag > maxSpeed {
		scale := maxSpeed / velMag
		platform.State.Velocity.North *= scale
		platform.State.Velocity.East *= scale
		platform.State.Velocity.Up *= scale
	}
}

func (pe *PhysicsEngine) applyAngularVelocityLimits(platform *models.UniversalPlatform) {
	maxRollRate := platform.TypeDef.Performance.MaxRollRate
	maxPitchRate := platform.TypeDef.Performance.MaxPitchRate
	maxYawRate := platform.TypeDef.Performance.MaxYawRate

	if maxRollRate > 0 {
		platform.State.Physics.AngularVelocity.RollRate = math.Max(-maxRollRate, math.Min(maxRollRate, platform.State.Physics.AngularVelocity.RollRate))
	}
	if maxPitchRate > 0 {
		platform.State.Physics.AngularVelocity.PitchRate = math.Max(-maxPitchRate, math.Min(maxPitchRate, platform.State.Physics.AngularVelocity.PitchRate))
	}
	if maxYawRate > 0 {
		platform.State.Physics.AngularVelocity.YawRate = math.Max(-maxYawRate, math.Min(maxYawRate, platform.State.Physics.AngularVelocity.YawRate))
	}
}

func (pe *PhysicsEngine) updatePositionFromVelocity(platform *models.UniversalPlatform, deltaSeconds float64) error {
	// Convert velocity to position change
	northDistance := platform.State.Velocity.North * deltaSeconds
	eastDistance := platform.State.Velocity.East * deltaSeconds
	upDistance := platform.State.Velocity.Up * deltaSeconds

	// Convert to lat/lon changes
	deltaLat := northDistance / pe.EarthRadius * 180.0 / math.Pi
	deltaLon := eastDistance / pe.EarthRadius * 180.0 / math.Pi / math.Cos(platform.State.Position.Latitude*math.Pi/180.0)

	// Update position
	platform.State.Position.Latitude += deltaLat
	platform.State.Position.Longitude += deltaLon
	platform.State.Position.Altitude += upDistance

	// Update speed and heading for compatibility
	platform.State.Speed = math.Sqrt(
		platform.State.Velocity.North*platform.State.Velocity.North +
			platform.State.Velocity.East*platform.State.Velocity.East,
	)

	if platform.State.Speed > 0 {
		platform.State.Heading = math.Atan2(platform.State.Velocity.East, platform.State.Velocity.North) * 180.0 / math.Pi
		platform.State.Heading = math.Mod(platform.State.Heading+360, 360)
	}

	return nil
}

// Enhanced physics calculations that integrate with 3D physics state

// Calculate3DForces computes comprehensive forces for any platform type
func (pe *PhysicsEngine) Calculate3DForces(platform *models.UniversalPlatform) error {
	physics := &platform.State.Physics

	// Environmental forces
	physics.Forces.Weight = physics.Mass * pe.GravityAccel

	// Platform-specific force calculations
	switch platform.PlatformType {
	case models.PlatformTypeAirborne:
		return pe.calculateAerodynamicForces(platform, physics)
	case models.PlatformTypeMaritime:
		return pe.calculateHydrodynamicForces(platform, physics)
	case models.PlatformTypeLand:
		return pe.calculateGroundForces(platform, physics)
	case models.PlatformTypeSpace:
		return pe.calculateOrbitalForces(platform, physics)
	default:
		return pe.calculateGenericForces(platform, physics)
	}
}

// calculateAerodynamicForces handles aircraft-specific forces
func (pe *PhysicsEngine) calculateAerodynamicForces(platform *models.UniversalPlatform, physics *models.PhysicsState) error {
	// Air density at altitude (simplified exponential model)
	altitude := physics.Position.Altitude
	airDensity := pe.AirDensity * math.Exp(-altitude/8400) // Scale height ~8.4km

	// Current speed and velocity
	velocity := math.Sqrt(physics.Velocity.North*physics.Velocity.North +
		physics.Velocity.East*physics.Velocity.East +
		physics.Velocity.Up*physics.Velocity.Up)

	// Wing area and aerodynamic coefficients
	wingArea := platform.TypeDef.Physical.WingArea
	if wingArea == 0 {
		// Estimate based on mass (rule of thumb)
		wingArea = platform.TypeDef.Physical.Mass / 500 // kg/m²
	}

	// Dynamic pressure
	dynamicPressure := 0.5 * airDensity * velocity * velocity

	// Lift calculation (simplified)
	angleOfAttack := physics.Attitude.Pitch * math.Pi / 180.0
	liftCoeff := 2 * math.Pi * math.Sin(angleOfAttack) // Simplified thin airfoil theory
	physics.Forces.Lift = dynamicPressure * wingArea * liftCoeff

	// Drag calculation
	inducedDragCoeff := liftCoeff * liftCoeff / (math.Pi * 8) // Simplified induced drag
	parasiteDragCoeff := 0.02                                 // Typical value
	totalDragCoeff := parasiteDragCoeff + inducedDragCoeff
	physics.Forces.Drag = dynamicPressure * wingArea * totalDragCoeff

	// Thrust calculation based on throttle setting
	maxThrust := physics.Mass * pe.GravityAccel * platform.TypeDef.Performance.MaxThrustToWeight
	if maxThrust == 0 {
		maxThrust = physics.Mass * pe.GravityAccel * 0.3 // Default T/W ratio
	}

	// Throttle based on speed error
	targetSpeed := platform.TypeDef.Performance.CruiseSpeed
	speedError := targetSpeed - velocity
	throttle := math.Max(0, math.Min(1, 0.5+speedError/targetSpeed))
	physics.Forces.Thrust = maxThrust * throttle

	return nil
}

// calculateHydrodynamicForces handles ship-specific forces
func (pe *PhysicsEngine) calculateHydrodynamicForces(platform *models.UniversalPlatform, physics *models.PhysicsState) error {
	// Water density (constant at sea level)
	waterDensity := 1025.0 // kg/m³

	velocity := math.Sqrt(physics.Velocity.North*physics.Velocity.North +
		physics.Velocity.East*physics.Velocity.East)

	// Hydrodynamic calculations
	wetArea := platform.TypeDef.Physical.WetArea
	if wetArea == 0 {
		// Estimate based on length and beam
		wetArea = platform.TypeDef.Physical.Length * platform.TypeDef.Physical.Width
	}

	// Wave resistance (simplified Froude number calculation)
	froudeNumber := velocity / math.Sqrt(pe.GravityAccel*platform.TypeDef.Physical.Length)
	waveResistanceCoeff := 0.002 * (1 + 10*froudeNumber*froudeNumber)

	// Total drag
	viscousDragCoeff := 0.01 // Simplified
	totalDragCoeff := viscousDragCoeff + waveResistanceCoeff
	physics.Forces.Drag = 0.5 * waterDensity * velocity * velocity * wetArea * totalDragCoeff

	// Thrust from propulsion
	targetSpeed := platform.TypeDef.Performance.CruiseSpeed
	speedError := targetSpeed - velocity
	maxThrust := physics.Mass * 2.0 // Ships have high thrust capability
	throttle := math.Max(0, math.Min(1, 0.5+speedError/targetSpeed))
	physics.Forces.Thrust = maxThrust * throttle

	// Buoyancy equals weight
	physics.Forces.Normal = physics.Forces.Weight

	return nil
}

// calculateGroundForces handles land vehicle forces
func (pe *PhysicsEngine) calculateGroundForces(platform *models.UniversalPlatform, physics *models.PhysicsState) error {
	velocity := math.Sqrt(physics.Velocity.North*physics.Velocity.North +
		physics.Velocity.East*physics.Velocity.East)

	// Aerodynamic drag
	frontalArea := platform.TypeDef.Physical.FrontalArea
	if frontalArea == 0 {
		frontalArea = platform.TypeDef.Physical.Width * platform.TypeDef.Physical.Height
	}

	dragCoeff := 0.3 // Typical for vehicles
	physics.Forces.Drag = 0.5 * pe.AirDensity * velocity * velocity * frontalArea * dragCoeff

	// Rolling resistance
	rollingResistance := physics.Forces.Weight * 0.01 // 1% of weight
	physics.Forces.Drag += rollingResistance

	// Traction force from drive system
	targetSpeed := platform.TypeDef.Performance.CruiseSpeed
	speedError := targetSpeed - velocity
	maxTraction := physics.Mass * 5.0 // High traction capability
	throttle := math.Max(0, math.Min(1, 0.5+speedError/targetSpeed))
	physics.Forces.Thrust = maxTraction * throttle

	// Normal force from ground
	physics.Forces.Normal = physics.Forces.Weight

	return nil
}

// calculateOrbitalForces handles space platform forces
func (pe *PhysicsEngine) calculateOrbitalForces(platform *models.UniversalPlatform, physics *models.PhysicsState) error {
	// In orbit, forces are minimal except for station keeping
	altitude := physics.Position.Altitude

	if altitude > 100000 { // Above atmosphere
		// Minimal atmospheric drag
		physics.Forces.Drag = 0

		// Station keeping thrusters (very small)
		physics.Forces.Thrust = 0.01 // Minimal thrust for attitude control

		// No aerodynamic lift
		physics.Forces.Lift = 0

		// Gravitational force
		r := pe.EarthRadius + altitude
		gravitationalAccel := pe.GravityAccel * (pe.EarthRadius * pe.EarthRadius) / (r * r)
		physics.Forces.Weight = physics.Mass * gravitationalAccel

		// Centrifugal force balances gravity in circular orbit
		orbitalSpeed := math.Sqrt(pe.GravityAccel * pe.EarthRadius * pe.EarthRadius / r)
		centrifugalForce := physics.Mass * orbitalSpeed * orbitalSpeed / r
		physics.Forces.Normal = centrifugalForce
	}

	return nil
}

// calculateGenericForces provides fallback force calculations
func (pe *PhysicsEngine) calculateGenericForces(platform *models.UniversalPlatform, physics *models.PhysicsState) error {
	// Basic force model
	velocity := math.Sqrt(physics.Velocity.North*physics.Velocity.North +
		physics.Velocity.East*physics.Velocity.East)

	// Simple drag proportional to velocity squared
	physics.Forces.Drag = 0.5 * pe.AirDensity * velocity * velocity * 10.0

	// Simple thrust control
	targetSpeed := platform.TypeDef.Performance.CruiseSpeed
	speedError := targetSpeed - velocity
	maxThrust := physics.Mass * 3.0
	throttle := math.Max(0, math.Min(1, 0.5+speedError/targetSpeed))
	physics.Forces.Thrust = maxThrust * throttle

	return nil
}

// Calculate3DTorques computes rotational forces for platform control
func (pe *PhysicsEngine) Calculate3DTorques(platform *models.UniversalPlatform) error {
	physics := &platform.State.Physics

	// Control torques based on desired vs actual attitude
	desiredAttitude := pe.calculateDesiredAttitude(platform)

	// PID-like control for each axis
	rollError := desiredAttitude.Roll - physics.Attitude.Roll
	pitchError := desiredAttitude.Pitch - physics.Attitude.Pitch
	yawError := desiredAttitude.Yaw - physics.Attitude.Yaw

	// Normalize yaw error
	for yawError > 180 {
		yawError -= 360
	}
	for yawError < -180 {
		yawError += 360
	}

	// Control authority based on platform type
	var maxRollTorque, maxPitchTorque, maxYawTorque float64

	switch platform.PlatformType {
	case models.PlatformTypeAirborne:
		// Aircraft control surfaces
		maxRollTorque = physics.Mass * 10.0 // Ailerons
		maxPitchTorque = physics.Mass * 8.0 // Elevator
		maxYawTorque = physics.Mass * 5.0   // Rudder
	case models.PlatformTypeMaritime:
		// Ship rudder (primarily yaw control)
		maxRollTorque = physics.Mass * 1.0
		maxPitchTorque = physics.Mass * 1.0
		maxYawTorque = physics.Mass * 3.0
	case models.PlatformTypeLand:
		// Vehicle steering and suspension
		maxRollTorque = physics.Mass * 2.0
		maxPitchTorque = physics.Mass * 2.0
		maxYawTorque = physics.Mass * 8.0 // Steering
	case models.PlatformTypeSpace:
		// Reaction wheels and thrusters
		maxRollTorque = physics.Mass * 0.1
		maxPitchTorque = physics.Mass * 0.1
		maxYawTorque = physics.Mass * 0.1
	default:
		maxRollTorque = physics.Mass * 2.0
		maxPitchTorque = physics.Mass * 2.0
		maxYawTorque = physics.Mass * 2.0
	}

	// Proportional control (simplified)
	kp := 0.1
	physics.Torques.Roll = kp * rollError * maxRollTorque
	physics.Torques.Pitch = kp * pitchError * maxPitchTorque
	physics.Torques.Yaw = kp * yawError * maxYawTorque

	// Limit torques
	physics.Torques.Roll = math.Max(-maxRollTorque, math.Min(maxRollTorque, physics.Torques.Roll))
	physics.Torques.Pitch = math.Max(-maxPitchTorque, math.Min(maxPitchTorque, physics.Torques.Pitch))
	physics.Torques.Yaw = math.Max(-maxYawTorque, math.Min(maxYawTorque, physics.Torques.Yaw))

	return nil
}

// calculateDesiredAttitude determines the target attitude for navigation
func (pe *PhysicsEngine) calculateDesiredAttitude(platform *models.UniversalPlatform) models.Attitude {
	if platform.Destination == nil {
		return platform.State.Physics.Attitude
	}

	// Calculate bearing to destination
	bearing := pe.CalculateBearing(platform.State.Position, *platform.Destination)

	// Calculate desired pitch for altitude change
	altitudeDiff := platform.Destination.Altitude - platform.State.Position.Altitude
	distance := pe.CalculateGreatCircleDistance(platform.State.Position, *platform.Destination)

	desiredPitch := 0.0
	if distance > 0 {
		pitchAngle := math.Atan(altitudeDiff/distance) * 180.0 / math.Pi

		// Limit pitch based on platform capabilities
		maxPitch := platform.TypeDef.Performance.MaxPitchAngle
		if maxPitch == 0 {
			switch platform.PlatformType {
			case models.PlatformTypeAirborne:
				maxPitch = 15.0
			case models.PlatformTypeLand:
				maxPitch = 30.0
			default:
				maxPitch = 10.0
			}
		}
		desiredPitch = math.Max(-maxPitch, math.Min(maxPitch, pitchAngle))
	}

	// Calculate desired roll for coordinated turns (aircraft)
	desiredRoll := 0.0
	if platform.PlatformType == models.PlatformTypeAirborne {
		headingError := bearing - platform.State.Physics.Attitude.Yaw
		for headingError > 180 {
			headingError -= 360
		}
		for headingError < -180 {
			headingError += 360
		}

		maxBank := platform.TypeDef.Performance.MaxBankAngle
		if maxBank == 0 {
			maxBank = 30.0
		}
		desiredRoll = math.Max(-maxBank, math.Min(maxBank, headingError/3))
	}

	return models.Attitude{
		Roll:  desiredRoll,
		Pitch: desiredPitch,
		Yaw:   bearing,
	}
}

// Enhanced environmental effects

// ApplyEnvironmentalEffects modifies forces based on environmental conditions
func (pe *PhysicsEngine) ApplyEnvironmentalEffects(platform *models.UniversalPlatform, weather *WeatherConditions) error {
	if !pe.EnableWeather || weather == nil {
		return nil
	}

	physics := &platform.State.Physics

	// Wind effects
	windEffect := pe.calculateWindEffect(platform, weather)
	physics.Forces.Drag += windEffect.Drag
	physics.Torques.Roll += windEffect.RollTorque
	physics.Torques.Pitch += windEffect.PitchTorque
	physics.Torques.Yaw += windEffect.YawTorque

	// Temperature effects on performance
	tempRatio := (weather.Temperature + 273.15) / 288.15 // ISA standard
	switch platform.PlatformType {
	case models.PlatformTypeAirborne:
		// Engine performance varies with temperature
		physics.Forces.Thrust *= math.Sqrt(tempRatio)
	}

	// Precipitation effects
	if weather.Precipitation > 0 {
		// Increased drag due to rain/snow
		physics.Forces.Drag *= (1 + weather.Precipitation*0.1)
	}

	return nil
}

// WeatherConditions represents environmental conditions
type WeatherConditions struct {
	WindSpeed     float64 // m/s
	WindDirection float64 // degrees
	Temperature   float64 // Celsius
	Pressure      float64 // Pa
	Humidity      float64 // %
	Precipitation float64 // mm/hr
	Visibility    float64 // meters
}

// WindEffect represents wind-induced forces and torques
type WindEffect struct {
	Drag        float64
	RollTorque  float64
	PitchTorque float64
	YawTorque   float64
}

// calculateWindEffect computes wind effects on platform
func (pe *PhysicsEngine) calculateWindEffect(platform *models.UniversalPlatform, weather *WeatherConditions) WindEffect {
	// Convert wind from meteorological to mathematical convention
	windDir := math.Mod(270-weather.WindDirection, 360) * math.Pi / 180.0

	// Wind velocity components
	windNorth := weather.WindSpeed * math.Cos(windDir)
	windEast := weather.WindSpeed * math.Sin(windDir)

	// Relative wind
	relWindNorth := windNorth - platform.State.Physics.Velocity.North
	relWindEast := windEast - platform.State.Physics.Velocity.East
	relWindSpeed := math.Sqrt(relWindNorth*relWindNorth + relWindEast*relWindEast)

	// Wind effects scale with exposed area and relative wind speed
	var exposedArea float64
	switch platform.PlatformType {
	case models.PlatformTypeAirborne:
		exposedArea = platform.TypeDef.Physical.WingArea
	case models.PlatformTypeMaritime:
		exposedArea = platform.TypeDef.Physical.Length * platform.TypeDef.Physical.Height
	case models.PlatformTypeLand:
		exposedArea = platform.TypeDef.Physical.FrontalArea
	default:
		exposedArea = platform.TypeDef.Physical.Width * platform.TypeDef.Physical.Height
	}

	if exposedArea == 0 {
		exposedArea = 10.0 // Default
	}

	// Wind drag
	windDrag := 0.5 * pe.AirDensity * relWindSpeed * relWindSpeed * exposedArea * 0.5

	// Wind-induced torques (simplified)
	leverArm := platform.TypeDef.Physical.Length / 2
	windTorque := windDrag * leverArm * 0.1

	return WindEffect{
		Drag:        windDrag,
		RollTorque:  windTorque * 0.3,
		PitchTorque: windTorque * 0.2,
		YawTorque:   windTorque * 0.5,
	}
}

// Missing helper functions for physics calculations

// calculateEnvironmentalForces computes environmental forces like wind and current
func (pe *PhysicsEngine) calculateEnvironmentalForces(platform *models.UniversalPlatform) models.Acceleration {
	// Simplified environmental forces - can be expanded for weather effects
	var envForce models.Acceleration

	// For now, return zero forces - this can be enhanced with weather data
	return envForce
}

// applyEnvironmentalForces applies environmental effects to the platform
func (pe *PhysicsEngine) applyEnvironmentalForces(platform *models.UniversalPlatform, deltaSeconds float64) {
	// Apply basic environmental effects like atmospheric density changes
	// This is a placeholder for more complex environmental modeling
}

// updateDerivedParameters updates calculated values based on current state
func (pe *PhysicsEngine) updateDerivedParameters(platform *models.UniversalPlatform) {
	// Update speed from velocity components
	platform.State.Speed = math.Sqrt(
		platform.State.Velocity.North*platform.State.Velocity.North +
			platform.State.Velocity.East*platform.State.Velocity.East,
	)

	// Update heading from velocity direction
	if platform.State.Speed > 0 {
		platform.State.Heading = math.Atan2(platform.State.Velocity.East, platform.State.Velocity.North) * 180.0 / math.Pi
		platform.State.Heading = math.Mod(platform.State.Heading+360, 360)
	}
}

// normalizeAngle normalizes an angle to the range [-π, π]
func (pe *PhysicsEngine) normalizeAngle(angle float64) float64 {
	for angle > math.Pi {
		angle -= 2 * math.Pi
	}
	for angle < -math.Pi {
		angle += 2 * math.Pi
	}
	return angle
}
