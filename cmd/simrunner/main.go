package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/rhino11/trafficsim/internal/config"
	"github.com/rhino11/trafficsim/internal/models"
)

func main() {
	// Parse command line flags
	configFile := flag.String("config", "data/config.yaml", "Path to configuration file")
	scenario := flag.String("scenario", "east_coast_demo", "Scenario name to run")
	flag.Parse()

	fmt.Println("Global Traffic Simulator - Configuration-Driven Demo")
	fmt.Println("====================================================")

	// Load configuration
	fmt.Printf("Loading configuration from: %s\n", *configFile)
	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Display configuration summary
	fmt.Printf("Simulation Settings:\n")
	fmt.Printf("  Update Interval: %s\n", cfg.Simulation.UpdateInterval)
	fmt.Printf("  Time Scale: %.1fx\n", cfg.Simulation.TimeScale)
	fmt.Printf("  Max Duration: %s\n", cfg.Simulation.MaxDuration)
	fmt.Printf("  Server: %s:%d\n", cfg.Server.Host, cfg.Server.Port)
	fmt.Printf("  CoT Output: %s (rate: %s)\n", cfg.Output.CoT.Endpoint, cfg.Output.CoT.UpdateRate)

	// Display platform type registry
	fmt.Printf("\nPlatform Type Registry:\n")
	fmt.Printf("  Airborne Types: %d\n", len(cfg.Platforms.AirborneTypes))
	fmt.Printf("  Maritime Types: %d\n", len(cfg.Platforms.MaritimeTypes))
	fmt.Printf("  Land Types: %d\n", len(cfg.Platforms.LandTypes))
	fmt.Printf("  Space Types: %d\n", len(cfg.Platforms.SpaceTypes))
	fmt.Printf("  Available Scenarios: %d\n", len(cfg.Platforms.Scenarios))

	// List available scenarios
	fmt.Printf("\nAvailable Scenarios:\n")
	for name, scenarioCfg := range cfg.Platforms.Scenarios {
		fmt.Printf("  - %s: %s (%d platforms)\n", name, scenarioCfg.Description, len(scenarioCfg.Instances))
	}

	// Create platform factory
	factory := config.NewPlatformFactory(&cfg.Platforms)

	// Load the specified scenario
	fmt.Printf("\nLoading scenario: %s\n", *scenario)
	platforms, err := factory.CreateScenario(*scenario)
	if err != nil {
		log.Fatalf("Failed to create scenario: %v", err)
	}

	fmt.Printf("Created %d platform instances from configuration\n", len(platforms))

	// Display platform information
	fmt.Println("\nPlatform Instances:")
	for _, platform := range platforms {
		displayPlatformInfo(platform)
	}

	// Parse update interval
	updateInterval, err := time.ParseDuration(cfg.Simulation.UpdateInterval)
	if err != nil {
		log.Fatalf("Invalid update interval: %v", err)
	}

	// Parse max duration
	maxDuration, err := time.ParseDuration(cfg.Simulation.MaxDuration)
	if err != nil {
		log.Fatalf("Invalid max duration: %v", err)
	}

	// Run simulation
	fmt.Printf("\nRunning simulation for %s with %s updates...\n", maxDuration, updateInterval)

	startTime := time.Now()
	ticker := time.NewTicker(updateInterval)
	defer ticker.Stop()

	lastStatusTime := time.Now()
	statusInterval := 10 * time.Second

	for {
		select {
		case <-ticker.C:
			// Update all platforms
			for _, platform := range platforms {
				err := platform.Update(updateInterval)
				if err != nil {
					log.Printf("Error updating platform %s: %v", platform.GetID(), err)
				}
			}

			// Display status periodically
			if time.Since(lastStatusTime) >= statusInterval {
				fmt.Printf("\n--- Status after %.0f seconds ---\n", time.Since(startTime).Seconds())
				displayPlatformStatus(platforms)
				lastStatusTime = time.Now()
			}

			// Check if simulation should end
			if time.Since(startTime) >= maxDuration {
				fmt.Printf("\nSimulation completed after %s\n", time.Since(startTime))
				return
			}
		}
	}
}

func displayPlatformInfo(platform models.Platform) {
	state := platform.GetState()
	fmt.Printf("  %s (%s) - %s\n", platform.GetName(), platform.GetClass(), platform.GetCallSign())
	fmt.Printf("    Type: %s | Position: %.4f,%.4f @ %.0fm\n",
		platform.GetType(), state.Position.Latitude, state.Position.Longitude, state.Position.Altitude)
	fmt.Printf("    Specs: L=%.1fm, W=%.1fm, H=%.1fm, Mass=%.0fkg\n",
		platform.GetLength(), platform.GetWidth(), platform.GetHeight(), platform.GetMass())
	fmt.Printf("    Performance: Max=%.1fm/s (%.1fkts), Max Alt=%.0fm\n",
		platform.GetMaxSpeed(), platform.GetMaxSpeed()*1.944, platform.GetMaxAltitude())
}

func displayPlatformStatus(platforms []models.Platform) {
	for _, platform := range platforms {
		state := platform.GetState()
		fmt.Printf("  %s: Lat=%.4f, Lon=%.4f, Alt=%.0fm, Speed=%.1fm/s, Hdg=%.1f°\n",
			platform.GetCallSign(),
			state.Position.Latitude,
			state.Position.Longitude,
			state.Position.Altitude,
			state.Speed,
			state.Heading,
		)
	}
}
