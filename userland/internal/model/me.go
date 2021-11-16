package model

import "github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model/entity"

type ProfileResponse struct {
	Profile entity.Profile `json:"profile"`
}

type UpdateProfileRequest struct {
	Fullname   string `json:"fullname"`
	DosageType string `json:"dosage_type"`
}

type UpdateProfileResponse struct {
	Success bool `json:"success"`
}

type EmailResponse struct {
	Email string `json:"email"`
}

type ChangeEmailRequest struct {
	Email string `json:"email"`
}

type ChangeEmailResponse struct {
	Success bool `json:"success"`
}

type ChangePasswordRequest struct {
	PasswordCurrent string `json:"password_current"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

type ChangePasswordResponse struct {
	Success bool `json:"success"`
}

type DeleteAccountRequest struct {
	Password string `json:"password"`
}

type DeleteAccountResponse struct {
	Success bool `json:"success"`
}

type SetPictureResponse struct {
	Success bool `json:"success"`
}

type DeletePictureResponse struct {
	Success bool `json:"success"`
}