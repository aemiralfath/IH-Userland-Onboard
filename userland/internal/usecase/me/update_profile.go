package me

import (
	"context"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/myerror"
	"github.com/rs/zerolog/log"
)

func (u *UsecaseMe) UpdateProfile(ctx context.Context, userId string, body model.UpdateProfileRequest) (model.UpdateProfileResponse, error) {
	var result model.UpdateProfileResponse

	err := u.me.UpdateProfile(ctx, userId, body)
	if err != nil {
		log.Error().Err(err).Stack().Msg("Error update profile")
		return result, myerror.New(err.Error(), "STATUS-USC-01")
	}

	result.Success = true

	return result, nil
}
