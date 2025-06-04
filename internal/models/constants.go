package models

// Platform type constants
const (
	// Aircraft types
	AircraftTypeF16       = "F-16 Fighting Falcon"
	AircraftTypeF22       = "F-22 Raptor"
	AircraftTypeBoeing777 = "Boeing 777"
	AircraftTypeBoeing737 = "Boeing 737"
	AircraftTypeA320      = "Airbus A320"
	AircraftTypeA380      = "Airbus A380"
	AircraftTypeC130      = "C-130 Hercules"

	// Vehicle types
	VehicleTypeAbrams       = "M1A2 Abrams"
	VehicleTypeBradley      = "M2 Bradley"
	VehicleTypeHumvee       = "HMMWV"
	VehicleTypeTeslaS       = "Tesla Model S"
	VehicleTypeCamry        = "Toyota Camry"
	VehicleTypeF150         = "Ford F-150"
	VehicleTypeFreightliner = "Freightliner Cascadia"

	// Maritime types
	MaritimeTypeArleighBurke = "Arleigh Burke-class"
	MaritimeTypeTiconderoga  = "Ticonderoga-class"
	MaritimeTypeContainer    = "Container Ship"

	// Space types
	SpaceTypeISS       = "ISS Module"
	SpaceTypeStarlink  = "Starlink Satellite"
	SpaceTypeGPS       = "GPS Block III"
	SpaceTypeTelescope = "Space Telescope"
	SpaceTypeDragon    = "Dragon 2 Capsule"

	// Propulsion types
	PropulsionJet        = "Jet Engine"
	PropulsionProp       = "Propeller"
	PropulsionDiesel     = "Diesel"
	PropulsionGasoline   = "Gasoline"
	PropulsionElectric   = "Electric"
	PropulsionGasTurbine = "Gas Turbine"
	PropulsionNuclear    = "Nuclear"

	// Mission types
	MissionTypeCombat        = "combat"
	MissionTypeTransport     = "transport"
	MissionTypeCommunication = "communication"
	MissionTypeNavigation    = "navigation"
	MissionTypeObservation   = "observation"
	MissionTypeCrewTransport = "crew transport"
	MissionTypeSpaceStation  = "space station"

	// Earth constants for orbital mechanics
	EarthRadius      = 6371000.0 // meters
	EarthGM          = 3.986e14  // m³/s² (gravitational parameter)
	SeaLevel         = 0.0       // meters
	ISSAltitude      = 408000    // meters
	StarlinkAltitude = 550000    // meters
	GPSAltitude      = 20200000  // meters
	HubbleAltitude   = 547000    // meters

	// Common call sign prefixes
	CallSignNavy     = "NAVY"
	CallSignAir      = "AIR"
	CallSignGround   = "GROUND"
	CallSignISS      = "ISS"
	CallSignStarlink = "STARLINK"
	CallSignGPS      = "GPS"
	CallSignHubble   = "HUBBLE"
	CallSignDragon   = "DRAGON"
)

// Weapon system constants
var (
	AegisWeapons       = []string{"Aegis Combat System", "VLS Missiles", "5-inch Gun"}
	TiconderogaWeapons = []string{"Aegis Combat System", "VLS Missiles", "5-inch Gun", "Tomahawk Missiles"}
	F16Weapons         = []string{"M61 Vulcan Cannon", "AIM-120 AMRAAM", "AIM-9 Sidewinder"}
	F22Weapons         = []string{"M61A2 Cannon", "AIM-120C AMRAAM", "AIM-9X Sidewinder", "GBU-32 JDAM"}
	AbramsSystems      = []string{"M256 120mm Gun", "M240 7.62mm MG", "M2 .50 cal MG"}
	BradleySystems     = []string{"M242 25mm Cannon", "TOW Missiles", "M240C 7.62mm MG"}
)
