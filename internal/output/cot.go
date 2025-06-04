package output

import (
	"encoding/xml"
	"fmt"
	"math"
	"time"
)

// CoTEvent represents a Cursor on Target event message
type CoTEvent struct {
	XMLName xml.Name  `xml:"event"`
	Version string    `xml:"version,attr"`
	UID     string    `xml:"uid,attr"`
	Type    string    `xml:"type,attr"`
	How     string    `xml:"how,attr"`
	Time    string    `xml:"time,attr"`
	Start   string    `xml:"start,attr"`
	Stale   string    `xml:"stale,attr"`
	Point   CoTPoint  `xml:"point"`
	Detail  CoTDetail `xml:"detail"`
}

// CoTPoint represents the geographical point in a CoT message
type CoTPoint struct {
	Lat float64 `xml:"lat,attr"`
	Lon float64 `xml:"lon,attr"`
	Hae float64 `xml:"hae,attr"`
	CE  float64 `xml:"ce,attr"` // Circular Error (meters)
	LE  float64 `xml:"le,attr"` // Linear Error (meters)
}

// CoTDetail contains platform-specific details
type CoTDetail struct {
	Contact CoTContact `xml:"contact"`
	Track   CoTTrack   `xml:"track"`
	Precis  CoTPrecis  `xml:"precisionlocation"`
}

// CoTContact contains contact information
type CoTContact struct {
	Callsign string `xml:"callsign,attr"`
	Endpoint string `xml:"endpoint,attr,omitempty"`
}

// CoTTrack contains movement information
type CoTTrack struct {
	Speed  float64 `xml:"speed,attr"`  // m/s
	Course float64 `xml:"course,attr"` // degrees true
}

// CoTPrecis contains precision location data
type CoTPrecis struct {
	Geopointsrc string `xml:"geopointsrc,attr"`
	Altsrc      string `xml:"altsrc,attr"`
}

// PlatformState represents the current state of a platform
type PlatformState struct {
	ID          string
	Callsign    string
	Latitude    float64
	Longitude   float64
	Altitude    float64
	Speed       float64
	Course      float64
	CoTType     string
	Affiliation string
}

// CoTGenerator generates Cursor on Target messages
type CoTGenerator struct {
	staleTime time.Duration
}

// NewCoTGenerator creates a new CoT message generator
func NewCoTGenerator() *CoTGenerator {
	return &CoTGenerator{
		staleTime: 15 * time.Minute, // Default stale time
	}
}

// GenerateCoTMessage creates a CoT XML message from platform state
func (g *CoTGenerator) GenerateCoTMessage(state PlatformState) ([]byte, error) {
	now := time.Now().UTC()
	staleTime := now.Add(g.staleTime)

	event := CoTEvent{
		Version: "2.0",
		UID:     fmt.Sprintf("TRAFFICSIM-%s", state.ID),
		Type:    state.CoTType,
		How:     "m-g", // machine-generated
		Time:    now.Format("2006-01-02T15:04:05.000Z"),
		Start:   now.Format("2006-01-02T15:04:05.000Z"),
		Stale:   staleTime.Format("2006-01-02T15:04:05.000Z"),
		Point: CoTPoint{
			Lat: state.Latitude,
			Lon: state.Longitude,
			Hae: state.Altitude,
			CE:  10.0, // 10 meter circular error
			LE:  10.0, // 10 meter linear error
		},
		Detail: CoTDetail{
			Contact: CoTContact{
				Callsign: state.Callsign,
				Endpoint: fmt.Sprintf("trafficsim:%s", state.ID),
			},
			Track: CoTTrack{
				Speed:  state.Speed,
				Course: state.Course,
			},
			Precis: CoTPrecis{
				Geopointsrc: "GPS",
				Altsrc:      "GPS",
			},
		},
	}

	xmlData, err := xml.MarshalIndent(event, "", "  ")
	if err != nil {
		return nil, err
	}

	// Add XML declaration
	xmlDeclaration := []byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	return append(xmlDeclaration, xmlData...), nil
}

// GenerateMILSTD2525Type generates MIL-STD-2525D type codes based on platform category
func GenerateMILSTD2525Type(category, affiliation, dimension string) string {
	var affiliationCode string
	switch affiliation {
	case "friend":
		affiliationCode = "f"
	case "hostile":
		affiliationCode = "h"
	case "neutral":
		affiliationCode = "n"
	case "unknown":
		affiliationCode = "u"
	default:
		affiliationCode = "u"
	}

	var dimensionCode string
	switch dimension {
	case "air":
		dimensionCode = "A"
	case "ground":
		dimensionCode = "G"
	case "sea":
		dimensionCode = "S"
	case "space":
		dimensionCode = "P"
	default:
		dimensionCode = "G"
	}

	// Generate type based on category and dimension
	switch {
	case dimension == "air" && category == "fighter_aircraft":
		return fmt.Sprintf("a-%s-%s-M-F", affiliationCode, dimensionCode)
	case dimension == "air" && category == "unmanned_aircraft":
		return fmt.Sprintf("a-%s-%s-M-U", affiliationCode, dimensionCode)
	case dimension == "air" && category == "commercial_aircraft":
		return fmt.Sprintf("a-%s-%s-C-F", affiliationCode, dimensionCode)
	case dimension == "ground" && category == "main_battle_tank":
		return fmt.Sprintf("a-%s-%s-U-C-I", affiliationCode, dimensionCode)
	case dimension == "ground" && category == "tactical_vehicle":
		return fmt.Sprintf("a-%s-%s-U-C-V", affiliationCode, dimensionCode)
	case dimension == "ground" && category == "commercial_vehicle":
		return fmt.Sprintf("a-%s-%s-U-C-V", affiliationCode, dimensionCode)
	case dimension == "sea" && category == "destroyer":
		return fmt.Sprintf("a-%s-%s-U-W-D", affiliationCode, dimensionCode)
	case dimension == "sea" && category == "cargo_vessel":
		return fmt.Sprintf("a-%s-%s-U-C-V", affiliationCode, dimensionCode)
	case dimension == "space":
		return fmt.Sprintf("a-%s-%s-U-S", affiliationCode, dimensionCode)
	default:
		return fmt.Sprintf("a-%s-%s-U", affiliationCode, dimensionCode)
	}
}

// CalculateCourse calculates course between two points
func CalculateCourse(lat1, lon1, lat2, lon2 float64) float64 {
	// Convert to radians
	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLon := (lon2 - lon1) * math.Pi / 180

	y := math.Sin(deltaLon) * math.Cos(lat2Rad)
	x := math.Cos(lat1Rad)*math.Sin(lat2Rad) - math.Sin(lat1Rad)*math.Cos(lat2Rad)*math.Cos(deltaLon)

	bearing := math.Atan2(y, x) * 180 / math.Pi

	// Normalize to 0-360 degrees
	if bearing < 0 {
		bearing += 360
	}

	return bearing
}

// SetStaleTime sets the stale time for generated CoT messages
func (g *CoTGenerator) SetStaleTime(duration time.Duration) {
	g.staleTime = duration
}
