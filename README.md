# My Helsinki City Bike Single Page App

Greetings and welcome! This is my solution for [Solita's pre-assignment](https://github.com/solita/dev-academy-2023-exercise "Go to repo") for the Dev Academy 2023. My name is Danilo and I am thrilled to present my work.

Before diving into the code, I would like to express my gratitude to my colleague Iuliia Chipsanova for sharing her coding knowledge with me throughout the past year at Grit:Lab. I would also like to thank my boyfriend Jan Korte for his patience and valuable inputs on this project, despite hearing about "the bikes" every day.

Creating this app has been a challenging month-long journey, and I am proud and excited to finally share the result with you. Throughout the process, I worked on it almost every day, updating both this GitHub repository and my [Trello dashboard](https://trello.com/b/ZfZX3lh6/tasks).

Now, let's take a closer look at the code!

## Table of Contents

1. [Introduction](#introduction)
2. [Getting Started](#getting-started)
3. [Data Import](#data-import)
4. [Running the App](#running-the-app)
5. [Stations](#stations)
   1) [Stations list](#stations-list)
   2) [Single station view](#single-station-view)
   3) [Add station](#add-station)
6. [Journeys](#journeys)
7. [Testing](#testing)
8. [Running the application in Docker](#running-the-application-in-docker)
9. [Cloud-based Backend](#cloud-based-backend)
10. [Contributing](#contributing)

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

This command will create the database for the application and download the datasets of Journeys and Stations provided in the pre-assignment. The resulting data will be saved to hsk-city-bike-app.db. During the import process:

- Bike journeys that lasted less than 10 seconds and covered distances shorter than 10m have been removed;
- The number of journeys that started and ended from each station has been counted; and
- Duplicate data has been removed.

Note: The database contains tables with raw data, which are used while running Go tests.

## Running the App

To start the server, run the following command from the root directory:

```
go run main.go
```

Then, open the application by accessing http://localhost:8080/ in your browser.

The following endpoints are implemented:  

``/index``  
Serves the index.html file

``/``  
Permanently redirects to /index and includes error handling for non-existent paths

``/stations``  
API handling with the option to include ?id=STATION_ID

``/journeys``  
API handling with the option to include ?lastJourneyId=JOURNEY_ID

``/addstation``  
Accepts POST requests to add new stations to the database.

For all the endpoint handling, please refer to the ``handlers`` directory.  
To access the database handling code, navigate to the ``database`` directory and open the files ``db.go`` and ``dbcreate.sh``.

## Stations

This section includes the list of all stations, a view of each individual station, and a form to add a station.

### Stations list

In this section of the page, all bike stations imported and added to the database, and available at ``/stations`` can be seen. A scrolling pagination of 20 stations per scrolling has been implemented. Clicking on each station will display information of the individul station.

To read about the functions that handle the stations list, please refer to the ``getStations`` and ``showAllStations`` functions in ``src/Stations.js``.

### Single station view

When you view a single station, the following information will be retrieved by sending a request to the endpoint ``/stations?id=STATION_ID``:

- Station names in Finnish and Swedish;
- Station address in Finnish and Swedish;
- Total number of journeys that start from this station;
- Total number of journeys that end at this station; and
- Station location on the map with a popup showing the station's name and bikes capacity.

To learn more about the functions that handle the stations list, please refer to the ``getStationData``, ``showSingleStation``, and ``renderMap`` functions in ``src/Stations.js``.

Please note that in order to display the maps properly, you will need to insert the Mapbox API key that has been provided in the application in the ``renderMap`` function. If you encounter any issues, please do not hesitate to contact me via phone or email, both of which can be found within the application.

### Add Station

In this subsection, you can add a new station using the provided form. The fields that need to be filled out include:

- ID: The number of the station. Please note that there is no discernible pattern to the station IDs in the stations dataset, so the user must input a number.
- Station Names: The station name should be provided in Finnish, Swedish, and English.
- Addresses: The address should be provided in Finnish and English.
- Operator (Optional): Operator (Optional): While the dataset only includes "CityBike Finland" or no operator at all, users can choose to leave this field empty or enter an operator of their choice.
- Capacity: The capacity field has a maximum value of 44, which is based on the largest capacity present in the dataset.
- Count of journeys that originated from this station
- Count of journeys that ended at this station.

The last two fields were included because even if a station has not been added to the system, the company may have this information and give it to the user.

The form will be validaded and, if the validation succeeds, a new station will be added in the database. The station will be automatically rendered in the Single station view subsection.

To view how this section is implemented, please refer to the functions in ``src/AddStation.js``.

Note: The fields for city names ("Kaupunki" in Finnish and "Stad" in Swedish), x (Latitude), and y (Longitude) in the database will be populated automatically once the addresses are validated through a Google Cloud API request. To ensure that these requests are successful, please insert the API KEY provided in the application into the ``validateAddresses`` function located in ``src/AddStation.js``. If you encounter any issues, please do not hesitate to contact me.

## Journeys

Here, you can look at all journeys from the database. Since there's a lot of data, I chose to import batches of 3000 journeys and then a scrolling pagination of 50 journeys per scrolling. In the table of this section, the following information can be seen:

- Departure (date and time);
- From (departure station);
- To (return station);
- Distance (in km); and
- Duration (of the journey, in hours, minutes and seconds).

## Testing

### Golang unit tests

To run the tests located in the ``tests`` directory of the root directory in verbose mode, use the following command:

```
go test -v ./tests
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
Checks the availability of the localhost at ``http://localhost:8080/`` and runs three tests on the bike sharing app API that returns station information: ``TestGetInvalidStationID``, ``TestGetNonExistingStation``, and ``TestGetValidStationInfo``. These tests are available only on Mac and Windows platforms.

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

## Running the application in Docker
A Dockerfile is included in the project and you can dockerize this application. To build the Docker image, execute the following command from the primary directory:
```
docker build -t hsk-bikeapp-solita .
```

After building the Docker image, you can run the application using the following command:

```
docker run -p 8080:8080 hsk-bikeapp-solita
```

When the application is running, you can access it by navigating to http://localhost:8080 in our web browser.

## Cloud-based Backend
The backend of this application has been migrated to the cloud using Amazon Web Services (AWS). To enable this, I created a separate repository for the cloud implementation, which can be found [here](https://github.com/danilocangucu/hsk-bikeapp-solita-cloud).

## Contributing

Loved this project and want to contibute? Please fork this repository and submit a pull request with your changes.
