package output

import (
	"strings"

	"github.com/rhino11/trafficsim/internal/models"
)

// PlatformToCoTState converts a platform interface to a CoT PlatformState
// This function maps platform data to the structure needed for CoT generation
func PlatformToCoTState(platform models.Platform) PlatformState {
	state := platform.GetState()

	// Extract CoT type and affiliation from platform category and type
	cotType, affiliation := determinePlatformCoTInfo(platform)

	return PlatformState{
		ID:          platform.GetID(),
		Callsign:    platform.GetCallSign(),
		Latitude:    state.Position.Latitude,
		Longitude:   state.Position.Longitude,
		Altitude:    state.Position.Altitude,
		Speed:       state.Speed,
		Course:      state.Heading,
		CoTType:     cotType,
		Affiliation: affiliation,
	}
}

// determinePlatformCoTInfo extracts CoT type and affiliation from platform data
func determinePlatformCoTInfo(platform models.Platform) (cotType, affiliation string) {
	platformType := strings.ToLower(string(platform.GetType()))
	class := strings.ToLower(platform.GetClass())

	// Try to get category from the platform if it's a UniversalPlatform
	category := "unknown"
	if up, ok := platform.(*models.UniversalPlatform); ok {
		category = strings.ToLower(up.TypeDef.Category)
	}

	// Determine affiliation based on category and class
	affiliation = "neutral" // Default
	if strings.Contains(category, "military") || strings.Contains(class, "military") ||
		strings.Contains(class, "f-16") || strings.Contains(class, "f-22") ||
		strings.Contains(class, "abrams") || strings.Contains(class, "arleigh burke") ||
		strings.Contains(class, "destroyer") || strings.Contains(class, "hmmwv") {
		affiliation = "friend"
	} else if strings.Contains(category, "commercial") || strings.Contains(category, "civilian") ||
		strings.Contains(class, "boeing") || strings.Contains(class, "airbus") ||
		strings.Contains(class, "container") || strings.Contains(class, "tesla") {
		affiliation = "neutral"
	}

	// Determine dimension
	var dimension string
	switch platformType {
	case "airborne":
		dimension = DimensionAir
	case "land":
		dimension = DimensionGround
	case "maritime":
		dimension = DimensionSea
	case "space":
		dimension = DimensionSpace
	default:
		dimension = DimensionGround
	}

	// Map specific platform categories to more specific CoT categories for better type generation
	mappedCategory := mapCategoryForCoT(category, class)

	// Generate CoT type using the existing function
	cotType = GenerateMILSTD2525Type(mappedCategory, affiliation, dimension)

	return cotType, affiliation
}

// mapCategoryForCoT maps platform categories/classes to CoT-specific categories
func mapCategoryForCoT(category, class string) string {
	classLower := strings.ToLower(class)
	categoryLower := strings.ToLower(category)

	// Aircraft types
	if strings.Contains(classLower, "f-16") || strings.Contains(classLower, "f-22") ||
		strings.Contains(classLower, "fighter") {
		return "fighter_aircraft"
	}
	if strings.Contains(classLower, "boeing") || strings.Contains(classLower, "airbus") ||
		strings.Contains(categoryLower, "commercial") {
		return "commercial_aircraft"
	}
	if strings.Contains(classLower, "c-130") || strings.Contains(classLower, "hercules") {
		return "transport_aircraft"
	}
	if strings.Contains(classLower, "mq-9") || strings.Contains(classLower, "reaper") ||
		strings.Contains(categoryLower, "drone") {
		return "unmanned_aircraft"
	}

	// Ground vehicle types
	if strings.Contains(classLower, "abrams") || strings.Contains(categoryLower, "tank") {
		return "main_battle_tank"
	}
	if strings.Contains(classLower, "hmmwv") || strings.Contains(categoryLower, "tactical") {
		return "tactical_vehicle"
	}
	if strings.Contains(classLower, "tesla") || strings.Contains(classLower, "truck") ||
		strings.Contains(categoryLower, "commercial") || strings.Contains(categoryLower, "civilian") {
		return "commercial_vehicle"
	}

	// Maritime types
	if strings.Contains(classLower, "destroyer") || strings.Contains(classLower, "arleigh burke") {
		return "destroyer"
	}
	if strings.Contains(classLower, "container") || strings.Contains(classLower, "cargo") {
		return "cargo_vessel"
	}
	if strings.Contains(classLower, "submarine") {
		return "submarine"
	}

	// Space types
	if strings.Contains(classLower, "satellite") || strings.Contains(classLower, "starlink") {
		return "satellite"
	}
	if strings.Contains(classLower, "station") || strings.Contains(classLower, "iss") {
		return "space_station"
	}
	if strings.Contains(classLower, "gps") {
		return "navigation_satellite"
	}

	// Fall back to the original category
	return categoryLower
}

// ConvertPlatformListToCoTStates converts a list of platforms to CoT states
func ConvertPlatformListToCoTStates(platforms []models.Platform) []PlatformState {
	states := make([]PlatformState, len(platforms))
	for i, platform := range platforms {
		states[i] = PlatformToCoTState(platform)
	}
	return states
}
