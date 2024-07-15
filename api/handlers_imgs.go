package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func postImg(w http.ResponseWriter, r *http.Request, userId int64, db *sql.DB) {
	filePath, err := uploadImg(r)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}
	var img Img
	img.Link = strings.TrimPrefix(filePath, "assets")
	img.UserId = userId
	err = createImg(db, img)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to save file information to database.")
	}
	sendJson(w, http.StatusCreated, map[string]string{"msg": "Img created successfully"})
}

func getImg(w http.ResponseWriter, r *http.Request, userId int64, db *sql.DB) {
	q := r.URL.Query()
	limit := 10
	offset := 0
	if q.Has("limit") {
		var err error = nil
		limit, err = strconv.Atoi(q.Get("limit"))
		if err != nil {
			sendError(w, http.StatusInternalServerError, "error")
			return
		}
	}
	if q.Has("offset") {
		var err error = nil
		offset, err = strconv.Atoi(q.Get("offset"))
		if err != nil {
			sendError(w, http.StatusInternalServerError, "error")
			return
		}
	}
	imgs, er := getImgs(db, userId, limit, offset)
	if er != nil {
		sendError(w, http.StatusInternalServerError, "error getting imgs")
		return
	}
	msg, er := json.Marshal(imgs)
	if er != nil {
		sendError(w, http.StatusBadRequest, "error marshaling")
		return
	}
	sendJsonBytes(w, http.StatusOK, msg)
}
