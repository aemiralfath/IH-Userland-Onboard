package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/aemiralfath/IH-Userland-Onboard/api/helper"
	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
	"github.com/go-chi/render"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(jwtAuth helper.JWTAuth, userStore datastore.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		req := &loginRequest{}

		if err := render.Bind(r, req); err != nil {
			render.Render(w, r, helper.BadRequestErrorRenderer(err))
			return
		}

		usr, err := userStore.GetUserByEmail(ctx, req.Email)
		if usr == nil {
			render.Render(w, r, helper.BadRequestErrorRenderer(err))
			return
		} else if err != nil {
			fmt.Println(render.Render(w, r, helper.InternalServerErrorRenderer(err)))
			return
		}

		if err := helper.ConfirmPassword(usr.Password, req.Password); err != nil {
			render.Render(w, r, helper.BadRequestErrorRenderer(err))
			return
		}

		accessToken, err := jwtAuth.CreateToken(usr.ID, usr.Email, helper.AccessTokenExpiration)
		if err != nil {
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}

		helper.CustomRender(w, http.StatusOK, map[string]interface{}{
			"require_tfa":  false,
			"access_token": accessToken,
		})
	}
}

func (login *loginRequest) Bind(r *http.Request) error {
	if strings.TrimSpace(login.Email) == "" {
		return fmt.Errorf("required email")
	}

	if strings.TrimSpace(login.Password) == "" {
		return fmt.Errorf("required password")
	}

	return nil
}

func (*loginRequest) Render(w http.ResponseWriter, r *http.Request) {}
