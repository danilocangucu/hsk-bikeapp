package handlers

import (
	"fmt"
	"hsk-bikeapp-solita/database"
	"io"
	"log"
	"net/http"
	"os"
)

type handlerFunction func(w http.ResponseWriter, r *http.Request)

type Handler struct {
	Endpoint     string
	Method       string
	GetFunction  handlerFunction
	PostFunction handlerFunction
}

var DB database.Db
var err error

func Start(collection []Handler) {
	DB, err = database.OpenDatabase()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening database: %v\n", err)
		os.Exit(1)
	}
	defer DB.CloseDatabase()

	mux := http.NewServeMux()

	for _, handler := range collection {
		mux.Handle(handler.Endpoint, GetFunc(handler))
	}

	jsFs := http.FileServer(http.Dir("./src"))
	mux.Handle("/src/", http.StripPrefix("/src/", jsFs))

	log.Printf("Server listening on http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func GetFunc(handler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && handler.GetFunction != nil {
			handler.GetFunction(w, r)
		} else if r.Method == "POST" && handler.PostFunction != nil {
			handler.PostFunction(w, r)
		} else {
			w.WriteHeader(http.StatusNotFound)
			io.WriteString(w, "Not found")
		}
	}
}
