package main

import (
	"net/http"
)

func (app *application) internalServer(w http.ResponseWriter, r *http.Request, err error) {

	app.logger.Errorw("internal server error", "method", r.Method, "path", r.URL.Path, "error", err)

	writeJSONError(w, http.StatusInternalServerError, "internal server error")
}

func (app *application) badRequest(w http.ResponseWriter, r *http.Request, err error) {

	app.logger.Warnw("bad request", "method", r.Method, "path", r.URL.Path, "error", err)

	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFound(w http.ResponseWriter, r *http.Request, err error) {

	app.logger.Warnw("resource not found", "method", r.Method, "path", r.URL.Path, "error", err)

	writeJSONError(w, http.StatusNotFound, "not found")
}

// func (app *application) conflictError(w http.ResponseWriter, r *http.Request, err error) {

// 	// data := map[string]string {
// 	// 	"status" : "error",
// 	// 	"message" : "conflict error",
// 	// 	"error" : err.Error(),
// 	// }

// 	app.logger.Errorw("conflict error", "method", r.Method, "path", r.URL.Path, "err", err)

// 	writeJSONError(w, http.StatusConflict, err.Error())
// }

func (app *application) unauthorizedError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("unauthorized error", "method", r.Method, "path", r.URL.Path, "error", err)

	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	http.Error(w, "Unauthorized", http.StatusUnauthorized)

	writeJSONError(w, http.StatusConflict, "unauthorized error")
}

func (app *application) unauthorizedErrorOthers(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("unauthorized error", "method", r.Method, "path", r.URL.Path, "error", err)
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	writeJSONError(w, http.StatusConflict, "unauthorized error")
}

func (app *application) forbiddenResponse(w http.ResponseWriter, r *http.Request) {
	app.logger.Warnf("forbidden error", "method", r.Method, "path", r.URL.Path, "error")

	writeJSONError(w, http.StatusForbidden, "forbidden error")
}

func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request, retryAfter string) {
	app.logger.Warnw("rate limit exceeded", "method", r.Method, "path", r.URL.Path)

	w.Header().Set("Retry-After", retryAfter)

	writeJSONError(w, http.StatusTooManyRequests, "rate limit exceeded, retry after: "+retryAfter)
}
