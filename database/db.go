package database

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type Db struct {
	connection *sql.DB
}

type Station struct {
	FID          int
	ID           int
	Nimi         string
	Namn         string
	Name         string
	Osoite       string
	Adress       string
	Kaupunki     string
	Stad         string
	Operaattor   string
	Kapasiteet   int
	Latitude     float32
	Longitude    float32
	JourneysFrom int
	JourneysTo   int
}

type Journey struct {
	ID                   int
	Departure            string
	Return               string
	DepartureStationId   int
	DepartureStationName string
	ReturnStationId      int
	ReturnStationName    string
	CoveredDistanceM     float64
	DurationSec          int
}

type StationFilter struct {
	ID        int
	Nimi      string
	Namn      string
	Name      string
	Osoite    string
	Adress    string
	Latitude  float32
	Longitude float32
}

type JourneyFilter struct {
	BatchFromId int
	Limit       int
}

func OpenDatabase() (Db, error) {
	db, err := sql.Open("sqlite3", "./database/hsk-city-bike-app.db")
	if err != nil {
		return Db{}, err
	}

	return Db{connection: db}, nil
}

func (db *Db) CloseDatabase() {
	db.connection.Close()
}

func (db *Db) GetAllStations() (stations []Station, err error) {
	var station Station

	query := "SELECT FID, ID, Nimi, Namn, Name, Osoite, Adress, Kaupunki, Stad, Operaattor, Kapasiteet, x, y, JourneysFrom, JourneysTo FROM stations ORDER BY FID DESC"

	rows, err := db.connection.Query(query)
	if err != nil {
		return stations, err
	}

	for rows.Next() {
		err := rows.Scan(&station.FID, &station.ID, &station.Nimi, &station.Namn, &station.Name, &station.Osoite, &station.Adress, &station.Kaupunki, &station.Stad, &station.Operaattor, &station.Kapasiteet, &station.Latitude, &station.Longitude, &station.JourneysFrom, &station.JourneysTo)
		if err != nil {
			return stations, err
		}
		stations = append(stations, station)
	}

	defer rows.Close()
	return stations, err
}

func (db *Db) GetSingleStation(filter StationFilter) (station Station, err error) {
	query := "SELECT FID, ID, Nimi, Namn, Name, Osoite, Adress, Kaupunki, Stad, Operaattor, Kapasiteet, x, y, JourneysFrom, JourneysTo FROM stations WHERE "

	// Build the query based on the presence of filter values
	var args []interface{}
	if filter.ID != 0 {
		query += "ID = $1"
		args = append(args, filter.ID)
	} else {
		conditions := []string{}
		if filter.Nimi != "" {
			conditions = append(conditions, "Nimi = $"+strconv.Itoa(len(args)+1))
			args = append(args, filter.Nimi)
		}
		if filter.Namn != "" {
			conditions = append(conditions, "Namn = $"+strconv.Itoa(len(args)+1))
			args = append(args, filter.Namn)
		}
		if filter.Name != "" {
			conditions = append(conditions, "Name = $"+strconv.Itoa(len(args)+1))
			args = append(args, filter.Name)
		}
		if filter.Osoite != "" {
			conditions = append(conditions, "Osoite = $"+strconv.Itoa(len(args)+1))
			args = append(args, filter.Osoite)
		}
		if filter.Adress != "" {
			conditions = append(conditions, "Adress = $"+strconv.Itoa(len(args)+1))
			args = append(args, filter.Adress)
		}
		if filter.Latitude != 0 {
			conditions = append(conditions, "x = $"+strconv.Itoa(len(args)+1))
			args = append(args, filter.Latitude)
		}
		if filter.Longitude != 0 {
			conditions = append(conditions, "y = $"+strconv.Itoa(len(args)+1))
			args = append(args, filter.Longitude)
		}
		query += strings.Join(conditions, " AND ")
	}

	err = db.connection.QueryRow(query, args...).Scan(&station.FID, &station.ID, &station.Nimi, &station.Namn, &station.Name, &station.Osoite, &station.Adress, &station.Kaupunki, &station.Stad, &station.Operaattor, &station.Kapasiteet, &station.Latitude, &station.Longitude, &station.JourneysFrom, &station.JourneysTo)
	if err != nil {
		return station, err
	}

	return station, nil
}

func (db *Db) ValidateNewStation(newStation Station) (errors []error) {
	// Check that the station ID is unique
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if _, err := db.GetSingleStation(StationFilter{ID: newStation.ID}); err == nil {
			errors = append(errors, fmt.Errorf("Station with ID %d already exists", newStation.ID))
		}
	}()

	// Check that the station names are unique
	wg.Add(3)
	go func() {
		defer wg.Done()
		if _, err := db.GetSingleStation(StationFilter{Nimi: newStation.Nimi}); err == nil {
			errors = append(errors, fmt.Errorf("Station with name '%s' already exists", newStation.Nimi))
		}
	}()
	go func() {
		defer wg.Done()
		if _, err := db.GetSingleStation(StationFilter{Namn: newStation.Namn}); err == nil {
			errors = append(errors, fmt.Errorf("Station with name '%s' already exists", newStation.Namn))
		}
	}()
	go func() {
		defer wg.Done()
		if _, err := db.GetSingleStation(StationFilter{Name: newStation.Name}); err == nil {
			errors = append(errors, fmt.Errorf("Station with name '%s' already exists", newStation.Name))
		}
	}()

	// Check that the station address is unique
	wg.Add(2)
	go func() {
		defer wg.Done()
		if _, err := db.GetSingleStation(StationFilter{Osoite: newStation.Osoite}); err == nil {
			errors = append(errors, fmt.Errorf("Station with address '%s' already exists", newStation.Osoite))
		}
	}()
	go func() {
		defer wg.Done()
		if _, err := db.GetSingleStation(StationFilter{Adress: newStation.Adress}); err == nil {
			errors = append(errors, fmt.Errorf("Station with address '%s' already exists", newStation.Adress))
		}
	}()

	// Check that the station coordinates are unique
	wg.Add(1)
	go func() {
		defer wg.Done()
		if _, err := db.GetSingleStation(StationFilter{Latitude: newStation.Latitude, Longitude: newStation.Longitude}); err == nil {
			errors = append(errors, fmt.Errorf("Station with coordinates (%f, %f) already exists", newStation.Latitude, newStation.Longitude))
		}
	}()

	wg.Wait()

	if len(errors) > 0 {
		return errors
	}
	return errors
}

func (db *Db) AddNewStation(newStation Station) error {
	var lastFid int
	var lastId int

	err := db.connection.QueryRow("SELECT MAX(FID) FROM stations").Scan(&lastFid)
	if err != nil {
		return err
	}

	err = db.connection.QueryRow("SELECT MAX(ID) FROM stations").Scan(&lastId)
	if err != nil {
		return err
	}

	newStation.FID = lastFid + 1
	newStation.ID = lastId + 1

	query := `INSERT INTO stations (FID, ID, Nimi, Namn, Name, Osoite, Adress, Kaupunki, Stad, Operaattor, Kapasiteet, x, y, JourneysFrom, JourneysTo) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`

	_, err = db.connection.Exec(query, newStation.FID, newStation.ID, newStation.Nimi, newStation.Namn, newStation.Name, newStation.Osoite, newStation.Adress, newStation.Kaupunki, newStation.Stad, newStation.Operaattor, newStation.Kapasiteet, newStation.Longitude, newStation.Latitude, newStation.JourneysFrom, newStation.JourneysTo)

	return err
}

func (db *Db) GetLastJourneyId() (lastJourneyId int, err error) {

	row, err := db.connection.Query("SELECT MAX(id) FROM all_journeys;")
	if err != nil {
		return lastJourneyId, err
	}

	for row.Next() {
		err := row.Scan(&lastJourneyId)
		if err != nil {
			return lastJourneyId, err
		}
	}

	defer row.Close()

	return lastJourneyId, err
}

func (db *Db) GetJourneys(filter JourneyFilter) (journeys []Journey, err error) {
	var journey Journey

	query := fmt.Sprintf("SELECT * FROM 'all_journeys' WHERE id >= %v ORDER BY id LIMIT %v", filter.BatchFromId, filter.Limit)
	rows, err := db.connection.Query(query)

	if err != nil {
		return journeys, err
	}

	for rows.Next() {
		err := rows.Scan(&journey.ID, &journey.Departure, &journey.Return, &journey.DepartureStationId, &journey.DepartureStationName, &journey.ReturnStationId, &journey.ReturnStationName, &journey.CoveredDistanceM, &journey.DurationSec)
		if err != nil {
			return journeys, err
		}
		journeys = append(journeys, journey)
	}

	defer rows.Close()
	return journeys, err
}
