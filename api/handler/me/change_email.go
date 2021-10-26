package me

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/aemiralfath/IH-Userland-Onboard/api/helper"
	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
	"github.com/go-chi/render"
)

type changeEmailRequest struct {
	Email string `json:"email"`
}

func ChangeEmail(jwtAuth helper.JWTAuth, userStore datastore.UserStore, token datastore.TokenStore) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		req := &changeEmailRequest{}

		if err := render.Bind(r, req); err != nil {
			render.Render(rw, r, helper.BadRequestErrorRenderer(err))
			return
		}

		_, claims, err := helper.FromContext(ctx)
		if err != nil {
			fmt.Println(render.Render(rw, r, helper.BadRequestErrorRenderer(err)))
			return
		}

		usr, _ := userStore.CheckUserEmailExist(ctx, req.Email)
		if usr != nil {
			render.Render(rw, r, helper.BadRequestErrorRenderer(fmt.Errorf("Email already exist!")))
			return
		}

		tokenCode, err := helper.GenerateOTP(6)
		if err != nil {
			render.Render(rw, r, helper.InternalServerErrorRenderer(err))
			return
		}

		userId := claims["id"]
		value := fmt.Sprintf("%f-%s", userId.(float64), req.Email)
		if err := token.SetToken(ctx, "user", value, tokenCode); err != nil {
			render.Render(rw, r, helper.InternalServerErrorRenderer(err))
			return
		}

		subject := "Userland Change Email Verification!"
		msg := fmt.Sprintf("Use this otp for verify your new email: %s", tokenCode)

		go helper.SendEmail(req.Email, subject, msg)

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
