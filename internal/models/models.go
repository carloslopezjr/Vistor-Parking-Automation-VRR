package models

// Vehicle represents a visitor parking entry loaded from CSV.
type Vehicle struct {
	ID                    int
	ApartmentName         string
	LicensePlate          string
	ConfirmedLicensePlate string
	VehicleMake           string
	VehicleModel          string
	ResidentName          string
	UnitNumber            string
	VisitorName           string
	ConfirmationEmail     string
}

// AutomationResult represents the outcome of a registration automation attempt.
type AutomationResult struct {
	Success bool
	Error   string
	Logs    []string
}
