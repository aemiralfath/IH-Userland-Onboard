package me

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/api/crypto"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/api/helper"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/api/jwt"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/datastore"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

type ChangePasswordRequest struct {
	PasswordCurrent string `json:"password_current"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

func ChangePassword(jwtAuth jwt.JWT, crypto crypto.Crypto, userStore datastore.UserStore, passwordStore datastore.PasswordStore) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		req := &ChangePasswordRequest{}

		if err := render.Bind(r, req); err != nil {
			render.Render(rw, r, helper.BadRequestErrorRenderer(err))
			return
		}

		_, claims, err := jwtAuth.FromContext(ctx)
		if err != nil {
			fmt.Println(render.Render(rw, r, helper.BadRequestErrorRenderer(err)))
			return
		}

		emailUser := claims["email"]
		usr, err := userStore.GetUserByEmail(ctx, emailUser.(string))
		if err != nil {
			if err == sql.ErrNoRows {
				render.Render(rw, r, helper.BadRequestErrorRenderer(fmt.Errorf("User not found")))
				return
			} else {
				log.Error().Err(err).Stack().Msg(err.Error())
				render.Render(rw, r, helper.InternalServerErrorRenderer(err))
				return
			}
		}

		if err := crypto.ConfirmPassword(usr.Password, req.PasswordCurrent); !err {
			render.Render(rw, r, helper.BadRequestErrorRenderer(fmt.Errorf("Wrong Password")))
			return
		}

		lastThreePassword, err := passwordStore.GetLastThreePassword(ctx, usr.ID)
		if err != nil {
			log.Error().Err(err).Stack().Msg(err.Error())
			render.Render(rw, r, helper.InternalServerErrorRenderer(err))
			return
		}

		for _, e := range lastThreePassword {
			if err := crypto.ConfirmPassword(e, req.Password); err {
				render.Render(rw, r, helper.BadRequestErrorRenderer(fmt.Errorf("Password must different from last 3 password")))
				return
			}
		}

		if err := updateStore(ctx, crypto, req, usr, userStore, passwordStore); err != nil {
			log.Error().Err(err).Stack().Msg(err.Error())
			render.Render(rw, r, helper.InternalServerErrorRenderer(err))
			return
		}

		if err := render.Render(rw, r, helper.SuccesRenderer()); err != nil {
			log.Error().Err(err).Stack().Msg(err.Error())
			render.Render(rw, r, helper.InternalServerErrorRenderer(err))
			return
		}

	}
}

func updateStore(ctx context.Context, crypto crypto.Crypto, req *ChangePasswordRequest, usr *datastore.User, userStore datastore.UserStore, passwordStore datastore.PasswordStore) error {
	hashPassword, err := crypto.HashPassword(req.Password)
	if err != nil {
		return err
	}

	usr.Password = string(hashPassword)
	if err := userStore.ChangePassword(ctx, usr); err != nil {
		return err
	}

	req.Password = string(hashPassword)
	if err := passwordStore.AddNewPassword(ctx, parseChangeRequestPassword(req), usr.ID); err != nil {
		return err
	}
	return nil
}

func parseChangeRequestPassword(u *ChangePasswordRequest) *datastore.Password {
	return &datastore.Password{
		Password: u.Password,
	}
}

func (request *ChangePasswordRequest) Bind(r *http.Request) error {
	if strings.TrimSpace(request.PasswordCurrent) == "" {
		return fmt.Errorf("required current password")
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

func (*ChangePasswordRequest) Render(w http.ResponseWriter, r *http.Request) {}