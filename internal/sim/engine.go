package sim

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/rhino11/trafficsim/internal/config"
	"github.com/rhino11/trafficsim/internal/models"
)

// isTestMode checks if we're running in test mode
func isTestMode() bool {
	return strings.Contains(os.Args[0], ".test") ||
		strings.HasSuffix(os.Args[0], "/test") ||
		os.Getenv("GO_TESTING") == "1"
}

// logf is a conditional logger that respects test verbosity
func logf(format string, args ...interface{}) {
	if isTestMode() {
		// In test mode, suppress output to reduce noise
		return
	} else {
		// In production mode, always log
		fmt.Printf(format+"\n", args...)
	}
}

// Enhanced logging for simulation debugging
func logSimulationStart(platformCount int, updateInterval time.Duration) {
	logf("[SIM-START] Starting simulation with %d platforms at %v update interval", platformCount, updateInterval)
}

func logSimulationStop(reason string) {
	logf("[SIM-STOP] Simulation stopped: %s", reason)
}

func logPlatformOperation(operation string, platformID string, details interface{}) {
	logf("[PLATFORM-OP] %s - ID: %s - Details: %+v", operation, platformID, details)
}

func logSimulationPerformance(updateCount int, avgUpdateTime time.Duration, platformCount int) {
	logf("[SIM-PERF] Updates: %d, Avg Update Time: %v, Platforms: %d", updateCount, avgUpdateTime, platformCount)
}

func logSimulationError(context string, err error, platformID string) {
	logf("[SIM-ERROR] Context: %s, Platform: %s, Error: %v", context, platformID, err)
}

// Engine represents the main simulation engine that orchestrates all platform movement
type Engine struct {
	config         *config.Config
	physics        *PhysicsEngine
	platforms      map[string]models.Platform
	platformsMux   sync.RWMutex
	isRunning      bool
	runningMux     sync.RWMutex
	updateTicker   *time.Ticker
	stopCh         chan struct{}
	simulationTime float64
	timeMux        sync.RWMutex
	updateInterval time.Duration

	// Performance tracking
	updateCount     int64
	totalUpdateTime time.Duration
	lastPerfLog     time.Time
}

// NewEngine creates a new simulation engine
func NewEngine(cfg *config.Config) *Engine {
	updateInterval := time.Second / 60 // 60 FPS default
	if cfg != nil && cfg.Simulation.UpdateInterval != "" {
		if parsed, err := time.ParseDuration(cfg.Simulation.UpdateInterval); err == nil {
			updateInterval = parsed
		}
	}

	return &Engine{
		config:         cfg,
		physics:        NewPhysicsEngine(),
		platforms:      make(map[string]models.Platform),
		stopCh:         make(chan struct{}),
		updateInterval: updateInterval,
	}
}

// Start begins the simulation loop
func (e *Engine) Start() error {
	e.runningMux.Lock()
	defer e.runningMux.Unlock()

	if e.isRunning {
		return fmt.Errorf("simulation is already running")
	}

	e.isRunning = true
	e.stopCh = make(chan struct{})
	e.updateTicker = time.NewTicker(e.updateInterval)

	// Start the simulation loop in a goroutine
	go e.simulationLoop()

	logSimulationStart(len(e.platforms), e.updateInterval)
	return nil
}

// Stop halts the simulation
func (e *Engine) Stop() {
	e.runningMux.Lock()
	defer e.runningMux.Unlock()

	if !e.isRunning {
		return
	}

	e.isRunning = false
	if e.updateTicker != nil {
		e.updateTicker.Stop()
	}
	close(e.stopCh)

	logSimulationStop("User requested stop")
}

// Reset resets the simulation to initial state
func (e *Engine) Reset() error {
	wasRunning := e.IsRunning()
	if wasRunning {
		e.Stop()
	}

	e.timeMux.Lock()
	e.simulationTime = 0
	e.timeMux.Unlock()

	// Reset all platforms to their initial positions
	e.platformsMux.Lock()
	for _, platform := range e.platforms {
		if universalPlatform, ok := platform.(*models.UniversalPlatform); ok {
			universalPlatform.State.Position = universalPlatform.Config.StartPosition
			universalPlatform.State.Speed = 0
			universalPlatform.State.Heading = 0
			universalPlatform.State.Velocity = models.Velocity{}
			universalPlatform.MissionTime = 0
			universalPlatform.State.LastUpdated = time.Now()
		}
	}
	e.platformsMux.Unlock()

	if wasRunning {
		return e.Start()
	}

	logf("Simulation reset")
	return nil
}

// IsRunning returns whether the simulation is currently running
func (e *Engine) IsRunning() bool {
	e.runningMux.RLock()
	defer e.runningMux.RUnlock()
	return e.isRunning
}

// GetSimulationTime returns the current simulation time in seconds
func (e *Engine) GetSimulationTime() float64 {
	e.timeMux.RLock()
	defer e.timeMux.RUnlock()
	return e.simulationTime
}

// AddPlatform adds a platform to the simulation
func (e *Engine) AddPlatform(platform models.Platform) error {
	e.platformsMux.Lock()
	defer e.platformsMux.Unlock()

	id := platform.GetID()
	if _, exists := e.platforms[id]; exists {
		return fmt.Errorf("platform with ID %s already exists", id)
	}

	e.platforms[id] = platform
	logPlatformOperation("ADD", id, platform)
	return nil
}

// RemovePlatform removes a platform from the simulation
func (e *Engine) RemovePlatform(id string) error {
	e.platformsMux.Lock()
	defer e.platformsMux.Unlock()

	if _, exists := e.platforms[id]; !exists {
		return fmt.Errorf("platform with ID %s not found", id)
	}

	delete(e.platforms, id)
	logPlatformOperation("REMOVE", id, nil)
	return nil
}

// GetPlatform returns a platform by ID
func (e *Engine) GetPlatform(id string) (models.Platform, error) {
	e.platformsMux.RLock()
	defer e.platformsMux.RUnlock()

	platform, exists := e.platforms[id]
	if !exists {
		return nil, fmt.Errorf("platform with ID %s not found", id)
	}

	return platform, nil
}

// GetAllPlatforms returns all platforms in the simulation
func (e *Engine) GetAllPlatforms() []models.Platform {
	e.platformsMux.RLock()
	defer e.platformsMux.RUnlock()

	platforms := make([]models.Platform, 0, len(e.platforms))
	for _, platform := range e.platforms {
		platforms = append(platforms, platform)
	}

	return platforms
}

// LoadPlatformsFromConfig loads platforms from configuration
func (e *Engine) LoadPlatformsFromConfig() error {
	if e.config == nil {
		return fmt.Errorf("no configuration provided")
	}

	// This would load platforms from the configuration
	// For now, we'll create some example platforms
	if err := e.createExamplePlatforms(); err != nil {
		return fmt.Errorf("failed to create example platforms: %w", err)
	}

	return nil
}

// createExamplePlatforms creates some example platforms for testing
func (e *Engine) createExamplePlatforms() error {
	// Create example aircraft
	boeing737 := models.NewBoeing737_800Universal(
		"UA123",
		"United 123",
		models.Position{Latitude: 40.7128, Longitude: -74.0060, Altitude: 10000}, // NYC
	)
	if err := boeing737.SetDestination(models.Position{Latitude: 34.0522, Longitude: -118.2437, Altitude: 10000}); err != nil {
		logf("Error setting Boeing 737 destination: %v", err)
	}

	if err := e.AddPlatform(boeing737); err != nil {
		return err
	}

	// Create example ship
	destroyer := models.NewArleighBurkeDestroyerUniversal(
		"DDG-89",
		"Mustin",
		models.Position{Latitude: 36.8485, Longitude: -76.2951, Altitude: 0}, // Norfolk, VA
	)
	if err := destroyer.SetDestination(models.Position{Latitude: 25.7617, Longitude: -80.1918, Altitude: 0}); err != nil {
		logf("Error setting destroyer destination: %v", err)
	}

	if err := e.AddPlatform(destroyer); err != nil {
		return err
	}

	// Create example tank
	tank := models.NewM1A2AbramsUniversal(
		"ARMOR-01",
		"Alpha Company",
		models.Position{Latitude: 31.8720, Longitude: -106.3744, Altitude: 1200}, // El Paso, TX
	)
	if err := tank.SetDestination(models.Position{Latitude: 31.8800, Longitude: -106.3600, Altitude: 1250}); err != nil {
		logf("Error setting tank destination: %v", err)
	}

	if err := e.AddPlatform(tank); err != nil {
		return err
	}

	// Create example satellite
	satellite := models.NewStarlinkSatelliteUniversal(
		"STARLINK-1234",
		"1234",
		models.Position{Latitude: 0, Longitude: 0, Altitude: 550000},
	)

	if err := e.AddPlatform(satellite); err != nil {
		return err
	}

	logf("Created %d example platforms", len(e.platforms))
	return nil
}

// Update performs a single simulation step
func (e *Engine) Update(deltaTime time.Duration) error {
	if !e.IsRunning() {
		return fmt.Errorf("simulation is not running")
	}

	e.platformsMux.RLock()
	platforms := make([]models.Platform, 0, len(e.platforms))
	for _, platform := range e.platforms {
		platforms = append(platforms, platform)
	}
	e.platformsMux.RUnlock()

	// Update all platforms using physics engine
	for _, platform := range platforms {
		if err := e.physics.CalculateMovement(platform, deltaTime); err != nil {
			logf("Error updating platform %s: %v", platform.GetID(), err)
		}
	}

	// Update simulation time
	e.timeMux.Lock()
	e.simulationTime += deltaTime.Seconds()
	e.timeMux.Unlock()

	// Performance tracking
	e.updateCount++
	e.totalUpdateTime += deltaTime
	if time.Since(e.lastPerfLog) > time.Second {
		avgUpdateTime := e.totalUpdateTime / time.Duration(e.updateCount)
		logSimulationPerformance(int(e.updateCount), avgUpdateTime, len(e.platforms))
		e.lastPerfLog = time.Now()
	}

	return nil
}

// simulationLoop runs the main simulation update loop
func (e *Engine) simulationLoop() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Simulation loop panic: %v", r)
		}
	}()

	lastUpdate := time.Now()

	for {
		select {
		case <-e.stopCh:
			return
		case currentTime := <-e.updateTicker.C:
			deltaTime := currentTime.Sub(lastUpdate)
			lastUpdate = currentTime

			if err := e.Update(deltaTime); err != nil {
				log.Printf("Simulation update error: %v", err)
			}
		}
	}
}

// SetUpdateInterval changes the simulation update frequency
func (e *Engine) SetUpdateInterval(interval time.Duration) {
	e.updateInterval = interval

	// If running, restart with new interval
	if e.IsRunning() {
		e.Stop()
		time.Sleep(100 * time.Millisecond) // Brief pause
		if err := e.Start(); err != nil {
			log.Printf("Error restarting simulation with new interval: %v", err)
		}
	}
}

// GetPlatformCount returns the number of platforms in the simulation
func (e *Engine) GetPlatformCount() int {
	e.platformsMux.RLock()
	defer e.platformsMux.RUnlock()
	return len(e.platforms)
}

// GetPlatformsByType returns platforms filtered by type
func (e *Engine) GetPlatformsByType(platformType models.PlatformType) []models.Platform {
	e.platformsMux.RLock()
	defer e.platformsMux.RUnlock()

	var filtered []models.Platform
	for _, platform := range e.platforms {
		if universalPlatform, ok := platform.(*models.UniversalPlatform); ok {
			if universalPlatform.PlatformType == platformType {
				filtered = append(filtered, platform)
			}
		}
	}

	return filtered
}

// GetStatistics returns simulation statistics
func (e *Engine) GetStatistics() SimulationStatistics {
	e.platformsMux.RLock()
	defer e.platformsMux.RUnlock()

	stats := SimulationStatistics{
		TotalPlatforms: len(e.platforms),
		SimulationTime: e.GetSimulationTime(),
		IsRunning:      e.IsRunning(),
		UpdateInterval: e.updateInterval,
	}

	// Count by type
	for _, platform := range e.platforms {
		if universalPlatform, ok := platform.(*models.UniversalPlatform); ok {
			switch universalPlatform.PlatformType {
			case models.PlatformTypeAirborne:
				stats.AirbornePlatforms++
			case models.PlatformTypeMaritime:
				stats.MaritimePlatforms++
			case models.PlatformTypeLand:
				stats.LandPlatforms++
			case models.PlatformTypeSpace:
				stats.SpacePlatforms++
			}
		}
	}

	return stats
}

// SimulationStatistics contains simulation metrics
type SimulationStatistics struct {
	TotalPlatforms    int           `json:"total_platforms"`
	AirbornePlatforms int           `json:"airborne_platforms"`
	MaritimePlatforms int           `json:"maritime_platforms"`
	LandPlatforms     int           `json:"land_platforms"`
	SpacePlatforms    int           `json:"space_platforms"`
	SimulationTime    float64       `json:"simulation_time"`
	IsRunning         bool          `json:"is_running"`
	UpdateInterval    time.Duration `json:"update_interval"`
}

// SetDestinationForPlatform sets a destination for a specific platform
func (e *Engine) SetDestinationForPlatform(id string, destination models.Position) error {
	platform, err := e.GetPlatform(id)
	if err != nil {
		return err
	}

	if universalPlatform, ok := platform.(*models.UniversalPlatform); ok {
		if err := universalPlatform.SetDestination(destination); err != nil {
			return fmt.Errorf("failed to set destination for platform %s: %w", id, err)
		}
		logPlatformOperation("SET_DESTINATION", id, destination)
		return nil
	}

	return fmt.Errorf("platform %s does not support destination setting", id)
}

// GetPlatformStatus returns detailed status for a platform
func (e *Engine) GetPlatformStatus(id string) (*PlatformStatus, error) {
	platform, err := e.GetPlatform(id)
	if err != nil {
		return nil, err
	}

	universalPlatform, ok := platform.(*models.UniversalPlatform)
	if !ok {
		return nil, fmt.Errorf("platform %s is not a universal platform", id)
	}

	status := &PlatformStatus{
		ID:            universalPlatform.ID,
		Name:          universalPlatform.Config.Name,
		Type:          string(universalPlatform.PlatformType),
		Position:      universalPlatform.State.Position,
		Velocity:      universalPlatform.State.Velocity,
		Speed:         universalPlatform.State.Speed,
		Heading:       universalPlatform.State.Heading,
		FuelRemaining: universalPlatform.FuelRemaining,
		SystemStatus:  universalPlatform.SystemStatus,
		LastUpdated:   universalPlatform.State.LastUpdated,
	}

	if universalPlatform.Destination != nil {
		status.Destination = universalPlatform.Destination
		status.DistanceToDestination = e.physics.CalculateGreatCircleDistance(
			universalPlatform.State.Position,
			*universalPlatform.Destination,
		)
	}

	return status, nil
}

// PlatformStatus represents detailed platform status information
type PlatformStatus struct {
	ID                    string              `json:"id"`
	Name                  string              `json:"name"`
	Type                  string              `json:"type"`
	Position              models.Position     `json:"position"`
	Destination           *models.Position    `json:"destination,omitempty"`
	DistanceToDestination float64             `json:"distance_to_destination,omitempty"`
	Velocity              models.Velocity     `json:"velocity"`
	Speed                 float64             `json:"speed"`
	Heading               float64             `json:"heading"`
	FuelRemaining         float64             `json:"fuel_remaining"`
	SystemStatus          models.SystemStatus `json:"system_status"`
	LastUpdated           time.Time           `json:"last_updated"`
}
