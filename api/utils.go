package main

import (
	"encoding/json"
	"net/http"
)

func sendJson(w http.ResponseWriter, status int, msg map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(msg)
}

func sendJsonBytes(w http.ResponseWriter, status int, msg []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(msg)
}

func sendError(w http.ResponseWriter, status int, msg string) {
	sendJson(w, status, map[string]string{"err": msg})
}
