package database

import (
	"database/sql"
	"fmt"

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
	StationId int
}

type JourneyFilter struct {
	LastId int
	Limit  int
}

func OpenDatabase() (Db, error) {
	db, err := sql.Open("sqlite3", "./database/hsk-city-bike-app.db")
	if err != nil {
		return Db{}, err
	}

	return Db{connection: db}, nil
}

func (db *Db) Close() {
	db.connection.Close()
}

func (db *Db) GetStations(filter StationFilter) (stations []Station, err error) {
	var station Station
	var query string

	if filter != (StationFilter{}) {
		query = fmt.Sprintf("SELECT FID, ID, Nimi, Namn, Name, Osoite, Adress, Kaupunki, Stad, Operaattor, Kapasiteet, x, y, JourneysFrom, JourneysTo FROM stations WHERE ID = %v", filter.StationId)

	} else {
		query = "select FID,ID,Nimi,Namn,Name,Osoite,Adress,Kaupunki,Stad,Operaattor,Kapasiteet,x,y,JourneysFrom,JourneysTo from stations"
	}

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

func (db *Db) GetLastJourneyId() (lastJourney JourneyFilter, err error) {

	row, err := db.connection.Query("SELECT MAX(id) FROM all_journeys;")
	if err != nil {
		return lastJourney, err
	}

	for row.Next() {
		err := row.Scan(&lastJourney.LastId)
		if err != nil {
			return lastJourney, err
		}
	}

	defer row.Close()

	return lastJourney, err
}

func (db *Db) GetJourneys(filter JourneyFilter) (journeys []Journey, err error) {
	var journey Journey

	query := fmt.Sprintf("SELECT * FROM 'all_journeys' WHERE id > %v ORDER BY id LIMIT %v", filter.LastId, filter.Limit)
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
