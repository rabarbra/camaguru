package main

import (
	"net/http"
	"orm"
	"reflect"
	"strconv"
)

func handleRequest(model orm.Model) ProtectedDBRequestHandler {
	return func(w http.ResponseWriter, r *http.Request, userId int64, db *orm.Orm) {
		switch r.Method {
		case http.MethodGet:
			GetManyHandler(w, r, model, db, userId)
		case http.MethodPost, http.MethodPut:
			CreateHandler(w, r, model, db, userId)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func handleRequestId(model orm.Model) ProtectedDBRequestHandler {
	return func(w http.ResponseWriter, r *http.Request, userId int64, db *orm.Orm) {
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
		case http.MethodDelete:
			DeleteHandler(w, r, model, db, userId, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func AddCrudRoutes(model orm.Model, db *orm.Orm) {
	routeName := orm.ToSnakeCase(reflect.TypeOf(model).Name())
	handler := CorsM(AuthDBM(db, handleRequest(model)))
	handlerId := CorsM(AuthDBM(db, handleRequestId(model)))
	http.HandleFunc("/"+routeName, handler)
	http.HandleFunc("/"+routeName+"/", handler)
	http.HandleFunc("/"+routeName+"/{id}", handlerId)
	http.HandleFunc("/"+routeName+"/{id}/", handlerId)
}
