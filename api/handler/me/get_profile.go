package me

import (
	"fmt"
	"net/http"

	"github.com/aemiralfath/IH-Userland-Onboard/api/handler"
	"github.com/aemiralfath/IH-Userland-Onboard/api/helper"
	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
	"github.com/go-chi/render"
)

func GetProfile(jwtAuth helper.JWTAuth, profileStore datastore.ProfileStore) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		_, claims, err := helper.FromContext(ctx)
		if err != nil {
			fmt.Println(render.Render(rw, r, handler.BadRequestErrorRenderer(err)))
			return
		}

		userId := claims["id"]
		profile, err := profileStore.GetProfile(ctx, userId.(float64))
		if err != nil {
			fmt.Println(render.Render(rw, r, handler.InternalServerErrorRenderer(err)))
			return
		}

		handler.CustomRender(rw, http.StatusOK, map[string]interface{}{
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
