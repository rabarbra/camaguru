package orm

import (
	"database/sql"
	"fmt"
	"reflect"
)

func getModelFields(model Model, withId bool) string {
	sqlReq := ""
	values := reflect.ValueOf(model).Elem()
	types := values.Type()
	for i := 0; i < values.NumField(); i++ {
		val := ToSnakeCase(types.Field(i).Name)
		if val == "base_model" {
			val = "id"
		}
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

func (o *Orm) GetOne(model Model, filter []Filter) error {
	filters := ""
	if len(filter) != 0 {
		filters = filter[0].ToSQL(filter...)
	}
	err := o.db.QueryRow(
		fmt.Sprintf("SELECT %s FROM %s %s;",
			getModelFields(model, true),
			model.TableName(),
			filters,
		),
	).Scan(Fields(model)...)
	return err
}

func (o *Orm) GetOneById(model Model, id int64) error {
	err := o.db.QueryRow(
		fmt.Sprintf("SELECT %s FROM %s WHERE id = %d;",
			getModelFields(model, true),
			model.TableName(),
			id,
		),
	).Scan(Fields(model)...)
	return err
}

func (o *Orm) GetMany(
	model Model,
	filter []Filter,
	sort []Sort,
	pagination Pagination,
) ([]Model, error) {
	filters := ""
	if len(filter) != 0 {
		filters = filter[0].ToSQL(filter...)
	}
	sorts := ""
	if len(sort) != 0 {
		sorts = sort[0].ToSQL(sort...)
	}
	rows, err := o.db.Query(
		fmt.Sprintf("SELECT %s FROM %s %s %s %s;",
			getModelFields(model, true),
			model.TableName(),
			filters,
			sorts,
			pagination.ToSQL(),
		),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var r []Model
	for rows.Next() {
		newModel := reflect.New(reflect.TypeOf(model).Elem()).Interface().(Model)
		if err := rows.Scan(Fields(newModel)...); err != nil {
			return nil, err
		}
		r = append(r, newModel)
	}
	return r, nil
}

func (o *Orm) Create(model Model) (int64, error) {
	sqlValues := ""
	values := reflect.ValueOf(model).Elem()
	types := values.Type()
	for i := 0; i < values.NumField(); i++ {
		if types.Field(i).Name == "BaseModel" {
			continue
		}
		sqlValues += ToSQL(values.Field(i))
		if i < values.NumField()-1 {
			sqlValues += ", "
		}
	}
	var newId int64
	err := o.db.QueryRow(
		fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s) RETURNING id;",
			model.TableName(),
			getModelFields(model, false),
			sqlValues,
		),
	).Scan(&newId)
	values.FieldByName("BaseModel").FieldByName("Id").SetInt(newId)
	return newId, err
}

func (o *Orm) Update(model Model, id int64) error {
	sqlValues := ""
	values := reflect.ValueOf(model).Elem()
	types := values.Type()
	for i := 0; i < values.NumField(); i++ {
		if types.Field(i).Name == "BaseModel" {
			continue
		}
		sqlValues += fmt.Sprintf("%s = %s",
			ToSnakeCase(types.Field(i).Name),
			ToSQL(values.Field(i)),
		)
		if i < values.NumField()-1 {
			sqlValues += ", "
		}
	}
	_, err := o.db.Exec(
		fmt.Sprintf("UPDATE %s SET %s WHERE id = %d;",
			model.TableName(),
			sqlValues,
			id,
		),
	)
	return err
}

func (o *Orm) Delete(model Model, id int64) error {
	sqlReq := fmt.Sprintf("DELETE FROM %s WHERE id = %d",
		model.TableName(),
		id,
	)
	_, err := o.db.Exec(sqlReq)
	return err
}

func (o *Orm) Exists(model Model, filter []Filter) bool {
	filters := ""
	if len(filter) != 0 {
		filters = filter[0].ToSQL(filter...)
	}
	err := o.db.QueryRow(
		fmt.Sprintf("SELECT %s FROM %s %s;",
			getModelFields(model, true),
			model.TableName(),
			filters,
		),
	).Scan(Fields(model)...)
	return err != sql.ErrNoRows
}
