package csv

import (
	"encoding/csv"
	"fmt"
	"os"

	"vistor-parking-automation-vrr/internal/models"
)

// LoadVehicles reads a CSV file and returns a slice of Vehicle records.
// CSV columns expected (in order):
// 0: ApartmentName
// 1: LicensePlate
// 2: ConfirmedLicensePlate
// 3: VehicleMake
// 4: VehicleModel
// 5: ResidentName
// 6: UnitNumber
// 7: VisitorName
// 8: ConfirmationEmail
func LoadVehicles(filePath string) ([]models.Vehicle, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV file: %w", err)
	}

	var vehicles []models.Vehicle
	for i, record := range records {
		// Skip header row
		if i == 0 {
			continue
		}

		// Skip incomplete rows
		if len(record) < 9 {
			continue
		}

		vehicle := models.Vehicle{
			ID:                    i - 1, // Adjust index to exclude header
			ApartmentName:         record[0],
			LicensePlate:          record[1],
			ConfirmedLicensePlate: record[2],
			VehicleMake:           record[3],
			VehicleModel:          record[4],
			ResidentName:          record[5],
			UnitNumber:            record[6],
			VisitorName:           record[7],
			ConfirmationEmail:     record[8],
		}
		vehicles = append(vehicles, vehicle)
	}

	if len(vehicles) == 0 {
		return nil, fmt.Errorf("no valid vehicles found in CSV file")
	}

	return vehicles, nil
}
