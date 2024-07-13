package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func checkUser(db *sql.DB, u User) string {
	if u.Username == "" {
		return "username required"
	} else if checkUserExists(db, "username", u.Username, u.Id) {
		return "username exists"
	}

	if u.Email == "" {
		return "email required"
	} else if checkUserExists(db, "email", u.Email, u.Id) {
		return "email exists"
	}

	if u.Pass == "" {
		return "password required"
	}

	return ""
}

func postUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var u User
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&u); err != nil {
		sendError(w, http.StatusBadRequest, "error decoding json")
		return
	}

	u.Id = 0
	msg := checkUser(db, u)
	if msg != "" {
		sendError(w, http.StatusBadRequest, msg)
		return
	}

	u.EmailVerified = false
	u.LikeNotify = true
	u.CommNotify = true

	hash, er := HashPass(u.Pass)
	if er != nil {
		sendError(w, http.StatusBadRequest, "error hashing password")
		return
	}

	u.Pass = hash

	err := createUser(db, u)
	if err != nil {
		log.Println("Error inserting user:", err)
		sendError(w, http.StatusBadRequest, "error creating user")
		return
	}

	sendVerificationEmail(u.Email, db)
	sendJson(w, http.StatusCreated, map[string]string{"message": "User created successfully"})
}

func getUser(w http.ResponseWriter, r *http.Request, userId int64, db *sql.DB) {
	user, err := getUserById(db, userId)
	if err != nil {
		sendError(w, http.StatusBadRequest, "user not found")
		return
	}
	msg, er := json.Marshal(user)
	if er != nil {
		sendError(w, http.StatusBadRequest, "error marshaling")
		return
	}
	sendJsonBytes(w, http.StatusOK, msg)
}

func putUser(w http.ResponseWriter, r *http.Request, userId int64, db *sql.DB) {
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

	if u.Id != userId {
		sendError(w, http.StatusBadRequest, "unauthorized")
		return
	}

	msg := checkUser(db, u)
	if msg != "" {
		sendError(w, http.StatusBadRequest, msg)
		return
	}

	hash, er := HashPass(u.Pass)
	if er != nil {
		sendError(w, http.StatusBadRequest, "error hashing password")
		return
	}

	u.Pass = hash

	err := updateUser(db, u)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}
	sendJson(w, http.StatusOK, map[string]string{"message": "User updated successfully"})
}
