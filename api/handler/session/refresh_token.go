package session

import (
	"fmt"
	"net/http"

	"github.com/aemiralfath/IH-Userland-Onboard/api/helper"
	"github.com/go-chi/render"
)

func GetRefreshToken(jwtAuth helper.JWTAuth) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		_, claims, err := helper.FromContext(ctx)
		if err != nil {
			fmt.Println(render.Render(rw, r, helper.BadRequestErrorRenderer(err)))
			return
		}

		userId := claims["userID"]
		userEmail := claims["email"]

		refreshToken, err := jwtAuth.CreateToken(userId.(float64), userEmail.(string), helper.RefreshTokenExpiration)
		if err != nil {
			render.Render(rw, r, helper.InternalServerErrorRenderer(err))
			return
		}

		helper.CustomRender(rw, http.StatusOK, map[string]interface{}{
			"refresh_token": refreshToken,
		})
	}
}
