package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxSizeDataTORead := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxSizeDataTORead))
	decoded := json.NewDecoder(r.Body)
	decoded.DisallowUnknownFields()

	return decoded.Decode(data)

}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	encoded := json.NewEncoder(w).Encode(data)
	return encoded
}

func writeJSONError(w http.ResponseWriter, status int, data any) error {
	type errFmt struct {
		Error any `json:"error"`
	}

	return writeJSON(w, status, &errFmt{Error: data})

}

func jsonResponse(w http.ResponseWriter, status int, data any) error {

	type envelope struct {
		Data any `json:"data"`
	}

	return writeJSON(w, status, &envelope{Data: data})
}
