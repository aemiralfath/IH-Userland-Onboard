package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/aemiralfath/IH-Userland-Onboard/api/helper"
	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
	"github.com/go-chi/render"
)

type forgotPasswordRequest struct {
	Email string `json:"Email"`
}

func ForgotPassword(email helper.Email, userStore datastore.UserStore, token datastore.TokenStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		req := &forgotPasswordRequest{}

		if err := render.Bind(r, req); err != nil {
			render.Render(w, r, helper.BadRequestErrorRenderer(err))
			return
		}

		usr, err := userStore.GetUserByEmail(ctx, req.Email)
		if usr == nil {
			render.Render(w, r, helper.BadRequestErrorRenderer(err))
			return
		}

		if err != nil {
			fmt.Println(render.Render(w, r, helper.InternalServerErrorRenderer(err)))
			return
		}

		tokenCode, err := helper.GenerateOTP(6)
		if err != nil {
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}

		if err := token.SetToken(ctx, "password", req.Email, tokenCode); err != nil {
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}

		subject := "Userland Reset Password!"
		msg := fmt.Sprintf("Use this token for reset your password: %s", tokenCode)

		go email.SendEmail(req.Email, subject, msg)

		if err := render.Render(w, r, helper.SuccesRenderer()); err != nil {
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}
	}
}

func (request *forgotPasswordRequest) Bind(r *http.Request) error {
	if strings.TrimSpace(request.Email) == "" {
		return fmt.Errorf("required email")
	}
	return nil
}

func (*forgotPasswordRequest) Render(w http.ResponseWriter, r *http.Request) {}
