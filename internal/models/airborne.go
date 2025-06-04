package models

import (
	"fmt"
	"math"
	"time"
)

// AirbornePlatform represents aircraft platforms with enhanced physics
type AirbornePlatform struct {
	// Embed UniversalPlatform for base functionality
	UniversalPlatform

	// Enhanced flight characteristics (aircraft-specific)
	MaxRollRate     float64 // degrees/second
	MaxPitchRate    float64 // degrees/second
	MaxYawRate      float64 // degrees/second
	MaxBankAngle    float64 // degrees
	MaxPitchAngle   float64 // degrees
	MaxLoadFactor   float64 // g-force
	StallSpeed      float64 // m/s
	MaxAcceleration float64 // m/s²
	MaxDeceleration float64 // m/s²

	// Physical characteristics (additional to base)
	WingArea        float64         // m²
	WingLoading     float64         // kg/m²
	CenterOfGravity CenterOfGravity // center of mass location

	// Flight state
	FlightPhase FlightPhase // takeoff, climb, cruise, descent, approach, landing
}

// FlightPhase represents the current phase of flight
type FlightPhase string

const (
	FlightPhaseTakeoff  FlightPhase = "takeoff"
	FlightPhaseClimb    FlightPhase = "climb"
	FlightPhaseCruise   FlightPhase = "cruise"
	FlightPhaseDescent  FlightPhase = "descent"
	FlightPhaseApproach FlightPhase = "approach"
	FlightPhaseLanding  FlightPhase = "landing"
	FlightPhaseParked   FlightPhase = "parked"
)

// Core Platform interface implementation - delegate to embedded UniversalPlatform
func (a *AirbornePlatform) GetID() string           { return a.UniversalPlatform.GetID() }
func (a *AirbornePlatform) GetType() PlatformType   { return a.UniversalPlatform.GetType() }
func (a *AirbornePlatform) GetClass() string        { return a.UniversalPlatform.GetClass() }
func (a *AirbornePlatform) GetName() string         { return a.UniversalPlatform.GetName() }
func (a *AirbornePlatform) GetCallSign() string     { return a.UniversalPlatform.GetCallSign() }
func (a *AirbornePlatform) GetState() PlatformState { return a.UniversalPlatform.GetState() }
func (a *AirbornePlatform) GetMaxSpeed() float64    { return a.UniversalPlatform.GetMaxSpeed() }
func (a *AirbornePlatform) GetMaxAltitude() float64 { return a.UniversalPlatform.GetMaxAltitude() }
func (a *AirbornePlatform) GetLength() float64      { return a.UniversalPlatform.GetLength() }
func (a *AirbornePlatform) GetWidth() float64       { return a.UniversalPlatform.GetWidth() }
func (a *AirbornePlatform) GetHeight() float64      { return a.UniversalPlatform.GetHeight() }
func (a *AirbornePlatform) GetMass() float64        { return a.UniversalPlatform.GetMass() }

func (a *AirbornePlatform) UpdateState(state PlatformState) {
	a.UniversalPlatform.UpdateState(state)
}

func (a *AirbornePlatform) SetDestination(pos Position) error {
	return a.UniversalPlatform.SetDestination(pos)
}

// Enhanced 3D physics methods
func (a *AirbornePlatform) Initialize3DPhysics() {
	// Initialize physics state with aircraft characteristics
	a.UniversalPlatform.State.Physics = PhysicsState{
		Position:            a.UniversalPlatform.State.Position,
		Velocity:            a.UniversalPlatform.State.Velocity,
		Acceleration:        Acceleration{},
		Attitude:            Attitude{Yaw: a.UniversalPlatform.State.Heading},
		AngularVelocity:     AngularVelocity{},
		AngularAcceleration: AngularAcceleration{},
		Mass:                a.UniversalPlatform.GetMass(),
		MomentOfInertia:     a.calculateMomentOfInertia(),
		Forces:              Forces{Weight: a.UniversalPlatform.GetMass() * 9.81},
		Torques:             Torques{},
	}
}

func (a *AirbornePlatform) calculateMomentOfInertia() MomentOfInertia {
	// Calculate realistic moment of inertia for aircraft
	mass := a.UniversalPlatform.GetMass()
	length := a.UniversalPlatform.GetLength()
	width := a.UniversalPlatform.GetWidth()
	height := a.UniversalPlatform.GetHeight()

	return MomentOfInertia{
		Ixx: mass * (width*width + height*height) / 20,   // Roll axis
		Iyy: mass * (length*length + height*height) / 20, // Pitch axis
		Izz: mass * (length*length + width*width) / 20,   // Yaw axis
	}
}

func (a *AirbornePlatform) Update3DPhysics(deltaTime time.Duration) error {
	// Use the existing enhanced Update method which already includes 3D physics
	return a.Update(deltaTime)
}

func (a *AirbornePlatform) GetPhysicsState() PhysicsState {
	return a.UniversalPlatform.State.Physics
}

func (a *AirbornePlatform) SetPhysicsState(physics PhysicsState) {
	a.UniversalPlatform.State.Physics = physics
	// Update legacy compatibility fields
	a.UniversalPlatform.State.Position = physics.Position
	a.UniversalPlatform.State.Velocity = physics.Velocity
	a.UniversalPlatform.State.Heading = physics.Attitude.Yaw
}

// Update implements enhanced aircraft physics and flight dynamics
func (a *AirbornePlatform) Update(deltaTime time.Duration) error {
	if a.UniversalPlatform.Destination == nil {
		return nil
	}

	dt := deltaTime.Seconds()

	// Calculate 3D distance and bearing to destination
	deltaLat := a.UniversalPlatform.Destination.Latitude - a.UniversalPlatform.State.Position.Latitude
	deltaLon := a.UniversalPlatform.Destination.Longitude - a.UniversalPlatform.State.Position.Longitude
	deltaAlt := a.UniversalPlatform.Destination.Altitude - a.UniversalPlatform.State.Position.Altitude

	distance := math.Sqrt(deltaLat*deltaLat + deltaLon*deltaLon)
	distance3D := math.Sqrt(distance*distance + (deltaAlt/111320)*(deltaAlt/111320))

	if distance3D < 0.001 { // Close enough
		a.UniversalPlatform.Destination = nil
		return nil
	}

	// Initialize physics if not already done
	if a.UniversalPlatform.State.Physics.Mass == 0 {
		a.Initialize3DPhysics()
	}

	// Enhanced flight physics
	a.updateFlightPhase()
	a.updateFlightDynamics(dt)

	// Apply aerodynamic forces and moments
	a.calculateAerodynamicForces()
	a.updateAttitude(dt)

	// Update position based on velocity
	a.updatePosition(dt)

	a.UniversalPlatform.State.LastUpdated = time.Now()
	return nil
}

// updateFlightPhase determines the current phase of flight
func (a *AirbornePlatform) updateFlightPhase() {
	// If we have a destination and are at low altitude with low speed, start takeoff
	if a.UniversalPlatform.Destination != nil && a.UniversalPlatform.State.Position.Altitude < 100 && a.UniversalPlatform.State.Speed < 50 {
		a.FlightPhase = FlightPhaseTakeoff
	} else if a.UniversalPlatform.State.Position.Altitude < 100 && a.UniversalPlatform.State.Speed < 50 {
		a.FlightPhase = FlightPhaseParked
	} else if a.UniversalPlatform.State.Position.Altitude < 300 && a.UniversalPlatform.State.Velocity.Up > 0 {
		a.FlightPhase = FlightPhaseTakeoff
	} else if a.UniversalPlatform.State.Velocity.Up > 5 {
		a.FlightPhase = FlightPhaseClimb
	} else if a.UniversalPlatform.State.Velocity.Up < -5 {
		a.FlightPhase = FlightPhaseDescent
	} else if a.UniversalPlatform.State.Position.Altitude > 1000 {
		a.FlightPhase = FlightPhaseCruise
	} else {
		a.FlightPhase = FlightPhaseApproach
	}
}

// updateFlightDynamics applies realistic flight dynamics
func (a *AirbornePlatform) updateFlightDynamics(dt float64) {
	// Calculate target velocity based on flight phase
	var targetSpeed float64
	var targetAltitude float64

	cruiseSpeed := a.UniversalPlatform.GetMaxSpeed() * 0.9  // Use 90% of max speed as cruise
	cruiseAlt := a.UniversalPlatform.GetMaxAltitude() * 0.8 // Use 80% of max altitude as cruise

	switch a.FlightPhase {
	case FlightPhaseTakeoff:
		targetSpeed = math.Min(cruiseSpeed*0.8, a.UniversalPlatform.State.Speed+a.MaxAcceleration*dt)
		targetAltitude = a.UniversalPlatform.State.Position.Altitude + 10*dt // 10 m/s climb rate
	case FlightPhaseClimb:
		targetSpeed = cruiseSpeed
		targetAltitude = math.Min(a.UniversalPlatform.Destination.Altitude, a.UniversalPlatform.State.Position.Altitude+15*dt)
	case FlightPhaseCruise:
		targetSpeed = cruiseSpeed
		targetAltitude = cruiseAlt
	case FlightPhaseDescent:
		targetSpeed = cruiseSpeed * 0.9
		targetAltitude = math.Max(a.UniversalPlatform.Destination.Altitude, a.UniversalPlatform.State.Position.Altitude-12*dt)
	case FlightPhaseApproach:
		targetSpeed = math.Max(a.StallSpeed*1.3, cruiseSpeed*0.6)
		targetAltitude = a.UniversalPlatform.Destination.Altitude
	default:
		targetSpeed = 0
		targetAltitude = a.UniversalPlatform.State.Position.Altitude
	}

	// Apply acceleration constraints
	speedDiff := targetSpeed - a.UniversalPlatform.State.Speed
	maxSpeedChange := a.MaxAcceleration * dt

	if math.Abs(speedDiff) <= maxSpeedChange {
		a.UniversalPlatform.State.Speed = targetSpeed
	} else if speedDiff > 0 {
		a.UniversalPlatform.State.Speed += maxSpeedChange
	} else {
		a.UniversalPlatform.State.Speed -= math.Min(maxSpeedChange, a.MaxDeceleration*dt)
	}

	// Ensure minimum flying speed
	if a.UniversalPlatform.State.Position.Altitude > 50 && a.UniversalPlatform.State.Speed < a.StallSpeed {
		a.UniversalPlatform.State.Speed = a.StallSpeed
	}

	// Calculate climb/descent rate
	altDiff := targetAltitude - a.UniversalPlatform.State.Position.Altitude
	a.UniversalPlatform.State.Velocity.Up = math.Max(-15, math.Min(15, altDiff/dt)) // Limit vertical speed
}

// calculateAerodynamicForces computes forces acting on the aircraft
func (a *AirbornePlatform) calculateAerodynamicForces() {
	// Air density at altitude (simplified)
	airDensity := 1.225 * math.Exp(-a.UniversalPlatform.State.Position.Altitude/8400) // kg/m³

	// Dynamic pressure
	dynamicPressure := 0.5 * airDensity * a.UniversalPlatform.State.Speed * a.UniversalPlatform.State.Speed

	// Lift force (simplified)
	liftCoeff := 0.8 // Typical cruise lift coefficient
	a.UniversalPlatform.State.Physics.Forces.Lift = liftCoeff * dynamicPressure * a.WingArea

	// Drag force
	dragCoeff := 0.025 + (liftCoeff*liftCoeff)/(math.Pi*8*0.8) // Induced drag
	a.UniversalPlatform.State.Physics.Forces.Drag = dragCoeff * dynamicPressure * a.WingArea

	// Weight
	gravity := 9.81
	a.UniversalPlatform.State.Physics.Forces.Weight = a.UniversalPlatform.GetMass() * gravity

	// Thrust (to maintain speed)
	a.UniversalPlatform.State.Physics.Forces.Thrust = a.UniversalPlatform.State.Physics.Forces.Drag

	// Calculate accelerations
	mass := a.UniversalPlatform.GetMass()
	a.UniversalPlatform.State.Physics.Acceleration.North = (a.UniversalPlatform.State.Physics.Forces.Thrust - a.UniversalPlatform.State.Physics.Forces.Drag) / mass
	a.UniversalPlatform.State.Physics.Acceleration.Up = (a.UniversalPlatform.State.Physics.Forces.Lift - a.UniversalPlatform.State.Physics.Forces.Weight) / mass
}

// updateAttitude updates aircraft orientation based on flight dynamics
func (a *AirbornePlatform) updateAttitude(dt float64) {
	// Calculate desired bank angle for turns
	if a.UniversalPlatform.Destination != nil {
		desiredHeading := a.calculateBearing(a.UniversalPlatform.State.Position, *a.UniversalPlatform.Destination)
		headingError := desiredHeading - a.UniversalPlatform.State.Heading

		// Normalize heading error to [-180, 180]
		for headingError > 180 {
			headingError -= 360
		}
		for headingError < -180 {
			headingError += 360
		}

		// Calculate bank angle for coordinated turn - use a proportional controller
		// Larger heading errors require more bank angle
		maxBankForTurn := math.Min(a.MaxBankAngle, 30) // Limit to 30 degrees for normal flight

		// Use a gain factor to make the aircraft more responsive to heading errors
		bankGain := 1.5 // Adjust this to make turns more or less aggressive
		desiredBank := math.Max(-maxBankForTurn, math.Min(maxBankForTurn, headingError*bankGain))

		// Only bank if there's a significant heading error (more than 2 degrees)
		if math.Abs(headingError) > 2.0 {
			// Apply roll rate limits
			rollError := desiredBank - a.UniversalPlatform.State.Physics.Attitude.Roll
			maxRollChange := a.MaxRollRate * dt

			if math.Abs(rollError) <= maxRollChange {
				a.UniversalPlatform.State.Physics.Attitude.Roll = desiredBank
			} else if rollError > 0 {
				a.UniversalPlatform.State.Physics.Attitude.Roll += maxRollChange
			} else {
				a.UniversalPlatform.State.Physics.Attitude.Roll -= maxRollChange
			}
		} else {
			// Level off if close to desired heading
			rollError := -a.UniversalPlatform.State.Physics.Attitude.Roll
			maxRollChange := a.MaxRollRate * dt

			if math.Abs(rollError) <= maxRollChange {
				a.UniversalPlatform.State.Physics.Attitude.Roll = 0
			} else if rollError > 0 {
				a.UniversalPlatform.State.Physics.Attitude.Roll += maxRollChange
			} else {
				a.UniversalPlatform.State.Physics.Attitude.Roll -= maxRollChange
			}
		}

		// Calculate turn rate from bank angle (coordinated turn formula)
		if a.UniversalPlatform.State.Speed > 0 && math.Abs(a.UniversalPlatform.State.Physics.Attitude.Roll) > 1.0 {
			// Turn rate = (g * tan(bank_angle)) / velocity (in radians/second)
			bankRad := a.UniversalPlatform.State.Physics.Attitude.Roll * math.Pi / 180
			turnRate := (9.81 * math.Tan(bankRad)) / a.UniversalPlatform.State.Speed * 180 / math.Pi // Convert to degrees/second
			a.UniversalPlatform.State.Physics.AngularVelocity.YawRate = turnRate
			a.UniversalPlatform.State.Heading += turnRate * dt

			// Normalize heading to [0, 360)
			if a.UniversalPlatform.State.Heading >= 360 {
				a.UniversalPlatform.State.Heading -= 360
			} else if a.UniversalPlatform.State.Heading < 0 {
				a.UniversalPlatform.State.Heading += 360
			}
		}
	} else {
		// No destination - level off the aircraft
		rollError := -a.UniversalPlatform.State.Physics.Attitude.Roll
		maxRollChange := a.MaxRollRate * dt

		if math.Abs(rollError) <= maxRollChange {
			a.UniversalPlatform.State.Physics.Attitude.Roll = 0
		} else if rollError > 0 {
			a.UniversalPlatform.State.Physics.Attitude.Roll += maxRollChange
		} else {
			a.UniversalPlatform.State.Physics.Attitude.Roll -= maxRollChange
		}
	}

	// Calculate pitch angle based on climb/descent
	if a.UniversalPlatform.State.Speed > 0 {
		desiredPitch := math.Atan(a.UniversalPlatform.State.Velocity.Up/a.UniversalPlatform.State.Speed) * 180 / math.Pi
		desiredPitch = math.Max(-a.MaxPitchAngle, math.Min(a.MaxPitchAngle, desiredPitch))

		pitchError := desiredPitch - a.UniversalPlatform.State.Physics.Attitude.Pitch
		maxPitchChange := a.MaxPitchRate * dt

		if math.Abs(pitchError) <= maxPitchChange {
			a.UniversalPlatform.State.Physics.Attitude.Pitch = desiredPitch
		} else if pitchError > 0 {
			a.UniversalPlatform.State.Physics.Attitude.Pitch += maxPitchChange
		} else {
			a.UniversalPlatform.State.Physics.Attitude.Pitch -= maxPitchChange
		}
	}

	// Update yaw to match heading
	a.UniversalPlatform.State.Physics.Attitude.Yaw = a.UniversalPlatform.State.Heading
	// Update roll for compatibility
	a.UniversalPlatform.State.Roll = a.UniversalPlatform.State.Physics.Attitude.Roll
}

// updatePosition calculates new position based on current heading and speed
func (a *AirbornePlatform) updatePosition(deltaSeconds float64) {
	if a.UniversalPlatform.State.Speed <= 0 {
		return
	}

	// Convert heading to radians (0° = North, 90° = East)
	headingRad := a.UniversalPlatform.State.Heading * math.Pi / 180.0

	// Calculate distance moved
	distance := a.UniversalPlatform.State.Speed * deltaSeconds

	// Earth radius in meters
	earthRadius := 6371000.0

	// Calculate new position
	deltaLat := (distance * math.Cos(headingRad)) / earthRadius * 180.0 / math.Pi
	deltaLon := (distance * math.Sin(headingRad)) / earthRadius * 180.0 / math.Pi / math.Cos(a.UniversalPlatform.State.Position.Latitude*math.Pi/180.0)

	a.UniversalPlatform.State.Position.Latitude += deltaLat
	a.UniversalPlatform.State.Position.Longitude += deltaLon

	// Update altitude based on climb rate
	climbRate := a.UniversalPlatform.State.Velocity.Up
	a.UniversalPlatform.State.Position.Altitude += climbRate * deltaSeconds

	// Sync with physics state
	a.UniversalPlatform.State.Physics.Position = a.UniversalPlatform.State.Position
	a.UniversalPlatform.State.Physics.Velocity.North = a.UniversalPlatform.State.Speed * math.Cos(headingRad)
	a.UniversalPlatform.State.Physics.Velocity.East = a.UniversalPlatform.State.Speed * math.Sin(headingRad)
	a.UniversalPlatform.State.Physics.Velocity.Up = climbRate
}

// calculateBearing calculates bearing from one position to another
func (a *AirbornePlatform) calculateBearing(from, to Position) float64 {
	lat1 := from.Latitude * math.Pi / 180.0
	lat2 := to.Latitude * math.Pi / 180.0
	deltaLon := (to.Longitude - from.Longitude) * math.Pi / 180.0

	y := math.Sin(deltaLon) * math.Cos(lat2)
	x := math.Cos(lat1)*math.Sin(lat2) - math.Sin(lat1)*math.Cos(lat2)*math.Cos(deltaLon)

	bearing := math.Atan2(y, x) * 180.0 / math.Pi

	// Normalize to [0, 360)
	for bearing < 0 {
		bearing += 360
	}
	for bearing >= 360 {
		bearing -= 360
	}

	return bearing
}

// Aircraft factory functions with enhanced physics parameters

// NewBoeing737_800 creates a Boeing 737-800 with realistic physics
func NewBoeing737_800(id, flightNumber string, startPos Position) *AirbornePlatform {
	mass := 79010.0   // kg
	wingArea := 124.6 // m²

	// Create the base UniversalPlatform
	typeDef := &PlatformTypeDefinition{
		Class: "Boeing 737-800",
		Performance: PerformanceCharacteristics{
			MaxSpeed:      257,   // m/s (500 kts)
			CruiseSpeed:   230,   // m/s (447 kts)
			MaxAltitude:   12500, // meters (41,000 ft)
			ClimbRate:     15,    // m/s
			Acceleration:  1.5,   // m/s²
			TurningRadius: 1000,  // meters
		},
		Physical: PhysicalCharacteristics{
			Length: 39.5, // meters
			Width:  35.8, // wingspan
			Height: 12.5, // meters
			Mass:   mass,
		},
	}

	config := &PlatformConfiguration{
		ID:   id,
		Name: flightNumber,
		Type: "Boeing 737-800",
	}

	universalPlatform := UniversalPlatform{
		ID:           id,
		PlatformType: PlatformTypeAirborne,
		TypeDef:      typeDef,
		Config:       config,
		CallSign:     flightNumber,
		State: PlatformState{
			ID:          id,
			Position:    startPos,
			Velocity:    Velocity{},
			Heading:     0,
			Speed:       0,
			LastUpdated: time.Now(),
			Physics: PhysicsState{
				Position: startPos,
				Mass:     mass,
				MomentOfInertia: MomentOfInertia{
					Ixx: mass * 15.0 * 15.0,
					Iyy: mass * 20.0 * 20.0,
					Izz: mass * 20.0 * 20.0,
				},
			},
		},
	}

	return &AirbornePlatform{
		UniversalPlatform: universalPlatform,
		MaxRollRate:       15,  // degrees/second
		MaxPitchRate:      5,   // degrees/second
		MaxYawRate:        3,   // degrees/second
		MaxBankAngle:      30,  // degrees (normal ops)
		MaxPitchAngle:     15,  // degrees
		MaxLoadFactor:     2.5, // g-force
		StallSpeed:        77,  // m/s (150 kts)
		MaxAcceleration:   1.5, // m/s²
		MaxDeceleration:   3.0, // m/s²
		WingArea:          wingArea,
		WingLoading:       mass / wingArea,
		CenterOfGravity:   CenterOfGravity{X: 18.0, Y: 0, Z: 2.0},
		FlightPhase:       FlightPhaseParked,
	}
}

// NewAirbusA320 creates an Airbus A320 with realistic physics
func NewAirbusA320(id, flightNumber string, startPos Position) *AirbornePlatform {
	mass := 78000.0   // kg
	wingArea := 122.6 // m²

	// Create the base UniversalPlatform
	typeDef := &PlatformTypeDefinition{
		Class: "Airbus A320",
		Performance: PerformanceCharacteristics{
			MaxSpeed:      257,   // m/s (500 kts)
			CruiseSpeed:   230,   // m/s (447 kts)
			MaxAltitude:   12000, // meters (39,370 ft)
			ClimbRate:     15,    // m/s
			Acceleration:  1.5,   // m/s²
			TurningRadius: 1000,  // meters
		},
		Physical: PhysicalCharacteristics{
			Length: 37.6, // meters
			Width:  36.0, // wingspan
			Height: 11.8, // meters
			Mass:   mass,
		},
	}

	config := &PlatformConfiguration{
		ID:   id,
		Name: flightNumber,
		Type: "Airbus A320",
	}

	universalPlatform := UniversalPlatform{
		ID:           id,
		PlatformType: PlatformTypeAirborne,
		TypeDef:      typeDef,
		Config:       config,
		CallSign:     flightNumber,
		State: PlatformState{
			ID:          id,
			Position:    startPos,
			Velocity:    Velocity{},
			Heading:     0,
			Speed:       0,
			LastUpdated: time.Now(),
			Physics: PhysicsState{
				Position: startPos,
				Mass:     mass,
				MomentOfInertia: MomentOfInertia{
					Ixx: mass * 18.0 * 18.0,
					Iyy: mass * 19.0 * 19.0,
					Izz: mass * 19.0 * 19.0,
				},
			},
		},
	}

	return &AirbornePlatform{
		UniversalPlatform: universalPlatform,
		MaxRollRate:       15,  // degrees/second
		MaxPitchRate:      5,   // degrees/second
		MaxYawRate:        3,   // degrees/second
		MaxBankAngle:      30,  // degrees
		MaxPitchAngle:     15,  // degrees
		MaxLoadFactor:     2.5, // g-force
		StallSpeed:        77,  // m/s (150 kts)
		MaxAcceleration:   1.5, // m/s²
		MaxDeceleration:   3.0, // m/s²
		WingArea:          wingArea,
		WingLoading:       mass / wingArea,
		CenterOfGravity:   CenterOfGravity{X: 17.0, Y: 0, Z: 2.0},
		FlightPhase:       FlightPhaseParked,
	}
}

// NewF16FightingFalcon creates an F-16 with enhanced military flight dynamics
func NewF16FightingFalcon(id, tailNumber string, startPos Position) *AirbornePlatform {
	mass := 19187.0   // kg
	wingArea := 27.87 // m²

	// Create the base UniversalPlatform
	typeDef := &PlatformTypeDefinition{
		Class: "F-16 Fighting Falcon",
		Performance: PerformanceCharacteristics{
			MaxSpeed:      617,   // m/s (Mach 2.0+ at altitude)
			CruiseSpeed:   257,   // m/s (500 kts)
			MaxAltitude:   15240, // meters (50,000 ft)
			ClimbRate:     25,    // m/s
			Acceleration:  8.0,   // m/s²
			TurningRadius: 500,   // meters
		},
		Physical: PhysicalCharacteristics{
			Length: 15.0, // meters
			Width:  10.0, // wingspan
			Height: 5.1,  // meters
			Mass:   mass,
		},
	}

	config := &PlatformConfiguration{
		ID:   id,
		Name: tailNumber,
		Type: "F-16 Fighting Falcon",
	}

	callSign := fmt.Sprintf("VIPER%s", tailNumber[len(tailNumber)-3:])

	universalPlatform := UniversalPlatform{
		ID:           id,
		PlatformType: PlatformTypeAirborne,
		TypeDef:      typeDef,
		Config:       config,
		CallSign:     callSign,
		State: PlatformState{
			ID:          id,
			Position:    startPos,
			Velocity:    Velocity{},
			Heading:     0,
			Speed:       0,
			LastUpdated: time.Now(),
			Physics: PhysicsState{
				Position: startPos,
				Mass:     mass,
				MomentOfInertia: MomentOfInertia{
					Ixx: mass * 5.0 * 5.0,
					Iyy: mass * 7.5 * 7.5,
					Izz: mass * 7.5 * 7.5,
				},
			},
		},
	}

	return &AirbornePlatform{
		UniversalPlatform: universalPlatform,
		MaxRollRate:       720,  // degrees/second (highly maneuverable)
		MaxPitchRate:      40,   // degrees/second
		MaxYawRate:        20,   // degrees/second
		MaxBankAngle:      90,   // degrees (fighter aircraft)
		MaxPitchAngle:     60,   // degrees
		MaxLoadFactor:     9.0,  // g-force
		StallSpeed:        93,   // m/s (180 kts)
		MaxAcceleration:   8.0,  // m/s² (high thrust-to-weight)
		MaxDeceleration:   15.0, // m/s² (air brakes)
		WingArea:          wingArea,
		WingLoading:       mass / wingArea,
		CenterOfGravity:   CenterOfGravity{X: 7.5, Y: 0, Z: 1.0},
		FlightPhase:       FlightPhaseParked,
	}
}

// NewC130Hercules creates a C-130 military transport with cargo aircraft characteristics
func NewC130Hercules(id, tailNumber string, startPos Position) *AirbornePlatform {
	mass := 70300.0   // kg
	wingArea := 162.1 // m²

	// Create the base UniversalPlatform
	typeDef := &PlatformTypeDefinition{
		Class: "C-130 Hercules",
		Performance: PerformanceCharacteristics{
			MaxSpeed:      190,   // m/s (370 kts)
			CruiseSpeed:   160,   // m/s (310 kts)
			MaxAltitude:   10060, // meters (33,000 ft)
			ClimbRate:     8,     // m/s
			Acceleration:  1.0,   // m/s²
			TurningRadius: 1500,  // meters
		},
		Physical: PhysicalCharacteristics{
			Length: 29.8, // meters
			Width:  40.4, // wingspan
			Height: 11.7, // meters
			Mass:   mass,
		},
	}

	config := &PlatformConfiguration{
		ID:   id,
		Name: tailNumber,
		Type: "C-130 Hercules",
	}

	callSign := fmt.Sprintf("HERKY%s", tailNumber[len(tailNumber)-2:])

	universalPlatform := UniversalPlatform{
		ID:           id,
		PlatformType: PlatformTypeAirborne,
		TypeDef:      typeDef,
		Config:       config,
		CallSign:     callSign,
		State: PlatformState{
			ID:          id,
			Position:    startPos,
			Velocity:    Velocity{},
			Heading:     0,
			Speed:       0,
			LastUpdated: time.Now(),
			Physics: PhysicsState{
				Position: startPos,
				Mass:     mass,
				MomentOfInertia: MomentOfInertia{
					Ixx: mass * 20.0 * 20.0,
					Iyy: mass * 15.0 * 15.0,
					Izz: mass * 25.0 * 25.0,
				},
			},
		},
	}

	return &AirbornePlatform{
		UniversalPlatform: universalPlatform,
		MaxRollRate:       10,  // degrees/second (large aircraft)
		MaxPitchRate:      3,   // degrees/second
		MaxYawRate:        2,   // degrees/second
		MaxBankAngle:      25,  // degrees (transport aircraft)
		MaxPitchAngle:     10,  // degrees
		MaxLoadFactor:     2.5, // g-force
		StallSpeed:        61,  // m/s (119 kts)
		MaxAcceleration:   1.0, // m/s²
		MaxDeceleration:   2.5, // m/s²
		WingArea:          wingArea,
		WingLoading:       mass / wingArea,
		CenterOfGravity:   CenterOfGravity{X: 14.9, Y: 0, Z: 3.0},
		FlightPhase:       FlightPhaseParked,
	}
}
