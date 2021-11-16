package me

import (
	"context"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/myerror"
	"github.com/rs/zerolog/log"
)

func (u *UsecaseMe) ChangeEmail(ctx context.Context, userId string, body model.ChangeEmailRequest) (model.ChangeEmailResponse, error) {
	var result model.ChangeEmailResponse

	err := u.me.ChangeEmail(ctx, userId, body)
	if err != nil {
		log.Error().Err(err).Stack().Msg("Error update profile")
		return result, myerror.New(err.Error(), "STATUS-USC-01")
	}

	result.Success = true

	return result, nil
}
