package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"jwt"
	"log"
	"net/http"
	"net/smtp"
	"time"
)

type LoginReq struct {
	Username string `json:"username"`
	Pass     string `json:"pass"`
}

func verifyEmail(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Redirect(w, r, EMAIL_VERIFIED_REDIR_URL+"?err=token invalid", http.StatusFound)
		return
	}
	p, err := jwt.VerifyJWT(token, JWT_SECRET)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, EMAIL_VERIFIED_REDIR_URL+"?err=token invalid", http.StatusFound)
		return
	}
	if time.Now().Unix() > p.Exp {
		http.Redirect(w, r, EMAIL_VERIFIED_REDIR_URL+"?err=token expired", http.StatusFound)
		return
	}
	user, err := getUserById(db, p.Id)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, EMAIL_VERIFIED_REDIR_URL+"?err=error getting user", http.StatusFound)
		return
	}
	user.EmailVerified = true
	err = updateUser(db, user)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, EMAIL_VERIFIED_REDIR_URL+"?err=error updating user", http.StatusFound)
		return
	}
	http.Redirect(w, r, EMAIL_VERIFIED_REDIR_URL, http.StatusFound)
}

func sendVerificationEmail(email string, db *sql.DB) error {
	verifyLink := BACKEND_HOST + "/verify?token="
	user, err := getUserByEmail(db, email)
	if err != nil {
		return err
	}
	token, err := jwt.CreateJWT(JWT_SECRET, user.Id, time.Hour*100)
	if err != nil {
		return err
	}
	verifyLink += token
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\nConfirm email: %s\r\n",
		email, "Confirm email", verifyLink))
	auth := smtp.PlainAuth("", SMTP_FROM, SMTP_PASSWORD, SMTP_HOST)
	err = smtp.SendMail(SMTP_HOST+":"+SMTP_PORT, auth, SMTP_FROM, []string{email}, msg)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func signin(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var req LoginReq
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "error decoding json")
		return
	}
	if req.Username == "" {
		sendError(w, http.StatusBadRequest, "username required")
		return
	}
	if req.Pass == "" {
		sendError(w, http.StatusBadRequest, "password required")
		return
	}

	user, err := getUserByUsename(db, req.Username)
	if err != nil {
		sendError(w, http.StatusBadRequest, "no such user")
		return
	}

	if !user.EmailVerified {
		sendError(w, http.StatusBadRequest, "please, verify your email")
		return
	}

	if CheckPassHash(req.Pass, user.Pass) {
		token, er := jwt.CreateJWT(JWT_SECRET, user.Id, time.Hour*24)
		if er != nil {
			sendError(w, http.StatusInternalServerError, "cannot create token")
			return
		}
		sendJson(w, http.StatusOK, map[string]string{"token": token})
		return
	}

	sendError(w, http.StatusBadRequest, "unauthorized")
}

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
