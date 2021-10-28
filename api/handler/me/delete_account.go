package me

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/aemiralfath/IH-Userland-Onboard/api/helper"
	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
	"github.com/go-chi/render"
)

type deleteAccountRequest struct {
	Password string `json:"password"`
}

func DeleteAccount(jwtAuth helper.JWTAuth, userStore datastore.UserStore) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		req := &deleteAccountRequest{}

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

		if err := helper.ConfirmPassword(usr.Password, req.Password); err != nil {
			render.Render(rw, r, helper.BadRequestErrorRenderer(err))
			return
		}

		if err := userStore.SafeDeleteUser(ctx, emailUser.(string)); err != nil {
			render.Render(rw, r, helper.InternalServerErrorRenderer(err))
			return
		}

		if err := render.Render(rw, r, helper.SuccesRenderer()); err != nil {
			render.Render(rw, r, helper.InternalServerErrorRenderer(err))
			return
		}
	}
}

func (request *deleteAccountRequest) Bind(r *http.Request) error {
	if strings.TrimSpace(request.Password) == "" {
		return fmt.Errorf("required password")
	}

	return nil
}

func (*deleteAccountRequest) Render(w http.ResponseWriter, r *http.Request) {}
