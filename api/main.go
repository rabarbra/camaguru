package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type Handler func(http.ResponseWriter, *http.Request, *sql.DB)

func DBWrapper(db *sql.DB, f Handler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		f(w, r, db)
	}
}

func main() {
	db, err := sql.Open("sqlite3", "./camaguru.db")
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	http.HandleFunc("/signin", DBWrapper(db, signin))
	http.HandleFunc("/me", DBWrapper(db, me))
	http.HandleFunc("/img", DBWrapper(db, img))

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
