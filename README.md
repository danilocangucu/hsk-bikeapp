# My Helsinki city bike app
Hello and welcome! I am Danilo and this is the README for my solution of [Solita's pre-assignment](https://github.com/solita/dev-academy-2023-exercise "Go to repo") for the Dev Academy 2023.

##### Table of contents
[1. Running the app locally](#running-the-app-locally)
[2. Data import](#data-import)
[3. Stations](#stations)


## Running the app locally
The project uses `SQLite3`, `wget` command, `Go` and `internet connection`. Make sure you have them all running smoothly in your computer. Clone this repo in your IDE.

## Data import
Run the script that will create the database:

> ./dbcreate.sh

The data from the CSV files will be downloaded and saved to hsk-city-bike-app.db. Through queries: 

 - Bike journeys that lasted for less than ten seconds and covered distances shorter than 10m has been removed;
 - How many journeys has started and ended from each station has been counted and;
 - Duplicated data has been removed.

## Stations

In this section of the page, all bike stations can be seen. A `pagination` of 10 stations per page has been implemented. Clicking on each station, you will access the following information from them:

- Station name;
- Station address;
- Total number of journeys starting from this station and;
- Total number of journeys ending at this station.
