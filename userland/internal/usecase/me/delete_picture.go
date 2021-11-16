package me

import (
	"context"
	"os"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/myerror"
	"github.com/rs/zerolog/log"
)

func (u *UsecaseMe) DeletePicture(ctx context.Context, userId string) (model.DeletePictureResponse, error) {
	var result model.DeletePictureResponse

	filename, err := u.me.DeletePicture(ctx, userId)
	if err != nil {
		log.Error().Err(err).Stack().Msg("Error remove picture")
		return result, myerror.New(err.Error(), "STATUS-USC-01")
	}

	if err := os.Remove(filename); err != nil {
		log.Error().Err(err).Stack().Msg("Error remove picture")
		return result, myerror.New(err.Error(), "STATUS-USC-02")
	}

	result.Success = true

	return result, nil
}
