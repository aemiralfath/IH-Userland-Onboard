package me

import (
	"context"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/myerror"
	"github.com/rs/zerolog/log"
)

func (u *UsecaseMe) Profile(ctx context.Context, userId string) (model.ProfileResponse, error) {
	var result model.ProfileResponse

	profile, err := u.me.Profile(ctx, userId)
	if err != nil {
		log.Error().Err(err).Stack().Msg("Error get profile")
		return result, myerror.New(err.Error(), "STATUS-USC-01")
	}

	result.Profile = profile

	return result, nil
}
