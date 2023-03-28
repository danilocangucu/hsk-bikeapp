package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Db struct {
	connection *sql.DB
}

type Station struct {
	FID string
	ID string
	Nimi string
	Namn string
	Name string
	Osoite string
	Adress string
	Kaupunki string
	Stad string
	Operaattor string
	Kapasiteet string
	x string
	y string
	JourneysFrom string
	JourneysTo string
}

type StationFilter struct {
	StationId int
}

func OpenDatabase() Db {
	db, err := sql.Open("sqlite3", "./database/hsk-city-bike-app.db")
	if err != nil{
		fmt.Fprintln(os.Stderr, err)
		//todo error?
	}
	fmt.Println("db opened")

	return Db{connection: db}
}

func (db *Db) Close() {
	db.connection.Close()
}

func (db *Db) GetStations(filter StationFilter) (stations []Station, err error) {
	var station Station
	var query string

	if filter != (StationFilter{}){
		query = fmt.Sprintf("SELECT FID, ID, Nimi, Namn, Name, Osoite, Adress, Kaupunki, Stad, Operaattor, Kapasiteet, JourneysFrom, JourneysTo FROM stations WHERE ID = %v", filter.StationId)

	} else {
		query = "select FID,ID,Nimi,Namn,Name,Osoite,Adress,Kaupunki,Stad,Operaattor,Kapasiteet,JourneysFrom,JourneysTo from stations"
	}

	rows, err := db.connection.Query(query)
	if err != nil {
		fmt.Println(1)
		return stations, err
	}

	for rows.Next() {
		err := rows.Scan(&station.FID, &station.ID, &station.Nimi, &station.Namn, &station.Name, &station.Osoite, &station.Adress, &station.Kaupunki, &station.Stad, &station.Operaattor, &station.Kapasiteet, &station.JourneysFrom, &station.JourneysTo)
		if err != nil {
			fmt.Println(2)
			return stations, err
		}
		stations = append(stations, station)
	}
	defer rows.Close()
	return stations, err
}

