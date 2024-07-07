package main

import (
	"jwt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func HashPass(pass string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), 14)
	return string(bytes), err
}

func CheckPassHash(pass string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	return err == nil
}

func CheckAuthorized(r *http.Request) int64 {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return 0
	}
	tokenString = tokenString[len("Bearer "):]

	p, err := jwt.VerifyJWT(tokenString, JWT_SECRET)
	if err != nil {
		return 0
	}
	return p.Id
}
