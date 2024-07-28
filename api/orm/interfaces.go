package orm

import (
	"database/sql"
	"io"
	"log"
	"net/http"
	"os"
)

type Model interface {
	TableName() string
	PrimaryKey() string
	NewItem(r *http.Request, userId int64) error
}

type BaseModel struct {
	Id int64 `json:"id"`
}

func (b BaseModel) PrimaryKey() string {
	return "id"
}

type IOrm interface {
	Connect(connString string) error
	Migrate(migrationPath string) error
	Close()
	Exec(query string) (sql.Result, error)
	// Operations
	GetOne(model Model, filter []Filter) error
	GetOneById(model Model, id int64) error
	GetMany(
		model Model,
		filter []Filter,
		sort []Sort,
		pagination Pagination,
	) ([]Model, error)
	Create(model Model) (int64, error)
	Update(model Model, id int64) error
	Delete(model Model, id int64) error
	//
	Exists(model Model, filter []Filter) bool
}

type Orm struct {
	db *sql.DB
}

func (o *Orm) Connect(connString string) error {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Println("Error in Connect: ", err)
		return err
	}
	o.db = db
	return nil
}

func (o *Orm) Close() {
	if o.db != nil {
		o.db.Close()
	}
}

func (o *Orm) Migrate(migrationPath string) error {
	sqlFile, err := os.Open(migrationPath)
	if err != nil {
		log.Println("Error in Migrate, opening sql file: ", err)
		return err
	}
	defer sqlFile.Close()

	sqlBytes, err := io.ReadAll(sqlFile)
	if err != nil {
		log.Println("Error in Migrate, reading sql file: ", err)
		return err
	}

	sqlCommands := string(sqlBytes)
	_, err = o.db.Exec(sqlCommands)
	if err != nil {
		log.Println("Error in Migrate, executing migrations: ", err)
		return err
	}
	return nil
}

func (o *Orm) Exec(query string) (sql.Result, error) {
	return o.db.Exec(query)
}
