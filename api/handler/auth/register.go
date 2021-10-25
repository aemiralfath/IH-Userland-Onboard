package auth

import (
	"fmt"
	"net/http"
	"strings"
	"unicode"

	"github.com/aemiralfath/IH-Userland-Onboard/api/helper"
	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
)

type registerRequest struct {
	Fullname        string `json:"fullname"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

func Register(userStore datastore.UserStore, profileStore datastore.ProfileStore, passwordStore datastore.PasswordStore, otp datastore.OTPStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		req := &registerRequest{}

		if err := render.Bind(r, req); err != nil {
			render.Render(w, r, helper.BadRequestErrorRenderer(err))
			return
		}

		hashPassword, err := hash(req.Password)
		if err != nil {
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}

		req.Password = string(hashPassword)
		if err := userStore.AddNewUser(ctx, parseRegisterUser(req), parseRegisterProfile(req), parseRegisterPassword(req)); err != nil {
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}

		token, err := helper.GenerateOTP(6)
		if err != nil {
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}

		res, err := otp.GetOTP(ctx, req.Email, token)
		if err != nil {
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}

		subject := "Userland Email Verification!"
		msg := fmt.Sprintf("Use this otp for verify your email: %s", res)

		go helper.SendEmail(req.Email, subject, msg)

		if err := render.Render(w, r, helper.SuccesRenderer()); err != nil {
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}
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

func (*registerRequest) Render(w http.ResponseWriter, r *http.Request) {}
