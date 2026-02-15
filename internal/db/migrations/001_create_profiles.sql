-- 001_create_profiles.sql

CREATE TABLE IF NOT EXISTS profiles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    apartment_name TEXT NOT NULL,
    license_plate TEXT NOT NULL,
    vehicle_make TEXT NOT NULL,
    vehicle_model TEXT NOT NULL,
    resident_name TEXT NOT NULL,
    unit_number TEXT NOT NULL,
    visitor_name TEXT NOT NULL,
    resident_email TEXT NOT NULL,
    profile_name TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    last_registration_at DATETIME
);

CREATE INDEX IF NOT EXISTS idx_profiles_email ON profiles (resident_email);
CREATE INDEX IF NOT EXISTS idx_profiles_name ON profiles (profile_name);
