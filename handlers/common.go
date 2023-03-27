package handlers

import (
	"encoding/json"
	"fmt"
	"hsk-bikeapp-solita/database"
	"html/template"
	"io"
	"log"
	"net/http"
)

const responseError = "Response error"

type handlerFunction func(w http.ResponseWriter, r *http.Request)

type Handler struct {
	Endpoint string
	Method string
	GetFunction handlerFunction
	PostFunction handlerFunction
}

type Response struct {
	Status string
	Data string
}

type ResponseError struct {
	Status string
	Error string
}

var DB database.Db

func Start(collection []Handler) {
	DB = database.OpenDatabase()
	defer DB.Close()

	mux := http.NewServeMux()

	for _, handler := range collection {
		mux.Handle(handler.Endpoint, GetFunc(handler))
	}
	fmt.Println("server started")
	jsFs := http.FileServer(http.Dir("./src"))
	mux.Handle("/src/", http.StripPrefix("/src/", jsFs))
	log.Fatal(http.ListenAndServe(":8080", mux))

}

func GetFunc(handler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		if r.Method == "GET" && handler.GetFunction != nil {
			handler.GetFunction(w, r)
		} else if r.Method == "POST" && handler.PostFunction != nil {
			handler.PostFunction(w, r)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func GetErrorResponse(w http.ResponseWriter, errorMessage string, statusCode int) {
	w.WriteHeader(statusCode)
	response := ResponseError{Status: responseError, Error: errorMessage}
	result, err := json.Marshal(response)
	if err != nil {
		//todo handle error
		fmt.Println("geterrresponse error")
	}
	io.WriteString(w, string(result))
}

func IndexGet(w http.ResponseWriter, r *http.Request) {
	templ := template.Must(template.New("index.html").ParseFiles("index.html"))
	err := templ.ExecuteTemplate(w, "index.html", "") //todo handle error
	if err != nil {
		fmt.Println(err)
	}
}