#!/bin/bash

# Change the file paths and names to match your own
DATASET_1="2021-05.csv"
DATASET_2="2021-06.csv"
DATASET_3="2021-07.csv"
DATASET_4="726277c507ef4914b0aec3cbcfcbfafc_0.csv"
DB_NAME="hsk-city-bike-app.db"

# Download the datasets
echo "Downloading $DATASET_1 ..."
wget -q https://dev.hsl.fi/citybikes/od-trips-2021/$DATASET_1
echo "Downloading $DATASET_2 ..."
wget -q https://dev.hsl.fi/citybikes/od-trips-2021/$DATASET_2
echo "Downloading $DATASET_3 ..."
wget -q https://dev.hsl.fi/citybikes/od-trips-2021/$DATASET_3
echo "Downloading $DATASET_4 ..."
wget -q https://opendata.arcgis.com/datasets/$DATASET_4

# Run the SQLite3 commands
echo "Creating database ..."
sqlite3 $DB_NAME << EOF
.mode csv

CREATE TABLE raw_journeys ("Departure","Return","Departure_station_id","Departure_station_name","Return_station_id","Return_station_name","Covered_distance_m","Duration_sec");
.import $DATASET_1 raw_journeys
CREATE TABLE journeys AS SELECT DISTINCT "Departure","Return","Departure_station_id","Departure_station_name","Return_station_id","Return_station_name","Covered_distance_m","Duration_sec" FROM raw_journeys WHERE "Duration_sec" >= 10 AND "Covered_distance_m" >= 10 ORDER BY "Departure" ASC;
DROP TABLE raw_journeys;
.import $DATASET_2 raw_journeys
INSERT INTO journeys SELECT DISTINCT "Departure","Return","Departure_station_id","Departure_station_name","Return_station_id","Return_station_name","Covered_distance_m","Duration_sec" FROM raw_journeys WHERE "Duration_sec" >= 10 AND "Covered_distance_m" >= 10 ORDER BY "Departure" ASC;
DROP TABLE raw_journeys;
.import $DATASET_3 raw_journeys
INSERT INTO journeys SELECT DISTINCT "Departure","Return","Departure_station_id","Departure_station_name","Return_station_id","Return_station_name","Covered_distance_m","Duration_sec" FROM raw_journeys WHERE "Duration_sec" >= 10 AND "Covered_distance_m" >= 10 ORDER BY "Departure" ASC;
DROP TABLE raw_journeys;
CREATE TABLE stations_raw ("FID","ID","Nimi","Namn","Name","Osoite","Adress","Kaupunki","Stad","Operaattor","Kapasiteet","x","y");
.import $DATASET_4 raw_stations
CREATE TABLE stations AS SELECT "FID","ID","Nimi","Namn","Name","Osoite","Adress","Kaupunki","Stad","Operaattor","Kapasiteet","x","y" FROM stations_raw ORDER BY "Name" ASC;
DROP TABLE stations_raw;
.save $DB_NAME
.exit
EOF

sleep 1 # Wait for the database to finish writing

echo "Data imported and saved to $DB_NAME"

rm $DATASET_1 $DATASET_2 $DATASET_3 $DATASET_4
echo "$DATASET_1 $DATASET_2 $DATASET_3 $DATASET_4 removed"
