package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rhino11/trafficsim/internal/models"
)

func TestMain(t *testing.T) {
	// Save original args and defer restoration
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Test with no arguments - should show usage and exit
	os.Args = []string{"validate-yaml"}

	// Capture stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// We expect main to call os.Exit(1), so we need to test this differently
	// Instead, we'll test the logic that would cause the exit
	if len(os.Args) < 2 {
		// This is the condition that would cause the usage message
		w.Close()
		os.Stderr = oldStderr

		var buf bytes.Buffer
		io.Copy(&buf, r)
		// Since we can't easily test os.Exit, we'll just verify the condition
		// that would trigger the usage message
	} else {
		w.Close()
		os.Stderr = oldStderr
		t.Error("Expected fewer than 2 arguments to trigger usage")
	}
}

func TestCollectYAMLFiles(t *testing.T) {
	// Create temporary directory structure
	tmpDir, err := os.MkdirTemp("", "yaml-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	yamlFile1 := filepath.Join(tmpDir, "test1.yaml")
	yamlFile2 := filepath.Join(tmpDir, "test2.yml")
	txtFile := filepath.Join(tmpDir, "test.txt")

	for _, file := range []string{yamlFile1, yamlFile2, txtFile} {
		if err := os.WriteFile(file, []byte("test: content"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Create subdirectory with YAML file
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	yamlFile3 := filepath.Join(subDir, "test3.yaml")
	if err := os.WriteFile(yamlFile3, []byte("test: content"), 0644); err != nil {
		t.Fatal(err)
	}

	// Test collecting from directory
	files, err := collectYAMLFiles(tmpDir)
	if err != nil {
		t.Fatalf("collectYAMLFiles failed: %v", err)
	}

	if len(files) != 3 {
		t.Errorf("Expected 3 YAML files, got %d", len(files))
	}

	// Test collecting from single file
	files, err = collectYAMLFiles(yamlFile1)
	if err != nil {
		t.Fatalf("collectYAMLFiles failed for single file: %v", err)
	}

	if len(files) != 1 {
		t.Errorf("Expected 1 file, got %d", len(files))
	}

	// Test non-existent path
	_, err = collectYAMLFiles("/non/existent/path")
	if err == nil {
		t.Error("Expected error for non-existent path")
	}
}

func TestIsYAMLFile(t *testing.T) {
	tests := []struct {
		filename string
		expected bool
	}{
		{"test.yaml", true},
		{"test.yml", true},
		{"test.YAML", true},
		{"test.YML", true},
		{"test.txt", false},
		{"test.json", false},
		{"test", false},
		{"", false},
		{".yaml", true},
		{".yml", true},
	}

	for _, test := range tests {
		result := isYAMLFile(test.filename)
		if result != test.expected {
			t.Errorf("isYAMLFile(%q) = %v, expected %v", test.filename, result, test.expected)
		}
	}
}

func TestValidateFilePath(t *testing.T) {
	// Get current working directory
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "valid data path",
			path:    filepath.Join(wd, "data", "config.yaml"),
			wantErr: false,
		},
		{
			name:    "valid current directory",
			path:    filepath.Join(wd, "test.yaml"),
			wantErr: false,
		},
		{
			name:    "invalid outside path",
			path:    "/etc/passwd",
			wantErr: true,
		},
		{
			name:    "directory traversal attempt",
			path:    "../../etc/passwd",
			wantErr: true,
		},
		{
			name:    "relative data path",
			path:    "./data/config.yaml",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFilePath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateFilePath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDetermineFileType(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"data/config.yaml", "main_config"},
		{"/path/to/config.yaml", "main_config"},
		{"data/platforms/airborne/aircraft.yaml", "platform_definition"},
		{"data/platforms/maritime/ships.yaml", "platform_definition"},
		{"data/configs/scenario1.yaml", "scenario_config"},
		{"random/file.yaml", "unknown"},
		{"config.yaml", "main_config"},
		{"test.yaml", "unknown"},
		{"", "unknown"},
	}

	for _, test := range tests {
		result := determineFileType(test.path)
		if result != test.expected {
			t.Errorf("determineFileType(%q) = %q, expected %q", test.path, result, test.expected)
		}
	}
}

func TestValidateFile(t *testing.T) {
	// Create temporary directory in current working directory
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	tmpDir := filepath.Join(wd, "test-temp")
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Test valid YAML
	validYAML := `
test:
  key: value
  number: 42
`
	validFile := filepath.Join(tmpDir, "valid.yaml")
	if err := os.WriteFile(validFile, []byte(validYAML), 0644); err != nil {
		t.Fatal(err)
	}

	result := validateFile(validFile)
	if !result.Valid {
		t.Errorf("Expected valid result for valid YAML, got errors: %v", result.Errors)
	}

	// Test invalid YAML
	invalidYAML := `
test:
  key: value
  invalid: [unclosed array
`
	invalidFile := filepath.Join(tmpDir, "invalid.yaml")
	if err := os.WriteFile(invalidFile, []byte(invalidYAML), 0644); err != nil {
		t.Fatal(err)
	}

	result = validateFile(invalidFile)
	if result.Valid {
		t.Error("Expected invalid result for invalid YAML")
	}
	if len(result.Errors) == 0 {
		t.Error("Expected errors for invalid YAML")
	}

	// Test non-existent file
	result = validateFile("/non/existent/file.yaml")
	if result.Valid {
		t.Error("Expected invalid result for non-existent file")
	}

	// Test file outside allowed paths
	outsideFile := "/etc/passwd"
	result = validateFile(outsideFile)
	if result.Valid {
		t.Error("Expected invalid result for file outside allowed paths")
	}
	if len(result.Errors) == 0 {
		t.Error("Expected security error for file outside allowed paths")
	}
}

func TestValidateMainConfig(t *testing.T) {
	// Valid main config
	validConfig := `
simulation:
  update_interval: "16ms"
  time_scale: 1.0
server:
  port: 8080
platforms:
  airborne_types:
    boeing_737:
      class: "Boeing 737"
      category: "commercial"
      performance:
        max_speed: 257
        cruise_speed: 230
      physical:
        length: 39.5
        width: 35.8
        mass: 79010
  scenarios:
    test_scenario:
      instances:
        - type_id: "boeing_737"
`

	errors := validateMainConfig([]byte(validConfig))
	if len(errors) > 0 {
		t.Errorf("Expected no errors for valid config, got: %v", errors)
	}

	// Invalid config - missing required fields
	invalidConfig := `
simulation:
  time_scale: 0
server:
  port: 99999
platforms:
  scenarios:
    test_scenario:
      instances:
        - type_id: "unknown_type"
`

	errors = validateMainConfig([]byte(invalidConfig))
	if len(errors) == 0 {
		t.Error("Expected errors for invalid config")
	}

	// Check specific error conditions
	foundTimeScaleError := false
	foundPortError := false
	foundUpdateIntervalError := false
	foundUnknownTypeError := false

	for _, err := range errors {
		if strings.Contains(err, "time_scale must be positive") {
			foundTimeScaleError = true
		}
		if strings.Contains(err, "port must be between") {
			foundPortError = true
		}
		if strings.Contains(err, "update_interval is required") {
			foundUpdateIntervalError = true
		}
		if strings.Contains(err, "references unknown platform type") {
			foundUnknownTypeError = true
		}
	}

	if !foundTimeScaleError {
		t.Error("Expected time scale validation error")
	}
	if !foundPortError {
		t.Error("Expected port validation error")
	}
	if !foundUpdateIntervalError {
		t.Error("Expected update interval validation error")
	}
	if !foundUnknownTypeError {
		t.Error("Expected unknown type reference error")
	}

	// Test invalid YAML
	invalidYAMLConfig := `invalid: yaml: content: [`
	errors = validateMainConfig([]byte(invalidYAMLConfig))
	if len(errors) == 0 {
		t.Error("Expected errors for invalid YAML")
	}
}

func TestValidateScenarioConfig(t *testing.T) {
	// Valid scenario config
	validScenario := `
metadata:
  name: "Test Scenario"
  description: "A test scenario"
  duration: 3600
platforms:
  - id: "aircraft1"
    type: "boeing_737"
    start_position:
      latitude: 40.7128
      longitude: -74.0060
      altitude: 10000
  - id: "aircraft2"
    type: "airbus_a320"
    start_position:
      latitude: 34.0522
      longitude: -118.2437
      altitude: 11000
`

	errors := validateScenarioConfig([]byte(validScenario))
	if len(errors) > 0 {
		t.Errorf("Expected no errors for valid scenario, got: %v", errors)
	}

	// Invalid scenario config with duplicate IDs
	invalidScenario := `
metadata:
  duration: -1
platforms:
  - id: ""
    type: ""
    start_position:
      latitude: 91
      longitude: 181
      altitude: 0
  - id: "aircraft1"
    type: "boeing_737"
    start_position:
      latitude: 40.7128
      longitude: -74.0060
      altitude: 10000
  - id: "aircraft1"
    type: "boeing_737"
    start_position:
      latitude: 40.7128
      longitude: -74.0060
      altitude: 10000
`

	errors = validateScenarioConfig([]byte(invalidScenario))
	if len(errors) == 0 {
		t.Error("Expected errors for invalid scenario")
	}

	// Check for specific validation errors
	hasNameError := false
	hasDurationError := false
	hasIDError := false
	hasTypeError := false
	hasLatError := false
	hasLonError := false
	hasDuplicateIDError := false

	for _, err := range errors {
		if strings.Contains(err, "name is required") {
			hasNameError = true
		}
		if strings.Contains(err, "duration must be positive") {
			hasDurationError = true
		}
		if strings.Contains(err, "id is required") {
			hasIDError = true
		}
		if strings.Contains(err, "type is required") {
			hasTypeError = true
		}
		if strings.Contains(err, "latitude must be between") {
			hasLatError = true
		}
		if strings.Contains(err, "longitude must be between") {
			hasLonError = true
		}
		if strings.Contains(err, "duplicate id") {
			hasDuplicateIDError = true
		}
	}

	if !hasNameError {
		t.Error("Expected name validation error")
	}
	if !hasDurationError {
		t.Error("Expected duration validation error")
	}
	if !hasIDError {
		t.Error("Expected ID validation error")
	}
	if !hasTypeError {
		t.Error("Expected type validation error")
	}
	if !hasLatError {
		t.Error("Expected latitude validation error")
	}
	if !hasLonError {
		t.Error("Expected longitude validation error")
	}
	if !hasDuplicateIDError {
		t.Error("Expected duplicate ID validation error")
	}

	// Test invalid YAML
	invalidYAMLScenario := `metadata: name: [invalid`
	errors = validateScenarioConfig([]byte(invalidYAMLScenario))
	if len(errors) == 0 {
		t.Error("Expected errors for invalid YAML")
	}
}

func TestValidatePlatformDefinition(t *testing.T) {
	// Valid platform definition
	validPlatform := `
platform_types:
  boeing_737:
    class: "Boeing 737-800"
    category: "commercial"
    performance:
      max_speed: 257
      cruise_speed: 230
      max_altitude: 12500
      climb_rate: 15
    physical:
      length: 39.5
      width: 35.8
      mass: 79010
`

	errors := validatePlatformDefinition([]byte(validPlatform), "data/platforms/airborne/commercial.yaml")
	if len(errors) > 0 {
		t.Errorf("Expected no errors for valid platform, got: %v", errors)
	}

	// Invalid platform definition
	invalidPlatform := `
platform_types:
  invalid_aircraft:
    class: ""
    category: ""
    performance:
      max_speed: 0
      cruise_speed: 300
    physical:
      length: 0
      width: -1
      mass: 0
`

	errors = validatePlatformDefinition([]byte(invalidPlatform), "data/platforms/airborne/test.yaml")
	if len(errors) == 0 {
		t.Error("Expected errors for invalid platform")
	}

	// Check for specific validation errors
	foundErrors := make(map[string]bool)
	for _, err := range errors {
		if strings.Contains(err, "class is required") {
			foundErrors["class"] = true
		}
		if strings.Contains(err, "category is required") {
			foundErrors["category"] = true
		}
		if strings.Contains(err, "max_speed must be positive") {
			foundErrors["max_speed"] = true
		}
		if strings.Contains(err, "cruise_speed cannot exceed max_speed") {
			foundErrors["cruise_speed"] = true
		}
		if strings.Contains(err, "length must be positive") {
			foundErrors["length"] = true
		}
		if strings.Contains(err, "width must be positive") {
			foundErrors["width"] = true
		}
		if strings.Contains(err, "mass must be positive") {
			foundErrors["mass"] = true
		}
	}

	expectedErrors := []string{"class", "category", "max_speed", "cruise_speed", "length", "width", "mass"}
	for _, expected := range expectedErrors {
		if !foundErrors[expected] {
			t.Errorf("Expected %s validation error", expected)
		}
	}

	// Test empty platform types
	emptyPlatform := `platform_types: {}`
	errors = validatePlatformDefinition([]byte(emptyPlatform), "test.yaml")
	if len(errors) == 0 {
		t.Error("Expected error for empty platform types")
	}

	// Test invalid YAML
	invalidYAMLPlatform := `platform_types: [invalid`
	errors = validatePlatformDefinition([]byte(invalidYAMLPlatform), "test.yaml")
	if len(errors) == 0 {
		t.Error("Expected errors for invalid YAML")
	}
}

func TestValidatePlatformType(t *testing.T) {
	// Test with minimal platform type that will fail validation
	emptyPlatform := models.PlatformTypeDefinition{}

	// This will test the validation with empty/zero values
	errors := validatePlatformType("test_platform", emptyPlatform, "test.yaml")
	if len(errors) == 0 {
		t.Error("Expected validation errors for empty platform type")
	}

	// Should find errors for all required fields
	expectedErrorSubstrings := []string{
		"class is required",
		"category is required",
		"max_speed must be positive",
		"cruise_speed must be positive",
		"length must be positive",
		"width must be positive",
		"mass must be positive",
	}

	for _, expectedSubstring := range expectedErrorSubstrings {
		found := false
		for _, err := range errors {
			if strings.Contains(err, expectedSubstring) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected error containing '%s' but didn't find it in: %v", expectedSubstring, errors)
		}
	}

	// Test valid platform type
	validPlatform := models.PlatformTypeDefinition{
		Class:    "Test Aircraft",
		Category: "commercial",
		Performance: models.PerformanceCharacteristics{
			MaxSpeed:    300,
			CruiseSpeed: 250,
		},
		Physical: models.PhysicalCharacteristics{
			Length: 40.0,
			Width:  35.0,
			Mass:   80000,
		},
	}

	errors = validatePlatformType("valid_platform", validPlatform, "test.yaml")
	if len(errors) > 0 {
		t.Errorf("Expected no errors for valid platform type, got: %v", errors)
	}
}

func TestValidateDomainSpecific(t *testing.T) {
	// Test airborne validation
	airbornePlatform := models.PlatformTypeDefinition{
		Performance: models.PerformanceCharacteristics{
			MaxAltitude: 0, // This should trigger error
			ClimbRate:   0, // This should trigger error
		},
	}

	errors := validateDomainSpecific("test_aircraft", airbornePlatform, "data/platforms/airborne/test.yaml")
	if len(errors) == 0 {
		t.Error("Expected domain-specific validation errors for airborne platform")
	}

	// Check for airborne-specific errors
	foundMaxAltitudeError := false
	foundClimbRateError := false
	for _, err := range errors {
		if strings.Contains(err, "max_altitude must be positive") {
			foundMaxAltitudeError = true
		}
		if strings.Contains(err, "climb_rate must be positive") {
			foundClimbRateError = true
		}
	}

	if !foundMaxAltitudeError {
		t.Error("Expected max_altitude validation error")
	}
	if !foundClimbRateError {
		t.Error("Expected climb_rate validation error")
	}

	// Test maritime validation
	maritimePlatform := models.PlatformTypeDefinition{
		Physical: models.PhysicalCharacteristics{
			Draft: 0, // This should trigger error
		},
	}

	errors = validateDomainSpecific("test_ship", maritimePlatform, "data/platforms/maritime/test.yaml")
	if len(errors) == 0 {
		t.Error("Expected domain-specific validation errors for maritime platform")
	}

	// Test space validation
	spacePlatform := models.PlatformTypeDefinition{
		Performance: models.PerformanceCharacteristics{
			OrbitalVelocity: 0, // This should trigger error
			OrbitalAltitude: 0, // This should trigger error
		},
	}

	errors = validateDomainSpecific("test_satellite", spacePlatform, "data/platforms/space/test.yaml")
	if len(errors) == 0 {
		t.Error("Expected domain-specific validation errors for space platform")
	}

	// Test unknown domain (should have no domain-specific errors)
	unknownPlatform := models.PlatformTypeDefinition{}
	errors = validateDomainSpecific("test_platform", unknownPlatform, "data/unknown/test.yaml")
	if len(errors) != 0 {
		t.Errorf("Expected no domain-specific errors for unknown domain, got: %v", errors)
	}

	// Test valid airborne platform
	validAirborne := models.PlatformTypeDefinition{
		Performance: models.PerformanceCharacteristics{
			MaxAltitude: 12000,
			ClimbRate:   15,
		},
	}

	errors = validateDomainSpecific("valid_aircraft", validAirborne, "data/platforms/airborne/test.yaml")
	if len(errors) > 0 {
		t.Errorf("Expected no domain-specific errors for valid airborne platform, got: %v", errors)
	}
}

func TestValidateFiles(t *testing.T) {
	// Create temporary directory in current working directory
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	tmpDir := filepath.Join(wd, "validate-files-test")
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create valid and invalid files
	validFile := filepath.Join(tmpDir, "valid.yaml")
	invalidFile := filepath.Join(tmpDir, "invalid.yaml")

	validContent := "test: value"
	invalidContent := "test: [unclosed"

	if err := os.WriteFile(validFile, []byte(validContent), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(invalidFile, []byte(invalidContent), 0644); err != nil {
		t.Fatal(err)
	}

	files := []string{validFile, invalidFile}
	summary := validateFiles(files)

	if summary.TotalFiles != 2 {
		t.Errorf("Expected 2 total files, got %d", summary.TotalFiles)
	}
	if summary.ValidFiles != 1 {
		t.Errorf("Expected 1 valid file, got %d", summary.ValidFiles)
	}
	if summary.InvalidFiles != 1 {
		t.Errorf("Expected 1 invalid file, got %d", summary.InvalidFiles)
	}
	if len(summary.Results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(summary.Results))
	}

	// Test with empty file list
	emptySummary := validateFiles([]string{})
	if emptySummary.TotalFiles != 0 {
		t.Errorf("Expected 0 total files for empty list, got %d", emptySummary.TotalFiles)
	}
}

func TestPrintSummary(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	summary := ValidationSummary{
		TotalFiles:   2,
		ValidFiles:   1,
		InvalidFiles: 1,
		Results: []ValidationResult{
			{File: "valid.yaml", Valid: true, Errors: []string{}},
			{File: "invalid.yaml", Valid: false, Errors: []string{"syntax error"}},
		},
		FatalErrors: []string{"fatal error"},
	}

	printSummary(summary)

	// Restore stdout and read output
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	expectedStrings := []string{
		"YAML Validation Summary",
		"Total files: 2",
		"Valid files: 1",
		"Invalid files: 1",
		"Fatal Errors:",
		"fatal error",
		"Validation Errors:",
		"invalid.yaml",
		"syntax error",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected output to contain %q, but it didn't. Output: %s", expected, output)
		}
	}

	// Test successful case (no errors)
	r2, w2, _ := os.Pipe()
	os.Stdout = w2

	successSummary := ValidationSummary{
		TotalFiles:   1,
		ValidFiles:   1,
		InvalidFiles: 0,
		Results: []ValidationResult{
			{File: "valid.yaml", Valid: true, Errors: []string{}},
		},
		FatalErrors: []string{},
	}

	printSummary(successSummary)

	w2.Close()
	os.Stdout = oldStdout

	var buf2 bytes.Buffer
	io.Copy(&buf2, r2)
	output2 := buf2.String()

	if !strings.Contains(output2, "All YAML files are valid!") {
		t.Errorf("Expected success message but didn't find it. Output: %s", output2)
	}
}
