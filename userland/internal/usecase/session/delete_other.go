package session

import (
	"context"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/myerror"
	"github.com/rs/zerolog/log"
)

func (u *UsecaseSession) DeleteOtherSession(ctx context.Context, id string) (model.DeleteOtherResponse, error) {
	var result model.DeleteOtherResponse

	err := u.session.DeleteOtherSession(ctx, id)
	if err != nil {
		log.Error().Err(err).Stack().Msg("Error delete other sessions")
		return result, myerror.New(err.Error(), "STATUS-USC-01")
	}

	result.Success = true

	return result, nil
}
