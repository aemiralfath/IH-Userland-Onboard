package me

import (
	"context"
	"fmt"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/security"
)

func (r *Repository) ChangePassword(ctx context.Context, userId string, req model.ChangePasswordRequest) error {
	user, err := r.UserStore.GetUserById(ctx, userId)
	if err != nil {
		return err
	}

	exist, user, err := r.UserStore.CheckEmailExist(ctx, user.Email)
	if err != nil {
		return err
	}

	if !exist {
		return fmt.Errorf("User not found")
	}

	if err := security.ConfirmPassword(user.Password, req.PasswordCurrent); !err {
		return fmt.Errorf("Wrong Password")
	}

	lastThreePassword, err := r.PasswordStore.GetLastThreePassword(ctx, user.ID)
	if err != nil {
		return err
	}

	for _, e := range lastThreePassword {
		if err := security.ConfirmPassword(e, req.Password); err {
			return fmt.Errorf("Password must different from last 3 password")
		}
	}

	hashPassword, err := security.HashPassword(req.Password)
	if err != nil {
		return err
	}

	if err := r.UserStore.ChangePassword(ctx, user.ID, hashPassword); err != nil {
		return err
	}

	if err := r.PasswordStore.AddNewPassword(ctx, user.ID, hashPassword); err != nil {
		return err
	}

	return nil
}
