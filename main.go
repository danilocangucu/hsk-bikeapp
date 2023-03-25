package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func OpenDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "./database/hsk-city-bike.db")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("database opened")
	return db
}

func main() {
	OpenDatabase()
}
