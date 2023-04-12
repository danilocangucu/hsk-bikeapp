# My Helsinki city bike app
Hello and welcome! I am Danilo and this is the README for my solution of [Solita's pre-assignment](https://github.com/solita/dev-academy-2023-exercise "Go to repo") for the Dev Academy 2023.

##### Table of contents
[1. Running the app locally](#running-the-app-locally)  
[2. Data import](#data-import)  
[3. Server & application](#starting-the-server)  
[4. Stations](#stations)  
[5. Tests](#tests)


## Running the app locally
The project uses  `Go`, `JavaScript`, `SQLite3`, `wget` command, `Mapbox GL JS library`, `Font Awesome library` and `internet connection`. Make sure you have them all running smoothly in your computer. Clone this repo in your IDE.

## Data import
Run the script from the root directory, it will create the database for the application:

> ./database/dbcreate.sh

The data from the CSV files will be downloaded and saved to hsk-city-bike-app.db. Through queries: 

 - Bike journeys that lasted for less than ten seconds and covered distances shorter than 10m has been removed;
 - How many journeys has started and ended from each station has been counted and;
 - Duplicated data has been removed.

 Note: The database contains tables with raw data, they will be used while running Go tests.

## Server & application
Start the server with the following command in the root directory:

> go run main.go

Access the application opening the link in your browser:

> http://localhost:8080/

## Stations

In this section of the page, all bike stations can be seen. A pagination of 20 stations per page/scrolling has been implemented. Clicking on each station, you will access the following information from them:

- Station name;
- Station address;
- Total number of journeys starting from this station;
- Total number of journeys ending at this station and;
- Station location on the map.

## Tests
Go tests (dbimportjourneys_test.go and dbstations_test.go) can be run using the command in the root directory:

> go test

Note: If tests have passed, unnecessary data to run the project will be deleted (csv files and unused tables, for example).