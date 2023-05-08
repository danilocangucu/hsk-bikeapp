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
   3) [Add a station](#add-a-station)
6. [Journeys](#journeys)
7. [Testing](#testing)
   1) [Golang unit tests](#golang-unit-tests)
   2) [Cypress E2E tests](#cypress-e2e-tests)
9. [Running the application in Docker](#running-the-application-in-docker)
10. [Cloud-based Backend](#cloud-based-backend)
11. [Reflections](#reflections)

## Introduction

This is a Helsinki city bike app that displays information about bike stations and bike journeys during May, June and July of 2021 in Helsinki and Espoo.

## Getting Started

To get started with this project, you will need to have the following installed on your computer:
- [wget command](https://www.gnu.org/software/wget/)
- [SQLite3](https://www.sqlite.org/index.html)
- [Go](https://golang.org/)
  - [goquery package](https://github.com/PuerkitoBio/goquery)
- JavaScript (most likely your browser has it, if not, check the [JS documentation](https://developer.mozilla.org/en-US/docs/Web/JavaScript) for help)
  - [Mapbox GL JS library](https://docs.mapbox.com/mapbox-gl-js/api/)
  - [Cypress](https://docs.cypress.io/guides/overview/why-cypress)
- [Docker](https://www.docker.com/)

If you are using Windows:
- [Git Bash](https://gitforwindows.org/)
- [Node.js and npm](https://docs.npmjs.com/downloading-and-installing-node-js-and-npm)

Once you have these technologies installed, clone this repository to your IDE.

Insert the following from the second page of my cover letter document:
- Mapbox API Key to file src/Stations.js, line 110
- Google Cloud API Key to file src/AddStation.js, line 151
- Video.mp4 to src/

If you encounter any issues, please do not hesitate to contact me via phone or email, both of which can be found within the application.

## Data Import

To import the data, run the following script from the root directory with the command below.

Note: If you are using Windows, please run the command in Git Bash.

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
Serves the index.html file.

``/``  
Permanently redirects to /index and includes error handling for non-existent paths.

``/stations``  
API handling with the option to include the query parameter "id".  

   ``/stations?id=STATION_ID``  
   The STATION_ID must be an integer between 1 and 405. If new stations are added, the allowed id will increase.

``/journeys``  
API handling with the option to include the query parameter "batchfromid".  

   ``/journeys?batchfromid=JOURNEY_ID``  
   The JOURNEY_ID must be an integer between 1 and 1500580. A batch of maximum 500 journeys starting from the JOURNEY_ID will be returned.

``/addstation``  
Accepts POST requests to add new stations to the database.

For all the endpoint handling, please refer to the ``handlers`` directory.  
To access the database handling code, navigate to the ``database`` directory and open the files ``db.go`` and ``dbcreate.sh``.

## Stations

This section includes a counter of stations, a the list of all stations, a view of each individual station, and a form to add a station.

### Stations list

In this subsection of the page, all bike stations imported and added to the database, and available at ``/stations`` can be seen. A scrolling pagination of 20 stations per scrolling has been implemented. Clicking on each station will display information of the individul station.

To read about the functions that handle the stations list, please refer to the ``getStations`` and ``showAllStations`` functions in ``src/Stations.js`` and ``handlers/stations.go``.

### Single station view

When you view a single station, the following information will be retrieved by sending a GET request to the endpoint ``/stations?id=STATION_ID``:

- Station names in Finnish and Swedish;
- Station address in Finnish and Swedish;
- Total number of journeys that start from this station;
- Total number of journeys that end at this station; and
- Station location on the map with a popup showing the station's name and bikes capacity.

To learn more about the functions that handle the stations list, please refer to the ``getStationData``, ``showSingleStation``, and ``renderMap`` functions in ``src/Stations.js`` and ``handlers/addstation.go``.

Please note that in order to display the maps properly, you will need to insert the Mapbox API key that has been provided in the application in the ``renderMap`` function in ``src/Stations.js``, line 110.

### Add a station

In this subsection, you can add a new station using the provided form. The fields that need to be filled out include:

- Station Names: The station name should be provided in Finnish, Swedish, and English;
- Addresses: The address should be provided in Finnish and English;
- Operator (Optional): While the dataset only includes "CityBike Finland" or no operator at all, users can choose to leave this field empty or enter an operator of their choice;
- Capacity: The capacity field has a maximum value of 44, which is based on the largest capacity present in the dataset;
- Count of journeys that originated from this station;
- Count of journeys that ended at this station.

The last two fields were included because even if a station has not been added to the system, the company may have this information and give it to the user.

The form will be validaded and, if the validation succeeds, a new station will be added in the database. The station will be automatically rendered in the Single station view subsection.

Notes:
1. The fields for city names ("Kaupunki" in Finnish and "Stad" in Swedish), x (Latitude), and y (Longitude) in the database will be populated automatically once the addresses are validated through a Google Cloud API request. To ensure that these requests are successful, please insert the API KEY provided in the application into the ``validateAddresses`` function located in ``src/AddStation.js``, line 151.
2. The fields "FID" and "ID" will be automatically filled.

To view how this section is implemented, please refer to the functions in ``src/AddStation.js`` and ``handlers/addstation.go``.

## Journeys

Here, you can look at all journeys from the database. A counter of loaded journeys is also displayed. Since there's a lot of data, I chose to import batches of 500 journeys and then a scrolling pagination of 50 journeys per scrolling. In the table of this section, the following information can be seen:

- Departure (date and time);
- From (Finnish name of departure station);
- To (Finnish name of return station);
- Distance (in km); and
- Duration (of the journey, in hours, minutes and seconds).

## Testing

### Golang unit tests

To run the tests located in the ``tests`` directory of the root directory in verbose mode, use the following command:

```
go test -v ./tests
```

The following tests will be performed:

``handle_indexget_test.go``  
Tests an HTTP server's response to a GET request by verifying the presence of specific HTML elements and assets in the response body.

``dbstations_test.go``  
Imports station data from a CSV file, ``/datasets/726277c507ef4914b0aec3cbcfcbfafc_0.csv``, into a SQLite database, creates a table for the data, inserts the data into the table, and verifies that the data was inserted correctly by checking the number of rows in the table and comparing it with the number of records in the CSV file.

``dbimportjourneys_test.go``  
Imports data from multiple CSV files in ``/datasets`` (2021-05.csv, 2021-06.csv, 2021-07.csv) into a SQLite database table, verifies that the table exists or creates it if necessary, inserts the data into a temporary table, and checks if the imported data matches the raw data.

``dbshortjourneys_test.go``  
Checks if there are any rows in ``/database/hsk-city-bike-app.db`` with covered distance or duration less than 10 (meters or seconds).

``apis_test.go``  
Checks the availability of the localhost at ``http://localhost:8080/`` and runs three tests on the bike sharing app API that returns station information: ``TestGetInvalidStationID``, ``TestGetNonExistingStation``, and ``TestGetValidStationInfo``. These tests are available only on Mac and Windows platforms.

### Cypress E2E tests

Cypress tests are located in ``cypress/e2e`` directory. To run these tests, you will need two terminal windows. You should run the server in one terminal and the following command in the other, from the project's root directory.

If you're using an unix-based OS:

```
node_modules/.bin/cypress open
```

If you're using Windows:
```
npx cypress open
```


The Cypress UI will open. Select ``E2E Testing``, choose your browser and click on  ``Start E2E Testing``.  A new window will be opened. Under the "Specs" sections there are two test files:

``journeys.cy.js``  
Ensures journey details in the web app match the API data by visiting the page, requesting data, and comparing the details in the table with the API response after extracting relevant data and converting timestamps.

``stations.cy.js``  
The first test on this file clicks on a station in a list, waits for the station details to be loaded, makes a request to the detail page API, and asserts that the details of the clicked station match the API data. The second test scrolls down a list of stations and verifies that 20 more station names are added per scroll.

To run the the tests, click on the previous files in the "Sepcs" section in Cypress UI.

## Running the application in Docker
A Dockerfile is included in the project and you can dockerize this application.

First open Docker in your computer. Then, to build the Docker image, execute the following command from the primary directory of this project:
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

## Reflections

I am very happy and proud of this project, having devoted substantial time and effort to it. My motivation arises not only from my aspiration to join the Dev Academy but also from my innate drive to excel in tasks - perhaps a characteristic of my Capricorn nature?

During this journey, I acquired valuable knowledge on various services such as AWS, Mapbox, and Google Cloud API. Moreover, I sharpened my skills in E2E testing, Go, JavaScript, and HTML/CSS. From the beginning, I aimed to challenge myself and broaden my expertise in specific areas. For instance, I was intrigued by the video effect visible upon page load. I encountered it on a website during my daily "web design exploration" activity through websites ([Awwwards](https://www.awwwards.com/) is my favourite!) and thought it would be an interesting addition to this project. At one point, I became enamored with the idea of using the official HSL maps with Mapbox - yes, I explored HSLdevcom's GitHub, and obtained a Digitransit API key. I spent a couple of days attempting to make it work, but success eluded me (for now!).

What began as a simple intention to complete one or two "extra section" tasks rapidly evolved into an insatiable thirst for exploring and experimenting with new features and technologies. And I accomplished them all! Given the wealth of learning opportunities presented by the assignment, I anticipate that the Academy will provide numerous challenges for me to explore and share with my peers - YAY!

### Future Improvements
While I take pride in this project, there are certain aspects that could use some improvement. Therefore, I plan to continue updating the project even after submitting it to Solita on Sunday, 07.05.2023. Here's a list of areas for improvement that I've noted down in my notebook:

- Journeys Batch Error  
    I discovered an issue while fetching batches of journeys. I'd like to better understand the problem and find a solution.

- Prevent Scrolling While Loading a New Station  
    This issue was quite frustrating! I've explored various options involving click events, focus, and CSS, but I'm not entirely satisfied with the current solution. It works for now, but I'd like to improve it.

- Styling for Small Devices  
    I haven't been able to thoroughly test the styling on smaller devices. I plan to test the application on my cellphone and iPad to identify what works well and what needs adjustment.

- Error Handling in the createdb Script  
    This is an urgent issue that needs to be addressed. I hope you don't encounter any errors while running the script!

- Testing the Project on More Systems  
    So far, I've tested the project on my MacBook and my boyfriend's Windows laptop. I'd like to test it on other systems to ensure everything works as expected.

- Database Test with Raw Data  
    I've been questioning the appropriateness of this test, as it only uses raw data from the stations dataset. I'd like to improve the test by incorporating filtered data (e.g., no journeys shorter than 10 meters and 10 seconds, no duplicated data).

- Loading Transition for the Index Page  
    Since it takes a few milliseconds to load all the content on the page, adding a smooth and elegant loading transition could enhance the user experience.

I hope you found this README helpful and informative! If you have any questions or suggestions, feel free to reach out to me. I appreciate your time and consideration, and am looking forward to seeing you at Solita soon! 😊

Happy coding! 🚀
