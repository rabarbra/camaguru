package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func sendJson(w http.ResponseWriter, status int, msg map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(msg)
}

func sendJsonBytes(w http.ResponseWriter, status int, msg []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(msg)
}

func sendError(w http.ResponseWriter, status int, msg string) {
	sendJson(w, status, map[string]string{"err": msg})
}

func uploadImg(r *http.Request) (string, error) {
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		log.Println(err, MAX_UPLOAD_SIZE)
		return "", err
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		log.Println(file, err)
		return "", err
	}
	defer file.Close()

	fileExtension := filepath.Ext(fileHeader.Filename)
	if fileExtension != ".jpg" && fileExtension != ".jpeg" && fileExtension != ".png" {
		return "", err
	}

	if _, err := os.Stat(UPLOAD_DIR); os.IsNotExist(err) {
		os.MkdirAll(UPLOAD_DIR, 0755)
	}

	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), fileExtension)
	filePath := filepath.Join(UPLOAD_DIR, fileName)

	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := file.Seek(0, 0); err != nil {
		return "", err
	}
	if _, err := dst.ReadFrom(file); err != nil {
		return "", err
	}
	return filePath, nil
}
