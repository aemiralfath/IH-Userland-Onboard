package entity

import "context"

type OTPStore interface {
	SetOTP(ctx context.Context, otpType, otpCode, otpValue string) error
	GetOTP(ctx context.Context, otpType, otpCode string) (string, error)
}
