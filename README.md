# My Helsinki city bike app
Hello and welcome! I am Danilo and the README of my solution of [Solita's pre-assignment](https://github.com/solita/dev-academy-2023-exercise "Go to repo") for the Dev Academy 2023.

##### Table of contents
[1. Running the app locally](#locally)
[2. Data import](#data)

<a name="locally"/>
## Running the app locally
The project uses SQLite3, `wget` command, Go and internet connection. Make sure you have them all in your computer. Clone this repo in your IDE.

<a name="data"/>
## Data import
Go to the database folder, with ./dbcreate.sh you will run the script that will create the database. The data from the CSV files will be downloaded and saved to hsk-city-bike-app.db. Journeys that lasted for less than ten seconds and covered distances shorter than 10m has been removed, as well as duplicate data.