package main

import (
	"database/sql"
	"net/http"
	"reflect"
	"strconv"
)

func handleRequest(model *BaseModel) ProtectedDBRequestHandler {
	return func(w http.ResponseWriter, r *http.Request, userId int64, db *sql.DB) {
		switch r.Method {
		// case http.MethodGet:
		// 	GetManyHandler(w, r, model, db, userId)
		case http.MethodPost, http.MethodPut:
			CreateHandler(w, r, model, db, userId)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func handleRequestId(model *BaseModel) ProtectedDBRequestHandler {
	return func(w http.ResponseWriter, r *http.Request, userId int64, db *sql.DB) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 0)
		if err != nil || id < 1 {
			http.NotFound(w, r)
			return
		}
		switch r.Method {
		case http.MethodGet:
			GetOneHandler(w, r, model, db, userId, id)
		// case http.MethodPut:
		// 	UpdateHandler(w, r, model, db, userId, id)
		// case http.MethodPatch:
		// 	PatchHandler(w, r, model, db, userId, id)
		// case http.MethodDelete:
		// 	DeleteHandler(w, r, model, db, userId, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func AddCrudRoutes(model *BaseModel, db *sql.DB) {
	routeName := ToSnakeCase(reflect.TypeOf(model).Elem().Name())
	handler := CorsM(AuthDBM(db, handleRequest(model)))
	handlerId := CorsM(AuthDBM(db, handleRequestId(model)))
	http.HandleFunc("/"+routeName, handler)
	http.HandleFunc("/"+routeName+"/", handler)
	http.HandleFunc("/"+routeName+"/{id}", handlerId)
	http.HandleFunc("/"+routeName+"/{id}/", handlerId)
}
