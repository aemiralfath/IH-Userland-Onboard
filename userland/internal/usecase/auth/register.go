package auth

import (
	"context"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/myerror"
	"github.com/rs/zerolog/log"
)

func (u *UsecaseAuth) Register(ctx context.Context, body model.RegisterRequest) (model.RegisterResponse, error) {
	var result model.RegisterResponse

	err := u.auth.Register(ctx, body)
	if err != nil {
		log.Error().Err(err).Stack().Msg("Error adding status")
		return result, myerror.New(err.Error(), "STATUS-USC-01")
	}

	result.Success = true

	return result, nil
}
