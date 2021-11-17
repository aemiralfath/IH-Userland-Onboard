package session

import (
	"fmt"
	"net/http"

	"github.com/aemiralfath/IH-Userland-Onboard/api/helper"
	"github.com/aemiralfath/IH-Userland-Onboard/api/jwt"
	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

type clientResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type sessionResponse struct {
	IsCurrent bool           `json:"is_current"`
	IP        string         `json:"ip"`
	CreatedAt string         `json:"created_at"`
	UpdatedAt string         `json:"updated_at"`
	Client    clientResponse `json:"client"`
}

type listSessionResponse struct {
	Data []sessionResponse `json:"data"`
}

func GetListSession(jwtAuth jwt.JWTAuth, sessionStore datastore.SessionStore) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		_, claims, err := jwt.FromContext(ctx)
		if err != nil {
			fmt.Println(render.Render(rw, r, helper.BadRequestErrorRenderer(err)))
			return
		}

		userId := claims["userID"]
		sessions, err := sessionStore.GetUserSession(ctx, userId.(float64))
		if err != nil {
			log.Error().Err(err).Stack().Msg(err.Error())
			render.Render(rw, r, helper.InternalServerErrorRenderer(err))
			return
		}

		listSessions := &listSessionResponse{}
		for _, e := range sessions {
			log.Info().Msg(fmt.Sprintf("%f", e.ClientId))
			listSessions.Data = append(listSessions.Data, sessionResponse{
				IsCurrent: e.IsCurrent,
				IP:        e.IP,
				CreatedAt: e.CreatedAt,
				UpdatedAt: e.UpdatedAt,
				Client: clientResponse{
					ID:   int64(e.Client.ID),
					Name: e.Client.Name,
				},
			})
		}

		helper.CustomRender(rw, http.StatusOK, &listSessions)
	}
}
