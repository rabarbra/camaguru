package orm

import (
	"database/sql"
	"io"
	"log"
	"os"
)

type BaseModel struct {
	Id int64
}

type IOrm interface {
	Connect(connString string) error
	Migrate(migrationPath string) error
	Close()
	Exec(sql string) (sql.Result, error)
	// Operations
	GetOne(model BaseModel, filter []Filter) (BaseModel, error)
	GetOneById(model BaseModel, id int64) (BaseModel, error)
	GetMany(model BaseModel, filter []Filter, sort []Sort, pagination Pagination) ([]BaseModel, error)
	Create(model BaseModel) (BaseModel, error)
	Update(model BaseModel, id int64) (BaseModel, error)
	Delete(model BaseModel, id int64) error
}

type Orm struct {
	db *sql.DB
}

func (o Orm) Connect(connString string) error {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Println(err)
		return err
	}
	o.db = db
	return nil
}

func (o Orm) Close() {
	o.db.Close()
}

func (o Orm) Migrate(migrationPath string) error {
	sqlFile, err := os.Open(migrationPath)
	if err != nil {
		log.Println(err)
		return err
	}
	defer sqlFile.Close()

	sqlBytes, err := io.ReadAll(sqlFile)
	if err != nil {
		log.Println(err)
		return err
	}

	sqlCommands := string(sqlBytes)
	_, err = o.db.Exec(sqlCommands)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (o Orm) Exec(sql string) (sql.Result, error) {
	return o.db.Exec(sql)
}
