package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"orm"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

func parseFilters(r *http.Request) []orm.Filter {
	var filters []orm.Filter
	query := r.URL.Query()
	for key, values := range query {
		if strings.HasPrefix(key, "[") && strings.HasSuffix(key, "]") {
			var op orm.FilterOperation
			field := strings.TrimPrefix(key, "[")
			field = strings.TrimSuffix(field, "]")
			if strings.HasSuffix(field, ">>") {
				op = ">="
				field = strings.TrimSuffix(field, ">>")
			} else if strings.HasSuffix(field, "<<") {
				op = "<="
				field = strings.TrimSuffix(field, "<<")
			} else if strings.HasSuffix(field, "<>") {
				op = "<>"
				field = strings.TrimSuffix(field, "<>")
			} else if strings.HasSuffix(field, ">") {
				op = ">"
				field = strings.TrimSuffix(field, ">")
			} else if strings.HasSuffix(field, "<") {
				op = "<"
				field = strings.TrimSuffix(field, "<")
			} else if strings.HasSuffix(field, "~") {
				op = "LIKE"
				field = strings.TrimSuffix(field, "~")
			} else {
				op = "="
			}
			for _, val := range values {
				var v any
				v = val
				if strings.Contains(val, ",") && op == "=" {
					op = "IN"
					v = strings.Split(val, ",")
				}
				filters = append(filters, orm.Filter{
					Key:       field,
					Operation: op,
					Value:     v,
				})
			}
		}
	}
	return filters
}

func parseSort(r *http.Request) []orm.Sort {
	var sorts []orm.Sort
	sortParam := r.URL.Query().Get("sort")
	if sortParam != "" {
		fields := strings.Split(sortParam, ",")
		for _, field := range fields {
			var direction orm.OrderDirection
			direction = "ASC"
			if strings.HasPrefix(field, "-") {
				direction = "DESC"
				field = strings.TrimPrefix(field, "-")
			}
			sorts = append(sorts, orm.Sort{
				Key:       field,
				Direction: direction,
			})
		}
	}
	return sorts
}

func parsePagination(r *http.Request) (orm.Pagination, error) {
	offsetParam := r.URL.Query().Get("offset")
	limitParam := r.URL.Query().Get("lim")

	var offset, limit int
	var err error

	if offsetParam != "" {
		offset, err = strconv.Atoi(offsetParam)
		if err != nil {
			return orm.Pagination{}, fmt.Errorf("invalid page parameter")
		}
	}

	if limitParam != "" {
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			return orm.Pagination{}, fmt.Errorf("invalid per_page parameter")
		}
	}

	return orm.Pagination{
		Limit:  limit,
		Offset: offset,
	}, nil
}
