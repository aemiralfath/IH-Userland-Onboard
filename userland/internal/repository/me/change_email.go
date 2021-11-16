package me

import (
	"context"
	"fmt"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/email"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/security"
)

func (r *Repository) ChangeEmail(ctx context.Context, userId string, req model.ChangeEmailRequest) error {
	exist, _, err := r.UserStore.CheckEmailExist(ctx, req.Email)
	if err != nil {
		return err
	}

	if exist {
		return fmt.Errorf("Email already used")
	}

	otpCode, err := security.GenerateOTP(6)
	if err != nil {
		return err
	}

	otpValue := fmt.Sprintf("%s %s", userId, req.Email)
	if err := r.OtpStore.SetOTP(ctx, "user", otpCode, otpValue); err != nil {
		return err
	}

	subject := "Userland Change Email Verification!"
	msg := fmt.Sprintf("Use this otp for verify your new email: %s", otpCode)

	go email.NewEmail().SendEmail(req.Email, subject, msg)

	return nil
}
