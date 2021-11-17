package me

import (
	"context"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/myerror"
	"github.com/rs/zerolog/log"
)

func (u *UsecaseMe) DeleteAccount(ctx context.Context, userId string, body model.DeleteAccountRequest) (model.DeleteAccountResponse, error) {
	var result model.DeleteAccountResponse

	err := u.me.DeleteAccount(ctx, userId, body)
	if err != nil {
		log.Error().Err(err).Stack().Msg("Error delete account")
		return result, myerror.New(err.Error(), "STATUS-USC-01")
	}

	result.Success = true

	return result, nil
}
