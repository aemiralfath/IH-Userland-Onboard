package auth

import (
	"context"
	"fmt"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/email"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/security"
)

func (r *Repository) ForgotPassword(ctx context.Context, req model.ForgotPasswordRequest) error {
	exist, _, err := r.UserStore.CheckEmailExist(ctx, req.Email)
	if err != nil {
		return err
	}

	if !exist {
		return fmt.Errorf("User not found")
	}

	otpCode, err := security.GenerateOTP(6)
	if err != nil {
		return err
	}

	if err := r.OtpStore.SetOTP(ctx, "password", otpCode, req.Email); err != nil {
		return err
	}

	subject := "Userland Reset Password!"
	msg := fmt.Sprintf("Use this otp for reset your password: %s", otpCode)

	go email.NewEmail().SendEmail(req.Email, subject, msg)

	return nil
}
