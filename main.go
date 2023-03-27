package main

import (
	"fmt"
	"hsk-bikeapp-solita/handlers"
)

func main() {
	fmt.Println("running main")
	collection := []handlers.Handler{
		{
			Endpoint: "/",
			GetFunction: handlers.IndexGet,
		},
		{
			Endpoint: "/stations",
			GetFunction: handlers.StationsGet,
		},
	}

	handlers.Start(collection)
}
