package me

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/aemiralfath/IH-Userland-Onboard/api/helper"
	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
	"github.com/go-chi/render"
)

type changePasswordRequest struct {
	PasswordCurrent string `json:"password_current"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

func ChangePassword(jwtAuth helper.JWTAuth, userStore datastore.UserStore, passwordStore datastore.PasswordStore) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		req := &changePasswordRequest{}

		if err := render.Bind(r, req); err != nil {
			render.Render(rw, r, helper.BadRequestErrorRenderer(err))
			return
		}

		_, claims, err := helper.FromContext(ctx)
		if err != nil {
			fmt.Println(render.Render(rw, r, helper.BadRequestErrorRenderer(err)))
			return
		}

		emailUser := claims["email"]
		usr, err := userStore.GetUserByEmail(ctx, emailUser.(string))
		if usr == nil {
			render.Render(rw, r, helper.BadRequestErrorRenderer(err))
			return
		} else if err != nil {
			fmt.Println(render.Render(rw, r, helper.InternalServerErrorRenderer(err)))
			return
		}

		if err := helper.ConfirmPassword(usr.Password, req.PasswordCurrent); err != nil {
			render.Render(rw, r, helper.BadRequestErrorRenderer(err))
			return
		}

		lastThreePassword, err := passwordStore.GetLastThreePassword(ctx, usr.ID)
		if err != nil {
			render.Render(rw, r, helper.InternalServerErrorRenderer(err))
			return
		}

		diffPassword := 0
		for _, e := range lastThreePassword {
			if err := helper.ConfirmPassword(e, req.Password); err != nil {
				diffPassword += 1
			}
		}

		if diffPassword == len(lastThreePassword) {

			if err := updateStore(ctx, req, usr, userStore, passwordStore); err != nil {
				render.Render(rw, r, helper.InternalServerErrorRenderer(err))
				return
			}

			if err := render.Render(rw, r, helper.SuccesRenderer()); err != nil {
				render.Render(rw, r, helper.InternalServerErrorRenderer(err))
				return
			}
		} else {
			render.Render(rw, r, helper.BadRequestErrorRenderer(fmt.Errorf("Password must different from last 3 password")))
			return
		}

	}
}

func updateStore(ctx context.Context, req *changePasswordRequest, usr *datastore.User, userStore datastore.UserStore, passwordStore datastore.PasswordStore) error {
	hashPassword, err := helper.HashPassword(req.Password)
	if err != nil {
		return err
	}

	usr.Password = string(hashPassword)
	if err := userStore.ChangePassword(ctx, usr); err != nil {
		return err
	}

	req.Password = string(hashPassword)
	if err := passwordStore.AddNewPassword(ctx, parseChangeRequestPassword(req), usr.ID); err != nil {
		return err
	}
	return nil
}

func parseChangeRequestPassword(u *changePasswordRequest) *datastore.Password {
	return &datastore.Password{
		Password: u.Password,
	}
}

func (request *changePasswordRequest) Bind(r *http.Request) error {
	if strings.TrimSpace(request.PasswordCurrent) == "" {
		return fmt.Errorf("required current password")
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

func (*changePasswordRequest) Render(w http.ResponseWriter, r *http.Request) {}
