package auth

import (
	"context"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/myerror"
	"github.com/rs/zerolog/log"
)

func (u *UsecaseAuth) ForgotPassword(ctx context.Context, body model.ForgotPasswordRequest) (model.ForgotPasswordResponse, error) {
	var result model.ForgotPasswordResponse

	err := u.auth.ForgotPassword(ctx, body)
	if err != nil {
		log.Error().Err(err).Stack().Msg("Error Forgot Password")
		return result, myerror.New(err.Error(), "STATUS-USC-01")
	}

	result.Success = true

	return result, nil
}
