package auth

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/aemiralfath/IH-Userland-Onboard/api/crypto"
	"github.com/aemiralfath/IH-Userland-Onboard/api/email"
	"github.com/aemiralfath/IH-Userland-Onboard/api/helper"
	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
	"github.com/go-chi/render"
)

type ForgotPasswordRequest struct {
	Email string `json:"Email"`
}

func ForgotPassword(email email.Email, crypto crypto.Crypto, userStore datastore.UserStore, otp datastore.OTPStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		req := &ForgotPasswordRequest{}

		if err := render.Bind(r, req); err != nil {
			render.Render(w, r, helper.BadRequestErrorRenderer(err))
			return
		}

		_, err := userStore.GetUserByEmail(ctx, req.Email)
		if err != nil {
			if err == sql.ErrNoRows {
				render.Render(w, r, helper.BadRequestErrorRenderer(fmt.Errorf("User not found")))
				return
			} else {
				render.Render(w, r, helper.InternalServerErrorRenderer(err))
				return
			}
		}

		otpCode, err := crypto.GenerateOTP(6)
		if err != nil {
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}

		if err := otp.SetOTP(ctx, "password", otpCode, req.Email); err != nil {
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}

		subject := "Userland Reset Password!"
		msg := fmt.Sprintf("Use this otp for reset your password: %s", otpCode)

		go email.SendEmail(req.Email, subject, msg)

		if err := render.Render(w, r, helper.SuccesRenderer()); err != nil {
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}
	}
}

func (request *ForgotPasswordRequest) Bind(r *http.Request) error {
	if strings.TrimSpace(request.Email) == "" {
		return fmt.Errorf("required email")
	}
	return nil
}

func (*ForgotPasswordRequest) Render(w http.ResponseWriter, r *http.Request) {}
