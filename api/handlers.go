package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func sendJson(w http.ResponseWriter, status int, msg map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(msg)
}

func sendError(w http.ResponseWriter, status int, msg string) {
	sendJson(w, status, map[string]string{"err": msg})
}

func signin(w http.ResponseWriter, r *http.Request, db *sql.DB) {

}

func checkUser(db *sql.DB, u User) string {
	if u.Username == "" {
		return "username required"
	} else if checkUserExists(db, "username", u.Username) {
		return "username exists"
	}

	if u.Email == "" {
		return "email required"
	} else if checkUserExists(db, "email", u.Email) {
		return "email exists"
	}

	if u.Pass == "" {
		return "password required"
	}

	return ""
}

func me(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == "GET" {
		fmt.Println(r)
	} else if r.Method == "POST" {
		var u User
		r.Body = http.MaxBytesReader(w, r.Body, 1048576)

		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&u); err != nil {
			sendError(w, http.StatusBadRequest, "error decoding json")
			return
		}

		msg := checkUser(db, u)
		if msg != "" {
			sendError(w, http.StatusBadRequest, msg)
			return
		}

		u.EmailVerified = false
		u.LikeNotify = true
		u.CommNotify = true

		err := createUser(db, u)
		if err != nil {
			log.Println("Error inserting user:", err)
			sendError(w, http.StatusBadRequest, "error creating user")
			return
		}

		sendJson(w, http.StatusCreated, map[string]string{"message": "User created successfully"})
		return
	} else if r.Method == "PUT" {
		var u User
		r.Body = http.MaxBytesReader(w, r.Body, 1048576)

		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&u); err != nil {
			sendError(w, http.StatusBadRequest, "error decoding json")
			return
		}

		if u.Id == 0 {
			sendError(w, http.StatusBadRequest, "id required")
			return
		}

		msg := checkUser(db, u)
		if msg != "" {
			sendError(w, http.StatusBadRequest, msg)
			return
		}

		err := updateUser(db, u)
		if err != nil {
			sendError(w, http.StatusBadRequest, err.Error())
			return
		}
		sendJson(w, http.StatusOK, map[string]string{"message": "User updated successfully"})
	}
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
