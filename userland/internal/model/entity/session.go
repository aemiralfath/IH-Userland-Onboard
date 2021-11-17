package entity

import "context"

type Session struct {
	JTI       string `json:"JTI"`
	UserId    string `json:"userId"`
	IP        string `json:"ip"`
	Client    string `json:"client"`
	CreatedAt string `json:"createdAt"`
}

type SessionStore interface {
	AddNewSession(ctx context.Context, userId, ip, client string) (string, error)
	GetUserSession(ctx context.Context, userId string) ([]Session, error)
	EndSession(ctx context.Context, jti string) error
	DeleteOtherSession(ctx context.Context, jti string) error
}
