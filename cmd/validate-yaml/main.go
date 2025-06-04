package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rhino11/trafficsim/internal/config"
	"github.com/rhino11/trafficsim/internal/models"
	"gopkg.in/yaml.v3"
)

// ValidationResult represents the result of validating a single file
type ValidationResult struct {
	File   string
	Valid  bool
	Errors []string
}

// ValidationSummary contains overall validation results
type ValidationSummary struct {
	TotalFiles   int
	ValidFiles   int
	InvalidFiles int
	Results      []ValidationResult
	FatalErrors  []string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <directory_or_file> [directory_or_file...]\n", os.Args[0])
		os.Exit(1)
	}

	var allFiles []string

	// Collect all YAML files from provided paths
	for _, path := range os.Args[1:] {
		files, err := collectYAMLFiles(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error collecting files from %s: %v\n", path, err)
			os.Exit(1)
		}
		allFiles = append(allFiles, files...)
	}

	if len(allFiles) == 0 {
		fmt.Println("No YAML files found to validate")
		return
	}

	summary := validateFiles(allFiles)
	printSummary(summary)

	if summary.InvalidFiles > 0 || len(summary.FatalErrors) > 0 {
		os.Exit(1)
	}
}

// collectYAMLFiles recursively finds all YAML files in the given path
func collectYAMLFiles(path string) ([]string, error) {
	var files []string

	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		if isYAMLFile(path) {
			files = append(files, path)
		}
		return files, nil
	}

	err = filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && isYAMLFile(filePath) {
			files = append(files, filePath)
		}

		return nil
	})

	return files, err
}

// isYAMLFile checks if a file has a YAML extension
func isYAMLFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".yaml" || ext == ".yml"
}

// validateFiles validates all provided YAML files
func validateFiles(files []string) ValidationSummary {
	summary := ValidationSummary{
		TotalFiles: len(files),
		Results:    make([]ValidationResult, 0, len(files)),
	}

	for _, file := range files {
		result := validateFile(file)
		summary.Results = append(summary.Results, result)

		if result.Valid {
			summary.ValidFiles++
		} else {
			summary.InvalidFiles++
		}
	}

	return summary
}

// validateFile validates a single YAML file based on its location and content
func validateFile(filePath string) ValidationResult {
	result := ValidationResult{
		File:   filePath,
		Valid:  true,
		Errors: []string{},
	}

	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to read file: %v", err))
		return result
	}

	// Basic YAML syntax validation
	var yamlData interface{}
	if err := yaml.Unmarshal(content, &yamlData); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, fmt.Sprintf("Invalid YAML syntax: %v", err))
		return result
	}

	// Determine file type and validate accordingly
	fileType := determineFileType(filePath)

	switch fileType {
	case "main_config":
		errors := validateMainConfig(content)
		result.Errors = append(result.Errors, errors...)
	case "platform_definition":
		errors := validatePlatformDefinition(content, filePath)
		result.Errors = append(result.Errors, errors...)
	case "scenario_config":
		errors := validateScenarioConfig(content)
		result.Errors = append(result.Errors, errors...)
	case "unknown":
		// For unknown files, just validate YAML syntax (already done above)
		break
	}

	if len(result.Errors) > 0 {
		result.Valid = false
	}

	return result
}

// determineFileType analyzes the file path to determine what type of config it is
func determineFileType(filePath string) string {
	// Normalize path separators
	normalizedPath := filepath.ToSlash(filePath)

	if strings.Contains(normalizedPath, "data/config.yaml") || strings.HasSuffix(normalizedPath, "config.yaml") {
		return "main_config"
	}

	if strings.Contains(normalizedPath, "data/platforms/") {
		return "platform_definition"
	}

	if strings.Contains(normalizedPath, "data/configs/") {
		return "scenario_config"
	}

	return "unknown"
}

// validateMainConfig validates the main configuration file structure
func validateMainConfig(content []byte) []string {
	var errors []string

	var cfg config.Config
	if err := yaml.Unmarshal(content, &cfg); err != nil {
		errors = append(errors, fmt.Sprintf("Failed to parse main config: %v", err))
		return errors
	}

	// Validate simulation config
	if cfg.Simulation.UpdateInterval == "" {
		errors = append(errors, "simulation.update_interval is required")
	}

	if cfg.Simulation.TimeScale <= 0 {
		errors = append(errors, "simulation.time_scale must be positive")
	}

	// Validate server config
	if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
		errors = append(errors, fmt.Sprintf("server.port must be between 1 and 65535, got %d", cfg.Server.Port))
	}

	// Validate platform types exist
	totalPlatformTypes := len(cfg.Platforms.AirborneTypes) +
		len(cfg.Platforms.MaritimeTypes) +
		len(cfg.Platforms.LandTypes) +
		len(cfg.Platforms.SpaceTypes)

	if totalPlatformTypes == 0 {
		errors = append(errors, "At least one platform type must be defined")
	}

	// Validate scenario references
	for scenarioName, scenario := range cfg.Platforms.Scenarios {
		for i, instance := range scenario.Instances {
			if !cfg.Platforms.HasType(instance.TypeID) {
				errors = append(errors, fmt.Sprintf("scenario %s, instance %d: references unknown platform type '%s'",
					scenarioName, i, instance.TypeID))
			}
		}
	}

	return errors
}

// validatePlatformDefinition validates platform definition files
func validatePlatformDefinition(content []byte, filePath string) []string {
	var errors []string

	// Try to parse as a platform type definition
	var platformData struct {
		PlatformTypes map[string]models.PlatformTypeDefinition `yaml:"platform_types"`
	}

	if err := yaml.Unmarshal(content, &platformData); err != nil {
		errors = append(errors, fmt.Sprintf("Failed to parse platform definition: %v", err))
		return errors
	}

	if len(platformData.PlatformTypes) == 0 {
		errors = append(errors, "No platform types defined in file")
		return errors
	}

	// Validate each platform type
	for typeName, platformType := range platformData.PlatformTypes {
		platformErrors := validatePlatformType(typeName, platformType, filePath)
		errors = append(errors, platformErrors...)
	}

	return errors
}

// validatePlatformType validates individual platform type definitions
func validatePlatformType(typeName string, platform models.PlatformTypeDefinition, filePath string) []string {
	var errors []string

	// Required fields validation
	if platform.Class == "" {
		errors = append(errors, fmt.Sprintf("platform '%s': class is required", typeName))
	}

	if platform.Category == "" {
		errors = append(errors, fmt.Sprintf("platform '%s': category is required", typeName))
	}

	// Performance validation
	if platform.Performance.MaxSpeed <= 0 {
		errors = append(errors, fmt.Sprintf("platform '%s': max_speed must be positive", typeName))
	}

	if platform.Performance.CruiseSpeed <= 0 {
		errors = append(errors, fmt.Sprintf("platform '%s': cruise_speed must be positive", typeName))
	}

	if platform.Performance.CruiseSpeed > platform.Performance.MaxSpeed {
		errors = append(errors, fmt.Sprintf("platform '%s': cruise_speed cannot exceed max_speed", typeName))
	}

	// Physical characteristics validation
	if platform.Physical.Length <= 0 {
		errors = append(errors, fmt.Sprintf("platform '%s': length must be positive", typeName))
	}

	if platform.Physical.Width <= 0 {
		errors = append(errors, fmt.Sprintf("platform '%s': width must be positive", typeName))
	}

	if platform.Physical.Mass <= 0 {
		errors = append(errors, fmt.Sprintf("platform '%s': mass must be positive", typeName))
	}

	// Domain-specific validation
	errors = append(errors, validateDomainSpecific(typeName, platform, filePath)...)

	return errors
}

// validateDomainSpecific performs domain-specific validation based on file location
func validateDomainSpecific(typeName string, platform models.PlatformTypeDefinition, filePath string) []string {
	var errors []string

	normalizedPath := filepath.ToSlash(filePath)

	if strings.Contains(normalizedPath, "airborne/") {
		// Airborne-specific validation
		if platform.Performance.MaxAltitude <= 0 {
			errors = append(errors, fmt.Sprintf("airborne platform '%s': max_altitude must be positive", typeName))
		}

		if platform.Performance.ClimbRate <= 0 {
			errors = append(errors, fmt.Sprintf("airborne platform '%s': climb_rate must be positive", typeName))
		}
	}

	if strings.Contains(normalizedPath, "maritime/") {
		// Maritime-specific validation
		if platform.Physical.Draft <= 0 {
			errors = append(errors, fmt.Sprintf("maritime platform '%s': draft must be positive", typeName))
		}
	}

	if strings.Contains(normalizedPath, "space/") {
		// Space-specific validation
		if platform.Performance.OrbitalVelocity <= 0 {
			errors = append(errors, fmt.Sprintf("space platform '%s': orbital_velocity must be positive", typeName))
		}

		if platform.Performance.OrbitalAltitude <= 0 {
			errors = append(errors, fmt.Sprintf("space platform '%s': orbital_altitude must be positive", typeName))
		}
	}

	return errors
}

// validateScenarioConfig validates scenario configuration files
func validateScenarioConfig(content []byte) []string {
	var errors []string

	var scenarioData struct {
		Metadata struct {
			Name        string `yaml:"name"`
			Description string `yaml:"description"`
			Duration    int    `yaml:"duration"`
		} `yaml:"metadata"`
		Platforms []struct {
			ID            string `yaml:"id"`
			Type          string `yaml:"type"`
			StartPosition struct {
				Latitude  float64 `yaml:"latitude"`
				Longitude float64 `yaml:"longitude"`
				Altitude  float64 `yaml:"altitude"`
			} `yaml:"start_position"`
		} `yaml:"platforms"`
	}

	if err := yaml.Unmarshal(content, &scenarioData); err != nil {
		errors = append(errors, fmt.Sprintf("Failed to parse scenario config: %v", err))
		return errors
	}

	// Validate metadata
	if scenarioData.Metadata.Name == "" {
		errors = append(errors, "metadata.name is required")
	}

	if scenarioData.Metadata.Duration <= 0 {
		errors = append(errors, "metadata.duration must be positive")
	}

	// Validate platforms
	platformIDs := make(map[string]bool)
	for i, platform := range scenarioData.Platforms {
		if platform.ID == "" {
			errors = append(errors, fmt.Sprintf("platform %d: id is required", i))
		} else {
			if platformIDs[platform.ID] {
				errors = append(errors, fmt.Sprintf("platform %d: duplicate id '%s'", i, platform.ID))
			}
			platformIDs[platform.ID] = true
		}

		if platform.Type == "" {
			errors = append(errors, fmt.Sprintf("platform %d (%s): type is required", i, platform.ID))
		}

		// Validate coordinates
		if platform.StartPosition.Latitude < -90 || platform.StartPosition.Latitude > 90 {
			errors = append(errors, fmt.Sprintf("platform %d (%s): latitude must be between -90 and 90", i, platform.ID))
		}

		if platform.StartPosition.Longitude < -180 || platform.StartPosition.Longitude > 180 {
			errors = append(errors, fmt.Sprintf("platform %d (%s): longitude must be between -180 and 180", i, platform.ID))
		}
	}

	return errors
}

// printSummary prints the validation results
func printSummary(summary ValidationSummary) {
	fmt.Printf("\n=== YAML Validation Summary ===\n")
	fmt.Printf("Total files: %d\n", summary.TotalFiles)
	fmt.Printf("Valid files: %d\n", summary.ValidFiles)
	fmt.Printf("Invalid files: %d\n", summary.InvalidFiles)

	if len(summary.FatalErrors) > 0 {
		fmt.Printf("\nFatal Errors:\n")
		for _, err := range summary.FatalErrors {
			fmt.Printf("  ❌ %s\n", err)
		}
	}

	if summary.InvalidFiles > 0 {
		fmt.Printf("\nValidation Errors:\n")
		for _, result := range summary.Results {
			if !result.Valid {
				fmt.Printf("\n📁 %s:\n", result.File)
				for _, err := range result.Errors {
					fmt.Printf("  ❌ %s\n", err)
				}
			}
		}
	}

	if summary.InvalidFiles == 0 && len(summary.FatalErrors) == 0 {
		fmt.Printf("\n✅ All YAML files are valid!\n")
	}
}
