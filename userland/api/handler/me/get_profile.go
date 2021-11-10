package me

import (
	"fmt"
	"net/http"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/api/helper"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/api/jwt"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/datastore"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

func GetProfile(jwtAuth jwt.JWT, profileStore datastore.ProfileStore) http.HandlerFunc {
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
			log.Error().Err(err).Stack().Msg(err.Error())
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