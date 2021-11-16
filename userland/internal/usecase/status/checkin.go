package status

import (
	"context"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/myerror"
	"github.com/rs/zerolog/log"
)

func (u *UsecaseStatus) CheckIn(ctx context.Context, body model.CheckInRequest) (model.CheckInResponse, error) {
	var result model.CheckInResponse

	status, profile, place, err := u.status.CheckIn(ctx, body)
	if err != nil {
		log.Error().Err(err).Stack().Msg("Error adding status")
		return result, myerror.New(err.Error(), "STATUS-USC-01")
	}

	result.Status = status
	result.Profile = profile
	result.Place = place

	return result, nil
}
