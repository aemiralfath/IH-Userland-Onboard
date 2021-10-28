package helper

import (
	"fmt"
	"net/smtp"
)

type EmailConfig struct {
	Host     string
	Port     string
	From     string
	Password string
}

type Email struct {
	Auth smtp.Auth
	Addr string
	From string
}

func NewEmail(config EmailConfig) *Email {
	return &Email{
		Auth: smtp.PlainAuth("", config.From, config.Password, config.Host),
		Addr: fmt.Sprintf("%s:%v", config.Host, config.Port),
		From: config.From,
	}
}

func (email *Email) SendEmail(toEmail string, subject string, msg string) {
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
