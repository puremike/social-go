package main

import "net/http"

func (app *application) health(w http.ResponseWriter, r *http.Request) {

	data := map[string]string {
		"status" : "OK",
		"message" : "Application is Healthy",
		"environment" : app.config.environment,
	}
	if  err := writeJSON(w, http.StatusOK, data); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to write JSON response")
        return   
	}
}