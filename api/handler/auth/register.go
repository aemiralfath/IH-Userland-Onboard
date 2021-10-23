package auth

import (
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"unicode"

	"github.com/aemiralfath/IH-Userland-Onboard/api/handler"
	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
	"github.com/go-chi/render"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type registerRequest struct {
	Fullname        string `json:"fullname"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

func Register(userStore datastore.UserStore, profileStore datastore.ProfileStore, passwordStore datastore.PasswordStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		req := &registerRequest{}

		if err := render.Bind(r, req); err != nil {
			fmt.Println(render.Render(w, r, handler.BadRequestErrorRenderer(err)))
			return
		}

		hashPassword, err := hash(req.Password)
		if err != nil {
			fmt.Println(render.Render(w, r, handler.InternalServerErrorRenderer(err)))
			return
		}

		req.Password = string(hashPassword)
		if err := userStore.AddNewUser(ctx, parseRegisterUser(req), parseRegisterProfile(req), parseRegisterPassword(req)); err != nil {
			fmt.Println(render.Render(w, r, handler.InternalServerErrorRenderer(err)))
			return
		}

		go sendEmailVerification(req.Email)

		if err := render.Render(w, r, handler.SuccesRenderer()); err != nil {
			fmt.Println(render.Render(w, r, handler.InternalServerErrorRenderer(err)))
			return
		}
	}
}

func sendEmailVerification(toEmail string) {
	from := "criptdestroyer@gmail.com"
	msg := []byte("To: " + toEmail + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: Userland Email Verification!\r\n" +
		"\r\n" +
		"This is the email is sent using golang and sendinblue.\r\n")

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

func parseRegisterUser(u *registerRequest) *datastore.User {
	return &datastore.User{
		Email:    u.Email,
		Password: u.Password,
	}
}

func parseRegisterProfile(u *registerRequest) *datastore.Profile {
	return &datastore.Profile{
		Fullname: u.Fullname,
	}
}

func parseRegisterPassword(u *registerRequest) *datastore.Password {
	return &datastore.Password{
		Password: u.Password,
	}
}

func hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func verifyPassword(s string) (eightOrMore, number, upper bool) {
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

func (register *registerRequest) Bind(r *http.Request) error {
	if strings.TrimSpace(register.Fullname) == "" {
		return fmt.Errorf("required fullname")
	}

	if strings.TrimSpace(register.Email) == "" {
		return fmt.Errorf("required email")
	}

	if register.Password != register.PasswordConfirm {
		return fmt.Errorf("password and confirm password must same!")
	}

	passLength, number, upper := verifyPassword(register.Password)
	if !passLength || !number || !upper {
		return fmt.Errorf("password must have lowercase, uppercase, number, and minimum 8 chars!")
	}

	return nil
}

func (*registerRequest) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
