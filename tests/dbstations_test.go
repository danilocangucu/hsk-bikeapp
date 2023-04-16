package tests

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestImportData1(t *testing.T) {
	// Open a connection to the database
	db, err := sql.Open("sqlite3", "../database/hsk-city-bike-app.db")
	if err != nil {
		t.Fatalf("failed to open database connection: %v", err)
	}
	defer db.Close()

	// Import the data from each CSV file and verify it was imported correctly
	for _, filename := range []string{"../datasets/726277c507ef4914b0aec3cbcfcbfafc_0.csv"} {
		// Open the CSV file
		file, err := os.Open(filename)
		if err != nil {
			t.Fatalf("failed to open CSV file %q: %v", filename, err)
		}
		defer file.Close()

		// Parse the CSV file
		r := csv.NewReader(file)

		// Ignore the first row
		if _, err := r.Read(); err != nil {
			t.Fatalf("failed to read CSV file %q: %v", filename, err)
		}

		records, err := r.ReadAll()
		if err != nil {
			t.Fatalf("failed to parse CSV file %q: %v", filename, err)
		}

		// Create a table for the data
		tablename := fmt.Sprintf("test_%s", strings.TrimSuffix(path.Base(filename), path.Ext(filename)))
		fmt.Println(tablename)

		// Drop the table if it exists
        _, err = db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", tablename))
        if err != nil {
            t.Fatalf("failed to drop table for CSV file %q: %v", filename, err)
        }

		_, err = db.Exec(fmt.Sprintf("CREATE TABLE %s (FID INT, ID INT, Nimi TEXT, Namn TEXT, Name TEXT, Osoite TEXT, Adress TEXT, Kaupunki TEXT, Stad TEXT, Operaattor TEXT, Kapasiteet INT, x REAL, y REAL)", tablename))
		if err != nil {
			t.Fatalf("failed to create table for CSV file %q: %v", filename, err)
		}

		// Insert the data into the table
		tx, err := db.Begin()
		if err != nil {
			t.Fatalf("failed to start transaction for CSV file %q: %v", filename, err)
		}
		stmt, err := tx.Prepare(fmt.Sprintf("INSERT INTO %s (FID, ID, Nimi, Namn, Name, Osoite, Adress, Kaupunki, Stad, Operaattor, Kapasiteet, x, y) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", tablename))
		if err != nil {
			t.Fatalf("failed to prepare statement for CSV file %q: %v", filename, err)
		}
		defer stmt.Close()
		for _, record := range records {
			_, err = stmt.Exec(record[0], record[1], record[2], record[3], record[4], record[5], record[6], record[7], record[8], record[9], record[10], record[11], record[12])
			if err != nil {
				t.Fatalf("failed to insert record into table for CSV file %q: %v", filename, err)
			}
		}
		err = tx.Commit()
		if err != nil {
			t.Fatalf("failed to commit transaction for CSV file %q: %v", filename, err)
		}

		// Verify the data was inserted correctly
		var count int
		err = db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", tablename)).Scan(&count)
		if err != nil {
			t.Fatalf("failed to query table for CSV file %q: %v", filename, err)
		}
		if count != len(records) {
			t.Fatalf("unexpected number of records in table for CSV file %q: got %d, want %d", filename, count, len(records))
		}

		var stationCount int
		err = db.QueryRow("SELECT COUNT(*) FROM stations").Scan(&stationCount)
		if err != nil {
			t.Fatalf("failed to query stations table: %v", err)
		}
		var importedCount int
		err = db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE ID NOT IN (SELECT id FROM stations) OR Name NOT IN (SELECT name FROM stations)", tablename)).Scan(&importedCount)
		if err != nil {
			t.Fatalf("failed to query imported data for CSV file %q: %v", filename, err)
		}
		if importedCount != 0 {
			t.Fatalf("imported data for CSV file %q does not match stations table: %d rows imported with invalid station IDs", filename, importedCount)
		}
		if count != stationCount {
			t.Fatalf("imported data for CSV file %q does not match stations table: %d rows imported, %d rows in stations table", filename, count, stationCount)
		}

		_, err = db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", tablename))
        if err != nil {
            t.Fatalf("failed to drop table for CSV file %q: %v", filename, err)
        }
	}
}