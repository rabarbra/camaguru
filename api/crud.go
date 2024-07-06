package main

import (
	"database/sql"
	"fmt"
)

func createUser(db *sql.DB, user User) error {
	res, err := db.Exec(`
		INSERT INTO users(username, email, pass, email_verified, like_notify, comm_notify)
		VALUES(?, ?, ?, ?, ?, ?)`,
		user.Username, user.Email, user.Pass, user.EmailVerified, user.LikeNotify, user.CommNotify,
	)
	fmt.Println(res, err)
	return err
}

func updateUser(db *sql.DB, user User) error {
	_, err := db.Exec("UPDATE users SET username=?, email=?, pass=?, email_verified=?, like_notify=?, comm_notify=? WHERE id=?",
		user.Username, user.Email, user.Pass, user.EmailVerified, user.LikeNotify, user.CommNotify, user.Id,
	)
	return err
}

func getUserById(db *sql.DB, id int64) (User, error) {
	var user User
	err := db.QueryRow(
		"SELECT id, username, email, pass FROM users WHERE id=?",
		id,
	).Scan(&user.Id, &user.Username, &user.Email, &user.Pass)
	return user, err
}

func getUserByUsename(db *sql.DB, username string) (User, error) {
	var user User
	err := db.QueryRow(
		"SELECT id, username, email, pass FROM users WHERE username=?",
		username,
	).Scan(&user.Id, &user.Username, &user.Email, &user.Pass)
	return user, err
}

func checkUserExists(db *sql.DB, field string, value string, id int64) bool {
	var u string
	query := fmt.Sprintf("SELECT %s FROM users WHERE %s=?", field, field)
	if id != 0 {
		query = fmt.Sprintf("SELECT %s FROM users WHERE %s=? AND id != %d", field, field, id)
	}
	err := db.QueryRow(query, value).Scan(&u)
	if err == nil {
		return true
	} else if err == sql.ErrNoRows {
		return false
	}
	return true
}
