package email

import (
	"fmt"
	"net/smtp"
	"os"
)

type Email interface {
	SendEmail(toEmail string, subject string, msg string)
}

type EmailSendInBlue struct {
	Auth smtp.Auth
	Addr string
	From string
}

func NewEmail() Email {
	return &EmailSendInBlue{
		Auth: smtp.PlainAuth("", os.Getenv("SMTP_FROM"), os.Getenv("SMTP_PASSWORD"), os.Getenv("SMTP_HOST")),
		Addr: fmt.Sprintf("%s:%v", os.Getenv("SMTP_HOST"), os.Getenv("SMTP_PORT")),
		From: os.Getenv("SMTP_FROM"),
	}
}

func (email *EmailSendInBlue) SendEmail(toEmail string, subject string, msg string) {
	from := email.From
	value := []byte("To: " + toEmail + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + msg + "\r\n")

	err := smtp.SendMail(email.Addr, email.Auth, from, []string{toEmail}, value)
	if err != nil {
		fmt.Println(fmt.Errorf(err.Error()))
	}
}
