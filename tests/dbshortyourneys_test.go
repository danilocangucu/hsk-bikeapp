package tests

import (
	"database/sql"
	"testing"
)

func TestShortJourneys(t *testing.T) {
	// Open a connection to the database
	db, err := sql.Open("sqlite3", "../database/hsk-city-bike-app.db")
	if err != nil {
		t.Fatalf("failed to open database connection: %v", err)
	}
	defer db.Close()

	// Query the "all_journeys" table
	rows, err := db.Query("SELECT * FROM all_journeys WHERE `Covered distance (m)` < 10 OR `Duration (sec.)` < 10")
	if err != nil {
		t.Fatalf("failed to query all_journeys table: %v", err)
	}
	defer rows.Close()

	// Check if any rows match the condition
	if rows.Next() {
		t.Errorf("found a row in all_journeys table with Covered distance (m) or Duration (sec.) less than 10")
	}
}
