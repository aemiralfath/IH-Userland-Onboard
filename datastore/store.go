package datastore

import (
	"context"

	"github.com/aemiralfath/IH-Userland-Onboard/datastore/models"
)

type ProfileStore interface {
	GetProfile(ctx context.Context) error
}

type UserStore interface {
	GetUser(ctx context.Context) error
	AddNewUser(ctx context.Context, user *models.User) error
}

type PasswordStore interface {
	GetPassword(ctx context.Context) error
}
