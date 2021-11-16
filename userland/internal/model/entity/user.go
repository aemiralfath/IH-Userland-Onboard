package entity

import "context"

type User struct {
	ID        string `json:"id"`
	ProfileId string `json:"profile_id"`
	Email     string `json:"email" validate:"required"`
	Password  string `json:"password" validate:"required"`
	Verified  bool   `json:"verified" validate:"required"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	DeletedAt string `json:"deletedAt"`
}

type UserStore interface {
	CheckEmailExist(ctx context.Context, email string) (bool, User, error)
	AddNewUser(ctx context.Context, profileId, email, password string) (string, error)
	ChangePassword(ctx context.Context, id, password string) error
	GetUserById(ctx context.Context, id string) (User, error)
	SoftDeleteUser(ctx context.Context, email string) error
}
