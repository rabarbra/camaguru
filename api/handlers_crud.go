package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetOneHandler(w http.ResponseWriter, r *http.Request, model BaseModel, db *sql.DB, userId int64, id int64) {
	err := get(&model, db, id)
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

func CreateHandler(w http.ResponseWriter, r *http.Request, model BaseModel, db *sql.DB, userId int64) {
	model, err := model.NewItem(r, userId)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
	}
	id, err := create(model, db)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
	}
	sendJson(w, http.StatusCreated, map[string]string{"id": fmt.Sprintf("%d", id)})
	w.WriteHeader(http.StatusCreated)
}
