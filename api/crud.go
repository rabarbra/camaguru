package main

import (
	"database/sql"
	"fmt"
)

func createUser(db *sql.DB, user User) error {
	_, err := db.Exec(`
		INSERT INTO users(username, email, pass, email_verified, like_notify, comm_notify)
		VALUES(?, ?, ?, ?, ?, ?)`,
		user.Username, user.Email, user.Pass, user.EmailVerified, user.LikeNotify, user.CommNotify,
	)
	return err
}

func updateUser(db *sql.DB, user User) error {
	_, err := db.Exec("UPDATE users SET username=?, email=?, pass=?, email_verified=?, like_notify=?, comm_notify=? WHERE id=?",
		user.Username, user.Email, user.Pass, user.EmailVerified, user.LikeNotify, user.CommNotify, user.Id,
	)
	return err
}

func getUser(db *sql.DB, id int64) (User, error) {
	var user User
	err := db.QueryRow(
		"SELECT id, username, email, pass FROM users WHERE id=?",
		id,
	).Scan(&user.Id, &user.Username, &user.Email, &user.Pass)
	return user, err
}

func checkUserExists(db *sql.DB, field string, value string) bool {
	var u string
	query := fmt.Sprintf("SELECT %s FROM users WHERE %s=?", field, field)
	err := db.QueryRow(query, value).Scan(&u)
	if err == nil {
		return true
	} else if err == sql.ErrNoRows {
		return false
	}
	return true
}
