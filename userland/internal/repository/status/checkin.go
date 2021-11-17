package status

import (
	"context"
	"fmt"
	"time"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model/entity"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/mytime"
)

func (r *Repository) CheckIn(ctx context.Context, req model.CheckInRequest) (entity.Status, entity.Profile, entity.Place, error) {
	var status entity.Status
	var profile entity.Profile
	var place entity.Place

	profileExist, profile, err := r.ProfileStore.CheckNIKExist(ctx, req.Profile.NIK)
	if err != nil {
		return status, profile, place, err
	}

	if !profileExist {
		profile, err = r.ProfileStore.AddNewProfile(ctx, req.Profile)
		if err != nil {
			return status, profile, place, err
		}
	}

	placeExist, place, err := r.PlaceStore.CheckPlaceExist(ctx, req.Place.ID)
	if err != nil {
		return status, profile, place, err
	}

	if !placeExist {
		place, err = r.PlaceStore.AddNewPlace(ctx, req.Place)
		if err != nil {
			return status, profile, place, err
		}
	}

	statusExist, status, err := r.StatusStore.CheckStatusExist(ctx, profile.ID, place.ID)
	if err != nil {
		return status, profile, place, err
	}

	if !statusExist {

		// Todo: Add goroutine later for race condition case
		if place.CurrentCapacity+1 > place.MaxCapacity {
			return status, profile, place, fmt.Errorf("Capacity is full")
		}

		status, err = r.StatusStore.AddNewStatus(ctx, profile.ID, place.ID)
		if err != nil {
			return status, profile, place, err
		}

		place, err = r.PlaceStore.UpdateCurrentCapacity(ctx, place.ID, place.CurrentCapacity+1)
		if err != nil {
			return status, profile, place, err
		}

	} else {
		if time.Since(status.CheckInAt.Time).Hours() > float64(place.MaxHours) || !mytime.DateEqual(time.Now(), status.CheckInAt.Time) {
			status, err = r.StatusStore.UpdateCheckOut(ctx, status.ID)
			if err != nil {
				return status, profile, place, err
			}

			place, err = r.PlaceStore.UpdateCurrentCapacity(ctx, place.ID, place.CurrentCapacity-1)
			if err != nil {
				return status, profile, place, err
			}

			return status, profile, place, fmt.Errorf("Checkin timelimit exceeded, checkout now!")
		}
	}

	return status, profile, place, nil
}
