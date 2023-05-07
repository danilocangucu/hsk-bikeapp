package handlers

import (
	"encoding/json"
	"fmt"
	db "hsk-bikeapp-solita/database"
	"log"
	"net/http"
	"strconv"
)

var lastJourneyId int

func JourneysGet(w http.ResponseWriter, r *http.Request) {
	var err error

	params := r.URL.Query()

	filter := db.JourneyFilter{}

	if id, isExist := params["batchfromid"]; isExist {
		filter.BatchFromId, err = strconv.Atoi(id[0])
	} else if len(params) == 0 {
		filter.BatchFromId = 1
	}

	if err != nil {
		log.Println("Invalid query parameter")
		http.Error(w, "Could not retreive journeys data", http.StatusBadRequest)
		return
	}

	if lastJourneyId == 0 {
		lastJourneyId, err = DB.GetLastJourneyId()
		if err != nil {
			log.Println("Error while last journey ID:", err)
			http.Error(w, "Error while getting journeys", http.StatusBadRequest)
			return
		}
	}

	if filter.BatchFromId < 1 || filter.BatchFromId > lastJourneyId {
		log.Println("Invalid batch from id:", filter.BatchFromId)
		http.Error(w, "Could not retreive journeys data", http.StatusBadRequest)
		return
	}

	filter.Limit = 3000
	remainingIds := lastJourneyId - filter.Limit
	if remainingIds < 3000 {
		filter.Limit = remainingIds
	}

	journeys, err := DB.GetJourneys(filter)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(journeys); err != nil {
		log.Println("Error encoding response:", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

}
