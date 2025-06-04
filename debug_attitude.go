package main

import (
	"fmt"
	"math"
	"time"

	"github.com/rhino11/trafficsim/internal/models"
)

func main() {
	startPos := models.Position{Latitude: 40.0, Longitude: -74.0, Altitude: 10000}
	dest := models.Position{Latitude: 40.1, Longitude: -74.0, Altitude: 10000}

	aircraft := models.NewBoeing737_800("TEST-004", "TEST101", startPos)
	if err := aircraft.SetDestination(dest); err != nil {
		fmt.Printf("Error setting destination: %v\n", err)
		return
	}
	aircraft.UniversalPlatform.State.Speed = 200

	fmt.Printf("Initial heading: %.1f°\n", aircraft.UniversalPlatform.State.Heading)

	// Calculate bearing manually to debug
	lat1 := startPos.Latitude * math.Pi / 180.0
	lat2 := dest.Latitude * math.Pi / 180.0
	deltaLon := (dest.Longitude - startPos.Longitude) * math.Pi / 180.0

	y := math.Sin(deltaLon) * math.Cos(lat2)
	x := math.Cos(lat1)*math.Sin(lat2) - math.Sin(lat1)*math.Cos(lat2)*math.Cos(deltaLon)

	bearing := math.Atan2(y, x) * 180.0 / math.Pi
	for bearing < 0 {
		bearing += 360
	}

	fmt.Printf("Bearing to destination: %.1f°\n", bearing)

	headingError := bearing - aircraft.UniversalPlatform.State.Heading
	fmt.Printf("Heading error: %.1f°\n", headingError)

	fmt.Printf("Initial roll: %.3f°\n", aircraft.UniversalPlatform.State.Physics.Attitude.Roll)

	// Test what happens during a full Update call
	if err := aircraft.Update(time.Second); err != nil {
		fmt.Printf("Error updating aircraft: %v\n", err)
		return
	}
	fmt.Printf("Roll after Update: %.3f°\n", aircraft.UniversalPlatform.State.Physics.Attitude.Roll)
	fmt.Printf("Heading after Update: %.1f°\n", aircraft.UniversalPlatform.State.Heading)
}
