package me

import (
	"context"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model"
)

type meUC interface {
	Profile(ctx context.Context, userId string) (model.ProfileResponse, error)
	UpdateProfile(ctx context.Context, userId string, body model.UpdateProfileRequest) (model.UpdateProfileResponse, error)
	Email(ctx context.Context, userId string) (model.EmailResponse, error)
	ChangeEmail(ctx context.Context, userId string, body model.ChangeEmailRequest) (model.ChangeEmailResponse, error)
	ChangePassword(ctx context.Context, userId string, body model.ChangePasswordRequest) (model.ChangePasswordResponse, error)
	DeleteAccount(ctx context.Context, userId string, body model.DeleteAccountRequest) (model.DeleteAccountResponse, error)
}

type DeliveryMe struct {
	me meUC
}

func NewMe(me meUC) *DeliveryMe {
	return &DeliveryMe{
		me: me,
	}
}
