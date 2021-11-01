package crypto

import (
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
