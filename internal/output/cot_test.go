package output

import (
	"encoding/xml"
	"strings"
	"testing"
	"time"
)

func TestCoTGenerator_GenerateCoTMessage(t *testing.T) {
	generator := NewCoTGenerator()

	testCases := []struct {
		name     string
		state    PlatformState
		expected map[string]string // key-value pairs to check in XML
	}{
		{
			name: "F-22 Fighter Aircraft",
			state: PlatformState{
				ID:          "F22_001",
				Callsign:    "FALCON01",
				Latitude:    37.4419,
				Longitude:   -122.1430,
				Altitude:    10000.0,
				Speed:       250.0,
				Course:      090.0,
				CoTType:     "a-f-A-M-F",
				Affiliation: "friend",
			},
			expected: map[string]string{
				"uid":     "TRAFFICSIM-F22_001",
				"type":    "a-f-A-M-F",
				"how":     "m-g",
				"version": "2.0",
			},
		},
		{
			name: "Commercial Aircraft",
			state: PlatformState{
				ID:          "UAL1234",
				Callsign:    "UNITED1234",
				Latitude:    40.7128,
				Longitude:   -74.0060,
				Altitude:    35000.0,
				Speed:       125.0,
				Course:      270.0,
				CoTType:     "a-n-A-C-F",
				Affiliation: "neutral",
			},
			expected: map[string]string{
				"uid":  "TRAFFICSIM-UAL1234",
				"type": "a-n-A-C-F",
			},
		},
		{
			name: "Ground Vehicle",
			state: PlatformState{
				ID:          "HMMWV_001",
				Callsign:    "ALPHA01",
				Latitude:    39.0458,
				Longitude:   -76.6413,
				Altitude:    50.0,
				Speed:       15.0,
				Course:      045.0,
				CoTType:     "a-f-G-U-C-V",
				Affiliation: "friend",
			},
			expected: map[string]string{
				"uid":  "TRAFFICSIM-HMMWV_001",
				"type": "a-f-G-U-C-V",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			xmlData, err := generator.GenerateCoTMessage(tc.state)
			if err != nil {
				t.Fatalf("Failed to generate CoT message: %v", err)
			}

			// Parse XML to verify structure
			var event CoTEvent
			if err := xml.Unmarshal(xmlData, &event); err != nil {
				t.Fatalf("Failed to parse generated XML: %v", err)
			}

			// Check expected values
			for key, expectedValue := range tc.expected {
				switch key {
				case "uid":
					if event.UID != expectedValue {
						t.Errorf("Expected UID %s, got %s", expectedValue, event.UID)
					}
				case "type":
					if event.Type != expectedValue {
						t.Errorf("Expected Type %s, got %s", expectedValue, event.Type)
					}
				case "how":
					if event.How != expectedValue {
						t.Errorf("Expected How %s, got %s", expectedValue, event.How)
					}
				case "version":
					if event.Version != expectedValue {
						t.Errorf("Expected Version %s, got %s", expectedValue, event.Version)
					}
				}
			}

			// Verify coordinate accuracy
			if event.Point.Lat != tc.state.Latitude {
				t.Errorf("Expected Latitude %f, got %f", tc.state.Latitude, event.Point.Lat)
			}
			if event.Point.Lon != tc.state.Longitude {
				t.Errorf("Expected Longitude %f, got %f", tc.state.Longitude, event.Point.Lon)
			}
			if event.Point.Hae != tc.state.Altitude {
				t.Errorf("Expected Altitude %f, got %f", tc.state.Altitude, event.Point.Hae)
			}

			// Verify track data
			if event.Detail.Track.Speed != tc.state.Speed {
				t.Errorf("Expected Speed %f, got %f", tc.state.Speed, event.Detail.Track.Speed)
			}
			if event.Detail.Track.Course != tc.state.Course {
				t.Errorf("Expected Course %f, got %f", tc.state.Course, event.Detail.Track.Course)
			}

			// Verify callsign
			if event.Detail.Contact.Callsign != tc.state.Callsign {
				t.Errorf("Expected Callsign %s, got %s", tc.state.Callsign, event.Detail.Contact.Callsign)
			}

			// Verify XML format is valid
			xmlString := string(xmlData)
			if !strings.Contains(xmlString, "<?xml") {
				t.Error("Generated XML should contain XML declaration")
			}
		})
	}
}

func TestGenerateMILSTD2525Type(t *testing.T) {
	testCases := []struct {
		name        string
		category    string
		affiliation string
		dimension   string
		expected    string
	}{
		{
			name:        "Friendly Fighter Aircraft",
			category:    "fighter_aircraft",
			affiliation: "friend",
			dimension:   "air",
			expected:    "a-f-A-M-F",
		},
		{
			name:        "Hostile UAV",
			category:    "unmanned_aircraft",
			affiliation: "hostile",
			dimension:   "air",
			expected:    "a-h-A-M-U",
		},
		{
			name:        "Neutral Commercial Aircraft",
			category:    "commercial_aircraft",
			affiliation: "neutral",
			dimension:   "air",
			expected:    "a-n-A-C-F",
		},
		{
			name:        "Friendly Main Battle Tank",
			category:    "main_battle_tank",
			affiliation: "friend",
			dimension:   "ground",
			expected:    "a-f-G-U-C-I",
		},
		{
			name:        "Unknown Ground Vehicle",
			category:    "tactical_vehicle",
			affiliation: "unknown",
			dimension:   "ground",
			expected:    "a-u-G-U-C-V",
		},
		{
			name:        "Neutral Commercial Vehicle",
			category:    "commercial_vehicle",
			affiliation: "neutral",
			dimension:   "ground",
			expected:    "a-n-G-U-C-V",
		},
		{
			name:        "Friendly Destroyer",
			category:    "destroyer",
			affiliation: "friend",
			dimension:   "sea",
			expected:    "a-f-S-U-W-D",
		},
		{
			name:        "Neutral Cargo Vessel",
			category:    "cargo_vessel",
			affiliation: "neutral",
			dimension:   "sea",
			expected:    "a-n-S-U-C-V",
		},
		{
			name:        "Unknown Space Asset",
			category:    "satellite",
			affiliation: "unknown",
			dimension:   "space",
			expected:    "a-u-P-U-S",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GenerateMILSTD2525Type(tc.category, tc.affiliation, tc.dimension)
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

func TestCalculateCourse(t *testing.T) {
	testCases := []struct {
		name      string
		lat1      float64
		lon1      float64
		lat2      float64
		lon2      float64
		expected  float64
		tolerance float64
	}{
		{
			name:      "North",
			lat1:      0.0,
			lon1:      0.0,
			lat2:      1.0,
			lon2:      0.0,
			expected:  0.0,
			tolerance: 1.0,
		},
		{
			name:      "East",
			lat1:      0.0,
			lon1:      0.0,
			lat2:      0.0,
			lon2:      1.0,
			expected:  90.0,
			tolerance: 1.0,
		},
		{
			name:      "South",
			lat1:      1.0,
			lon1:      0.0,
			lat2:      0.0,
			lon2:      0.0,
			expected:  180.0,
			tolerance: 1.0,
		},
		{
			name:      "West",
			lat1:      0.0,
			lon1:      1.0,
			lat2:      0.0,
			lon2:      0.0,
			expected:  270.0,
			tolerance: 1.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := CalculateCourse(tc.lat1, tc.lon1, tc.lat2, tc.lon2)

			// Handle the circular nature of bearing (0° = 360°)
			diff := result - tc.expected
			if diff > 180 {
				diff -= 360
			} else if diff < -180 {
				diff += 360
			}

			if abs(diff) > tc.tolerance {
				t.Errorf("Expected course %f±%f, got %f (diff: %f)", tc.expected, tc.tolerance, result, diff)
			}
		})
	}
}

func TestCoTGenerator_SetStaleTime(t *testing.T) {
	generator := NewCoTGenerator()

	// Test default stale time
	state := PlatformState{
		ID:       "TEST_001",
		Callsign: "TEST01",
		CoTType:  "a-f-A-M-F",
	}

	xmlData, err := generator.GenerateCoTMessage(state)
	if err != nil {
		t.Fatalf("Failed to generate CoT message: %v", err)
	}

	var event CoTEvent
	if err := xml.Unmarshal(xmlData, &event); err != nil {
		t.Fatalf("Failed to parse XML: %v", err)
	}

	startTime, err := time.Parse("2006-01-02T15:04:05.000Z", event.Start)
	if err != nil {
		t.Fatalf("Failed to parse start time: %v", err)
	}

	staleTime, err := time.Parse("2006-01-02T15:04:05.000Z", event.Stale)
	if err != nil {
		t.Fatalf("Failed to parse stale time: %v", err)
	}

	defaultDuration := staleTime.Sub(startTime)
	if defaultDuration != 15*time.Minute {
		t.Errorf("Expected default stale time of 15 minutes, got %v", defaultDuration)
	}

	// Test custom stale time
	customStaleTime := 30 * time.Minute
	generator.SetStaleTime(customStaleTime)

	xmlData2, err := generator.GenerateCoTMessage(state)
	if err != nil {
		t.Fatalf("Failed to generate CoT message with custom stale time: %v", err)
	}

	var event2 CoTEvent
	if err := xml.Unmarshal(xmlData2, &event2); err != nil {
		t.Fatalf("Failed to parse XML with custom stale time: %v", err)
	}

	startTime2, err := time.Parse("2006-01-02T15:04:05.000Z", event2.Start)
	if err != nil {
		t.Fatalf("Failed to parse start time: %v", err)
	}

	staleTime2, err := time.Parse("2006-01-02T15:04:05.000Z", event2.Stale)
	if err != nil {
		t.Fatalf("Failed to parse stale time: %v", err)
	}

	customDuration := staleTime2.Sub(startTime2)
	if customDuration != customStaleTime {
		t.Errorf("Expected custom stale time of %v, got %v", customStaleTime, customDuration)
	}
}

// Helper function for absolute value
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
