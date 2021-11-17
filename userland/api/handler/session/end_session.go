package session

import (
	"fmt"
	"net/http"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/api/helper"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/api/jwt"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/datastore"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

func EndCurrentSession(jwtAuth jwt.JWT, sessionStore datastore.SessionStore) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		_, claims, err := jwtAuth.FromContext(ctx)
		if err != nil {
			fmt.Println(render.Render(rw, r, helper.BadRequestErrorRenderer(err)))
			return
		}

		id := claims["id"]
		if err := sessionStore.EndSession(ctx, id.(string)); err != nil {
			log.Error().Err(err).Stack().Msg(err.Error())
			render.Render(rw, r, helper.InternalServerErrorRenderer(err))
			return
		}

		if err := render.Render(rw, r, helper.SuccesRenderer()); err != nil {
			log.Error().Err(err).Stack().Msg(err.Error())
			render.Render(rw, r, helper.InternalServerErrorRenderer(err))
			return
		}

	}
}
