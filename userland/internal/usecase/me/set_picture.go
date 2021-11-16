package me

import (
	"context"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/myerror"
	"github.com/rs/zerolog/log"
)

func (u *UsecaseMe) SetPicture(ctx context.Context, userId, fileName string) (model.SetPictureResponse, error) {
	var result model.SetPictureResponse

	err := u.me.SetPicture(ctx, userId, fileName)
	if err != nil {
		log.Error().Err(err).Stack().Msg("Error set picture")
		return result, myerror.New(err.Error(), "STATUS-USC-01")
	}

	result.Success = true

	return result, nil
}
