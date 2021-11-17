package me

import "context"

func (r *Repository) SetPicture(ctx context.Context, userId, fileName string) error {
	user, err := r.UserStore.GetUserById(ctx, userId)
	if err != nil {
		return err
	}

	err = r.ProfileStore.SetPicture(ctx, user.ProfileId, fileName)
	if err != nil {
		return err
	}

	return nil
}
