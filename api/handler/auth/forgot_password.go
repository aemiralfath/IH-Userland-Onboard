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

func ForgotPassword(userStore datastore.UserStore, otp datastore.OTPStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		req := &forgotPasswordRequest{}

		if err := render.Bind(r, req); err != nil {
			render.Render(w, r, helper.BadRequestErrorRenderer(err))
			return
		}

		usr, err := userStore.GetUser(ctx, parseForgotPasswordRequest(req))
		if usr == nil {
			render.Render(w, r, helper.BadRequestErrorRenderer(err))
			return
		}

		if err != nil {
			fmt.Println(render.Render(w, r, helper.InternalServerErrorRenderer(err)))
			return
		}

		token, err := helper.GenerateOTP(9)
		if err != nil {
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}

		res, err := otp.GetTokenPassword(ctx, req.Email, token)
		if err != nil {
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}

		subject := "Userland Reset Password!"
		msg := fmt.Sprintf("Use this token for reset your password: %s", res)

		go helper.SendEmail(req.Email, subject, msg)

		if err := render.Render(w, r, helper.SuccesRenderer()); err != nil {
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}
	}
}

func parseForgotPasswordRequest(u *forgotPasswordRequest) *datastore.User {
	return &datastore.User{
		Email: u.Email,
	}
}

func (request *forgotPasswordRequest) Bind(r *http.Request) error {
	if strings.TrimSpace(request.Email) == "" {
		return fmt.Errorf("required email")
	}
	return nil
}

func (*forgotPasswordRequest) Render(w http.ResponseWriter, r *http.Request) {}
