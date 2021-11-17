package entity

import "context"

type Password struct {
	ID        string `json:"id"`
	UserId    string `json:"user_id"`
	Password  string `json:"password" validate:"required"`
	CreatedAt string `json:"createdAt"`
}

type PasswordStore interface {
	AddNewPassword(ctx context.Context, userId, password string) error
	GetLastThreePassword(ctx context.Context, userId string) ([]string, error)
}
