package status

import (
	"context"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/myerror"
	"github.com/rs/zerolog/log"
)

func (u *UsecaseStatus) CheckOut(ctx context.Context, body model.CheckOutRequest) (model.CheckOutResponse, error) {
	var result model.CheckOutResponse

	status, err := u.status.CheckOut(ctx, body)
	if err != nil {
		log.Error().Err(err).Stack().Msg("Error adding status")
		return result, myerror.New(err.Error(), "STATUS-USC-01")
	}

	result.Status = status

	return result, nil
}
