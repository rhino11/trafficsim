package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rhino11/trafficsim/internal/config"
	"github.com/rhino11/trafficsim/internal/models"
	"github.com/rhino11/trafficsim/internal/server"
	"github.com/rhino11/trafficsim/internal/sim"
)

func main() {
	// Command line flags
	var (
		configPath = flag.String("config", "data/config.yaml", "Path to configuration file")
		webMode    = flag.Bool("web", false, "Run in web server mode")
		port       = flag.String("port", "8080", "Port for web server")
	)
	flag.Parse()

	fmt.Println("Global Traffic Simulator - Configuration-Driven Demo")
	fmt.Println("====================================================")

	// Load configuration
	fmt.Printf("Loading configuration from: %s\n", *configPath)
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create simulation engine
	engine := sim.NewEngine(cfg)

	if *webMode {
		// Run web server mode
		fmt.Printf("Starting web server on port %s...\n", *port)

		// Load platforms from configuration (needed for web mode)
		fmt.Println("Loading platforms for web simulation...")
		if err := engine.LoadPlatformsFromConfig(); err != nil {
			log.Fatalf("Failed to load platforms: %v", err)
		}

		platforms := engine.GetAllPlatforms()
		fmt.Printf("Loaded %d platforms for web interface\n", len(platforms))

		// Start simulation engine for web mode
		if err := engine.Start(); err != nil {
			log.Fatalf("Failed to start simulation: %v", err)
		}

		// Create server
		srv := server.NewServer(cfg, engine)

		// Start server
		if err := srv.Start(*port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	} else {
		// Run command-line mode
		runCLISimulation(engine, cfg)
	}
}

func runCLISimulation(engine *sim.Engine, cfg *config.Config) {
	fmt.Println("Starting traffic simulation...")

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nReceived interrupt signal, shutting down...")
		cancel()
	}()

	// Load platforms from configuration or create examples
	if err := engine.LoadPlatformsFromConfig(); err != nil {
		cancel() // Cancel context before fatal exit
		log.Fatalf("Failed to load platforms: %v", err)
	}

	platforms := engine.GetAllPlatforms()
	fmt.Printf("Loaded %d platforms\n", len(platforms))

	// Display platform information
	for _, platform := range platforms {
		displayPlatformInfo(platform)
	}

	// Start simulation
	if err := engine.Start(); err != nil {
		log.Fatalf("Failed to start simulation: %v", err)
	}

	// Run simulation monitoring loop
	ticker := time.NewTicker(1 * time.Second) // Status updates every second
	defer ticker.Stop()

	startTime := time.Now()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Simulation stopped")
			engine.Stop()
			return
		case <-ticker.C:
			// Display status
			elapsed := time.Since(startTime)
			simTime := engine.GetSimulationTime()
			platforms := engine.GetAllPlatforms()

			fmt.Printf("Real time: %.1fs, Sim time: %.1fs, Platforms: %d\n",
				elapsed.Seconds(), simTime, len(platforms))

			// Display platform positions
			if len(platforms) > 0 {
				displayPlatformStatus(platforms)
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
