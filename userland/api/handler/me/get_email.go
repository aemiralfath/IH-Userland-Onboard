package me

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/api/helper"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/api/jwt"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/datastore"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

func GetEmail(jwtAuth jwt.JWT, userStore datastore.UserStore) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		_, claims, err := jwtAuth.FromContext(ctx)
		if err != nil {
			fmt.Println(render.Render(rw, r, helper.BadRequestErrorRenderer(err)))
			return
		}

		userId := claims["userID"]
		email, err := userStore.GetEmailByID(ctx, userId.(float64))
		if err != nil {
			if err == sql.ErrNoRows {
				render.Render(rw, r, helper.BadRequestErrorRenderer(fmt.Errorf("User not found")))
				return
			} else {
				log.Error().Err(err).Stack().Msg(err.Error())
				render.Render(rw, r, helper.InternalServerErrorRenderer(err))
				return
			}
		}

		helper.CustomRender(rw, http.StatusOK, map[string]interface{}{
			"email": email,
		})
	}
}
