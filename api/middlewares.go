package main

import (
	"net/http"
	"orm"
)

type RequestHandler func(http.ResponseWriter, *http.Request)
type DBRequestHandler func(http.ResponseWriter, *http.Request, *orm.Orm)
type ProtectedRequestHandler func(w http.ResponseWriter, r *http.Request, userId int64)
type ProtectedDBRequestHandler func(w http.ResponseWriter, r *http.Request, userId int64, db *orm.Orm)

func congigureCors(w *http.ResponseWriter, r *http.Request) bool {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
	(*w).Header().Set(
		"Access-Control-Allow-Headers",
		`	Access-Control-Allow-Headers,
		Origin,Accept,
		X-Requested-With,
		Content-Type,
		Access-Control-Request-Method,
		Access-Control-Request-Headers,
		Authorization
		`,
	)
	if r.Method == "OPTIONS" {
		(*w).WriteHeader(http.StatusOK)
		return true
	}
	(*w).Header().Set("Content-Type", "application/json")
	return false
}

func CorsM(f RequestHandler) RequestHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		if congigureCors(&w, r) {
			return
		}
		f(w, r)
	}
}

func AuthM(f ProtectedRequestHandler) RequestHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := CheckAuthorized(r)
		if userId == 0 {
			sendError(w, http.StatusBadRequest, "unauthorized")
			return
		}
		f(w, r, userId)
	}
}

func DBM(db *orm.Orm, f DBRequestHandler) RequestHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		f(w, r, db)
	}
}

func AuthDBM(db *orm.Orm, f ProtectedDBRequestHandler) RequestHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := CheckAuthorized(r)
		if userId == 0 {
			sendError(w, http.StatusBadRequest, "unauthorized")
			return
		}
		f(w, r, userId, db)
	}
}
