package session

import (
	"fmt"
	"net/http"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/api/helper"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/api/jwt"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

func GetRefreshToken(jwtAuth jwt.JWT) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		_, claims, err := jwtAuth.FromContext(ctx)
		if err != nil {
			fmt.Println(render.Render(rw, r, helper.BadRequestErrorRenderer(err)))
			return
		}

		userId := claims["userID"]
		userEmail := claims["email"]

		refreshToken, _, err := jwtAuth.CreateToken(userId.(float64), userEmail.(string), jwt.RefreshTokenExpiration)
		if err != nil {
			log.Error().Err(err).Stack().Msg(err.Error())
			render.Render(rw, r, helper.InternalServerErrorRenderer(err))
			return
		}

		helper.CustomRender(rw, http.StatusOK, map[string]interface{}{
			"refresh_token": refreshToken,
		})
	}
}
