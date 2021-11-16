package me

import (
	"context"
	"fmt"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/security"
)

func (r *Repository) DeleteAccount(ctx context.Context, userId string, req model.DeleteAccountRequest) error {
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

	if err := security.ConfirmPassword(user.Password, req.Password); !err {
		return fmt.Errorf("Wrong Password")
	}

	if err := r.UserStore.SoftDeleteUser(ctx, user.Email); err != nil {
		return err
	}

	return nil
}
