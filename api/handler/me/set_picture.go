package me

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/aemiralfath/IH-Userland-Onboard/api/helper"
	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
	"github.com/go-chi/render"
)

func SetPicture(jwtAuth helper.JWTAuth, profileStore datastore.ProfileStore) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		_, claims, err := helper.FromContext(ctx)
		if err != nil {
			fmt.Println(render.Render(rw, r, helper.BadRequestErrorRenderer(err)))
			return
		}

		userId := claims["userID"]
		profile, err := profileStore.GetProfile(ctx, userId.(float64))
		if err != nil {
			fmt.Println(render.Render(rw, r, helper.UnauthorizedErrorRenderer(err)))
			return
		}

		if err := r.ParseMultipartForm(32 << 20); err != nil {
			render.Render(rw, r, helper.BadRequestErrorRenderer(err))
			return
		}

		file, handler, err := r.FormFile("file")
		if err != nil {
			render.Render(rw, r, helper.BadRequestErrorRenderer(err))
			return
		}
		defer file.Close()

		fileName := fmt.Sprintf("/ih-userland-onboard/assets/profile/%f-%s", userId, handler.Filename)
		localFile, err := os.Create(fileName)
		if err != nil {
			render.Render(rw, r, helper.BadRequestErrorRenderer(err))
			return
		}
		defer localFile.Close()

		if _, err := io.Copy(localFile, file); err != nil {
			render.Render(rw, r, helper.InternalServerErrorRenderer(err))
			return
		}

		profile.Picture = fileName
		if err := profileStore.UpdatePicture(ctx, profile, userId.(float64)); err != nil {
			render.Render(rw, r, helper.InternalServerErrorRenderer(err))
			return
		}

		if err := render.Render(rw, r, helper.SuccesRenderer()); err != nil {
			render.Render(rw, r, helper.InternalServerErrorRenderer(err))
			return
		}

	}
}
