package datastore

import (
	"context"
)

type User struct {
	ID        float64 `json:"id"`
	Email     string  `json:"email" validate:"required"`
	Password  string  `json:"password" validate:"required"`
	Verified  bool    `json:"verified" validate:"required"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
	DeletedAt string  `json:"deletedAt"`
}

type Profile struct {
	ID        float64 `json:"id"`
	UserId    float64 `json:"userId"`
	Fullname  string  `json:"fullname" validate:"required"`
	Location  string  `json:"location"`
	Bio       string  `json:"bio"`
	Web       string  `json:"web"`
	Picture   string  `json:"picture"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
}

type Password struct {
	ID        float64 `json:"id"`
	UserId    float64 `json:"userId"`
	Password  string  `json:"password" validate:"required"`
	CreatedAt string  `json:"createdAt"`
}

type Session struct {
	JTI       string  `json:"JTI"`
	UserId    float64 `json:"userId"`
	ClientId  float64 `json:"clientId"`
	IsCurrent bool    `json:"isCurrent" validate:"required"`
	Event     string  `json:"event"`
	UserAgent string  `json:"userAgent"`
	IP        string  `json:"ip"`
	Client    Client  `json:"client"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
}

type Client struct {
	ID   float64 `json:"id"`
	Name string  `json:"name" validate:"required"`
}

type SessionStore interface {
	GetUserSession(ctx context.Context, userId float64) ([]Session, error)
	AddNewSession(ctx context.Context, session *Session, clientId float64) error
	EndSession(ctx context.Context, jti string) error
	DeleteOtherSession(ctx context.Context, jti string) error
}

type ClientStore interface {
	GetClientByName(ctx context.Context, name string) (*Client, error)
	AddNewClient(ctx context.Context, name string) (*Client, error)
}

type ProfileStore interface {
	GetProfile(ctx context.Context, userId float64) (*Profile, error)
	AddNewProfile(ctx context.Context, profile *Profile, userId float64) error
	UpdateProfile(ctx context.Context, profile *Profile, userId float64) error
	UpdatePicture(ctx context.Context, profile *Profile, userId float64) error
}

type UserStore interface {
	GetEmailByID(ctx context.Context, id float64) (string, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	AddNewUser(ctx context.Context, user *User) (float64, error)
	ChangePassword(ctx context.Context, user *User) error
	CheckUserEmailExist(ctx context.Context, email string) (*User, error)
	SoftDeleteUser(ctx context.Context, email string) error
}

type PasswordStore interface {
	GetLastThreePassword(ctx context.Context, userId float64) ([]string, error)
	AddNewPassword(ctx context.Context, password *Password, userId float64) error
}

type OTPStore interface {
	SetOTP(ctx context.Context, otpType, otpCode, otpValue string) error
	GetOTP(ctx context.Context, otpType, otpCode string) (string, error)
}
