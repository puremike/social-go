package main

import (
	"net/http"
)


func (app *application) internalServer(w http.ResponseWriter, r *http.Request, err error) {
	
	data := map[string]string {
        "status" : "error",
        "message" : "Internal Server Error",
		"error": err.Error(),
    }

	app.logger.Errorw("internal server error", "method", r.Method, "path", r.URL.Path, "error", err)

	writeJSONError(w, http.StatusInternalServerError, data)
}

func (app *application) badRequest(w http.ResponseWriter, r *http.Request, err error) {

	data := map[string]string {
        "status" : "error",
        "message" : "Bad Request",
		"error" : err.Error(),
    }
	
	app.logger.Warnw("bad request", "method", r.Method, "path", r.URL.Path, "error", err)

	writeJSONError(w, http.StatusBadRequest, data)
}

func (app *application) notFound(w http.ResponseWriter, r *http.Request, err error) {

	data := map[string]string {
        "status" : "error",
        "message" : "Resource not found",
		"error" : err.Error(),
    }

	app.logger.Warnw("resource not found", "method", r.Method, "path", r.URL.Path, "error", err)

	writeJSONError(w, http.StatusNotFound, data)
}

func (app *application) conflictError (w http.ResponseWriter, r *http.Request, err error) {

	data := map[string]string {
		"status" : "error",
		"message" : "conflict error",
		"error" : err.Error(),
	}

	app.logger.Errorw("conflict error", "method", r.Method, "path", r.URL.Path, "err", err)

	writeJSONError(w, http.StatusConflict, data)
} 