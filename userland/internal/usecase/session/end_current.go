package session

import (
	"context"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/myerror"
	"github.com/rs/zerolog/log"
)

func (u *UsecaseSession) EndCurrentSession(ctx context.Context, id string) (model.EndCurrentResponse, error) {
	var result model.EndCurrentResponse

	err := u.session.EndCurrentSession(ctx, id)
	if err != nil {
		log.Error().Err(err).Stack().Msg("Error get sessions")
		return result, myerror.New(err.Error(), "STATUS-USC-01")
	}

	result.Success = true

	return result, nil
}
