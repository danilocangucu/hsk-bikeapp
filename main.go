package main

import (
	"hsk-bikeapp-solita/handlers"
)

func main() {
	collection := []handlers.Handler{
		{
			Endpoint:    "/index",
			GetFunction: handlers.IndexGet,
		},
		{
			Endpoint:    "/stations",
			GetFunction: handlers.StationsGet,
		},
		{
			Endpoint:    "/journeys",
			GetFunction: handlers.JourneysGet,
		},
		{
			Endpoint:    "/",
			GetFunction: handlers.HandleNotFound,
		},
	}

	handlers.Start(collection)
}
