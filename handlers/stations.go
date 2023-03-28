package handlers

import (
	"encoding/json"
	"fmt"
	db "hsk-bikeapp-solita/database"
	"io"
	"net/http"
	"strconv"
)

func StationsGet(w http.ResponseWriter, r *http.Request) {
	var err error

	params := r.URL.Query()

	fmt.Println(params)

	filter := db.StationFilter{}

	if id, isExist := params["id"]; isExist {
		fmt.Println(id)
		filter.StationId, err = strconv.Atoi(id[0])
	}

	if err != nil {
		errorMessage := "Invalid station id"
		GetErrorResponse(w, errorMessage, http.StatusBadRequest)
		return
	}

	stations, err := DB.GetStations(filter)
	if err != nil {
		errorMessage := "Error while getting station" // temp err message
		GetErrorResponse(w, errorMessage, http.StatusBadRequest)
		return
	}
	
	w.WriteHeader(http.StatusAccepted)
	result, err := json.Marshal(stations)
	if err != nil {
		//todo error handling
		fmt.Println("error json")
	}

	io.WriteString(w, string(result))
}