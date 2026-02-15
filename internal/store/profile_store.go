package store

import (
	"context"
	"database/sql"
	"time"

	"vistor-parking-automation-vrr/internal/models"
)

// ProfileStore provides CRUD access to profiles.
type ProfileStore struct {
	db *sql.DB
}

func NewProfileStore(db *sql.DB) *ProfileStore {
	return &ProfileStore{db: db}
}

func (s *ProfileStore) List(ctx context.Context) ([]models.Profile, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, apartment_name, license_plate, vehicle_make, vehicle_model,
		 resident_name, unit_number, visitor_name, resident_email, profile_name,
		 created_at, updated_at, last_registration_at
		 FROM profiles
		 ORDER BY created_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []models.Profile
	for rows.Next() {
		var p models.Profile
		var lastReg sql.NullTime
		if err := rows.Scan(&p.ID, &p.ApartmentName, &p.LicensePlate, &p.VehicleMake, &p.VehicleModel,
			&p.ResidentName, &p.UnitNumber, &p.VisitorName, &p.ResidentEmail, &p.ProfileName,
			&p.CreatedAt, &p.UpdatedAt, &lastReg); err != nil {
			return nil, err
		}
		if lastReg.Valid {
			p.LastRegistrationAt = &lastReg.Time
		}
		res = append(res, p)
	}
	return res, rows.Err()
}

func (s *ProfileStore) Get(ctx context.Context, id int64) (*models.Profile, error) {
	var p models.Profile
	var lastReg sql.NullTime
	err := s.db.QueryRowContext(ctx,
		`SELECT id, apartment_name, license_plate, vehicle_make, vehicle_model,
		 resident_name, unit_number, visitor_name, resident_email, profile_name,
		 created_at, updated_at, last_registration_at
		 FROM profiles WHERE id = ?`, id,
	).Scan(&p.ID, &p.ApartmentName, &p.LicensePlate, &p.VehicleMake, &p.VehicleModel,
		&p.ResidentName, &p.UnitNumber, &p.VisitorName, &p.ResidentEmail, &p.ProfileName,
		&p.CreatedAt, &p.UpdatedAt, &lastReg)
	if err != nil {
		return nil, err
	}
	if lastReg.Valid {
		p.LastRegistrationAt = &lastReg.Time
	}
	return &p, nil
}

// Create inserts a new profile.
func (s *ProfileStore) Create(ctx context.Context, p *models.Profile) (int64, error) {
	now := time.Now().UTC()
	res, err := s.db.ExecContext(ctx,
		`INSERT INTO profiles (
			apartment_name, license_plate, vehicle_make, vehicle_model,
			resident_name, unit_number, visitor_name, resident_email, profile_name,
			created_at, updated_at, last_registration_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NULL)`,
		p.ApartmentName, p.LicensePlate, p.VehicleMake, p.VehicleModel,
		p.ResidentName, p.UnitNumber, p.VisitorName, p.ResidentEmail, p.ProfileName,
		now, now,
	)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// Update updates an existing profile.
func (s *ProfileStore) Update(ctx context.Context, p *models.Profile) error {
	now := time.Now().UTC()
	_, err := s.db.ExecContext(ctx,
		`UPDATE profiles SET
			apartment_name = ?,
			license_plate = ?,
			vehicle_make = ?,
			vehicle_model = ?,
			resident_name = ?,
			unit_number = ?,
			visitor_name = ?,
			resident_email = ?,
			profile_name = ?,
			updated_at = ?
		 WHERE id = ?`,
		p.ApartmentName, p.LicensePlate, p.VehicleMake, p.VehicleModel,
		p.ResidentName, p.UnitNumber, p.VisitorName, p.ResidentEmail, p.ProfileName,
		now, p.ID,
	)
	return err
}

// Delete removes a profile.
func (s *ProfileStore) Delete(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM profiles WHERE id = ?`, id)
	return err
}
