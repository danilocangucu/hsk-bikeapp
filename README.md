# My Helsinki City Bike App

Hello and welcome! I am Danilo and this is the README for my solution of [Solita's pre-assignment](https://github.com/solita/dev-academy-2023-exercise "Go to repo") for the Dev Academy 2023.

## Table of Contents

1. [Introduction](#introduction)
2. [Getting Started](#getting-started)
3. [Data Import](#data-import)
4. [Running the App](#running-the-app)
5. [Stations](#stations)
6. [Journeys](#journeys)
7. [Testing](#testing)
8. [Contributing](#contributing)

## Introduction

This is a Helsinki city bike app that displays information about bike stations and bike journeys during May, June and July of 2021 in Helsinki and Espoo.

## Getting Started

To get started with this project, you will need to have the following installed on your computer:
- [Go](https://golang.org/)
  - [goquery package](https://github.com/PuerkitoBio/goquery)
- JavaScript (most likely your browser has it, if not, check the [JS documentation](https://developer.mozilla.org/en-US/docs/Web/JavaScript) for help)
  - [Mapbox GL JS library](https://docs.mapbox.com/mapbox-gl-js/api/)
  - [Cypress](https://docs.cypress.io/guides/overview/why-cypress)
- [SQLite3](https://www.sqlite.org/index.html)
- [wget command](https://www.gnu.org/software/wget/)
- Internet connection

Once you have these technologies installed, clone this repository to your IDE.

## Data Import

To import the data, run the following script from the root directory:

```
./database/dbcreate.sh
```

This will create the database for the application and download the data from the CSV files. The data will be saved to `hsk-city-bike-app.db`. During the import process:

- Bike journeys that lasted less than 10 seconds and covered distances shorter than 10m have been removed;
- The number of journeys that started and ended from each station has been counted; and
- Duplicate data has been removed.

Note: The database contains tables with raw data, which are used while running Go tests.

## Running the App

To start the server, run the following command from the root directory:

```
go run main.go
```

Then, access the application by opening the following link in your browser:

```
http://localhost:8080/
```

## Stations

In this section of the page, all bike stations imported to the database can be seen. A scrolling pagination of 20 stations per scrolling has been implemented. Clicking on each station will display the following information:

- Station name;
- Station address;
- Total number of journeys starting from this station;
- Total number of journeys ending at this station; and
- Station location on the map.

## Journeys

Here, you can look at all journeys from the database. Since there's a lot of data, I chose to import batches of 3000 journeys and then a scrolling pagination of 50 journeys per scrolling. In the table of this section, the following information can be seen:

- Departure (date and time)
- From (departure station)
- To (return station)
- Distance (in km)
- Duration (of the journey, in hours, minutes and seconds)

## Testing

### Golang unit tests

To run the tests located in the ``test`` directory of the root directory in verbose mode, use the following command:

```
go test -v ./test
```

``handle_indexget_test.go``  
Tests an HTTP server's response to a GET request by verifying the presence of specific HTML elements and assets in the response body.

``dbstations_test.go``  
Imports station data from a CSV file into a SQLite database, creates a table for the data, inserts the data into the table, and verifies that the data was inserted correctly by checking the number of rows in the table and comparing it with the number of records in the CSV file.

``dbimportjourneys_test.go``  
Imports data from multiple CSV files into a SQLite database table, verifies that the table exists or creates it if necessary, inserts the data into a temporary table, and checks if the imported data matches the raw data.

``dbshortjourneys_test.go``  
Checks if there are any rows in a specified database table with covered distance or duration less than 10 (meters or seconds).

``apis_test.go``  
Checks the availability of the localhost at http://localhost:8080/ and runs three tests on the bike sharing app API that returns station information: TestGetInvalidStationID, TestGetNonExistingStation, and TestGetValidStationInfo

### Cypress E2E tests

Cypress tests are located in ``cypress/e2e`` directory. To run them, execute the command:

```
node_modules/.bin/cypress open
```

The Cypress UI will open. Select ``E2E Testing``, choose your browser and click on  ``Start E2E Testing``.  A new window will be opened. Under the "Specs" sections there are two test files:

``journeys.cy.js``  
Ensures journey details in the web app match the API data by visiting the page, requesting data, and comparing the details in the table with the API response after extracting relevant data and converting timestamps.

``stations.cy.js``  
The first test on this file clicks on a station in a list, waits for the station details to be loaded, makes a request to the detail page API, and asserts that the details of the clicked station match the API data. The second test scrolls down a list of stations and verifies that 20 more station names are added per scroll.


## Contributing

Loved this project and want to contibute? Please fork this repository and submit a pull request with your changes.
