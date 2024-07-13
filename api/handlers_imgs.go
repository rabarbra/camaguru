package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func getMaxUploadSize() int64 {
	var maxUploadSize, err = strconv.Atoi(os.Getenv("MAX_UPLOAD_SIZE"))
	if err != nil {
		maxUploadSize = 10 * 1024 * 1024
	}
	return int64(maxUploadSize)
}

var (
	maxUploadSize = getMaxUploadSize()
	uploadDir     = os.Getenv("UPLOAD_DIRECTORY")
)

func postImg(w http.ResponseWriter, r *http.Request, userId int64, db *sql.DB) {
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		log.Println(err, maxUploadSize)
		sendError(w, http.StatusBadRequest, "The uploaded file is too big.")
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		log.Println(file, err)
		sendError(w, http.StatusBadRequest, "Couldn't retrieve file.")
		return
	}
	defer file.Close()

	fileExtension := filepath.Ext(fileHeader.Filename)
	if fileExtension != ".jpg" && fileExtension != ".jpeg" && fileExtension != ".png" {
		sendError(w, http.StatusBadRequest, "The provided file format is not allowed. Please upload a JPEG or PNG image.")
		return
	}

	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, 0755)
	}

	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), fileExtension)
	filePath := filepath.Join(uploadDir, fileName)

	dst, err := os.Create(filePath)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Unable to create the file for writing. Check your write access privilege.")
		return
	}
	defer dst.Close()

	if _, err := file.Seek(0, 0); err != nil {
		sendError(w, http.StatusInternalServerError, "Unable to read the file.")
		return
	}
	if _, err := dst.ReadFrom(file); err != nil {
		sendError(w, http.StatusInternalServerError, "Unable to save the file.")
		return
	}

	var img Img
	img.Link = filePath
	img.UserId = userId
	err = createImg(db, img)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to save file information to database.")
	}

	sendJson(w, http.StatusCreated, map[string]string{"message": "Img created successfully"})
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
