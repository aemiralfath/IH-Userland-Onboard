package auth

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/aemiralfath/IH-Userland-Onboard/api/helper"
	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
	"github.com/go-chi/render"
)

type resetPasswordRequest struct {
	Token           string `json:"token"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

func ResetPassword(userStore datastore.UserStore, passwordStore datastore.PasswordStore, otp datastore.OTPStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		req := &resetPasswordRequest{}

		if err := render.Bind(r, req); err != nil {
			render.Render(w, r, helper.BadRequestErrorRenderer(err))
			return
		}

		email, err := otp.GetOTP(ctx, "password", req.Token)
		if err != nil {
			render.Render(w, r, helper.BadRequestErrorRenderer(err))
			return
		}

		usr, err := userStore.GetUserByEmail(ctx, email)
		if err != nil {
			if err == sql.ErrNoRows {
				render.Render(w, r, helper.BadRequestErrorRenderer(fmt.Errorf("User not found")))
				return
			} else {
				render.Render(w, r, helper.InternalServerErrorRenderer(err))
				return
			}
		}

		lastThreePassword, err := passwordStore.GetLastThreePassword(ctx, usr.ID)
		if err != nil {
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}

		for _, e := range lastThreePassword {
			if err := helper.ConfirmPassword(e, req.Password); err == nil {
				render.Render(w, r, helper.BadRequestErrorRenderer(fmt.Errorf("Password must different from last 3 password")))
				return
			}
		}

		if err := updateStore(ctx, req, usr, userStore, passwordStore); err != nil {
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}

		if err := render.Render(w, r, helper.SuccesRenderer()); err != nil {
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}

	}
}

func updateStore(
	ctx context.Context,
	req *resetPasswordRequest, 
	usr *datastore.User, 
	userStore datastore.UserStore, 
	passwordStore datastore.PasswordStore) error {

	hashPassword, err := helper.HashPassword(req.Password)
	if err != nil {
		return err
	}

	usr.Password = string(hashPassword)
	if err := userStore.ChangePassword(ctx, usr); err != nil {
		return err
	}

	req.Password = string(hashPassword)
	if err := passwordStore.AddNewPassword(ctx, parseResetRequestPassword(req), usr.ID); err != nil {
		return err
	}
	return nil
}

func parseResetRequestPassword(u *resetPasswordRequest) *datastore.Password {
	return &datastore.Password{
		Password: u.Password,
	}
}

func (request *resetPasswordRequest) Bind(r *http.Request) error {
	if strings.TrimSpace(request.Token) == "" {
		return fmt.Errorf("required token")
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

func (*resetPasswordRequest) Render(w http.ResponseWriter, r *http.Request) {}
