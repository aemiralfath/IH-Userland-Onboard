package helper

import (
	"crypto/rand"
	"unicode"

	"github.com/thanhpk/randstr"
)

func VerifyPassword(s string) (eightOrMore, number, upper bool) {
	letters := 0
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
			letters++
		case unicode.IsUpper(c):
			upper = true
			letters++
		case unicode.IsLetter(c) || c == ' ':
			letters++
		}
	}
	eightOrMore = letters >= 8
	return
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

func GenerateRandomID() string {
	randomIDBytes := 128 / 8 // 128-bit
	return randstr.Hex(randomIDBytes)
}
