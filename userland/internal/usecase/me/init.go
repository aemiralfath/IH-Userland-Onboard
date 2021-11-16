package me

import (
	"context"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model/entity"
)

type meRepo interface {
	Profile(ctx context.Context, userId string) (entity.Profile, error)
	UpdateProfile(ctx context.Context, userId string, body model.UpdateProfileRequest) error
	Email(ctx context.Context, userId string) (string, error)
	ChangeEmail(ctx context.Context, userId string, body model.ChangeEmailRequest) error
	ChangePassword(ctx context.Context, userId string, body model.ChangePasswordRequest) error
	DeleteAccount(ctx context.Context, userId string, body model.DeleteAccountRequest) error
	SetPicture(ctx context.Context, userId, fileName string) error
	DeletePicture(ctx context.Context, userId string) (string, error)
}

type UsecaseMe struct {
	me meRepo
}

func New(repo meRepo) *UsecaseMe {
	return &UsecaseMe{
		me: repo,
	}
}
