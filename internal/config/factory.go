package config

import (
	"fmt"
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

// CreatePlatform creates a platform instance from configuration
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

	// Create platform based on type
	switch typeDef.Type {
	case "airborne":
		return f.createAirbornePlatform(typeDef, instance, startPos, callSign)
	case "maritime":
		return f.createMaritimePlatform(typeDef, instance, startPos, callSign)
	case "land":
		return f.createLandPlatform(typeDef, instance, startPos, callSign)
	case "space":
		return f.createSpacePlatform(typeDef, instance, startPos, callSign)
	default:
		return nil, fmt.Errorf("unknown platform type: %s", typeDef.Type)
	}
}

// CreateScenario creates all platforms for a given scenario
func (f *PlatformFactory) CreateScenario(scenarioName string) ([]models.Platform, error) {
	scenario, exists := f.registry.Scenarios[scenarioName]
	if !exists {
		return nil, fmt.Errorf("scenario not found: %s", scenarioName)
	}

	var platforms []models.Platform
	for _, instance := range scenario.Instances {
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

	// Fallback to instance name or ID
	return instanceID
}

// createAirbornePlatform creates an airborne platform from configuration
func (f *PlatformFactory) createAirbornePlatform(typeDef *PlatformTypeDefinition, instance PlatformInstance, startPos models.Position, callSign string) (*models.AirbornePlatform, error) {
	return &models.AirbornePlatform{
		ID:       instance.ID,
		Class:    typeDef.Class,
		Name:     instance.Name,
		CallSign: callSign,
		State: models.PlatformState{
			ID:          instance.ID,
			Position:    startPos,
			Velocity:    models.Velocity{},
			Heading:     0,
			Speed:       0,
			LastUpdated: time.Now(),
		},
		MaxSpeed:    typeDef.MaxSpeed,
		MaxAltitude: typeDef.MaxAltitude,
		CruiseSpeed: typeDef.CruiseSpeed,
		CruiseAlt:   typeDef.ServiceCeiling * 0.9, // 90% of service ceiling as default cruise
		ServiceCeil: typeDef.ServiceCeiling,
		Length:      typeDef.Length,
		Width:       typeDef.Width,
		Height:      typeDef.Height,
		Mass:        typeDef.Mass,
	}, nil
}

// createMaritimePlatform creates a maritime platform from configuration
func (f *PlatformFactory) createMaritimePlatform(typeDef *PlatformTypeDefinition, instance PlatformInstance, startPos models.Position, callSign string) (*models.MaritimePlatform, error) {
	// Ensure maritime platforms start at sea level
	startPos.Altitude = 0

	return &models.MaritimePlatform{
		ID:       instance.ID,
		Class:    typeDef.Class,
		Name:     instance.Name,
		CallSign: callSign,
		State: models.PlatformState{
			ID:          instance.ID,
			Position:    startPos,
			Velocity:    models.Velocity{},
			Heading:     0,
			Speed:       0,
			LastUpdated: time.Now(),
		},
		MaxSpeed:     typeDef.MaxSpeed,
		CruiseSpeed:  typeDef.CruiseSpeed,
		Draft:        typeDef.Draft,
		Displacement: typeDef.Displacement,
		Length:       typeDef.Length,
		Width:        typeDef.Width,
		Height:       typeDef.Height,
		Mass:         typeDef.Mass,
	}, nil
}

// createLandPlatform creates a land platform from configuration
func (f *PlatformFactory) createLandPlatform(typeDef *PlatformTypeDefinition, instance PlatformInstance, startPos models.Position, callSign string) (*models.LandPlatform, error) {
	return &models.LandPlatform{
		ID:       instance.ID,
		Class:    typeDef.Class,
		Name:     instance.Name,
		CallSign: callSign,
		State: models.PlatformState{
			ID:          instance.ID,
			Position:    startPos,
			Velocity:    models.Velocity{},
			Heading:     0,
			Speed:       0,
			LastUpdated: time.Now(),
		},
		MaxSpeed:     typeDef.MaxSpeed,
		CruiseSpeed:  typeDef.CruiseSpeed,
		MaxGradient:  typeDef.MaxGradient,
		FuelCapacity: typeDef.FuelCapacity,
		Range:        typeDef.Range,
		Length:       typeDef.Length,
		Width:        typeDef.Width,
		Height:       typeDef.Height,
		Mass:         typeDef.Mass,
	}, nil
}

// createSpacePlatform creates a space platform from configuration
func (f *PlatformFactory) createSpacePlatform(typeDef *PlatformTypeDefinition, instance PlatformInstance, startPos models.Position, callSign string) (*models.SpacePlatform, error) {
	// Calculate orbital velocity based on altitude
	orbitalVelocity := typeDef.MaxSpeed
	if orbitalVelocity == 0 {
		// Calculate from orbital mechanics if not specified
		earthRadius := 6371000.0 // meters
		orbitalRadius := earthRadius + startPos.Altitude
		GM := 3.986e14                               // Earth's gravitational parameter
		orbitalVelocity = (GM / orbitalRadius) * 0.5 // Simplified calculation
	}

	return &models.SpacePlatform{
		ID:       instance.ID,
		Class:    typeDef.Class,
		Name:     instance.Name,
		CallSign: callSign,
		State: models.PlatformState{
			ID:          instance.ID,
			Position:    startPos,
			Velocity:    models.Velocity{},
			Heading:     90, // Eastward
			Speed:       orbitalVelocity,
			LastUpdated: time.Now(),
		},
		MaxSpeed:      orbitalVelocity,
		OrbitalPeriod: typeDef.OrbitalPeriod,
		Apogee:        typeDef.Apogee,
		Perigee:       typeDef.Perigee,
		Inclination:   typeDef.Inclination,
		Length:        typeDef.Length,
		Width:         typeDef.Width,
		Height:        typeDef.Height,
		Mass:          typeDef.Mass,
	}, nil
}
