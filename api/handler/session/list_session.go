package session

import (
	"fmt"
	"net/http"

	"github.com/aemiralfath/IH-Userland-Onboard/api/helper"
	"github.com/aemiralfath/IH-Userland-Onboard/api/jwt"
	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
	"github.com/go-chi/render"
)

func GetListSession(jwtAuth jwt.JWTAuth, sessionStore datastore.SessionStore, clientStore datastore.ClientStore) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		_, claims, err := jwt.FromContext(ctx)
		if err != nil {
			fmt.Println(render.Render(rw, r, helper.BadRequestErrorRenderer(err)))
			return
		}

		userId := claims["userID"]
		fmt.Println(userId)
	}
}
