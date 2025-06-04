package models

import (
	"fmt"
	"time"
)

// UnifiedPlatformFactory provides factory functions for creating UniversalPlatform instances
// with realistic configurations for different platform types

// Helper function to create base platform structure
func createBasePlatform(id string, platformType PlatformType, typeDef *PlatformTypeDefinition,
	config *PlatformConfiguration, startPos Position, callSign string, fuelRemaining float64) *UniversalPlatform {
	return &UniversalPlatform{
		ID:           id,
		PlatformType: platformType,
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
				Position:        startPos,
				Mass:            typeDef.Physical.Mass,
				MomentOfInertia: calculateMomentOfInertia(typeDef.Physical.Mass, platformType),
			},
		},
		CallSign:      callSign,
		FuelRemaining: fuelRemaining,
		MissionTime:   0,
		SystemStatus:  createSystemStatus(platformType, typeDef.Class),
		lastPosition:  startPos,
		acceleration:  0,
	}
}

// Helper function to calculate moment of inertia based on platform type
func calculateMomentOfInertia(mass float64, platformType PlatformType) MomentOfInertia {
	switch platformType {
	case PlatformTypeAirborne:
		return MomentOfInertia{
			Ixx: mass * 5.0 * 5.0,
			Iyy: mass * 7.5 * 7.5,
			Izz: mass * 7.5 * 7.5,
		}
	case PlatformTypeLand:
		return MomentOfInertia{
			Ixx: mass * 1.8 * 1.8,
			Iyy: mass * 4.9 * 4.9,
			Izz: mass * 4.9 * 4.9,
		}
	case PlatformTypeMaritime:
		return MomentOfInertia{
			Ixx: mass * 10.0 * 10.0,
			Iyy: mass * 77.5 * 77.5,
			Izz: mass * 77.5 * 77.5,
		}
	case PlatformTypeSpace:
		return MomentOfInertia{
			Ixx: mass * 2.5 * 2.5,
			Iyy: mass * 2.5 * 2.5,
			Izz: mass * 2.5 * 2.5,
		}
	default:
		return MomentOfInertia{
			Ixx: mass * 5.0 * 5.0,
			Iyy: mass * 5.0 * 5.0,
			Izz: mass * 5.0 * 5.0,
		}
	}
}

// Helper function to create system status based on platform type and class
func createSystemStatus(platformType PlatformType, class string) SystemStatus {
	// Determine if this is a civilian/commercial platform
	isCivilian := isCivilianPlatform(class)

	switch platformType {
	case PlatformTypeAirborne, PlatformTypeLand, PlatformTypeMaritime:
		weaponStatus := WeaponStatusArmed
		weaponSystemOperational := true
		weaponSystemEfficiency := 1.0

		// Civilian platforms should not be armed
		if isCivilian {
			weaponStatus = WeaponStatusNA
			weaponSystemOperational = false
			weaponSystemEfficiency = 0.0
		}

		return SystemStatus{
			PowerSystem:         SystemState{Operational: true, Efficiency: 1.0},
			PropulsionSystem:    SystemState{Operational: true, Efficiency: 0.98},
			NavigationSystem:    SystemState{Operational: true, Efficiency: 1.0},
			CommunicationSystem: SystemState{Operational: true, Efficiency: 0.99},
			SensorSystem:        SystemState{Operational: true, Efficiency: 0.95},
			WeaponSystem:        SystemState{Operational: weaponSystemOperational, Efficiency: weaponSystemEfficiency},
			FuelSystem:          SystemState{Operational: true, Efficiency: 1.0},
			WeaponStatus:        weaponStatus,
		}
	case PlatformTypeSpace:
		return SystemStatus{
			PowerSystem:         SystemState{Operational: true, Efficiency: 1.0},
			PropulsionSystem:    SystemState{Operational: true, Efficiency: 0.95},
			NavigationSystem:    SystemState{Operational: true, Efficiency: 1.0},
			CommunicationSystem: SystemState{Operational: true, Efficiency: 0.98},
			SensorSystem:        SystemState{Operational: true, Efficiency: 0.97},
			WeaponSystem:        SystemState{Operational: false, Efficiency: 0.0},
			FuelSystem:          SystemState{Operational: true, Efficiency: 1.0},
			WeaponStatus:        WeaponStatusNA,
		}
	default:
		return SystemStatus{
			PowerSystem:         SystemState{Operational: true, Efficiency: 1.0},
			PropulsionSystem:    SystemState{Operational: true, Efficiency: 1.0},
			NavigationSystem:    SystemState{Operational: true, Efficiency: 1.0},
			CommunicationSystem: SystemState{Operational: true, Efficiency: 1.0},
			SensorSystem:        SystemState{Operational: true, Efficiency: 1.0},
			WeaponSystem:        SystemState{Operational: false, Efficiency: 0.0},
			FuelSystem:          SystemState{Operational: true, Efficiency: 1.0},
			WeaponStatus:        WeaponStatusNA,
		}
	}
}

// Helper function to determine if a platform class represents a civilian platform
func isCivilianPlatform(class string) bool {
	civilianClasses := []string{
		"Civilian Car",
		"Commercial Aircraft",
		"Commercial Ship",
		"Cargo Aircraft",
		"Passenger Aircraft",
		"Commercial Truck",
		"Civilian Vehicle",
	}

	for _, civilianClass := range civilianClasses {
		if class == civilianClass {
			return true
		}
	}
	return false
}

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

	platform := createBasePlatform(id, PlatformTypeAirborne, typeDef, config, startPos, flightNumber, 26000)

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

	platform := createBasePlatform(id, PlatformTypeAirborne, typeDef, config, startPos, fmt.Sprintf("VIPER%s", tailNumber[len(tailNumber)-3:]), 3200)

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

	platform := createBasePlatform(id, PlatformTypeLand, typeDef, config, startPos, fmt.Sprintf("ARMOR%s", id[len(id)-2:]), 1900)

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

	platform := createBasePlatform(id, PlatformTypeMaritime, typeDef, config, startPos, fmt.Sprintf("NAVY%s", id[len(id)-3:]), 1200000)

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

	platform := createBasePlatform(id, PlatformTypeSpace, typeDef, config, startPos, fmt.Sprintf("STARLINK%s", id[len(id)-3:]), 50)

	return platform
}

// Additional factory functions for comprehensive coverage

// NewCivilianCarUniversal creates a civilian passenger car
func NewCivilianCarUniversal(id, model string, startPos Position) *UniversalPlatform {
	mass := 1500.0 // kg typical passenger car

	typeDef := &PlatformTypeDefinition{
		Class:    "Civilian Car", // Changed from model to "Civilian Car"
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

	platform := createBasePlatform(id, PlatformTypeLand, typeDef, config, startPos, fmt.Sprintf("CAR%s", id[len(id)-3:]), 60)

	return platform
}

// NewContainerShipUniversal creates a large container vessel
func NewContainerShipUniversal(id, shipName string, startPos Position) *UniversalPlatform {
	startPos.Altitude = 0 // Sea level
	mass := 200000000.0   // kg (200,000 tonnes)

	typeDef := &PlatformTypeDefinition{
		Class:    "Commercial Ship", // Changed from "Ultra Large Container Vessel" to "Commercial Ship"
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

	platform := createBasePlatform(id, PlatformTypeMaritime, typeDef, config, startPos, fmt.Sprintf("CARGO%s", id[len(id)-3:]), 15000000)

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
