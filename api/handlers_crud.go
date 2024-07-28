package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"orm"
)

func DeleteHandler(w http.ResponseWriter, r *http.Request, model orm.Model, db *orm.Orm, userId int64, id int64) {
	err := db.Delete(model, id)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}
	sendJson(w, http.StatusOK, map[string]string{"msg": "deleted"})
}

func UpdateHandler(w http.ResponseWriter, r *http.Request, model orm.Model, db *orm.Orm, userId int64, id int64) {
	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	err := db.Update(model, id)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}
	sendJson(w, http.StatusOK, map[string]string{"msg": "updated"})
}

func PatchHandler(w http.ResponseWriter, r *http.Request, model orm.Model, db *orm.Orm, userId int64, id int64) {
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	err := db.Patch(model, id, updates)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}
	sendJson(w, http.StatusOK, map[string]string{"msg": "updated"})
}

func GetOneHandler(w http.ResponseWriter, r *http.Request, model orm.Model, db *orm.Orm, userId int64, id int64) {
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
