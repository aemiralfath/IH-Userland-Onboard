package me

import (
	"fmt"
	"net/http"

	"github.com/aemiralfath/IH-Userland-Onboard/api/helper"
	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
	"github.com/go-chi/render"
)

func GetEmail(jwtAuth helper.JWTAuth, userStore datastore.UserStore) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		_, claims, err := helper.FromContext(ctx)
		if err != nil {
			fmt.Println(render.Render(rw, r, helper.BadRequestErrorRenderer(err)))
			return
		}

		userId := claims["userID"]
		email, err := userStore.GetEmailByID(ctx, userId.(float64))
		if err != nil {
			fmt.Println(render.Render(rw, r, helper.InternalServerErrorRenderer(err)))
			return
		}

		helper.CustomRender(rw, http.StatusOK, map[string]interface{}{
			"email": email,
		})
	}
}
