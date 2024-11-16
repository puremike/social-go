package main

import (
	"log"
	"net/http"
)


func (app *application) internalServer(w http.ResponseWriter, r *http.Request, err error) {

	data := map[string]string {
        "status" : "error",
        "message" : "Internal Server Error",
    }

	log.Printf("internal server error: %s\n path: %s\n error: %s", r.Method, r.URL.Path, err)

	writeJSONError(w, http.StatusInternalServerError, data)
}

func (app *application) badRequest(w http.ResponseWriter, r *http.Request, err error) {

	data := map[string]string {
        "status" : "error",
        "message" : "Bad Request",
		"error" : err.Error(),
    }
	
	log.Printf("bad request error: %s\n path: %s\n error: %s", r.Method, r.URL.Path, err)

	writeJSONError(w, http.StatusBadRequest, data)
}

func (app *application) notFound(w http.ResponseWriter, r *http.Request, err error) {

	data := map[string]string {
        "status" : "error",
        "message" : "Resource not found",
		"error" : err.Error(),
    }
	
	log.Printf("resource not found: %s\n path: %s\n error: %s", r.Method, r.URL.Path, err)

	writeJSONError(w, http.StatusNotFound, data)
}