package main

import (
	"encoding/json"
	"jwt"
	"log"
	"net/http"
	"orm"
	"time"
)

func verifyEmail(w http.ResponseWriter, r *http.Request, db *orm.Orm) {
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
	var user User
	err = db.GetOneById(&user, p.Id)
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
	err = db.Update(&user, user.BaseModel.Id)
	if err != nil {
		log.Println(err)
		http.Redirect(
			w, r,
			EMAIL_VERIFIED_REDIR_URL+"?err=error updating user",
			http.StatusFound,
		)
		return
	}
	accessToken, err := jwt.CreateJWT(JWT_SECRET, user.BaseModel.Id, time.Hour*24)
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

func resetPassword(w http.ResponseWriter, r *http.Request, userId int64, db *orm.Orm) {
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
	var user User
	err := db.GetOneById(&user, userId)
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
	err = db.Update(&user, userId)
	if err != nil {
		sendError(w, http.StatusBadRequest, "error updating user")
		return
	}
	sendJson(w, http.StatusOK, map[string]string{"msg": "password updated successfully!"})
}

type resetPasswordUnauthReq struct {
	Email string `json:"email"`
}

func resetPasswordUnauth(w http.ResponseWriter, r *http.Request, db *orm.Orm) {
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
	var user User
	err := db.GetOne(&user, []orm.Filter{{Key: "email", Value: req.Email, Operation: "="}})
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

func signin(w http.ResponseWriter, r *http.Request, db *orm.Orm) {
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

	var user User
	err := db.GetOne(&user, []orm.Filter{{Key: "username", Value: req.Username, Operation: "="}})
	if err != nil {
		sendError(w, http.StatusBadRequest, "no such user")
		return
	}

	if !user.EmailVerified {
		sendError(w, http.StatusBadRequest, "please, verify your email")
		return
	}

	if CheckPassHash(req.Pass, user.Pass) {
		token, er := jwt.CreateJWT(JWT_SECRET, user.BaseModel.Id, time.Hour*24)
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
