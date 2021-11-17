package model

import "github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/jwt"

type ListSessionResponse struct {
	Sessions []SessionResponse `json:"sessions"`
}

type SessionResponse struct {
	JTI       string `json:"jti"`
	IP        string `json:"ip"`
	Client    string `json:"client"`
	CreatedAt string `json:"created_at"`
}

type EndCurrentResponse struct {
	Success bool `json:"success"`
}

type DeleteOtherResponse struct {
	Success bool `json:"success"`
}

type RefreshTokenResponse struct {
	RefreshToken jwt.Token `json:"refresh_token"`
}

type AccessTokenResponse struct {
	AccessToken jwt.Token `json:"access_token"`
}
