package me

import (
	"fmt"
	"net/http"
	"os"

	"github.com/aemiralfath/IH-Userland-Onboard/api/helper"
	"github.com/aemiralfath/IH-Userland-Onboard/api/jwt"
	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
	"github.com/go-chi/render"
)

func DeletePicture(jwtAuth jwt.JWT, profileStore datastore.ProfileStore) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		_, claims, err := jwtAuth.FromContext(ctx)
		if err != nil {
			fmt.Println(render.Render(rw, r, helper.BadRequestErrorRenderer(err)))
			return
		}

		userId := claims["userID"]
		profile, err := profileStore.GetProfile(ctx, userId.(float64))
		if err != nil {
			fmt.Println(render.Render(rw, r, helper.UnauthorizedErrorRenderer(err)))
			return
		}

		if err := os.Remove(profile.Picture); err != nil {
			render.Render(rw, r, helper.InternalServerErrorRenderer(err))
			return
		}

		profile.Picture = ""
		if err := profileStore.UpdatePicture(ctx, profile, userId.(float64)); err != nil {
			render.Render(rw, r, helper.InternalServerErrorRenderer(err))
			return
		}

		if err := render.Render(rw, r, helper.SuccesRenderer()); err != nil {
			render.Render(rw, r, helper.InternalServerErrorRenderer(err))
			return
		}
	}
}
