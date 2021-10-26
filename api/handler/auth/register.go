package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/aemiralfath/IH-Userland-Onboard/api/helper"
	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
	"github.com/go-chi/render"
)

type registerRequest struct {
	Fullname        string `json:"fullname"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

func Register(userStore datastore.UserStore, profileStore datastore.ProfileStore, passwordStore datastore.PasswordStore, token datastore.TokenStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		req := &registerRequest{}

		if err := render.Bind(r, req); err != nil {
			render.Render(w, r, helper.BadRequestErrorRenderer(err))
			return
		}

		hashPassword, err := helper.HashPassword(req.Password)
		if err != nil {
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}

		usr, _ := userStore.GetUserByEmail(ctx, req.Email)
		if usr != nil {
			render.Render(w, r, helper.BadRequestErrorRenderer(fmt.Errorf("Email already exist!")))
			return
		}

		req.Password = string(hashPassword)
		userId, err := userStore.AddNewUser(ctx, parseRegisterRequestUser(req))
		if err != nil {
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}

		if err := profileStore.AddNewProfile(ctx, parseRegisterRequestProfile(req), userId); err != nil {
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}

		if err := passwordStore.AddNewPassword(ctx, parseRegisterRequestPassword(req), userId); err != nil {
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}

		tokenCode, err := helper.GenerateOTP(6)
		if err != nil {
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}

		if err := token.SetToken(ctx, "user", req.Email, tokenCode); err != nil {
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}

		subject := "Userland Email Verification!"
		msg := fmt.Sprintf("Use this otp for verify your email: %s", tokenCode)

		go helper.SendEmail(req.Email, subject, msg)

		if err := render.Render(w, r, helper.SuccesRenderer()); err != nil {
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}
	}
}

func parseRegisterRequestUser(u *registerRequest) *datastore.User {
	return &datastore.User{
		Email:    u.Email,
		Password: u.Password,
	}
}

func parseRegisterRequestProfile(u *registerRequest) *datastore.Profile {
	return &datastore.Profile{
		Fullname: u.Fullname,
	}
}

func parseRegisterRequestPassword(u *registerRequest) *datastore.Password {
	return &datastore.Password{
		Password: u.Password,
	}
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

	passLength, number, upper := helper.VerifyPassword(register.Password)
	if !passLength || !number || !upper {
		return fmt.Errorf("password must have lowercase, uppercase, number, and minimum 8 chars!")
	}

	return nil
}

func (*registerRequest) Render(w http.ResponseWriter, r *http.Request) {}
