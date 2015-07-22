package main

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

func httpJsonBodyEndpoint(w http.ResponseWriter, r *http.Request, incoming interface{}, instruments *instrumentation, cb callback) error {
	var err error
	defer func(begin time.Time) {
		instruments.logger.Log("error", err, "took", time.Since(begin))
		instruments.total.Add(1)
		instruments.duration.Observe(time.Since(begin))
	}(time.Now())

	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}

	if err = json.NewDecoder(r.Body).Decode(incoming); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "{\"error\": {\"message\": \"Bad request\"}}")
		return err
	}
	r.Body.Close()

	err = httpJsonEndpointHandler(w, r, cb)
	return err
}

func httpJsonEndpoint(w http.ResponseWriter, r *http.Request, instruments *instrumentation, cb callback) error {
	var err error
	defer func(begin time.Time) {
		instruments.logger.Log("error", err, "took", time.Since(begin))
		instruments.total.Add(1)
		instruments.duration.Observe(time.Since(begin))
	}(time.Now())

	err = httpJsonEndpointHandler(w, r, cb)
	return err
}

func httpJsonEndpointHandler(w http.ResponseWriter, r *http.Request, cb callback) error {
	var err error
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}

	var resp interface{}
	if resp, err = cb(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "{\"error\": {\"message\": \"Internal server error\"}}")
		return err
	}

	if err = json.NewEncoder(w).Encode(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "{\"error\": {\"message\": \"Internal server error\"}}")
		return err
	}
	return err
}
