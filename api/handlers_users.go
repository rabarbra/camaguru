package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
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

	_, err := create(u, db)
	// err := createUser(db, u)
	if err != nil {
		log.Println("Error inserting user:", err)
		sendError(w, http.StatusBadRequest, "error creating user")
		return
	}

	err = sendVerificationEmail(u.Email, db)
	if err != nil {
		log.Println(err)
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendJson(w, http.StatusCreated, map[string]string{"msg": "User created successfully"})
}

func getUser(w http.ResponseWriter, r *http.Request, userId int64, db *sql.DB) {
	var user User
	err := get(&user, db, userId)
	// user, err := getUserById(db, userId)
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
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if updates["pass"] != nil {
		hash, er := HashPass(updates["pass"].(string))
		if er != nil {
			sendError(w, http.StatusBadRequest, "error hashing password")
			return
		}
		updates["pass"] = hash
	}

	err := partUpdateUser(db, userId, updates)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}
	sendJson(w, http.StatusOK, map[string]string{"msg": "User updated successfully"})
}

func postUserAvatar(w http.ResponseWriter, r *http.Request, userId int64, db *sql.DB) {
	u, err := getUserById(db, userId)
	if err != nil {
		log.Println(err)
		sendError(w, http.StatusBadRequest, "no such user")
		return
	}
	filePath, err := uploadImg(r)
	if err != nil {
		log.Println(err)
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}
	u.Avatar = strings.TrimPrefix(filePath, "assets")
	err = updateUser(db, u)
	if err != nil {
		log.Println(err)
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendJson(w, http.StatusOK, map[string]string{"msg": "Avatar updated successfully"})
}
