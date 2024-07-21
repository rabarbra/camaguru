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
	var v reflect.Value
	if rv, ok := val.(reflect.Value); ok {
		v = rv
	} else {
		v = reflect.ValueOf(val)
	}
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
		slice := make([]any, v.Len())
		for i := 0; i < v.Len(); i++ {
			slice[i] = v.Index(i).Interface()
		}
		return fmt.Sprintf("(%s)", PrintSlice(slice, ", "))
	default:
		return fmt.Sprintf("%v", val)
	}
}

func Fields(v interface{}) []interface{} {
	val := reflect.ValueOf(v).Elem()
	var fields []interface{}
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := val.Type().Field(i)

		if fieldType.Anonymous {
			for j := 0; j < field.NumField(); j++ {
				fields = append(fields, field.Field(j).Addr().Interface())
			}
		} else {
			fields = append(fields, field.Addr().Interface())
		}
	}
	return fields
}
