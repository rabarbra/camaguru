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

type FilterOperation string

const (
	OpEqual        FilterOperation = "="
	OpNotEqual     FilterOperation = "!="
	OpGreaterThan  FilterOperation = ">"
	OpLessThan     FilterOperation = "<"
	OpGreaterEqual FilterOperation = ">="
	OpLessEqual    FilterOperation = "<="
)

type Filter struct {
	Key       string
	Value     any
	Operation FilterOperation
}

type OrderDirection string

const (
	ASC  OrderDirection = "ASC"
	DESC OrderDirection = "DESC"
)

type Order struct {
	Key       string
	Direction OrderDirection
}

type Pagination struct {
	limit  int
	offset int
}

func buildPagination(pag Pagination) string {
	return fmt.Sprintf(" LIMIT %d OFFSET %d", pag.limit, pag.offset)
}

func buildQuery(baseQuery string, filters []Filter) (string, []interface{}) {
	var whereClauses []string
	var args []interface{}

	for i, filter := range filters {
		whereClauses = append(whereClauses, fmt.Sprintf("%s %s $%d", filter.Key, filter.Operation, i+1))
		args = append(args, filter.Value)
	}

	query := baseQuery
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	return query, args
}

func getMany[T any](r *[]T, db *sql.DB, filters []Filter, pag Pagination) error {
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
	sqlReq += fmt.Sprintf("FROM %s", tableName)
	query, args := buildQuery(sqlReq, filters)
	query += buildPagination(pag)

	rows, err := db.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var item T
		if err := rows.Scan(Fields(item)...); err != nil {
			return err
		}
		*r = append(*r, item)
	}
	return nil
}
