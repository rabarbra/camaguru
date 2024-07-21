package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"orm"
)

func GetOneHandler(w http.ResponseWriter, r *http.Request, model orm.Model, db *orm.Orm, userId int64, id int64) {
	err := db.GetOneById(model, id)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
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
	}
	id, err := db.Create(model)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
	}
	sendJson(w, http.StatusCreated, map[string]string{"id": fmt.Sprintf("%d", id)})
	w.WriteHeader(http.StatusCreated)
}

func GetManyHandler(w http.ResponseWriter, r *http.Request, model orm.Model, db *orm.Orm, userId int64) {
	// var u []User
	u, err := db.GetMany(model, nil, nil, orm.Pagination{})
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
	}
	msg, er := json.Marshal(u)
	if er != nil {
		sendError(w, http.StatusBadRequest, "error marshaling")
		return
	}
	sendJsonBytes(w, http.StatusOK, msg)
}
