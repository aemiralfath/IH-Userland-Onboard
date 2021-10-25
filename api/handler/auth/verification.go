package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/aemiralfath/IH-Userland-Onboard/api/helper"
	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
	"github.com/go-chi/render"
)

type verificationRequest struct {
	Type      string `json:"type"`
	Recipient string `json:"recipient"`
}

func Verification(otp datastore.OTPStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		req := &verificationRequest{}

		if err := render.Bind(r, req); err != nil {
			render.Render(w, r, helper.BadRequestErrorRenderer(err))
		}

		token, err := helper.GenerateOTP(6)
		if err != nil {
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}

		res, err := otp.GetOTP(ctx, req.Recipient, token)
		if err != nil {
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}

		subject := "Userland Email Verification!"
		msg := fmt.Sprintf("Use this otp for verify your email: %s", res)

		go helper.SendEmail(req.Recipient, subject, msg)

		if err := render.Render(w, r, helper.SuccesRenderer()); err != nil {
			fmt.Println(render.Render(w, r, helper.InternalServerErrorRenderer(err)))
			return
		}
	}
}

func (verification *verificationRequest) Bind(r *http.Request) error {
	if strings.TrimSpace(verification.Recipient) == "" {
		return fmt.Errorf("required recipient")
	}
	return nil
}

func (verification *verificationRequest) Render(w http.ResponseWriter, r *http.Request) {}
