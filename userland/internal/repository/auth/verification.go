package auth

import (
	"context"
	"fmt"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/email"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/security"
)

func (r *Repository) Verification(ctx context.Context, req model.VerificationRequest) error {
	exist, user, err := r.UserStore.CheckEmailExist(ctx, req.Email)
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

	otpValue := fmt.Sprintf("%s %s", user.ID, req.Email)
	if err := r.OtpStore.SetOTP(ctx, "user", otpCode, otpValue); err != nil {
		return err
	}

	subject := "Userland Email Verification!"
	msg := fmt.Sprintf("Use this otp for verify your email: %s", otpCode)

	go email.NewEmail().SendEmail(req.Email, subject, msg)

	return nil
}
