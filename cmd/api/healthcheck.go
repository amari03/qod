package main

import (
	
	"net/http"
)

//regular go structure for dependencies
func (a * application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	//panic("Apples & Oranges")   // deliberate panic
	
	// Create a map to hold the response data
	data := envelope {
		"status": "available",
		"system_info": map[string]string{
				"environment": a.config.env,
				"version": version,
	   },
}

	err := a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	   }	
}
