package auth

import (
	"context"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model"
)

type authUC interface {
	Register(ctx context.Context, body model.RegisterRequest) (model.RegisterResponse, error)
	Verification(ctx context.Context, body model.VerificationRequest) (model.VerificationResponse, error)
	ForgotPassword(ctx context.Context, body model.ForgotPasswordRequest) (model.ForgotPasswordResponse, error)
	ResetPassword(ctx context.Context, body model.ResetPasswordRequest) (model.ResetPasswordResponse, error)
	Login(ctx context.Context, ip, clientName string, body model.LoginRequest) (model.LoginResponse, error)
}

type DeliveryAuth struct {
	auth authUC
}

func NewAuth(auth authUC) *DeliveryAuth {
	return &DeliveryAuth{
		auth: auth,
	}
}
