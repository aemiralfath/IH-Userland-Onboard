package session

import (
	"context"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model"
)

type sessionUC interface {
	ListSession(ctx context.Context, userId string) (model.ListSessionResponse, error)
	EndCurrentSession(ctx context.Context, id string) (model.EndCurrentResponse, error)
	DeleteOtherSession(ctx context.Context, id string) (model.DeleteOtherResponse, error)
	RefreshToken(ctx context.Context, jti, id string) (model.RefreshTokenResponse, error)
	AccessToken(ctx context.Context, jti, id string) (model.AccessTokenResponse, error)
}

type DeliverySession struct {
	session sessionUC
}

func NewSession(session sessionUC) *DeliverySession {
	return &DeliverySession{
		session: session,
	}
}
