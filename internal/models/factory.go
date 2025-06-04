package models

import (
	"fmt"
	"time"
)

// UnifiedPlatformFactory provides factory functions for creating UniversalPlatform instances
// with realistic configurations for different platform types

// Aircraft Factories

// NewBoeing737_800Universal creates a Boeing 737-800 using UniversalPlatform
func NewBoeing737_800Universal(id, flightNumber string, startPos Position) *UniversalPlatform {
	mass := 79010.0   // kg
	wingArea := 124.6 // m²

	typeDef := &PlatformTypeDefinition{
		Class:    "Boeing 737-800",
		Category: "commercial",
		Physical: PhysicalCharacteristics{
			Length:      39.5,
			Width:       35.8, // wingspan
			Height:      12.5,
			Mass:        mass,
			WingArea:    wingArea,
			FrontalArea: 15.0, // Estimated
		},
		Performance: PerformanceCharacteristics{
			MaxSpeed:          257,     // m/s (500 kts)
			CruiseSpeed:       230,     // m/s (447 kts)
			MaxAltitude:       12500,   // meters (41,000 ft)
			FuelConsumption:   3.5,     // kg/s fuel consumption rate
			TurningRadius:     3000,    // meters turning radius
			Acceleration:      1.5,     // m/s²
			ClimbRate:         12.0,    // m/s climb rate
			Range:             5665000, // meters range
			MaxThrustToWeight: 0.28,    // Typical for commercial aircraft
			MaxRollRate:       15,      // degrees/second
			MaxPitchRate:      5,       // degrees/second
			MaxYawRate:        3,       // degrees/second
			MaxBankAngle:      30,      // degrees (normal ops)
			MaxPitchAngle:     15,      // degrees
			MaxLoadFactor:     2.5,     // g-force
			StallSpeed:        77,      // m/s (150 kts)
			MaxAcceleration:   1.5,     // m/s²
			MaxDeceleration:   3.0,     // m/s²
		},
		Sensors: SensorCharacteristics{
			HasGPS:        true,
			HasRadar:      true,
			HasCompass:    true,
			RadarRange:    100000, // 100km weather radar
			OpticalRange:  50000,  // 50km visibility
			InfraredRange: 20000,  // 20km
		},
		Operational: OperationalCharacteristics{
			Range:             5665000, // meters (3,060 nm)
			PassengerCapacity: 189,
			CrewCapacity:      6,
		},
	}

	config := &PlatformConfiguration{
		ID:            id,
		Type:          "airborne",
		Name:          flightNumber,
		StartPosition: startPos,
	}

	platform := &UniversalPlatform{
		ID:           id,
		PlatformType: PlatformTypeAirborne,
		TypeDef:      typeDef,
		Config:       config,
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
					Ixx: mass * 15.0 * 15.0, // wingspan
					Iyy: mass * 20.0 * 20.0, // length
					Izz: mass * 20.0 * 20.0, // length
				},
			},
		},
		CallSign:      flightNumber,
		FuelRemaining: 26000,
		MissionTime:   0,
		SystemStatus: SystemStatus{
			PowerSystem:         SystemState{Operational: true, Efficiency: 1.0},
			PropulsionSystem:    SystemState{Operational: true, Efficiency: 0.95},
			NavigationSystem:    SystemState{Operational: true, Efficiency: 1.0},
			CommunicationSystem: SystemState{Operational: true, Efficiency: 0.98},
			SensorSystem:        SystemState{Operational: true, Efficiency: 0.99},
			FuelSystem:          SystemState{Operational: true, Efficiency: 1.0},
			WeaponStatus:        "N/A", // Civilian aircraft
		},
		lastPosition: startPos,
		acceleration: 0,
	}

	return platform
}

// NewF16FightingFalconUniversal creates an F-16 using UniversalPlatform
func NewF16FightingFalconUniversal(id, tailNumber string, startPos Position) *UniversalPlatform {
	mass := 19187.0   // kg
	wingArea := 27.87 // m²

	typeDef := &PlatformTypeDefinition{
		Class:    "F-16 Fighting Falcon",
		Category: "military",
		Physical: PhysicalCharacteristics{
			Length:       15.0,
			Width:        10.0, // wingspan
			Height:       5.1,
			Mass:         mass,
			WingArea:     wingArea,
			FrontalArea:  8.0,  // Estimated
			FuelCapacity: 3200, // liters
		},
		Performance: PerformanceCharacteristics{
			MaxSpeed:          617,   // m/s (Mach 2.0+ at altitude)
			CruiseSpeed:       257,   // m/s (500 kts)
			MaxAltitude:       15240, // meters (50,000 ft)
			MaxThrustToWeight: 1.2,   // High performance fighter
			MaxRollRate:       720,   // degrees/second
			MaxPitchRate:      40,    // degrees/second
			MaxYawRate:        20,    // degrees/second
			MaxBankAngle:      90,    // degrees
			MaxPitchAngle:     60,    // degrees
			MaxLoadFactor:     9.0,   // g-force
			StallSpeed:        93,    // m/s (180 kts)
			MaxAcceleration:   8.0,   // m/s²
			MaxDeceleration:   15.0,  // m/s²
		},
		Sensors: SensorCharacteristics{
			HasGPS:        true,
			HasRadar:      true,
			HasCompass:    true,
			RadarRange:    200000, // 200km air-to-air radar
			OpticalRange:  50000,  // 50km
			InfraredRange: 100000, // 100km IRST
		},
		Operational: OperationalCharacteristics{
			Range:         2220000, // meters (1,200 nm)
			CrewCapacity:  1,
			WeaponSystems: []string{"AIM-120", "AIM-9", "20mm Cannon"},
		},
	}

	config := &PlatformConfiguration{
		ID:            id,
		Type:          "airborne",
		Name:          tailNumber,
		StartPosition: startPos,
	}

	platform := &UniversalPlatform{
		ID:           id,
		PlatformType: PlatformTypeAirborne,
		TypeDef:      typeDef,
		Config:       config,
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
		CallSign:      fmt.Sprintf("VIPER%s", tailNumber[len(tailNumber)-3:]),
		FuelRemaining: 3200,
		MissionTime:   0,
		SystemStatus: SystemStatus{
			PowerSystem:         SystemState{Operational: true, Efficiency: 1.0},
			PropulsionSystem:    SystemState{Operational: true, Efficiency: 0.98},
			NavigationSystem:    SystemState{Operational: true, Efficiency: 1.0},
			CommunicationSystem: SystemState{Operational: true, Efficiency: 0.99},
			SensorSystem:        SystemState{Operational: true, Efficiency: 0.95},
			WeaponSystem:        SystemState{Operational: true, Efficiency: 1.0},
			FuelSystem:          SystemState{Operational: true, Efficiency: 1.0},
			WeaponStatus:        "ARMED", // Military fighter
		},
		lastPosition: startPos,
		acceleration: 0,
	}

	return platform
}

// Land Vehicle Factories

// NewM1A2AbramsUniversal creates an M1A2 Abrams tank using UniversalPlatform
func NewM1A2AbramsUniversal(id, unitDesignation string, startPos Position) *UniversalPlatform {
	mass := 62000.0 // kg

	typeDef := &PlatformTypeDefinition{
		Class:    "M1A2 Abrams MBT",
		Category: "military",
		Physical: PhysicalCharacteristics{
			Length:       9.8,
			Width:        3.7,
			Height:       2.4,
			Mass:         mass,
			FrontalArea:  8.9,  // width × height
			FuelCapacity: 1900, // liters
		},
		Performance: PerformanceCharacteristics{
			MaxSpeed:        20,   // m/s (45 mph)
			CruiseSpeed:     13.4, // m/s (30 mph)
			MaxAltitude:     0,    // Ground vehicle
			MaxAcceleration: 2.0,  // m/s²
			MaxDeceleration: 8.0,  // m/s²
			MaxGradient:     30.0, // degrees
		},
		Sensors: SensorCharacteristics{
			HasGPS:        true,
			HasRadar:      false, // Ground vehicles don't typically have radar
			HasCompass:    true,
			RadarRange:    0,
			OpticalRange:  5000,  // 5km visual
			InfraredRange: 10000, // 10km thermal imaging
		},
		Operational: OperationalCharacteristics{
			Range:         426000, // meters (265 miles)
			CrewCapacity:  4,
			WeaponSystems: []string{"120mm M256", "M240 Machine Gun"},
		},
	}

	config := &PlatformConfiguration{
		ID:            id,
		Type:          "land",
		Name:          unitDesignation,
		StartPosition: startPos,
	}

	platform := &UniversalPlatform{
		ID:           id,
		PlatformType: PlatformTypeLand,
		TypeDef:      typeDef,
		Config:       config,
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
					Ixx: mass * 1.8 * 1.8, // width
					Iyy: mass * 4.9 * 4.9, // length
					Izz: mass * 4.9 * 4.9, // length
				},
			},
		},
		CallSign:      fmt.Sprintf("ARMOR%s", id[len(id)-2:]),
		FuelRemaining: 1900,
		MissionTime:   0,
		SystemStatus: SystemStatus{
			PowerSystem:         SystemState{Operational: true, Efficiency: 1.0},
			PropulsionSystem:    SystemState{Operational: true, Efficiency: 0.90},
			NavigationSystem:    SystemState{Operational: true, Efficiency: 1.0},
			CommunicationSystem: SystemState{Operational: true, Efficiency: 0.95},
			SensorSystem:        SystemState{Operational: true, Efficiency: 0.92},
			WeaponSystem:        SystemState{Operational: true, Efficiency: 1.0},
			FuelSystem:          SystemState{Operational: true, Efficiency: 1.0},
			WeaponStatus:        "ARMED", // Military tank
		},
		lastPosition: startPos,
		acceleration: 0,
	}

	return platform
}

// Maritime Factories

// NewArleighBurkeDestroyerUniversal creates an Arleigh Burke destroyer using UniversalPlatform
func NewArleighBurkeDestroyerUniversal(id, shipName string, startPos Position) *UniversalPlatform {
	startPos.Altitude = 0 // Sea level
	mass := 9200000.0     // kg

	typeDef := &PlatformTypeDefinition{
		Class:    "Arleigh Burke-class",
		Category: "naval",
		Physical: PhysicalCharacteristics{
			Length:       155,
			Width:        20, // beam
			Height:       18, // above waterline
			Mass:         mass,
			WetArea:      3100,    // underwater hull area
			Draft:        6.3,     // depth below waterline
			FuelCapacity: 1200000, // liters
		},
		Performance: PerformanceCharacteristics{
			MaxSpeed:        15.4, // m/s (30+ knots)
			CruiseSpeed:     10.3, // m/s (20 knots)
			MaxAcceleration: 0.5,  // m/s²
			MaxDeceleration: 2.0,  // m/s²
		},
		Sensors: SensorCharacteristics{
			HasGPS:        true,
			HasRadar:      true,
			HasCompass:    true,
			RadarRange:    300000, // 300km long-range naval radar
			OpticalRange:  50000,  // 50km
			InfraredRange: 100000, // 100km
		},
		Operational: OperationalCharacteristics{
			Range:         8000000, // meters (4,320 nm)
			CrewCapacity:  323,
			WeaponSystems: []string{"Tomahawk", "Standard Missile", "5\" Gun"},
		},
	}

	config := &PlatformConfiguration{
		ID:            id,
		Type:          "maritime",
		Name:          fmt.Sprintf("USS %s", shipName),
		StartPosition: startPos,
	}

	platform := &UniversalPlatform{
		ID:           id,
		PlatformType: PlatformTypeMaritime,
		TypeDef:      typeDef,
		Config:       config,
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
					Ixx: mass * 10.0 * 10.0, // beam
					Iyy: mass * 77.5 * 77.5, // length
					Izz: mass * 77.5 * 77.5, // length
				},
			},
		},
		CallSign:      fmt.Sprintf("NAVY%s", id[len(id)-3:]),
		FuelRemaining: 1200000,
		MissionTime:   0,
		SystemStatus: SystemStatus{
			PowerSystem:         SystemState{Operational: true, Efficiency: 1.0},
			PropulsionSystem:    SystemState{Operational: true, Efficiency: 0.95},
			NavigationSystem:    SystemState{Operational: true, Efficiency: 1.0},
			CommunicationSystem: SystemState{Operational: true, Efficiency: 0.98},
			SensorSystem:        SystemState{Operational: true, Efficiency: 0.96},
			WeaponSystem:        SystemState{Operational: true, Efficiency: 1.0},
			FuelSystem:          SystemState{Operational: true, Efficiency: 1.0},
			WeaponStatus:        "ARMED", // Naval destroyer
		},
		lastPosition: startPos,
		acceleration: 0,
	}

	return platform
}

// Space Platform Factories

// NewStarlinkSatelliteUniversal creates a Starlink satellite using UniversalPlatform
func NewStarlinkSatelliteUniversal(id, satelliteNumber string, startPos Position) *UniversalPlatform {
	startPos.Altitude = 550000 // 550 km altitude
	mass := 260.0              // kg

	typeDef := &PlatformTypeDefinition{
		Class:    "Starlink Satellite",
		Category: "communications",
		Physical: PhysicalCharacteristics{
			Length:         2.8,
			Width:          1.9,
			Height:         0.32,
			Mass:           mass,
			SolarPanelArea: 32.0, // m²
			FuelCapacity:   50,   // kg xenon propellant
		},
		Performance: PerformanceCharacteristics{
			MaxSpeed:        7590,   // m/s orbital velocity
			CruiseSpeed:     7590,   // same as max for satellites
			MaxAltitude:     550000, // meters
			OrbitalPeriod:   5760,   // seconds (96 minutes)
			OrbitalVelocity: 7590,   // m/s
			OrbitalAltitude: 550000, // meters
			Inclination:     53.0,   // degrees
			MaxAcceleration: 0.001,  // m/s² (ion thrusters)
			MaxDeceleration: 0.001,  // m/s²
		},
		Sensors: SensorCharacteristics{
			HasGPS:        false, // Satellites don't use GPS
			HasRadar:      false,
			HasCompass:    false, // Not useful in space
			RadarRange:    0,
			OpticalRange:  1000000, // 1000km space observation
			InfraredRange: 500000,  // 500km
		},
		Operational: OperationalCharacteristics{
			MissionLife:    5.0, // years
			FrequencyBands: []string{"Ku-band", "Ka-band"},
		},
	}

	config := &PlatformConfiguration{
		ID:            id,
		Type:          "space",
		Name:          fmt.Sprintf("Starlink-%s", satelliteNumber),
		StartPosition: startPos,
	}

	platform := &UniversalPlatform{
		ID:           id,
		PlatformType: PlatformTypeSpace,
		TypeDef:      typeDef,
		Config:       config,
		State: PlatformState{
			ID:          id,
			Position:    startPos,
			Velocity:    Velocity{},
			Heading:     90, // Eastward orbital motion
			Speed:       7590,
			LastUpdated: time.Now(),
			Physics: PhysicsState{
				Position: startPos,
				Mass:     mass,
				MomentOfInertia: MomentOfInertia{
					Ixx: mass * 0.95 * 0.95, // width
					Iyy: mass * 1.4 * 1.4,   // length
					Izz: mass * 1.4 * 1.4,   // length
				},
			},
		},
		CallSign:      fmt.Sprintf("STARLINK%s", id[len(id)-3:]),
		FuelRemaining: 50,
		MissionTime:   0,
		SystemStatus: SystemStatus{
			PowerSystem:         SystemState{Operational: true, Efficiency: 1.0},
			PropulsionSystem:    SystemState{Operational: true, Efficiency: 0.98},
			NavigationSystem:    SystemState{Operational: true, Efficiency: 1.0},
			CommunicationSystem: SystemState{Operational: true, Efficiency: 1.0},
			SensorSystem:        SystemState{Operational: true, Efficiency: 0.99},
			FuelSystem:          SystemState{Operational: true, Efficiency: 1.0},
			WeaponStatus:        "N/A", // Civilian satellite
		},
		lastPosition: startPos,
		acceleration: 0,
	}

	return platform
}

// Additional factory functions for comprehensive coverage

// NewCivilianCarUniversal creates a civilian passenger car
func NewCivilianCarUniversal(id, model string, startPos Position) *UniversalPlatform {
	mass := 1500.0 // kg typical passenger car

	typeDef := &PlatformTypeDefinition{
		Class:    model,
		Category: "civilian",
		Physical: PhysicalCharacteristics{
			Length:       4.5,
			Width:        1.8,
			Height:       1.4,
			Mass:         mass,
			FrontalArea:  2.5, // typical car
			FuelCapacity: 60,  // liters
		},
		Performance: PerformanceCharacteristics{
			MaxSpeed:        50,  // m/s (112 mph)
			CruiseSpeed:     25,  // m/s (56 mph)
			MaxAcceleration: 3.0, // m/s²
			MaxDeceleration: 8.0, // m/s²
		},
		Sensors: SensorCharacteristics{
			HasGPS:        true,
			HasRadar:      false, // Basic civilian car
			HasCompass:    true,
			RadarRange:    0,
			OpticalRange:  1000, // 1km visual
			InfraredRange: 0,
		},
		Operational: OperationalCharacteristics{
			Range:             600000, // meters (600 km)
			PassengerCapacity: 5,
		},
	}

	config := &PlatformConfiguration{
		ID:            id,
		Type:          "land",
		Name:          fmt.Sprintf("Vehicle %s", id),
		StartPosition: startPos,
	}

	platform := &UniversalPlatform{
		ID:           id,
		PlatformType: PlatformTypeLand,
		TypeDef:      typeDef,
		Config:       config,
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
			},
		},
		CallSign:      fmt.Sprintf("CAR%s", id[len(id)-3:]),
		FuelRemaining: 60,
		MissionTime:   0,
		SystemStatus: SystemStatus{
			PowerSystem:         SystemState{Operational: true, Efficiency: 1.0},
			PropulsionSystem:    SystemState{Operational: true, Efficiency: 0.95},
			NavigationSystem:    SystemState{Operational: true, Efficiency: 1.0},
			CommunicationSystem: SystemState{Operational: true, Efficiency: 0.98},
			SensorSystem:        SystemState{Operational: true, Efficiency: 0.99},
			FuelSystem:          SystemState{Operational: true, Efficiency: 1.0},
			WeaponStatus:        "N/A", // Civilian vehicle
		},
		lastPosition: startPos,
		acceleration: 0,
	}

	return platform
}

// NewContainerShipUniversal creates a large container vessel
func NewContainerShipUniversal(id, shipName string, startPos Position) *UniversalPlatform {
	startPos.Altitude = 0 // Sea level
	mass := 200000000.0   // kg (200,000 tonnes)

	typeDef := &PlatformTypeDefinition{
		Class:    "Ultra Large Container Vessel",
		Category: "commercial",
		Physical: PhysicalCharacteristics{
			Length:       400,
			Width:        59, // beam
			Height:       73, // above waterline
			Mass:         mass,
			WetArea:      24000,    // large underwater hull
			Draft:        16,       // deep draft
			FuelCapacity: 15000000, // liters (massive fuel capacity)
		},
		Performance: PerformanceCharacteristics{
			MaxSpeed:        12.9,     // m/s (25 knots)
			CruiseSpeed:     10.3,     // m/s (20 knots)
			MaxAcceleration: 0.1,      // m/s² (very slow acceleration)
			MaxDeceleration: 0.5,      // m/s²
			Range:           24000000, // meters (global range) - FIXED: Added missing Range field
		},
		Sensors: SensorCharacteristics{
			HasGPS:        true,
			HasRadar:      true,
			HasCompass:    true,
			RadarRange:    50000, // 50km navigation radar
			OpticalRange:  20000, // 20km
			InfraredRange: 10000, // 10km
		},
		Operational: OperationalCharacteristics{
			Range:         24000000, // meters (global range)
			CrewCapacity:  25,
			CargoCapacity: 24000, // containers
		},
	}

	config := &PlatformConfiguration{
		ID:            id,
		Type:          "maritime",
		Name:          shipName,
		StartPosition: startPos,
	}

	platform := &UniversalPlatform{
		ID:           id,
		PlatformType: PlatformTypeMaritime,
		TypeDef:      typeDef,
		Config:       config,
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
			},
		},
		CallSign:      fmt.Sprintf("CARGO%s", id[len(id)-3:]),
		FuelRemaining: 15000000,
		MissionTime:   0,
		SystemStatus: SystemStatus{
			PowerSystem:         SystemState{Operational: true, Efficiency: 1.0},
			PropulsionSystem:    SystemState{Operational: true, Efficiency: 0.92},
			NavigationSystem:    SystemState{Operational: true, Efficiency: 1.0},
			CommunicationSystem: SystemState{Operational: true, Efficiency: 0.96},
			SensorSystem:        SystemState{Operational: true, Efficiency: 0.98},
			FuelSystem:          SystemState{Operational: true, Efficiency: 1.0},
			WeaponStatus:        "N/A", // Commercial vessel
		},
		lastPosition: startPos,
		acceleration: 0,
	}

	return platform
}

// Utility functions for platform creation

// CreatePlatformFromConfig creates a UniversalPlatform from configuration data
func CreatePlatformFromConfig(config map[string]interface{}) (*UniversalPlatform, error) {
	// This would parse configuration and create appropriate platform
	// Implementation would depend on configuration format
	// For now, return a basic platform

	id, ok := config["id"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid platform ID")
	}

	platformType, ok := config["type"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid platform type")
	}

	var pos Position
	if posData, ok := config["position"].(map[string]interface{}); ok {
		if lat, ok := posData["latitude"].(float64); ok {
			pos.Latitude = lat
		}
		if lon, ok := posData["longitude"].(float64); ok {
			pos.Longitude = lon
		}
		if alt, ok := posData["altitude"].(float64); ok {
			pos.Altitude = alt
		}
	}

	// Create basic platform based on type
	switch platformType {
	case "airborne":
		return NewBoeing737_800Universal(id, id, pos), nil
	case "land":
		return NewCivilianCarUniversal(id, "Generic Car", pos), nil
	case "maritime":
		return NewContainerShipUniversal(id, "Generic Ship", pos), nil
	case "space":
		return NewStarlinkSatelliteUniversal(id, id, pos), nil
	default:
		return nil, fmt.Errorf("unknown platform type: %s", platformType)
	}
}
