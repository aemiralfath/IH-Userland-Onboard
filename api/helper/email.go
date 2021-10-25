package helper

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

func SendEmailVerification(toEmail string, otp string) {
	from := "criptdestroyer@gmail.com"
	msg := []byte("To: " + toEmail + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: Userland Email Verification!\r\n" +
		"\r\n" +
		"Use this otp for verify your email:" + otp + "\r\n")

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(fmt.Errorf("Error loading .env file"))
	}

	auth := smtp.PlainAuth("", from, os.Getenv("PASSWORD"), os.Getenv("SMTP_HOST"))
	smtpAddress := fmt.Sprintf("%s:%v", os.Getenv("SMTP_HOST"), os.Getenv("SMTP_PORT"))
	err = smtp.SendMail(smtpAddress, auth, from, []string{toEmail}, msg)
	if err != nil {
		fmt.Println(fmt.Errorf(err.Error()))
	}
}
