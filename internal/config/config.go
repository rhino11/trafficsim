package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the main configuration structure
type Config struct {
	Simulation SimulationConfig `yaml:"simulation"`
	Platforms  PlatformRegistry `yaml:"platforms"`
	Server     ServerConfig     `yaml:"server"`
	Output     OutputConfig     `yaml:"output"`
}

// SimulationConfig contains simulation runtime parameters
type SimulationConfig struct {
	UpdateInterval string       `yaml:"update_interval" default:"1s"`
	TimeScale      float64      `yaml:"time_scale" default:"1.0"`
	MaxDuration    string       `yaml:"max_duration" default:"1h"`
	StartTime      string       `yaml:"start_time,omitempty"`
	BoundingBox    *BoundingBox `yaml:"bounding_box,omitempty"`
}

// BoundingBox defines simulation area limits
type BoundingBox struct {
	North float64 `yaml:"north"`
	South float64 `yaml:"south"`
	East  float64 `yaml:"east"`
	West  float64 `yaml:"west"`
}

// ServerConfig contains web server settings
type ServerConfig struct {
	Port    int    `yaml:"port" default:"8080"`
	Host    string `yaml:"host" default:"localhost"`
	WebRoot string `yaml:"web_root" default:"web"`
}

// OutputConfig contains CoT and other output settings
type OutputConfig struct {
	CoT     CoTConfig     `yaml:"cot"`
	Logging LoggingConfig `yaml:"logging"`
}

// CoTConfig contains Cursor-on-Target output settings
type CoTConfig struct {
	Enabled    bool   `yaml:"enabled" default:"true"`
	Endpoint   string `yaml:"endpoint" default:"udp://239.2.3.1:6969"`
	UpdateRate string `yaml:"update_rate" default:"5s"`
}

// LoggingConfig contains logging settings
type LoggingConfig struct {
	Level  string `yaml:"level" default:"info"`
	Format string `yaml:"format" default:"text"`
	File   string `yaml:"file,omitempty"`
}

// PlatformRegistry contains all platform type definitions and scenarios
type PlatformRegistry struct {
	// Platform type definitions (the "database tables")
	AirborneTypes PlatformTypeDefinitions `yaml:"airborne_types"`
	MaritimeTypes PlatformTypeDefinitions `yaml:"maritime_types"`
	LandTypes     PlatformTypeDefinitions `yaml:"land_types"`
	SpaceTypes    PlatformTypeDefinitions `yaml:"space_types"`

	// Scenario definitions (instances to create)
	Scenarios map[string]ScenarioConfig `yaml:"scenarios"`
}

// PlatformTypeDefinitions represents a collection of platform type definitions
type PlatformTypeDefinitions map[string]PlatformTypeDefinition

// PlatformTypeDefinition defines a real-world platform type with all its characteristics
type PlatformTypeDefinition struct {
	// Basic identification
	Name     string `yaml:"name"`
	Class    string `yaml:"class"`
	Type     string `yaml:"type"`     // airborne, maritime, land, space
	Category string `yaml:"category"` // commercial, military, civilian

	// Performance characteristics
	MaxSpeed       float64 `yaml:"max_speed"`                 // m/s
	CruiseSpeed    float64 `yaml:"cruise_speed"`              // m/s
	MaxAltitude    float64 `yaml:"max_altitude"`              // meters
	ServiceCeiling float64 `yaml:"service_ceiling,omitempty"` // meters

	// Physical characteristics
	Length float64 `yaml:"length"` // meters
	Width  float64 `yaml:"width"`  // meters (wingspan for aircraft, beam for ships)
	Height float64 `yaml:"height"` // meters
	Mass   float64 `yaml:"mass"`   // kg

	// Type-specific characteristics
	Draft        float64 `yaml:"draft,omitempty"`         // meters (ships)
	Displacement float64 `yaml:"displacement,omitempty"`  // tonnes (ships)
	FuelCapacity float64 `yaml:"fuel_capacity,omitempty"` // liters (land/air)
	Range        float64 `yaml:"range,omitempty"`         // meters
	MaxGradient  float64 `yaml:"max_gradient,omitempty"`  // degrees (land)

	// Orbital characteristics (space)
	OrbitalPeriod float64 `yaml:"orbital_period,omitempty"` // seconds
	Apogee        float64 `yaml:"apogee,omitempty"`         // meters
	Perigee       float64 `yaml:"perigee,omitempty"`        // meters
	Inclination   float64 `yaml:"inclination,omitempty"`    // degrees

	// Call sign patterns
	CallSignPrefix string `yaml:"callsign_prefix,omitempty"`
	CallSignFormat string `yaml:"callsign_format,omitempty"` // e.g., "{prefix}{id}"
}

// ScenarioConfig defines a simulation scenario with platform instances
type ScenarioConfig struct {
	Name        string             `yaml:"name"`
	Description string             `yaml:"description,omitempty"`
	Duration    string             `yaml:"duration,omitempty"`
	Instances   []PlatformInstance `yaml:"instances"`
}

// PlatformInstance defines a specific platform instance in a scenario
type PlatformInstance struct {
	ID          string          `yaml:"id"`
	TypeID      string          `yaml:"type_id"`            // References PlatformTypeDefinition
	Name        string          `yaml:"name"`               // Display name/flight number
	CallSign    string          `yaml:"callsign,omitempty"` // Override callsign
	StartPos    Position        `yaml:"start_position"`
	Destination *Position       `yaml:"destination,omitempty"`
	Route       []Position      `yaml:"route,omitempty"`
	Behavior    *BehaviorConfig `yaml:"behavior,omitempty"`
}

// Position represents a 3D position
type Position struct {
	Latitude  float64 `yaml:"latitude"`
	Longitude float64 `yaml:"longitude"`
	Altitude  float64 `yaml:"altitude"`
}

// BehaviorConfig defines platform-specific behavior parameters
type BehaviorConfig struct {
	Patrol        *PatrolBehavior     `yaml:"patrol,omitempty"`
	CircuitFlight *CircuitBehavior    `yaml:"circuit,omitempty"`
	RandomWalk    *RandomWalkBehavior `yaml:"random_walk,omitempty"`
}

// PatrolBehavior defines patrol pattern behavior
type PatrolBehavior struct {
	Pattern   string     `yaml:"pattern"` // "line", "box", "circle"
	Points    []Position `yaml:"points"`
	Speed     float64    `yaml:"speed,omitempty"`
	LoopCount int        `yaml:"loop_count,omitempty"` // -1 for infinite
}

// CircuitBehavior defines circuit flight pattern
type CircuitBehavior struct {
	Center Position `yaml:"center"`
	Radius float64  `yaml:"radius"`
	Speed  float64  `yaml:"speed,omitempty"`
}

// RandomWalkBehavior defines random movement
type RandomWalkBehavior struct {
	Area        BoundingBox `yaml:"area"`
	MaxDistance float64     `yaml:"max_distance"`
	Speed       float64     `yaml:"speed,omitempty"`
}

// LoadConfig loads configuration from a YAML file
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Apply defaults
	if err := applyDefaults(&config); err != nil {
		return nil, fmt.Errorf("failed to apply defaults: %w", err)
	}

	// Validate configuration
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// applyDefaults applies default values to configuration
func applyDefaults(config *Config) error {
	// Apply simulation defaults
	if config.Simulation.UpdateInterval == "" {
		config.Simulation.UpdateInterval = "1s"
	}
	if config.Simulation.TimeScale == 0 {
		config.Simulation.TimeScale = 1.0
	}
	if config.Simulation.MaxDuration == "" {
		config.Simulation.MaxDuration = "1h"
	}

	// Apply server defaults
	if config.Server.Port == 0 {
		config.Server.Port = 8080
	}
	if config.Server.Host == "" {
		config.Server.Host = "localhost"
	}
	if config.Server.WebRoot == "" {
		config.Server.WebRoot = "web"
	}

	// Apply output defaults
	if config.Output.CoT.UpdateRate == "" {
		config.Output.CoT.UpdateRate = "5s"
	}
	if config.Output.CoT.Endpoint == "" {
		config.Output.CoT.Endpoint = "udp://239.2.3.1:6969"
	}
	if config.Output.Logging.Level == "" {
		config.Output.Logging.Level = "info"
	}
	if config.Output.Logging.Format == "" {
		config.Output.Logging.Format = "text"
	}

	return nil
}

// validateConfig validates the configuration
func validateConfig(config *Config) error {
	// Validate server port
	if config.Server.Port < 1 || config.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", config.Server.Port)
	}

	// Validate time scale
	if config.Simulation.TimeScale <= 0 {
		return fmt.Errorf("invalid time scale: %f", config.Simulation.TimeScale)
	}

	// Validate platform type references in scenarios
	for scenarioName, scenario := range config.Platforms.Scenarios {
		for i, instance := range scenario.Instances {
			if !config.Platforms.HasType(instance.TypeID) {
				return fmt.Errorf("scenario %s, instance %d: unknown platform type %s",
					scenarioName, i, instance.TypeID)
			}
		}
	}

	return nil
}

// HasType checks if a platform type exists in the registry
func (pr *PlatformRegistry) HasType(typeID string) bool {
	if _, exists := pr.AirborneTypes[typeID]; exists {
		return true
	}
	if _, exists := pr.MaritimeTypes[typeID]; exists {
		return true
	}
	if _, exists := pr.LandTypes[typeID]; exists {
		return true
	}
	if _, exists := pr.SpaceTypes[typeID]; exists {
		return true
	}
	return false
}

// GetType retrieves a platform type definition by ID
func (pr *PlatformRegistry) GetType(typeID string) (*PlatformTypeDefinition, error) {
	if def, exists := pr.AirborneTypes[typeID]; exists {
		return &def, nil
	}
	if def, exists := pr.MaritimeTypes[typeID]; exists {
		return &def, nil
	}
	if def, exists := pr.LandTypes[typeID]; exists {
		return &def, nil
	}
	if def, exists := pr.SpaceTypes[typeID]; exists {
		return &def, nil
	}
	return nil, fmt.Errorf("platform type not found: %s", typeID)
}
