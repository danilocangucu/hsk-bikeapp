package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func StationsGet(w http.ResponseWriter, r *http.Request) {
	var err error

	// params := r.URL.Query()

	// filter := db.StationsFilter{}

	// // if id, isExist := params["id"]; isExist {
	// 	filter.StationID, err = strconv.Atoi(stationId[0])
	// }

	// if err != nil {
	// 	errorMessage := "Invalid station id"
	// GetErrorResponse(w, errorMessage, http.StatusBadRequest)
	// return
	// } else {
		stations, err := DB.GetStation()
		if err != nil {
			errorMessage := "Error while getting station" // temp err message
			GetErrorResponse(w, errorMessage, http.StatusBadRequest)
			return
		}
	// }
	
	w.WriteHeader(http.StatusAccepted)
	result, err := json.Marshal(stations)
	if err != nil {
		//todo error handling
		fmt.Println("error json")
	}

	io.WriteString(w, string(result))
}