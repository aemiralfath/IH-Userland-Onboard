package datastore

import (
	"context"
)

type User struct {
	ID        int64  `json:"id" sql:"id"`
	Email     string `json:"email" validate:"required" sql:"email"`
	Password  string `json:"password" validate:"required" sql:"password"`
	Verified  bool   `json:"verified" validate:"required" sql:"verified"`
	CreatedAt string `json:"createdAt" sql:"created_at"`
	UpdatedAt string `json:"updatedAt" sql:"updated_at"`
	DeletedAt string `json:"deletedAt" sql:"deleted_at"`
}

type Profile struct {
	ID        int64  `json:"id" sql:"id"`
	UserId    int64  `json:"userId" sql:"user_id"`
	Fullname  string `json:"fullname" validate:"required" sql:"fullname"`
	Location  string `json:"location" sql:"location"`
	Bio       string `json:"bio" sql:"bio"`
	Web       string `json:"web" sql:"web"`
	Picture   string `json:"picture" sql:"picture"`
	CreatedAt string `json:"createdAt" sql:"created_at"`
	UpdatedAt string `json:"updatedAt" sql:"updated_at"`
}

type Password struct {
	ID        int64  `json:"id" sql:"id"`
	UserId    int64  `json:"userId" sql:"user_id"`
	Password  string `json:"password" validate:"required" sql:"password"`
	CreatedAt string `json:"createdAt" sql:"created_at"`
}

type ProfileStore interface {
	GetProfile(ctx context.Context, userId float64) (*Profile, error)
}

type UserStore interface {
	GetUser(ctx context.Context, user *User) (*User, error)
	AddNewUser(ctx context.Context, user *User, profile *Profile, password *Password) error
}

type PasswordStore interface {
	GetPassword(ctx context.Context) error
}

type OTPStore interface {
	GetOTP(ctx context.Context, email string, otp string) (string, error)
	GetTokenPassword(ctx context.Context, email string, token string) (string, error)
}
