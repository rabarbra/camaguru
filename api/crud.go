package main

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func create(r any, db *sql.DB) error {
	tableName := ToSnakeCase(reflect.TypeOf(r).Name()) + "s"
	sqlReq := "INSERT INTO " + tableName + "("
	sqlValues := "VALUES("
	values := reflect.ValueOf(r)
	types := values.Type()
	for i := 0; i < values.NumField(); i++ {
		if types.Field(i).Name == "Id" {
			continue
		}
		sqlReq += ToSnakeCase(types.Field(i).Name)
		format := "'%s'"
		switch values.Field(i).Type().Kind() {
		case reflect.Int64, reflect.Int32, reflect.Int:
			format = "%d"
		case reflect.Float64, reflect.Float32:
			format = "%f"
		case reflect.Bool:
			format = "%t"
		}
		sqlValues += fmt.Sprintf(format, values.Field(i))
		if i < values.NumField()-1 {
			sqlReq += ", "
			sqlValues += ", "
		} else {
			sqlReq += ") "
			sqlValues += ");"
		}
	}
	res, err := db.Exec(sqlReq + sqlValues)
	log.Println(res, err)
	return err
}

func delete(r any, db *sql.DB, id int64) error {
	tableName := ToSnakeCase(reflect.TypeOf(r).Name()) + "s"
	sqlReq := fmt.Sprintf("DELETE FROM %s WHERE id = %d", tableName, id)
	_, err := db.Exec(sqlReq)
	return err
}

func Fields(v interface{}) []interface{} {
	val := reflect.ValueOf(v).Elem()
	fields := make([]interface{}, val.NumField())
	for i := 0; i < val.NumField(); i++ {
		fields[i] = val.Field(i).Addr().Interface()
	}
	return fields
}

func get[T any](r *T, db *sql.DB, id int64) error {
	tableName := ToSnakeCase(reflect.TypeOf(r).Elem().Name()) + "s"
	sqlReq := "SELECT "
	values := reflect.ValueOf(r).Elem()
	types := values.Type()
	for i := 0; i < values.NumField(); i++ {
		sqlReq += ToSnakeCase(types.Field(i).Name)
		if i < values.NumField()-1 {
			sqlReq += ", "
		} else {
			sqlReq += " "
		}
	}
	sqlReq += fmt.Sprintf("FROM %s WHERE id = %d", tableName, id)
	err := db.QueryRow(sqlReq).Scan(Fields(r)...)
	return err
}
