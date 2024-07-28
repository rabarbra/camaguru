package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"orm"
	"reflect"
)

func DeleteHandler(w http.ResponseWriter, r *http.Request, model orm.Model, db *orm.Orm, userId int64, id int64) {
	err := db.GetOneById(model, id)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}
	ownerId := reflect.ValueOf(model).Elem().FieldByName("user_id").Int()
	if ownerId != 0 && ownerId != userId {
		sendError(w, http.StatusForbidden, "not yours")
		return
	}
	err = db.Delete(model, id)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}
	sendJson(w, http.StatusOK, map[string]string{"msg": "deleted"})
}

func GetOneHandler(w http.ResponseWriter, r *http.Request, model orm.Model, db *orm.Orm, userId int64, id int64) {
	err := db.GetOneById(model, id)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}
	msg, er := json.Marshal(model)
	if er != nil {
		sendError(w, http.StatusBadRequest, "error marshaling")
		return
	}
	sendJsonBytes(w, http.StatusOK, msg)
}

func CreateHandler(w http.ResponseWriter, r *http.Request, model orm.Model, db *orm.Orm, userId int64) {
	err := model.NewItem(r, userId)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}
	id, err := db.Create(model)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}
	sendJson(w, http.StatusCreated, map[string]string{"id": fmt.Sprintf("%d", id)})
	w.WriteHeader(http.StatusCreated)
}

func GetManyHandler(w http.ResponseWriter, r *http.Request, model orm.Model, db *orm.Orm, userId int64) {
	filters := parseFilters(r)
	sorts := parseSort(r)
	pagination, err := parsePagination(r)
	if err != nil {
		sendError(w, http.StatusBadRequest, "invalid pagination parameters")
		return
	}
	u, err := db.GetMany(model, filters, sorts, pagination)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}
	msg, er := json.Marshal(u)
	if er != nil {
		sendError(w, http.StatusBadRequest, "error marshaling")
		return
	}
	sendJsonBytes(w, http.StatusOK, msg)
}
