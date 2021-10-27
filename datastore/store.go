package datastore

import (
	"context"
)

type User struct {
	ID        float64 `json:"id" sql:"id"`
	Email     string  `json:"email" validate:"required" sql:"email"`
	Password  string  `json:"password" validate:"required" sql:"password"`
	Verified  bool    `json:"verified" validate:"required" sql:"verified"`
	CreatedAt string  `json:"createdAt" sql:"created_at"`
	UpdatedAt string  `json:"updatedAt" sql:"updated_at"`
	DeletedAt string  `json:"deletedAt" sql:"deleted_at"`
}

type Profile struct {
	ID        float64 `json:"id" sql:"id"`
	UserId    float64 `json:"userId" sql:"user_id"`
	Fullname  string  `json:"fullname" validate:"required" sql:"fullname"`
	Location  string  `json:"location" sql:"location"`
	Bio       string  `json:"bio" sql:"bio"`
	Web       string  `json:"web" sql:"web"`
	Picture   string  `json:"picture" sql:"picture"`
	CreatedAt string  `json:"createdAt" sql:"created_at"`
	UpdatedAt string  `json:"updatedAt" sql:"updated_at"`
}

type Password struct {
	ID        float64 `json:"id" sql:"id"`
	UserId    float64 `json:"userId" sql:"user_id"`
	Password  string  `json:"password" validate:"required" sql:"password"`
	CreatedAt string  `json:"createdAt" sql:"created_at"`
}

type Session struct {
	JTI       string  `json:"JTI" sql:"jti"`
	UserId    float64 `json:"userId" sql:"user_id"`
	IsCurrent bool    `json:"isCurrent" validate:"required" sql:"is_current"`
}

type Client struct {
	ID   float64 `json:"id" sql:"id"`
	Name string  `json:"name" validate:"required" sql:"name"`
}

type Event struct {
	ID        float64 `json:"id" sql:"id"`
	SessionId string  `json:"JTI" sql:"jti"`
	ClientId  float64 `json:"clientId" sql:"client_id"`
	Event     string  `json:"event" sql:"event"`
	UserAgent string  `json:"userAgent" sql:"user_agent"`
	IP        string  `json:"ip" sql:"ip"`
	CreatedAt string  `json:"createdAt" sql:"created_at"`
	UpdatedAt string  `json:"updatedAt" sql:"updated_at"`
}

type SessionStore interface {
	GetUserSession(ctx context.Context, userId float64) (*Session, error)
	AddNewSession(ctx context.Context, session *Session) error
}

type ClientStore interface {
	GetClientByName(ctx context.Context, name string) (*Client, error)
	AddNewClient(ctx context.Context, name string) (*Client, error)
}

type EventStore interface {
	GetEventBySession(ctx context.Context, sessionId string) (*Event, error)
	AddNewEvent(ctx context.Context, sessionId string, clientId float64) error
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
	SafeDeleteUser(ctx context.Context, email string) error
}

type PasswordStore interface {
	GetLastThreePassword(ctx context.Context, userId float64) ([]string, error)
	AddNewPassword(ctx context.Context, password *Password, userId float64) error
}

type TokenStore interface {
	SetToken(ctx context.Context, tokenType, email, token string) error
	GetToken(ctx context.Context, tokenType, token string) (string, error)
}
