package status

import (
	"context"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/myerror"
	"github.com/rs/zerolog/log"
)

func (u *UsecaseStatus) Report(ctx context.Context, placeId string) (model.ReportResponse, error) {
	var result model.ReportResponse

	place, reports, err := u.status.Report(ctx, placeId)
	if err != nil {
		log.Error().Err(err).Stack().Msg("Error getting report")
		return result, myerror.New(err.Error(), "STATUS-USC-01")
	}

	result.Place = place
	result.Reports = reports

	return result, nil
}
