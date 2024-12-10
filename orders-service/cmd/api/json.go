package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func readJSON(r *http.Request, d interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(d)
}

func writeJSON(w http.ResponseWriter, jsonData any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(jsonData)
	if err != nil {
		log.Println("error encoding data: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
