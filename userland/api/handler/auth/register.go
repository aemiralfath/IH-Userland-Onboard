package auth

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/api/crypto"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/api/email"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/api/helper"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/datastore"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

type RegisterRequest struct {
	Fullname        string `json:"fullname"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

func Register(email email.Email, crypto crypto.Crypto, userStore datastore.UserStore, profileStore datastore.ProfileStore, passwordStore datastore.PasswordStore, otp datastore.OTPStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		req := &RegisterRequest{}

		if err := render.Bind(r, req); err != nil {
			render.Render(w, r, helper.BadRequestErrorRenderer(err))
			return
		}

		_, err := userStore.CheckUserEmailExist(ctx, req.Email)
		if err != nil && err != sql.ErrNoRows {
			log.Error().Err(err).Stack().Msg(err.Error())
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		} else if err == nil {
			fmt.Println(render.Render(w, r, helper.BadRequestErrorRenderer(fmt.Errorf("Email already used"))))
			return
		}

		hashPassword, err := crypto.HashPassword(req.Password)
		if err != nil {
			log.Error().Err(err).Stack().Msg(err.Error())
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}

		req.Password = hashPassword
		userId, err := userStore.AddNewUser(ctx, parseRegisterRequestUser(req))
		if err != nil {
			log.Error().Err(err).Stack().Msg(err.Error())
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}

		if err := profileStore.AddNewProfile(ctx, parseRegisterRequestProfile(req), userId); err != nil {
			log.Error().Err(err).Stack().Msg(err.Error())
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}

		if err := passwordStore.AddNewPassword(ctx, parseRegisterRequestPassword(req), userId); err != nil {
			log.Error().Err(err).Stack().Msg(err.Error())
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}

		otpCode, err := crypto.GenerateOTP(6)
		if err != nil {
			log.Error().Err(err).Stack().Msg(err.Error())
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}

		otpValue := fmt.Sprintf("%f-%s", userId, req.Email)
		if err := otp.SetOTP(ctx, "user", otpCode, otpValue); err != nil {
			log.Error().Err(err).Stack().Msg(err.Error())
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}

		subject := "Userland Email Verification!"
		msg := fmt.Sprintf("Use this otp for verify your email: %s", otpCode)

		go email.SendEmail(req.Email, subject, msg)

		if err := render.Render(w, r, helper.SuccesRenderer()); err != nil {
			log.Error().Err(err).Stack().Msg(err.Error())
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}
	}
}

func parseRegisterRequestUser(u *RegisterRequest) *datastore.User {
	return &datastore.User{
		Email:    u.Email,
		Password: u.Password,
	}
}

func parseRegisterRequestProfile(u *RegisterRequest) *datastore.Profile {
	return &datastore.Profile{
		Fullname: u.Fullname,
	}
}

func parseRegisterRequestPassword(u *RegisterRequest) *datastore.Password {
	return &datastore.Password{
		Password: u.Password,
	}
}

func (register *RegisterRequest) Bind(r *http.Request) error {
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

func (*RegisterRequest) Render(w http.ResponseWriter, r *http.Request) {}