package auth

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/aemiralfath/IH-Userland-Onboard/api/helper"
	"github.com/aemiralfath/IH-Userland-Onboard/api/jwt"
	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
	"github.com/go-chi/render"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	RequireTFA  bool       `json:"require_tfa"`
	AccessToken *jwt.Token `json:"access_token"`
}

func Login(jwtAuth jwt.JWTAuth, crypto datastore.Crypto, userStore datastore.UserStore, sessionStore datastore.SessionStore, clientStore datastore.ClientStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		req := &loginRequest{}
		clientName := r.Header.Get("X-API-ClientID")

		if err := render.Bind(r, req); err != nil {
			render.Render(w, r, helper.BadRequestErrorRenderer(err))
			return
		}

		usr, err := userStore.GetUserByEmail(ctx, req.Email)
		if err != nil {
			if err == sql.ErrNoRows {
				render.Render(w, r, helper.BadRequestErrorRenderer(fmt.Errorf("User not found")))
				return
			} else {
				render.Render(w, r, helper.InternalServerErrorRenderer(err))
				return
			}
		}

		if err := crypto.ConfirmPassword(usr.Password, req.Password); !err {
			render.Render(w, r, helper.BadRequestErrorRenderer(fmt.Errorf("Wrong Password")))
			return
		}

		accessToken, jti, err := jwtAuth.CreateToken(usr.ID, usr.Email, jwt.AccessTokenExpiration)
		if err != nil {
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}

		client, err := clientStore.GetClientByName(ctx, clientName)
		if err != nil {
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}

		if err := sessionStore.AddNewSession(ctx, &datastore.Session{JTI: jti, UserId: usr.ID, IsCurrent: true}, client.ID); err != nil {
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}

		helper.CustomRender(w, http.StatusOK, loginResponse{
			RequireTFA:  false,
			AccessToken: accessToken,
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
