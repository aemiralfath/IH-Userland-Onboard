package crypto

import (
	"crypto/rand"

	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
	"golang.org/x/crypto/bcrypt"
)

type AppCrypto struct {
}

func NewAppCrypto() datastore.Crypto {
	return &AppCrypto{}
}

func (c *AppCrypto) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (c *AppCrypto) ConfirmPassword(hashedPassword, password string) bool {
	byteHash := []byte(hashedPassword)
	err := bcrypt.CompareHashAndPassword(byteHash, []byte(password))

	return err == nil
}

func (c *AppCrypto) GenerateOTP(length int) (string, error) {
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
