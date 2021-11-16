package me

import (
	"context"
)

func (r *Repository) Email(ctx context.Context, userId string) (string, error) {
	user, err := r.UserStore.GetUserById(ctx, userId)
	if err != nil {
		return "", err
	}

	return user.Email, nil
}
