#!/bin/bash

# Define the file paths and names
DATASET_1="2021-05.csv"
DATASET_2="2021-06.csv"
DATASET_3="2021-07.csv"
DATASET_4="726277c507ef4914b0aec3cbcfcbfafc_0.csv"
DB_NAME="database/hsk-city-bike-app.db"

# Delete CSV files if they exist
if [ -f "$DATASET_1" ]; then
  rm $DATASET_1
fi
if [ -f "$DATASET_2" ]; then
  rm $DATASET_2
fi
if [ -f "$DATASET_3" ]; then
  rm $DATASET_3
fi
if [ -f "$DATASET_4" ]; then
  rm $DATASET_4
fi

# Download the datasets
echo "Downloading $DATASET_1 ..."
wget -q https://dev.hsl.fi/citybikes/od-trips-2021/$DATASET_1
echo "Downloading $DATASET_2 ..."
wget -q https://dev.hsl.fi/citybikes/od-trips-2021/$DATASET_2
echo "Downloading $DATASET_3 ..."
wget -q https://dev.hsl.fi/citybikes/od-trips-2021/$DATASET_3
echo "Downloading $DATASET_4 ..."
wget -q https://opendata.arcgis.com/datasets/$DATASET_4

# Create the database and import data
echo "Creating database ..."
sqlite3 $DB_NAME << EOF
PRAGMA journal_mode = MEMORY;
.mode csv

DROP TABLE IF EXISTS journeys;
DROP TABLE IF EXISTS stations;
DROP TABLE IF EXISTS raw_stations;

CREATE TABLE raw_journeys (
  "Departure",
  "Return",
  "Departure_station_id" INTEGER,
  "Departure_station_name",
  "Return_station_id" INTEGER,
  "Return_station_name",
  "Covered_distance_m" INTEGER,
  "Duration_sec" INTEGER
);

.import $DATASET_1 raw_journeys
CREATE TABLE journeys AS SELECT DISTINCT
  "Departure",
  "Return",
  "Departure_station_id",
  "Departure_station_name",
  "Return_station_id",
  "Return_station_name",
  "Covered_distance_m",
  "Duration_sec"
FROM raw_journeys
WHERE "Duration_sec" >= 10 AND "Covered_distance_m" >= 10
ORDER BY "Departure" ASC;

DROP TABLE raw_journeys;
.import $DATASET_2 raw_journeys
INSERT INTO journeys SELECT DISTINCT
  "Departure",
  "Return",
  "Departure_station_id",
  "Departure_station_name",
  "Return_station_id",
  "Return_station_name",
  "Covered_distance_m",
  "Duration_sec"
FROM raw_journeys
WHERE "Duration_sec" >= 10 AND "Covered_distance_m" >= 10
ORDER BY "Departure" ASC;

DROP TABLE raw_journeys;
.import $DATASET_3 raw_journeys
INSERT INTO journeys SELECT DISTINCT
  "Departure",
  "Return",
  "Departure_station_id",
  "Departure_station_name",
  "Return_station_id",
  "Return_station_name",
  "Covered_distance_m",
  "Duration_sec"
FROM raw_journeys
WHERE "Duration_sec" >= 10 AND "Covered_distance_m" >= 10
ORDER BY "Departure" ASC;

DROP TABLE raw_journeys;
CREATE TABLE stations (
  "FID" INTEGER,
  "ID" INTEGER,
  "Nimi",
  "Namn",
  "Name",
  "Osoite",
  "Adress",
  "Kaupunki",
  "Stad",
  "Operaattor",
  "Kapasiteet" INTEGER,
  "x",
  "y"
);
.import $DATASET_4 stations
ALTER TABLE stations ADD COLUMN JourneysFrom INTEGER;
ALTER TABLE stations ADD COLUMN JourneysTo INTEGER;
UPDATE stations
SET JourneysFrom = (
    SELECT COUNT(*)
    FROM journeys
    WHERE journeys.Departure_station_id = stations.ID
),
JourneysTo = (
    SELECT COUNT(*)
    FROM journeys
    WHERE journeys.Return_station_id = stations.ID
)
WHERE EXISTS (
    SELECT 1
    FROM journeys
    WHERE journeys.Departure_station_id = stations.ID
    OR journeys.Return_station_id = stations.ID
);
DELETE FROM stations WHERE ROWID = 1;
.save $DB_NAME
.sleep 2
.exit
EOF

echo "Data imported and saved to $DB_NAME"

rm $DATASET_1 $DATASET_2 $DATASET_3 $DATASET_4
echo "$DATASET_1, $DATASET_2, $DATASET_3 and $DATASET_4 removed"
