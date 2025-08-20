package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type API struct {
	DB *AppDB
}

func (api *API) getMessages(w http.ResponseWriter, r *http.Request) {
	messages, err := api.DB.GetAllMessages()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func (api *API) postMessage(w http.ResponseWriter, r *http.Request) {
	var msg Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if msg.Content == "" {
		http.Error(w, "content is required", http.StatusBadRequest)
		return
	}

	if err := api.DB.InsertMessage("api-client", msg.Content); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"status":"ok"}`))
}

func NewRouter(db *AppDB) http.Handler {
	api := &API{DB: db}
	r := mux.NewRouter()

	// Routes
	r.HandleFunc("/messages", api.getMessages).Methods("GET")
	r.HandleFunc("/messages", api.postMessage).Methods("POST")

	return r
}
