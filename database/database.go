package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// InitDB initializes the database connection and creates necessary tables
func InitDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if err = createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return db, nil
}

// createTables creates the necessary tables if they don't exist
func createTables(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS ships (
			id SERIAL PRIMARY KEY,
			mmsi VARCHAR(255) UNIQUE NOT NULL,
			name VARCHAR(255),
			last_updated TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS ship_data (
			id SERIAL PRIMARY KEY,
			ship_id INT REFERENCES ships(id),
			latitude FLOAT,
			longitude FLOAT,
			speed FLOAT,
			heading FLOAT,
			timestamp TIMESTAMP
		);
	`)
	return err
}

// StoreShipData stores or updates ship data in the database
func StoreShipData(db *sql.DB, mmsi string, name string, lat, lon, speed, heading float64, timestamp time.Time) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var shipID int
	err = tx.QueryRow(`
		INSERT INTO ships (mmsi, name, last_updated)
		VALUES ($1, $2, $3)
		ON CONFLICT (mmsi) DO UPDATE SET name = $2, last_updated = $3
		RETURNING id
	`, mmsi, name, timestamp).Scan(&shipID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO ship_data (ship_id, latitude, longitude, speed, heading, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, shipID, lat, lon, speed, heading, timestamp)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetShipData retrieves the latest ship data for a given MMSI
func GetShipData(db *sql.DB, mmsi string) ([]map[string]interface{}, error) {
	rows, err := db.Query(`
		SELECT s.mmsi, s.name, sd.latitude, sd.longitude, sd.speed, sd.heading, sd.timestamp
		FROM ships s
		JOIN ship_data sd ON s.id = sd.ship_id
		WHERE s.mmsi = $1
		ORDER BY sd.timestamp DESC
		LIMIT 100
	`, mmsi)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var mmsi, name string
		var lat, lon, speed, heading float64
		var timestamp time.Time
		err := rows.Scan(&mmsi, &name, &lat, &lon, &speed, &heading, &timestamp)
		if err != nil {
			return nil, err
		}
		results = append(results, map[string]interface{}{
			"mmsi":      mmsi,
			"name":      name,
			"latitude":  lat,
			"longitude": lon,
			"speed":     speed,
			"heading":   heading,
			"timestamp": timestamp,
		})
	}
	return results, nil
}