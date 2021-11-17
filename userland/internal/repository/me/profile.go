package me

import (
	"context"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model/entity"
)

func (r *Repository) Profile(ctx context.Context, userId string) (entity.Profile, error) {
	var profile entity.Profile
	user, err := r.UserStore.GetUserById(ctx, userId)
	if err != nil {
		return profile, err
	}

	profile, err = r.ProfileStore.GetProfileById(ctx, user.ProfileId)
	if err != nil {
		return profile, err
	}

	return profile, nil
}
