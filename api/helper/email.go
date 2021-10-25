package helper

import (
	"crypto/rand"
	"fmt"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

func SendEmail(toEmail string, subject string, msg string) {
	from := "criptdestroyer@gmail.com"
	err := godotenv.Load(".env")

	email := []byte("To: " + toEmail + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + msg + "\r\n")

	if err != nil {
		fmt.Println(fmt.Errorf("Error loading .env file"))
		fmt.Println(fmt.Errorf("Error loading .env file"))
	}

	auth := smtp.PlainAuth("", from, os.Getenv("PASSWORD"), os.Getenv("SMTP_HOST"))
	smtpAddress := fmt.Sprintf("%s:%v", os.Getenv("SMTP_HOST"), os.Getenv("SMTP_PORT"))
	err = smtp.SendMail(smtpAddress, auth, from, []string{toEmail}, email)
	if err != nil {
		fmt.Println(fmt.Errorf(err.Error()))
	}
}

func GenerateOTP(length int) (string, error) {
	otpChars := "1234567890"
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	otpCharsLength := len(otpChars)
	for i := 0; i < length; i++ {
		buffer[i] = otpChars[int(buffer[i])%otpCharsLength]
	}

	return string(buffer), nil
}
