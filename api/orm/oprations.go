package orm

import (
	"fmt"
	"reflect"
)

func getModelFields(model BaseModel, withId bool) string {
	sqlReq := ""
	values := reflect.ValueOf(model)
	types := values.Type()
	for i := 0; i < values.NumField(); i++ {
		val := ToSnakeCase(types.Field(i).Name)
		if !withId && val == "id" {
			continue
		}
		sqlReq += val
		if i < values.NumField()-1 {
			sqlReq += ", "
		}
	}
	return sqlReq
}

func (o Orm) GetOne(model BaseModel, filter []Filter) (BaseModel, error) {
	tableName := ToSnakeCase(reflect.TypeOf(model).Name()) + "s"
	err := o.db.QueryRow(
		fmt.Sprintf("SELECT %s FROM %s %s;", getModelFields(model, true), tableName, filter[0].ToSQL(filter...)),
	).Scan(Fields(model)...)
	return model, err
}

func (o Orm) GetOneById(model BaseModel, id int64) (BaseModel, error) {
	tableName := ToSnakeCase(reflect.TypeOf(model).Name()) + "s"
	err := o.db.QueryRow(
		fmt.Sprintf("SELECT %s FROM %s WHERE id = %d;", getModelFields(model, true), tableName, id),
	).Scan(Fields(model)...)
	return model, err
}

func (o Orm) GetMany(model BaseModel, filter []Filter, sort []Sort, pagination Pagination) ([]BaseModel, error) {
	tableName := ToSnakeCase(reflect.TypeOf(model).Name()) + "s"
	rows, err := o.db.Query(
		fmt.Sprintf("SELECT %s FROM %s %s %s %s;",
			getModelFields(model, true),
			tableName,
			filter[0].ToSQL(filter...),
			sort[0].ToSQL(sort...),
			pagination.ToSQL(),
		),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var r []BaseModel
	for rows.Next() {
		if err := rows.Scan(Fields(model)...); err != nil {
			return nil, err
		}
		r = append(r, model)
	}
	return r, nil
}

func (o Orm) Create(model BaseModel) (BaseModel, error) {
	tableName := ToSnakeCase(reflect.TypeOf(model).Name()) + "s"
	sqlValues := ""
	values := reflect.ValueOf(model)
	types := values.Type()
	for i := 0; i < values.NumField(); i++ {
		if types.Field(i).Name == "Id" {
			continue
		}
		sqlValues += ToSQL(values.Field(i))
		if i < values.NumField()-1 {
			sqlValues += ", "
		}
	}
	var newId int64
	err := o.db.QueryRow(
		fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s) RETURNING id;", tableName, getModelFields(model, false), sqlValues),
	).Scan(&newId)
	model.Id = newId
	return model, err
}

func (o Orm) Update(model BaseModel, id int64) (BaseModel, error) {
	tableName := ToSnakeCase(reflect.TypeOf(model).Name()) + "s"
	sqlValues := ""
	values := reflect.ValueOf(model)
	types := values.Type()
	for i := 0; i < values.NumField(); i++ {
		if types.Field(i).Name == "Id" {
			continue
		}
		sqlValues += fmt.Sprintf("%s = %s", ToSnakeCase(types.Field(i).Name), ToSQL(values.Field(i)))
		if i < values.NumField()-1 {
			sqlValues += ", "
		}
	}
	_, err := o.db.Exec(
		fmt.Sprintf("UPDATE %s SET %s WHERE id = %d;", tableName, getModelFields(model, false), id),
	)
	return model, err
}

func (o Orm) Delete(model BaseModel, id int64) error {
	tableName := ToSnakeCase(reflect.TypeOf(model).Name()) + "s"
	sqlReq := fmt.Sprintf("DELETE FROM %s WHERE id = %d", tableName, id)
	_, err := o.db.Exec(sqlReq)
	return err
}
