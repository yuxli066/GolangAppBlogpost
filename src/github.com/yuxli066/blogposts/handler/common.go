package handler

import (
	"encoding/json"
	"net/http"
)

// The responsdJson function ensures all responses are in JSON format
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

// The respondError function ensures all errors are returned in JSON format
func respondError(w http.ResponseWriter, code int, message string) {
	respondJSON(w, code, map[string]string{"error": message})
}

// The healthCheck function performs a basic heart beat health check on this service
func healthCheck(w http.ResponseWriter, status int, payload interface{}) {
	respondJSON(w, status, map[string]bool{"success": status == 200})
}
