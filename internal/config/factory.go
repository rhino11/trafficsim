package config

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/rhino11/trafficsim/internal/models"
)

// PlatformFactory creates platform instances from configuration data
type PlatformFactory struct {
	registry *PlatformRegistry
}

// NewPlatformFactory creates a new platform factory
func NewPlatformFactory(registry *PlatformRegistry) *PlatformFactory {
	return &PlatformFactory{
		registry: registry,
	}
}

// CreatePlatform creates a universal platform instance from configuration
func (f *PlatformFactory) CreatePlatform(instance PlatformInstance) (models.Platform, error) {
	// Get the platform type definition
	typeDef, err := f.registry.GetType(instance.TypeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get platform type %s: %w", instance.TypeID, err)
	}

	// Convert config position to models position
	startPos := models.Position{
		Latitude:  instance.StartPos.Latitude,
		Longitude: instance.StartPos.Longitude,
		Altitude:  instance.StartPos.Altitude,
	}

	// Generate call sign if not provided
	callSign := instance.CallSign
	if callSign == "" {
		callSign = f.generateCallSign(typeDef, instance.ID)
	}

	// Determine platform type from configuration
	platformType := f.determinePlatformType(typeDef.Type)

	// Convert PlatformTypeDefinition to models.PlatformTypeDefinition
	modelTypeDef := f.convertToModelTypeDefinition(typeDef)

	// Create platform configuration
	platformConfig := &models.PlatformConfiguration{
		ID:            instance.ID,
		Type:          typeDef.Type,
		Name:          instance.Name,
		StartPosition: startPos,
		Mission: models.MissionConfiguration{
			Type:       "standard",
			Parameters: make(map[string]interface{}),
		},
	}

	// Create universal platform
	platform := &models.UniversalPlatform{
		ID:           instance.ID,
		PlatformType: platformType,
		TypeDef:      modelTypeDef,
		Config:       platformConfig,
		CallSign:     callSign,
		State: models.PlatformState{
			ID:          instance.ID,
			Position:    startPos,
			Velocity:    models.Velocity{},
			Heading:     0,
			Speed:       0,
			LastUpdated: time.Now(),
		},
		FuelRemaining: modelTypeDef.Physical.FuelCapacity,
		MissionTime:   0,
	}

	return platform, nil
}

// CreateScenario creates all platforms for a given scenario - only loads specified platforms
func (f *PlatformFactory) CreateScenario(scenarioName string) ([]models.Platform, error) {
	scenario, exists := f.registry.Scenarios[scenarioName]
	if !exists {
		return nil, fmt.Errorf("scenario not found: %s", scenarioName)
	}

	var platforms []models.Platform

	// Only create platforms that are explicitly defined in the scenario
	for _, instance := range scenario.Instances {
		// Verify the platform type exists in configuration
		if !f.registry.HasType(instance.TypeID) {
			return nil, fmt.Errorf("platform type %s not found in configuration for instance %s",
				instance.TypeID, instance.ID)
		}

		platform, err := f.CreatePlatform(instance)
		if err != nil {
			return nil, fmt.Errorf("failed to create platform %s: %w", instance.ID, err)
		}

		// Set destination if specified
		if instance.Destination != nil {
			dest := models.Position{
				Latitude:  instance.Destination.Latitude,
				Longitude: instance.Destination.Longitude,
				Altitude:  instance.Destination.Altitude,
			}
			if err := platform.SetDestination(dest); err != nil {
				return nil, fmt.Errorf("failed to set destination for %s: %w", instance.ID, err)
			}
		}

		platforms = append(platforms, platform)
	}

	return platforms, nil
}

// GetAvailablePlatformTypes returns only the platform types that are actually configured
func (f *PlatformFactory) GetAvailablePlatformTypes() map[string][]string {
	available := make(map[string][]string)

	// Only include platform types that are defined in the configuration
	if len(f.registry.AirborneTypes) > 0 {
		var airborneTypes []string
		for typeID := range f.registry.AirborneTypes {
			airborneTypes = append(airborneTypes, typeID)
		}
		available["airborne"] = airborneTypes
	}

	if len(f.registry.MaritimeTypes) > 0 {
		var maritimeTypes []string
		for typeID := range f.registry.MaritimeTypes {
			maritimeTypes = append(maritimeTypes, typeID)
		}
		available["maritime"] = maritimeTypes
	}

	if len(f.registry.LandTypes) > 0 {
		var landTypes []string
		for typeID := range f.registry.LandTypes {
			landTypes = append(landTypes, typeID)
		}
		available["land"] = landTypes
	}

	if len(f.registry.SpaceTypes) > 0 {
		var spaceTypes []string
		for typeID := range f.registry.SpaceTypes {
			spaceTypes = append(spaceTypes, typeID)
		}
		available["space"] = spaceTypes
	}

	return available
}

// ValidateScenario validates that all platforms in a scenario reference valid types
func (f *PlatformFactory) ValidateScenario(scenarioName string) error {
	scenario, exists := f.registry.Scenarios[scenarioName]
	if !exists {
		return fmt.Errorf("scenario not found: %s", scenarioName)
	}

	for i, instance := range scenario.Instances {
		if !f.registry.HasType(instance.TypeID) {
			return fmt.Errorf("scenario %s, instance %d (%s): unknown platform type %s",
				scenarioName, i, instance.ID, instance.TypeID)
		}
	}

	return nil
}

// generateCallSign generates a call sign based on the platform type definition
func (f *PlatformFactory) generateCallSign(typeDef *PlatformTypeDefinition, instanceID string) string {
	if typeDef.CallSignPrefix != "" {
		if typeDef.CallSignFormat != "" {
			// Use custom format
			callSign := strings.ReplaceAll(typeDef.CallSignFormat, "{prefix}", typeDef.CallSignPrefix)
			callSign = strings.ReplaceAll(callSign, "{id}", instanceID)
			return callSign
		}
		// Default format: prefix + last 3 chars of ID
		suffix := instanceID
		if len(instanceID) > 3 {
			suffix = instanceID[len(instanceID)-3:]
		}
		return typeDef.CallSignPrefix + suffix
	}

	// Fallback to instance ID
	return instanceID
}

// determinePlatformType converts string type to PlatformType enum
func (f *PlatformFactory) determinePlatformType(typeStr string) models.PlatformType {
	switch typeStr {
	case "airborne":
		return models.PlatformTypeAirborne
	case "maritime":
		return models.PlatformTypeMaritime
	case "land":
		return models.PlatformTypeLand
	case "space":
		return models.PlatformTypeSpace
	default:
		return models.PlatformTypeAirborne // Default fallback
	}
}

// convertToModelTypeDefinition converts config type definition to models type definition
func (f *PlatformFactory) convertToModelTypeDefinition(configDef *PlatformTypeDefinition) *models.PlatformTypeDefinition {
	return &models.PlatformTypeDefinition{
		Class:    configDef.Class,
		Category: configDef.Category,
		Performance: models.PerformanceCharacteristics{
			MaxSpeed:        configDef.MaxSpeed,
			CruiseSpeed:     configDef.CruiseSpeed,
			MaxAltitude:     configDef.MaxAltitude,
			FuelConsumption: calculateFuelConsumption(configDef),
			TurningRadius:   calculateTurningRadius(configDef),
			Acceleration:    calculateAcceleration(configDef),
			MaxGradient:     configDef.MaxGradient,
			ClimbRate:       calculateClimbRate(configDef),

			// Orbital characteristics
			OrbitalVelocity: configDef.MaxSpeed, // Use max speed as orbital velocity for space platforms
			OrbitalPeriod:   configDef.OrbitalPeriod,
			OrbitalAltitude: configDef.MaxAltitude, // Use max altitude as orbital altitude
			Inclination:     configDef.Inclination,
			Eccentricity:    0.0, // Assume circular orbits for simplicity
		},
		Physical: models.PhysicalCharacteristics{
			Length:       configDef.Length,
			Width:        configDef.Width,
			Height:       configDef.Height,
			Mass:         configDef.Mass,
			FuelCapacity: configDef.FuelCapacity,
			Draft:        configDef.Draft,
		},
		Operational: models.OperationalCharacteristics{
			Range: configDef.Range,
		},
		CallsignConf: models.CallsignConfiguration{
			Prefix: configDef.CallSignPrefix,
			Format: configDef.CallSignFormat,
		},
	}
}

// Helper functions to calculate missing performance characteristics

func calculateFuelConsumption(def *PlatformTypeDefinition) float64 {
	// Estimate fuel consumption based on platform characteristics
	if def.FuelCapacity > 0 && def.Range > 0 {
		// Fuel consumption rate in liters per meter
		fuelRate := def.FuelCapacity / def.Range
		// Convert to liters per second at cruise speed
		return fuelRate * def.CruiseSpeed
	}
	return 0.1 // Default minimal fuel consumption
}

func calculateTurningRadius(def *PlatformTypeDefinition) float64 {
	// Estimate turning radius based on platform type and size
	switch def.Type {
	case "airborne":
		// Aircraft turning radius based on speed and banking
		if def.CruiseSpeed > 0 {
			// Assume 30 degree bank angle for commercial aircraft
			bankAngle := 30.0 * math.Pi / 180.0
			return (def.CruiseSpeed * def.CruiseSpeed) / (9.81 * math.Tan(bankAngle))
		}
		return def.Length * 10 // Fallback: 10x length
	case "maritime":
		// Ship turning radius
		return def.Length * 5 // Ships typically turn in 5-10 ship lengths
	case "land":
		// Vehicle turning radius
		return def.Length * 2 // Vehicles can turn tighter
	case "space":
		// Orbital mechanics - very large turning radius
		return def.MaxAltitude // Use orbital altitude as turning "radius"
	default:
		return def.Length * 3
	}
}

func calculateAcceleration(def *PlatformTypeDefinition) float64 {
	// Estimate acceleration based on platform type
	switch def.Type {
	case "airborne":
		if def.Category == "military" {
			return 5.0 // Military aircraft: higher acceleration
		}
		return 2.0 // Commercial aircraft: moderate acceleration
	case "maritime":
		return 0.5 // Ships: slow acceleration
	case "land":
		if def.Category == "military" {
			return 3.0 // Military vehicles: good acceleration
		}
		return 2.0 // Civilian vehicles: moderate acceleration
	case "space":
		return 0.01 // Space platforms: very low acceleration
	default:
		return 1.0
	}
}

func calculateClimbRate(def *PlatformTypeDefinition) float64 {
	// Estimate climb rate for aircraft
	if def.Type == "airborne" {
		if def.Category == "military" {
			return 20.0 // Military aircraft: higher climb rate
		}
		return 10.0 // Commercial aircraft: standard climb rate
	}
	return 0.0 // Non-aircraft don't climb
}
