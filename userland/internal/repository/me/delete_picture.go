package me

import "context"

func (r *Repository) DeletePicture(ctx context.Context, userId string) (string, error) {
	user, err := r.UserStore.GetUserById(ctx, userId)
	if err != nil {
		return "", err
	}

	profile, err := r.ProfileStore.GetProfileById(ctx, user.ProfileId)
	if err != nil {
		return "", err
	}

	err = r.ProfileStore.SetPicture(ctx, user.ProfileId, "")
	if err != nil {
		return "", err
	}

	return profile.Picture, nil
}
