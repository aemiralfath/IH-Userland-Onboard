package me

import (
	"context"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model"
)

func (r *Repository) UpdateProfile(ctx context.Context, userId string, body model.UpdateProfileRequest) error {
	user, err := r.UserStore.GetUserById(ctx, userId)
	if err != nil {
		return err
	}

	err = r.ProfileStore.UpdateProfile(ctx, user.ProfileId, body.Fullname, body.DosageType)
	if err != nil {
		return err
	}

	return nil
}
