package auth

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/api/crypto"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/api/helper"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/datastore"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

type ResetPasswordRequest struct {
	Token           string `json:"token"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

func ResetPassword(crypto crypto.Crypto, userStore datastore.UserStore, passwordStore datastore.PasswordStore, otp datastore.OTPStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		req := &ResetPasswordRequest{}

		if err := render.Bind(r, req); err != nil {
			render.Render(w, r, helper.BadRequestErrorRenderer(err))
			return
		}

		email, err := otp.GetOTP(ctx, "password", req.Token)
		if err != nil {
			render.Render(w, r, helper.BadRequestErrorRenderer(err))
			return
		}

		usr, err := userStore.GetUserByEmail(ctx, email)
		if err != nil {
			if err == sql.ErrNoRows {
				render.Render(w, r, helper.BadRequestErrorRenderer(fmt.Errorf("User not found")))
				return
			} else {
				log.Error().Err(err).Stack().Msg(err.Error())
				render.Render(w, r, helper.InternalServerErrorRenderer(err))
				return
			}
		}

		lastThreePassword, err := passwordStore.GetLastThreePassword(ctx, usr.ID)
		if err != nil {
			log.Error().Err(err).Stack().Msg(err.Error())
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}

		for _, e := range lastThreePassword {
			if err := crypto.ConfirmPassword(e, req.Password); err {
				render.Render(w, r, helper.BadRequestErrorRenderer(fmt.Errorf("Password must different from last 3 password")))
				return
			}
		}

		if err := updateStore(ctx, crypto, req, usr, userStore, passwordStore); err != nil {
			log.Error().Err(err).Stack().Msg(err.Error())
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}

		if err := render.Render(w, r, helper.SuccesRenderer()); err != nil {
			log.Error().Err(err).Stack().Msg(err.Error())
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}

	}
}

func updateStore(
	ctx context.Context,
	crypto crypto.Crypto,
	req *ResetPasswordRequest,
	usr *datastore.User,
	userStore datastore.UserStore,
	passwordStore datastore.PasswordStore) error {

	hashPassword, err := crypto.HashPassword(req.Password)
	if err != nil {
		return err
	}

	usr.Password = string(hashPassword)
	if err := userStore.ChangePassword(ctx, usr); err != nil {
		return err
	}

	req.Password = hashPassword
	if err := passwordStore.AddNewPassword(ctx, parseResetRequestPassword(req), usr.ID); err != nil {
		return err
	}
	return nil
}

func parseResetRequestPassword(u *ResetPasswordRequest) *datastore.Password {
	return &datastore.Password{
		Password: u.Password,
	}
}

func (request *ResetPasswordRequest) Bind(r *http.Request) error {
	if strings.TrimSpace(request.Token) == "" {
		return fmt.Errorf("required token")
	}

	if request.Password != request.PasswordConfirm {
		return fmt.Errorf("password and confirm password must same!")
	}

	passLength, number, upper := helper.VerifyPassword(request.Password)
	if !passLength || !number || !upper {
		return fmt.Errorf("password must have lowercase, uppercase, number, and minimum 8 chars!")
	}

	return nil
}

func (*ResetPasswordRequest) Render(w http.ResponseWriter, r *http.Request) {}