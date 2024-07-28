package main

import (
	"database/sql"
	"log"
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
			log.Println("id handler: ", err)
			http.NotFound(w, r)
			return
		}
		err = db.GetOneById(model, id)
		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			sendError(w, http.StatusBadRequest, err.Error())
			return
		}
		if r.Method == http.MethodGet {
			GetOneHandler(w, r, model, db, userId, id)
			return
		}
		field := reflect.ValueOf(model).Elem().FieldByName("UserId")
		if field.Kind() != 0 {
			ownerId := field.Int()
			if ownerId != 0 && ownerId != userId {
				sendError(w, http.StatusForbidden, "not yours")
				return
			}
		}
		switch r.Method {
		case http.MethodPut:
			UpdateHandler(w, r, model, db, userId, id)
		case http.MethodPatch:
			PatchHandler(w, r, model, db, userId, id)
		case http.MethodDelete:
			DeleteHandler(w, r, model, db, userId, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func AddCrudRoutes(model orm.Model, db *orm.Orm) {
	routeName := orm.ToSnakeCase(reflect.TypeOf(model).Elem().Name())
	log.Printf("Adding CRUD routes for /%s", routeName)
	handler := CorsM(AuthDBM(db, handleRequest(model)))
	handlerId := CorsM(AuthDBM(db, handleRequestId(model)))
	http.HandleFunc("/"+routeName, handler)
	http.HandleFunc("/"+routeName+"/", handler)
	http.HandleFunc("/"+routeName+"/{id}", handlerId)
	http.HandleFunc("/"+routeName+"/{id}/", handlerId)
}
