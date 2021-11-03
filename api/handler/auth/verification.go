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

type VerificationRequest struct {
	Type      string `json:"type"`
	Recipient string `json:"recipient"`
}

func Verification(email email.Email, crypto crypto.Crypto, otp datastore.OTPStore, userStore datastore.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		req := &VerificationRequest{}

		if err := render.Bind(r, req); err != nil {
			render.Render(w, r, helper.BadRequestErrorRenderer(err))
		}

		usr, err := userStore.CheckUserEmailExist(ctx, req.Recipient)
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

		otpValue := fmt.Sprintf("%f-%s", usr.ID, req.Recipient)
		if err := otp.SetOTP(ctx, "user", otpCode, otpValue); err != nil {
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}

		subject := "Userland Email Verification!"
		msg := fmt.Sprintf("Use this otp for verify your email: %s", otpCode)

		go email.SendEmail(req.Recipient, subject, msg)

		if err := render.Render(w, r, helper.SuccesRenderer()); err != nil {
			fmt.Println(render.Render(w, r, helper.InternalServerErrorRenderer(err)))
			return
		}
	}
}

func (verification *VerificationRequest) Bind(r *http.Request) error {
	if strings.TrimSpace(verification.Recipient) == "" {
		return fmt.Errorf("required recipient")
	}
	return nil
}

func (verification *VerificationRequest) Render(w http.ResponseWriter, r *http.Request) {}
