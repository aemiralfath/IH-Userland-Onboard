package session

import (
	"context"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model/entity"
)

type sessionRepo interface {
	ListSession(ctx context.Context, userId string) ([]entity.Session, error)
	EndCurrentSession(ctx context.Context, id string) error
	DeleteOtherSession(ctx context.Context, id string) error
}

type UsecaseSession struct {
	session sessionRepo
}

func New(repo sessionRepo) *UsecaseSession {
	return &UsecaseSession{
		session: repo,
	}
}
