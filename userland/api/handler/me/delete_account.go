package me

import (
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

type deleteAccountRequest struct {
	Password string `json:"password"`
}

func DeleteAccount(jwtAuth jwt.JWT, crypto crypto.Crypto, userStore datastore.UserStore) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		req := &deleteAccountRequest{}

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

		if err := crypto.ConfirmPassword(usr.Password, req.Password); !err {
			render.Render(rw, r, helper.BadRequestErrorRenderer(fmt.Errorf("Wrong Password")))
			return
		}

		if err := userStore.SoftDeleteUser(ctx, emailUser.(string)); err != nil {
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

func (request *deleteAccountRequest) Bind(r *http.Request) error {
	if strings.TrimSpace(request.Password) == "" {
		return fmt.Errorf("required password")
	}

	return nil
}

func (*deleteAccountRequest) Render(w http.ResponseWriter, r *http.Request) {}
