package handlers

import (
	"encoding/json"
	db "hsk-bikeapp-solita/database"
	"io"
	"net/http"
)

func AddStationPost(w http.ResponseWriter, r *http.Request) {
	var errors []error

	body, err := io.ReadAll(r.Body)
	if err != nil {
		errors = append(errors, err)
		sendErrorResponse(w, http.StatusBadRequest, errors)
		return
	}

	var newStation db.Station
	if err := json.Unmarshal(body, &newStation); err != nil {
		errors = append(errors, err)
		sendErrorResponse(w, http.StatusBadRequest, errors)
		return
	}

	validationErrors := DB.ValidateNewStation(newStation)

	if len(validationErrors) > 0 {
		errors = append(errors, validationErrors...)
		sendErrorResponse(w, http.StatusBadRequest, errors)
		return
	}

	addNewStationError := DB.AddNewStation(newStation)
	if addNewStationError != nil {
		errors = append(errors, []error{addNewStationError}...)
		sendErrorResponse(w, http.StatusBadRequest, errors)
	}

	response := map[string]string{
		"message": "Station added to database",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func sendErrorResponse(w http.ResponseWriter, status int, errors []error) {
	response := make(map[string][]string)
	errorMessages := make([]string, 0, len(errors))

	for _, err := range errors {
		errorMessages = append(errorMessages, err.Error())
	}

	response["errors"] = errorMessages
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
