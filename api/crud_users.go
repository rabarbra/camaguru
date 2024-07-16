package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

func createUser(db *sql.DB, user User) error {
	res, err := db.Exec(fmt.Sprintf(`
		INSERT INTO users(username, email, pass, email_verified, like_notify, comm_notify)
		VALUES('%s', '%s', '%s', %t, %t, %t)`,
		user.Username, user.Email, user.Pass, user.EmailVerified, user.LikeNotify, user.CommNotify,
	))
	fmt.Println(res, err)
	return err
}

func updateUser(db *sql.DB, user User) error {
	_, err := db.Exec(fmt.Sprintf("UPDATE users SET username='%s', email='%s', pass='%s', avatar='%s', email_verified=%t, like_notify=%t, comm_notify=%t WHERE id=%d",
		user.Username, user.Email, user.Pass, user.Avatar, user.EmailVerified, user.LikeNotify, user.CommNotify, user.Id,
	))
	return err
}

func partUpdateUser(db *sql.DB, userId int64, updates map[string]interface{}) error {
	var setClauses []string
	var args []interface{}
	argId := 1
	for field, value := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", field, argId))
		args = append(args, value)
		argId++
	}
	args = append(args, userId)

	query := fmt.Sprintf("UPDATE users SET %s WHERE id = $%d", strings.Join(setClauses, ", "), argId)
	_, err := db.Exec(query, args...)
	if err != nil {
		log.Println("Update user error:", err)
		return err
	}
	return nil
}

func getUserById(db *sql.DB, id int64) (User, error) {
	var user User
	err := db.QueryRow(fmt.Sprintf(
		"SELECT id, username, email, pass, avatar, email_verified, like_notify, comm_notify FROM users WHERE id=%d;",
		id,
	)).Scan(&user.Id, &user.Username, &user.Email, &user.Pass, &user.Avatar, &user.EmailVerified, &user.LikeNotify, &user.CommNotify)
	return user, err
}

func getUserByUsename(db *sql.DB, username string) (User, error) {
	var user User
	err := db.QueryRow(fmt.Sprintf(
		"SELECT id, username, email, pass, avatar, email_verified, like_notify, comm_notify FROM users WHERE username='%s'",
		username,
	)).Scan(&user.Id, &user.Username, &user.Email, &user.Pass, &user.Avatar, &user.EmailVerified, &user.LikeNotify, &user.CommNotify)
	return user, err
}

func getUserByEmail(db *sql.DB, email string) (User, error) {
	var user User
	err := db.QueryRow(fmt.Sprintf(
		"SELECT id, username, email, pass, avatar, email_verified, like_notify, comm_notify FROM users WHERE email='%s'",
		email,
	)).Scan(&user.Id, &user.Username, &user.Email, &user.Pass, &user.Avatar, &user.EmailVerified, &user.LikeNotify, &user.CommNotify)
	return user, err
}

func checkUserExists(db *sql.DB, field string, value string, id int64) bool {
	var u string
	query := fmt.Sprintf("SELECT %s FROM users WHERE %s='%s';", field, field, value)
	if id != 0 {
		query = fmt.Sprintf("SELECT %s FROM users WHERE %s='%s' AND id != %d;", field, field, value, id)
	}
	err := db.QueryRow(query).Scan(&u)
	if err == nil {
		return true
	} else if err == sql.ErrNoRows {
		return false
	}
	return true
}
