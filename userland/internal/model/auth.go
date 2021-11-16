package model

import (
	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model/entity"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/jwt"
)

type RegisterRequest struct {
	Email           string         `json:"email"`
	Password        string         `json:"password"`
	PasswordConfirm string         `json:"password_confirm"`
	Profile         entity.Profile `json:"profile"`
}

type RegisterResponse struct {
	Success bool `json:"success"`
}

type VerificationRequest struct {
	Email string `json:"recipient"`
	Type  string `json:"type"`
}

type VerificationResponse struct {
	Success bool `json:"success"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ForgotPasswordResponse struct {
	Success bool `json:"success"`
}

type ResetPasswordRequest struct {
	Token           string `json:"token"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

type ResetPasswordResponse struct {
	Success bool `json:"success"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	RequireTFA  bool      `json:"require_tfa"`
	AccessToken jwt.Token `json:"access_token"`
}
