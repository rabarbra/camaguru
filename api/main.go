package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var JWT_SECRET = os.Getenv("JWT_SECRET")

func migrate(db *sql.DB) {
	sqlFile, err := os.Open("./assets/migrations/01_create_tables.up.sql")
	if err != nil {
		log.Fatal(err)
	}
	defer sqlFile.Close()

	sqlBytes, err := ioutil.ReadAll(sqlFile)
	if err != nil {
		log.Fatal(err)
	}

	sqlCommands := string(sqlBytes)
	_, err = db.Exec(sqlCommands)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_DB"),
	)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	migrate(db)

	http.HandleFunc("/signin", CorsM(DBM(db, signin)))
	http.HandleFunc("/me", CorsM(AuthDBM(db, getUser)))
	http.HandleFunc("POST /me", CorsM(DBM(db, postUser)))
	http.HandleFunc("PUT /me", CorsM(AuthDBM(db, putUser)))
	http.HandleFunc("/img", CorsM(DBM(db, img)))

	fs := http.FileServer(http.Dir("./assets/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Fatal(http.ListenAndServe(":8000", nil))
}
