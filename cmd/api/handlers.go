package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

func (app *application) writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func (app *application) readJSON(r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(nil, r.Body, 1_048_576)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}

func (app *application) errorJSON(w http.ResponseWriter, err error, status ...int) {
	code := http.StatusBadRequest
	if len(status) > 0 {
		code = status[0]
	}
	_ = app.writeJSON(w, code, map[string]string{"error": err.Error()})
}

func (app *application) notFound(w http.ResponseWriter) {
	app.errorJSON(w, errors.New("not found"), http.StatusNotFound)
}
