package output

import (
	"encoding/xml"
	"testing"
	"time"
)

func TestMulticastPublisher_Integration(t *testing.T) {
	// Skip this test in CI or when network multicast isn't available
	if testing.Short() {
		t.Skip("Skipping multicast integration test in short mode")
	}

	// Test that we can create a publisher and it generates valid CoT messages
	publisher, err := NewMulticastPublisher("239.2.3.1", 6969)
	if err != nil {
		t.Fatalf("Failed to create multicast publisher: %v", err)
	}
	defer publisher.Close()

	// Test platform states
	testStates := []PlatformState{
		{
			ID:          "F16_001",
			Callsign:    "VIPER01",
			Latitude:    39.0458,
			Longitude:   -76.6413,
			Altitude:    10000.0,
			Speed:       250.0,
			Course:      090.0,
			CoTType:     "a-f-A-M-F",
			Affiliation: "friend",
		},
		{
			ID:          "TANK001",
			Callsign:    "ARMOR01",
			Latitude:    39.0500,
			Longitude:   -76.6500,
			Altitude:    100.0,
			Speed:       15.0,
			Course:      045.0,
			CoTType:     "a-f-G-U-C-I",
			Affiliation: "friend",
		},
	}

	// Test that publishing doesn't fail (we can't easily test reception in unit tests)
	for _, state := range testStates {
		err := publisher.PublishPlatformState(state)
		if err != nil {
			t.Errorf("Failed to publish platform state for %s: %v", state.ID, err)
		}

		// Verify the generated message is valid CoT XML
		cotMessage, err := publisher.generator.GenerateCoTMessage(state)
		if err != nil {
			t.Errorf("Failed to generate CoT message for %s: %v", state.ID, err)
			continue
		}

		// Parse the XML to ensure it's valid
		var event CoTEvent
		if err := xml.Unmarshal(cotMessage, &event); err != nil {
			t.Errorf("Generated CoT message is not valid XML for %s: %v", state.ID, err)
			continue
		}

		// Verify the UID format
		expectedUID := "TRAFFICSIM-" + state.ID
		if event.UID != expectedUID {
			t.Errorf("Expected UID %s, got %s", expectedUID, event.UID)
		}

		// Verify the CoT type
		if event.Type != state.CoTType {
			t.Errorf("Expected type %s, got %s", state.CoTType, event.Type)
		}

		t.Logf("Successfully generated and validated CoT message for %s (Type: %s)", state.Callsign, state.CoTType)
	}
}

func TestMulticastPublisher_ContinuousPublishing(t *testing.T) {
	// Test continuous publishing with a channel
	publisher, err := NewMulticastPublisher("239.2.3.1", 6970) // Different port to avoid conflicts
	if err != nil {
		t.Fatalf("Failed to create multicast publisher: %v", err)
	}
	defer publisher.Close()

	// Set a short publish interval for testing
	publisher.SetPublishInterval(100 * time.Millisecond)

	// Create channels for communication
	platformStates := make(chan PlatformState, 10)
	stopChan := make(chan bool, 1)

	// Start publishing in a goroutine
	go publisher.StartPublishing(platformStates, stopChan)

	// Send test data
	testState := PlatformState{
		ID:          "TEST001",
		Callsign:    "TEST01",
		Latitude:    40.0,
		Longitude:   -75.0,
		Altitude:    1000.0,
		Speed:       100.0,
		Course:      180.0,
		CoTType:     "a-f-A-M-F",
		Affiliation: "friend",
	}

	// Send a few states
	for i := 0; i < 3; i++ {
		platformStates <- testState
		time.Sleep(50 * time.Millisecond)
	}

	// Stop publishing
	stopChan <- true

	// Wait a bit for cleanup
	time.Sleep(200 * time.Millisecond)

	// Test passed if no panics occurred
}

func TestCoTGenerator_GenerateFromPlatformData(t *testing.T) {
	// Test generating CoT from realistic platform data
	generator := NewCoTGenerator()

	// Test different platform types with proper CoT types
	testCases := []struct {
		name         string
		state        PlatformState
		expectedType string
	}{
		{
			name: "F-16 Fighter",
			state: PlatformState{
				ID:          "F16_001",
				Callsign:    "VIPER01",
				Latitude:    39.0458,
				Longitude:   -76.6413,
				Altitude:    10000.0,
				Speed:       250.0,
				Course:      090.0,
				CoTType:     "a-f-A-M-F",
				Affiliation: "friend",
			},
			expectedType: "a-f-A-M-F",
		},
		{
			name: "M1A2 Abrams Tank",
			state: PlatformState{
				ID:          "TANK001",
				Callsign:    "ARMOR01",
				Latitude:    39.0500,
				Longitude:   -76.6500,
				Altitude:    100.0,
				Speed:       15.0,
				Course:      045.0,
				CoTType:     "a-f-G-U-C-I",
				Affiliation: "friend",
			},
			expectedType: "a-f-G-U-C-I",
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
			expectedType: "a-n-A-C-F",
		},
		{
			name: "Arleigh Burke Destroyer",
			state: PlatformState{
				ID:          "DDG051",
				Callsign:    "USS_COLE",
				Latitude:    36.8467,
				Longitude:   -76.2950,
				Altitude:    0.0,
				Speed:       10.3,
				Course:      180.0,
				CoTType:     "a-f-S-U-W-D",
				Affiliation: "friend",
			},
			expectedType: "a-f-S-U-W-D",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			xmlData, err := generator.GenerateCoTMessage(tc.state)
			if err != nil {
				t.Fatalf("Failed to generate CoT message: %v", err)
			}

			// Parse the generated XML
			var event CoTEvent
			if err := xml.Unmarshal(xmlData, &event); err != nil {
				t.Fatalf("Failed to parse generated XML: %v", err)
			}

			// Verify the type
			if event.Type != tc.expectedType {
				t.Errorf("Expected type %s, got %s", tc.expectedType, event.Type)
			}

			// Verify the UID format
			expectedUID := "TRAFFICSIM-" + tc.state.ID
			if event.UID != expectedUID {
				t.Errorf("Expected UID %s, got %s", expectedUID, event.UID)
			}

			// Verify coordinates
			if event.Point.Lat != tc.state.Latitude {
				t.Errorf("Expected latitude %f, got %f", tc.state.Latitude, event.Point.Lat)
			}
			if event.Point.Lon != tc.state.Longitude {
				t.Errorf("Expected longitude %f, got %f", tc.state.Longitude, event.Point.Lon)
			}
			if event.Point.Hae != tc.state.Altitude {
				t.Errorf("Expected altitude %f, got %f", tc.state.Altitude, event.Point.Hae)
			}

			// Verify track data
			if event.Detail.Track.Speed != tc.state.Speed {
				t.Errorf("Expected speed %f, got %f", tc.state.Speed, event.Detail.Track.Speed)
			}
			if event.Detail.Track.Course != tc.state.Course {
				t.Errorf("Expected course %f, got %f", tc.state.Course, event.Detail.Track.Course)
			}

			// Verify callsign
			if event.Detail.Contact.Callsign != tc.state.Callsign {
				t.Errorf("Expected callsign %s, got %s", tc.state.Callsign, event.Detail.Contact.Callsign)
			}
		})
	}
}
