package me

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/aemiralfath/IH-Userland-Onboard/api/helper"
	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
	"github.com/go-chi/render"
)

type updateProfileRequest struct {
	Fullname string `json:"fullname"`
	Location string `json:"location"`
	Bio      string `json:"bio"`
	Web      string `json:"web"`
}

func UpdateProfile(jwtAuth helper.JWTAuth, profileStore datastore.ProfileStore) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		req := &updateProfileRequest{}

		if err := render.Bind(r, req); err != nil {
			render.Render(rw, r, helper.BadRequestErrorRenderer(err))
			return
		}

		_, claims, err := helper.FromContext(ctx)
		if err != nil {
			fmt.Println(render.Render(rw, r, helper.BadRequestErrorRenderer(err)))
			return
		}

		userId := claims["id"]
		if err := profileStore.UpdateProfile(ctx, parseUpdateRequestProfile(req), userId.(float64)); err != nil {
			render.Render(rw, r, helper.InternalServerErrorRenderer(err))
			return
		}

		if err := render.Render(rw, r, helper.SuccesRenderer()); err != nil {
			render.Render(rw, r, helper.InternalServerErrorRenderer(err))
			return
		}
	}
}

func parseUpdateRequestProfile(p *updateProfileRequest) *datastore.Profile {
	return &datastore.Profile{
		Fullname: p.Fullname,
		Location: p.Location,
		Bio:      p.Bio,
		Web:      p.Web,
	}
}

func (request *updateProfileRequest) Bind(r *http.Request) error {
	if strings.TrimSpace(request.Fullname) == "" {
		return fmt.Errorf("required Fullname")
	}
	return nil
}

func (*updateProfileRequest) Render(w http.ResponseWriter, r *http.Request) {}
