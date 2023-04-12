package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestImportData(t *testing.T) {
	// Open a connection to the database
	db, err := sql.Open("sqlite3", "./database/hsk-city-bike-app.db")
	if err != nil {
		t.Fatalf("failed to open database connection: %v", err)
	}
	defer db.Close()

	// Import the data from each CSV file and verify it was imported correctly
	for _, filedata := range []struct {
		filename   string
		tablename  string
		rawjourney string
	}{
		{"2021-05.csv", "test_2021_05", "raw_journeys202105"},
		{"2021-06.csv", "test_2021_06", "raw_journeys202106"},
		{"2021-07.csv", "test_2021_07", "raw_journeys202107"},
	} {
		// Open the CSV file
		file, err := os.Open(filedata.filename)
		if err != nil {
			t.Fatalf("failed to open CSV file %q: %v", filedata.filename, err)
		}
		defer file.Close()

		// Parse the CSV file
		r := csv.NewReader(file)
		records, err := r.ReadAll()
		if err != nil {
			t.Fatalf("failed to parse CSV file %q: %v", filedata.filename, err)
		}

		// Drop the table if it exists
		 _, err = db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", filedata.tablename))
		if err != nil {
			t.Fatalf("failed to drop table for CSV file %q: %v", filedata.filename, err)
		}

		// Create a table for the data
		_, err = db.Exec(fmt.Sprintf("CREATE TABLE `%s` (Departure TEXT, Return TEXT, `Departure station id` INTEGER, `Departure station name` TEXT, `Return station id` INTEGER, `Return station name` TEXT, `Covered distance (m)` INTEGER, `Duration (sec.)` INTEGER)", filedata.tablename))
		if err != nil {
			t.Fatalf("failed to create table for CSV file %q: %v", filedata.filename, err)
		}

        // Insert the data into the table
        tx, err := db.Begin()
        if err != nil {
            t.Fatalf("failed to start transaction for CSV file %q: %v", filedata.filename, err)
        }
        stmt, err := tx.Prepare(fmt.Sprintf("INSERT INTO %s (Departure, Return, `Departure station id`, `Departure station name`, `Return station id`, `Return station name`, `Covered distance (m)`, `Duration (sec.)`) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", filedata.tablename))
        if err != nil {
            t.Fatalf("failed to prepare statement for CSV file %q: %v", filedata.filename, err)
        }
        defer stmt.Close()

        for _, record := range records {
            _, err = stmt.Exec(record[0], record[1], record[2], record[3], record[4], record[5], record[6], record[7])
            if err != nil {
                t.Fatalf("failed to insert record into table for CSV file %q: %v", filedata.filename, err)
            }
        }

        err = tx.Commit()
        if err != nil {
            t.Fatalf("failed to commit transaction for CSV file %q: %v", filedata.filename, err)
        }

        // Verify the imported data
        var importedCount int
        err = db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE `Departure station id` NOT IN (SELECT `Departure station id` FROM %s) OR `Return station id` NOT IN (SELECT `Return station id` FROM %s)", filedata.tablename, filedata.rawjourney, filedata.rawjourney)).Scan(&importedCount)
        if err != nil {
            t.Fatalf("failed to query imported data for CSV file %q: %v", filedata.filename, err)
        }
        if importedCount != 0 {
            t.Errorf("failed to import all data for CSV file %q: %d rows not imported", filedata.filename, importedCount)
        }

        // Drop the rawjourneys table
        _, err = db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", filedata.rawjourney))
        if err != nil {
            t.Fatalf("failed to drop rawjourneys table for CSV file %q: %v", filedata.filename, err)
        }

        // Drop the test table
        _, err = db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", filedata.tablename))
        if err != nil {
            t.Fatalf("failed to drop table for CSV file %q: %v", filedata.filename, err)
        }
    }
}