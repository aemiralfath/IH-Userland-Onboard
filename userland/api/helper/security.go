package helper

import (
	"net/http"
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

func GenerateRandomID() string {
	randomIDBytes := 128 / 8 // 128-bit
	return randstr.Hex(randomIDBytes)
}

func GetIPAddress(r *http.Request) string {
	ip := r.Header.Get("X-Real-Ip")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}

	if ip == "" {
		ip = r.RemoteAddr
	}
	return ip
}
