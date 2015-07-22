package main

import (
	"encoding/json"
	"io"
	"net/http"
)

func httpJsonBodyEndpoint(w http.ResponseWriter, r *http.Request, request interface{}, cb callback) {
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}

	var err error
	if err = json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "{\"error\": {\"message\": \"Bad request\"}}")
		return
	}
	r.Body.Close()

	httpJsonEndpoint(w, r, cb)
}

func httpJsonEndpoint(w http.ResponseWriter, r *http.Request, cb callback) {
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}

	var resp interface{}
	var err error
	if resp, err = cb(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "{\"error\": {\"message\": \"Internal server error\"}}")
		return
	}

	if err = json.NewEncoder(w).Encode(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "{\"error\": {\"message\": \"Internal server error\"}}")
	}
}
