package auth

import (
	"context"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/myerror"
	"github.com/rs/zerolog/log"
)

func (u *UsecaseAuth) ResetPassword(ctx context.Context, body model.ResetPasswordRequest) (model.ResetPasswordResponse, error) {
	var result model.ResetPasswordResponse

	err := u.auth.ResetPassword(ctx, body)
	if err != nil {
		log.Error().Err(err).Stack().Msg("Error Reset Password")
		return result, myerror.New(err.Error(), "STATUS-USC-01")
	}

	result.Success = true

	return result, nil
}
