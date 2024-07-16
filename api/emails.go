package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
	"jwt"
	"log"
	"net/smtp"
	"strings"
	"time"
)

type EmailData struct {
	Username string
	Link     string
}

func sendVerificationEmail(email string, db *sql.DB) error {
	verifyLink := BACKEND_HOST + "/auth/verify?token="
	user, err := getUserByEmail(db, email)
	if err != nil {
		log.Println("sendVerificationEmail: " + err.Error())
		return err
	}
	token, err := jwt.CreateJWT(JWT_SECRET, user.Id, time.Hour*100)
	if err != nil {
		return err
	}
	verifyLink += token

	data := EmailData{
		Username: user.Username,
		Link:     verifyLink,
	}
	t, err := template.ParseFiles("./assets/templates/email_verify.html")
	if err != nil {
		log.Println(err)
		return err
	}

	var body bytes.Buffer
	err = t.Execute(&body, data)
	if err != nil {
		log.Println(err)
		return err
	}

	headers := make(map[string]string)
	headers["From"] = SMTP_FROM
	headers["To"] = email
	headers["Subject"] = "Camaguru Verify Email"
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	var message strings.Builder
	for key, value := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}
	message.WriteString("\r\n")
	message.WriteString(body.String())

	auth := smtp.PlainAuth("", SMTP_FROM, SMTP_PASSWORD, SMTP_HOST)
	err = smtp.SendMail(SMTP_HOST+":"+SMTP_PORT, auth, SMTP_FROM, []string{email}, []byte(message.String()))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func sendResetPasswordEmail(email string, db *sql.DB) error {
	verifyLink := FRONTEND_HOST + "/reset?token="
	user, err := getUserByEmail(db, email)
	if err != nil {
		return err
	}
	token, err := jwt.CreateJWT(JWT_SECRET, user.Id, time.Hour*3)
	if err != nil {
		return err
	}
	verifyLink += token

	data := EmailData{
		Username: user.Username,
		Link:     verifyLink,
	}
	t, err := template.ParseFiles("./assets/templates/email_reset.html")
	if err != nil {
		log.Println(err)
		return err
	}

	var body bytes.Buffer
	err = t.Execute(&body, data)
	if err != nil {
		log.Println(err)
		return err
	}

	headers := make(map[string]string)
	headers["From"] = SMTP_FROM
	headers["To"] = email
	headers["Subject"] = "Camaguru Reset Password"
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	var message strings.Builder
	for key, value := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}
	message.WriteString("\r\n")
	message.WriteString(body.String())

	auth := smtp.PlainAuth("", SMTP_FROM, SMTP_PASSWORD, SMTP_HOST)
	err = smtp.SendMail(SMTP_HOST+":"+SMTP_PORT, auth, SMTP_FROM, []string{email}, []byte(message.String()))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
