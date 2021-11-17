package me

import (
	"context"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/myerror"
	"github.com/rs/zerolog/log"
)

func (u *UsecaseMe) Email(ctx context.Context, userId string) (model.EmailResponse, error) {
	var result model.EmailResponse

	email, err := u.me.Email(ctx, userId)
	if err != nil {
		log.Error().Err(err).Stack().Msg("Error get email")
		return result, myerror.New(err.Error(), "STATUS-USC-01")
	}

	result.Email = email

	return result, nil
}
