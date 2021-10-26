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

func Verification(token datastore.TokenStore, userStore datastore.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		req := &verificationRequest{}

		if err := render.Bind(r, req); err != nil {
			render.Render(w, r, helper.BadRequestErrorRenderer(err))
		}

		usr, err := userStore.CheckUserEmailExist(ctx, req.Recipient)
		if err != nil {
			render.Render(w, r, helper.BadRequestErrorRenderer(err))
			return
		}

		tokenCode, err := helper.GenerateOTP(6)
		if err != nil {
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}

		value := fmt.Sprintf("%f-%s", usr.ID, req.Recipient)
		if err := token.SetToken(ctx, "user", value, tokenCode); err != nil {
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}

		subject := "Userland Email Verification!"
		msg := fmt.Sprintf("Use this otp for verify your email: %s", tokenCode)

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
