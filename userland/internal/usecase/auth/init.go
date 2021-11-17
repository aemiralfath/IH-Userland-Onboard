package auth

import (
	"context"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model/entity"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/kafka"
)

type authRepo interface {
	Register(ctx context.Context, req model.RegisterRequest) error
	Verification(ctx context.Context, req model.VerificationRequest) error
	ForgotPassword(ctx context.Context, req model.ForgotPasswordRequest) error
	ResetPassword(ctx context.Context, req model.ResetPasswordRequest) error
	Login(ctx context.Context, ip, clientName string, req model.LoginRequest) (string, entity.User, error)
}

type UsecaseAuth struct {
	auth  authRepo
	kafka kafka.Kafka
}

func New(repo authRepo) (*UsecaseAuth, error) {
	kafka, err := kafka.NewKafka()
	if err != nil {
		return &UsecaseAuth{}, err
	}

	return &UsecaseAuth{
		auth:  repo,
		kafka: kafka,
	}, nil
}
