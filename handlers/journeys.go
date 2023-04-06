package handlers

import (
	"encoding/json"
	"fmt"
	db "hsk-bikeapp-solita/database"
	"log"
	"net/http"
	"strconv"
)

func JourneysGet(w http.ResponseWriter, r *http.Request) {
	var err error

	params := r.URL.Query()

	filter := db.JourneyFilter{}

	if id, isExist := params["lastJourneyId"]; isExist {
		filter.LastId, err = strconv.Atoi(id[0])
	}
	fmt.Println(filter.LastId)

	if err != nil {
		log.Println("Invalid last journey id")
		http.Error(w, "Could not retreive journeys data", http.StatusBadRequest)
		return
	}

	lastJourney, err := DB.GetLastJourneyId()
	if err != nil {
		log.Println("Error while getting journeys:", err)
		http.Error(w, "Error while getting journeys", http.StatusBadRequest)
		return
	}

	filter.Limit = 3000
	remainingIds := lastJourney.LastId - filter.Limit
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