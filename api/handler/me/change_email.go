package me

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/aemiralfath/IH-Userland-Onboard/api/crypto"
	"github.com/aemiralfath/IH-Userland-Onboard/api/email"
	"github.com/aemiralfath/IH-Userland-Onboard/api/helper"
	"github.com/aemiralfath/IH-Userland-Onboard/api/jwt"
	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

type changeEmailRequest struct {
	Email string `json:"email"`
}

func ChangeEmail(jwtAuth jwt.JWT, crypto crypto.Crypto, email email.Email, userStore datastore.UserStore, otp datastore.OTPStore) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		req := &changeEmailRequest{}

		if err := render.Bind(r, req); err != nil {
			render.Render(rw, r, helper.BadRequestErrorRenderer(err))
			return
		}

		_, claims, err := jwt.FromContext(ctx)
		if err != nil {
			fmt.Println(render.Render(rw, r, helper.BadRequestErrorRenderer(err)))
			return
		}

		_, err = userStore.CheckUserEmailExist(ctx, req.Email)
		if err != nil && err != sql.ErrNoRows {
			log.Error().Err(err).Stack().Msg(err.Error())
			render.Render(rw, r, helper.InternalServerErrorRenderer(err))
			return
		} else if err == nil {
			fmt.Println(render.Render(rw, r, helper.BadRequestErrorRenderer(fmt.Errorf("Email already used"))))
			return
		}

		otpCode, err := crypto.GenerateOTP(6)
		if err != nil {
			render.Render(rw, r, helper.InternalServerErrorRenderer(err))
			return
		}

		userId := claims["userID"]
		otpValue := fmt.Sprintf("%f-%s", userId.(float64), req.Email)
		if err := otp.SetOTP(ctx, "user", otpCode, otpValue); err != nil {
			render.Render(rw, r, helper.InternalServerErrorRenderer(err))
			return
		}

		subject := "Userland Change Email Verification!"
		msg := fmt.Sprintf("Use this otp for verify your new email: %s", otpCode)

		go email.SendEmail(req.Email, subject, msg)

		if err := render.Render(rw, r, helper.SuccesRenderer()); err != nil {
			render.Render(rw, r, helper.InternalServerErrorRenderer(err))
			return
		}

	}
}

func (request *changeEmailRequest) Bind(r *http.Request) error {
	if strings.TrimSpace(request.Email) == "" {
		return fmt.Errorf("required Email")
	}
	return nil
}

func (*changeEmailRequest) Render(w http.ResponseWriter, r *http.Request) {}
