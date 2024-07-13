package main

import (
	"os"
	"strconv"
)

func getMaxUploadSize() int64 {
	var maxUploadSize, err = strconv.Atoi(os.Getenv("MAX_UPLOAD_SIZE"))
	if err != nil {
		maxUploadSize = 10 * 1024 * 1024
	}
	return int64(maxUploadSize)
}

var (
	MAX_UPLOAD_SIZE = getMaxUploadSize()
	UPLOAD_DIR      = os.Getenv("UPLOAD_DIRECTORY")

	JWT_SECRET               = os.Getenv("JWT_SECRET")
	EMAIL_VERIFIED_REDIR_URL = os.Getenv("EMAIL_VERIFIED_REDIR_URL")
	SMTP_FROM                = os.Getenv("SMTP_FROM")
	SMTP_PASSWORD            = os.Getenv("SMTP_PASSWORD")
	SMTP_HOST                = os.Getenv("SMTP_HOST")
	SMTP_PORT                = os.Getenv("SMTP_PORT")
	BACKEND_HOST             = "http://localhost:8000"
)
