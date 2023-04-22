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
	var stations []db.Station
	var station db.Station

	params := r.URL.Query()
	filter := db.StationFilter{}

	if id, idExist := params["id"]; idExist {
		filter.ID, err = strconv.Atoi(id[0])
		if err != nil {
			log.Println("Invalid query parameter:", err)
			http.Error(w, "Error while getting stations", http.StatusBadRequest)
			return
		}
	}

	if filter.ID != 0 {
		station, err = DB.GetSingleStation(filter)
		if err != nil {
			log.Println("Error while getting station:", err)
			http.Error(w, "Error while getting station", http.StatusBadRequest)
			return
		}
	} else {
		stations, err = DB.GetAllStations()
		if err != nil {
			log.Println("Error while getting stations:", err)
			http.Error(w, "Error while getting stations", http.StatusBadRequest)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if filter.ID != 0 {
		if err := json.NewEncoder(w).Encode(station); err != nil {
			log.Println("Error encoding response:", err)
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
	} else {
		if err := json.NewEncoder(w).Encode(stations); err != nil {
			log.Println("Error encoding response:", err)
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
	}
}
