package main

import (
	"database/sql"
	"encoding/json"
	"jwt"
	"log"
	"net/http"
	"time"
)

func verifyEmail(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Redirect(
			w, r,
			EMAIL_VERIFIED_REDIR_URL+"?err=token invalid",
			http.StatusFound,
		)
		return
	}
	p, err := jwt.VerifyJWT(token, JWT_SECRET)
	if err != nil {
		log.Println(err)
		http.Redirect(
			w, r,
			EMAIL_VERIFIED_REDIR_URL+"?err=token invalid",
			http.StatusFound,
		)
		return
	}
	if time.Now().Unix() > p.Exp {
		http.Redirect(
			w, r,
			EMAIL_VERIFIED_REDIR_URL+"?err=token expired",
			http.StatusFound,
		)
		return
	}
	user, err := getUserById(db, p.Id)
	if err != nil {
		log.Println(err)
		http.Redirect(
			w, r,
			EMAIL_VERIFIED_REDIR_URL+"?err=error getting user",
			http.StatusFound,
		)
		return
	}
	user.EmailVerified = true
	err = updateUser(db, user)
	if err != nil {
		log.Println(err)
		http.Redirect(
			w, r,
			EMAIL_VERIFIED_REDIR_URL+"?err=error updating user",
			http.StatusFound,
		)
		return
	}
	accessToken, err := jwt.CreateJWT(JWT_SECRET, user.Id, time.Hour*24)
	if err != nil {
		log.Println(err)
		http.Redirect(
			w, r,
			EMAIL_VERIFIED_REDIR_URL+"?err=error authorizing user",
			http.StatusFound,
		)
		return
	}
	http.Redirect(
		w, r,
		EMAIL_VERIFIED_REDIR_URL+"?token="+accessToken,
		http.StatusFound,
	)
}

type resetPasswordReq struct {
	Pass string `json:"pass"`
}

func resetPassword(w http.ResponseWriter, r *http.Request, userId int64, db *sql.DB) {
	var req resetPasswordReq
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "error decoding json")
		return
	}
	if req.Pass == "" {
		sendError(w, http.StatusBadRequest, "password required")
		return
	}
	user, err := getUserById(db, userId)
	if err != nil {
		sendError(w, http.StatusBadRequest, "no such user")
		return
	}
	if !user.EmailVerified {
		user.EmailVerified = true
	}
	hash, er := HashPass(req.Pass)
	if er != nil {
		sendError(w, http.StatusBadRequest, "error hashing password")
		return
	}

	user.Pass = hash
	err = updateUser(db, user)
	if err != nil {
		sendError(w, http.StatusBadRequest, "error updating user")
		return
	}
	sendJson(w, http.StatusOK, map[string]string{"msg": "password updated successfully!"})
}

type resetPasswordUnauthReq struct {
	Email string `json:"email"`
}

func resetPasswordUnauth(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var req resetPasswordUnauthReq
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "error decoding json")
		return
	}
	if req.Email == "" {
		sendError(w, http.StatusBadRequest, "email required")
		return
	}
	user, err := getUserByEmail(db, req.Email)
	if err != nil {
		sendError(w, http.StatusBadRequest, "no such user")
		return
	}
	err = sendResetPasswordEmail(user.Email, db)
	if err != nil {
		sendError(w, http.StatusBadRequest, "error sending email")
		return
	}
	sendJson(w, http.StatusOK, map[string]string{"msg": "email sent"})
}

type LoginReq struct {
	Username string `json:"username"`
	Pass     string `json:"pass"`
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
		log.Println(er)
		if er != nil {
			sendError(w, http.StatusInternalServerError, "cannot create token")
			return
		}
		sendJson(w, http.StatusOK, map[string]string{"token": token})
		return
	}

	sendError(w, http.StatusBadRequest, "unauthorized")
}
