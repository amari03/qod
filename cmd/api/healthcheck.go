package main

import (
	
	"net/http"
)

//regular go structure for dependencies
func (a * application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// Create a map to hold the response data
	response := map[string]string{
		"status":      "available",
		"environment": a.config.env,
		"version":     version,
	}

	err := a.writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		a.logger.Error(err.Error())
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
		return
	}
}
