package auth

import (
	"context"
	"encoding/json"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/jwt"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/myerror"
	"github.com/rs/zerolog/log"
)

func (u *UsecaseAuth) Login(ctx context.Context, ip, clientName string, body model.LoginRequest) (model.LoginResponse, error) {
	var result model.LoginResponse

	jti, user, err := u.auth.Login(ctx, ip, clientName, body)
	if err != nil {
		log.Error().Err(err).Stack().Msg("Error adding status")
		return result, myerror.New(err.Error(), "STATUS-USC-01")
	}

	accessToken, err := jwt.New().CreateToken(jti, user.ID, jwt.AccessTokenExpiration)
	if err != nil {
		log.Error().Err(err).Stack().Msg("Error create jwt access token")
		return result, myerror.New("Error create jwt access token.", "AUTH-USC-03")
	}

	topic := "login-succeed"
	message, _ := json.Marshal(map[string]interface{}{
		"remote-ip": ip,
		"username":  user.Email,
		"userid":    user.ID,
	})

	if err := u.kafka.SendMessage(topic, message); err != nil {
		log.Error().Err(err).Stack().Msg(err.Error())
		return result, myerror.New("Error create kafka message.", "AUTH-USC-04")
	}

	result.RequireTFA = false
	result.AccessToken = *accessToken

	return result, nil
}
