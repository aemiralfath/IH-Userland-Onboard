package auth

import (
	"context"
	"fmt"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model/entity"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/security"
)

func (r *Repository) Login(ctx context.Context, ip, clientName string, req model.LoginRequest) (string, entity.User, error) {
	exist, user, err := r.UserStore.CheckEmailExist(ctx, req.Email)
	if err != nil {
		return "", user, err
	}

	if !exist {
		return "", user, fmt.Errorf("User not found")
	}

	if err := security.ConfirmPassword(user.Password, req.Password); !err {
		return "", user, fmt.Errorf("Wrong Password")
	}

	jti, err := r.SessionStore.AddNewSession(ctx, user.ID, ip, clientName)
	if err != nil {
		return "", user, err
	}

	return jti, user, nil
}
