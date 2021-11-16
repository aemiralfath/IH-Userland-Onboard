package auth

import (
	"context"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model/entity"
)

type authRepo interface {
	Register(ctx context.Context, req model.RegisterRequest) error
	Verification(ctx context.Context, req model.VerificationRequest) error
	ForgotPassword(ctx context.Context, req model.ForgotPasswordRequest) error
	ResetPassword(ctx context.Context, req model.ResetPasswordRequest) error
	Login(ctx context.Context, ip, clientName string, req model.LoginRequest) (string, entity.User, error)
}

type UsecaseAuth struct {
	auth authRepo
}

func New(repo authRepo) *UsecaseAuth {
	return &UsecaseAuth{
		auth: repo,
	}
}
