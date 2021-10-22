package auth

import (
	"fmt"
	"net/http"
	"strings"
	"unicode"

	"github.com/aemiralfath/IH-Userland-Onboard/api/handler"
	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
	"github.com/aemiralfath/IH-Userland-Onboard/datastore/models"
	"github.com/go-chi/render"
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
			render.Render(w, r, handler.BadRequestErrorRenderer(err))
			return
		}

		hashPassword, err := hash(req.Password)
		if err != nil {
			render.Render(w, r, handler.InternalServerErrorRenderer(err))
			return
		}

		req.Password = string(hashPassword)
		if err := userStore.AddNewUser(ctx, parseHandlerUser(req), parseHandlerProfile(req), parseHandlerPassword(req)); err != nil {
			render.Render(w, r, handler.InternalServerErrorRenderer(err))
			return
		}

		if err := render.Render(w, r, handler.SuccesRenderer("Success")); err != nil {
			render.Render(w, r, handler.InternalServerErrorRenderer(err))
			return
		}
	}
}

func parseHandlerUser(u *registerRequest) *models.User {
	return &models.User{
		Email:    u.Email,
		Password: u.Password,
	}
}

func parseHandlerProfile(u *registerRequest) *models.Profile {
	return &models.Profile{
		Fullname: u.Fullname,
	}
}

func parseHandlerPassword(u *registerRequest) *models.Password {
	return &models.Password{
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
