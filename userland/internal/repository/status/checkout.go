package status

import (
	"context"
	"fmt"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model/entity"
)

func (r *Repository) CheckOut(ctx context.Context, req model.CheckOutRequest) (entity.Status, error) {
	var status entity.Status

	statusExist, status, err := r.StatusStore.CheckStatusExistById(ctx, req.StatusID)
	if err != nil {
		return status, err
	}

	if !statusExist {
		return status, fmt.Errorf("Status not exist")
	}

	_, place, err := r.PlaceStore.CheckPlaceExist(ctx, status.PlaceId)
	if err != nil {
		return status, err
	}

	statusExist, statusTemp, err := r.StatusStore.CheckStatusExist(ctx, status.ProfileId, status.PlaceId)
	if err != nil {
		return status, err
	}

	if statusExist && status.ID == statusTemp.ID {
		status, err = r.StatusStore.UpdateCheckOut(ctx, req.StatusID)
		if err != nil {
			return status, err
		}

		place, err = r.PlaceStore.UpdateCurrentCapacity(ctx, place.ID, place.CurrentCapacity-1)
		if err != nil {
			return status, err
		}
	}

	return status, nil
}
