package main

import (
	"net/http"
	"time"
)

func (app *application) health(w http.ResponseWriter, r *http.Request) {

	data := map[string]string{
		"status":      "OK",
		"message":     "Application is Healthy",
		"environment": app.config.environment,
	}

	time.Sleep(time.Second * 3)

	if err := jsonResponse(w, http.StatusOK, data); err != nil {
		app.internalServer(w, r, err)
		return
	}
}
