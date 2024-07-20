package orm

import (
	"bytes"
	"fmt"
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

func PrintSlice(values []any, separator string) string {
	var buffer bytes.Buffer
	for i, val := range values {
		buffer.WriteString(ToSQL(val))
		if i < len(values)-1 {
			buffer.WriteString(separator)
		}
	}
	return buffer.String()
}

func ToSQL(val any) string {
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.Bool:
		return fmt.Sprintf("%t", val)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", val)
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%f", val)
	case reflect.String:
		return fmt.Sprintf("'%s'", val)
	case reflect.Slice:
		return fmt.Sprintf("(%s)", PrintSlice(val.([]any), ", "))
	default:
		return fmt.Sprintf("%v", val)
	}
}

func Fields(v interface{}) []interface{} {
	val := reflect.ValueOf(v).Elem()
	fields := make([]interface{}, val.NumField())
	for i := 0; i < val.NumField(); i++ {
		fields[i] = val.Field(i).Addr().Interface()
	}
	return fields
}
