package main

import (
	"bytes"
	"fmt"
	"html/template"
	"jwt"
	"log"
	"net/smtp"
	"orm"
	"strings"
	"time"
)

func SendEmail(email string, body string, subject string) error {
	headers := make(map[string]string)
	headers["From"] = SMTP_FROM
	headers["To"] = email
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	var message strings.Builder
	for key, value := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}
	message.WriteString("\r\n")
	message.WriteString(body)

	auth := smtp.PlainAuth("", SMTP_FROM, SMTP_PASSWORD, SMTP_HOST)
	err := smtp.SendMail(SMTP_HOST+":"+SMTP_PORT, auth, SMTP_FROM, []string{email}, []byte(message.String()))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

type EmailData struct {
	Username string
	Link     string
}

func sendVerificationEmail(email string, db *orm.Orm) error {
	verifyLink := BACKEND_HOST + "/auth/verify?token="
	var user User
	err := db.GetOne(&user, []orm.Filter{{Key: "email", Value: email, Operation: "="}})
	if err != nil {
		log.Println("sendVerificationEmail: " + err.Error())
		return err
	}
	token, err := jwt.CreateJWT(JWT_SECRET, user.BaseModel.Id, time.Hour*100)
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
	return SendEmail(email, body.String(), "Camaguru Verify Email")
}

func sendResetPasswordEmail(email string, db *orm.Orm) error {
	verifyLink := FRONTEND_HOST + "/reset?token="
	var user User
	err := db.GetOne(&user, []orm.Filter{{Key: "email", Value: email, Operation: "="}})
	if err != nil {
		return err
	}
	token, err := jwt.CreateJWT(JWT_SECRET, user.BaseModel.Id, time.Hour*3)
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
	return SendEmail(email, body.String(), "Camaguru Reset Password")
}
