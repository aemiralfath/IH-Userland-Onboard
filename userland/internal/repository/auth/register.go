package auth

import (
	"context"
	"fmt"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/email"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/security"
)

func (r *Repository) Register(ctx context.Context, req model.RegisterRequest) error {
	exist, _, err := r.UserStore.CheckEmailExist(ctx, req.Email)
	if err != nil {
		return err
	}

	if exist {
		return fmt.Errorf("Email already exists")
	}

	hashPassword, err := security.HashPassword(req.Password)
	if err != nil {
		return err
	}

	profileExist, profile, err := r.ProfileStore.CheckNIKExist(ctx, req.Profile.NIK)
	if err != nil {
		return err
	}

	if !profileExist {
		profile, err = r.ProfileStore.AddNewProfile(ctx, req.Profile)
		if err != nil {
			return err
		}
	}

	userId, err := r.UserStore.AddNewUser(ctx, profile.ID, req.Email, hashPassword)
	if err != nil {
		return err
	}

	if err := r.PasswordStore.AddNewPassword(ctx, userId, hashPassword); err != nil {
		return err
	}

	otpCode, err := security.GenerateOTP(6)
	if err != nil {
		return err
	}

	otpValue := fmt.Sprintf("%s %s", userId, req.Email)
	if err := r.OtpStore.SetOTP(ctx, "user", otpCode, otpValue); err != nil {
		return err
	}

	subject := "Userland Email Verification!"
	msg := fmt.Sprintf("Use this otp for verify your email: %s", otpCode)

	go email.NewEmail().SendEmail(req.Email, subject, msg)

	return nil
}
