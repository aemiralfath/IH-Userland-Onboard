package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/api/crypto"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/api/helper"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/api/jwt"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/api/kafka"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/datastore"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	RequireTFA  bool       `json:"require_tfa"`
	AccessToken *jwt.Token `json:"access_token"`
}

func Login(jwtAuth jwt.JWT, crypto crypto.Crypto, kafka kafka.Kafka, userStore datastore.UserStore, sessionStore datastore.SessionStore, clientStore datastore.ClientStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		req := &LoginRequest{}
		clientName := r.Header.Get("X-API-ClientID")

		ip := helper.GetIPAddress(r)

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
				log.Error().Err(err).Stack().Msg(err.Error())
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
			log.Error().Err(err).Stack().Msg(err.Error())
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}

		client, err := clientStore.GetClientByName(ctx, clientName)
		if err != nil {
			log.Error().Err(err).Stack().Msg(err.Error())
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}

		if err := sessionStore.AddNewSession(ctx, &datastore.Session{JTI: jti, UserId: usr.ID, IsCurrent: true}, client.ID); err != nil {
			log.Error().Err(err).Stack().Msg(err.Error())
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}

		topic := "login-succeed"
		message, _ := json.Marshal(map[string]interface{}{
			"remote-ip": ip,
			"username":  usr.Email,
			"userid":    usr.ID,
		})

		if err := kafka.SendMessage(topic, message); err != nil {
			log.Error().Err(err).Stack().Msg(err.Error())
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}

		helper.CustomRender(w, http.StatusOK, loginResponse{
			RequireTFA:  false,
			AccessToken: accessToken,
		})
	}
}

func (login *LoginRequest) Bind(r *http.Request) error {
	if strings.TrimSpace(login.Email) == "" {
		return fmt.Errorf("required email")
	}

	if strings.TrimSpace(login.Password) == "" {
		return fmt.Errorf("required password")
	}

	return nil
}

func (*LoginRequest) Render(w http.ResponseWriter, r *http.Request) {}
