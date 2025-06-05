package server

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"github.com/rhino11/trafficsim/internal/config"
	"github.com/rhino11/trafficsim/internal/models"
	"github.com/rhino11/trafficsim/internal/sim"
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
		// In test mode, suppress output unless explicitly enabled
		if os.Getenv("VERBOSE_TESTS") == "1" {
			log.Printf(format, args...)
		}
		// Otherwise suppress output in test mode
	} else {
		// In production mode, always log
		log.Printf(format, args...)
	}
}

// Enhanced logging for web interface debugging
func logWebRequest(r *http.Request, status string) {
	logf("[WEB] %s %s - %s - User-Agent: %s - RemoteAddr: %s",
		r.Method, r.URL.Path, status, r.UserAgent(), r.RemoteAddr)
}

func logWebError(context string, err error) {
	logf("[WEB-ERROR] %s: %v", context, err)
}

func logWebSocket(action string, clientCount int) {
	logf("[WEBSOCKET] %s - Active clients: %d", action, clientCount)
}

func logJSLoad(filename string, status string) {
	logf("[JS-LOAD] %s - %s", filename, status)
}

func logInitialization(component string, status string, duration time.Duration) {
	logf("[INIT] %s - %s (took %v)", component, status, duration)
}

func logSimulationEvent(event string, details interface{}) {
	logf("[SIM] %s - %+v", event, details)
}

func logPerformance(metric string, value interface{}) {
	logf("[PERF] %s: %v", metric, value)
}

func logDebug(component string, message string, data interface{}) {
	logf("[DEBUG] [%s] %s - %+v", component, message, data)
}

func logClientMessage(msgType string, clientAddr string, data interface{}) {
	logf("[CLIENT-MSG] Type: %s, From: %s, Data: %+v", msgType, clientAddr, data)
}

func logDataStream(component string, action string, details interface{}) {
	logf("[STREAM] [%s] %s - %+v", component, action, details)
}

func logPlatformUpdate(platformCount int, action string) {
	logf("[PLATFORM] %s - Count: %d", action, platformCount)
}

// Server represents the web server for the traffic simulation
type Server struct {
	config     *config.Config
	simulation *sim.Engine
	router     *mux.Router
	upgrader   websocket.Upgrader
	clients    map[*websocket.Conn]bool
	clientsMux sync.RWMutex
	broadcast  chan []byte
	ctx        context.Context
	cancel     context.CancelFunc
}

// Client represents a connected WebSocket client
type Client struct {
	conn   *websocket.Conn
	send   chan []byte
	server *Server
}

// Message represents a WebSocket message
type Message struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

// PlatformUpdate represents a platform update message
type PlatformUpdate struct {
	Type      string            `json:"type"`
	Platforms []models.Platform `json:"platforms"`
	Timestamp int64             `json:"timestamp"`
}

// SimulationStatus represents simulation status
type SimulationStatus struct {
	Running       bool    `json:"running"`
	Time          float64 `json:"time"`
	PlatformCount int     `json:"platform_count"`
	Speed         float64 `json:"speed"`
}

// NewServer creates a new web server instance
func NewServer(cfg *config.Config, simulation *sim.Engine) *Server {
	ctx, cancel := context.WithCancel(context.Background())

	server := &Server{
		config:     cfg,
		simulation: simulation,
		router:     mux.NewRouter(),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for development
			},
		},
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan []byte, 256),
		ctx:       ctx,
		cancel:    cancel,
	}

	server.setupRoutes()
	return server
}

// setupRoutes configures all HTTP routes
func (s *Server) setupRoutes() {
	// Static files with logging
	staticHandler := http.StripPrefix("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logWebRequest(r, "SERVING_STATIC")

		// Serve the file
		fileServer := http.FileServer(http.Dir("web/static/"))
		fileServer.ServeHTTP(w, r)

		duration := time.Since(start)
		logPerformance("static_file_serve", map[string]interface{}{
			"file":     r.URL.Path,
			"duration": duration,
		})

		// Log JavaScript file loads specifically
		if strings.HasSuffix(r.URL.Path, ".js") {
			logJSLoad(r.URL.Path, "LOADED")
		}
	}))
	s.router.PathPrefix("/static/").Handler(staticHandler)

	// WebSocket endpoint
	s.router.HandleFunc("/ws", s.handleWebSocket)

	// API endpoints with logging middleware
	api := s.router.PathPrefix("/api").Subrouter()
	api.Use(s.loggingMiddleware)
	api.HandleFunc("/platforms", s.handleGetPlatforms).Methods("GET")
	api.HandleFunc("/platform-types", s.handleGetPlatformTypes).Methods("GET")
	api.HandleFunc("/simulation/start", s.handleStartSimulation).Methods("POST")
	api.HandleFunc("/simulation/stop", s.handleStopSimulation).Methods("POST")
	api.HandleFunc("/simulation/reset", s.handleResetSimulation).Methods("POST")
	api.HandleFunc("/simulation/status", s.handleSimulationStatus).Methods("GET")
	api.HandleFunc("/stream/platforms", s.handleSSEPlatforms).Methods("GET")
	// Performance monitoring endpoint
	api.HandleFunc("/metrics", s.handleMetrics).Methods("GET")
	// Client logging endpoint for debugging
	api.HandleFunc("/log", s.handleClientLog).Methods("POST")
	// Scenario creation endpoint
	api.HandleFunc("/scenarios", s.handleCreateScenario).Methods("POST")

	// Main page
	s.router.HandleFunc("/", s.handleIndex).Methods("GET")
	// Scenario Builder page
	s.router.HandleFunc("/scenario-builder", s.handleScenarioBuilder).Methods("GET")

	logInitialization("Router", "CONFIGURED", 0)
}

// Start starts the web server
func (s *Server) Start(port string) error {
	log.Printf("Starting web server on port %s", port)

	// Start the broadcast goroutine
	go s.handleBroadcast()

	// Start simulation updates if simulation is running
	go s.streamSimulationUpdates()

	// Create HTTP server with proper timeouts for security
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return server.ListenAndServe()
}

// Stop stops the web server
func (s *Server) Stop() {
	s.cancel()

	// Close all WebSocket connections with proper error handling
	s.clientsMux.Lock()
	for client := range s.clients {
		if err := client.Close(); err != nil {
			log.Printf("Error closing WebSocket connection: %v", err)
		}
	}
	s.clientsMux.Unlock()
}

// handleIndex serves the main HTML page
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/templates/index.html")
	if err != nil {
		http.Error(w, "Error loading template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Title string
	}{
		Title: "Traffic Simulation",
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleScenarioBuilder serves the scenario builder HTML page
func (s *Server) handleScenarioBuilder(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/templates/scenario-builder.html")
	if err != nil {
		http.Error(w, "Error loading scenario builder template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Title string
	}{
		Title: "Scenario Builder - Traffic Simulation",
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Error executing scenario builder template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleWebSocket handles WebSocket connections
func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{
		conn:   conn,
		send:   make(chan []byte, 256),
		server: s,
	}

	s.clientsMux.Lock()
	s.clients[conn] = true
	s.clientsMux.Unlock()

	logWebSocket("New connection", len(s.clients))

	// Send initial data
	go s.sendInitialData(client)

	// Start goroutines for this client
	go client.writePump()
	go client.readPump()
}

// sendInitialData sends initial platform data to a new client
func (s *Server) sendInitialData(client *Client) {
	platforms := s.simulation.GetAllPlatforms()

	message := PlatformUpdate{
		Type:      "platform_update",
		Platforms: platforms,
		Timestamp: time.Now().UnixMilli(),
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling initial data: %v", err)
		return
	}

	select {
	case client.send <- data:
	case <-time.After(time.Second):
		log.Printf("Timeout sending initial data to client")
	}
}

// handleGetPlatforms returns all current platforms
func (s *Server) handleGetPlatforms(w http.ResponseWriter, r *http.Request) {
	platforms := s.simulation.GetAllPlatforms()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(platforms); err != nil {
		http.Error(w, "Error encoding platforms: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleGetPlatformTypes returns all available platform types from configuration
func (s *Server) handleGetPlatformTypes(w http.ResponseWriter, r *http.Request) {
	// Load platform types from the configuration
	platformTypes, err := s.loadPlatformTypesFromConfig()
	if err != nil {
		log.Printf("Error loading platform types: %v", err)
		// Return fallback data if config loading fails
		fallbackTypes := s.getFallbackPlatformTypes()
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(fallbackTypes); err != nil {
			http.Error(w, "Error encoding platform types: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(platformTypes); err != nil {
		http.Error(w, "Error encoding platform types: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// PlatformTypeInfo represents a platform type for the scenario builder
type PlatformTypeInfo struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Class       string                 `json:"class"`
	Category    string                 `json:"category"`
	Domain      string                 `json:"domain"`
	Description string                 `json:"description"`
	Performance map[string]interface{} `json:"performance"`
}

// loadPlatformTypesFromConfig loads platform types from YAML configuration files
func (s *Server) loadPlatformTypesFromConfig() ([]PlatformTypeInfo, error) {
	// Load configuration from file
	cfg, err := config.LoadConfig("data/config.yaml")
	if err != nil {
		log.Printf("Could not load configuration: %v", err)
		return s.getDefaultPlatformTypes(), nil
	}

	var platformTypes []PlatformTypeInfo

	// Convert airborne types
	for typeID, typeDef := range cfg.Platforms.AirborneTypes {
		typeInfo := PlatformTypeInfo{
			ID:          typeID,
			Name:        typeDef.Name,
			Class:       typeDef.Class,
			Category:    typeDef.Category,
			Domain:      "airborne",
			Description: s.generateDescription(typeID),
			Performance: map[string]interface{}{
				"max_speed":    typeDef.MaxSpeed,
				"cruise_speed": typeDef.CruiseSpeed,
				"max_altitude": typeDef.MaxAltitude,
			},
		}
		platformTypes = append(platformTypes, typeInfo)
	}

	// Convert maritime types
	for typeID, typeDef := range cfg.Platforms.MaritimeTypes {
		typeInfo := PlatformTypeInfo{
			ID:          typeID,
			Name:        typeDef.Name,
			Class:       typeDef.Class,
			Category:    typeDef.Category,
			Domain:      "maritime",
			Description: s.generateDescription(typeID),
			Performance: map[string]interface{}{
				"max_speed":    typeDef.MaxSpeed,
				"cruise_speed": typeDef.CruiseSpeed,
			},
		}
		platformTypes = append(platformTypes, typeInfo)
	}

	// Convert land types
	for typeID, typeDef := range cfg.Platforms.LandTypes {
		typeInfo := PlatformTypeInfo{
			ID:          typeID,
			Name:        typeDef.Name,
			Class:       typeDef.Class,
			Category:    typeDef.Category,
			Domain:      "land",
			Description: s.generateDescription(typeID),
			Performance: map[string]interface{}{
				"max_speed":    typeDef.MaxSpeed,
				"cruise_speed": typeDef.CruiseSpeed,
			},
		}
		platformTypes = append(platformTypes, typeInfo)
	}

	// Convert space types
	for typeID, typeDef := range cfg.Platforms.SpaceTypes {
		typeInfo := PlatformTypeInfo{
			ID:          typeID,
			Name:        typeDef.Name,
			Class:       typeDef.Class,
			Category:    typeDef.Category,
			Domain:      "space",
			Description: s.generateDescription(typeID),
			Performance: map[string]interface{}{
				"orbital_velocity": typeDef.MaxSpeed,
				"orbital_altitude": typeDef.MaxAltitude,
			},
		}
		platformTypes = append(platformTypes, typeInfo)
	}

	if len(platformTypes) == 0 {
		return s.getDefaultPlatformTypes(), nil
	}

	return platformTypes, nil
}

// determineDomain determines the domain based on platform type
func (s *Server) determineDomain(platformType string) string {
	airborneTypes := []string{"airbus_a320", "boeing_777_300er", "f16_fighting_falcon", "commercial_aircraft", "fighter_aircraft"}
	maritimeTypes := []string{"container_ship", "arleigh_burke_destroyer", "cargo_vessel", "guided_missile_destroyer"}
	landTypes := []string{"truck", "tank", "convoy"}
	spaceTypes := []string{"satellite", "space_station"}

	for _, t := range airborneTypes {
		if strings.Contains(platformType, t) {
			return "airborne"
		}
	}
	for _, t := range maritimeTypes {
		if strings.Contains(platformType, t) {
			return "maritime"
		}
	}
	for _, t := range landTypes {
		if strings.Contains(platformType, t) {
			return "land"
		}
	}
	for _, t := range spaceTypes {
		if strings.Contains(platformType, t) {
			return "space"
		}
	}

	return "unknown"
}

// generateDescription generates a description for a platform type
func (s *Server) generateDescription(platformType string) string {
	descriptions := map[string]string{
		"airbus_a320":             "Short to medium-range commercial airliner",
		"boeing_777_300er":        "Long-range wide-body commercial airliner",
		"f16_fighting_falcon":     "Multi-role fighter aircraft",
		"container_ship":          "Large cargo container vessel",
		"arleigh_burke_destroyer": "US Navy guided missile destroyer",
	}

	if desc, exists := descriptions[platformType]; exists {
		return desc
	}

	// Use proper text casing instead of deprecated strings.Title
	caser := cases.Title(language.English)
	return fmt.Sprintf("%s platform", caser.String(strings.ReplaceAll(platformType, "_", " ")))
}

// getFallbackPlatformTypes returns hardcoded platform types as fallback
func (s *Server) getFallbackPlatformTypes() []PlatformTypeInfo {
	return []PlatformTypeInfo{
		{
			ID:          "airbus_a320",
			Name:        "Airbus A320",
			Class:       "Airbus A320",
			Category:    "commercial_aircraft",
			Domain:      "airborne",
			Description: "Short to medium-range commercial airliner",
			Performance: map[string]interface{}{
				"max_speed":    257.0,
				"cruise_speed": 230.0,
				"max_altitude": 12000.0,
			},
		},
		{
			ID:          "boeing_777_300er",
			Name:        "Boeing 777-300ER",
			Class:       "Boeing 777-300ER",
			Category:    "wide_body_airliner",
			Domain:      "airborne",
			Description: "Long-range wide-body commercial airliner",
			Performance: map[string]interface{}{
				"max_speed":    290.0,
				"cruise_speed": 257.0,
				"max_altitude": 13100.0,
			},
		},
		{
			ID:          "f16_fighting_falcon",
			Name:        "F-16 Fighting Falcon",
			Class:       "F-16 Fighting Falcon",
			Category:    "fighter_aircraft",
			Domain:      "airborne",
			Description: "Multi-role fighter aircraft",
			Performance: map[string]interface{}{
				"max_speed":    588.89,
				"cruise_speed": 261.11,
				"max_altitude": 15240.0,
			},
		},
		{
			ID:          "container_ship",
			Name:        "Container Ship",
			Class:       "Container Ship",
			Category:    "cargo_vessel",
			Domain:      "maritime",
			Description: "Large cargo container vessel",
			Performance: map[string]interface{}{
				"max_speed":    12.9,
				"cruise_speed": 10.8,
			},
		},
		{
			ID:          "arleigh_burke_destroyer",
			Name:        "Arleigh Burke Destroyer",
			Class:       "Arleigh Burke-class Destroyer",
			Category:    "guided_missile_destroyer",
			Domain:      "maritime",
			Description: "US Navy guided missile destroyer",
			Performance: map[string]interface{}{
				"max_speed":    15.4,
				"cruise_speed": 10.3,
			},
		},
	}
}

// getDefaultPlatformTypes returns default platform types when config loading fails
func (s *Server) getDefaultPlatformTypes() []PlatformTypeInfo {
	return []PlatformTypeInfo{
		{
			ID:          "boeing_737_800",
			Name:        "Boeing 737-800",
			Class:       "Boeing 737-800",
			Category:    "commercial_aircraft",
			Domain:      "airborne",
			Description: "Short to medium-range commercial airliner",
			Performance: map[string]interface{}{
				"max_speed":    257.0,
				"cruise_speed": 230.0,
				"max_altitude": 12000.0,
			},
		},
		{
			ID:          "f16_fighting_falcon",
			Name:        "F-16 Fighting Falcon",
			Class:       "F-16 Fighting Falcon",
			Category:    "fighter_aircraft",
			Domain:      "airborne",
			Description: "Multi-role fighter aircraft",
			Performance: map[string]interface{}{
				"max_speed":    588.89,
				"cruise_speed": 261.11,
				"max_altitude": 15240.0,
			},
		},
		{
			ID:          "m1a2_abrams",
			Name:        "M1A2 Abrams",
			Class:       "M1A2 Abrams",
			Category:    "main_battle_tank",
			Domain:      "land",
			Description: "Main battle tank",
			Performance: map[string]interface{}{
				"max_speed":    18.0,
				"cruise_speed": 12.0,
			},
		},
		{
			ID:          "civilian_car",
			Name:        "Civilian Car",
			Class:       "Civilian Car",
			Category:    "civilian_vehicle",
			Domain:      "land",
			Description: "Standard civilian automobile",
			Performance: map[string]interface{}{
				"max_speed":    36.1,
				"cruise_speed": 25.0,
			},
		},
		{
			ID:          "arleigh_burke_destroyer",
			Name:        "Arleigh Burke Destroyer",
			Class:       "Arleigh Burke-class Destroyer",
			Category:    "guided_missile_destroyer",
			Domain:      "maritime",
			Description: "US Navy guided missile destroyer",
			Performance: map[string]interface{}{
				"max_speed":    15.4,
				"cruise_speed": 10.3,
			},
		},
		{
			ID:          "container_ship",
			Name:        "Container Ship",
			Class:       "Container Ship",
			Category:    "cargo_vessel",
			Domain:      "maritime",
			Description: "Large cargo container vessel",
			Performance: map[string]interface{}{
				"max_speed":    12.9,
				"cruise_speed": 10.8,
			},
		},
		{
			ID:          "starlink_satellite",
			Name:        "Starlink Satellite",
			Class:       "Starlink Satellite",
			Category:    "communication_satellite",
			Domain:      "space",
			Description: "Low Earth orbit communication satellite",
			Performance: map[string]interface{}{
				"orbital_velocity": 7.66,
				"orbital_altitude": 550000.0,
			},
		},
	}
}

// handleStartSimulation starts the simulation
func (s *Server) handleStartSimulation(w http.ResponseWriter, r *http.Request) {
	if err := s.simulation.Start(); err != nil {
		http.Error(w, "Error starting simulation: "+err.Error(), http.StatusInternalServerError)
		return
	}

	s.broadcastSimulationStatus()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "started"}); err != nil {
		logWebError("Start simulation response encoding", err)
	}
}

// handleStopSimulation stops the simulation
func (s *Server) handleStopSimulation(w http.ResponseWriter, r *http.Request) {
	s.simulation.Stop()
	s.broadcastSimulationStatus()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "stopped"}); err != nil {
		logWebError("Stop simulation response encoding", err)
	}
}

// handleResetSimulation resets the simulation
func (s *Server) handleResetSimulation(w http.ResponseWriter, r *http.Request) {
	if err := s.simulation.Reset(); err != nil {
		http.Error(w, "Error resetting simulation: "+err.Error(), http.StatusInternalServerError)
		return
	}

	s.broadcastSimulationStatus()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "reset"}); err != nil {
		logWebError("Reset simulation response encoding", err)
	}
}

// handleSimulationStatus returns simulation status
func (s *Server) handleSimulationStatus(w http.ResponseWriter, r *http.Request) {
	status := SimulationStatus{
		Running:       s.simulation.IsRunning(),
		Time:          s.simulation.GetSimulationTime(),
		PlatformCount: len(s.simulation.GetAllPlatforms()),
		Speed:         1.0, // Default speed
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(status); err != nil {
		log.Printf("Error encoding status response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// handleSSEPlatforms handles Server-Sent Events for platform updates
func (s *Server) handleSSEPlatforms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Send initial data
	platforms := s.simulation.GetAllPlatforms()
	if len(platforms) > 0 {
		data, _ := json.Marshal(platforms)
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
		log.Printf("SSE: Sent initial data with %d platforms", len(platforms))
	}

	// Create a ticker to send regular updates
	ticker := time.NewTicker(100 * time.Millisecond) // 10 FPS
	defer ticker.Stop()

	// Keep connection alive ticker
	heartbeatTicker := time.NewTicker(30 * time.Second)
	defer heartbeatTicker.Stop()

	for {
		select {
		case <-r.Context().Done():
			log.Printf("SSE: Client disconnected")
			return
		case <-heartbeatTicker.C:
			fmt.Fprintf(w, ": heartbeat\n\n")
			flusher.Flush()
		case <-ticker.C:
			if s.simulation.IsRunning() {
				platforms := s.simulation.GetAllPlatforms()
				if len(platforms) > 0 {
					message := PlatformUpdate{
						Type:      "platform_update",
						Platforms: platforms,
						Timestamp: time.Now().UnixMilli(),
					}

					data, err := json.Marshal(message)
					if err != nil {
						log.Printf("Error marshaling SSE platform update: %v", err)
						continue
					}

					fmt.Fprintf(w, "data: %s\n\n", data)
					flusher.Flush()
				}
			}
		}
	}
}

// streamSimulationUpdates continuously streams simulation updates to WebSocket clients
func (s *Server) streamSimulationUpdates() {
	ticker := time.NewTicker(100 * time.Millisecond) // 10 FPS
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			if s.simulation.IsRunning() {
				platforms := s.simulation.GetAllPlatforms()

				message := PlatformUpdate{
					Type:      "platform_update",
					Platforms: platforms,
					Timestamp: time.Now().UnixMilli(),
				}

				data, err := json.Marshal(message)
				if err != nil {
					log.Printf("Error marshaling platform update: %v", err)
					continue
				}

				select {
				case s.broadcast <- data:
				default:
					// Channel is full, skip this update
				}
			}
		}
	}
}

// broadcastSimulationStatus broadcasts simulation status to all clients
func (s *Server) broadcastSimulationStatus() {
	status := SimulationStatus{
		Running:       s.simulation.IsRunning(),
		Time:          s.simulation.GetSimulationTime(),
		PlatformCount: len(s.simulation.GetAllPlatforms()),
		Speed:         1.0,
	}

	message := Message{
		Type:      "simulation_status",
		Data:      status,
		Timestamp: time.Now().UnixMilli(),
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling simulation status: %v", err)
		return
	}

	select {
	case s.broadcast <- data:
	default:
		// Channel is full, skip this update
	}
}

// handleBroadcast handles broadcasting messages to all connected clients
func (s *Server) handleBroadcast() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case message := <-s.broadcast:
			s.clientsMux.RLock()
			clientsToRemove := make([]*websocket.Conn, 0)

			for client := range s.clients {
				err := client.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Printf("Error writing to WebSocket client: %v", err)
					clientsToRemove = append(clientsToRemove, client)
				}
			}
			s.clientsMux.RUnlock()

			// Remove failed clients
			if len(clientsToRemove) > 0 {
				s.clientsMux.Lock()
				for _, client := range clientsToRemove {
					delete(s.clients, client)
					if err := client.Close(); err != nil {
						log.Printf("Error closing failed WebSocket client: %v", err)
					}
				}
				s.clientsMux.Unlock()
				logWebSocket("Removed failed clients", len(s.clients))
			}
		}
	}
}

// Client methods for WebSocket handling

// readPump handles reading messages from the WebSocket connection
func (c *Client) readPump() {
	defer func() {
		c.server.clientsMux.Lock()
		delete(c.server.clients, c.conn)
		c.server.clientsMux.Unlock()
		if err := c.conn.Close(); err != nil {
			log.Printf("Error closing WebSocket connection in readPump: %v", err)
		}
	}()

	c.conn.SetReadLimit(512)
	if err := c.conn.SetReadDeadline(time.Now().Add(60 * time.Second)); err != nil {
		log.Printf("Error setting read deadline: %v", err)
	}
	c.conn.SetPongHandler(func(string) error {
		if err := c.conn.SetReadDeadline(time.Now().Add(60 * time.Second)); err != nil {
			log.Printf("Error setting pong read deadline: %v", err)
		}
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Handle incoming messages (ping, control messages, etc.)
		c.handleMessage(message)
	}
}

// writePump handles writing messages to the WebSocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		if err := c.conn.Close(); err != nil {
			log.Printf("Error closing WebSocket connection in writePump: %v", err)
		}
	}()

	for {
		select {
		case message, ok := <-c.send:
			if err := c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
				log.Printf("Error setting write deadline: %v", err)
			}
			if !ok {
				if err := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					log.Printf("Error sending close message: %v", err)
				}
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			if err := c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
				log.Printf("Error setting ping write deadline: %v", err)
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage handles incoming WebSocket messages
func (c *Client) handleMessage(data []byte) {
	clientAddr := c.conn.RemoteAddr().String()

	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		logWebError("Message unmarshaling", err)
		return
	}

	logClientMessage(msg.Type, clientAddr, msg.Data)

	switch msg.Type {
	case "ping":
		// Respond with pong
		response := Message{
			Type:      "pong",
			Timestamp: time.Now().UnixMilli(),
		}
		responseData, _ := json.Marshal(response)
		select {
		case c.send <- responseData:
		default:
		}
		logDebug("WEBSOCKET", "Pong sent", map[string]interface{}{
			"client":  clientAddr,
			"latency": time.Now().UnixMilli() - msg.Timestamp,
		})

	case "viewport_update":
		// Handle viewport changes for server-side filtering
		logDataStream("VIEWPORT", "Update received", msg.Data)
		// Future: Implement viewport-based platform filtering

	case "filter_update":
		// Handle platform filter changes
		logDataStream("FILTER", "Update received", msg.Data)
		// Future: Implement server-side platform filtering

	case "request_initial_data":
		// Send current platform data
		logDataStream("CLIENT", "Initial data requested", clientAddr)
		go c.server.sendInitialData(c)

	case "start_simulation":
		// Handle simulation start command
		logSimulationEvent("START_REQUESTED", map[string]interface{}{
			"client":    clientAddr,
			"timestamp": msg.Timestamp,
		})
		if err := c.server.simulation.Start(); err != nil {
			logWebError("Simulation start", err)
		} else {
			c.server.broadcastSimulationStatus()
		}

	case "stop_simulation":
		// Handle simulation stop command
		logSimulationEvent("STOP_REQUESTED", map[string]interface{}{
			"client":    clientAddr,
			"timestamp": msg.Timestamp,
		})
		c.server.simulation.Stop()
		c.server.broadcastSimulationStatus()

	case "control":
		// Handle other simulation control messages
		logSimulationEvent("CONTROL_MESSAGE", msg.Data)

	default:
		// Log unknown message types with full context for debugging
		logDebug("WEBSOCKET", "Unknown message type received", map[string]interface{}{
			"type":      msg.Type,
			"client":    clientAddr,
			"data":      msg.Data,
			"timestamp": msg.Timestamp,
		})
	}
}

// handleClientLog handles client-side logging for debugging
func (s *Server) handleClientLog(w http.ResponseWriter, r *http.Request) {
	var logData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&logData); err != nil {
		logWebError("Client log decode", err)
		http.Error(w, "Invalid log data", http.StatusBadRequest)
		return
	}

	// Log client-side events with context
	logType, _ := logData["type"].(string)
	step, _ := logData["step"].(string)
	message, _ := logData["message"].(string)
	timestamp, _ := logData["timestamp"].(string)
	userAgent, _ := logData["userAgent"].(string)

	switch logType {
	case "client_log":
		log.Printf("[CLIENT-LOG] [%s] %s: %s - UA: %s", step, timestamp, message, userAgent)
	case "client_error":
		errorMsg, _ := logData["error"].(string)
		context, _ := logData["context"].(string)
		stack, _ := logData["stack"].(string)
		log.Printf("[CLIENT-ERROR] [%s] %s: %s - Context: %s - Stack: %s - UA: %s",
			step, timestamp, errorMsg, context, stack, userAgent)
	default:
		log.Printf("[CLIENT-UNKNOWN] %+v", logData)
	}

	w.WriteHeader(http.StatusNoContent)
}

// loggingMiddleware adds request/response logging to API endpoints
func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logWebRequest(r, "API_REQUEST_START")

		// Create a response writer wrapper to capture status code
		wrapper := &responseWriter{ResponseWriter: w, statusCode: 200}

		next.ServeHTTP(wrapper, r)

		duration := time.Since(start)
		logPerformance("api_request", map[string]interface{}{
			"endpoint": r.URL.Path,
			"method":   r.Method,
			"status":   wrapper.statusCode,
			"duration": duration,
		})
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// handleMetrics returns comprehensive performance metrics
func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	stats := s.simulation.GetStatistics()

	metrics := map[string]interface{}{
		"simulation": stats,
		"server": map[string]interface{}{
			"active_websocket_clients": len(s.clients),
			"uptime_seconds":           time.Since(time.Now()).Seconds(), // Will be corrected with actual start time
		},
		"platforms": map[string]interface{}{
			"total": stats.TotalPlatforms,
			"by_type": map[string]interface{}{
				"airborne": stats.AirbornePlatforms,
				"maritime": stats.MaritimePlatforms,
				"land":     stats.LandPlatforms,
				"space":    stats.SpacePlatforms,
			},
		},
		"timestamp": time.Now().Unix(),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(metrics); err != nil {
		logWebError("Metrics encoding", err)
		http.Error(w, "Error encoding metrics", http.StatusInternalServerError)
		return
	}
}

// handleCreateScenario handles scenario creation requests
func (s *Server) handleCreateScenario(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string                   `json:"name"`
		Description string                   `json:"description"`
		Duration    int                      `json:"duration"`
		Platforms   []map[string]interface{} `json:"platforms"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Create platforms from the request
	var platforms []*models.UniversalPlatform
	for _, platformConfig := range req.Platforms {
		platform, err := models.CreatePlatformFromConfig(platformConfig)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create platform: %v", err), http.StatusBadRequest)
			return
		}
		platforms = append(platforms, platform)
	}

	// Create scenario configuration
	scenario := map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
		"duration":    req.Duration,
		"platforms":   platforms,
		"created_at":  time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Scenario created successfully",
		"scenario": scenario,
	}); err != nil {
		logWebError("Scenario creation response encoding", err)
	}
}
