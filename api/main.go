package main

import (
	"fmt"
	"log"
	"net/http"
	"orm"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_DB"),
	)

	var db orm.Orm
	db.Connect(connString)
	defer db.Close()
	db.Migrate("./assets/migrations/01_create_tables.up.sql")

	http.HandleFunc("/me/avatar", CorsM(AuthDBM(&db, postUserAvatar)))
	http.HandleFunc("/me", CorsM(AuthDBM(&db, getUser)))
	http.HandleFunc("POST /me", CorsM(DBM(&db, postUser)))
	http.HandleFunc("PUT /me", CorsM(AuthDBM(&db, putUser)))

	http.HandleFunc("/auth/signin", CorsM(DBM(&db, signin)))
	http.HandleFunc("/auth/verify", CorsM(DBM(&db, verifyEmail)))
	http.HandleFunc("/auth/pass", CorsM(DBM(&db, resetPasswordUnauth)))
	http.HandleFunc("/auth/reset", CorsM(AuthDBM(&db, resetPassword)))

	AddCrudRoutes(&Img{}, &db)
	// AddCrudRoutes(&Like{}, &db)
	// AddCrudRoutes(&Comment{}, db)

	fs := http.FileServer(http.Dir("./assets/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Fatal(http.ListenAndServe(":8000", nil))
}
