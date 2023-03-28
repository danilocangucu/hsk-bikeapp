package handlers

import (
	"encoding/json"
	db "hsk-bikeapp-solita/database"
	"log"
	"net/http"
	"strconv"
)

func StationsGet(w http.ResponseWriter, r *http.Request) {
	var err error

	params := r.URL.Query()

	filter := db.StationFilter{}

	if id, isExist := params["id"]; isExist {
		filter.StationId, err = strconv.Atoi(id[0])
	}

	if err != nil {
		log.Println("Invalid station id")
		http.Error(w, "Invalid station id", http.StatusBadRequest)
		return
	}

	stations, err := DB.GetStations(filter)
	if err != nil {
		log.Println("Error while getting stations:", err)
		http.Error(w, "Error while getting stations", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(stations); err != nil {
        log.Println("Error encoding response:", err)
        http.Error(w, "Error encoding response", http.StatusInternalServerError)
        return
    }

}