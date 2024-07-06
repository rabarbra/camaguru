package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func signin(w http.ResponseWriter, r *http.Request, db *sql.DB) {

}

func me(w http.ResponseWriter, r *http.Request, db *sql.DB) {

}

func img(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	data := map[string]string{"status": "ok"}
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
