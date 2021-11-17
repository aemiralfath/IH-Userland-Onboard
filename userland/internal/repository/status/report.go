package status

import (
	"context"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model/entity"
)

func (r *Repository) Report(ctx context.Context, placeId string) (entity.Place, []entity.Report, error) {
	var place entity.Place
	var reports []entity.Report

	_, place, err := r.PlaceStore.CheckPlaceExist(ctx, placeId)
	if err != nil {
		return place, reports, err
	}

	reports, err = r.StatusStore.GetTodayStatus(ctx, placeId)
	if err != nil {
		return place, reports, err
	}

	return place, reports, nil
}
