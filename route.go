package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Router() http.Handler {
	cfg := LoadConfig()
	r := mux.NewRouter()

	// Routes
	r.HandleFunc("/ip", func(w http.ResponseWriter, r *http.Request) {
		w.Write(deviceIp(cfg))
	}).Methods("GET")

	return r
}

