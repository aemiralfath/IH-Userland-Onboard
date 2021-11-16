package auth

import (
	"context"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/myerror"
	"github.com/rs/zerolog/log"
)

func (u *UsecaseAuth) Verification(ctx context.Context, body model.VerificationRequest) (model.VerificationResponse, error) {
	var result model.VerificationResponse

	err := u.auth.Verification(ctx, body)
	if err != nil {
		log.Error().Err(err).Stack().Msg("Error adding status")
		return result, myerror.New(err.Error(), "STATUS-USC-01")
	}

	result.Success = true

	return result, nil
}
