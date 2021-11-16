package session

import (
	"context"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/jwt"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/myerror"
	"github.com/rs/zerolog/log"
)

func (u *UsecaseSession) AccessToken(ctx context.Context, jti, userId string) (model.AccessTokenResponse, error) {
	var result model.AccessTokenResponse

	accessToken, err := jwt.New().CreateToken(jti, userId, jwt.AccessTokenExpiration)
	if err != nil {
		log.Error().Err(err).Stack().Msg("Error create jwt access token")
		return result, myerror.New("Error create jwt access token.", "AUTH-USC-03")
	}

	result.AccessToken = *accessToken

	return result, nil
}
