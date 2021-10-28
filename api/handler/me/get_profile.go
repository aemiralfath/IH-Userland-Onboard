package me

import (
	"fmt"
	"net/http"

	"github.com/aemiralfath/IH-Userland-Onboard/api/helper"
	"github.com/aemiralfath/IH-Userland-Onboard/api/jwt"
	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
	"github.com/go-chi/render"
)

func GetProfile(jwtAuth jwt.JWTAuth, profileStore datastore.ProfileStore) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		_, claims, err := jwt.FromContext(ctx)
		if err != nil {
			fmt.Println(render.Render(rw, r, helper.BadRequestErrorRenderer(err)))
			return
		}

		userId := claims["userID"]
		profile, err := profileStore.GetProfile(ctx, userId.(float64))
		if err != nil {
			fmt.Println(render.Render(rw, r, helper.InternalServerErrorRenderer(err)))
			return
		}

		helper.CustomRender(rw, http.StatusOK, map[string]interface{}{
			"id":         profile.ID,
			"fullname":   profile.Fullname,
			"location":   profile.Location,
			"bio":        profile.Bio,
			"web":        profile.Web,
			"picture":    profile.Picture,
			"created_at": profile.CreatedAt,
		})
	}
}
